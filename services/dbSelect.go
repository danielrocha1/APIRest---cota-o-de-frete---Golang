package services

import (
	auth "FreteRapido/auth"
	structure "FreteRapido/structur"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
)

func Http_request(writer http.ResponseWriter, reader *http.Request) ([]byte, int) {
	var info structure.Api 

	err := json.NewDecoder(reader.Body).Decode(&info)
	auth.Check_error(err)
	a,err := json.MarshalIndent(info, "", "\t")
	auth.Check_error(err)
	
	resp, err := http.Post("https://freterapido.com/api/external/embarcador/v1/quote-simulator", "application/json", bytes.NewBuffer(a))
	
	if resp.StatusCode == 400 {
		check_body(writer, resp)
	}
	if resp.StatusCode != 200 {
		checkHttp(writer, resp.StatusCode)
	}

	auth.Check_error(err)

	body, err := ioutil.ReadAll(resp.Body)
	auth.Check_error(err)
	
	return body, resp.StatusCode
}

//===========================================================================

func InsertTransportadora(writer http.ResponseWriter,db *sql.DB, value []byte) {
	var values structure.Api_request
	json.Unmarshal(value, &values)
	auth.Check_TokenOferta(auth.Filtro_TokenOfertaT(db), values.TokenOferta)

	for item := range values.Transportadoras {
		_, err := db.Exec("INSERT INTO transportadoras(Cnpj,Logotipo,Nome,Servico,DescricaoServico,PrazoEntrega,EntregaEstimada,Validade,CustoFrete,PrecoFrete,TokenOferta,Oferta) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)",
			values.Transportadoras[item].Cnpj,
			values.Transportadoras[item].Logotipo,
			values.Transportadoras[item].Nome,
			values.Transportadoras[item].Servico,
			values.Transportadoras[item].DescricaoServico,
			values.Transportadoras[item].PrazoEntrega,
			values.Transportadoras[item].EntregaEstimada,
			values.Transportadoras[item].Validade,
			values.Transportadoras[item].CustoFrete,
			values.Transportadoras[item].PrecoFrete,
			values.TokenOferta,
			values.Transportadoras[item].Oferta)

		auth.Check_error(err)
	}
	fmt.Fprintf(writer,"\nEste é o Token para consulta:%s\n",values.TokenOferta)
} 

func InsertVolumes(db *sql.DB, value []byte) {
	var values structure.Api_request
	json.Unmarshal(value, &values)
	auth.Check_TokenOferta(auth.Filtro_TokenOfertaV(db), values.TokenOferta)

	for item := range values.Volumes {
		_, err := db.Exec("INSERT INTO volumes (Tipo,Sku,Tag,Descricao,Quantidade,Altura,Largura,Comprimento,Peso,Valor,VolumesProduto,TokenOferta) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)",
			values.Volumes[item].Tipo,
			values.Volumes[item].Sku,
			values.Volumes[item].Tag,
			values.Volumes[item].Descricao,
			values.Volumes[item].Quantidade,
			values.Volumes[item].Altura,
			values.Volumes[item].Largura,
			values.Volumes[item].Comprimento,
			values.Volumes[item].Peso,
			values.Volumes[item].Valor,
			values.Volumes[item].VolumesProduto,
			values.TokenOferta)

		auth.Check_error(err)
	}
}

//===========================================================================


func Db_CountNome(db *sql.DB, writer http.ResponseWriter, reader *http.Request, TokenOferta string) map[string]int {
	result, err := db.Query("SELECT `Nome` FROM `transportadoras` WHERE TokenOferta=(?) ", TokenOferta)
	auth.Check_error(err)

	var name []string
	for result.Next() {
		var db_nome string
		err := result.Scan(&db_nome)
		auth.Check_error(err)
		name = append(name, db_nome)
	}

	var check = make(map[string]int)
	for _, item := range name {
		_, exist := check[item]
		if exist {
			check[item] += 1
		} else {
			check[item] = 1
		}
	}

	fmt.Fprintf(writer,"\nResultados Por Transportadoras\n")
	for key, value := range check {
		fmt.Fprintf(writer,"Transportadora:%s  | Quantidade de Cotação:%d\n", key, value)
	}
	return check
}

func Db_TotalPreco(db *sql.DB, writer http.ResponseWriter, reader *http.Request,TokenOferta string) map[string]float64 {
	result, err := db.Query("SELECT `Nome`, `PrecoFrete` FROM `transportadoras` WHERE TokenOferta=(?)", TokenOferta)
	auth.Check_error(err)

	var name []string
	var price []float64
	for result.Next() {
		var db_nome string
		var db_preco float64
		erro := result.Scan(&db_nome, &db_preco)

		auth.Check_error(erro)
		name = append(name, db_nome)
		price = append(price, db_preco)
	}
	var check = make(map[string]float64)
	for i, str := range name {
		_, exist := check[str]
		if exist {
			sum := check[str]
			check[str] = sum + price[i]
		} else{
			check[str] = price[i]
		}
	}

	fmt.Fprintf(writer,"\nPreço Total das Transportadoras\n")
	for str, value := range check {
		fmt.Fprintf(writer,"Transportadora:%s  | Preço Total:%.2f\n", str, value)
	}
	return check
} 

