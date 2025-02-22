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

package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/app/root/deps"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/coordinator/root"
	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/internal/mock"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/sql"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
)

func TestExecuteAPI_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)
	repoFct := state.NewMockRepositoryFactory(ctrl)
	stateMgr := root.NewMockStateManager(ctrl)

	cfg := `{\"config\":{\"namespace\":\"test\",\"timeout\":10,\"dialTimeout\":10,`
	cfg += `\"leaseTTL\":10,\"endpoints\":[\"http://localhost:2379\"]}}`
	api := NewExecuteAPI(&deps.HTTPDeps{
		Ctx:         context.Background(),
		Repo:        repo,
		RepoFactory: repoFct,
		StateMgr:    stateMgr,
		Cfg: &config.Root{
			HTTP: config.HTTP{ReadTimeout: ltoml.Duration(time.Second * 10)},
		},
		QueryLimiter: concurrent.NewLimiter(
			context.TODO(),
			2,
			time.Second*5,
			metrics.NewLimitStatistics("exec", linmetric.RootRegistry),
		),
	})
	r := gin.New()
	api.Register(r)

	cases := []struct {
		name    string
		reqBody string
		prepare func()
		assert  func(resp *httptest.ResponseRecorder)
	}{
		{
			name: "param invalid",
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "parse sql failure",
			reqBody: `{"sql":"show a"}`,
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "unknown statement type",
			reqBody: `{"sql":"show brokers"}`,
			prepare: func() {
				sqlParseFn = func(sql string) (stmt stmtpkg.Statement, err error) {
					return &stmtpkg.State{}, nil
				}
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "create broker json err",
			reqBody: `{"sql":"create broker ` + cfg + `"}`,
			prepare: func() {
				sqlParseFn = func(sql string) (stmt stmtpkg.Statement, err error) {
					return &stmtpkg.Broker{Type: stmtpkg.BrokerOpCreate, Value: "xx"}, nil
				}
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "create broker, config validate failure",
			reqBody: `{"sql":"create broker ` + cfg + `"}`,
			prepare: func() {
				sqlParseFn = func(sql string) (stmt stmtpkg.Statement, err error) {
					return &stmtpkg.Broker{Type: stmtpkg.BrokerOpCreate, Value: `{"config":{}}`}, nil
				}
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "create broker successfully, broker not exist",
			reqBody: `{"sql":"create broker ` + cfg + `"}`,
			prepare: func() {
				repoFct.EXPECT().CreateBrokerRepo(gomock.Any()).Return(repo, nil)
				repo.EXPECT().Close().Return(nil)
				repo.EXPECT().PutWithTX(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, _ []byte, check func([]byte) error) (bool, error) {
						if err := check([]byte{1, 2, 3}); err != nil {
							return false, err
						}
						return true, nil
					})
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "create broker successfully, broker exist",
			reqBody: `{"sql":"create broker ` + cfg + `"}`,
			prepare: func() {
				repoFct.EXPECT().CreateBrokerRepo(gomock.Any()).Return(repo, nil)
				repo.EXPECT().Close().Return(nil)
				repo.EXPECT().PutWithTX(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, _ []byte, check func([]byte) error) (bool, error) {
						cfg1 := strings.ReplaceAll(cfg, `\"`, `"`)
						data := []byte(cfg1)
						broker := &config.BrokerCluster{}
						err := encoding.JSONUnmarshal(data, broker)
						assert.NoError(t, err)
						data = encoding.JSONMarshal(broker)
						if err := check(data); err != nil {
							return false, err
						}
						return true, nil
					})
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "create broker failure with err",
			reqBody: `{"sql":"create broker ` + cfg + `"}`,
			prepare: func() {
				repoFct.EXPECT().CreateBrokerRepo(gomock.Any()).Return(repo, nil)
				repo.EXPECT().Close().Return(nil)
				repo.EXPECT().PutWithTX(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false, fmt.Errorf("err"))
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "create broker failure",
			reqBody: `{"sql":"create broker ` + cfg + `"}`,
			prepare: func() {
				repoFct.EXPECT().CreateBrokerRepo(gomock.Any()).Return(repo, nil)
				repo.EXPECT().Close().Return(nil)
				repo.EXPECT().PutWithTX(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "create broker repo failure",
			reqBody: `{"sql":"create broker ` + cfg + `"}`,
			prepare: func() {
				repoFct.EXPECT().CreateBrokerRepo(gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "create broker, close repo failure",
			reqBody: `{"sql":"create broker ` + cfg + `"}`,
			prepare: func() {
				repoFct.EXPECT().CreateBrokerRepo(gomock.Any()).Return(repo, nil)
				repo.EXPECT().Close().Return(fmt.Errorf("err"))
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "unknown broker op type",
			reqBody: `{"sql":"show brokers"}`,
			prepare: func() {
				sqlParseFn = func(sql string) (stmt stmtpkg.Statement, err error) {
					return &stmtpkg.Broker{Type: stmtpkg.BrokerOpUnknown}, nil
				}
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name:    "show brokers, get brokers failure",
			reqBody: `{"sql":"show brokers"}`,
			prepare: func() {
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "show brokers, list broker successfully, but unmarshal failure",
			reqBody: `{"sql":"show brokers"}`,
			prepare: func() {
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(
					[]state.KeyValue{{Key: "", Value: []byte("[]")}}, nil)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "show brokers successfully",
			reqBody: `{"sql":"show brokers"}`,
			prepare: func() {
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(
					[]state.KeyValue{{Key: "", Value: []byte(`{ "config": {"namespace":"xxx"}}`)}}, nil)
				stateMgr.EXPECT().GetBrokerState("xxx").Return(models.BrokerState{}, true)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "show brokers successfully,but state not found",
			reqBody: `{"sql":"show brokers"}`,
			prepare: func() {
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(
					[]state.KeyValue{{Key: "", Value: []byte(`{ "config": {"namespace":"xxx"}}`)}}, nil)
				stateMgr.EXPECT().GetBrokerState("xxx").Return(models.BrokerState{}, false)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				sqlParseFn = sql.Parse
			}()
			if tt.prepare != nil {
				tt.prepare()
			}
			resp := mock.DoRequest(t, r, http.MethodPut, ExecutePath, tt.reqBody)
			if tt.assert != nil {
				tt.assert(resp)
			}
		})
	}
}
