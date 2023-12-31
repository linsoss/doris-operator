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

func TestStringFallback(t *testing.T) {
	assert.Equal(t, "foo", StringFallback("foo", "bar"))
	assert.Equal(t, "bar", StringFallback("", "bar"))
}

func TestPointerFallback(t *testing.T) {
	type StringPtr struct {
		S string
	}
	assert.Equal(t, "foo", PointerFallback(&StringPtr{"foo"}, &StringPtr{"bar"}).S)
	assert.Equal(t, "bar", PointerFallback(nil, &StringPtr{"bar"}).S)
}

func TestPointerFallbackAndDeRefer(t *testing.T) {
	type StringPtr struct {
		S string
	}
	assert.Equal(t, "foo", PointerFallbackAndDeRefer(&StringPtr{"foo"}, &StringPtr{"bar"}, StringPtr{"baz"}).S)
	assert.Equal(t, "baz", PointerFallbackAndDeRefer(nil, nil, StringPtr{"baz"}).S)
}

func TestElvis(t *testing.T) {
	// normal value
	assert.Equal(t, "foo", Elvis(true, "foo", "bar"))
	assert.Equal(t, "bar", Elvis(false, "foo", "bar"))

	// function value
	fn1 := func() string {
		return "foo"
	}
	fn2 := func() string {
		return "bar"
	}
	assert.Equal(t, "foo", Elvis(true, fn1, fn2)())
}

func TestMd5Hash(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"c": 3,
		"b": 2,
	}
	assert.Equal(t, "9e0bf104708effc55357dc36f9426ce7", Md5HashOr(m, ""))
}
