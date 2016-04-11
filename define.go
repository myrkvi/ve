package main

import (
	"database/sql"
	"fmt"
)

//AddEntry adds a dictionary entry to the table adn database specified.
func AddEntry(word string, ipa string, class string, description string, conn *sql.DB, tbl string) int {
	resAdd, err := conn.Query("INSERT INTO " + tbl + " VALUES (NULL, '" + word + "', '" + ipa + "', '" + class + "', '" + description + "'); SELECT last_insert_rowid();")
	if err != nil {
		panic(err)
	}
	var id int
	resAdd.Scan(&id)
	return id
}

//AddConlangEntry adds an entry to the Conlang table, to later be linked with Natlang entries.
func AddConlangEntry(word string, ipa string, class string, description string, conn *sql.DB) {
	id := AddEntry(word, ipa, class, description, conn, "Conlang")

	fmt.Printf("Successfully added new entry: %s\nID: %d\n", word, id)
}

//AddNatlangEntry adds an entry to the Natlang table, to later be linked with the Conlang table.
func AddNatlangEntry(word string, class string, description string, conn *sql.DB) {
	id := AddEntry(word, "", class, description, conn, "Natlang")

	fmt.Printf("Successfully added new entry: %s\nID: %d\n", word, id)
}
