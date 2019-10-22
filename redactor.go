package loggerus

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

type RedactingLogger interface {
	GetRedactor() *Redactor
}

type Redactor struct {
	output                 io.Writer
	redactions             []string
	valueRedactions        []string
	replacementString      string
	valueReplacementString string
}

func NewRedactor(output io.Writer) *Redactor {
	return &Redactor{
		output:                 output,
		redactions:             []string{},
		valueRedactions:        []string{},
		replacementString:      "*****",
		valueReplacementString: "[redacted]",
	}
}

func (r *Redactor) AddValueRedactions(valueRedactions []string) {
	r.valueRedactions = append(r.valueRedactions, valueRedactions...)
	r.valueRedactions = r.removeDuplicates(r.valueRedactions)
}

func (r *Redactor) AddRedactions(redactions []string) {
	var nonEmptyRedactions []string

	for _, redaction := range redactions {
		if redaction != "" {
			nonEmptyRedactions = append(nonEmptyRedactions, redaction)
		}
	}

	r.redactions = append(r.redactions, nonEmptyRedactions...)
	r.redactions = r.removeDuplicates(r.redactions)
}

func (r *Redactor) Write(p []byte) (n int, err error) {
	redactedPrint := r.redact(string(p[:]))
	return r.output.Write([]byte(redactedPrint))
}

func (r *Redactor) redact(input string) string {
	redacted := input

	// golang regex doesn't support lookarounds, so we will check things manually
	matchKeyWithSeparatorTemplate := `\\*[\'"]?(?i)%s\\*[\'"]?\s*[=:]\s*`
	matchValue := `\'[^\']*?\'|\"[^\"]*\"|\S*`

	// redact values of either strings of type `valueRedaction=[value]` or `valueRedaction: [value]`
	// w/wo single/double quotes
	for _, redactionField := range r.valueRedactions {
		matchKeyWithSeparator := fmt.Sprintf(matchKeyWithSeparatorTemplate, redactionField)
		re := regexp.MustCompile(fmt.Sprintf(`(%s)(%s)`, matchKeyWithSeparator, matchValue))
		redacted = re.ReplaceAllString(redacted, fmt.Sprintf(`$1%s`, r.valueReplacementString))
	}

	// replace the simple string redactions
	for _, redactionField := range r.redactions {
		redacted = strings.Replace(redacted, redactionField, r.replacementString, -1)
	}

	return redacted
}

func (r *Redactor) removeDuplicates(elements []string) []string {
	encountered := map[string]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	var result []string
	for key := range encountered {
		result = append(result, key)
	}
	return result
}
