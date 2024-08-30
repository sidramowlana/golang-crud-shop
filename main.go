package main

import (
	"encoding/json"
	"fmt"
	m "golang-crud/models"
	"net/http"
	"strconv"
	"strings"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "End point to the Home Page")
}

func getAllProducts(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Get all products end point called")
	json.NewEncoder(w).Encode(m.ProductsList)
}
func getAProductById(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Get A product By its Id end point called")

	//get the product id from the url - convert the id string from url to int
	idStr := strings.TrimPrefix(r.URL.Path, "/product/")
	id, err := strconv.Atoi(idStr)
	// check if the taken product id is available
	if (err != nil) || id < 1 {
		http.Error(w, "Invalid ID, product cannot be found", http.StatusBadRequest)
		return
	}

	// if available send the product details
	for _, product := range m.ProductsList {
		if product.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(product)
			return
		}
	}
	http.Error(w, "Product ID cannot be found", http.StatusNotFound)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	// Make sure if the end point method is post
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid reuqest method", http.StatusMethodNotAllowed)
	}

	// if it is a post method u have to decode the body to product struct
	var newProduct m.Product
	// Get the body from the request and assign it to the new variable
	err := json.NewDecoder(r.Body).Decode(&newProduct)

	//Check if the body is valid
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
	}

	// if it is valid -> create a new id
	newProduct.ID = getNextProductID()
	//Append the newly created product
	m.ProductsList = append(m.ProductsList, newProduct)

	// Return the newly created product in the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newProduct)

}
func getNextProductID() int {
	if len(m.ProductsList) == 0 {
		return 1
	}
	lastProduct := m.ProductsList[len(m.ProductsList)-1]
	return lastProduct.ID + 1
}
func deleteAProductbyId(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete a product by ID end point called")
	idStr := strings.TrimPrefix(r.URL.Path, "/delete-product/")
	id, err := strconv.Atoi(idStr)
	// id there is an error on the converted id
	if err != nil || id < 1 {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
	}

	index := -1
	//remove the product if the id is available
	for i, product := range m.ProductsList {
		if product.ID == id {
			index = i
			break
		}
	}
	// if the id is not available say no product
	if index == -1 {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Remove the product from the list
	m.ProductsList = append(m.ProductsList[:index], m.ProductsList[index+1:]...)
	// Return a success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted successfully"})
}

func updateAProductbyId(w http.ResponseWriter, r *http.Request) {}

// get all the request end points here
func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/all-products", getAllProducts)
	http.HandleFunc("/product/", getAProductById)
	http.HandleFunc("/create-product", createProduct)
	http.HandleFunc("/delete-product/", deleteAProductbyId)
	http.HandleFunc("/update-product/", updateAProductbyId)
	http.ListenAndServe(":8081", nil)
}

func main() {
	fmt.Printf("Hey sidra starting here!")
	handleRequests()
}
