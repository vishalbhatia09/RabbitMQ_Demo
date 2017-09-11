package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	_ "math/rand"
	"sync"
	"time"
)

var wg sync.WaitGroup
var db *sql.DB
var stmtIns *sql.Stmt

const (
	//insertRecord = `insert into cases(Case_id, Year_id, State_id, County_id, Hospital_id, Disease_id, Num_of_cases, Num_of_Deaths) Values($1,$2,$3,$4,$5,$6,$7,$8)`
	insertRecord = `insert into cases(Case_id, Year_id, State_id, County_id, Hospital_id, Disease_id, Num_of_cases, Num_of_Deaths) Values(?,?,?,?,?,?,?,?)`
)

const MAX_JOBS = 100

type sqldata struct {
	Case_id       int
	Year_id       int
	State_id      int
	County_id     int
	Hospital_id   int
	Disease_id    int
	Num_of_cases  int
	Num_of_Deaths int
}

var (
	jobs = make(chan sqldata, MAX_JOBS)
)

func main() {

	db, err := sql.Open("mysql", "root:password@/test")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Set the max number of connections
	db.SetMaxOpenConns(MAX_JOBS + 10)
	db.SetMaxIdleConns(0)
	db.SetConnMaxLifetime(1 * time.Nanosecond)

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Connection Successful!!")

	// Spin up workers up to the size of the queue
	for i := 0; i < cap(jobs); i++ {
		// Add this worker to the worker count
		wg.Add(1)

		// This is the goroutine that'll create them in parallel
		go func() {
			for {
				// This will pick a job from the worker pool
				// and it'll also check whether no more jobs were
				// sent back here
				info, ok := <-jobs

				// If no more jobs were sent, then mark the job as complete
				// else just grab more jobs to perform
				if !ok {
					wg.Done()
					return
				}

				// This is the actual task, which can be a function or just a sql statement
				// like here

				_, err := db.Exec(insertRecord, info.Case_id, info.County_id, info.Disease_id, info.Hospital_id, info.Num_of_Deaths, info.Num_of_cases, info.State_id, info.Year_id)
				//_, err := db.Exec(insertRecord, 245678, 345678, 5678, 123456, 345678, 9876, 10890, 20050)
				if err != nil {
					log.Println(err.Error())
				}
			}
		}()
	}

	// Create new jobs for the workers
	for i := 0; i < 1000000; i++ {
		jobs <- sqldata{Case_id: i + 200, Hospital_id: i + 300, Disease_id: i + 400, County_id: i + 250, State_id: 5678, Num_of_Deaths: 890, Num_of_cases: 789, Year_id: 2345}
	}

	// Since all the tasks are going to be sent there, now
	// we just close the channel so the worker can stop
	// doing what it's doing: the infinite for loop
	close(jobs)

	// Then wait until all workers finish
	wg.Wait()

}
