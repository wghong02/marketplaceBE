package main

import (
	sqlMethods "appBE/database"
	"appBE/handler"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("started-service")
	// initialize database
    sqlMethods.InitSQLDatabase()
	// deploy to port 8080
	log.Fatal(http.ListenAndServe(":8080", handler.InitRouter()))
}
