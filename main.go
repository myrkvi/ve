package main

import (
	"database/sql"
	"fmt"
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
		lookupWordNat(args["<word>"].(string), conn)

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
			Class TEXT
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

func lookupWord(q string, conn *sql.DB, tbl string) (int, string, string) {
	res, err := conn.Query("SELECT * FROM ? WHERE Word='?';", tbl, q)
	if err != nil {
		panic(err)
	}

	var ID int
	var Word string
	var Class string

	for res.Next() {
		err = res.Scan(&ID, &Word, &Class)
		if err != nil {
			panic(err)
		}
		fmt.Println(ID, Word, Class)
	}

	return ID, Word, Class
}

func lookupWordNat(q string, conn *sql.DB) {
	IDN, WordN, ClassN := lookupWord(q, conn, "Natlang")
	fmt.Println(IDN, WordN, ClassN)
}
