package app

// CollectionResponseDTO is a data transfer object needed for the RESTful API implementation.
// It contains short information about the collection and the urls for each item in the collection.
type CollectionResponseDTO struct {
	Info string   `json:"info"`
	Urls []string `json:"urls"`
}

// StockDTO is a data transfer object that can be used for marshaling and unmarshaling
// an existing stock item
type StockDTO struct {
	ID string `json:"id"`

	Name string    `json:"name"`
	Type stockType `json:"type"`

	ExpirationDate string `json:"expirationDate"`
	MinQuantity    string `json:"minQuantity"`

	Distributor string `json:"distributor"`
}

// NewStockDTO is a data transfer object that can be used between
// reading a JSON with data for a new stock item and
// creating the new stock item with NewStock(*NewStockDTO) (Stock, error)
type NewStockDTO struct {
	Name string    `json:"name"`
	Type stockType `json:"type"`

	ExpirationDate string `json:"expirationDate"`
	MinQuantity    string `json:"minQuantity"`

	Distributor string `json:"distributor"`
}
