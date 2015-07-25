package validate

import (
	"fmt"
	"strings"
)

/// Checks whether its argument is a valid UUID, per
/// https://webbluetoothcg.github.io/web-bluetooth/#dfn-valid-uuid
func ValidUUID(uuid string) bool {
	if len(uuid) != 36 {
		return false
	}
	for index, char := range uuid {
		switch index {
		case 8, 13, 18, 23:
			if char != '-' {
				return false
			}
		default:
			switch {
			case char >= '0' && char <= '9':
			case char >= 'a' && char <= 'f':
			default:
				return false
			}
		}
	}
	return true
}

/// Checks whether its argument is a blacklist file that's usable in the algorithm at
/// https://webbluetoothcg.github.io/web-bluetooth/#dfn-parsing-the-blacklist
func ValidateBlacklist(blacklist string) error {
	for index, line := range strings.Split(blacklist, "\n") {
		line_num := index + 1
		if len(line) == 0 || line[0] == '#' {
			// Comment or blank line.
			continue
		}
		var tokens = strings.Split(line, " ")
		switch len(tokens) {
		case 0:
			panic("strings.Split(\" \") never returns an empty array")
		case 1, 2:
			uuid := tokens[0]
			if !ValidUUID(uuid) {
				return fmt.Errorf("line %v: '%v' is not a valid UUID", line_num, uuid)
			}
			if len(tokens) == 2 {
				exclusion := tokens[1]
				switch exclusion {
				case "exclude-reads", "exclude-writes":
				default:
					return fmt.Errorf("line %v: '%v' should be 'exclude', 'exclude-reads', or 'exclude-writes'",
						line_num, exclusion)
				}
			}
		default:
			return fmt.Errorf("line %v: Too many tokens", line_num)
		}
	}
	return nil
}
