package regexp

import (
	"fmt"
	"regexp"
)

type RegexpMatch map[string]string

func Match(pattern string, text string) (RegexpMatch, error) {
	matches := MatchAll(pattern, text)

	if len(matches) == 0 {
		return nil, nil
	}

	if len(matches) > 1 {
		return nil, fmt.Errorf("multiple matches for the pattern: %s", pattern)
	}

	return matches[0], nil
}

func MatchAll(pattern string, text string) []RegexpMatch {
	result := make([]RegexpMatch, 0)
	expr := regexp.MustCompile(pattern)

	for _, m := range expr.FindAllStringSubmatch(text, -1) {
		t := RegexpMatch{}

		for i, name := range expr.SubexpNames() {
			if name != "" {
				t[name] = m[i]
			}
		}

		result = append(result, t)
	}

	return result
}

func Remove(pattern string, s string) string {
	expr := regexp.MustCompile(pattern)
	return expr.ReplaceAllString(s, "")
}
