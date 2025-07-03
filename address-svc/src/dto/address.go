package dto

// Entrada para crear
type AddressCreate struct {
	Floor   string `json:"floor"`
	Unit    string `json:"unit"`
	Street  string `json:"street"  binding:"required"`
	Number  int    `json:"number"  binding:"required"`
	City    string `json:"city"    binding:"required"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
	Country string `json:"country" binding:"required,len=2"`
}

// Entrada para update (campos opcionales)
type AddressUpdate struct {
	Floor   *string `json:"floor"`
	Unit    *string `json:"unit"`
	Street  *string `json:"street"`
	Number  *int    `json:"number"`
	City    *string `json:"city"`
	State   *string `json:"state"`
	Zip     *string `json:"zip"`
	Country *string `json:"country"`
}
