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

package main

import (
	"time"
	"os"
	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"github.com/yadayadajaychan/playlog/internal/update"
	"github.com/yadayadajaychan/playlog/database"
	"github.com/yadayadajaychan/playlog/internal/context"
	"github.com/pborman/getopt/v2"
)


func main() {
	ctx := &context.PlaylogCtx{}

	help := getopt.BoolLong("help", 'h', "display help")
	songdbFilename := getopt.StringLong("songdb", 's', "songs.db", "filename of song db")
	playdbFilename := getopt.StringLong("playdb", 'p', "plays.db", "filename of play db")

	verbose := getopt.CounterLong("verbose", 'v', "verbosity level (errors only, info, debug)")
	listenPort := getopt.IntLong("listen-port", 'l', 5000, "port to listen on")

	updateInterval := getopt.IntLong("update-interval", 't', 900, "seconds to wait between updates")
	apiInterval := getopt.IntLong("api-interval", 'a', 3, "seconds to wait between api requests")

	getopt.FlagLong(&ctx.UpdateOnly, "update-only", 'u', "only update the play db & exit").SetGroup("action")
	getopt.FlagLong(&ctx.BackendOnly, "backend-only", 'b', "only run the backend").SetGroup("action")

	getopt.Parse()

	ctx.Verbose = *verbose
	ctx.ListenPort = *listenPort

	ctx.UpdateInterval = time.Duration(*updateInterval) * time.Second
	ctx.ApiInterval = time.Duration(*apiInterval) * time.Second

	if !ctx.UpdateOnly && !ctx.BackendOnly {
		ctx.UpdateAndBackend = true
	} else {
		ctx.UpdateAndBackend = false
	}

	if *help {
		getopt.Usage()
		os.Exit(0)
	}

	if ctx.UpdateInterval <= 0 {
		log.Fatal("update interval must be greater than 0")
	}
	if ctx.ApiInterval <= 0 {
		log.Fatal("api interval must be greater than 0")
	}


	db, err := sql.Open("sqlite3", *playdbFilename)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ctx.Playdb, err = database.NewPlayDB(db)
	if err != nil {
		panic(err)
	}

	if ctx.UpdateAndBackend || ctx.UpdateOnly {
		ctx.AccessCode = os.Getenv("PLAYLOG_ACCESS_CODE")
		if ctx.AccessCode == "" {
			log.Fatal("missing 'PLAYLOG_ACCESS_CODE' environment variable")
		}

		if ctx.UpdateOnly {
			update.Update(ctx.Playdb, ctx.AccessCode, ctx.ApiInterval)
		} else {
			go updateLoop(ctx)
		}
	}

	if ctx.UpdateAndBackend || ctx.BackendOnly {
		db2, err := sql.Open("sqlite3", *songdbFilename)
		if err != nil {
			panic(err)
		}
		defer db2.Close()

		ctx.Songdb, err = database.NewSongDB(db2)
		if err != nil {
			panic(err)
		}

		backendLoop(ctx)
	}
}

func updateLoop(ctx *context.PlaylogCtx) {
	for {
		update.Update(ctx.Playdb, ctx.AccessCode, ctx.ApiInterval)
		time.Sleep(ctx.UpdateInterval)
	}
}

func backendLoop(ctx *context.PlaylogCtx) {
	for {
		time.Sleep(time.Hour)
	}
}
