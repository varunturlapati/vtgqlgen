package main

import (
	"fmt"
	"github.com/varunturlapati/vtgqlgen/dataloaders"
	"net/http"
	"os"

	ds "github.com/varunturlapati/vtgqlgen/datasource"
	"github.com/varunturlapati/vtgqlgen/datasource/db"
	"github.com/varunturlapati/vtgqlgen/gqlgen"
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
	repo := ds.NewRepository(dbObj)
	// initialize the dataloaders
	dl := dataloaders.NewRetriever()

	// configure the server
	mux := http.NewServeMux()
	mux.Handle("/", gqlgen.NewPlaygroundHandler("/query"))
	dlMiddleware := dataloaders.Middleware(repo)
	queryHandler := gqlgen.NewHandler(repo, dl)
	mux.Handle("/query", dlMiddleware(queryHandler))

	// run the server
	port := ":7777"
	fmt.Fprintf(os.Stdout, "🚀 Server ready at http://localhost%s\n", port)
	fmt.Fprintln(os.Stderr, http.ListenAndServe(port, mux))
}
