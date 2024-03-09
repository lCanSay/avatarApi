# Golang-Application-Project
The project is based on Avatar: Aang cartoon-serial.

Khan Kensey 22B030608 \n
Kenzhegul Rashid 22B030148

/users method POST

/users/{userld:[0-9]+} method GET

/users/{userld:[0-9]+} method PUT

/users/{userld:[0-9]+} method DELETE

## Postgres DB structers

### Tables

#### `characters`

| Column       | Type    | Description                       |
|--------------|---------|-----------------------------------|
| id           | SERIAL  | Unique identifier for the character|
| name         | VARCHAR | Full name of the character        |
| age          | INTEGER | Age of the character              |
| gender       | VARCHAR | Gender of the character           |
| affiliation  | VARCHAR | Affiliation of the character      |
| abilities    | VARCHAR | Special abilities of the character|
| image        | VARCHAR | URL of the character's image      |

#### `places`

| Column | Type    | Description                   |
|--------|---------|-------------------------------|
| id     | SERIAL  | Unique identifier for the place|
| name   | VARCHAR | Name of the location           |
| type   | VARCHAR | Type or affiliation of the place|
| image  | VARCHAR | URL of the place's image       |

#### `elements`

| Column      | Type    | Description                    |
|-------------|---------|--------------------------------|
| id          | SERIAL  | Unique identifier for the element|
| name        | VARCHAR | Name of the bending element     |
| description | TEXT    | Description of the bending element|
| image       | VARCHAR | URL of the element's image      |

### Relationships

- The `characters` table may have a foreign key `place_id` that relates to the `places` table, indicating the location affiliation of a character.

- If characters have specific bending abilities listed in the `abilities` column, you might consider creating a separate table for abilities and establishing a many-to-many relationship between `characters` and `abilities`.

Feel free to adjust the table structure and relationships based on your project's specific needs.

## API Endpoints

### Base URL

The base URL for all API endpoints is `https://your-api-domain.com`.

### Characters

#### `GET /api/characters`

Get a list of all characters.

**Response:**
```json
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
```

### Places

#### `GET /api/places`

Get a list of all places.

**Response:**
```json
[
  {
    "name": "Air Temple",
    "type": "Air Nomads",
    "image": "https://example.com/air_temple.jpg"
  },
  // ... other places
]
```

### Elements

#### `GET /api/elements`

Get a list of all elements.

**Response:**
```json
[
  {
    "name": "Airbending",
    "discription": "Manipulation of air currents",
    "image": "https://example.com/airbending.jpg"
  },
  // ... other elements
]
```
