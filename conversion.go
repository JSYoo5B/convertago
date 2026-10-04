// Package convertago converts tagged Go structs into messenger-native messages.
// Each messenger has its own tag namespace and builder roles.
package convertago

import (
	"github.com/JSYoo5B/convertago/internal/conversion"
	"github.com/JSYoo5B/convertago/internal/validation"
)

// Option configures any of the message conversion and validation functions.
type Option = conversion.Option

// Diagnostic identifies a rule violation or a skipped optional feature.
// Conversion and validation errors are Diagnostic values and can be inspected with errors.As.
// Code is a rule ID such as "header.text.length" or a tag error code such as "invalid_tag".
type Diagnostic = conversion.Diagnostic

// Severity ranks how a diagnostic affects message delivery.
type Severity = validation.Severity

const (
	// Fatal diagnostics describe requests the platform rejects. They are always errors.
	Fatal = validation.Fatal
	// Warning diagnostics describe requests the platform accepts with a possibly unintended result.
	// They are errors with WithWarningAsError.
	Warning = validation.Warning
	// Advisory diagnostics describe UX guidance or client limits. They are never errors.
	Advisory = validation.Advisory
)

// WithWarningAsError turns Warning diagnostics into errors, including skipped optional features.
// Absent tags, '-' tags, nil fields, and omitempty remain intentional omissions.
func WithWarningAsError() Option { return conversion.WithWarningAsError() }

// WithDiagnostics receives Warning and Advisory diagnostics that do not stop processing, in order.
func WithDiagnostics(fn func(Diagnostic)) Option { return conversion.WithDiagnostics(fn) }
