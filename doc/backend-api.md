# Endpoints

`GET /api/playlog`
------------------
- **Description**: Retrieve the playlog
- **Query Parameters**:

|   Name    | Type |                 Description                 | Required | Default |
|-----------|------|---------------------------------------------|----------|---------|
| ascending | bool | sort dates by ascending or descending order | no       | false   |
| page      | int  | page no.                                    | no       | 1       |
| count     | int  | no. of entries per page                     | no       | 50      |

- **JSON Response**:

|  Field  |      Type      |
|---------|----------------|
| MaxPage | int            |
| Playlog | []playlogEntry |

- **playlogEntry**:

|       Field       |   Type   |
|-------------------|----------|
| SongInfo          | SongInfo |
| PlayInfo          | PlayInfo |
| PreviousBestScore | int      |

# Types

SongInfo
--------

|  Field   |    Type     |
|----------|-------------|
| SongId   | int         |
| Name     | string      |
| Artist   | string      |
| Type     | string      |
| Bpm      | int         |
| Category | string      |
| Version  | string      |
| Sort     | string      |
| Charts   | []ChartInfo |

ChartInfo
---------

|     Field     |          Type           |
|---------------|-------------------------|
| Difficulty    | Difficulty              |
| Level         | int                     |
| InternalLevel | int // multiplied by 10 |
| NotesDesigner | string                  |
| MaxNotes      | int                     |

Difficulty
----------

| Value | Description |
|-------|-------------|
|     0 | basic       |
|     1 | advanced    |
|     2 | expert      |
|     3 | master      |
|     4 | remaster    |
|     5 | utage       |
