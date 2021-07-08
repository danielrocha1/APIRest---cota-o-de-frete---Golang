package main

import (
	auth "FreteRapido/auth"
	dbSelect "FreteRapido/services"

	// "database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

//==========================================================================================================

func main() {

	
	configurarServidor()

}

func rotaPrincipal(writer http.ResponseWriter, reader *http.Request) {
	fmt.Fprintf(writer, "Bem vindo, CONECTADO!")
}

func configurarRotas() {
	router := mux.NewRouter()

	http.Handle("/", router)
	router.HandleFunc("/", rotaPrincipal)
	router.HandleFunc("/quote", dbSelect.Rota_1).Methods("POST")
	router.HandleFunc("/metrics", dbSelect.LastQuotes).Methods("GET")
}

func configurarServidor() {
	
	configurarRotas()

	fmt.Println("Servidor está rodando na porta 1337")
	err := http.ListenAndServe(":1337", nil)
	auth.Check_error(err)
}
