/*
 * Minio Cloud Storage, (C) 2017 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"strings"
	"testing"
)

func testConstant(t *testing.T, constants, expected string) {

	table := defineTable(strings.Split(constants, "\n"), "LCTABLE")

	if table.Constants != expected {
		t.Errorf("TestConstants(): \nexpected %s\ngot      %s", expected, table.Constants)
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
        .align  1
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
        .align  32
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

	constant3 := `LCPI0_0:
	.space	32,1
LCPI0_1:
	.short	2                       ## 0x2
	.short	2                       ## 0x2
	.short	2                       ## 0x2
	.short	2                       ## 0x2
	.short	2                       ## 0x2
	.short	2                       ## 0x2
	.short	2                       ## 0x2
	.short	2                       ## 0x2
	.short	2                       ## 0x2
	.short	2                       ## 0x2
	.short	2                       ## 0x2
	.short	2                       ## 0x2
	.short	2                       ## 0x2
	.short	2                       ## 0x2
	.short	2                       ## 0x2
	.short	2                       ## 0x2`

	table3 := `DATA LCTABLE<>+0x000(SB)/8, $0x0101010101010101
DATA LCTABLE<>+0x008(SB)/8, $0x0101010101010101
DATA LCTABLE<>+0x010(SB)/8, $0x0101010101010101
DATA LCTABLE<>+0x018(SB)/8, $0x0101010101010101
DATA LCTABLE<>+0x020(SB)/8, $0x0002000200020002
DATA LCTABLE<>+0x028(SB)/8, $0x0002000200020002
DATA LCTABLE<>+0x030(SB)/8, $0x0002000200020002
DATA LCTABLE<>+0x038(SB)/8, $0x0002000200020002
GLOBL LCTABLE<>(SB), 8, $64`

	testConstant(t, constant3, table3)

	constant4 := `.LCPI1_0:
	.byte	255                     # 0xff
	.byte	0                       # 0x0
	.byte	255                     # 0xff
	.byte	0                       # 0x0
	.byte	255                     # 0xff
	.byte	0                       # 0x0
	.byte	255                     # 0xff
	.byte	0                       # 0x0
	.byte	255                     # 0xff
	.byte	0                       # 0x0
	.byte	255                     # 0xff
	.byte	0                       # 0x0
	.byte	255                     # 0xff
	.byte	0                       # 0x0
	.byte	255                     # 0xff
	.byte	0                       # 0x0
.LCPI1_2:
	.quad	281474976776192         # 0x1000000010000
	.quad	281474976776192         # 0x1000000010000
.LCPI1_3:
	.short	1606                    # 0x646
	.short	4211                    # 0x1073
	.short	1606                    # 0x646
	.short	4211                    # 0x1073
	.short	1606                    # 0x646
	.short	4211                    # 0x1073
	.short	1606                    # 0x646
	.short	4211                    # 0x1073
.LCPI1_12:
	.zero	16`

	table4 := `DATA LCTABLE<>+0x000(SB)/8, $0x00ff00ff00ff00ff
DATA LCTABLE<>+0x008(SB)/8, $0x00ff00ff00ff00ff
DATA LCTABLE<>+0x010(SB)/8, $0x0001000000010000
DATA LCTABLE<>+0x018(SB)/8, $0x0001000000010000
DATA LCTABLE<>+0x020(SB)/8, $0x1073064610730646
DATA LCTABLE<>+0x028(SB)/8, $0x1073064610730646
DATA LCTABLE<>+0x030(SB)/8, $0x0000000000000000
DATA LCTABLE<>+0x038(SB)/8, $0x0000000000000000
GLOBL LCTABLE<>(SB), 8, $64`

	testConstant(t, constant4, table4)

	constant5 := `        .p2align        4
.LCPI0_0:
        .long   1127219200              # 0x43300000
        .long   1160773632              # 0x45300000
        .long   0                       # 0x0
        .long   0                       # 0x0
.LCPI0_1:
        .quad   4841369599423283200     # double 4503599627370496
        .quad   4985484787499139072     # double 1.9342813113834067E+25
        .section        .rodata.cst8,"aM",@progbits,8
        .p2align        3
.LCPI0_2:
        .quad   4602678819172646912     # double 0.5
.LCPI0_3:
        .quad   -4620693217682128896    # double -0.5
        .section        .rodata.cst4,"aM",@progbits,4
        .p2align        5, 0x12
.LCPI0_4:
        .long   1098907648              # float 16`

	table5 := `DATA LCTABLE<>+0x000(SB)/8, $0x4530000043300000
DATA LCTABLE<>+0x008(SB)/8, $0x0000000000000000
DATA LCTABLE<>+0x010(SB)/8, $0x4330000000000000
DATA LCTABLE<>+0x018(SB)/8, $0x4530000000000000
DATA LCTABLE<>+0x020(SB)/8, $0x3fe0000000000000
DATA LCTABLE<>+0x028(SB)/8, $0xbfe0000000000000
DATA LCTABLE<>+0x030(SB)/8, $0x1212121212121212
DATA LCTABLE<>+0x038(SB)/8, $0x1212121212121212
DATA LCTABLE<>+0x040(SB)/8, $0x0000000041800000
GLOBL LCTABLE<>(SB), 8, $72`

	testConstant(t, constant5, table5)

	constant6 := `        .p2align        4
.LCPI1_0:
        .zero   16,1
.LCPI1_1:
        .short  4                       # 0x4
        .short  4                       # 0x4
        .short  4                       # 0x4
        .short  4                       # 0x4
        .short  4                       # 0x4
        .short  4                       # 0x4
        .short  4                       # 0x4
        .short  4                       # 0x4`

	table6 := `DATA LCTABLE<>+0x000(SB)/8, $0x0101010101010101
DATA LCTABLE<>+0x008(SB)/8, $0x0101010101010101
DATA LCTABLE<>+0x010(SB)/8, $0x0004000400040004
DATA LCTABLE<>+0x018(SB)/8, $0x0004000400040004
GLOBL LCTABLE<>(SB), 8, $32`

	testConstant(t, constant6, table6)
}
