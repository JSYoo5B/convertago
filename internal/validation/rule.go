package validation

import (
	"net/url"
	"strings"
	"unicode/utf8"
)

// Severity ranks how a rule violation affects message delivery.
type Severity uint8

const (
	// Fatal means the platform rejects the request.
	Fatal Severity = iota + 1
	// Warning means the platform accepts the request, but the result may differ from the user's intent.
	Warning
	// Advisory means UX guidance or client limits may keep the message from displaying as requested.
	Advisory
)

func (s Severity) String() string {
	switch s {
	case Fatal:
		return "fatal"
	case Warning:
		return "warning"
	case Advisory:
		return "advisory"
	}
	return "unknown"
}

// Rule is one platform constraint with the evidence for its severity.
type Rule struct {
	ID        string
	Severity  Severity
	Verified  bool
	Reference string
	Message   string
}

// Violation reports a broken rule at a field relative to the checked value.
type Violation struct {
	Rule  Rule
	Field string
}

// Check runs the rules of one native value.
// Report receives each violation; Visit, when set, receives each child value.
type Check struct {
	Report func(Violation)
	Visit  func(field string, check func(*Check))
}

// Fail reports a violation of rule at field. An empty field means the checked value itself.
func (c *Check) Fail(rule Rule, field string) {
	if c.Report != nil {
		c.Report(Violation{Rule: rule, Field: field})
	}
}

// When reports a violation of rule at field if failed is true.
func (c *Check) When(failed bool, rule Rule, field string) {
	if failed {
		c.Fail(rule, field)
	}
}

// Child passes a child value's checks to the traversal, if any.
func (c *Check) Child(field string, check func(*Check)) {
	if c.Visit != nil {
		c.Visit(field, check)
	}
}

// Walk checks a value and all children it visits, reporting fields as paths below root.
func Walk(root string, check func(*Check), report func(path string, violation Violation)) {
	var visit func(path string, check func(*Check))
	visit = func(path string, check func(*Check)) {
		check(&Check{
			Report: func(v Violation) { report(Join(path, v.Field), v) },
			Visit:  func(field string, child func(*Check)) { visit(Join(path, field), child) },
		})
	}
	visit(root, check)
}

// Join appends a field such as "text" or "blocks[0]" to a path.
func Join(path, field string) string {
	if field == "" {
		return path
	}
	if strings.HasPrefix(field, "[") {
		return path + field
	}
	return path + "." + field
}

// Runes counts the characters of text.
func Runes(text string) int { return utf8.RuneCountInString(text) }

// QueryURIs reports whether text is a query string whose keys all belong to keys
// and whose values are single absolute URIs. At least one key is required.
func QueryURIs(text string, keys ...string) bool {
	values, err := url.ParseQuery(text)
	if err != nil || len(values) == 0 {
		return false
	}
	for key, destinations := range values {
		allowed := false
		for _, candidate := range keys {
			allowed = allowed || key == candidate
		}
		if !allowed || len(destinations) != 1 || !AbsoluteURI(destinations[0]) {
			return false
		}
	}
	return true
}
