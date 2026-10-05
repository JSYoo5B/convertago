// Package conversion reads tagged source structs and assembles them into
// platform builder nodes. It parses tags against the profile each messenger
// package registers, reads values through generated accessors or cached
// reflection plans, and reports problems as severity-ranked diagnostics.
// Messenger packages turn the nodes into native models and check them with
// their own rules.
package conversion
