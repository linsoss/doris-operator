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
	"fmt"
	"github.com/stretchr/testify/assert"
	"k8s.io/utils/strings/slices"
	"testing"
)

func TestMergeErrors(t *testing.T) {
	err1 := MergeErrors(fmt.Errorf("foo"), fmt.Errorf("bar"))
	assert.Equal(t, "foo; bar", err1.Error())

	err2 := MergeErrors(nil, fmt.Errorf("foo"), fmt.Errorf("bar"))
	assert.Equal(t, "foo; bar", err2.Error())

	err3 := MergeErrors(err1, fmt.Errorf("baz"))
	assert.Equal(t, "foo; bar; baz", err3.Error())

	err4 := MergeErrors(nil, nil)
	assert.Nil(t, err4)
}

func TestMergeErrorsWithTags(t *testing.T) {
	err1 := MergeErrorsWithTag(map[string]error{
		"tag1": fmt.Errorf("foo"),
		"tag2": fmt.Errorf("bar"),
	})
	assert.True(t, slices.Contains([]string{"[tag1] foo; [tag2] bar", "[tag2] bar; [tag1] foo"}, err1.Error()))

	err2 := MergeErrorsWithTag(map[string]error{
		"tag1": fmt.Errorf("foo"),
		"tag2": nil,
	})
	assert.Equal(t, "[tag1] foo", err2.Error())

	err3 := MergeErrorsWithTag(map[string]error{
		"tag1": nil,
		"tag2": nil,
	})
	assert.Nil(t, err3)
}
