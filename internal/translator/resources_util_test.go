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

package translator

import "testing"

func TestDumpJavaBasedComponentConf(t *testing.T) {
	eval := func(configs map[string]string, expected string) {
		result := dumpJavaBasedComponentConf(configs)
		if result != expected {
			t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, result)
		}
	}

	// normal configs
	eval(
		map[string]string{
			"http_port":     "8030",
			"rpc_port":      "9020",
			"sys_log_level": "INFO",
		},
		`http_port=8030
rpc_port=9020
sys_log_level=INFO
JAVA_OPTS="-XX:MaxRAMPercentage=75 -XX:InitialRAMPercentage=75 -XX:MinRAMPercentage=75"`)

	// with java opts
	eval(
		map[string]string{
			"http_port":           "8030",
			"JAVA_OPTS":           "-Djavax.security.auth.useSubjectCredsOnly=false -Xss4m -Xmx8192m -XX:+UseMembar -XX:SurvivorRatio=8",
			"JAVA_OPTS_FOR_JDK_9": "-Djavax.security.auth.useSubjectCredsOnly=false -Xss4m -Xmx8192m -XX:SurvivorRatio=8",
		},
		`JAVA_OPTS="-Djavax.security.auth.useSubjectCredsOnly=false -XX:+UseMembar -XX:SurvivorRatio=8 -XX:MaxRAMPercentage=75 -XX:InitialRAMPercentage=75 -XX:MinRAMPercentage=75"
JAVA_OPTS_FOR_JDK_9="-Djavax.security.auth.useSubjectCredsOnly=false -XX:SurvivorRatio=8 -XX:MaxRAMPercentage=75 -XX:InitialRAMPercentage=75 -XX:MinRAMPercentage=75"
http_port=8030`)

	// empty configs
	eval(
		map[string]string{},
		`JAVA_OPTS="-XX:MaxRAMPercentage=75 -XX:InitialRAMPercentage=75 -XX:MinRAMPercentage=75"`)

}

func TestDumpCppBasedComponentConf(t *testing.T) {
	eval := func(configs map[string]string, expected string) {
		result := dumpCppBasedComponentConf(configs)
		if result != expected {
			t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, result)
		}
	}
	// normal config
	eval(
		map[string]string{
			"be_port":        "9060",
			"webserver_port": "8040",
			"enable_https":   "false",
		},
		`be_port=9060
enable_https=false
webserver_port=8040`)

	// with java opts
	eval(
		map[string]string{
			"be_port":             "9060",
			"JAVA_OPTS":           "-Djavax.security.auth.useSubjectCredsOnly=false -Xss4m -Xmx8192m -XX:+UseMembar -XX:SurvivorRatio=8",
			"JAVA_OPTS_FOR_JDK_9": "-Djavax.security.auth.useSubjectCredsOnly=false -Xss4m -Xmx8192m -XX:SurvivorRatio=8",
		},
		`JAVA_OPTS="-Djavax.security.auth.useSubjectCredsOnly=false -Xss4m -Xmx8192m -XX:+UseMembar -XX:SurvivorRatio=8"
JAVA_OPTS_FOR_JDK_9="-Djavax.security.auth.useSubjectCredsOnly=false -Xss4m -Xmx8192m -XX:SurvivorRatio=8"
be_port=9060`)

}
