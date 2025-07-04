package models

type Secret struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
}
