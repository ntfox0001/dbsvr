package database

import (
	"blockchainproj/dbsvr/config"
	"blockchainproj/dbsvr/dberror"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/inconshreveable/log15"
)

type Database struct {
	ip       string
	port     string
	user     string
	password string
	database string

	dataOptCh    chan *DataOperation
	quitCh       chan struct{}
	operationMap map[int32]*DataOperation
	idBase       int32

	sqldb *sql.DB
}

func NewDatabase(ip, port, user, password, database string) (*Database, error) {
	opCnt := config.GetValue("DatabaseOpCnt", 10).(int)
	dbConnStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4", user, password, ip, port, database)
	if sqldb, err := sql.Open("mysql", dbConnStr); err == nil {

		db := &Database{
			ip:           ip,
			port:         port,
			user:         user,
			password:     password,
			database:     database,
			dataOptCh:    make(chan *DataOperation, opCnt),
			quitCh:       make(chan struct{}),
			operationMap: make(map[int32]*DataOperation),
			idBase:       0,
			sqldb:        sqldb,
		}
		go db.run()

		return db, nil
	} else {
		log.Error("database", "failed to connect database:", ip)
	}

	return nil, dberror.NewStringErr("failed to create database.")
}

func (d *Database) CreateOperation(sql string, args ...interface{}) *DataOperation {
	op := NewOperation(d.idBase, sql, args...)
	d.addOpt(op)
	return op
}

func (d *Database) addOpt(op *DataOperation) {
	d.idBase++
	d.operationMap[op.id] = op
	return
}

func (d *Database) ExecOperation(op *DataOperation) (<-chan *DataResult, error) {
	if !op.executing {
		op.executing = true
		go func() {

			d.dataOptCh <- op
		}()
		return op.resultOptCh, nil
	}
	log.Error("database", "operation is executing.")
	return nil, dberror.NewStringErr("operation is executing.")
}

func (d *Database) Close() {
	d.quitCh <- struct{}{}
}

func (d *Database) run() {
running:
	for {
		select {
		case op := <-d.dataOptCh:
			{
				d.exec(op)
			}
		case <-d.quitCh:
			break running
		}
	}
	return
}

func (d *Database) exec(op *DataOperation) {
	op.exec(d.sqldb)
}
