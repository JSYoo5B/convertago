package conversion

import (
	"fmt"

	"github.com/JSYoo5B/convertago/internal/validation"
)

// Severity ranks how a diagnostic affects message delivery.
type Severity = validation.Severity

// Diagnostic identifies a conversion or validation problem at a source or native path.
type Diagnostic struct {
	Platform string
	Path     string
	Code     string
	Message  string
	Severity Severity
}

// Error formats the diagnostic with its platform, path, message, severity, and code.
func (d Diagnostic) Error() string {
	return fmt.Sprintf("convertago: %s %s: %s [%s %s]", d.Platform, d.Path, d.Message, d.Severity, d.Code)
}

// Option configures a message conversion.
type Option func(*Options)

// Options holds the configuration shared by all converters.
type Options struct {
	WarningAsError bool
	Diagnostic     func(Diagnostic)
}

// WithWarningAsError turns Warning diagnostics into errors.
func WithWarningAsError() Option { return func(o *Options) { o.WarningAsError = true } }

// WithDiagnostics receives each Warning and Advisory diagnostic that does not stop processing.
func WithDiagnostics(fn func(Diagnostic)) Option {
	return func(o *Options) { o.Diagnostic = fn }
}

// Configure applies options in order, ignoring nil options.
func Configure(options []Option) Options {
	var result Options
	for _, option := range options {
		if option != nil {
			option(&result)
		}
	}
	return result
}

// Handle returns d as an error when its severity stops processing and otherwise reports it.
func (o Options) Handle(d Diagnostic) error {
	if d.Severity == validation.Fatal || (d.Severity == validation.Warning && o.WarningAsError) {
		return d
	}
	if o.Diagnostic != nil {
		o.Diagnostic(d)
	}
	return nil
}

// Error returns a Fatal diagnostic.
func Error(platform, path, code, message string) error {
	return Diagnostic{Platform: platform, Path: path, Code: code, Message: message, Severity: validation.Fatal}
}
