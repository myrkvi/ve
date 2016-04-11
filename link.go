package main

import (
	"database/sql"
	"fmt"
)

//LinkWords adds a new link between a Conlang word and a Natlang word, to display
//them together in search queries.
func LinkWords(natlangID string, conlangID string, conn *sql.DB) {
	_, err := conn.Query("INSERT INTO Conlang_Natlang_relation VALUES (NULL, " + natlangID + ", " + conlangID + ");")
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Successfully linked definitions!")
	}
}

//UnlinkWords removes a link between two definitions.
func UnlinkWords(natlangID string, conlangID string, conn *sql.DB) {
	res, err := conn.Query("SELECT Id FROM Conlang_Natlang_relation WHERE Natlang_Id=" + natlangID + " AND Conlang_Id=" + conlangID + ";")
	if err != nil {
		panic(err)
	}

	for res.Next() {
		var id string
		res.Scan(&id)

		conn.Exec("REMOVE FROM Conlang_Natlang_relation WHERE Id=" + id + ";")
	}
}
