package clickhouse

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/cheggaaa/pb/v3"
)

type Batch struct {
	driver.Batch
	totalRows   int // Total number of rows in the batch
	progressBar *pb.ProgressBar
}

func Prepare(conn driver.Conn, databaseName, tableName string) (*Batch, error) {
	batch, err := conn.PrepareBatch(context.Background(), fmt.Sprintf("INSERT INTO %s.%s", databaseName, tableName))
	if err != nil {
		return nil, err
	}
	b := &Batch{}
	b.Batch = batch
	return b, nil
}
func (b *Batch) NewProgressBar(size int) {
	b.progressBar = pb.StartNew(size)
}

func (b *Batch) Increment() {
	if b.progressBar != nil {
		b.progressBar.Increment()
	}
}
func (b *Batch) Finish() {
	if b.progressBar != nil {
		b.progressBar.Finish()
	}
}

// AppendStruct appends a struct to the batch and updates the total rows count
func (b *Batch) AppendStruct(s interface{}) error {
	err := b.Batch.AppendStruct(s)
	if err == nil {
		//b.Increment()
		b.totalRows++
	}
	return err
}

// TotalRows returns the total number of rows in the batch
func (b *Batch) TotalRows() int {
	return b.totalRows
}

// Send sends the batch for execution and resets the total rows count
func (b *Batch) Send() error {
	if b.Batch.IsSent() {
		return nil
	}
	err := b.Batch.Send()
	if err == nil {
		b.totalRows = 0
	}
	return err
}
