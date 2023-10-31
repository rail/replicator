// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package all

import (
	"github.com/cockroachdb/cdc-sink/internal/sinktest/base"
	"github.com/cockroachdb/cdc-sink/internal/staging/memo"
	"github.com/cockroachdb/cdc-sink/internal/staging/stage"
	"github.com/cockroachdb/cdc-sink/internal/staging/version"
	"github.com/cockroachdb/cdc-sink/internal/target/apply"
	"github.com/cockroachdb/cdc-sink/internal/target/dlq"
	"github.com/cockroachdb/cdc-sink/internal/target/schemawatch"
	"github.com/cockroachdb/cdc-sink/internal/util/applycfg"
	"github.com/cockroachdb/cdc-sink/internal/util/diag"
)

// Injectors from injector.go:

// NewFixture constructs a self-contained test fixture for all services
// in the target sub-packages.
func NewFixture() (*Fixture, func(), error) {
	context, cleanup := base.ProvideContext()
	diagnostics, cleanup2 := diag.New(context)
	sourcePool, cleanup3, err := base.ProvideSourcePool(context, diagnostics)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	sourceSchema, cleanup4, err := base.ProvideSourceSchema(context, sourcePool)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	stagingPool, cleanup5, err := base.ProvideStagingPool(context)
	if err != nil {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	stagingSchema, cleanup6, err := base.ProvideStagingSchema(context, stagingPool)
	if err != nil {
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	targetPool, cleanup7, err := base.ProvideTargetPool(context, sourcePool, diagnostics)
	if err != nil {
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	targetStatements, cleanup8 := base.ProvideTargetStatements(targetPool)
	targetSchema, cleanup9, err := base.ProvideTargetSchema(context, diagnostics, targetPool, targetStatements)
	if err != nil {
		cleanup8()
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	fixture := &base.Fixture{
		Context:      context,
		SourcePool:   sourcePool,
		SourceSchema: sourceSchema,
		StagingPool:  stagingPool,
		StagingDB:    stagingSchema,
		TargetCache:  targetStatements,
		TargetPool:   targetPool,
		TargetSchema: targetSchema,
	}
	configs, err := applycfg.ProvideConfigs(diagnostics)
	if err != nil {
		cleanup9()
		cleanup8()
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	config, err := ProvideDLQConfig()
	if err != nil {
		cleanup9()
		cleanup8()
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	watchers, cleanup10, err := schemawatch.ProvideFactory(targetPool, diagnostics)
	if err != nil {
		cleanup9()
		cleanup8()
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	dlQs := dlq.ProvideDLQs(config, targetPool, watchers)
	appliers, cleanup11, err := apply.ProvideFactory(targetStatements, configs, diagnostics, dlQs, targetPool, watchers)
	if err != nil {
		cleanup10()
		cleanup9()
		cleanup8()
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	memoMemo, err := memo.ProvideMemo(context, stagingPool, stagingSchema)
	if err != nil {
		cleanup11()
		cleanup10()
		cleanup9()
		cleanup8()
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	stagers := stage.ProvideFactory(stagingPool, stagingSchema, context)
	checker := version.ProvideChecker(stagingPool, memoMemo)
	watcher, err := ProvideWatcher(context, targetSchema, watchers)
	if err != nil {
		cleanup11()
		cleanup10()
		cleanup9()
		cleanup8()
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	allFixture := &Fixture{
		Fixture:        fixture,
		Appliers:       appliers,
		Configs:        configs,
		Diagnostics:    diagnostics,
		DLQConfig:      config,
		DLQs:           dlQs,
		Memo:           memoMemo,
		Stagers:        stagers,
		VersionChecker: checker,
		Watchers:       watchers,
		Watcher:        watcher,
	}
	return allFixture, func() {
		cleanup11()
		cleanup10()
		cleanup9()
		cleanup8()
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}
