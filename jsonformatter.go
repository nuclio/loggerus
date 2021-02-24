package loggerus

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
)

const jsonDefaultTimestampFormat = "2006-01-02T15:04:05.000000"

type JSONFormatter struct {

	// timestampFormat sets the format used for marshaling timestamps.
	timestampFormat string

	// Log time zone
	TimeZone string
}

func newJSONFormatter(timestampFormat string, timeZone string) (*JSONFormatter, error) {
	return &JSONFormatter{
		timestampFormat: timestampFormat,
		TimeZone:        timeZone,
	}, nil
}

func (f *JSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields, len(entry.Data)+6)

	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:

			// Otherwise errors are ignored by `encoding/json`
			// https://github.com/Sirupsen/logrus/issues/137
			data[k] = v.Error()
		case []byte:
			data[k] = string(v)

		default:
			switch reflect.Indirect(reflect.ValueOf(v)).Kind() {
			case reflect.Slice, reflect.Map, reflect.Struct:
				fieldValueBytes, _ := json.Marshal(v)
				data[k] = string(fieldValueBytes)
			default:
				data[k] = v
			}
		}
	}

	timestampFormat := f.timestampFormat

	if timestampFormat == "" {
		timestampFormat = jsonDefaultTimestampFormat
	}

	// "when": "2016-06-19T09:56:29.043641"
	switch f.TimeZone {
	case "utc":
		data["when"] = entry.Time.UTC().Format(timestampFormat)

	default:
		data["when"] = entry.Time.Format(timestampFormat)

	}

	// "who": logger name
	data["who"] = entry.Data["who"]

	// "severity": log lvl
	data["severity"] = strings.ToUpper(entry.Level.String())

	// "what": message
	data["what"] = entry.Message

	// "more": json string
	data["more"] = buildMoreValue(&data)

	// extract context as first-class citizen
	ctx, ok := entry.Data["ctx"]
	if !ok {
		ctx = ""
	}

	// "ctx": some-uuid
	data["ctx"] = ctx

	serialized, err := json.Marshal(data)

	if err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON, %v", err)
	}

	// we append the rune (byte) '\n' rather than the string "\n"
	return append(serialized, '\n'), nil
}

// Build data["more"] value
func buildMoreValue(data *logrus.Fields) map[string]string {
	additionalData := make(map[string]string)

	for key, value := range *data {
		switch key {
		case "when":
		case "who":
		case "severity":
		case "what":
		case "more":
		case "ctx":
			// don't include these inside the more value
		default:

			formattedValue := convertValueToString(value)
			additionalData[key] = formattedValue

			// The key was copied to additional_data (No need for duplication)
			delete(*data, key)
		}
	}

	return additionalData
}

// Convert the given value to string
func convertValueToString(value interface{}) string {
	switch value := value.(type) {
	case string:
		return value
	case error:

		//return error message
		return value.Error()
	default:
		return fmt.Sprintf("%v", value)
	}
}
