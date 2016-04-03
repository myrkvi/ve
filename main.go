package main

import (
	"database/sql"
	"os/user"

	"github.com/docopt/docopt-go"
	_ "github.com/mattn/go-sqlite3"
)

var dataDir string

const usage string = `
Vé version 0.0.1

Usage:
	ve lookup <word>
	ve define <word>
	ve modify <word>
	ve link	  <word>
	ve remove <word>

Options:
	-h --help			Shows this help text.`

func main() {
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

	args, _ := docopt.Parse(usage, nil, true, "0.0.1", false, true)

	switch true {
	case args["lookup"]:
		//TODO: Add functionality.

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
			Ipa TEXT
			Class TEXT,
			Description TEXT
		);

		CREATE TABLE IF NOT EXISTS Natlang (
			Id INTEGER NOT NULL PRIMARY KEY,
			Word TEXT,
			Class TEXT
		);

		CREATE TABLE IF NOT EXISTS Conlang_Natlang_relation (
			Id INTEGER NOT NULL PRIMARY KEY,
			Conlang_Id INTEGER,
			Natlang_Id INTEGER,
			FOREIGN KEY (Conlang_Id) REFERENCES Conlang (Id),
			FOREIGN KEY (Natlang_Id) REFERENCES Natlang (Id)
		)
		`)
	return err
}
