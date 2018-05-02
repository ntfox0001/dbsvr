//package main

// import (
// 	"fmt"

// 	"github.com/ntfox0001/dbsvr/config"
// 	"github.com/ntfox0001/dbsvr/database"

// 	"github.com/ethereum/go-ethereum/log"
// )

// //"database/sql"
// //"github.com/go-sql-driver/mysql"

// func main() {
// 	if err := config.LoadConfigFile("config.json"); err != nil {
// 		log.Error("dbsvr", err.Error())
// 	}

// 	dbip := config.GetValue("dbip", "").(string)
// 	dbport := config.GetValue("dbport", 0)
// 	dbuser := config.GetValue("user", "").(string)
// 	dbpw := config.GetValue("password", "").(string)
// 	dbName := config.GetValue("database", "").(string)

// 	// if db, err := database.NewDatabase(dbip, fmt.Sprint(dbport), dbuser, dbpw, dbName); err != nil {
// 	// 	log.Error("main", err.Error())
// 	// } else {
// 	// 	op := db.CreateOperation("call BlockchainIPList_insert(?, ?, ?, ?, ?, ?)", "test1", "test1", "test1", "test1", "test1", "test1")

// 	// 	db.ExecOperation(op)
// 	// 	if err != nil {
// 	// 		return
// 	// 	}

// 	// 	for {
// 	// 	}

// 	// }

// 	database.InitialDB(dbip, fmt.Sprint(dbport), dbuser, dbpw, dbName)

// 	database.Query("call BlockchainIPList_insert(?, ?, ?, ?, ?, ?)", "test1", "test1", "test1", "test1", "test1", "test1")
// 	database.Query("call BlockchainIPList_insert(?, ?, ?, ?, ?, ?)", "test2", "test1", "test1", "test1", "test1", "test1")

// 	for {
// 	}
// }
// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/ntfox0001/dbsvr/network"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {

	flag.Parse()

	svr := network.NewServer("127.0.0.1", 8080)

	fr := network.NewRouterFileHandler("/", "home.html")
	svr.RegisterRouter("/", fr)

	wsr := network.NewRouterWSHandler()
	svr.RegisterRouter("/ws", wsr)

	wsr.RegisterJsonMsg("testmsg", func(msg interface{}) {
		fmt.Println(msg)
		json := make(map[string]interface{})
		json["msgId"] = "test1"
		json["value"] = "ffff"
		wsr.SendJsonMsg(json)
	})

	svr.Start()

	for {
	}
}
