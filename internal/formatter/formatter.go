package formatter

import (
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
