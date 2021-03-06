package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"math/rand"
	"sync"
)

var wg sync.WaitGroup
var db *sql.DB
var stmtIns *sql.Stmt

func inserting_data(i int, db *sql.DB, stmtIns *sql.Stmt){
	fmt.Println(i)
	db.SetMaxOpenConns(100)
	
	_, err := stmtIns.Exec(i,"owner name", "business name", "site address one", "site address two", "site city", "site state", rand.Intn(250000), "mail address", "mail city", "mail state", rand.Intn(45000), "phone number", "start date", i+200, time.Date(2017, time.September, 05, 0, 0, 0, 0, time.UTC), i+300, "authorized user", "original user", "external source info", "tax payer id", "business lic number", i +10, "export date", i+2 ,"car report date", rand.Intn(1000), time.Now(), rand.Intn(50), "notice 2 date", rand.Intn(60000), "notice 3 date", "billing number", "notice 4 date", "notice 5 date", "billing 1 date", "billing 2 date", "billing 3 date", "billing 4 date", "billing 5 date", "billing 6 date", "billing 7 date", "amount", "emailid" )
	

	if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
	}
	wg.Done()

}


func main() {

	db, err := sql.Open("mysql", "user:password@/test")
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

	
	stmtIns, err := db.Prepare("INSERT INTO tblCaseInventory VALUES( ?, ?, ?,?,?,?,?,?,?,?,?,?, ?, ?,?,?,?,?,?,?,?,?,?, ?, ?,?,?,?,?,?,?,?,?,?, ?, ?,?,?,?,?,?,?,?,? )") 
	if err != nil {
		panic(err.Error()) 
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates


	startTime := time.Now()
	// Insert data in table case inventory
	for i := 1 ; i <= 1000000 ; i++ {
		wg.Add(1)
		go inserting_data(i, db, stmtIns)		
	}
	
	wg.Wait()
	timeConsumed := time.Since(startTime)
	fmt.Println(timeConsumed)

}
