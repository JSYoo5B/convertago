// Package convertago converts tagged Go structs into messenger-native messages.
// Each messenger has its own tag namespace and builder roles.
package convertago

import "github.com/JSYoo5B/convertago/internal/conversion"

// Option configures any of the message conversion functions.
type Option = conversion.Option

// Diagnostic identifies an error or an explicitly skipped optional feature.
// Conversion errors can be inspected with errors.As.
type Diagnostic = conversion.Diagnostic

// WithStrict rejects recognized unavailable features even when tagged optional.
// Absent tags, '-' tags, nil fields, and omitempty remain intentional omissions.
func WithStrict() Option { return conversion.WithStrict() }

// WithDiagnostics receives diagnostics for skipped optional features in source order.
func WithDiagnostics(fn func(Diagnostic)) Option { return conversion.WithDiagnostics(fn) }
