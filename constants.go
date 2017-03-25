package main

import (
	"fmt"
	"strconv"
	"strings"
	"regexp"
)

type Table struct {
	Name      string
	Constants string
	Labels    []Label
}

func (t *Table) IsPresent() bool {
	return len(t.Labels) > 0
}

type Label struct {
	Name   string
	Offset uint
}

func DefineTable(constants []string, tableName string) Table {

	labels := []Label{}
	bytes := make([]byte, 0, 1000)

	for _, line := range constants {

		if strings.HasSuffix(line, ":") {
			labels = append(labels, Label{Name: line[:len(line)-1], Offset: uint(len(bytes))})
		} else if strings.Contains(line, ".byte") {
			v, _ := strconv.Atoi(strings.Fields(line)[1])
			bytes = append(bytes, byte(v))
		} else if strings.Contains(line, ".short") {
			v, _ := strconv.Atoi(strings.Fields(line)[1])
			bytes = append(bytes, byte(v))
			bytes = append(bytes, byte(v>>8))
		} else if strings.Contains(line, ".long") {
			v, _ := strconv.Atoi(strings.Fields(line)[1])
			bytes = append(bytes, byte(v))
			bytes = append(bytes, byte(v>>8))
			bytes = append(bytes, byte(v>>16))
			bytes = append(bytes, byte(v>>24))
		} else if strings.Contains(line, ".quad") {
			v, err := strconv.Atoi(strings.Fields(line)[1])
			if err != nil {
				panic(fmt.Sprintf("Atoi error for .quad: %v", err))
			}
			bytes = append(bytes, byte(v))
			bytes = append(bytes, byte(v>>8))
			bytes = append(bytes, byte(v>>16))
			bytes = append(bytes, byte(v>>24))
			bytes = append(bytes, byte(v>>32))
			bytes = append(bytes, byte(v>>40))
			bytes = append(bytes, byte(v>>48))
			bytes = append(bytes, byte(v>>56))
		} else if strings.Contains(line, ".align") {
			bits, _ := strconv.Atoi(strings.Fields(line)[1])
			align := 1 << uint(bits)
			for len(bytes)%align != 0 {
				bytes = append(bytes, 0)
			}
		} else if strings.Contains(line, ".zero") {
			length, _ := strconv.Atoi(strings.Fields(line)[1])
			for i := 0; i < length; i++ {
				bytes = append(bytes, byte(0))
			}
		} else if strings.Contains(line, ".space") {
			argument := strings.Fields(line)[1]
			args := strings.Split(argument, ",")
			length, _ := strconv.Atoi(args[0])
			b := 0
			if len(args) > 1 {
				b, _ = strconv.Atoi(args[1])
			}
			for i := 0; i < length; i++ {
				bytes = append(bytes, byte(b))
			}
		} else if strings.Contains(line, ".section") {
			// ignore
		} else if strings.Contains(line, ".text") {
			// ignore
		} else {
			panic(fmt.Sprintf("Unknown line for table: %s", line))
		}
	}

	// Pad onto a multiple of 8 bytes for aligned outputting
	for len(bytes)%8 != 0 {
		bytes = append(bytes, 0)
	}

	table := []string{}

	for i := 0; i < len(bytes); i += 8 {
		offset := fmt.Sprintf("%03x", i)
		hex := ""
		for j := i; j < i+8 && j < len(bytes); j++ {
			hex = fmt.Sprintf("%02x", bytes[j]) + hex
		}
		table = append(table, fmt.Sprintf("DATA %s<>+0x%s(SB)/8, $0x%s", tableName, offset, hex))
	}
	table = append(table, fmt.Sprintf("GLOBL %s<>(SB), 8, $%d", tableName, len(bytes)))

	return Table{Name: tableName, Constants: strings.Join(table, "\n"), Labels: labels}
}

func SplitOnGlobals(lines []string) []int {

	var result []int

	for index, line := range lines {
		if strings.Contains(line, ".globl") {
			result = append(result, index)
		}
	}

	return result
}

var regexpLabelConstant = regexp.MustCompile(`^\.?LCPI[0-9]+_0:`)

func GetFirstLabelConstants(lines []string) int {

	for iline, line := range lines {
		if match := regexpLabelConstant.FindStringSubmatch(line); len(match) > 0 {
			return iline
		}
	}

	return -1
}

func SegmentConsts(lines []string) []Table {

	consts := []Segment{}

	globals := SplitOnGlobals(lines)

	splitBegin := 0
	for _, splitEnd := range globals {
		start := GetFirstLabelConstants(lines[splitBegin:splitEnd])
		if start != -1 {
			// Add set of lines when a constant table has been found
			consts = append(consts, Segment{Name: fmt.Sprintf("LCDATA%d", len(consts)+1), Start: splitBegin + start, End: splitEnd})
		}
		splitBegin = splitEnd + 1
	}

	tables := []Table{}

	for _, c := range consts {

		tables = append(tables, DefineTable(lines[c.Start:c.End], c.Name))
	}

	return tables
}

func GetCorrespondingTable(lines []string, tables []Table) Table {

	concat := strings.Join(lines, "\n")

	for _, t := range tables {
		// Easy test -- we assume that if we find one label, we would find the others as well...
		if strings.Contains(concat, t.Labels[0].Name) {
			return t
		}
	}

	return Table{}
}
