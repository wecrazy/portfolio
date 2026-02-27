// Package sanitize wraps bluemonday with a shared, pre-configured policy so
// callers do not each allocate their own policy instance. It is safe to import
// from any package in the project.
package sanitize

import "github.com/microcosm-cc/bluemonday"

// strict is the shared bluemonday strict policy (strips all HTML).
var strict = bluemonday.StrictPolicy()

// Strict sanitizes input using bluemonday's StrictPolicy, which strips all
// HTML tags and attributes, returning only the safe plain text.
func Strict(input string) string {
	return strict.Sanitize(input)
}
