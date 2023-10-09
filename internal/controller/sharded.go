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

package controller

import (
	"github.com/al-assad/doris-operator/internal/util"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
)

// StCtrlErrSet is the standard controller error container
type StCtrlErrSet struct {
	Rec    error
	Sync   error
	Update error
}

func (r *StCtrlErrSet) AsResult() (ctrl.Result, error) {
	// Silent update conflict error
	updateConflict := false
	if errors.IsConflict(r.Update) {
		r.Update = nil
		updateConflict = true
	}
	mergedErr := util.MergeErrorsWithTag(map[string]error{
		"rec":           r.Rec,
		"sync":          r.Sync,
		"update-status": r.Update,
	})
	if mergedErr == nil {
		if updateConflict {
			return ctrl.Result{Requeue: true}, nil
		} else {
			return ctrl.Result{}, nil
		}
	} else {
		return ctrl.Result{Requeue: true}, mergedErr
	}
}
