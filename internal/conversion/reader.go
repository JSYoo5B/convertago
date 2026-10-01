package conversion

import (
	"math"
	"net/url"
	"strconv"
)

// Reader validates native scalar properties while retaining the first source error.
type Reader struct {
	Platform string
	Node     Node
	Err      error
}

func (r *Reader) Fail(code, message string) {
	if r.Err == nil {
		r.Err = Error(r.Platform, r.Node.Path, code, message)
	}
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
func (r *Reader) Text(slot string, min, max int) string {
	text := r.Node.Text(slot)
	if r.Node.Has(slot) {
		if err := ValidateText(r.Platform, r.path(slot), text, min, max); err != nil && r.Err == nil {
			r.Err = err
		}
	}
	return text
}
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
func (r *Reader) Int(slot string, min, max int) int {
	if !r.Node.Has(slot) {
		return 0
	}
	value, err := strconv.Atoi(r.Node.Text(slot))
	if err != nil || value < min || (max != 0 && value > max) {
		r.invalid(slot, slot+" is outside its integer range")
	}
	return value
}
func (r *Reader) Float(slot string, min, max float64) float64 {
	if !r.Node.Has(slot) {
		return 0
	}
	value, err := strconv.ParseFloat(r.Node.Text(slot), 64)
	if err != nil || math.IsNaN(value) || math.IsInf(value, 0) || value < min || (max != 0 && value > max) {
		r.invalid(slot, slot+" is outside its numeric range")
	}
	return value
}
func (r *Reader) Enum(slot string, values ...string) string {
	value := r.Node.Text(slot)
	if r.Node.Has(slot) && !contains(values, value) {
		r.invalid(slot, "unknown "+slot+" value")
	}
	return value
}
func (r *Reader) URL(slot string, httpsOnly bool) string {
	value := r.Node.Text(slot)
	if r.Node.Has(slot) {
		if err := ValidateURL(r.Platform, r.path(slot), value, httpsOnly); err != nil && r.Err == nil {
			r.Err = err
		}
	}
	return value
}

// URI validates a nonempty absolute URI, including application, mail, and phone schemes.
func (r *Reader) URI(slot string, schemes ...string) string {
	value := r.Node.Text(slot)
	if !r.Node.Has(slot) {
		return value
	}
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme == "" || (parsed.Host == "" && parsed.Opaque == "" && parsed.Path == "") || (len(schemes) != 0 && !contains(schemes, parsed.Scheme)) {
		r.invalid(slot, "invalid absolute URI")
	} else if parsed.Scheme == "http" || parsed.Scheme == "https" {
		r.URL(slot, false)
	}
	return value
}
func (r *Reader) Count(slot string, min, max int) []Node {
	children := r.Node.Children(slot)
	if len(children) < min || (max != 0 && len(children) > max) {
		r.Fail("limit_exceeded", slot+" has an invalid number of children")
	}
	return children
}
