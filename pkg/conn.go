//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package pkg

import (
	"os"
	"strings"
	"time"

	ck "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

func getConn(addr string) (driver.Conn, error) {
	options := &ck.Options{
		Addr: strings.Split(addr, ","),
		Auth: ck.Auth{
			Username: os.Getenv("CLICKHOUSE_USER"),
			Password: os.Getenv("CLICKHOUSE_PASSWORD"),
		},
		DialTimeout:     getDurationEnv("DIAL_TIME_OUT", 10*time.Second),
		Debug:           debugFlag,
		MaxIdleConns:    getIntEnv("MAX_IDLE_CONNS", 5),
		MaxOpenConns:    getIntEnv("MAX_OPEN_CONNS", 10),
		ConnMaxLifetime: getDurationEnv("CONN_MAX_LIFE_TIME", 1*time.Hour),
	}

	conn, err := ck.Open(options)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
