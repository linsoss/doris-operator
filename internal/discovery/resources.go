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

package discovery

import (
	"errors"
	"fmt"
	dapi "github.com/al-assad/doris-operator/api/v1beta1"
	tran "github.com/al-assad/doris-operator/internal/transformer"
	u "github.com/rjNemo/underscore"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"strings"
)

// check if the fe service is alive and ready
func (r *DorisDiscovery) checkFeSvcReady() *RecErr {
	statefulset := &appv1.StatefulSet{}
	exist, err := r.Exist(tran.GetFeStatefulSetKey(r.CR.ObjKey()), statefulset)
	if err != nil {
		return NewRecErr(err)
	}
	if !exist {
		return NewRecErr(errors.New("FE statefulset has not been created yet"))
	}
	if statefulset.Status.ReadyReplicas < 1 {
		return NewRecErr(errors.New("FE statefulset has not been ready yet"))
	}
	return nil
}

func (r *DorisDiscovery) createSqlConnConf() (*DorisSqlConnConf, *RecErr) {
	// find sql account from secret
	sqlAcc, accErr := r.getOprSqlAccount()
	if accErr != nil {
		return nil, NewRecErr(accErr)
	}
	// create sql connection
	sqlConnConf := DorisSqlConnConf{
		Host:     tran.GetFeServiceDNS(r.CR.ObjKey()),
		Port:     tran.GetFeQueryPort(r.CR),
		User:     sqlAcc.User,
		Password: sqlAcc.Password,
	}
	return &sqlConnConf, nil
}

func (r *DorisDiscovery) getOprSqlAccount() (SqlAccount, error) {
	secretRef := tran.GetOprSqlAccountSecretKey(r.CR.ObjKey())
	secret := &corev1.Secret{}
	exist, err := r.Exist(secretRef, secret)
	if err != nil {
		return SqlAccount{}, err
	}
	if !exist {
		return SqlAccount{}, nil
	}
	sqlAccount := SqlAccount{
		User:     string(secret.Data["user"]),
		Password: string(secret.Data["password"]),
	}
	return sqlAccount, nil
}

func GetFeExpectedHosts(cr *dapi.DorisCluster) []string {
	podNames := tran.GetFeExpectPodNames(cr.ObjKey(), cr.Spec.FE.Replicas)
	peerSvcName := tran.GetFePeerServiceKey(cr.ObjKey()).Name
	res := u.Map(podNames, func(podName string) string {
		return fmt.Sprintf("%s.%s.%s.svc.cluster.local", podName, peerSvcName, cr.Namespace)
	})
	return res
}

func GetBeExpectedHosts(cr *dapi.DorisCluster) []string {
	if cr.Spec.BE == nil {
		return []string{}
	}
	podNames := tran.GetBeExpectPodNames(cr.ObjKey(), cr.Spec.BE.Replicas)
	peerSvcName := tran.GetBePeerServiceKey(cr.ObjKey()).Name
	res := u.Map(podNames, func(podName string) string {
		return fmt.Sprintf("%s.%s.%s.svc.cluster.local", podName, peerSvcName, cr.Namespace)
	})
	return res
}

func GetCnExpectedHosts(cr *dapi.DorisCluster) []string {
	if cr.Spec.CN == nil {
		return []string{}
	}
	podNames := tran.GetCnExpectPodNames(cr.ObjKey(), cr.Spec.CN.Replicas)
	peerSvcName := tran.GetCnPeerServiceKey(cr.ObjKey()).Name
	res := u.Map(podNames, func(podName string) string {
		return fmt.Sprintf("%s.%s.%s.svc.cluster.local", podName, peerSvcName, cr.Namespace)
	})
	return res
}

func GetBrokerExpectedHosts(cr *dapi.DorisCluster) []string {
	if cr.Spec.Broker == nil {
		return []string{}
	}
	podNames := tran.GetBrokerExpectPodNames(cr.ObjKey(), cr.Spec.Broker.Replicas)
	peerSvcName := tran.GetBrokerPeerServiceKey(cr.ObjKey()).Name
	res := u.Map(podNames, func(podName string) string {
		return fmt.Sprintf("%s.%s.%s.svc.cluster.local", podName, peerSvcName, cr.Namespace)
	})
	return res
}

func GetBrokerExpectedNames(cr *dapi.DorisCluster) []string {
	if cr.Spec.Broker == nil {
		return []string{}
	}
	podNames := tran.GetBrokerExpectPodNames(cr.ObjKey(), cr.Spec.Broker.Replicas)
	res := u.Map(podNames, func(podName string) string {
		return GetBrokerNameByPodName(podName)
	})
	return res
}

func GetBrokerNameByPodName(bkPodName string) string {
	return strings.ReplaceAll(bkPodName, "-", "_")
}

func GetBrokerPodNameByName(bkPodName string) string {
	return strings.ReplaceAll(bkPodName, "_", "-")
}
