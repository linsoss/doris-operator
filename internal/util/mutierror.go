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
	"errors"
	"fmt"
	"strings"
)

// MultiError is a list of errors.
type MultiError struct {
	Errors []error
}

func (e *MultiError) Error() string {
	errStrs := make([]string, 0, len(e.Errors))
	for _, err := range e.Errors {
		errStrs = append(errStrs, err.Error())
	}
	return strings.Join(errStrs, "; ")
}

func (e *MultiError) Collect(err error) bool {
	if err != nil {
		e.Errors = append(e.Errors, err)
		return true
	} else {
		return false
	}
}

// CollectFnErr collects error from fn function, then execute rightFn function when the fn returns nil error.
func CollectFnErr[T any](errContainer *MultiError, fn func() (T, error), rightFn func(T)) {
	res, err := fn()
	if err != nil {
		errContainer.Collect(err)
	}
	rightFn(res)
}

// Dry returns nil if there is no error in the container.
func (e *MultiError) Dry() error {
	if len(e.Errors) == 0 {
		return nil
	}
	if len(e.Errors) == 1 {
		return e.Errors[0]
	}
	return e
}

// MergeErrors merges multiple errors into one.
func MergeErrors(errs ...error) *MultiError {
	var errorList []error
	for _, err := range errs {
		if err != nil {
			errorList = append(errorList, err)
		}
	}
	if len(errorList) == 0 {
		return nil
	}
	return &MultiError{Errors: errorList}
}

func AppendErrMsg(err error, message string) *MultiError {
	if err == nil {
		return nil
	}
	return MergeErrors(errors.New(message), err)
}

// MultiTaggedError is a list of errors with tags.
type MultiTaggedError struct {
	Errors map[string]error
}

func (e *MultiTaggedError) Error() string {
	errStrs := make([]string, 0, len(e.Errors))
	for tag, err := range e.Errors {
		errStrs = append(errStrs, fmt.Sprintf("[%s] %s", tag, err.Error()))
	}
	return strings.Join(errStrs, "; ")
}
