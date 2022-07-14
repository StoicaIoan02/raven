package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const (
	// Initialize connection constants.
	HOST     = "kdgpostgresql.postgres.database.azure.com"
	DATABASE = "raven_beta_ioan"
	USER     = "kdgdbadmin@kdgpostgresql"
	PASSWORD = "z6@U3ns!@RpP"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

type StartingJob struct {
	id                       int
	declaration_id           *int
	date                     time.Time
	status                   string
	tip                      string
	job_id                   *int
	reporting_declaration_id *int
	account_id               int
}

func checkJobs(tipCautat string) {

	// Initialize connection string.
	var connectionString string = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=require", HOST, USER, PASSWORD, DATABASE)

	// Initialize connection object.
	db, err := sql.Open("postgres", connectionString)
	checkError(err)

	// Ping
	err = db.Ping()
	checkError(err)
	fmt.Println("Successfully created connection to database")

	// Selectam prima declaratie in progres.
	sql_statement := fmt.Sprintf("SELECT id from declaration_queue  where status = 'on progress' and type = '%s' limit 1;", tipCautat)
	rows, err := db.Query(sql_statement)
	checkError(err)
	defer rows.Close()

	// Daca nu avem nici-o declaratie in progres
	if !rows.Next() {
		fmt.Println("Nu avem declaratii in progres")

		// Selectam prima declaratie din coada
		sql_statement := fmt.Sprintf("SELECT * from declaration_queue  where status = 'queued' and type = '%s' limit 1;", tipCautat)
		rows, err := db.Query(sql_statement)
		checkError(err)
		defer rows.Close()

		// Daca avem cel putin o declaratie in coada
		if rows.Next() {
			// Prima declaratie din coada
			var declaratie StartingJob
			err = rows.Scan(&declaratie.id, &declaratie.declaration_id, &declaratie.date, &declaratie.status, &declaratie.tip, &declaratie.job_id, &declaratie.reporting_declaration_id, &declaratie.account_id)
			checkError(err)
			fmt.Println("Data row = (", declaratie.id, ", ", declaratie.status, ")")

			/*var url_pad sql.NullString
			//var lrn sql.NullString

			if tipCautat == "report" {
				///Setare pad_url
				sql_statement := fmt.Sprintf("select pad_url from users where id = (SELECT users_id _id from reporting_declaration where id = %d );", *declaratie.reporting_declaration_id)
				rows, err := db.Query(sql_statement)
				checkError(err)
				defer rows.Close()

				rows.Next()
				err = rows.Scan(&url_pad)
				checkError(err)

				if url_pad.Valid == false {
					fmt.Println("error: url_pad.Valid == false ")
				}
				fmt.Println("Pad_url:", url_pad.String)

			}*/

		} else {
			fmt.Println("Nu avem declaratii in coada")
		}
	} else {
		fmt.Println("Exista deja o declaratie in progres")
	}
}

func main() {
	checkJobs("report")
	println("Program rulat cu succes")
}
