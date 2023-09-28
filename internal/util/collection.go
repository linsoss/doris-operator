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
	"crypto/md5"
	"encoding/json"
	"fmt"
	"sort"
)

// MergeMaps merges two maps into one,
// when there is a key conflict, prioritize using the value of m2.
func MergeMaps[K comparable, V any](m1, m2 map[K]V) map[K]V {
	if len(m1) == 0 {
		return m2
	}
	if len(m2) == 0 {
		return m1
	}
	merged := make(map[K]V, len(m1)+len(m2))
	for k, v := range m1 {
		merged[k] = v
	}
	for k, v := range m2 {
		merged[k] = v
	}
	return merged
}

// MapEqual checks if two maps are equal.
func MapEqual[K comparable, V comparable](m1, m2 map[K]V) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v1 := range m1 {
		v2, ok := m2[k]
		if !ok || v1 != v2 {
			return false
		}
	}
	return true
}

// MapSortedKeys returns all sorted keys of a map.
func MapSortedKeys[K string | int, V any](data map[K]V) []K {
	keys := make([]K, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}

// MapMd5 returns the md5 string of a map.
// When the map is empty, return an empty string.
func MapMd5[K comparable, V any](data map[K]V) string {
	if len(data) == 0 {
		return ""
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	hash := md5.Sum(jsonData)
	md5String := fmt.Sprintf("%x", hash)
	return md5String
}
