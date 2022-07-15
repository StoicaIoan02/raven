package start

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"

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

type Account struct {
	userName string
	password string
	office   string
}

type Declaration_queue struct {
	id                       int
	declaration_id           *int
	reporting_declaration_id *int
	curent_declaration_id    int // == valoare nenula dintre declaraiton_id si reporting_declaration_id
	tip                      string
	declaration_table        string
}

func checkJobs(tipCautat string) { // Primeste ca parametrii "export" sau "import"
	// Verificare parametrii
	if (tipCautat != "export") && ("report" != tipCautat) {
		panic("Tip gresit de declaratie. Tipuri permise: export/report")
	}

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

		// Prima declaratie din coada
		var declaration_queue Declaration_queue

		// Setare denumire declarie si tabela declaratie corespunzator cu denumirea campului din baza de date
		if tipCautat == "report" {
			declaration_queue.tip = "reporting_declaration_id"
			declaration_queue.declaration_table = "reporting_declaration"
		} else {
			declaration_queue.tip = "declaration_id"
			declaration_queue.declaration_table = "declaration"
		}

		// Selectam prima declaratie din coada
		sql_statement := fmt.Sprintf("SELECT id, %s from declaration_queue  where status = 'queued' and type = '%s' limit 1;", declaration_queue.tip, tipCautat)
		rows, err := db.Query(sql_statement)
		checkError(err)
		defer rows.Close()

		// Daca avem cel putin o declaratie in coada
		if rows.Next() {
			var pad_url string
			var lrn string

			if tipCautat == "report" {
				// Extragere reporting_declaration_id  rows -> int
				err = rows.Scan(&declaration_queue.id, &declaration_queue.reporting_declaration_id)
				declaration_queue.curent_declaration_id = *declaration_queue.reporting_declaration_id
				checkError(err)
				fmt.Printf("Declaratie report in asteptare gasita: (%d)\n", *declaration_queue.reporting_declaration_id)

				// Selectare pad_url
				sql_statement := fmt.Sprintf("select pad_url from users where id = (SELECT users_id _id from reporting_declaration where id = %d );", *declaration_queue.reporting_declaration_id)
				rows, err := db.Query(sql_statement)
				checkError(err)
				defer rows.Close()

				// Extragere pad_url rows -> string
				rows.Next()
				err = rows.Scan(&pad_url)
				checkError(err)
				fmt.Println("Pad_url gasit:", pad_url)

				// Selectare lrn
				sql_statement = fmt.Sprintf("SELECT lrn from reporting_declaration where id = %d ;", *declaration_queue.reporting_declaration_id)
				rows, err = db.Query(sql_statement)
				checkError(err)
				defer rows.Close()

				// Extragere lrn -> string
				rows.Next()
				err = rows.Scan(&lrn)
				checkError(err)
				fmt.Println("Lrn gasit:", lrn)
			} else { // tipCautat == "export"
				// Extragere declaration_id  rows -> int
				err = rows.Scan(&declaration_queue.id, &declaration_queue.declaration_id)
				declaration_queue.curent_declaration_id = *declaration_queue.declaration_id
				checkError(err)
				fmt.Printf("Declaratie export in asteptare gasita: (%d)\n", *declaration_queue.declaration_id)

				// Selectare pad_url
				sql_statement := fmt.Sprintf("select pad_url from users u where id = (SELECT users_id  from declaration where id = %d);", *declaration_queue.declaration_id)
				rows, err := db.Query(sql_statement)
				checkError(err)
				defer rows.Close()

				// Extragere pad_url rows -> string
				rows.Next()
				err = rows.Scan(&pad_url)
				checkError(err)
				fmt.Println("Pad_url gasit:", pad_url)
			}

			// Selectare detalii cont
			sql_statement = fmt.Sprintf("select username, password, office from account where id = (SELECT users_id  from %s where id = %d);", declaration_queue.declaration_table, declaration_queue.curent_declaration_id)
			rows, err = db.Query(sql_statement)
			checkError(err)
			defer rows.Close()

			// Account curent
			var account Account
			// Extragere detalii cont
			rows.Next()
			err = rows.Scan(&account.userName, &account.password, &account.office)
			checkError(err)
			fmt.Printf("Account gasit: (%s, %s, %s)\n", account.userName, account.password, account.office)

			// Setare json
			/*var jsonData = []byte(`{
				"USERNAME" : account.userName,
				"PASSWORD" : account.password,
				"BIROU_VAMAL" : account.office,
				"LRN" : lrn
			}`)*/
			var jsonData = []byte(`{
				"name": "morpheus",
				"job": "leader"
			}`)

			// Request catre robot
			request, err := http.NewRequest("POST", pad_url, bytes.NewBuffer(jsonData))
			checkError(err)
			request.Header.Set("Content-Type", "application/json; charset=UTF-8")
			client := &http.Client{}
			response, err := client.Do(request)
			checkError(err)
			defer response.Body.Close()

			// Afisare raspuns
			fmt.Println("response Status:", response.Status)
			//fmt.Println("response Headers:", response.Header)
			body, err := ioutil.ReadAll(response.Body)
			checkError(err)
			fmt.Println("response Body:", string(body))

			// Setare status nou
			/*var status string
			if response.Status == "202" {
				status = "on progress"
			} else {
				status = "error"
			}

			// Actualizare status in bd
			sql_statement = "UPDATE declaration"*/

		} else {
			fmt.Println("Nu avem declaratii in coada")
		}

	} else {
		fmt.Println("Exista deja o declaratie in progres")
	}

}

func main() {
	checkJobs("export") // Primeste ca parametrii "export" sau "report"
	//checkJobs("report")
	println("Exit code 0")
}
