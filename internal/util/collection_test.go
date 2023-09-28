/*
 *
 * Copyright 2023 @ Linying Assad <linying@apache.org>
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
 * /
 */

package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMergeMaps(t *testing.T) {
	m1 := map[string]int{
		"a": 1,
		"b": 2,
	}
	m2 := map[string]int{
		"b": 3,
		"c": 4,
	}
	result := MergeMaps(m1, m2)
	expect := map[string]int{
		"a": 1,
		"b": 3,
		"c": 4,
	}
	assert.Equal(t, expect, result)
}

func TestMapEqual(t *testing.T) {
	m1 := map[string]int{
		"a": 1,
		"c": 3,
		"b": 2,
	}
	m2 := map[string]int{
		"c": 3,
		"b": 2,
		"a": 1,
	}
	assert.True(t, MapEqual(m1, m2))
	m1 = map[string]int{
		"a": 1,
		"b": 2,
	}
	m2 = map[string]int{
		"b": 3,
	}
	assert.False(t, MapEqual(m1, m2))
}

func TestMapSortedKeys(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"c": 3,
		"b": 2,
	}
	assert.Equal(t, []string{"a", "b", "c"}, MapSortedKeys(m))
}

func TestMapMd5(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"c": 3,
		"b": 2,
	}
	assert.Equal(t, "9e0bf104708effc55357dc36f9426ce7", MapMd5(m))
	m = make(map[string]int)
	assert.Equal(t, "", MapMd5(m))
}
