package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/elgris/sqrl"
	_ "github.com/mattn/go-sqlite3"
)

//LookupWord looks up a word in the specified table and returns a DictionaryEntry slice.
func LookupWord(q string, conn *sql.DB, tbl string) []DictionaryEntry {
	res, err := sqrl.Select("*").From(tbl).Where("Word = ?", q).RunWith(conn).Query()
	if err != nil {
		panic(err)
	}

	var dictionaryEntries []DictionaryEntry

	for res.Next() {
		x := DictionaryEntry{}
		err = res.Scan(&x.ID, &x.Word, &x.IPA, &x.Class, &x.Description)
		if err != nil {
			panic(err)
		}

		dictionaryEntries = append(dictionaryEntries, x)
	}

	return dictionaryEntries
}

//LookupWordNat looks up words matching q, finds matching words from Conlang and lists them.
func LookupWordNat(q string, conn *sql.DB) {
	dictionaryEntriesNat := LookupWord(q, conn, "Natlang")
	var entries []DictionaryEntry

	for _, entry := range dictionaryEntriesNat {
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

			conRes, err := conn.Query("SELECT * FROM Conlang WHERE Id IN (?)", IDCstring)
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

	if entries == nil {
		entries = dictionaryEntriesNat
	}

	for i, natEntry := range entries {
		fmt.Printf("(%d.) -- %s    %s --\n\t%s\n\tID: %d\n", i+1, natEntry.Word, natEntry.Class, natEntry.Description, natEntry.ID)

		for i, translation := range natEntry.Translations {
			fmt.Printf("\t(%d.) %s [%s]    %s\n\t\t%s\n\t\tID: %d\n", i+1, translation.Word, translation.IPA, translation.Class, translation.Description, translation.ID)
		}
	}

}

//LookupWordCon looks up words matching q, finds matching words from Natlang and lists them.
func LookupWordCon(q string, conn *sql.DB) {
	dictionaryEntriesCon := LookupWord(q, conn, "Conlang")
	var entries []DictionaryEntry

	for _, entry := range dictionaryEntriesCon {
		relRes, err := conn.Query("SELECT Natlang_Id FROM Conlang_Natlang_relation WHERE Conlang_Id = ?", entry.ID)
		if err != nil {
			panic(err)
		}

		var IDN []int

		for relRes.Next() {
			var x int
			err = relRes.Scan(&x)
			if err != nil {
				panic(err)
			}

			IDN = append(IDN, x)
			var IDNs []string

			for _, v := range IDN {
				IDNs = append(IDNs, strconv.Itoa(v))
			}
			IDNstring := strings.Join(IDNs, ",")

			natRes, err := conn.Query("SELECT * FROM Natlang WHERE Id IN (?)", IDNstring)
			if err != nil {
				panic(err)
			}

			var dictionaryEntriesNat []DictionaryEntry

			for natRes.Next() {
				var x DictionaryEntry
				err = natRes.Scan(&x.ID, &x.Word, &x.IPA, &x.Class, &x.Description)

				if err != nil {
					panic(err)
				}

				dictionaryEntriesNat = append(dictionaryEntriesNat, x)

			}

			entry.Translations = dictionaryEntriesNat
			entries = append(entries, entry)

		}
	}
	if entries == nil {
		entries = dictionaryEntriesCon
	}

	for i, conEntry := range entries {
		fmt.Printf("(%d.) -- %s [%s]    %s --\n\t%s\n\tID: %d\n", i+1, conEntry.Word, conEntry.IPA, conEntry.Class, conEntry.Description, conEntry.ID)

		for i, translation := range conEntry.Translations {
			fmt.Printf("\t(%d.) %s    %s\n\t\t%s\n\t\tID: %d\n", i+1, translation.Word, translation.Class, translation.Description, translation.ID)
		}
	}

}
