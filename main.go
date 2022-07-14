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

	// Daca nu avem nici-o declaratie in progres:
	if !rows.Next() {
		fmt.Println("Nu avem declaratii in progres")

		// Declaratii in asteptare.
		sql_statement := fmt.Sprintf("SELECT count(id)  from declaration_queue  where status = 'queued' and type = '%s';", tipCautat)
		rows, err := db.Query(sql_statement)
		checkError(err)
		defer rows.Close()

		var queuedJobs int
		rows.Next()

		switch err := rows.Scan(&queuedJobs); err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned")
		case nil:
			fmt.Println("Numar declaratii in queuedJobs:", queuedJobs)
		default:
			checkError(err)
		}

		if queuedJobs != 0 {
			// Prima declaratie in asteptare va fi analizata
			sql_statement := fmt.Sprintf("SELECT * from declaration_queue  where status = 'queued' and type = '%s' limit 1;", tipCautat)
			rows, err := db.Query(sql_statement)
			checkError(err)
			defer rows.Close()

			var declaratie StartingJob
			// Prima declaratie din coada
			fmt.Println("Prima declaratie: ")
			rows.Next()
			switch err := rows.Scan(&declaratie.id, &declaratie.declaration_id, &declaratie.date, &declaratie.status, &declaratie.tip, &declaratie.job_id, &declaratie.reporting_declaration_id, &declaratie.account_id); err {
			case sql.ErrNoRows:
				fmt.Println("No rows were returned")
			case nil:
				fmt.Println("Data row = (", declaratie.id, ", ", declaratie.status, ")")
			default:
				checkError(err)
			}

			rows, err = db.Query("count(1);")

			fmt.Println(*declaratie.reporting_declaration_id)

			var url_pad sql.NullString
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

			}
		}
	} else {
		fmt.Println("Exista deja o declaratie in progres")
	}
}

func main() {
	checkJobs("report")
	println("Program rulat cu succes")
}
