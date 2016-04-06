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
		err = res.Scan(&x.ID, &x.Word, &x.IPA, &x.Class, &x.Description)
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

func lookupWordCon(q string, conn *sql.DB) {
	dictionaryEntriesCon := lookupWord(q, conn, "Conlang")
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

			natRes, err := conn.Query("SELECT * FROM Natlang WHERE Id = ?", IDNstring)
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

	for i, conEntry := range entries {
		fmt.Printf("(%d.) -- %s [%s]    %s --\n\t%s\n", i+1, conEntry.Word, conEntry.IPA, conEntry.Class, conEntry.Description)

		for i, translation := range conEntry.Translations {
			fmt.Printf("\t(%d.) %s    %s\n\t\t%s", i+1, translation.Word, translation.Class, translation.Description)
		}
	}

}
