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
| Rank              | string   |

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

Difficulty (int)
----------------

| Value | Description |
|-------|-------------|
|     0 | Basic       |
|     1 | Advanced    |
|     2 | Expert      |
|     3 | Master      |
|     4 | Remaster    |
|     5 | Utage       |

PlayInfo
--------

|        Field         |          Type           |
|----------------------|-------------------------|
| UserPlayDate         | int64 // Unix timestamp |
| SongId               | int                     |
| Difficulty           | Difficulty              |
| Score                | int                     |
| DxScore              | int                     |
| ComboStatus          | ComboStatus             |
| SyncStatus           | SyncStatus              |
| IsClear              | bool                    |
| IsNewRecord          | bool                    |
| IsDxNewRecord        | bool                    |
| Track                | int                     |
| MatchingUsers        | []string                |
| MaxCombo             | int                     |
| TotalCombo           | int                     |
| MaxSync              | int                     |
| TotalSync            | int                     |
| FastCount            | int                     |
| LateCount            | int                     |
| BeforeRating         | int                     |
| AfterRating          | int                     |
| TapCriticalPerfect   | int                     |
| TapPerfect           | int                     |
| TapGreat             | int                     |
| TapGood              | int                     |
| TapMiss              | int                     |
| HoldCriticalPerfect  | int                     |
| HoldPerfect          | int                     |
| HoldGreat            | int                     |
| HoldGood             | int                     |
| HoldMiss             | int                     |
| SlideCriticalPerfect | int                     |
| SlidePerfect         | int                     |
| SlideGreat           | int                     |
| SlideGood            | int                     |
| SlideMiss            | int                     |
| TouchCriticalPerfect | int                     |
| TouchPerfect         | int                     |
| TouchGreat           | int                     |
| TouchGood            | int                     |
| TouchMiss            | int                     |
| BreakCriticalPerfect | int                     |
| BreakPerfect         | int                     |
| BreakGreat           | int                     |
| BreakGood            | int                     |
| BreakMiss            | int                     |
| TotalCriticalPerfect | int                     |
| TotalPerfect         | int                     |
| TotalGreat           | int                     |
| TotalGood            | int                     |
| TotalMiss            | int                     |

ComboStatus (int)
-----------------

| Value |  Description   |
|-------|----------------|
|     0 | NoCombo        |
|     1 | FullCombo      |
|     2 | FullComboPlus  |
|     3 | AllPerfect     |
|     4 | AllPerfectPlus |

SyncStatus (int)
----------------

| Value |  Description   |
|-------|----------------|
|     0 | NoSync         |
|     1 | FullSync       |
|     2 | FullSyncPlus   |
|     3 | FullSyncDx     |
|     4 | FullSyncDxPlus |

# Errors

| Status Code | Description  |
|-------------|--------------|
|         200 | OK           |
|         400 | Bad Request  |
|         404 | Not Found    |
|         500 | Server Error |
