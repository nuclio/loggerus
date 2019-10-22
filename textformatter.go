package loggerus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"golang.org/x/tools/container/intsets"
)

type TextFormatter struct {
	maxVariableLen int
	enrichWhoField bool
}

func newTextFormatter(maxVariableLen int, enrichWhoField bool) (*TextFormatter, error) {
	color.NoColor = false
	return &TextFormatter{
		maxVariableLen: maxVariableLen,
		enrichWhoField: enrichWhoField,
	}, nil
}

func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})

	// write date
	buffer.WriteString(color.WhiteString(entry.Time.Format("02-01-06 15:04:05.000"))) // nolint: errcheck

	// write logger name
	if f.enrichWhoField {
		buffer.WriteString(" " + color.CyanString(f.getFormattedWho(entry.Data))) // nolint: errcheck
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
		return color.HiRedString("(E)")

	case logrus.WarnLevel:
		return color.HiYellowString("(W)")

	case logrus.InfoLevel:
		return color.HiBlueString("(I)")

	case logrus.DebugLevel:
		return color.GreenString("(D)")

	case logrus.TraceLevel:
		return color.GreenString("(T)")
	}

	return color.RedString("(?)")
}

func (f *TextFormatter) getFieldsOutput(fields logrus.Fields) string {
	maxVariableLen := f.maxVariableLen
	if maxVariableLen == 0 {
		maxVariableLen = intsets.MaxInt
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
		fieldsOutput = color.WhiteString(" :: ")
	}

	separator := color.WhiteString(" || ")

	for singleLineKey, singleLineValue := range singleLineKV {
		fieldsOutput += fmt.Sprintf("%s=%s%s", color.BlueString(singleLineKey), singleLineValue, separator)
	}

	// remove last ||
	fieldsOutput = strings.TrimSuffix(fieldsOutput, separator)

	if len(blockKV) != 0 {
		for blockKey, blockValue := range blockKV {
			fieldsOutput += fmt.Sprintf("\n* %s:\n", color.BlueString(blockKey))
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
