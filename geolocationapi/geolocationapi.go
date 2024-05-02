package geolocationapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"time"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Location represents a geographical location
type Location struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Membership represents a memebership
type Membership struct {
	ID          string `json:"id"`
	CommunityID string `json:"communityId"`
	Role        string `json:"role"`
}

// Community represents a community
type Community struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Location Location     `json:"location"`
	Members  []Membership `json:"members"`
}

var location []Location
var membership []Membership
var community []Community

// Endpoints For Location

// CreateLocation godoc
// @Summary Create a new location
// @Description Creates a new location and adds it to the MongoDB collection
// @Tags locations
// @Accept json
// @Produce json
// @Param Location body Location true "Location object to be created"
// @Success 201 {object} Location "location created"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /geolocationapi/location [post]
func CreateLocation(w http.ResponseWriter, r *http.Request) {
	var newItem Location
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding request body: %v", err)
		return
	}

	// Add the new item to the items slice
	location = append(location, newItem)

	collection := client.Database("geolocapi").Collection("locations")
	// Insert a new document into the collection
	_, err = collection.InsertOne(ctx, newItem)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error inserting into database: %v", err)
		return
	}
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

// GetLocationByID godoc
// @Summary Get a location by ID
// @Description Retrieves a location from the MongoDB collection by its ID
// @Tags locations
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} Location "location found"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /geolocationapi/location/{id} [get]
func GetLocationByID(w http.ResponseWriter, r *http.Request) {
	// Get the ID parameter from the URL
	id := chi.URLParam(r, "id")

	// Check if the MongoDB client is nil
	if client == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "MongoDB client is not initialized")
		return
	}

	// Get the MongoDB context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the locations collection
	collection := client.Database("geolocapi").Collection("locations")

	// Find the document by ID in the collection
	var foundItem Location
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&foundItem)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// If no document is found, return 404 Not Found
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Item not found"))
			return
		}
		// If an error occurs during the find operation, return internal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error retrieving location: %v", err)
		return
	}

	// Marshal item to JSON
	jsonData, err := json.Marshal(foundItem)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marshaling JSON: %v", err)
		return
	}

	// Set response header
	w.Header().Set("Content-Type", "application/json")
	// Write JSON response
	w.Write(jsonData)
}

// GetLocation godoc
// @Summary Get all locations
// @Description Retrieves all locations from the MongoDB collection
// @Tags locations
// @Accept  json
// @Produce  json
// @Success 200 {object} []Location
// @Failure 500 {string} string "Internal Server Error"
// @Router /geolocationapi/location [get]
func GetLocation(w http.ResponseWriter, r *http.Request) {
	// Check if the MongoDB client is nil
	if client == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "MongoDB client is not initialized")
		return
	}

	// Get the MongoDB context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the locations collection
	collection := client.Database("geolocapi").Collection("locations")

	// Find all documents in the collection
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error retrieving locations: %v", err)
		return
	}
	defer cursor.Close(ctx)

	// Define a slice to store retrieved locations
	var locations []Location

	// Iterate over the cursor and decode each document into a Location struct
	for cursor.Next(ctx) {
		var location Location
		if err := cursor.Decode(&location); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error decoding location: %v", err)
			return
		}
		locations = append(locations, location)
	}

	// Check if any error occurred during iteration
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error iterating over locations: %v", err)
		return
	}

	// Marshal the slice of locations to JSON and send it in the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(locations); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error encoding JSON response: %v", err)
		return
	}
}

// UpdateLocationbyID godoc
// @Summary Update a location by ID
// @Description Updates a location name in the MongoDB collection by its ID
// @Tags locations
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param updateData body Location true "Updated location data"
// @Success 200 {object} Location "location updated"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /geolocationapi/location/{id} [put]
func UpdateLocationByID(w http.ResponseWriter, r *http.Request) {
	// Get the ID parameter from the URL
	id := chi.URLParam(r, "id")

	// Check if the MongoDB client is nil
	if client == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "MongoDB client is not initialized")
		return
	}

	// Get the MongoDB context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the locations collection
	collection := client.Database("geolocapi").Collection("locations")

	// Find the document by ID in the collection
	var foundItem Location
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&foundItem)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// If no document is found, return 404 Not Found
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Item not found"))
			return
		}
		// If an error occurs during the find operation, return internal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error retrieving location: %v", err)
		return
	}

	// Decode the request body into updatedData
	var updatedData Location
	err = json.NewDecoder(r.Body).Decode(&updatedData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding request body: %v", err)
		return
	}

	// Update the foundItem fields
	foundItem.Name = updatedData.Name
	// Update other fields as needed

	// Update the document in the collection
	_, err = collection.ReplaceOne(ctx, bson.M{"id": id}, foundItem)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error updating location: %v", err)
		return
	}

	// Marshal updated item to JSON
	jsonData, err := json.Marshal(foundItem)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marshaling JSON: %v", err)
		return
	}

	// Set response status code and header
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	// Write JSON response
	w.Write(jsonData)
}

