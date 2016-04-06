package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/elgris/sqrl"
	_ "github.com/mattn/go-sqlite3"
)

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
	var entries []DictionaryEntry

	//THE ISSUE IS HERE AND YOU ARE A DUMB MOTHERFUCKER, VEGARD.
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
			entries = append(entries, entry)

		}
	}

	for i, natEntry := range entries {
		fmt.Printf("(%d.) -- %s    %s --\n", i+1, natEntry.Word, natEntry.Class)

		for i, translation := range natEntry.Translations {
			fmt.Printf("\t(%d.) %s [%s]    %s\n\t\t%s", i+1, translation.Word, translation.IPA, translation.Class, translation.Description)
		}
	}

}
