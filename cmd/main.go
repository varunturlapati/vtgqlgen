package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/varunturlapati/vtgqlgen/datasources/db" // update the username
	"github.com/varunturlapati/vtgqlgen/gqlgen"         // update the username
)

func main() {
	// initialize the db
	dsn := "reg:pass@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	dbObj, err := db.Open(dsn)
	if err != nil {
		panic(err)
	}
	// defer dbObj.Close()

	// initialize the repository
	repo := db.NewRepository(dbObj)

	// configure the server
	mux := http.NewServeMux()
	mux.Handle("/", gqlgen.NewPlaygroundHandler("/query"))
	mux.Handle("/query", gqlgen.NewHandler(repo))

	// run the server
	port := ":7777"
	fmt.Fprintf(os.Stdout, "🚀 Server ready at http://localhost%s\n", port)
	fmt.Fprintln(os.Stderr, http.ListenAndServe(port, mux))
}