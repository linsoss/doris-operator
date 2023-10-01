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
	dapi "github.com/al-assad/doris-operator/api/v1beta1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ReconcileContext is the context for reconciling CRD.
type ReconcileContext struct {
	client.Client
	Schema *runtime.Scheme
	Ctx    context.Context
	Log    logr.Logger
}

func NewReconcileContext(client client.Client, schema *runtime.Scheme, ctx context.Context) ReconcileContext {
	return ReconcileContext{
		Client: client,
		Schema: schema,
		Ctx:    ctx,
		Log:    log.FromContext(ctx),
	}
}

// Exist checks if the kubernetes object exists.
func (r *ReconcileContext) Exist(key types.NamespacedName, objType client.Object) (bool, error) {
	if err := r.Get(r.Ctx, key, objType); err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Find finds the kubernetes object, when it does not exist, the obj pointer will be set to nil.
func (r *ReconcileContext) Find(key types.NamespacedName, obj client.Object) error {
	err := r.Get(r.Ctx, key, obj)
	if errors.IsNotFound(err) {
		obj = nil
		return nil
	} else {
		return err
	}
}

// CreateWhenNotExist creates the kubernetes object if it does not exist.
func (r *ReconcileContext) CreateWhenNotExist(obj client.Object) error {
	key := client.ObjectKeyFromObject(obj)
	objType := reflect.TypeOf(obj)
	objTypeInstance := reflect.New(objType).Elem().Interface().(client.Object)

	exist, err := r.Exist(key, objTypeInstance)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	return r.Create(r.Ctx, obj)
}

// DeleteWhenExist deletes the kubernetes object if it exists.
func (r *ReconcileContext) DeleteWhenExist(key types.NamespacedName, objType client.Object) error {
	exist, err := r.Exist(key, objType)
	if err != nil {
		return err
	}
	if exist {
		if err := r.Delete(r.Ctx, objType); err != nil {
			return err
		}
	}
	return nil
}

// CreateOrUpdate creates or updates the kubernetes object.
func (r *ReconcileContext) CreateOrUpdate(obj client.Object) error {
	key := client.ObjectKeyFromObject(obj)
	objType := reflect.TypeOf(obj)
	prevObj := reflect.New(objType).Elem().Interface().(client.Object)

	exist, err := r.Exist(key, prevObj)
	if err != nil {
		return err
	}
	if !exist {
		return r.Create(r.Ctx, obj)
	} else {
		return r.Update(r.Ctx, obj)
	}
}

// FindRefDorisAutoScaler finds the DorisAutoscaler CR that refer to the DorisCluster CR.
// A DorisCluster CR can only be bound to one additional DorisAutoScaler CR.
func (r *ReconcileContext) FindRefDorisAutoScaler(dorisClusterRef client.ObjectKey) (*dapi.DorisAutoscaler, error) {
	crList := &dapi.DorisAutoscalerList{}
	if err := r.List(r.Ctx, crList, &client.ListOptions{Namespace: dorisClusterRef.Namespace}); err != nil {
		return nil, err
	}
	for _, item := range crList.Items {
		if item.Spec.Cluster == dorisClusterRef.Name {
			return &item, nil
		}
	}
	return nil, nil
}

// FindRefDorisInitializer finds the DorisInitializer CR that refer to the DorisCluster CR.
// A DorisCluster CR can only be bound to one additional DorisInitializer CR.
func (r *ReconcileContext) FindRefDorisInitializer(dorisClusterRef client.ObjectKey) (*dapi.DorisInitializer, error) {
	crList := &dapi.DorisInitializerList{}
	if err := r.List(r.Ctx, crList, &client.ListOptions{Namespace: dorisClusterRef.Namespace}); err != nil {
		return nil, err
	}
	for _, item := range crList.Items {
		if item.Spec.Cluster == dorisClusterRef.Name {
			return &item, nil
		}
	}
	return nil, nil
}
