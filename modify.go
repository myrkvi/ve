package main

import (
	"database/sql"
	"fmt"
	"strings"
)

//ModifyEntry modifies an existing entry specified by Id
func ModifyEntry(id string, word string, class string, description string, ipa string, conn *sql.DB, tbl string) {
	query := "UPDATE " + tbl + " SET "
	if word != "" {
		query += "Word='" + word + "', "
	}
	if class != "" {
		query += "Class='" + class + "', "
	}
	if description != "" {
		query += "Description='" + description + "', "
	}
	if ipa != "" {
		query += "Ipa='" + ipa + "', "
	}

	query = strings.Trim(query, ", ") + ";"

	_, err := conn.Exec(query)
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Modified entry with ID %s.\n", id)
	}
}
