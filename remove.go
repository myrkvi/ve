package main

import (
	"database/sql"
	"fmt"
)

//RemoveEntry removes an entry from the specified table with the specified Id.
func RemoveEntry(id string, conn *sql.DB, tbl string) {
	_, err := conn.Exec("DELETE FROM " + tbl + " WHERE   Id=" + id + ";")
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Deleted entry with ID %s from table %s.\n", id, tbl)
	}

}
