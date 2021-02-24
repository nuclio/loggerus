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
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/logrusorgru/aurora/v3"
	"github.com/sirupsen/logrus"
)

type TextFormatter struct {
	maxVariableLen int
	enrichWhoField bool
	auroraInstance aurora.Aurora
}

func newTextFormatter(maxVariableLen int, enrichWhoField bool, color bool) (*TextFormatter, error) {
	return &TextFormatter{
		maxVariableLen: maxVariableLen,
		enrichWhoField: enrichWhoField,
		auroraInstance: aurora.NewAurora(color),
	}, nil
}

func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})

	// write date
	buffer.WriteString(f.auroraInstance.White(entry.Time.Format("02-01-06 15:04:05.000")).String()) // nolint: errcheck

	// write logger name
	if f.enrichWhoField {
		buffer.WriteString(" " + f.auroraInstance.Cyan(f.getFormattedWho(entry.Data)).String()) // nolint: errcheck
	}

	// write level
	buffer.WriteString(" " + f.getLevelOutput(entry.Level)) // nolint: errcheck

	// write message
	buffer.WriteString(" " + entry.Message) // nolint: errcheck

	// write fields
	buffer.WriteString(f.getFieldsOutput(entry.Data)) // nolint: errcheck

	// add newline
	buffer.WriteByte('\n') // nolint: errcheck

	return buffer.Bytes(), nil
}

func (f *TextFormatter) getLevelOutput(level logrus.Level) string {

	switch level {

	case logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel:
		return f.auroraInstance.Red("(E)").String()

	case logrus.WarnLevel:
		return f.auroraInstance.Yellow("(W)").String()

	case logrus.InfoLevel:
		return f.auroraInstance.Blue("(I)").String()

	case logrus.DebugLevel:
		return f.auroraInstance.Green("(D)").String()

	case logrus.TraceLevel:
		return f.auroraInstance.Green("(T)").String()
	}

	return f.auroraInstance.BrightRed("(?)").String()
}

func (f *TextFormatter) getFieldsOutput(fields logrus.Fields) string {
	maxVariableLen := f.maxVariableLen
	if maxVariableLen == 0 {
		maxVariableLen = math.MaxInt64
	}

	// remove context - it shouldn't be printed
	delete(fields, "ctx")

	singleLineKV := map[string]string{}
	blockKV := map[string]string{}

	for fieldKey, fieldValue := range fields {

		// if we're dealing with a struct, use json
		switch reflect.Indirect(reflect.ValueOf(fieldValue)).Kind() {
		case reflect.Slice, reflect.Map, reflect.Struct:
			fieldValueBytes, _ := json.Marshal(fieldValue)

			// if it's short - add to single line. otherwise to block
			if len(fieldValueBytes) <= maxVariableLen {
				singleLineKV[fieldKey] = string(fieldValueBytes)
			} else {
				blockBuffer := bytes.NewBuffer([]byte{})

				if err := json.Indent(blockBuffer, fieldValueBytes, "", "\t"); err != nil {
					blockBuffer.WriteString(fmt.Sprintf("Failed to encode: %s", err.Error())) // nolint: errcheck
				}

				blockKV[fieldKey] = blockBuffer.String()
			}

		case reflect.String:
			stringFieldValue := fmt.Sprintf("%s", fieldValue)

			// if there are newlines in output, add to block
			if strings.Contains(stringFieldValue, "\n") {
				blockKV[fieldKey] = stringFieldValue
			} else {
				singleLineKV[fieldKey] = fmt.Sprintf(`"%s"`, fieldValue)
			}

		default:
			singleLineKV[fieldKey] = fmt.Sprintf("%v", fieldValue)
		}
	}

	fieldsOutput := ""
	if len(singleLineKV) != 0 {
		fieldsOutput = f.auroraInstance.White(" :: ").String()
	}

	separator := f.auroraInstance.White(" || ").String()

	for singleLineKey, singleLineValue := range singleLineKV {
		fieldsOutput += fmt.Sprintf("%s=%s%s", f.auroraInstance.Blue(singleLineKey).String(), singleLineValue, separator)
	}

	// remove last ||
	fieldsOutput = strings.TrimSuffix(fieldsOutput, separator)

	if len(blockKV) != 0 {
		for blockKey, blockValue := range blockKV {
			fieldsOutput += fmt.Sprintf("\n* %s:\n", f.auroraInstance.Blue(blockKey).String())
			fieldsOutput += blockValue
			fieldsOutput += "\n"
		}
	}

	return fieldsOutput
}

func (f *TextFormatter) getFormattedWho(data logrus.Fields) string {
	who, ok := data["who"]
	if ok {
		whoStr := fmt.Sprintf("%20s", who)
		return fmt.Sprintf("%20s", whoStr[len(whoStr)-20:])
	}

	return fmt.Sprintf("%20s", "")
}
