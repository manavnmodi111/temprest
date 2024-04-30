package geolocationapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"temprest/models"
	"time"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	var foundItem models.Location
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
	var locations []models.Location

	// Iterate over the cursor and decode each document into a Location struct
	for cursor.Next(ctx) {
		var location models.Location
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
	var foundItem models.Location
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
	var updatedData models.Location
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
func CreateMembership(w http.ResponseWriter, r *http.Request) {
	// Initialize a new Membership object
	var newMember models.Membership

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
	var foundMember models.Membership

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
	var memberships []models.Membership

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
		var member models.Membership
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

func UpdateMembershipByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Decode the request body into updated membership data
	var updatedData models.Membership
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
func CreateCommunity(w http.ResponseWriter, r *http.Request) {
	// Initialize a new community variable
	var com models.Community

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
	var foundItem models.Community
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
	var communities []models.Community

	// Iterate over the cursor and decode each community document
	for cursor.Next(ctx) {
		var community models.Community
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
	var updatedData models.Community
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
	updatedCommunity := models.Community{}
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
