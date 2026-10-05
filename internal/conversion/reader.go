package conversion

import (
	"math"
	"strconv"
)

// Reader parses the scalar slots of a node while retaining the first source error.
// It does not apply native rules; converters check the values they build instead.
type Reader struct {
	Platform string
	Node     Node
	Err      error
}

func (r *Reader) path(slot string) string {
	for _, p := range r.Node.Parts {
		if p.Slot == slot {
			return p.Path
		}
	}
	return r.Node.Path
}

func (r *Reader) invalid(slot, message string) {
	if r.Err == nil {
		r.Err = Error(r.Platform, r.path(slot), "invalid_value", message)
	}
}

// String returns a slot's text.
func (r *Reader) String(slot string) string { return r.Node.Text(slot) }

// ParseInt parses a slot as an integer. An absent slot is zero.
func (r *Reader) ParseInt(slot string) int {
	if !r.Node.Has(slot) {
		return 0
	}
	value, err := strconv.Atoi(r.Node.Text(slot))
	if err != nil {
		r.invalid(slot, slot+" requires an integer")
	}
	return value
}

// ParseFloat parses a slot as a finite number, which JSON can represent. An absent slot is zero.
func (r *Reader) ParseFloat(slot string) float64 {
	if !r.Node.Has(slot) {
		return 0
	}
	value, err := strconv.ParseFloat(r.Node.Text(slot), 64)
	if err != nil || math.IsNaN(value) || math.IsInf(value, 0) {
		r.invalid(slot, slot+" requires a finite number")
	}
	return value
}

// Bool parses a slot as a boolean. An absent slot is false.
func (r *Reader) Bool(slot string) bool {
	if !r.Node.Has(slot) {
		return false
	}
	value, err := strconv.ParseBool(r.Node.Text(slot))
	if err != nil {
		r.invalid(slot, slot+" requires a boolean")
	}
	return value
}
