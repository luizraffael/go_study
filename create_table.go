package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, erroOpen := sql.Open("mysql", "root:root@tcp(localhost:3306)/crud_go")

	if erroOpen != nil {
		log.Fatal(erroOpen.Error())
	}
	erroPing := db.Ping()
	if erroPing != nil {
		log.Fatal(erroPing.Error())
	}

	fmt.Println("Inserindo dados...")

	_, erroInsert := db.Exec("INSERT INTO books (author, title) VALUES " +
		"('Mel Kiranis', 'Os rios huways'), " +
		"('Geremias', 'A batalha'), " +
		"('Keanas','The Book'); ")

	if erroInsert != nil {
		log.Fatal(erroInsert.Error())
	}
}
