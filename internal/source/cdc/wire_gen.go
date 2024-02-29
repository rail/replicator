// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package cdc

import (
	"github.com/cockroachdb/cdc-sink/internal/script"
	"github.com/cockroachdb/cdc-sink/internal/sequencer/besteffort"
	"github.com/cockroachdb/cdc-sink/internal/sequencer/immediate"
	"github.com/cockroachdb/cdc-sink/internal/sequencer/retire"
	script2 "github.com/cockroachdb/cdc-sink/internal/sequencer/script"
	"github.com/cockroachdb/cdc-sink/internal/sequencer/serial"
	"github.com/cockroachdb/cdc-sink/internal/sequencer/shingle"
	"github.com/cockroachdb/cdc-sink/internal/sequencer/switcher"
	"github.com/cockroachdb/cdc-sink/internal/sinktest/all"
	"github.com/cockroachdb/cdc-sink/internal/staging/checkpoint"
	"github.com/cockroachdb/cdc-sink/internal/staging/leases"
	"github.com/cockroachdb/cdc-sink/internal/target/apply"
	"github.com/cockroachdb/cdc-sink/internal/target/dlq"
	"github.com/cockroachdb/cdc-sink/internal/target/schemawatch"
	"github.com/cockroachdb/cdc-sink/internal/util/auth/trust"
	"github.com/cockroachdb/cdc-sink/internal/util/diag"
)

// Injectors from test_fixture.go:

func newTestFixture(fixture *all.Fixture, config *Config) (*testFixture, error) {
	sequencerConfig := ProvideSequencerConfig(config)
	baseFixture := fixture.Fixture
	context := baseFixture.Context
	stagingPool := baseFixture.StagingPool
	stagingSchema := baseFixture.StagingDB
	typesLeases, err := leases.ProvideLeases(context, stagingPool, stagingSchema)
	if err != nil {
		return nil, err
	}
	stagers := fixture.Stagers
	targetPool := baseFixture.TargetPool
	diagnostics := diag.New(context)
	watchers, err := schemawatch.ProvideFactory(context, targetPool, diagnostics)
	if err != nil {
		return nil, err
	}
	bestEffort := besteffort.ProvideBestEffort(sequencerConfig, typesLeases, stagingPool, stagers, targetPool, watchers)
	authenticator := trust.New()
	targetStatements := baseFixture.TargetCache
	configs := fixture.Configs
	dlqConfig := ProvideDLQConfig(config)
	dlQs := dlq.ProvideDLQs(dlqConfig, targetPool, watchers)
	acceptor, err := apply.ProvideAcceptor(context, targetStatements, configs, diagnostics, dlQs, targetPool, watchers)
	if err != nil {
		return nil, err
	}
	checkpoints, err := checkpoint.ProvideCheckpoints(context, stagingPool, stagingSchema)
	if err != nil {
		return nil, err
	}
	retireRetire := retire.ProvideRetire(sequencerConfig, stagingPool, stagers)
	immediateImmediate := &immediate.Immediate{}
	scriptConfig := ProvideScriptConfig(config)
	loader, err := script.ProvideLoader(context, configs, scriptConfig, diagnostics)
	if err != nil {
		return nil, err
	}
	sequencer := script2.ProvideSequencer(loader, targetPool, watchers)
	serialSerial := serial.ProvideSerial(sequencerConfig, typesLeases, stagers, stagingPool, targetPool)
	shingleShingle := shingle.ProvideShingle(sequencerConfig, stagers, stagingPool, targetPool)
	switcherSwitcher := switcher.ProvideSequencer(bestEffort, diagnostics, immediateImmediate, sequencer, serialSerial, shingleShingle, stagingPool, targetPool)
	targets, err := ProvideTargets(context, acceptor, config, checkpoints, retireRetire, stagingPool, switcherSwitcher, watchers)
	if err != nil {
		return nil, err
	}
	handler := &Handler{
		Authenticator: authenticator,
		Config:        config,
		TargetPool:    targetPool,
		Targets:       targets,
	}
	cdcTestFixture := &testFixture{
		Fixture:    fixture,
		BestEffort: bestEffort,
		Handler:    handler,
		Targets:    targets,
	}
	return cdcTestFixture, nil
}

// test_fixture.go:

type testFixture struct {
	*all.Fixture
	BestEffort *besteffort.BestEffort
	Handler    *Handler
	Targets    *Targets
}
