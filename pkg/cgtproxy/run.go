package cgtproxy

import (
	"context"

	. "github.com/black-desk/lib/go/errwrap"
)

func (c *CGTProxy) Run() (err error) {
	defer Wrap(&err, "running cgtproxy core")

	c.components, err = injectedComponents(c.cfg, c.log)
	if err != nil {
		return
	}

	c.pool.Go(c.waitStop)
	c.pool.Go(c.runWatcher)
	c.pool.Go(c.runRuleManager)

	return c.pool.Wait()
}

func (c *CGTProxy) Stop(err error) {
	c.stopCh <- err
}

func (c *CGTProxy) waitStop(ctx context.Context) (err error) {
	defer Wrap(&err)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err = <-c.stopCh:
		c.log.Debug("Stopped.")
		return err
	}
}

func (c *CGTProxy) runRuleManager(ctx context.Context) (err error) {
	defer c.log.Debugw("Rule manager exited.")

	c.log.Debugw("Start nft rule manager.")

	err = c.components.r.Run()
	if err != nil {
		return
	}

	return ctx.Err()
}

func (c *CGTProxy) runWatcher(ctx context.Context) (err error) {
	defer c.log.Debugw("Watcher exited.")

	c.log.Debugw("Start filesystem watcher.")

	err = c.components.m.Run(ctx)
	if err != nil {
		return
	}

	return ctx.Err()
}
