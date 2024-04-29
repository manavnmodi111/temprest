package models

// Add more fields as needed

type Location struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Membership struct {
	ID          string `json:"id"`
	CommunityID string `json:"communityId"`
	Role        string `json:"role"`
}
type Community struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Location Location     `json:"location"`
	Members  []Membership `json:"members"`
}
