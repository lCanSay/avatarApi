CREATE TABLE IF NOT EXISTS character_ability (
    character_id INTEGER NOT NULL REFERENCES character ON DELETE CASCADE,
    ability_id INTEGER NOT NULL REFERENCES ability,
    PRIMARY KEY (character_id, ability_id)
);
