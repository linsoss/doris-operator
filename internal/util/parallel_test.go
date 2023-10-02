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

func TestParallelRunCase1(t *testing.T) {
	fns := []func() string{
		func() string { return "r1" },
		func() string { return "r2" },
		func() string { return "r3" },
	}
	res := ParallelRun(fns...)
	assert.Equal(t, []string{"r1", "r2", "r3"}, res)
}

func TestParallelRunCase2(t *testing.T) {
	fns := []func() func([]string) []string{
		func() func([]string) []string {
			return func(arr []string) []string {
				return append(arr, "r1")
			}
		},
		func() func([]string) []string {
			return func(arr []string) []string {
				return append(arr, "r2")
			}
		},
		func() func([]string) []string {
			return func(arr []string) []string {
				return append(arr, "r3")
			}
		},
	}
	resFns := ParallelRun(fns...)
	var res []string
	for _, fn := range resFns {
		res = fn(res)
	}
	assert.Equal(t, []string{"r1", "r2", "r3"}, res)
}
