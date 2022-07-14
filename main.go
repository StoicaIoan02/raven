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

	// Read rows from table; tipul variabilelor este sql.NullString deoarece putem gasii NULL in orice camp
	var id sql.NullString
	var declaration_id sql.NullString
	var date sql.NullString
	var status sql.NullString
	var tip sql.NullString
	var job_id sql.NullString
	var reporting_declaration_id sql.NullString
	var account_id sql.NullString

	err = db.Ping()
	checkError(err)
	fmt.Println("Successfully created connection to database")

	sql_statement := fmt.Sprintf("SELECT * from declaration_queue  where status = 'on progress' and type = '%s';", tipCautat)
	rows, err := db.Query(sql_statement)
	checkError(err)
	defer rows.Close()

	for rows.Next() {
		switch err := rows.Scan(&id, &declaration_id, &date, &status, &tip, &job_id, &reporting_declaration_id, &account_id); err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned")
		case nil:
			fmt.Println("Data row = (", id, ", ", status, ")\n")
		default:
			checkError(err)
		}
	}
}

func main() {
	checkJobs("report")
	println("final functie")
}
