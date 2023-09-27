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

// StringFallback returns the first string if it is not empty, otherwise the second string.
func StringFallback(str string, fallback string) string {
	if str != "" {
		return str
	} else {
		return fallback
	}
}

// PointerFallback returns the first pointer if it is not nil, otherwise the second pointer.
func PointerFallback[T interface{}](pointer *T, fallback *T) *T {
	if pointer != nil {
		return pointer
	} else {
		return fallback
	}
}

// PointerFallbackAndDeRefer returns the first pointer if it is not nil, otherwise the second pointer,
// then dereference the pointer fallback with the defaultValue.
func PointerFallbackAndDeRefer[T interface{}](pointer *T, fallback *T, defaultValue T) T {
	result := PointerFallback(pointer, fallback)
	if result != nil {
		return *result
	} else {
		return defaultValue
	}
}

// ArrayFallback returns the first array if it is not empty, otherwise the second array.
func ArrayFallback[T interface{}](array []T, fallback []T) []T {
	if len(array) != 0 {
		return array
	} else {
		return fallback
	}
}
