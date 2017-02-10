// package text is a utility library of various text/string manipulations
package text

import (
	"regexp"
	"strings"
	"path"
)

// Sanitize filename string for FILE/URL output but removing non-alphanumerics and trimming space
func SanitizeFilename(s string) string {
	fileNoNos := regexp.MustCompile(`[^[:alnum:]-]`)
	s = strings.Trim(s, " ")
	s = strings.Replace(s, " ", "-", -1)
	s = fileNoNos.ReplaceAllString(s, "-")
	return s
}

// Remove the filename extension
func RemoveExt(src string) string {
	return strings.TrimSuffix(src, path.Ext(src))
}
