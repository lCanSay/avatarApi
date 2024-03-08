package models

type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Place struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Image string `json:"image"`
}

var Places = []Place{
	{Name: "Air Temple", Type: "Air Nomads", Image: "https://example.com/air_temple.jpg"},
	{Name: "Southern Water Tribe", Type: "Water Tribe", Image: "https://example.com/southern_water _tribe.jpg"},
	{Name: "Northern Water Tribe", Type: "Water Tribe", Image: "https://example.com/northern_water_tribe.jpg"},
	{Name: "Ba Sing Se", Type: "Earth Kingdom", Image: "https://example.com/ba_sing_se.jpg"},
	{Name: "Fire Nation Capital", Type: "Fire Nation", Image: "https://example.com/fire_nation_capital.jpg"},
}
