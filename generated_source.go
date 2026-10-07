package convertago

import (
	"encoding"

	"github.com/JSYoo5B/convertago/internal/conversion"
)

// SourceField is the source field representation used by generated accessors.
type SourceField = conversion.Field

// SourceValue holds a scalar, tagged object, or ordered list from generated code.
type SourceValue = conversion.Value

// SourceTag is parsed tag metadata emitted by the generator.
// Its Style and Fixed slices may reference shared, read-only metadata; copy them before editing.
type SourceTag = conversion.Tag

// SourceFixed is a slot value fixed by a tag, emitted by the generator.
type SourceFixed = conversion.Fixed

// SourceState tracks traversal for generated accessors and dynamic fields.
// Its zero value is ready to use for one conversion.
type SourceState = conversion.State

// SourceObject constructs a tagged object for a generated accessor.
func SourceObject(fields []SourceField) SourceValue { return conversion.Object(fields) }

// SourceMarshaled reads a field's explicitly implemented text representation.
func SourceMarshaled(input encoding.TextMarshaler) SourceValue { return conversion.Marshaled(input) }
