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

package context

import (
	"time"
	"github.com/yadayadajaychan/playlog/database"
)

type DataSource int
const (
	Solips DataSource = iota
	Kamai
)

type PlaylogCtx struct {
	DataSource DataSource
	AccessCode string // Mythos Access Code
	KamaiUser string // kamaitachi username

	Playdb *database.PlayDB
	Songdb *database.SongDB

	Verbose        int
	ListenPort     int

	UpdateInterval time.Duration
	ApiInterval    time.Duration

	UpdateOnly       bool
	BackendOnly      bool
	UpdateAndBackend bool
}
