package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Label struct {
	Name   string
	Offset uint
}

func DefineTable(constants, tableName string) (string, []Label) {

	lines := strings.Split(constants, "\n")

	labels := []Label{}
	bytes := make([]byte, 0, 1000)

	for _, line := range lines {

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
		} else if strings.Contains(line, ".section") {
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

	//fmt.Println(strings.Join(table, "\n"))
	//fmt.Println(labels)

	return strings.Join(table, "\n"), labels
}
