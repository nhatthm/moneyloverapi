// Code generated by github.com/swaggest/json-cli v1.8.3, DO NOT EDIT.

// Package category contains JSON mapping structures.
package category

// Category structure is generated from "#/components/schemas/Category".
type Category struct {
	ID       string    `json:"_id"`  // Required.
	Name     string    `json:"name"` // Required.
	Icon     string    `json:"icon,omitempty"`
	Account  string    `json:"account,omitempty"`
	Type     int64     `json:"type,omitempty"`
	Metadata string    `json:"metadata,omitempty"`
	Parent   *Category `json:"parent,omitempty"`
}

// FlatCategory structure is generated from "#/components/schemas/FlatCategory".
type FlatCategory struct {
	ID       string `json:"_id"`  // Required.
	Name     string `json:"name"` // Required.
	Icon     string `json:"icon,omitempty"`
	Account  string `json:"account,omitempty"`
	Type     int64  `json:"type,omitempty"`
	Metadata string `json:"metadata,omitempty"`
	Parent   string `json:"parent,omitempty"`
}
