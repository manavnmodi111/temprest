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

// Endpoints For Location
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

func UpdateLocationByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var updatedItem *models.Location
	var updatedItemIndex int
	for i, item := range location {
		if item.ID == id {
			updatedItem = &location[i]
			updatedItemIndex = i
			break
		}
	}

	if updatedItem == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Item not found"))
		return
	}

	var updatedData models.Location
	err := json.NewDecoder(r.Body).Decode(&updatedData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding request body: %v", err)
		return
	}

	// Update item fields
	updatedItem.Name = updatedData.Name
	// Update other fields as needed

	// Marshal updated item to JSON
	jsonData, err := json.Marshal(updatedItem)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error marshaling JSON"))
		return
	}

	// Update item in the slice
	location[updatedItemIndex] = *updatedItem

	// Set response status code and header
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	// Write JSON response
	w.Write(jsonData)
}

func DeleteLocationByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var found bool
	for i, item := range location {
		if item.ID == id {
			// Remove item from the slice
			location = append(location[:i], location[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Item not found"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Endpoints For Membership
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

func UpdateMembershipByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var updatedItem *models.Membership
	var updatedItemIndex int
	for i, item := range membership {
		if item.ID == id {
			updatedItem = &membership[i]
			updatedItemIndex = i
			break
		}
	}

	if updatedItem == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Item not found"))
		return
	}

	var updatedData models.Membership
	err := json.NewDecoder(r.Body).Decode(&updatedData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding request body: %v", err)
		return
	}

	// Update item fields
	updatedItem.Role = updatedData.Role
	// Update other fields as needed

	// Marshal updated item to JSON
	jsonData, err := json.Marshal(updatedItem)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error marshaling JSON"))
		return
	}

	// Update item in the slice
	membership[updatedItemIndex] = *updatedItem

	// Set response status code and header
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	// Write JSON response
	w.Write(jsonData)
}

func DeleteMembershipByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var found bool
	for i, item := range membership {
		if item.ID == id {
			// Remove item from the slice
			membership = append(membership[:i], membership[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Item not found"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Endpoints For Community
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

func UpdateCommunityByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var updatedItem *models.Community
	var updatedItemIndex int
	for i, item := range community {
		if item.ID == id {
			updatedItem = &community[i]
			updatedItemIndex = i
			break
		}
	}

	if updatedItem == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Item not found"))
		return
	}

	var updatedData models.Community
	err := json.NewDecoder(r.Body).Decode(&updatedData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding request body: %v", err)
		return
	}

	// Update item fields
	updatedItem.Name = updatedData.Name
	// Update other fields as needed

	// Marshal updated item to JSON
	jsonData, err := json.Marshal(updatedItem)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error marshaling JSON"))
		return
	}

	// Update item in the slice
	community[updatedItemIndex] = *updatedItem

	// Set response status code and header
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	// Write JSON response
	w.Write(jsonData)
}

func DeleteCommunityByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var found bool
	for i, item := range community {
		if item.ID == id {
			// Remove item from the slice
			community = append(community[:i], community[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Item not found"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
