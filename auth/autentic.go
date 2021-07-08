package auth

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const (
	driver = "mysql"
	user  = "root"
	net = "tcp"
	host = "172.18.0.2"
	port = "3306"
	db = "freterapido_api"

)

func Connect() *sql.DB {
	source := fmt.Sprintf("%s@%s(%s:%s)/%s",user,net,host,port,db)
	db, err := sql.Open(driver,source)
	Check_error(err)
	return db
}

func Check_error(err error) {
	if err != nil {
		log.Fatal(err)
	}
} 

func Filtro_TokenOfertaT(db *sql.DB) []string {
	var token_oferta []string
	var filtro []string
	var check string

	result, err := db.Query("SELECT `TokenOferta` FROM `transportadoras`")
	Check_error(err)

	for result.Next() {
		erro := result.Scan(&check)
		token_oferta = append(token_oferta, check)
		Check_error(erro)

	}
	for item := range token_oferta {
		if token_oferta[item] != check {
			check = token_oferta[item]
			filtro = append(filtro, check)

		}
	}
	if len(filtro) == 0 {
		filtro = append(filtro, check)
	}
	return filtro
}

func Filtro_TokenOfertaV(db *sql.DB) []string {
	var token_oferta []string
	var filtro []string
	var check string

	result, err := db.Query("SELECT `TokenOferta` FROM `volumes`")
	Check_error(err)

	for result.Next() {
		erro := result.Scan(&check)
		token_oferta = append(token_oferta, check)
		Check_error(erro)
	}
	for item := range token_oferta {
		if token_oferta[item] != check {
			check = token_oferta[item]
			filtro = append(filtro, token_oferta[item])
		}
	}
	if len(filtro) == 0 {
		filtro = append(filtro, check)
	}
	return filtro
}

func Check_TokenOferta(value []string, check string) {
	for item := range value {
		if check == value[item] {
			log.Fatal("Token ja inserido no Banco de Dados!")
		}
	}
}


