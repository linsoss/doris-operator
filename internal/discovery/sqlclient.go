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
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// DorisSqlConnConf is the Doris SQL connection configuration
type DorisSqlConnConf struct {
	Host     string
	Port     int32
	User     string
	Password string
}

type SqlAccount struct {
	User     string
	Password string
}

func (e *DorisSqlConnConf) HostPort() string {
	return fmt.Sprintf("%s:%d", e.Host, e.Port)
}

func (e *DorisSqlConnConf) Connect() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", e.User, e.Password, e.Host, e.Port)
	return sql.Open("mysql", dsn)
}

type RowMap map[string]string

// ReadAllRowsAsString reads all rows from sql.Rows
// and returns a slice of map[string]string
func ReadAllRowsAsString(rows *sql.Rows) []RowMap {
	cols, _ := rows.Columns()
	var rowMap []RowMap

	for rows.Next() {
		columns := make([]any, len(cols))
		columnPointers := make([]any, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}
		_ = rows.Scan(columnPointers...)

		m := make(map[string]string)
		for i, colName := range cols {
			valPointer := columnPointers[i].(*any)
			val := *valPointer
			if val == nil {
				continue
			}
			valByte, ok := val.([]byte)
			if !ok {
				continue
			}
			m[colName] = string(valByte)
		}
		rowMap = append(rowMap, m)
	}
	return rowMap
}
