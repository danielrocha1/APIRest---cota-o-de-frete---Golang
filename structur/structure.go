package structur

type Api_request struct {
	TokenOferta     string            `json:"token_oferta"`
	Transportadoras []Transportadoras `json:"transportadoras"`
	Volumes         []Volumes         `json:"volumes"`
}

type Transportadoras struct {
	Oferta           int     `json:"oferta"`
	Cnpj             string  `json:"cnpj"`
	Logotipo         string  `json:"logotipo"`
	Nome             string  `json:"nome"`
	Servico          string  `json:"servico"`
	DescricaoServico string  `json:"descricao_servico"`
	PrazoEntrega     int     `json:"prazo_entrega"`
	EntregaEstimada  string  `json:"entrega_estimada"`
	Validade         string  `json:"validade"`
	CustoFrete       float64 `json:"custo_frete"`
	PrecoFrete       float64 `json:"preco_frete"`
}
type Volumes struct {
	Tipo           int     `json:"tipo"`
	Sku            string  `json:"sku"`
	Tag            string  `json:"tag"`
	Descricao      string  `json:"descricao"`
	Quantidade     int     `json:"quantidade"`
	Altura         float64 `json:"altura"`
	Largura        float64 `json:"largura"`
	Comprimento    float64 `json:"comprimento"`
	Peso           int     `json:"peso"`
	Valor          int     `json:"valor"`
	VolumesProduto int     `json:"volumes_produto"`
}


//------------------------------------------------------------------------------------------------------------------------


type Api struct {
	Remetente        Remetente    `json:"remetente"`
	Destinatario     Destinatario `json:"destinatario"`
	Carga            []Volume     `json:"volumes"`
	CodigoPlataforma interface{}       `json:"codigo_plataforma"`
	Token            interface{}       `json:"token"`
}

type Volume struct {
	Tipo        interface{}     `json:"tipo"`
	Quantidade  interface{}     `json:"quantidade"`
	Peso        interface{}     `json:"peso"`
	Valor       interface{}     `json:"valor"`
	Sku         interface{}  `json:"sku"`
	Altura      interface{} `json:"altura"`
	Largura     interface{} `json:"largura"`
	Comprimento interface{} `json:"comprimento"`
}

type Destinatario struct {
	TipoPessoa        interface{}    `json:"tipo_pessoa"`
	CnpjCpf           interface{} `json:"cnpj_cpf"`
	InscricaoEstadual interface{} `json:"inscricao_estadual"`
	Endereco         Endereco     `json:"endereco"`
}

type Remetente struct {
	Cnpj interface{} `json:"cnpj"`
}

type Endereco struct {
	Cep interface{} `json:"cep"`
}
