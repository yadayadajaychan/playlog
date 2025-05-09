// Copyright (C) 2025 Ethan Cheng <ethan@nijika.org>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Package backend implements the backend http api
package backend

import (
	"log"
	"fmt"
	"strings"
	"strconv"
	"encoding/json"
	"net/http"
	"github.com/yadayadajaychan/playlog/internal/context"
	"github.com/yadayadajaychan/playlog/database"
)

var ctx context.PlaylogCtx

func Entrypoint(c context.PlaylogCtx) {
	ctx = c

	if ctx.Verbose >= 1 {
		log.Printf("starting backend server on port %d", ctx.ListenPort)
	}

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/api/playlog", playlogHandler)

	http.ListenAndServe(fmt.Sprintf(":%d", ctx.ListenPort), nil)
}

func logRequest(r *http.Request, statusCode int) {
	if ctx.Verbose >= 1 {
		log.Printf(`%s "%s %s %s" %d "%s" "%s"`, r.RemoteAddr, r.Method, r.RequestURI, r.Proto, statusCode, r.Host, r.UserAgent())
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(404)
	fmt.Fprintln(w, "404 Not Found")
	logRequest(r, 404)
}

type playlog struct {
	Playlog []playlogEntry
}

type playlogEntry struct {
	SongInfo database.SongInfo
	PlayInfo database.PlayInfo
}

func playlogHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(500)
			fmt.Fprintln(w, err)
			logRequest(r, 500)
			log.Print(err)
			return
		}
	}()

	values := r.URL.Query()

	var asc bool
	ascending := strings.ToLower(values.Get("ascending"))
	switch (ascending) {
	case "true":
		asc = true
	case "false":
		asc = false
	default:
		asc = false
	}

	page, err := strconv.Atoi(values.Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	count, err := strconv.Atoi(values.Get("count"))
	if err != nil || count < 0 {
		count = 50
	}
	offset := (page - 1) * count

	plays, err := ctx.Playdb.GetPlays(asc, count, offset)
	if err != nil {
		panic(err)
	}

	pl := playlog{
		Playlog: make([]playlogEntry, 0, 100),
	}

	for _, play := range plays {
		song, err := ctx.Songdb.GetSong(play.SongId)
		if err != nil {
			panic(err)
		}

		entry := playlogEntry{
			SongInfo: song,
			PlayInfo: play,
		}

		pl.Playlog = append(pl.Playlog, entry)
	}

	j, err := json.Marshal(pl)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintln(w, string(j))
	logRequest(r, 200)
	return
}
