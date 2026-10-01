package conversion

import "fmt"

// Diagnostic identifies a conversion problem at a source field.
type Diagnostic struct {
	Platform string
	Path     string
	Code     string
	Message  string
}

func (d Diagnostic) Error() string {
	return fmt.Sprintf("convertago: %s %s: %s", d.Platform, d.Path, d.Message)
}

// Option configures a message conversion.
type Option func(*Options)

// Options holds the configuration shared by all converters.
type Options struct {
	Strict     bool
	Diagnostic func(Diagnostic)
}

// WithStrict rejects explicitly optional features that cannot be converted.
func WithStrict() Option { return func(o *Options) { o.Strict = true } }

// WithDiagnostics receives a diagnostic for each skipped optional feature.
func WithDiagnostics(fn func(Diagnostic)) Option {
	return func(o *Options) { o.Diagnostic = fn }
}

func Configure(options []Option) Options {
	var result Options
	for _, option := range options {
		if option != nil {
			option(&result)
		}
	}
	return result
}

func Error(platform, path, code, message string) error {
	return Diagnostic{Platform: platform, Path: path, Code: code, Message: message}
}
