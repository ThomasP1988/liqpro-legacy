package kraken

// OrderDescription - structure of order description
type OrderDescription struct {
	Pair           string  `json:"pair"`
	Side           string  `json:"type"`
	OrderType      string  `json:"ordertype"`
	Price          float64 `json:"price,string"`
	Price2         float64 `json:"price2,string"`
	Leverage       string  `json:"leverage"`
	Info           string  `json:"order"`
	CloseCondition string  `json:"close"`
}

// AddOrderResponse - response on AddOrder request
type AddOrderResponse struct {
	Description    OrderDescription `json:"descr"`
	TransactionIds []string         `json:"txid"`
}

// DataParser struct to avoid garbage when parsing data
type DataParser struct {
	Result      OrderBookUpdate
	body        map[string]interface{}
	ok          bool
	k           string
	v           interface{}
	checkSum    string
	err         error
	itemsParser ItemsParser
}

// ItemsParser struct to avoid garbage when parsing data
type ItemsParser struct {
	items         []OrderBookItem
	orderBookItem OrderBookItem
	updates       []interface{}
	entity        []interface{}
	item          interface{}
	ok            bool
}
