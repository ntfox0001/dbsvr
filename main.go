package main

import (
	"fmt"

	"github.com/ntfox0001/dbsvr/config"
	"github.com/ntfox0001/dbsvr/database"

	"github.com/ethereum/go-ethereum/log"
)

//"database/sql"
//"github.com/go-sql-driver/mysql"

func main() {
	if err := config.LoadConfigFile("config.json"); err != nil {
		log.Error("dbsvr", err.Error())
	}

	dbip := config.GetValue("dbip", "").(string)
	dbport := config.GetValue("dbport", 0)
	dbuser := config.GetValue("user", "").(string)
	dbpw := config.GetValue("password", "").(string)
	dbName := config.GetValue("database", "").(string)

	// if db, err := database.NewDatabase(dbip, fmt.Sprint(dbport), dbuser, dbpw, dbName); err != nil {
	// 	log.Error("main", err.Error())
	// } else {
	// 	op := db.CreateOperation("call BlockchainIPList_insert(?, ?, ?, ?, ?, ?)", "test1", "test1", "test1", "test1", "test1", "test1")

	// 	db.ExecOperation(op)
	// 	if err != nil {
	// 		return
	// 	}

	// 	for {
	// 	}

	// }

	database.InitialDB(dbip, fmt.Sprint(dbport), dbuser, dbpw, dbName)

	database.Query("call BlockchainIPList_insert(?, ?, ?, ?, ?, ?)", "test1", "test1", "test1", "test1", "test1", "test1")
	database.Query("call BlockchainIPList_insert(?, ?, ?, ?, ?, ?)", "test2", "test1", "test1", "test1", "test1", "test1")

	for {
	}
}
