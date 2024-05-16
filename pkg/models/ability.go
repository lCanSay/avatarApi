package models

type Ability struct {
	Name        string `json:"name"`
	Element     string `json:"element`
	Description string `json:"description"`
	Image       string `json:"image"`
}

var Abilities = []Ability{
	{Name: "Airbending", Description: "Manipulation of air currents", Image: "https://example.com/airbending.jpg"},
	{Name: "Waterbending", Description: "Manipulation of water in various forms", Image: "https://example.com/waterbending.jpg"},
	{Name: "Earthbending", Description: "Manipulation earth and rock", Image: "https://example.com/earthbending.jpg"},
	{Name: "Firebending", Description: "Manipulation of fire and lightning", Image: "https://example.com/firebending.jpg"},
}
