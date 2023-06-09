package pkg

import (
	"time"

	ck "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

func getConn(addr string) (driver.Conn, error) {
	options := &ck.Options{
		Addr: []string{addr},
		Auth: ck.Auth{
			//Database: databaseName,
			Username: "",
			Password: "",
		},
		DialTimeout:     getDurationEnv("DIAL_TIME_OUT", 10*time.Second),
		Debug:           getBoolEnv("DEBUG", false),
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
