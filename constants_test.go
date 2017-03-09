package main

import (
	"testing"
)

func testConstant(t *testing.T, constants, expected string) {

	table, _ := DefineTable(constants, "LCTABLE")

	if table != expected {
		t.Errorf("TestConstants(): \nexpected %s\ngot      %s", expected, table)
	}
}

func TestConstants(t *testing.T) {

	constant1 := `LCPI0_0:
        .byte   255                     ## 0xff
        .byte   0                       ## 0x0
        .byte   255                     ## 0xff
        .byte   0                       ## 0x0
        .byte   255                     ## 0xff
        .byte   0                       ## 0x0
        .byte   255                     ## 0xff
        .byte   0                       ## 0x0
        .byte   255                     ## 0xff
        .byte   0                       ## 0x0
        .byte   255                     ## 0xff
        .byte   0                       ## 0x0
        .byte   255                     ## 0xff
        .byte   0                       ## 0x0
        .byte   255                     ## 0xff
        .byte   0                       ## 0x0
        .byte   255                     ## 0xff
        .byte   0                       ## 0x0
        .byte   255                     ## 0xff
        .byte   0                       ## 0x0
        .byte   255                     ## 0xff
        .byte   0                       ## 0x0
        .byte   255                     ## 0xff
        .byte   0                       ## 0x0
        .byte   255                     ## 0xff
        .byte   0                       ## 0x0
        .byte   255                     ## 0xff
        .byte   0                       ## 0x0
        .byte   255                     ## 0xff
        .byte   0                       ## 0x0
        .byte   255                     ## 0xff
        .byte   0                       ## 0x0
LCPI0_1:
        .short  9617                    ## 0x2591
        .short  0                       ## 0x0
        .short  9617                    ## 0x2591
        .short  0                       ## 0x0
        .short  9617                    ## 0x2591
        .short  0                       ## 0x0
        .short  9617                    ## 0x2591
        .short  0                       ## 0x0
        .short  9617                    ## 0x2591
        .short  0                       ## 0x0
        .short  9617                    ## 0x2591
        .short  0                       ## 0x0
        .short  9617                    ## 0x2591
        .short  0                       ## 0x0
        .short  9617                    ## 0x2591
        .short  0                       ## 0x0
LCPI0_2:
        .short  1868                    ## 0x74c
        .short  4899                    ## 0x1323
        .short  1868                    ## 0x74c
        .short  4899                    ## 0x1323
        .short  1868                    ## 0x74c
        .short  4899                    ## 0x1323
        .short  1868                    ## 0x74c
        .short  4899                    ## 0x1323
        .short  1868                    ## 0x74c
        .short  4899                    ## 0x1323
        .short  1868                    ## 0x74c
        .short  4899                    ## 0x1323
        .short  1868                    ## 0x74c
        .short  4899                    ## 0x1323
        .short  1868                    ## 0x74c
        .short  4899                    ## 0x1323
        .section        __TEXT,__literal4,4byte_literals
        .align  2
LCPI0_3:
        .long   8192                    ## 0x2000`

	table1 := `DATA LCTABLE<>+0x000(SB)/8, $0x00ff00ff00ff00ff
DATA LCTABLE<>+0x008(SB)/8, $0x00ff00ff00ff00ff
DATA LCTABLE<>+0x010(SB)/8, $0x00ff00ff00ff00ff
DATA LCTABLE<>+0x018(SB)/8, $0x00ff00ff00ff00ff
DATA LCTABLE<>+0x020(SB)/8, $0x0000259100002591
DATA LCTABLE<>+0x028(SB)/8, $0x0000259100002591
DATA LCTABLE<>+0x030(SB)/8, $0x0000259100002591
DATA LCTABLE<>+0x038(SB)/8, $0x0000259100002591
DATA LCTABLE<>+0x040(SB)/8, $0x1323074c1323074c
DATA LCTABLE<>+0x048(SB)/8, $0x1323074c1323074c
DATA LCTABLE<>+0x050(SB)/8, $0x1323074c1323074c
DATA LCTABLE<>+0x058(SB)/8, $0x1323074c1323074c
DATA LCTABLE<>+0x060(SB)/8, $0x0000000000002000
GLOBL LCTABLE<>(SB), 8, $104`

	testConstant(t, constant1, table1)

	constant2 := `LCPI0_0:
        .quad   72340172838076673       ## 0x101010101010101
LCPI0_2:
        .quad   4294967297              ## 0x100000001
        .section        __TEXT,__const
        .align  5
LCPI0_1:
        .long   0                       ## 0x0
        .long   2                       ## 0x2
        .long   4                       ## 0x4
        .long   6                       ## 0x6
        .long   1                       ## 0x1
        .long   3                       ## 0x3
        .long   5                       ## 0x5
        .long   7                       ## 0x7
        .section        __TEXT,__literal4,4byte_literals
        .align  2
LCPI0_3:
        .long   1065353216              ## 0x3f800000`

	table2 := `DATA LCTABLE<>+0x000(SB)/8, $0x0101010101010101
DATA LCTABLE<>+0x008(SB)/8, $0x0000000100000001
DATA LCTABLE<>+0x010(SB)/8, $0x0000000000000000
DATA LCTABLE<>+0x018(SB)/8, $0x0000000000000000
DATA LCTABLE<>+0x020(SB)/8, $0x0000000200000000
DATA LCTABLE<>+0x028(SB)/8, $0x0000000600000004
DATA LCTABLE<>+0x030(SB)/8, $0x0000000300000001
DATA LCTABLE<>+0x038(SB)/8, $0x0000000700000005
DATA LCTABLE<>+0x040(SB)/8, $0x000000003f800000
GLOBL LCTABLE<>(SB), 8, $72`

	testConstant(t, constant2, table2)

}
