package formatter

import (
	"html/template"
	"strconv"
	"strings"
)

func FormatRupiah(amount int) string {
	digits := strconv.Itoa(amount)
	var groups []string
	for len(digits) > 3 {
		groups = append([]string{digits[len(digits)-3:]}, groups...)
		digits = digits[:len(digits)-3]
	}
	groups = append([]string{digits}, groups...)

	result := "Rp. " + strings.Join(groups, ",")
	return result
}

// SafeQuery marks an already percent-encoded query string (e.g. from
// url.Values.Encode()) as safe, so html/template embeds it verbatim in a
// href instead of re-escaping its "=" and "&" separators.
func SafeQuery(query string) template.URL {
	return template.URL(query)
}
