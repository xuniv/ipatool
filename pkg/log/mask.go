package log

import (
	"io"
	"regexp"
)

var maskedFields = []*regexp.Regexp{
	regexp.MustCompile(`(?i)("password"\s*:\s*")[^"]*(")`),
	regexp.MustCompile(`(?i)("authCode"\s*:\s*")[^"]*(")`),
	regexp.MustCompile(`(?i)("passwordToken"\s*:\s*")[^"]*(")`),
	regexp.MustCompile(`(?i)("directoryServicesIdentifier"\s*:\s*")[^"]*(")`),
}

func maskSecrets(p []byte) []byte {
	masked := append([]byte{}, p...)

	for _, pattern := range maskedFields {
		masked = pattern.ReplaceAll(masked, []byte(`${1}<redacted>${2}`))
	}

	return masked
}

type maskingWriter struct {
	writer io.Writer
}

func NewMaskingWriter(writer io.Writer) io.Writer {
	return &maskingWriter{
		writer: writer,
	}
}

func (w *maskingWriter) Write(p []byte) (int, error) {
	return w.writer.Write(maskSecrets(p))
}
