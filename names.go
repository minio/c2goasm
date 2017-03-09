package main

import (
	"strconv"
	"strings"
	"unicode"
)

func extractPart(part string) (int, string) {

	digits := 0
	for _, d := range part {
		if unicode.IsDigit(d) {
			digits += 1
		} else {
			break
		}
	}
	length, _ := strconv.Atoi(part[:digits])
	return digits + length, part[digits:(digits + length)]
}

func ExtractName(name string) string {

	var parts []string

	for index, ch := range name {
		if unicode.IsDigit(ch) {

			for index < len(name) {
				size, part := extractPart(name[index:])
				if size == 0 {
					break
				}

				parts = append(parts, part)
				index += size
			}

			break
		}
	}

	return strings.Join(parts, "")
}
