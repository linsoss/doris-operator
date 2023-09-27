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

package transformer

import (
	"crypto/rand"
	"fmt"
	dapi "github.com/al-assad/doris-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"math/big"
	"sort"
	"strconv"
	"strings"
)

const (
	K8sNameLabelKey      = "app.kubernetes.io/name"
	K8sInstanceLabelKey  = "app.kubernetes.io/instance"
	K8sManagedByLabelKey = "app.kubernetes.io/managed-by"
	K8sComponentLabelKey = "app.kubernetes.io/component"

	DorisK8sNameLabelValue      = "doris-cluster"
	DorisK8sManagedByLabelValue = "doris-operator"

	PrometheusPathAnnoKey   = "prometheus.io/path"
	PrometheusPortAnnoKey   = "prometheus.io/port"
	PrometheusScrapeAnnoKey = "prometheus.io/scrape"
)

// MakeResourceLabels make the k8s label meta for the managed resource
func MakeResourceLabels(dorisName string, component string) map[string]string {
	labels := map[string]string{
		K8sNameLabelKey:      DorisK8sNameLabelValue,
		K8sManagedByLabelKey: DorisK8sManagedByLabelValue,
		K8sInstanceLabelKey:  dorisName,
		K8sComponentLabelKey: component,
	}
	return labels
}

const DorisPasswordChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-=_+[]{}"

func GenerateRandomDorisPassword(length int) string {
	password := make([]byte, length)
	maxIndex := big.NewInt(int64(len(DorisPasswordChars)))

	for i := range password {
		randomIndex, _ := rand.Int(rand.Reader, maxIndex)
		password[i] = DorisPasswordChars[randomIndex.Int64()]
	}
	return string(password)
}

const (
	JvmOptKey        = "JAVA_OPTS"
	JvmOpt9Key       = "JAVA_OPTS_FOR_JDK_9"
	JvmRamPercentage = 75
)

// Dump the doris component(FE, Broker) KV configs into plain text
func dumpJavaBasedComponentConf(config map[string]string) string {
	containerJvmRamOpt := fmt.Sprintf(
		"-XX:MaxRAMPercentage=%d -XX:InitialRAMPercentage=%d -XX:MinRAMPercentage=%d",
		JvmRamPercentage, JvmRamPercentage, JvmRamPercentage)
	// order by key
	var keys []string
	for key := range config {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// generate config file content
	var lines []string
	hasJvmOpt := false

	for _, k := range keys {
		key := strings.TrimSpace(k)
		value := strings.TrimSpace(config[k])
		// handle JVM opt config
		if key == JvmOptKey {
			hasJvmOpt = true
		}
		if key == JvmOptKey || key == JvmOpt9Key {
			parts := strings.Split(value, " ")
			var filterParts []string
			for _, part := range parts {
				if !strings.HasPrefix(part, "-Xss") && !strings.HasPrefix(part, "-Xmx") {
					filterParts = append(filterParts, part)
				}
			}
			filterParts = append(filterParts, containerJvmRamOpt)
			value = fmt.Sprintf(`"%s"`, strings.Join(filterParts, " "))
		}
		lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	}
	if !hasJvmOpt {
		lines = append(lines, fmt.Sprintf("%s=%s", JvmOptKey, fmt.Sprintf(`"%s"`, containerJvmRamOpt)))
	}
	return strings.Join(lines, "\n")
}

// Dump the doris component(BE, CN) KV configs into plain text
func dumpCppBasedComponentConf(config map[string]string) string {
	// order by key
	var keys []string
	for key := range config {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// generate config file content
	var lines []string
	for _, k := range keys {
		key := strings.TrimSpace(k)
		value := strings.TrimSpace(config[k])
		// handle JVM opt config
		if key == JvmOptKey || key == JvmOpt9Key {
			value = fmt.Sprintf(`"%s"`, value)
		}
		lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(lines, "\n")
}

// Get the port value from the kv config map
func getPortValueFromRawConf(config map[string]string, key string, defaultValue int32) int32 {
	strValue := config[key]
	if strValue == "" {
		return defaultValue
	}
	intValue, err := strconv.ParseInt(strValue, 10, 32)
	if err != nil {
		return defaultValue
	}
	return int32(intValue)
}

// Merge HostAlias info in HostnameIpItem into HostAlias slice
func mergeHostAlias(items []dapi.HostnameIpItem, hostAlias []corev1.HostAlias) []corev1.HostAlias {
	if len(items) == 0 {
		return hostAlias
	}
	var result = make(map[string]*corev1.HostAlias)
	for _, item := range items {
		v, ok := result[item.IP]
		if !ok {
			result[item.IP] = &corev1.HostAlias{
				IP:        item.IP,
				Hostnames: []string{item.Name},
			}
		} else {
			v.Hostnames = append(v.Hostnames, item.Name)
		}
	}
	newHostAlias := make([]corev1.HostAlias, len(result))
	for _, value := range result {
		newHostAlias = append(newHostAlias, *value)
	}
	return newHostAlias
}
