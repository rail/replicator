// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package pglogical

import (
	"github.com/cockroachdb/cdc-sink/internal/script"
	"github.com/cockroachdb/cdc-sink/internal/sequencer/bypass"
	"github.com/cockroachdb/cdc-sink/internal/sequencer/chaos"
	script2 "github.com/cockroachdb/cdc-sink/internal/sequencer/script"
	"github.com/cockroachdb/cdc-sink/internal/sinkprod"
	"github.com/cockroachdb/cdc-sink/internal/staging/memo"
	"github.com/cockroachdb/cdc-sink/internal/target/apply"
	"github.com/cockroachdb/cdc-sink/internal/target/dlq"
	"github.com/cockroachdb/cdc-sink/internal/target/schemawatch"
	"github.com/cockroachdb/cdc-sink/internal/util/applycfg"
	"github.com/cockroachdb/cdc-sink/internal/util/diag"
	"github.com/cockroachdb/cdc-sink/internal/util/stopper"
)

// Injectors from injector.go:

// Start creates a PostgreSQL logical replication loop using the
// provided configuration.
func Start(context *stopper.Context, config *Config) (*PGLogical, error) {
	diagnostics := diag.New(context)
	configs, err := applycfg.ProvideConfigs(diagnostics)
	if err != nil {
		return nil, err
	}
	scriptConfig := &config.Script
	loader, err := script.ProvideLoader(configs, scriptConfig, diagnostics)
	if err != nil {
		return nil, err
	}
	eagerConfig := ProvideEagerConfig(config, loader)
	targetConfig := &eagerConfig.Target
	targetPool, err := sinkprod.ProvideTargetPool(context, targetConfig, diagnostics)
	if err != nil {
		return nil, err
	}
	targetStatements, err := sinkprod.ProvideStatementCache(context, targetConfig, targetPool, diagnostics)
	if err != nil {
		return nil, err
	}
	dlqConfig := &eagerConfig.DLQ
	watchers, err := schemawatch.ProvideFactory(context, targetPool, diagnostics)
	if err != nil {
		return nil, err
	}
	dlQs := dlq.ProvideDLQs(dlqConfig, targetPool, watchers)
	acceptor, err := apply.ProvideAcceptor(context, targetStatements, configs, diagnostics, dlQs, targetPool, watchers)
	if err != nil {
		return nil, err
	}
	sequencerConfig := &eagerConfig.Sequencer
	chaosChaos := &chaos.Chaos{
		Config: sequencerConfig,
	}
	bypassBypass := &bypass.Bypass{}
	stagingConfig := &eagerConfig.Staging
	stagingPool, err := sinkprod.ProvideStagingPool(context, stagingConfig, diagnostics, targetConfig)
	if err != nil {
		return nil, err
	}
	stagingSchema, err := sinkprod.ProvideStagingDB(stagingConfig)
	if err != nil {
		return nil, err
	}
	memoMemo, err := memo.ProvideMemo(context, stagingPool, stagingSchema)
	if err != nil {
		return nil, err
	}
	sequencer := script2.ProvideSequencer(loader, watchers)
	conn, err := ProvideConn(context, acceptor, chaosChaos, config, bypassBypass, memoMemo, sequencer, stagingPool, targetPool, watchers)
	if err != nil {
		return nil, err
	}
	pgLogical := &PGLogical{
		Conn:        conn,
		Diagnostics: diagnostics,
		Memo:        memoMemo,
	}
	return pgLogical, nil
}
