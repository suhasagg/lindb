// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package deps

import (
	"context"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/coordinator/root"
	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/pkg/state"
)

// HTTPDeps represents http server handler's dependency.
type HTTPDeps struct {
	Ctx          context.Context
	Cfg          *config.Root
	QueryLimiter *concurrent.Limiter
	Repo         state.Repository
	RepoFactory  state.RepositoryFactory
	StateMgr     root.StateManager
}

func (deps *HTTPDeps) WithTimeout() (context.Context, context.CancelFunc) {
	// choose the shorter duration
	// TODO: need modify
	timeout := deps.Cfg.Coordinator.Timeout.Duration()
	if deps.Cfg.HTTP.ReadTimeout.Duration() < timeout {
		timeout = deps.Cfg.HTTP.ReadTimeout.Duration()
	}
	return context.WithTimeout(deps.Ctx, timeout)
}
