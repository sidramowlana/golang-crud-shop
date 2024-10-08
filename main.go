package main

import (
	"encoding/json"
	"fmt"
	m "golang-crud/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/afex/hystrix-go/hystrix"
)

// The idea is to execute some extra steps before or after calling the original handler,
// or even replace the original handler's behavior in case of specific conditions
func HystrixHandler(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hystrix.Do("command_name", func() error {
			handlerFunc(w, r) //calling the original handler
			return nil
		},
			func(err error) error {
				http.Error(w, "Service is unavailable", http.StatusServiceUnavailable)
				return nil
			})
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "End point to the Home Page")
}

func getAllProducts(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Get all products end point called")
	json.NewEncoder(w).Encode(m.ProductsList)
}

func getProductByID(id int) (*m.Product, int) {
	for i, product := range m.ProductsList {
		if product.ID == id {
			return &m.ProductsList[i], i
		}
	}
	return nil, -1
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
	product, _ := getProductByID(id)
	if product == nil {
		http.Error(w, "Product ID cannot be found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)

}

func createProduct(w http.ResponseWriter, r *http.Request) {
	// Make sure if the end point method is post
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid reuqest method", http.StatusMethodNotAllowed)
		return
	}

	// if it is a post method u have to decode the body to product struct
	var newProduct m.Product
	// Get the body from the request and assign it to the new variable
	err := json.NewDecoder(r.Body).Decode(&newProduct)

	//Check if the body is valid
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
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

func updateAProductbyId(w http.ResponseWriter, r *http.Request) {
	// make sure the request is put
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request paylod", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/update-product/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid request paylod", http.StatusMethodNotAllowed)
		return
	}
	//search for the product id and get the product details -> call the getProductById
	product, index := getProductByID(id)
	if product == nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	//decode the body into product struct
	var updateProduct m.Product

	//Check if the body is valid
	err = json.NewDecoder(r.Body).Decode(&updateProduct)
	if err != nil {
		http.Error(w, "Invalid request paylod", http.StatusMethodNotAllowed)
		return
	}

	//update the product details
	m.ProductsList[index].Name = updateProduct.Name
	m.ProductsList[index].Description = updateProduct.Description
	m.ProductsList[index].Price = updateProduct.Price
	m.ProductsList[index].Quantity = updateProduct.Quantity

	// Return the updated product
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(m.ProductsList[index])
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

// get all the request end points here
func handleRequests() {
	http.HandleFunc("/", HystrixHandler(homePage))
	http.HandleFunc("/all-products", HystrixHandler(getAllProducts))
	http.HandleFunc("/product/", HystrixHandler(getAProductById))
	http.HandleFunc("/create-product", HystrixHandler(createProduct))
	http.HandleFunc("/delete-product/", HystrixHandler(deleteAProductbyId))
	http.HandleFunc("/update-product/", HystrixHandler(updateAProductbyId))
	http.ListenAndServe(":8081", nil)
}

func main() {
	fmt.Printf("Hey sidra starting here!")
	handleRequests()
}
