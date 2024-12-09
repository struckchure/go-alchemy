package internals

import (
	"strings"

	"github.com/samber/lo"
)

func Parse(input string) (*string, error) {
	input = strings.ReplaceAll(input, "// @alchemy block", "")
	lines := strings.Split(input, "\n")

	for lineIdx, line := range lines {
		line = strings.TrimSpace(line)

		// Match @alchemy statement
		if strings.HasPrefix(line, "// @alchemy statement") {
			lines[lineIdx+1] = strings.Replace(line, "// @alchemy statement", "", 1)
			lines = append(lines[:lineIdx], lines[lineIdx+1:]...)

			continue
		}

		// Match @alchemy replace
		if strings.HasPrefix(line, "// @alchemy replace") {
			replacement := strings.TrimSpace(line[len("// @alchemy replace"):])

			lines[lineIdx+1] = replacement
			lines[lineIdx] = strings.Replace(lines[lineIdx], line, replacement, 1)
			lines = append(lines[:lineIdx], lines[lineIdx+1:]...)

			continue
		}
	}

	return lo.ToPtr(strings.Join(lines, "\n")), nil
}
