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
	dapi "github.com/al-assad/doris-operator/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

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
