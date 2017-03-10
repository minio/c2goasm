package main

import (
	"strings"
	"unicode"
)

func isLower(str string) bool {

	for _, r := range str {
		return unicode.IsLower(r)
	}
	return false
}

func removeUndefined(line, undef string) string {

	if parts := strings.SplitN(line, undef, 2); len(parts) > 1 {
		line = parts[0] + strings.TrimSpace(parts[1])
	}
	return line
}

func assemblify(lines []string, table Table) ([]string, error) {

	var result []string

	for _, line := range lines {

		// Remove ## comments
		if parts := strings.SplitN(line, `##`, 2); len(parts) > 1 {
			if strings.TrimSpace(parts[0]) == "" {
				continue
			}
			line = parts[0]
		}

		// Skip lines with aligns
		if strings.Contains(line, ".align") {
			continue
		}

		// Make jmps uppercase
		if parts := strings.SplitN(line, `LBB0`, 2); len(parts) > 1 {
			// unless it is a label
			if !strings.Contains(parts[1], ":") {
				// make jmp statement uppercase
				line = strings.ToUpper(parts[0]) + "LBB0" + parts[1]
			}
		}

		fields := strings.Fields(line)
		// Test for any non-jmp instruction (lower case mnemonic)
		if len(fields) > 0 && !strings.Contains(fields[0], ":") && isLower(fields[0]) {
			// prepend line with comment for subsequent asm2plan9s assembly
			line = "                                 // " + strings.TrimSpace(line)
		}

		line = removeUndefined(line, "ptr")
		line = removeUndefined(line, "xmmword")
		line = removeUndefined(line, "ymmword")

		// TODO
		// shr/sar without arg --> add , 1
		// replace PIC load ([rip] based)
		// strip header
		// strip footer
		// add golang header
		// add golang footer
		// consistent use of rbp & rsp
		result = append(result, line)
	}

	return result, nil
}

