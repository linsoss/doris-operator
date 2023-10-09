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
	"database/sql"
	"errors"
	"fmt"
	ut "github.com/al-assad/doris-operator/internal/util"
	u "github.com/rjNemo/underscore"
)

func ShowFrontendHosts(db *sql.DB) ([]string, error) {
	rows, err := db.Query("show frontends")
	defer rows.Close()

	if err != nil {
		return []string{}, ut.MergeErrors(errors.New("failed to execute sql 'show frontends'"), err)
	}
	rowSet := ReadAllRowsAsString(rows)
	hosts := u.Map(rowSet, func(row RowMap) string {
		return row["Host"]
	})
	return hosts, nil
}

func ShowBackendHosts(db *sql.DB) ([]string, error) {
	rows, err := db.Query("show backends")
	defer rows.Close()

	if err != nil {
		return []string{}, ut.MergeErrors(errors.New("failed to execute sql 'show backends'"), err)
	}
	rowSet := ReadAllRowsAsString(rows)
	hosts := u.Map(rowSet, func(row RowMap) string {
		return row["Host"]
	})
	return hosts, nil
}

// ShowBrokerNameHosts returns map structure: key is broker name, value is broker host
func ShowBrokerNameHosts(db *sql.DB) (map[string]string, error) {
	rows, err := db.Query("show broker")
	defer rows.Close()

	if err != nil {
		return map[string]string{}, ut.MergeErrors(errors.New("failed to execute sql 'show broker'"), err)
	}
	rowSet := ReadAllRowsAsString(rows)
	nameHosts := make(map[string]string)
	for _, row := range rowSet {
		nameHosts[row["Name"]] = row["Host"]
	}
	return nameHosts, nil
}

func AddFrontend(db *sql.DB, feHostPort string) error {
	addSql := fmt.Sprintf(`alter system add follower "%s"`, feHostPort)
	_, err := db.Exec(addSql)
	if err != nil {
		return ut.MergeErrors(errors.New(fmt.Sprintf("failed to execute sql '%s'", addSql)), err)
	}
	return nil
}

func AddBackend(db *sql.DB, beHostPort string) error {
	addSql := fmt.Sprintf(`alter system add backend "%s"`, beHostPort)
	_, err := db.Exec(addSql)
	if err != nil {
		return ut.MergeErrors(errors.New(fmt.Sprintf("failed to execute sql '%s'", addSql)), err)
	}
	return nil
}

func AddBroker(db *sql.DB, brokerName string, brokerHost string) error {
	addSql := fmt.Sprintf(`alter system add broker %s "%s"`, brokerName, brokerHost)
	_, err := db.Exec(addSql)
	if err != nil {
		return ut.MergeErrors(errors.New(fmt.Sprintf("failed to execute sql '%s'", addSql)), err)
	}
	return nil
}

func DropFrontend(db *sql.DB, feHostPort string) error {
	addSql := fmt.Sprintf(`alter system drop follower "%s"`, feHostPort)
	_, err := db.Exec(addSql)
	if err != nil {
		return ut.MergeErrors(errors.New(fmt.Sprintf("failed to execute sql '%s'", addSql)), err)
	}
	return nil
}

func DropBackend(db *sql.DB, beHostPort string) error {
	addSql := fmt.Sprintf(`alter system drop backend "%s"`, beHostPort)
	_, err := db.Exec(addSql)
	if err != nil {
		return ut.MergeErrors(errors.New(fmt.Sprintf("failed to execute sql '%s'", addSql)), err)
	}
	return nil
}

func DropBroker(db *sql.DB, brokerName string) error {
	addSql := fmt.Sprintf(`alter system all broker %s`, brokerName)
	_, err := db.Exec(addSql)
	if err != nil {
		return ut.MergeErrors(errors.New(fmt.Sprintf("failed to execute sql '%s'", addSql)), err)
	}
	return nil
}
