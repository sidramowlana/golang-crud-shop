package models

type Product struct {
	ID          int    `json:"id"`
	Name        string "json: Name"
	Description string "json: Description"
	Price       int    "json: Price"
	Quantity    int    "json: Quantity"
}

// type ProductsList []Product
var ProductsList = []Product{
	{ID: 1, Name: "Product A", Description: "Desc A", Price: 40, Quantity: 340},
	{ID: 2, Name: "Product B", Description: "Desc B", Price: 10, Quantity: 60},
	{ID: 3, Name: "Product C", Description: "Desc C", Price: 20, Quantity: 240},
	{ID: 4, Name: "Product D", Description: "Desc D", Price: 55, Quantity: 140},
}
