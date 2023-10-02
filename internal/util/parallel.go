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

// ParallelRun executes the functions in parallel and wait for all functions to complete.
func ParallelRun[T any](fns ...func() T) []T {
	chans := make([]chan T, len(fns))
	for i := range chans {
		chans[i] = make(chan T)
	}
	res := make([]T, len(fns))

	for i, fn := range fns {
		f := fn
		ch := chans[i]
		go func() {
			r := f()
			ch <- r
		}()
	}
	for i, ch := range chans {
		r := <-ch
		res[i] = r
	}
	return res
}