// DeleteLocationByID godoc
// @Summary Delete a location by ID
// @Description Deletes a location from the MongoDB collection by its ID
// @Tags locations
// @Param id path string true "ID"
// @Success 204 "No Content"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Int
// @Router /geolocationapi/location/{id} [delete]
func DeleteLocationByID(w http.ResponseWriter, r *http.Request) {
	// Get the ID parameter from the URL
	id := chi.URLParam(r, "id")

	// Check if the MongoDB client is nil
	if client == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "MongoDB client is not initialized")
		return
	}

	// Get the MongoDB context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the locations collection
	collection := client.Database("geolocapi").Collection("locations")

	// Delete the document by ID from the collection
	_, err := collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// If no document is found, return 404 Not Found
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Item not found"))
			return
		}
		// If an error occurs during the delete operation, return internal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting location: %v", err)
		return
	}

	// Set response status code
	w.WriteHeader(http.StatusNoContent)
}

// Endpoints For Membership

// CreateMembership godoc
// @Summary Create a new membership
// @Description Creates a new membership and adds it to the MongoDB collection
// @Tags membership
// @Accept json
// @Produce json
// @Param Membership body Membership true "Membership object to be created"
// @Success 201 {object} Membership "membership created"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /geolocationapi/membership [post]
func CreateMembership(w http.ResponseWriter, r *http.Request) {
	// Initialize a new Membership object
	var newMember Membership

	// Decode the request body into the newMember variable
	err := json.NewDecoder(r.Body).Decode(&newMember)
	if err != nil {
		// If there's an error decoding the request body, return a bad request response
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding request body: %v", err)
		return
	}

	// Check if the MongoDB client is nil
	if client == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "MongoDB client is not initialized")
		return
	}

	// Get the MongoDB context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the memberships collection
	collection := client.Database("geolocapi").Collection("memberships")

	// Insert the newMember document into the memberships collection
	_, err = collection.InsertOne(ctx, newMember)
	if err != nil {
		// If an error occurs during the insert operation, return internal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating membership: %v", err)
		return
	}

	// Marshal newMember to JSON
	jsonData, err := json.Marshal(newMember)
	if err != nil {
		// If an error occurs during marshaling, return internal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marshaling JSON: %v", err)
		return
	}

	// Set response status code and header
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	// Write JSON response
	w.Write(jsonData)
}

// GetMembershipByID godoc
// @Summary Get a membership by ID
// @Description Retrieves a membership from the MongoDB collection by its ID
// @Tags membership
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} Membership "membership found"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /geolocationapi/membership/{id} [get]
func GetMembershipByID(w http.ResponseWriter, r *http.Request) {
	// Get the ID parameter from the URL
	id := chi.URLParam(r, "id")

	// Check if the MongoDB client is nil
	if client == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "MongoDB client is not initialized")
		return
	}

	// Get the MongoDB context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the memberships collection
	collection := client.Database("geolocapi").Collection("memberships")

	// Initialize a Membership object to store the found document
	var foundMember Membership

	// Find the document by ID in the memberships collection
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&foundMember)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// If no document is found, return 404 Not Found
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Item not found"))
			return
		}
		// If an error occurs during the find operation, return internal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error finding membership: %v", err)
		return
	}

	// Marshal the foundMember to JSON
	jsonData, err := json.Marshal(foundMember)
	if err != nil {
		// If an error occurs during marshaling, return internal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marshaling JSON: %v", err)
		return
	}

	// Set response header
	w.Header().Set("Content-Type", "application/json")
	// Write JSON response
	w.Write(jsonData)
}

