package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"appBE/model"
	"appBE/service"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
)

func uploadProductHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one upload request")
	// check data type
	if r.Header.Get("Content-Type") != "application/json" {
        http.Error(w, "Content-Type is not application/json", http.StatusUnsupportedMediaType)
        return
    }
    // check auth
    token := r.Context().Value("user")
    claims := token.(*jwt.Token).Claims
    userIDFloat, ok := claims.(jwt.MapClaims)["userID"].(float64)
    if !ok {
        http.Error(w, "Invalid user ID", http.StatusInternalServerError)
        return
    }
    userID := int64(userIDFloat)
    
    // 1. process data
    // Parse from body of request to get a json object.
    decoder := json.NewDecoder(r.Body)
    product := model.Product{
        SellerID: userID,
        PutOutTime: time.Now(),
        Views:0,
    }

    // 2. call service to save product
    if err := decoder.Decode(&product); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        fmt.Fprintf(w, "Error decoding upload request: %v", err)
        return
    }

    if err := service.UploadProduct(&product); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // 3. response

    fmt.Fprintf(w, "Product saved successfully\n")
    fmt.Fprintf(w, "Uploaded %s by %d \n", product.Title, userID)
}

func deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one delete request")
    // check auth
    token := r.Context().Value("user")
    claims := token.(*jwt.Token).Claims
    userIDFloat, ok := claims.(jwt.MapClaims)["userID"].(float64)
	productIDStr := mux.Vars(r)["productID"]
	
    if !ok {
        http.Error(w, "Invalid user ID", http.StatusInternalServerError)
        return
    }

    // 1. process data
    userID := int64(userIDFloat)
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid product ID provided", http.StatusBadRequest)
        return
    }

    // 2. call service level to delete product
    if err := service.DeleteProduct(productID, userID); err != nil {
        http.Error(w, "Failed to delete products from backend", http.StatusInternalServerError)
        return
    }

    // 3. response
    fmt.Fprintf(w, "Product deleted successfully\n")
}

func searchProductHandler(w http.ResponseWriter, r *http.Request) {
    // 1. process data
    // description here contains both product description and product title
	fmt.Println("Received one search request")
	w.Header().Set("Content-Type", "application/json")
	description := r.URL.Query().Get("description")
	batchStr := r.URL.Query().Get("batch")
	totalSizeStr := r.URL.Query().Get("totalSize")

	batch, err := strconv.Atoi(batchStr)
	if err != nil || batch < 1 {
		batch = 1 // default to first page
	}
	totalSize, err := strconv.Atoi(totalSizeStr)
	if err != nil || totalSize < 1 {
		totalSize = 60 // default total size to load from server
	}

	var products []model.Product

	// products, err = service.SearchProductsByTitle(title)
	// if err != nil {
	// 	http.Error(w, "Failed to read Apps from backend", http.StatusInternalServerError)
	// 	return
	// }
    
    // 2. call service to handle search
	products, err = service.SearchProductsByDescription(description, batch, totalSize)
	if err != nil {
		http.Error(w, "Failed to read products from backend", http.StatusInternalServerError)
		return
	}
    
    // 3. format json response
	js, err := json.Marshal(products)
	if err != nil {
		http.Error(w, "Failed to parse products into JSON format", http.StatusInternalServerError)
		return
	}
	w.Write(js)
}

func getProductHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one get product request")

	productIDStr := mux.Vars(r)["productID"]
	
    // 1. process data
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid product ID provided", http.StatusBadRequest)
        return
    }

    // 2. call service level to get product info
    product, err := service.SearchProductByID(productID)
    if err != nil {
        // Check if the error is due to the product not being found
        if err.Error() == fmt.Sprintf("no product found with ID %d", productID) {
            http.Error(w, err.Error(), http.StatusNotFound)
        } else {
            // For all other errors, return internal server error
            http.Error(w, "Failed to search product by ID from backend", http.StatusInternalServerError)
        }
        return
    }

    // 3. format json response
    js, err := json.Marshal(product)
	if err != nil {
		http.Error(w, "Failed to parse products into JSON format", http.StatusInternalServerError)
		return
	}
	w.Write(js)
}

func getMostViewedProductsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one get most viewed products request")

	w.Header().Set("Content-Type", "application/json")
	batchStr := r.URL.Query().Get("batch")
	totalSizeStr := r.URL.Query().Get("totalSize")

	batch, err := strconv.Atoi(batchStr)
	if err != nil || batch < 1 {
		batch = 1 // default to first page
	}
	totalSize, err := strconv.Atoi(totalSizeStr)
	if err != nil || totalSize < 1 {
		totalSize = 60 // default total size to load from server
	}

	var products []model.Product

    // 2. call service to handle search
	products, err = service.GetMostViewedProducts(batch, totalSize)
	if err != nil {
		http.Error(w, "Failed to read products from backend", http.StatusInternalServerError)
		return
	}
    
    // 3. format json response
	js, err := json.Marshal(products)
	if err != nil {
		http.Error(w, "Failed to parse products into JSON format", http.StatusInternalServerError)
		return
	}
	w.Write(js)
}

func getUserProductsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one get user products request")

    w.Header().Set("Content-Type", "application/json")
	userIDStr := mux.Vars(r)["userID"]
    batchStr := r.URL.Query().Get("batch")
	totalSizeStr := r.URL.Query().Get("totalSize")

	batch, err := strconv.Atoi(batchStr)
	if err != nil || batch < 1 {
		batch = 1 // default to first page
	}
	totalSize, err := strconv.Atoi(totalSizeStr)
	if err != nil || totalSize < 1 {
		totalSize = 60 // default total size to load from server
	}
	
    // 1. process data
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid userID provided", http.StatusBadRequest)
        return
    }

    // 2. call service level to get product info
    products, err := service.GetProductsByUserID(userID, batch, totalSize)
    if err != nil {
        // Check if the error is due to the product not being found
        if err.Error() == fmt.Sprintf("no user found with ID %d", userID) {
            http.Error(w, err.Error(), http.StatusNotFound)
        } else {
            // For all other errors, return internal server error
            http.Error(w, "Failed to search user products by ID from backend", http.StatusInternalServerError)
        }
        return
    }

    // 3. format json response
    js, err := json.Marshal(products)
	if err != nil {
		http.Error(w, "Failed to parse user into JSON format", http.StatusInternalServerError)
		return
	}
	w.Write(js)
}