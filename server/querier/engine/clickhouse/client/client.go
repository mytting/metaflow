/*
 * Copyright (c) 2022 Yunshan Networks
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
 */

package client

import (
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/jmoiron/sqlx"
	//"github.com/k0kubun/pp"
	logging "github.com/op/go-logging"
	"time"
	"unsafe"
)

var log = logging.MustGetLogger("clickhouse.client")

type Client struct {
	Host       string
	Port       int
	UserName   string
	Password   string
	connection *sqlx.DB
	DB         string
	Debug      *Debug
}

func (c *Client) init(query_uuid string) error {
	if c.Debug == nil {
		c.Debug = &Debug{
			QueryUUID: query_uuid,
			IP:        c.Host,
		}
	}
	url := fmt.Sprintf("clickhouse://%s:%s@%s:%d/%s?&query_id=%s", c.UserName, c.Password, c.Host, c.Port, c.DB, query_uuid)
	conn, err := sqlx.Open(
		"clickhouse", url,
	)
	if err != nil {
		log.Errorf("connect clickhouse failed: %s, url: %s, query_uuid: %s", err, url, query_uuid)
		return err
	}
	c.connection = conn
	return nil
}

func (c *Client) Close() error {
	return c.connection.Close()
}

func (c *Client) DoQuery(sql string, callbacks []func(columns []interface{}, values []interface{}) []interface{}, query_uuid string) (map[string][]interface{}, error) {
	err := c.init(query_uuid)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	rows, err := c.connection.Queryx(sql)
	c.Debug.Sql = sql
	if err != nil {
		log.Errorf("query clickhouse Error: %s, sql: %s, query_uuid: %s", err, sql, c.Debug.QueryUUID)
		c.Debug.Error = fmt.Sprintf("%s", err)
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.ColumnTypes()
	resColumns := len(columns)
	if err != nil {
		c.Debug.Error = fmt.Sprintf("%s", err)
		return nil, err
	}
	result := make(map[string][]interface{})
	var columnNames []interface{}
	var columnTypes []string
	// 获取列名和列类型
	for _, column := range columns {
		columnNames = append(columnNames, column.Name())
		columnTypes = append(columnTypes, column.DatabaseTypeName())
	}
	result["columns"] = columnNames
	var values []interface{}
	resSize := 0
	start := time.Now()
	for rows.Next() {
		row, err := rows.SliceScan()
		if err != nil {
			c.Debug.Error = fmt.Sprintf("%s", err)
			return nil, err
		}
		var record []interface{}
		for i, rawValue := range row {
			value, err := TransType(columnTypes[i], rawValue)
			if err != nil {
				c.Debug.Error = fmt.Sprintf("%s", err)
				return nil, err
			}
			resSize += int(unsafe.Sizeof(value))
			record = append(record, value)
		}
		values = append(values, record)
	}
	resRows := len(values)
	queryTime := time.Since(start)
	c.Debug.QueryTime = int64(queryTime)
	for _, callback := range callbacks {
		values = callback(columnNames, values)
	}
	result["values"] = values
	log.Debugf("sql: %s, query_uuid: %s", sql, c.Debug.QueryUUID)
	log.Debugf("res_rows: %v, res_columns: %v, res_size: %v", resRows, resColumns, resSize)
	return result, nil
}
