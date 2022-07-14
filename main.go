package main

import (
	"database/sql"
	"fmt"

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

	// Numarare declaratii in progres.
	sql_statement := fmt.Sprintf("SELECT count(id)  from declaration_queue  where status = 'on progress' and type = '%s';", tipCautat)
	rows, err := db.Query(sql_statement)
	checkError(err)
	defer rows.Close()

	var runningJobs int
	rows.Next()

	switch err := rows.Scan(&runningJobs); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned")
	case nil:
		fmt.Println("Numar declaratii in runningJobs:", runningJobs)
	default:
		checkError(err)
	}

	if runningJobs == 0 {
		fmt.Println("Nu avem declaratii in progres.")

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

			// startingJob
			var (
				id                       sql.NullString
				declaration_id           sql.NullString
				date                     sql.NullString
				status                   sql.NullString
				tip                      sql.NullString
				job_id                   sql.NullString
				reporting_declaration_id sql.NullString
				account_id               sql.NullString
			)

			// Prima declaratie din coada
			fmt.Println("Prima declaratie: ")
			rows.Next()
			switch err := rows.Scan(&id, &declaration_id, &date, &status, &tip, &job_id, &reporting_declaration_id, &account_id); err {
			case sql.ErrNoRows:
				fmt.Println("No rows were returned")
			case nil:
				fmt.Println("Data row = (", id.String, ", ", status.String, ")")
			default:
				checkError(err)
			}

			var url_pad sql.NullString
			//lrn := ""

			if tipCautat == "report" {
				sql_statement := fmt.Sprintf("select pad_url from users where id = (SELECT users_id _id from reporting_declaration where id = %s );", reporting_declaration_id.String)
				rows, err := db.Query(sql_statement)
				checkError(err)
				defer rows.Close()

				rows.Next()
				err = rows.Scan(&url_pad)
				checkError(err)

				fmt.Println("Pad_url:", url_pad.String)

			}
		}
	}
}

func main() {
	checkJobs("report")
	println("Program rulat cu succes")
}
