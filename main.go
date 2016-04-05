package main

import (
	"database/sql"
	"fmt"
	"log"
	"os/user"
	"strconv"
	"strings"

	"github.com/elgris/sqrl"

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
	ve define <word>
	ve modify <word>
	ve link	  <word>
	ve remove <word>
	ve -h | --help

Options:
	-h --help  Shows this help text.`

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
			lookupWordNat(args["<word>"].(string), conn)
		} else if args["-c"].(bool) {
			//TODO: add lookupWordCon()
		} else {
			fmt.Println("===== Results from Natlang dictionary =====")
			lookupWordNat(args["<word>"].(string), conn)
			//TODO: add lookupWordCon()
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

func lookupWord(q string, conn *sql.DB, tbl string) []DictionaryEntry {
	res, err := sqrl.Select("*").From(tbl).Where("Word = ?", q).RunWith(conn).Query()
	if err != nil {
		panic(err)
	}

	var dictionaryEntries []DictionaryEntry

	for res.Next() {
		x := DictionaryEntry{}
		err = res.Scan(&x.ID, &x.Word, &x.Class)
		if err != nil {
			panic(err)
		}

		dictionaryEntries = append(dictionaryEntries, x)
	}

	return dictionaryEntries
}

func lookupWordNat(q string, conn *sql.DB) {
	dictionaryEntriesNat := lookupWord(q, conn, "Natlang")

	for _, entry := range dictionaryEntriesNat {
		/*q := sqrl.Select("Conlang_Id").From("Conlang_Natlang_relation").Where("Natlang_Id = ?", entry.ID)
		fmt.Println(q.ToSql())
		relRes, err := q.RunWith(conn).Query()*/
		relRes, err := conn.Query("SELECT Conlang_Id FROM Conlang_Natlang_relation WHERE Natlang_Id = ?", entry.ID)
		if err != nil {
			panic(err)
		}

		var IDC []int

		for relRes.Next() {
			var x int
			err = relRes.Scan(&x)
			if err != nil {
				panic(err)
			}

			IDC = append(IDC, x)
			var IDCs []string

			for _, v := range IDC {
				IDCs = append(IDCs, strconv.Itoa(v))
			}
			IDCstring := strings.Join(IDCs, ",")
			fmt.Println(IDCstring)

			/*q := sqrl.Select("*").From("Conlang").Where("Id = ?", IDCstring)

			fmt.Println(q.ToSql())

			conRes, err := q.RunWith(conn).Query()*/
			conRes, err := conn.Query("SELECT * FROM Conlang WHERE Id = ?", IDCstring)
			if err != nil {
				panic(err)
			}

			var dictionaryEntriesCon []DictionaryEntry

			for conRes.Next() {
				var x DictionaryEntry
				err = conRes.Scan(&x.ID, &x.Word, &x.IPA, &x.Class, &x.Description)

				if err != nil {
					panic(err)
				}

				dictionaryEntriesCon = append(dictionaryEntriesCon, x)

			}

			entry.Translations = dictionaryEntriesCon
		}
	}

	for i, natEntry := range dictionaryEntriesNat {
		fmt.Printf("(%d.) -- %s    %s --\n", i, natEntry.Word, natEntry.Class)

		for i, translation := range natEntry.Translations {
			fmt.Printf("\t(%d.) %s [%s]    %s\n\t\t%s", i, translation.Word, translation.IPA, translation.Class, translation.Description)
		}
	}

}
