package app

/*
 * JSON format for the GET requests for "/<items>/"
 */
type CollectionResponseDTO struct {
	Info string   `json:"info"`
	Urls []string `json:"urls"`
}

/*
 * JSON format for the GET requests for "/<items>/<id>"
 */
type StockDTO struct {
	Id string `json:"id"`

	Name string    `json:"name"`
	Type stockType `json:"type"`

	ExpirationDate string `json:"expirationDate"`
	MinQuantity    string `json:"minQuantity"`

	Distributor string `json:"distributor"`
}

/*
 * new stock items' JSON format
 * IDs of new stock objects should be handled only by the back-end
 */
type NewStockDTO struct {
	Name string    `json:"name"`
	Type stockType `json:"type"`

	ExpirationDate string `json:"expirationDate"`
	MinQuantity    string `json:"minQuantity"`

	Distributor string `json:"distributor"`
}
