// Copyright 2024 The Cockroach Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

//go:build !(cgo && linux && (target_oracle || target_all))

package stdpool

import (
	"github.com/cockroachdb/replicator/internal/types"
	"github.com/cockroachdb/replicator/internal/util/stopper"
	"github.com/pkg/errors"
)

// OpenOracleAsTarget returns an unsupported error.
func OpenOracleAsTarget(
	ctx *stopper.Context, connectString string, options ...Option,
) (*types.TargetPool, error) {
	return nil, errors.New("this build does not support Oracle Database")
}
