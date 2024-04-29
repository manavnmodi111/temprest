package models

type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Add more fields as needed
}
type Location struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