// GetMembership godoc
// @Summary Get all membership
// @Description Retrieves all membership from the MongoDB collection
// @Tags membership
// @Accept  json
// @Produce  json
// @Success 200 {object} []Membership
// @Failure 500 {string} string "Internal Server Error"
// @Router /geolocationapi/membership [get]
func GetMembership(w http.ResponseWriter, r *http.Request) {
	// Check if the MongoDB client is nil
	if client == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "MongoDB client is not initialized")
		return
	}

	// Get the MongoDB context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the memberships collection
	collection := client.Database("geolocapi").Collection("memberships")

	// Initialize a slice to store Membership objects
	var memberships []Membership

	// Find all documents in the memberships collection
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		// If an error occurs during the find operation, return internal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error finding memberships: %v", err)
		return
	}
	defer cursor.Close(ctx)

	// Iterate over the cursor and decode each document into a Membership object
	for cursor.Next(ctx) {
		var member Membership
		if err := cursor.Decode(&member); err != nil {
			// If an error occurs during decoding, return internal server error
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error decoding membership: %v", err)
			return
		}
		memberships = append(memberships, member)
	}
	if err := cursor.Err(); err != nil {
		// If an error occurs with the cursor, return internal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error with cursor: %v", err)
		return
	}

	// Marshal the memberships slice to JSON
	jsonData, err := json.Marshal(memberships)
	if err != nil {
		// If an error occurs during marshaling, return internal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marshaling JSON: %v", err)
		return
	}

	// Set response header
	w.Header().Set("Content-Type", "application/json")
	// Write JSON response
	w.Write(jsonData)
}

// UpdateMembershipbyID godoc
// @Summary Update a membership by ID
// @Description Updates a membership role in the MongoDB collection by its ID
// @Tags membership
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param updateData body Membership true "Updated Membership data"
// @Success 200 {object} Membership "Membership updated"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /geolocationapi/membership/{id} [put]
func UpdateMembershipByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Decode the request body into updated membership data
	var updatedData Membership
	err := json.NewDecoder(r.Body).Decode(&updatedData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding request body: %v", err)
		return
	}

	// Check if the MongoDB client is nil
	if client == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "MongoDB client is not initialized")
		return
	}

	// Get the MongoDB context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the memberships collection
	collection := client.Database("geolocapi").Collection("memberships")

	// Update the membership in the database
	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"role": updatedData.Role}} // Update the role field, you can add more fields as needed
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error updating membership: %v", err)
		return
	}

	// Respond with updated membership data
	updatedData.ID = id // Ensure that the ID remains the same
	jsonData, err := json.Marshal(updatedData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marshaling JSON: %v", err)
		return
	}

	// Set response status code and header
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	// Write JSON response
	w.Write(jsonData)
}

// DeleteMembershipByID godoc
// @Summary Delete a membership by ID
// @Description Deletes a location from the MongoDB collection by its ID
// @Tags membership
// @Param id path string true "ID"
// @Success 204 "No Content"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Int
// @Router /geolocationapi/membership/{id} [delete]
func DeleteMembershipByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Check if the MongoDB client is nil
	if client == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "MongoDB client is not initialized")
		return
	}

	// Get the MongoDB context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the memberships collection
	collection := client.Database("geolocapi").Collection("memberships")

	// Delete the document by ID from the collection
	_, err := collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// If no document is found, return 404 Not Found
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Item not found"))
			return
		}
		// If an error occurs during the delete operation, return internal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting membership: %v", err)
		return
	}

	// Set response status code
	w.WriteHeader(http.StatusNoContent)
}

// Endpoints For Community

// CreateCommunity godoc
// @Summary Create a new community
// @Description Creates a new community and adds it to the MongoDB collection
// @Tags Community
// @Accept json
// @Produce json
// @Param community body Community true "Community object to be created"
// @Success 201 {object} Community "community created"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /geolocationapi/community [post]
func CreateCommunity(w http.ResponseWriter, r *http.Request) {
	// Initialize a new community variable
	var com Community

	// Decode the request body into the community variable
	err := json.NewDecoder(r.Body).Decode(&com)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding request body: %v", err)
		return
	}

	// Check if the MongoDB client is nil
	if client == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "MongoDB client is not initialized")
		return
	}

	// Get the MongoDB context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the communities collection
	collection := client.Database("geolocapi").Collection("communities")

	// Insert the newMember document into the memberships collection
	_, err = collection.InsertOne(ctx, com)
	if err != nil {
		// If an error occurs during the insert operation, return internal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating communities: %v", err)
		return
	}

	// Marshal newMember to JSON
	jsonData, err := json.Marshal(com)
	if err != nil {
		// If an error occurs during marshaling, return internal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marshaling JSON: %v", err)
		return
	}

	// Set response status code and header
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	// Write JSON response
	w.Write(jsonData)
}

