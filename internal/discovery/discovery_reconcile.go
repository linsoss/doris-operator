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
	"fmt"
	dapi "github.com/al-assad/doris-operator/api/v1beta1"
	rec "github.com/al-assad/doris-operator/internal/reconciler"
	tran "github.com/al-assad/doris-operator/internal/transformer"
	u "github.com/rjNemo/underscore"
	"golang.org/x/exp/maps"
)

type DorisDiscovery struct {
	rec.ReconcileContext
	CR *dapi.DorisCluster
}

// Notice: The discovery reconcile process is not enabled for the time being for the following reasons:
//  - Directly dropping the backend is an operation that is not secure.
//  - When CN is controlled by HPA, DorisCluster.Spec.CN.Replicas can no longer represent the actual replica
//    number expectation and requires additional inference mechanisms.

func (r *DorisDiscovery) Reconcile() *RecErr {
	if err := r.recDorisFrontends(); err != nil {
		return err
	}
	if err := r.recDorisBackends(); err != nil {
		return err
	}
	if err := r.recDorisBroker(); err != nil {
		return err
	}
	return nil
}

func (r *DorisDiscovery) recDorisFrontends() *RecErr {
	if r.CR.Spec.FE == nil || r.CR.Spec.FE.Replicas == 0 {
		return nil
	}
	if err := r.checkFeSvcReady(); err != nil {
		return err
	}
	// create sql connection
	sqlConnConf, err := r.createSqlConnConf()
	if err != nil {
		return err
	}
	db, connErr := sqlConnConf.Connect()
	if connErr != nil {
		return NewRecSqlErr(connErr)
	}
	// select fe node meta from doris cluster
	feHosts, showErr := ShowFrontendHosts(db)
	if showErr != nil {
		return NewRecSqlErr(showErr)
	}
	// calculate the fe node that should be added or evicted
	expectFeHosts := GetFeExpectedHosts(r.CR)
	addFeHosts := u.Difference(expectFeHosts, feHosts)
	evictFeHosts := u.Difference(feHosts, expectFeHosts)

	// add fe to doris cluster
	for _, host := range addFeHosts {
		hostPort := fmt.Sprintf("%s:%d", host, tran.GetFeEditLogPort(r.CR))
		if err := AddFrontend(db, hostPort); err != nil {
			return NewRecSqlErr(err)
		}
		r.Log.Info(fmt.Sprintf("add frontend[%s] to doris cluster[%s] via connection: %s",
			host, r.CR.ObjKey().String(), sqlConnConf.HostPort()))
	}
	// drop fe from doris cluster
	for _, host := range evictFeHosts {
		hostPort := fmt.Sprintf("%s:%d", host, tran.GetFeEditLogPort(r.CR))
		if err := DropFrontend(db, hostPort); err != nil {
			return NewRecSqlErr(err)
		}
		r.Log.Info(fmt.Sprintf("drop frontend[%s] from doris cluster[%s] via connection: %s",
			host, r.CR.ObjKey().String(), sqlConnConf.HostPort()))
	}
	return nil
}

func (r *DorisDiscovery) recDorisBackends() *RecErr {
	if err := r.checkFeSvcReady(); err != nil {
		return err
	}
	// create sql connection
	sqlConnConf, err := r.createSqlConnConf()
	if err != nil {
		return err
	}
	db, connErr := sqlConnConf.Connect()
	if connErr != nil {
		return NewRecSqlErr(connErr)
	}
	// select be node meta from doris cluster
	beHosts, showErr := ShowBackendHosts(db)
	if showErr != nil {
		return NewRecSqlErr(showErr)
	}
	// calculate the be node that should be added or evicted
	expectHosts := append(GetBeExpectedHosts(r.CR), GetCnExpectedHosts(r.CR)...)
	addBeHosts := u.Difference(expectHosts, beHosts)
	evictBeHosts := u.Difference(beHosts, expectHosts)

	// add be to doris cluster
	for _, host := range addBeHosts {
		hostPort := fmt.Sprintf("%s:%d", host, tran.GetBeHeartbeatServicePort(r.CR))
		if err := AddBackend(db, hostPort); err != nil {
			return NewRecSqlErr(err)
		}
		r.Log.Info(fmt.Sprintf("add backend[%s] to doris cluster[%s] via connection: %s",
			host, r.CR.ObjKey().String(), sqlConnConf.HostPort()))
	}
	// drop be from doris cluster
	for _, host := range evictBeHosts {
		hostPort := fmt.Sprintf("%s:%d", host, tran.GetBeHeartbeatServicePort(r.CR))
		if err := DropBackend(db, hostPort); err != nil {
			return NewRecSqlErr(err)
		}
		r.Log.Info(fmt.Sprintf("drop backend[%s] from doris cluster[%s] via connection: %s",
			host, r.CR.ObjKey().String(), sqlConnConf.HostPort()))
	}
	return nil
}

func (r *DorisDiscovery) recDorisBroker() *RecErr {
	if err := r.checkFeSvcReady(); err != nil {
		return err
	}
	// create sql connection
	sqlConnConf, err := r.createSqlConnConf()
	if err != nil {
		return err
	}
	db, connErr := sqlConnConf.Connect()
	if connErr != nil {
		return NewRecSqlErr(connErr)
	}
	// select broker node meta from doris cluster
	bkNameHosts, showErr := ShowBrokerNameHosts(db)
	if showErr != nil {
		return NewRecSqlErr(showErr)
	}

	// calculate the broker node that should be added or evicted
	expectBkNames := GetBrokerExpectedNames(r.CR)
	actualBkNames := maps.Keys(bkNameHosts)

	addBkNames := u.Difference(expectBkNames, actualBkNames)
	addBkNameHosts := make(map[string]string)
	for _, name := range addBkNames {
		addBkNameHosts[name] = fmt.Sprintf(
			"%s.%s.%s.svc.cluster.local",
			GetBrokerPodNameByName(name), tran.GetBrokerPeerServiceKey(r.CR.ObjKey()), r.CR.Namespace)
	}
	evictBkNames := u.Difference(actualBkNames, expectBkNames)

	// add broker to doris cluster
	for name, host := range addBkNameHosts {
		hostPort := fmt.Sprintf("%s:%d", host, tran.GetBrokerIpcPort(r.CR))
		if err := AddBroker(db, name, hostPort); err != nil {
			return NewRecSqlErr(err)
		}
		r.Log.Info(fmt.Sprintf("add broker[%s] to doris cluster[%s] via connection: %s",
			host, r.CR.ObjKey().String(), sqlConnConf.HostPort()))
	}
	// drop broker from doris cluster
	for _, name := range evictBkNames {
		if err := DropBroker(db, name); err != nil {
			return NewRecSqlErr(err)
		}
		r.Log.Info(fmt.Sprintf("drop broker[%s] from doris cluster[%s] via connection: %s",
			name, r.CR.ObjKey().String(), sqlConnConf.HostPort()))
	}

	return nil
}
