package database

import (
	"blockchainproj/dbsvr/dberror"
	"container/list"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/inconshreveable/log15"
)

type DataOperation struct {
	id          int32
	sql         string
	args        []interface{}
	usePrepare  bool
	dataSet     *list.List
	resultOptCh chan *DataResult
	errorOptCh  chan string
	executing   bool
	transaction bool
}

type DataResult struct {
	opt    *DataOperation
	result *list.List
	err    error
}

type dbOperation interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

func NewOperation(id int32, sql string, args ...interface{}) *DataOperation {
	op := &DataOperation{
		id:          id,
		sql:         sql,
		args:        args,
		usePrepare:  false,
		dataSet:     nil,
		resultOptCh: make(chan *DataResult),
		errorOptCh:  make(chan string, 10),
		executing:   false,
		transaction: false,
	}

	return op
}

func (d *DataOperation) exec(db *sql.DB) error {
	if d.transaction {
		if tx, err := db.Begin(); err == nil {

			d.callData(tx)

			// 发送结束标识
			d.resultOptCh <- &DataResult{
				opt:    d,
				result: nil,
				err:    nil,
			}
			if err := tx.Commit(); err != nil {
				log.Error("database", "Failed to commit:", d.sql)
			}

			return nil
		} else {
			log.Error("database", "Failed to begin transaction:", d.sql)
		}

	} else {
		d.callData(db)
		// 发送结束标识
		d.resultOptCh <- &DataResult{
			opt:    d,
			result: nil,
			err:    nil,
		}
	}

	return dberror.NewStringErr("failed to create database.")
}

func (d *DataOperation) callData(opt dbOperation) {
	fmt.Println(d.args)
	if row, err := opt.Query(d.sql, d.args...); err != nil {
		d.resultOptCh <- &DataResult{
			opt:    d,
			result: nil,
			err:    err,
		}
	} else {
		//返回所有列
		cols, _ := row.Columns()
		//这里表示一行所有列的值，用[]byte表示
		vals := make([][]byte, len(cols))
		//这里表示一行填充数据
		scans := make([]interface{}, len(cols))
		//这里scans引用vals，把数据填充到[]byte里
		for k, _ := range vals {
			scans[k] = &vals[k]
		}

		i := 0
		result := list.New()
		for row.Next() {
			//填充数据
			row.Scan(scans...)
			//每行数据
			rowData := make(map[string]string, 10)
			//把vals中的数据复制到row中
			for k, v := range vals {
				key := cols[k]
				//这里把[]byte数据转成string
				rowData[key] = string(v)
			}
			//放入结果集
			result.PushBack(rowData)
			i++
		}

		if i > 0 {
			// 发送结果
			d.resultOptCh <- &DataResult{
				opt:    d,
				result: result,
				err:    nil,
			}

		}
	}
}

func (d *DataOperation) SetOperationData(data []interface{}) error {
	if d.executing {
		return dberror.NewStringErr("operation has executed.")
	}
	d.usePrepare = true
	d.dataSet.PushBack(data)

	return nil
}

func (d *DataOperation) SetUsePrepare(use bool) error {
	if d.executing {
		return dberror.NewStringErr("operation has executed.")
	}
	d.usePrepare = use
	return nil
}

func (d *DataOperation) SetOperationTransaction(ts bool) error {
	if d.executing {
		return dberror.NewStringErr("operation has executed.")
	}
	d.transaction = ts
	return nil
}
