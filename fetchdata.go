package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/streadway/amqp"
	"log"
	_ "math/rand"
	"sync"
	"time"
)

var wg sync.WaitGroup
var db *sql.DB

const (
	fetchrecords = `select * from cases`
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

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

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

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	var (
		case_id       int
		year_id       int
		state_id      int
		county_id     int
		hospital_id   int
		disease_id    int
		num_of_cases  int
		num_of_Deaths int
	)

	// Spin up workers up to the size of the queue
	for i := 0; i < cap(jobs); i++ {

		fmt.Println("spinning worker:", i)
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

				// This is the actual task, where we are publihing the rows onto the rabbit queue
				// like here
				body, _ := json.Marshal(info)
				err = ch.Publish(
					"",     // exchange
					q.Name, // routing key
					false,  // mandatory
					false,  // immediate
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        body,
					})
				failOnError(err, "Failed to publish a message")

			}
		}()
	}

	rows, err := db.Query(fetchrecords)
	//_, err := db.Exec(insertRecord, 245678, 345678, 5678, 123456, 345678, 9876, 10890, 20050)
	if err != nil {
		log.Println(err.Error())
	}

	for rows.Next() {

		err := rows.Scan(&case_id, &year_id, &state_id, &county_id, &hospital_id, &disease_id, &num_of_cases, &num_of_Deaths)
		if err != nil {
			panic(err.Error())
		}

		//fmt.Println(case_id, year_id, state_id, county_id, hospital_id, disease_id, num_of_cases, num_of_Deaths)
		data := sqldata{Case_id: case_id,
			Year_id:       year_id,
			State_id:      state_id,
			County_id:     county_id,
			Hospital_id:   hospital_id,
			Disease_id:    disease_id,
			Num_of_cases:  num_of_cases,
			Num_of_Deaths: num_of_Deaths,
		}

		fmt.Println(data)
		jobs <- data
		fmt.Println(len(jobs))
	}

}