func Db_Media_Total(db *sql.DB, writer http.ResponseWriter, reader *http.Request, TokenOferta string) {
	check_count := Db_CountNome(db, writer, reader, TokenOferta)
	check_preco := Db_TotalPreco(db,writer,reader,TokenOferta)
	var media_transp = make(map[string]float64)
	for transp, i := range check_count {
		_, exist := check_preco[transp]
		if exist {
			sum := check_preco[transp]
			media_transp[transp] = sum / float64(i)
		}
	}
	fmt.Fprintf(writer, "\nPreço Médio das Transportadoras\n")
	for str, value := range media_transp {
		fmt.Fprintf(writer, "Transportadora:%s  | Preço médio:%.2f\n", str, value)
	}
	fmt.Println("")
}

func Db_Servico(db *sql.DB, writer http.ResponseWriter, reader *http.Request, TokenOferta string) {
	result, err := db.Query("SELECT `Nome`, `Servico`, `PrecoFrete` FROM `transportadoras` WHERE TokenOferta=(?) ", TokenOferta)
	auth.Check_error(err)

	var name []string
	var servico []string
	var price []float64
	for result.Next() {
		var db_nome string
		var db_servico string
		var db_preco float64
		err := result.Scan(&db_nome, &db_servico, &db_preco)
		auth.Check_error(err)

		name = append(name, db_nome)
		servico = append(servico, db_servico)
		price = append(price, db_preco)
	}

	max_offer := -1.0
	var point_max int
	for i, flt := range price {
		if max_offer <= flt {
			max_offer = flt
			point_max = i
		}
	}

	min_offer := 100000.00
	var point_min int
	for i, flt := range price {
		if min_offer >= flt {
			min_offer = flt
			point_min = i
		}
	}
	fmt.Fprintf(writer, "\nEssa é a maior oferta\n")
	fmt.Fprintf(writer, "Transportadora:%s |Serviço:%s | Preço:%.2f\n\n", name[point_max], servico[point_max], max_offer)
	fmt.Fprintf(writer, "Essa é a menor oferta\n")
	fmt.Fprintf(writer, "Transportadora:%s |Serviço:%s | Preço:%.2f\n", name[point_min], servico[point_min], min_offer)
}

//========================================= ROTAS TERMINADAS =========

func Rota_1(writer http.ResponseWriter, reader *http.Request) {
	body, value := Http_request(writer, reader)
	if value == 200 {
		InsertTransportadora(writer, auth.Connect(), body)
		InsertVolumes(auth.Connect(), body)
		fmt.Fprintf(writer, "\nINFORMAÇÕES FORAM INSERIDAS NO BANCO")
		fmt.Fprintf(writer, "")
	}
}

func Rota_2(db *sql.DB, writer http.ResponseWriter, reader *http.Request, TokenOferta string) {
	Db_Media_Total(db, writer, reader, TokenOferta)
	Db_Servico(db, writer, reader, TokenOferta)
	fmt.Println("")
}

func LastQuotes(writer http.ResponseWriter, reader *http.Request) {
	con := auth.Connect()
	
	TokenOferta := reader.URL.Query().Get("last_quotes")

	result, err := con.Query("SELECT `Cnpj`, `Logotipo`, `Nome`, `Servico`, `DescricaoServico`, `PrazoEntrega`, `Validade`, `CustoFrete`, `PrecoFrete`, `Oferta` FROM `transportadoras` WHERE `TokenOferta` = (?) ", TokenOferta)
	auth.Check_error(err)

	var info []structure.Transportadoras
	for result.Next() {
		var t structure.Transportadoras
		err = result.Scan(
			&t.Cnpj,
			&t.Logotipo,
			&t.Nome,
			&t.Servico,
			&t.DescricaoServico,
			&t.PrazoEntrega,
			&t.Validade,
			&t.CustoFrete,
			&t.PrecoFrete,
			&t.Oferta)
		auth.Check_error(err)
		info = append(info, t)
	}

	for item := len(info) - 1; item >= 0; item-- {
		fmt.Fprintf(writer, "|Nome:%s | Cnpj:%s | Serviço:%s | Descrição Serviço:%s|\n PrazoEntrega:%d | Validade:%s | Custo Frete:%.2f| Preço Frete:%.2f| Oferta:%d |\n\n",
			info[item].Nome, info[item].Cnpj, info[item].Servico, info[item].DescricaoServico, info[item].PrazoEntrega, info[item].Validade, info[item].CustoFrete, info[item].PrecoFrete, info[item].Oferta)
	}
	
	Rota_2(con, writer, reader, TokenOferta)
	defer con.Close()
}

//============================================== CONFIGURAÇÃO E ERROS ========

func checkHttp(writer http.ResponseWriter, erro int) {
	check := map[int]string{
		401: "Não autorizado (Token inválido) | Erro de autenticação com a API. ",
		403: "Não tem permissão para acesso ao recurso solicitado. Ou ausência de saldo para contratação de frete",
		404: "Não encontrado | Recurso ou informação inválida e não pode ser encontrada.",
		409: "Conflito não é permitido | Recurso ou informação inválida e não pode ser encontrada.",
		422: "Formato JSON com erro de sintaxe, valores faltando ou tipo de parâmetros inválidos.",
		500: "Erro interno no servidor do Frete Rápido. Você pode relatar isso para o suporte técnico"}
	fmt.Fprintf(writer, check[erro])
}

func check_body(writer http.ResponseWriter, http *http.Response) {
	byte_erro, _ := ioutil.ReadAll(http.Body)
	var byteJson bytes.Buffer

	err := json.Indent(&byteJson, byte_erro, "", "\t")
	auth.Check_error(err)

	fmt.Fprintf(writer, (byteJson.String()))
}

