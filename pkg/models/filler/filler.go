package filler

import (
	model "github.com/lCanSay/avatarApi/pkg/models"
)

func PopulateDatabase(models model.Models) error {
	for _, character := range characters {
		models.Characters.Insert(&character, 1)
	}
	// TODO: Implement restaurants pupulation
	// TODO: Implement the relationship between restaurants and menus
	return nil
}

var characters = []model.Character{
	{Name: "Kensey", Age: 112, Gender: "Male", Abilities: "Airbending, Energybending", Image: "https://example.com/aang.jpg", Affiliation_id: 1},
	{Name: "Kensey2", Age: 16, Gender: "Female", Abilities: "Waterbending, Healing", Image: "https://example.com/katara.jpg", Affiliation_id: 3},
}
