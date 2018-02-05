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

package c2goasmtest

import (
	"testing"
	"unsafe"
)

func testClibFloor32(t *testing.T, fl, expected float32) {

	got := _ClibFloor32(fl)
	if expected != got {
		t.Errorf("testClibFloor32(): \nexpected %v\ngot      %v", expected, got)
	}
}

func testClibFloor64(t *testing.T, fl, expected float64) {

	got := _ClibFloor64(fl)
	if expected != got {
		t.Errorf("testClibFloor64(): \nexpected %v\ngot      %v", expected, got)
	}
}

func TestClibFloor(t *testing.T) {

	testClibFloor32(t, 2.1, 2.0)
	testClibFloor32(t, 1.9, 1.0)
	testClibFloor32(t, 1.5, 1.0)
	testClibFloor32(t, 1.1, 1.0)
	testClibFloor32(t, 1.0, 1.0)
	testClibFloor32(t, 1.0-1e-6, 0.0)
	testClibFloor32(t, 0.0-1e-6, -1.0)

	testClibFloor64(t, 2.1, 2.0)
	testClibFloor64(t, 1.9, 1.0)
	testClibFloor64(t, 1.5, 1.0)
	testClibFloor64(t, 1.1, 1.0)
	testClibFloor64(t, 1.0, 1.0)
	testClibFloor64(t, 1.0-1e-6, 0.0)
	testClibFloor64(t, 0.0-1e-6, -1.0)

}

func TestClibMemcpy(t *testing.T) {

	src := make([]byte, 256)
	zero := make([]byte, 256)
	dst := make([]byte, 256)

	for i, _ := range src {
		src[i] = byte(i)
	}

	for count := 0; count < 256; count++ {

		copy(dst[:], zero[:])

		ptr := _ClibMemcpy(unsafe.Pointer(&dst[0]), unsafe.Pointer(&src[0]), uint(count))
		if unsafe.Pointer(&dst[0]) != ptr {
			t.Errorf("TestClibMemcpy(): \nexpected %v\ngot      %v", unsafe.Pointer(&dst[0]), ptr)
		}

		i := 0
		for ; i < count; i++ {
			if dst[i] != src[i] {
				t.Errorf("TestClibMemcpy(): \nexpected %d\ngot      %d", src[i], dst[i])
			}
		}
		for ; i < len(dst); i++ {
			if dst[i] != 0 {
				t.Errorf("TestClibMemcpy(): \nexpected %d\ngot      %d", 0, dst[i])
			}
		}
	}
}

func TestClibMemset(t *testing.T) {

	init := make([]byte, 256)
	dst := make([]byte, 256)

	for i, _ := range init {
		init[i] = byte(i)
	}

	for count := 0; count < 256; count++ {

		copy(dst[:], init[:])

		ptr := _ClibMemset(unsafe.Pointer(&dst[0]), count, uint(count))
		if unsafe.Pointer(&dst[0]) != ptr {
			t.Errorf("TestClibMemcpy(): \nexpected %v\ngot      %v", unsafe.Pointer(&dst[0]), ptr)
		}

		i := 0
		for ; i < count; i++ {
			if dst[i] != byte(count) {
				t.Errorf("1-TestClibMemcpy(%d): \nexpected %d\ngot      %d", i, count, dst[i])
			}
		}
		for ; i < len(dst); i++ {
			if dst[i] != init[i] {
				t.Errorf("2-TestClibMemcpy(%d): \nexpected %d\ngot      %d", i, init[i], dst[i])
			}
		}
	}
}
