package geolocationapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"temprest/models"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

var items []models.Item
var location []models.Location

// @title Sample API
// @version 1.0
// @description This is a sample API
// @BasePath /api/v1
// @host localhost:8080
// @schemes http
// @produce json

// @Summary Get an item by ID
// @Description Get an item by its name
// @ID get-item-by-id
// @Param name path string true "Item name"
// @Success 200 {object} Item
// @Router /items/{name} [get]
func GetItemByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request:", r)
	name := chi.URLParam(r, "name")
	fmt.Println("ID:", name)

	// Find the item by ID in the slice
	var foundItem models.Item
	found := false
	for _, item := range items {
		if item.Name == name {
			foundItem = item
			found = true
			break
		}
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Item not found"))
		return
	}

	// Marshal item to JSON
	jsonData, err := json.Marshal(foundItem)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error marshaling JSON"))
		return
	}

	// Set response header
	w.Header().Set("Content-Type", "application/json")
	// Write JSON response
	w.Write(jsonData)
}

// @Summary Get a location by ID
// @Description Get a location by its ID
// @ID get-location-by-id
// @Param id path string true "Location ID"
// @Success 200 {object} Location
// @Router /location/{id} [get]
func GetLocationByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request:", r)
	id := chi.URLParam(r, "id")
	fmt.Println("ID:", id)

	// Find the item by ID in the slice
	var foundItem models.Location
	found := false
	for _, item := range location {
		if item.ID == id {
			foundItem = item
			found = true
			break
		}
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Item not found"))
		return
	}

	// Marshal item to JSON
	jsonData, err := json.Marshal(foundItem)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error marshaling JSON"))
		return
	}

	// Set response header
	w.Header().Set("Content-Type", "application/json")
	// Write JSON response
	w.Write(jsonData)
}

// @Summary Get all locations
// @Description Get all locations
// @ID get-all-locations
// @Success 200 {array} Location
// @Router /location [get]
func GetLocation(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(location)
}

// @Summary Create a new item
// @Description Create a new item
// @ID create-item
// @Param request body Item true "Item object"
// @Success 201 {object} Item
// @Router /items [post]
func CreateItem(w http.ResponseWriter, r *http.Request) {
	var newItem models.Item
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding request body: %v", err)
		return
	}

	// Generate a new UUID for the item
	newItem.ID = uuid.New().String()

	// Add the new item to the items slice
	items = append(items, newItem)

	// Marshal item to JSON
	jsonData, err := json.Marshal(newItem)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error marshaling JSON"))
		return
	}

	// Set response status code and header
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	// Write JSON response
	w.Write(jsonData)
}

// @Summary Create a new location
// @Description Create a new location
// @ID create-location
// @Param request body Location true "Location object"
// @Success 201 {object} Location
// @Router /location [post]
func CreateLocation(w http.ResponseWriter, r *http.Request) {
	var newItem models.Location
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding request body: %v", err)
		return
	}

	// Generate a new UUID for the item
	//newItem.ID = uuid.New().String()

	// Add the new item to the items slice
	location = append(location, newItem)

	// Marshal item to JSON
	jsonData, err := json.Marshal(newItem)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error marshaling JSON"))
		return
	}

	// Set response status code and header
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	// Write JSON response
	w.Write(jsonData)
}
