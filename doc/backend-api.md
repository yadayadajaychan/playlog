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
