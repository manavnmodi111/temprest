package geolocationapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"temprest/models"

	"github.com/go-chi/chi"
)

var location []models.Location
var membership []models.Membership
var community []models.Community

func CreateLocation(w http.ResponseWriter, r *http.Request) {
	var newItem models.Location
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding request body: %v", err)
		return
	}

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

func GetLocation(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(location)
}

func CreateMembership(w http.ResponseWriter, r *http.Request) {
	var members models.Membership
	err := json.NewDecoder(r.Body).Decode(&members)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding request body: %v", err)
		return
	}

	// Add the new item to the items slice
	membership = append(membership, members)

	// Marshal item to JSON
	jsonData, err := json.Marshal(members)
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

func GetMembershipByID(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	// Find the item by ID in the slice
	var foundItem models.Membership
	found := false
	for _, item := range membership {
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

func GetMembership(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(membership)
}

func CreateCommunity(w http.ResponseWriter, r *http.Request) {
	var com models.Community
	err := json.NewDecoder(r.Body).Decode(&com)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding request body: %v", err)
		return
	}

	// Add the new item to the items slice
	community = append(community, com)

	// Marshal item to JSON
	jsonData, err := json.Marshal(com)
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

func GetCommunityByID(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	// Find the item by ID in the slice
	var foundItem models.Community
	found := false
	for _, item := range community {
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

func GetCommunity(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(community)
}
