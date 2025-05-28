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

// Package update handles updating the play database based on DataSource
package update

import (
	"log"
	"os"
	"errors"

	"github.com/yadayadajaychan/playlog/internal/context"
	"github.com/yadayadajaychan/playlog/internal/update/solips"
	"github.com/yadayadajaychan/playlog/internal/update/kamai"
)

// Update requires ctx.DataSource
func Update(ctx context.PlaylogCtx) error {
	if ctx.Verbose >= 1 {
		log.Print("starting update")
	}

	switch (ctx.DataSource) {
	case context.Solips:
		ctx.AccessCode = os.Getenv("PLAYLOG_ACCESS_CODE")
		if ctx.AccessCode == "" {
			log.Fatal("missing 'PLAYLOG_ACCESS_CODE' environment variable")
		}
		return solips.Update(ctx)

	case context.Kamai:
		ctx.KamaiUser = os.Getenv("PLAYLOG_KAMAI_USER")
		if ctx.KamaiUser == "" {
			log.Fatal("missing 'PLAYLOG_KAMAI_USER' environment variable")
		}
		return kamai.Update(ctx)

	default:
		return errors.New("invalid data source")
	}
}
