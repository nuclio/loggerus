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
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type redactorSuite struct {
	suite.Suite
	redactor *Redactor
}

func (suite *redactorSuite) TestKeyValueTypeValueRedactions() {
	buf := new(bytes.Buffer)

	// prepare redactor
	valueRedactions := []string{"artifactVersionManifestContents", "systemConfigContents"}
	suite.redactor = NewRedactor(buf)
	suite.redactor.AddValueRedactions(valueRedactions)

	// read file into byte string
	unredactedCommand, err := ioutil.ReadFile("test/key_value.txt")
	suite.Assert().Nil(err)

	// write it using the redactor write function
	bytesWritten, err := suite.redactor.Write(unredactedCommand)
	suite.Assert().Nil(err)
	suite.Assert().True(bytesWritten > 0)

	// verify that command was indeed redacted
	redactedCommand := buf.String()
	suite.Assert().True(strings.Contains(redactedCommand, "artifactVersionManifestContents=[redacted]"))
	suite.Assert().True(strings.Contains(redactedCommand, "systemConfigContents=[redacted]"))
}

func (suite *redactorSuite) TestDictTypeValueRedactions() {
	buf := new(bytes.Buffer)

	// prepare redactor
	valueRedactions := []string{"java_key_store"}
	suite.redactor = NewRedactor(buf)
	suite.redactor.AddValueRedactions(valueRedactions)

	// read file into byte string
	unredactedCommand, err := ioutil.ReadFile("test/dict.txt")
	suite.Assert().Nil(err)

	// write it using the redactor write function
	bytesWritten, err := suite.redactor.Write(unredactedCommand)
	suite.Assert().Nil(err)
	suite.Assert().True(bytesWritten > 0)

	// verify that command was indeed redacted
	redactedCommand := buf.String()
	suite.Assert().True(strings.Contains(redactedCommand, `"java_key_store":[redacted]`))
}

func (suite *redactorSuite) TestRegularRedactions() {
	buf := new(bytes.Buffer)

	// prepare redactor
	redactions := []string{"password"}
	suite.redactor = NewRedactor(buf)
	suite.redactor.AddRedactions(redactions)

	// push some string to writer
	unredactedCommand := "{asdhaksjd:\\ password \\ \n}"
	bytesWritten, err := suite.redactor.Write([]byte(unredactedCommand))
	suite.Assert().Nil(err)
	suite.Assert().True(bytesWritten > 0)

	// verify that command was indeed redacted
	redactedCommand := buf.String()
	suite.Assert().True(strings.Contains(redactedCommand, "{asdhaksjd:\\ ***** \\ \n}"))
}

func TestRedactorTestSuite(t *testing.T) {
	suite.Run(t, new(redactorSuite))
}
