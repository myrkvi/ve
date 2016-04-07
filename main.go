package main

import (
	"database/sql"
	"fmt"
	"log"
	"os/user"

	"github.com/docopt/docopt-go"
	_ "github.com/mattn/go-sqlite3"
)

var dataDir string

//DictionaryEntry holds information from a row in the database.
type DictionaryEntry struct {
	ID           int
	Word         string
	IPA          string
	Class        string
	Description  string
	Translations []DictionaryEntry
}

const usage string = `
Vé version 0.0.1

Usage:
	ve lookup [-c|-n] <word>
	ve define (-c|-n) <word> [ipa] [class] [description]
	ve modify (-c|-n) <id> <word> [ipa] [class] [description]
	ve link <n-id> <c-id>
	ve remove (-c|-n) <id>
	ve -h | --help

Options:
	-h --help  	Shows this help text.
	-c  		Apply this to the Conlang list.
	-n  		Apply this to the Natlang list.`

func main() {
	log.Println("Entered main function")
	currUser, _ := user.Current()
	dataDir = currUser.HomeDir + "/.ve/"

	conn, err := sql.Open("sqlite3", dataDir+"dictionary.db")
	if err != nil {
		panic(err)
	}

	err = initDb(conn)
	if err != nil {
		panic(err)
	}

	args, err := docopt.Parse(usage, nil, true, "0.0.1", false, true)
	if err != nil {
		panic(err)
	}

	switch true {
	case args["lookup"]:
		if args["-n"].(bool) {
			LookupWordNat(args["<word>"].(string), conn)
		} else if args["-c"].(bool) {
			LookupWordCon(args["<word>"].(string), conn)
		} else {
			fmt.Println("===== Results from Natlang dictionary =====")
			LookupWordNat(args["<word>"].(string), conn)
			fmt.Println("\n===== Results from Conlang dictionary =====")
			LookupWordCon(args["<word>"].(string), conn)
		}

	case args["define"]:
		//TODO: Add functionality.

	case args["modify"]:
		//TODO: Add functionality.

	case args["link"]:
		//TODO: Add functionality.

	case args["remove"]:
		//TODO: Add functionality.
	}
}

func initDb(conn *sql.DB) error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS Conlang (
			Id INTEGER NOT NULL PRIMARY KEY,
			Word TEXT,
			Ipa TEXT,
			Class TEXT,
			Description TEXT
		);

		CREATE TABLE IF NOT EXISTS Natlang (
			Id INTEGER NOT NULL PRIMARY KEY,
			Word TEXT,
			Ipa TEXT,
			Class TEXT,
			Description TEXT
		);

		CREATE TABLE IF NOT EXISTS Conlang_Natlang_relation (
			Id INTEGER NOT NULL PRIMARY KEY,
			Conlang_Id INTEGER,
			Natlang_Id INTEGER,
			FOREIGN KEY (Conlang_Id) REFERENCES Conlang (Id),
			FOREIGN KEY (Natlang_Id) REFERENCES Natlang (Id)
		);
		`)
	return err
}
