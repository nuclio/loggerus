/*
Copyright 2021 The Nuclio Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package loggerus

import (
	"context"
	"testing"

	"github.com/nuclio/logger"
	"github.com/stretchr/testify/suite"
)

type loggerSuite struct {
	suite.Suite
	logger logger.Logger
}

func (suite *loggerSuite) SetupSuite() {
	var err error

	// initialize logger for test
	suite.logger, err = NewLoggerusForTests("test")
	suite.Require().NoError(err)
}

func (suite *loggerSuite) TestLog() {
	suite.logger.Debug("test")
	suite.logger.Warn("test")
	suite.logger.Error("test")
	suite.logger.Info("test")
}

func (suite *loggerSuite) TestLogWith() {
	suite.logger.DebugWith("test", "with", "something")
	suite.logger.WarnWith("test", "with", "something")
	suite.logger.ErrorWith("test", "with", "something")
	suite.logger.InfoWith("test", "with", "something")
}

func (suite *loggerSuite) TestLogWithCtx() {
	ctx := context.WithValue(context.TODO(), "RequestID", "123") // nolint
	suite.logger.DebugWithCtx(ctx, "test", "with", "something")
	suite.logger.WarnWithCtx(ctx, "test", "with", "something")
	suite.logger.ErrorWithCtx(ctx, "test", "with", "something")
	suite.logger.InfoWithCtx(ctx, "test", "with", "something")
}

func TestLoggerTestSuite(t *testing.T) {
	suite.Run(t, new(loggerSuite))
}
