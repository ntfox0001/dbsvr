package database

import (
	log "github.com/inconshreveable/log15"
)

var (
	mDatabase     *Database
	mDataResultCh chan (<-chan *DataResult)
	mQuitCh       chan interface{}
)

func InitialDB(ip, port, user, password, database string) error {
	var err error = nil
	mDatabase, err = NewDatabase(ip, port, user, password, database)

	mDataResultCh = make(chan (<-chan *DataResult))
	if err != nil {
		log.Error("database", err.Error())
		return err
	}

	go run()

	return nil
}

func run() {
running:
	for {
		select {
		case drCh := <-mDataResultCh:
			{
				go func() {
					for {
						select {
						case dr := <-drCh:
							{
								if dr.result == nil && dr.err == nil {
									return
								}
							}
						}
					}
				}()
			}
		case <-mQuitCh:
			break running
		}
	}
}

func Query(sql string, args ...interface{}) error {

	op := mDatabase.CreateOperation(sql, args...)
	drCh, err := mDatabase.ExecOperation(op)

	mDataResultCh <- drCh

	return err
}

func Close() {
	mDatabase.Close()
	mQuitCh <- struct{}{}
}
