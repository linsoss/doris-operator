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

package reconciler

import (
	"context"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ReconcileContext struct {
	client.Client
	schema *runtime.Scheme
	req    ctrl.Request
	ctx    context.Context
	log    logr.Logger
}

// Exist checks if the kubernetes object exists.
func (r *ReconcileContext) Exist(key types.NamespacedName, objType client.Object) (bool, error) {
	if err := r.Get(r.ctx, key, objType); err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// DeleteWhenExist deletes the kubernetes object if it exists.
func (r *ReconcileContext) DeleteWhenExist(key types.NamespacedName, objType client.Object) error {
	exist, err := r.Exist(key, objType)
	if err != nil {
		return err
	}
	if exist {
		if err := r.Delete(r.ctx, objType); err != nil {
			return err
		}
	}
	return nil
}

// CreateOrUpdate creates or updates the kubernetes object.
func (r *ReconcileContext) CreateOrUpdate(obj client.Object) error {
	key := client.ObjectKeyFromObject(obj)
	objType := reflect.TypeOf(obj)
	emptyObj := reflect.New(objType).Elem().Interface().(client.Object)

	exist, err := r.Exist(key, emptyObj)
	if err != nil {
		return err
	}
	if !exist {
		return r.Create(r.ctx, obj)
	} else {
		return r.Update(r.ctx, obj)
	}
}
