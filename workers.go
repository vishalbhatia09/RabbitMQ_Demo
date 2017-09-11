package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math/rand"
	"sync"
	"time"
)

var wg sync.WaitGroup
var db *sql.DB
var stmtIns *sql.Stmt

const (
	insertRecord = `insert into tblcaseinventory(fldCaseNumber, fldOwnerName, fldBusinessName, fldSiteAddressOne, fldSiteAddresstwo,fldSiteCity,fldSiteState,fldSiteZipCode,fldMailAddress,	fldMailCity,fldMailState,fldMailZip,fldPhoneNumber,fldStartDate,fldClientCityID,fldInputDate,fldSourceID,fldAuthorizedUser,fldOriginalUser,fldExternalSourceInfo,fldTaxPayerID,fldBusinessLicNumber,fldResolutionCode,fldExportDate,fldCityReportBatchID,fldCARReportDate,fldDiscMail_1BatchID,fldNotice1Date,fldDiscMail_2BatchID,fldNotice2Date,fldDiscMail_3BatchID,fldNotice3Date,fldBillingNumber,fldNotice4Date,fldNotice5Date,fldBilling1Date,fldBilling2Date,fldBilling3Date,fldBilling4Date,fldBilling5Date,fldBilling6Date,fldBilling7Date,fldAmount,fldEmailID) VALUES( ?, ?, ?,?,?,?,?,?,?,?,?,?, ?, ?,?,?,?,?,?,?,?,?,?, ?, ?,?,?,?,?,?,?,?,?,?, ?, ?,?,?,?,?,?,?,?,?)`
)

const MAX_JOBS = 100

type sqldata struct {
	fldCaseNumber         int
	fldOwnerName          string
	fldBusinessName       string
	fldSiteAddressOne     string
	fldSiteAddressTwo     string
	fldSiteCity           string
	fldSiteState          string
	fldSiteZipCode        int
	fldMailAddress        string
	fldMailCity           string
	fldMailState          string
	fldMailZip            int
	fldPhoneNumber        string
	fldStartDate          string
	fldClientCityID       int
	fldInputDate          string
	fldSourceID           int
	fldAuthorizedUser     string
	fldOriginalUser       string
	fldExternalSourceInfo string
	fldTaxPayerID         string
	fldBusinessLicNumber  string
	fldResolutionCode     int
	fldExportDate         string
	fldCityReportBatchID  int
	fldCARReportDate      string
	fldDiscMail_1BatchID  int
	fldNotice1Date        string
	fldDiscMail_2BatchID  string
	fldNotice2Date        string
	fldDiscMail_3BatchID  string
	fldNotice3Date        string
	fldBillingNumber      string
	fldNotice4Date        string
	fldNotice5Date        string
	fldBilling1Date       string
	fldBilling2Date       string
	fldBilling3Date       string
	fldBilling4Date       string
	fldBilling5Date       string
	fldBilling6Date       string
	fldBilling7Date       string
	fldAmount             string
	fldEmailID            string
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
				_, err := db.Exec(insertRecord, info.fldCaseNumber, info.fldOwnerName, info.fldBusinessName, info.fldSiteAddressOne, info.fldSiteAddressTwo,
					info.fldSiteCity,
					info.fldSiteState,
					info.fldSiteZipCode,
					info.fldMailAddress,
					info.fldMailCity,
					info.fldMailState,
					info.fldMailZip,
					info.fldPhoneNumber,
					info.fldStartDate,
					info.fldClientCityID,
					info.fldInputDate,
					info.fldSourceID,
					info.fldAuthorizedUser,
					info.fldOriginalUser,
					info.fldExternalSourceInfo,
					info.fldTaxPayerID,
					info.fldBusinessLicNumber,
					info.fldResolutionCode,
					info.fldExportDate,
					info.fldCityReportBatchID,
					info.fldCARReportDate,
					info.fldDiscMail_1BatchID,
					info.fldNotice1Date,
					info.fldDiscMail_2BatchID,
					info.fldNotice2Date,
					info.fldDiscMail_3BatchID,
					info.fldNotice3Date,
					info.fldBillingNumber,
					info.fldNotice4Date,
					info.fldNotice5Date,
					info.fldBilling1Date,
					info.fldBilling2Date,
					info.fldBilling3Date,
					info.fldBilling4Date,
					info.fldBilling5Date,
					info.fldBilling6Date,
					info.fldBilling7Date,
					info.fldAmount,
					info.fldEmailID)
				if err != nil {
					log.Println(err.Error())
				}
			}
		}()
	}

	// Create new jobs for the workers
	for i := 0; i < 1000000; i++ {
		jobs <- sqldata{
			fldCaseNumber:         i,
			fldOwnerName:          "owner name",
			fldBusinessName:       "business name",
			fldSiteAddressOne:     "site address one",
			fldSiteAddressTwo:     "site address two",
			fldSiteCity:           "site city",
			fldSiteState:          "site state",
			fldSiteZipCode:        rand.Intn(250000),
			fldMailAddress:        "mail address",
			fldMailCity:           "mail city",
			fldMailState:          "mail state",
			fldMailZip:            rand.Intn(45000),
			fldPhoneNumber:        "phone number",
			fldStartDate:          "start date",
			fldClientCityID:       i + 200,
			fldInputDate:          "abc",
			fldSourceID:           i + 300,
			fldAuthorizedUser:     "authorized user",
			fldOriginalUser:       "original user",
			fldExternalSourceInfo: "external source info",
			fldTaxPayerID:         "tax payer id",
			fldBusinessLicNumber:  "business lic number",
			fldResolutionCode:     i + 10,
			fldExportDate:         "export date",
			fldDiscMail_1BatchID:  i + 2,
			fldCARReportDate:      "car report date",
			fldCityReportBatchID:  1,
			fldNotice1Date:        "wer",
			fldDiscMail_2BatchID:  "sdf",
			fldNotice2Date:        "notice 2 date",
			fldDiscMail_3BatchID:  "rand.Intn(60000)",
			fldNotice3Date:        "notice 3 date",
			fldBillingNumber:      "billing number",
			fldNotice4Date:        "notice 4 date",
			fldNotice5Date:        "notice 5 date",
			fldBilling1Date:       "billing 1 date",
			fldBilling2Date:       "billing 2 date",
			fldBilling3Date:       "billing 3 date",
			fldBilling4Date:       "billing 4 date",
			fldBilling5Date:       "billing 5 date",
			fldBilling6Date:       "billing 6 date",
			fldBilling7Date:       "billing 7 date",
			fldAmount:             "amount",
			fldEmailID:            "emailid",
		}
	}

	// Since all the tasks are going to be sent there, now
	// we just close the channel so the worker can stop
	// doing what it's doing: the infinite for loop
	close(jobs)

	// Then wait until all workers finish
	wg.Wait()

}
