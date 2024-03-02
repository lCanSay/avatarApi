package models

type Character struct {
	Id          int    `json:"id"`
	Name        string `json:"fname"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Affiliation string `json:"affiliation"`
	Abilities   string `json:"abilities"` // elements or technics
	Image       string `json:"image"`
}

var Characters = []Character{
	{Id: 1, Name: "Aang", Age: 112, Gender: "Male", Affiliation: "Air Nomads", Abilities: "Airbending, Energybending", Image: "https://example.com/aang.jpg"},
	{Id: 2, Name: "Katara", Age: 14, Gender: "Female", Affiliation: "Water Tribe", Abilities: "Waterbending, Healing", Image: "https://example.com/katara.jpg"},
	{Id: 3, Name: "Zuko", Age: 16, Gender: "Male", Affiliation: "Fire Nation", Abilities: "Firebending", Image: "https://example.com/zuko.jpg"},
	// will add more later
}
