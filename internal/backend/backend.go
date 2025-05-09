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
	"net/http"
	"github.com/yadayadajaychan/playlog/internal/context"
)

var ctx context.PlaylogCtx

func Entrypoint(c context.PlaylogCtx) {
	ctx = c

	if ctx.Verbose >= 1 {
		log.Printf("starting backend server on port %d", ctx.ListenPort)
	}

	http.HandleFunc("/", rootHandler)

	http.ListenAndServe(fmt.Sprintf(":%d", ctx.ListenPort), nil)
}

func logRequest(r *http.Request) {
	if ctx.Verbose >= 1 {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.String())
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	fmt.Fprintln(w, "hello world")
}
