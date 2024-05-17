# Golang-Application-Project
The project is based on Avatar: Aang cartoon-serial.

/users method POST

/users/{userld:[0-9]+} method GET

/users/{userld:[0-9]+} method PUT

/users/{userld:[0-9]+} method DELETE

## Postgres DB structers

### Tables

#### characters

The characters table contains the following columns:

| Column Name       | Data Type | Description                                  |
|-------------------|-----------|----------------------------------------------|
| id              | INT       | Primary key, unique identifier for each character. |
| name            | VARCHAR   | The name of the character.                   |
| description     | TEXT      | A detailed description of the character.     |
| age             | INT       | The age of the character.                    |
| affiliation_id  | INT       | Foreign key linking to the affiliation table. |
| created_at      | TIMESTAMP | Timestamp when the character was created.    |
| updated_at      | TIMESTAMP | Timestamp when the character was last updated. |

#### affiliation

The affiliation table contains the following columns:

| Column Name       | Data Type | Description                                  |
|-------------------|-----------|----------------------------------------------|
| id              | INT       | Primary key, unique identifier for each affiliation. |
| name            | VARCHAR   | The name of the affiliation.                 |
| description     | TEXT      | A detailed description of the affiliation.   |
| created_at      | TIMESTAMP | Timestamp when the affiliation was created.  |
| updated_at      | TIMESTAMP | Timestamp when the affiliation was last updated. |

#### ability

The ability table contains the following columns:

| Column Name       | Data Type | Description                                  |
|-------------------|-----------|----------------------------------------------|
| id              | INT       | Primary key, unique identifier for each ability. |
| name            | VARCHAR   | The name of the ability.                     |
| description     | TEXT      | A detailed description of the ability.       |
| created_at      | TIMESTAMP | Timestamp when the ability was created.      |
| updated_at      | TIMESTAMP | Timestamp when the ability was last updated. |

### Relationships

- The characters table may have a foreign key affiliation_id that relates to the affiliation table, indicating the location affiliation of a character.

- If characters have specific bending abilities listed in the ability column, you might consider creating a separate table for abilities and establishing a many-to-many relationship between characters and ability.


## API Endpoints

### Base URL

The base URL for all API endpoints is https://your-api-domain.com.

### Characters

- Get All Characters: /characters (GET)
- Get Character by ID: /characters/{id} (GET)
- Create a New Character: /characters (POST)
- Update a Character: /characters/{id} (PUT)
- Delete a Character: /characters/{id} (DELETE)

#### Permissions

Certain endpoints require specific permissions to be accessed:

- POST /characters: Requires characters:read permission.
- PUT /characters/{id}: Requires characters:write permission.
- DELETE /characters/{id}: Requires characters:write permission.

#### GET /api/characters

Get a list of all characters.

Response:
[
  {
    "id": 1,
    "name": "Aang",
    "age": 112,
    "gender": "Male",
    "affiliation": "Air Nomads",
    "abilities": "Airbending, Energybending",
    "image": "https://example.com/aang.jpg"
  },
  // ... other characters
]



#### GET /api/characters


### Affiliation

- Get All Affiliations: /affiliations (GET)
- Get Affiliation by ID: /affiliations/{id} (GET)
- Create a New Affiliation: /affiliations (POST)
- Update an Affiliation: `/affiliations