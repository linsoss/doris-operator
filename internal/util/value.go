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
)

// StringFallback returns the first string if it is not empty, otherwise the second string.
func StringFallback(str string, fallback string) string {
	if str != "" {
		return str
	} else {
		return fallback
	}
}

// ArrayFallback returns the first array if it is not empty, otherwise the second array.
func ArrayFallback[T any](array []T, fallback []T) []T {
	if len(array) != 0 {
		return array
	} else {
		return fallback
	}
}

func MapFallback[K comparable, V any](mapValue map[K]V, fallback map[K]V) map[K]V {
	if mapValue == nil || len(mapValue) == 0 {
		return fallback
	} else {
		return mapValue
	}
}

// PointerFallback returns the first pointer if it is not nil, otherwise the second pointer.
func PointerFallback[T any](pointer *T, fallback *T) *T {
	if pointer != nil {
		return pointer
	} else {
		return fallback
	}
}

// PointerFallbackAndDeRefer returns the first pointer if it is not nil, otherwise the second pointer,
// then dereference the pointer fallback with the defaultValue.
func PointerFallbackAndDeRefer[T any](pointer *T, fallback *T, defaultValue T) T {
	result := PointerFallback(pointer, fallback)
	if result != nil {
		return *result
	} else {
		return defaultValue
	}
}

func PointerDeRefer[T any](pointer *T, defaultValue T) T {
	if pointer != nil {
		return *pointer
	} else {
		return defaultValue
	}
}

// Elvis is a Groovy-like elvis expression.
func Elvis[T any](condition bool, leftValue T, rightValue T) T {
	if condition {
		return leftValue
	} else {
		return rightValue
	}
}

// Md5Hash returns the md5 hash of the given object base on json marshal.
func Md5Hash(obj any) (string, error) {
	if obj == nil {
		return "", nil
	}
	bytes, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	hashStr := fmt.Sprintf("%x", md5.Sum(bytes))
	return hashStr, nil
}

// Md5HashOr returns the md5 hash of the given object base on json marshal.
// when error occurs, return the fallback string.
func Md5HashOr(obj any, fallback string) string {
	hash, err := Md5Hash(obj)
	if err != nil {
		return fallback
	}
	return hash
}
