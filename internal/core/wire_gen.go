// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package core

import (
	"github.com/black-desk/deepin-network-proxy-manager/internal/core/monitor"
	"github.com/black-desk/deepin-network-proxy-manager/internal/core/repeater"
	"github.com/black-desk/deepin-network-proxy-manager/internal/core/rulemanager"
	"github.com/black-desk/deepin-network-proxy-manager/internal/core/watcher"
)

// Injectors from wire.go:

func injectedMonitor(core *Core) (*monitor.Monitor, error) {
	context := provideContext(core)
	v := provideOutputChan()
	config, err := provideConfig(core)
	if err != nil {
		return nil, err
	}
	cgroupRoot := provideCgroupRoot(config)
	watcher, err := provideWatcher(context, cgroupRoot)
	if err != nil {
		return nil, err
	}
	monitorMonitor, err := provideMonitor(context, v, watcher)
	if err != nil {
		return nil, err
	}
	return monitorMonitor, nil
}

func injectedRuleManager(core *Core) (*rulemanager.RuleManager, error) {
	conn, err := provideNftConn()
	if err != nil {
		return nil, err
	}
	config, err := provideConfig(core)
	if err != nil {
		return nil, err
	}
	rerouteMark := provideRerouteMark(config)
	cgroupRoot := provideCgroupRoot(config)
	table, err := provideTable(conn, rerouteMark, cgroupRoot)
	if err != nil {
		return nil, err
	}
	v := provideInputChan()
	ruleManager, err := provideRuleManager(table, config, v)
	if err != nil {
		return nil, err
	}
	return ruleManager, nil
}

func injectedRepeater(core *Core) (*repeater.Repeater, error) {
	repeaterRepeater, err := provideRepeater()
	if err != nil {
		return nil, err
	}
	return repeaterRepeater, nil
}

func injectedWatcher(core *Core) (*watcher.Watcher, error) {
	context := provideContext(core)
	config, err := provideConfig(core)
	if err != nil {
		return nil, err
	}
	cgroupRoot := provideCgroupRoot(config)
	watcherWatcher, err := provideWatcher(context, cgroupRoot)
	if err != nil {
		return nil, err
	}
	return watcherWatcher, nil
}
