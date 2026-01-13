package utils

import (
	"fmt"
	"github.com/gosimple/slug"
)

func GenerateSlug(title string) string {
	return slug.Make(title)
}

// GenerateUniqueSlug generates a unique slug by appending a counter if necessary
func GenerateUniqueSlug(title string, existsFn func(string) bool) string {
	s := slug.Make(title)
	if !existsFn(s) {
		return s
	}

	for i := 1; ; i++ {
		newSlug := fmt.Sprintf("%s-%d", s, i)
		if !existsFn(newSlug) {
			return newSlug
		}
	}
}
