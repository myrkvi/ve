package main

import "database/sql"

//AddEntry adds a dictionary entry to the table adn database specified.
func AddEntry(word string, ipa string, class string, description string, conn *sql.DB, tbl string) int {
	resAdd, err := conn.Query("INSERT INTO ? VALUES (?); SELECT last_insert_rowid()", tbl, word, ipa, class, description)
	if err != nil {
		panic(err)
	}
	var id int
	resAdd.Scan(&id)
	return id
}
