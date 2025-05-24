# Playlog

This program fetchs maimai play data from solips.app and saves it in an sqlite3 database.
This data is then served over an HTTP API.

## Compiling

Make both the frontend and backend:
```
$ make
```

Make only backend:
```
$ make backend
```

Make only frontend:
```
$ make frontend
```

The backend output is `playlog` and the frontend output is `build/`

## Running

### Backend

Get usage info by specifying `-h`:
```
$ ./playlog -h
Usage: playlog [-bhuv] [-a value] [-l value] [-p value] [-s value] [-t value] [parameters ...]
 -a, --api-interval=value
                    seconds to wait between api requests [3]
 -b, --backend-only
                    only run the backend {action}
 -h, --help         display help
 -l, --listen-port=value
                    port to listen on [5000]
 -p, --playdb=value
                    filename of play db [plays.db]
 -s, --songdb=value
                    filename of song db [songs.db]
 -t, --update-interval=value
                    seconds to wait between updates [900]
 -u, --update-only  only update the play db & exit {action}
 -v, --verbose      verbosity level (errors only, info, debug) [0]
```

#### Examples

Run the backend while updating the play database every 1000 seconds:
```
$ ./playlog -vt 1000
```

Only update the play database very verbosely,
waiting 10 seconds between each api request to solips:
```
$ ./playlog -uvva 10
```

Only run the backend on port 6969:
```
$ ./playlog -bvl 6969
```

Use a different file for the play database:
```
$ ./playlog -p plays2.db
```

### Frontend

```
$ node build
```

To bind to a specific port, do:
```
$ PORT=3001 node build
```

## HTTP API

The API is served under `/api/`.
See [doc/backend-api.md](doc/backend-api.md) for details.

## Compatibility

The play database file SHALL remain compatible with future versions of this software.
Future versions of the play database MAY be compatible with older versions of this software.

The backend HTTP API, the Go API, and the command line interface
SHALL be backwards compatible until the next major release.
