package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "math/rand"
	"sync"
	"time"
)

var wg sync.WaitGroup
var db *sql.DB
var stmtIns *sql.Stmt

func inserting_data(i int, db *sql.DB, stmtIns *sql.Stmt) {
	fmt.Println(i)
	db.SetMaxOpenConns(100)

	for {

		_, err := stmtIns.Exec(i)

		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		wg.Done()
	}

}

func main() {

	db, err := sql.Open("mysql", "root:password@/test")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	//db.SetMaxIdleConns(100)
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	fmt.Println("Connection Successful!!")

	stmtIns, err := db.Prepare("INSERT INTO testtable VALUES(?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	startTime := time.Now()
	// Insert data in table case inventory
	for i := 1; i <= 100; i++ {
		wg.Add(1)
		go inserting_data(i, db, stmtIns)
	}

	wg.Wait()
	timeConsumed := time.Since(startTime)
	fmt.Println(timeConsumed)

}