// GetCommunityByID godoc
// @Summary Get a community by ID
// @Description Retrieves a community from the MongoDB collection by its ID
// @Tags Community
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} Community "community found"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /geolocationapi/community/{id} [get]
func GetCommunityByID(w http.ResponseWriter, r *http.Request) {
	// Get the community ID from the URL parameter
	id := chi.URLParam(r, "id")

	// Check if the MongoDB client is nil
	if client == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "MongoDB client is not initialized")
		return
	}

	// Get the MongoDB context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the communities collection
	collection := client.Database("geolocapi").Collection("communities")

	// Find the community document by its ID
	var foundItem Community
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&foundItem)
	if err != nil {
		// Check if the error is due to document not found
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Item not found"))
			return
		}
		// Handle other errors
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error finding community document: %v", err)
		return
	}

	// Marshal the found community document to JSON
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

// GetCommunity godoc
// @Summary Get all Community
// @Description Retrieves all Community from the MongoDB collection
// @Tags Community
// @Accept  json
// @Produce  json
// @Success 200 {object} []Community
// @Failure 500 {string} string "Internal Server Error"
// @Router /geolocationapi/community [get]
func GetCommunity(w http.ResponseWriter, r *http.Request) {
	// Check if the MongoDB client is nil
	if client == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "MongoDB client is not initialized")
		return
	}

	// Get the MongoDB context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the communities collection
	collection := client.Database("geolocapi").Collection("communities")

	// Define a filter to find all community documents
	filter := bson.M{}

	// Find all community documents that match the filter
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error finding community documents: %v", err)
		return
	}
	defer cursor.Close(ctx)

	// Define a slice to store the retrieved community documents
	var communities []Community

	// Iterate over the cursor and decode each community document
	for cursor.Next(ctx) {
		var community Community
		if err := cursor.Decode(&community); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error decoding community document: %v", err)
			return
		}
		communities = append(communities, community)
	}

	// Check for errors during cursor iteration
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error iterating over community documents: %v", err)
		return
	}

	// Marshal the retrieved community documents to JSON
	jsonData, err := json.Marshal(communities)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marshaling JSON: %v", err)
		return
	}

	// Set response header
	w.Header().Set("Content-Type", "application/json")
	// Write JSON response
	w.Write(jsonData)
}

// UpdateCommunityByID godoc
// @Summary Update a community by ID
// @Description Updates a community name in the MongoDB collection by its ID
// @Tags Community
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param updatedCommunity body Community true "Updated community object"
// @Success 200 {object} Community "community updated"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /geolocationapi/community/{id} [put]
func UpdateCommunityByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Check if the MongoDB client is nil
	if client == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "MongoDB client is not initialized")
		return
	}

	// Get the MongoDB context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the communities collection
	collection := client.Database("geolocapi").Collection("communities")

	// Define a filter to find the community document by ID
	filter := bson.M{"id": id}

	// Define an update to set the fields of the updated community document
	var updatedData Community
	err := json.NewDecoder(r.Body).Decode(&updatedData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding request body: %v", err)
		return
	}

	update := bson.M{
		"$set": bson.M{
			"name": updatedData.Name,
			// Update other fields as needed
		},
	}

	// Perform the update operation
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error updating community document: %v", err)
		return
	}

	// Get the updated community document
	updatedCommunity := Community{}
	err = collection.FindOne(ctx, filter).Decode(&updatedCommunity)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error finding updated community document: %v", err)
		return
	}

	// Marshal the updated community document to JSON
	jsonData, err := json.Marshal(updatedCommunity)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marshaling JSON: %v", err)
		return
	}

	// Set response status code and header
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	// Write JSON response
	w.Write(jsonData)
}

// DeleteCommunityByID godoc
// @Summary Delete a community by ID
// @Description Deletes a community from the MongoDB collection by its ID
// @Tags Community
// @Param id path string true "ID"
// @Success 204 "No Content"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /geolocationapi/community/{id} [delete]
func DeleteCommunityByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Check if the MongoDB client is nil
	if client == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "MongoDB client is not initialized")
		return
	}

	// Get the MongoDB context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the communities collection
	collection := client.Database("geolocapi").Collection("communities")

	// Define a filter to find the community document by ID
	filter := bson.M{"id": id}

	// Delete the community document
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting community document: %v", err)
		return
	}

	// Set response status code
	w.WriteHeader(http.StatusNoContent)
}
