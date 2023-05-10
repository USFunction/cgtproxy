package monitor_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	. "github.com/black-desk/deepin-network-proxy-manager/internal/core/monitor"
	"github.com/black-desk/deepin-network-proxy-manager/internal/inject"
	. "github.com/black-desk/deepin-network-proxy-manager/internal/test/ginkgo-helper"
	"github.com/fsnotify/fsnotify"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sourcegraph/conc/pool"
)

var _ = Describe("Cgroup monitor create with fake fsnotify.Watcher", func() {
	var (
		watcher         *fsnotify.Watcher
		cgroupEventChan chan *CgroupEvent
		ctx             context.Context
		monitor         *Monitor
	)

	BeforeEach(func() {
		watcher = &fsnotify.Watcher{
			Events: make(chan fsnotify.Event),
			Errors: make(chan error),
		}

		cgroupEventChan = make(chan *CgroupEvent)

		var cgroupEventChanIn chan<- *CgroupEvent
		cgroupEventChanIn = cgroupEventChan

		ctx = context.Background()

		var err error

		container := inject.New()
		if err = container.Register(watcher); err != nil {
			Fail(fmt.Sprint(err.Error()))
		}
		if err = container.Register(cgroupEventChanIn); err != nil {
			Fail(fmt.Sprint(err.Error()))
		}
		if err = container.RegisterI(&ctx); err != nil {
			Fail(fmt.Sprint(err.Error()))
		}

		if monitor, err = New(container); err != nil {
			Fail(fmt.Sprintf("Failed to create monitor with fake fsnotify.Watcher:\n%s", err.Error()))
		}
	})

	ContextTable("receive %s", func(
		resultMsg string,
		events []fsnotify.Event, errs []error,
		expectResult []*CgroupEvent, expectErr error,
	) {
		var p *pool.ErrorPool

		BeforeEach(func() {
			p = new(pool.ErrorPool)

			p.Go(func() error {
				for i := range events {
					watcher.Events <- events[i]
				}
				close(watcher.Events)
				return nil
			})

			p.Go(func() error {
				// NOTE(black_desk): Errors from fsnotify is ignored for now.
				for i := range errs {
					watcher.Errors <- errs[i]
				}
				close(watcher.Errors)
				return nil
			})

			p.Go(monitor.Run)
		})

		AfterEach(func() {
			result := errors.Is(ctx.Err(), nil)
			Expect(result).To(BeTrue())
		})

		It(fmt.Sprintf("should %s", resultMsg), func() {
			var cgroupEvents []*CgroupEvent
			for cgroupEvent := range cgroupEventChan {
				cgroupEvents = append(cgroupEvents, cgroupEvent)
			}

			Expect(len(expectResult)).To(Equal(len(cgroupEvents)))

			for i := range cgroupEvents {
				Expect(*cgroupEvents[i]).To(Equal(*expectResult[i]))
			}

			err := p.Wait()
			if expectErr == nil {
				Expect(err).To(Succeed())
			} else {
				Expect(err).To(MatchError(expectErr))
			}
		})
	},
		ContextTableEntry(
			"send a `New` event, and exit with no error",
			[]fsnotify.Event{{
				Name: "/test/path/1",
				Op:   fsnotify.Create,
			}},
			[]error{},
			[]*CgroupEvent{{
				Path:      "/test/path/1",
				EventType: CgroupEventTypeNew,
			}},
			nil,
		).WithFmt("a fsnotify.Event with fsnotify.Create"),
		ContextTableEntry(
			"send a `Delete` event, and exit with no error",
			[]fsnotify.Event{{
				Name: "/test/path/2",
				Op:   fsnotify.Remove,
			}},
			[]error{},
			[]*CgroupEvent{{
				Path:      "/test/path/2",
				EventType: CgroupEventTypeDelete,
			}},
			nil,
		).WithFmt("a fsnotify.Event with fsnotify.Delete"),
		ContextTableEntry(
			"send nothing, and exit with no error",
			[]fsnotify.Event{},
			[]error{},
			[]*CgroupEvent{},
			nil,
		).WithFmt("nothing"),
		ContextTableEntry(
			"send nothing, and exit with error",
			[]fsnotify.Event{{
				Name: "/test/path/3",
				Op:   fsnotify.Op(0),
			}},
			[]error{},
			[]*CgroupEvent{},
			ErrUnexpectFsEventType,
		).WithFmt("invalid fsnotify.Event"),
	)
})

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Configuration Suite")
}
