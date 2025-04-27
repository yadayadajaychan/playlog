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

// Package update handles making api requests to solips.app and
// updating the database
package update

import (
	"os"
	"testing"
)

// TestGetPlaylog tests that the returned playlog has 100 entries
// with unique PlaylogApiIds and UserPlayDates
func TestGetPlaylog(t *testing.T) {
	accessCode := os.Getenv("PLAYLOG_ACCESS_CODE")
	if accessCode == "" {
		t.Fatal("missing 'PLAYLOG_ACCESS_CODE' environment variable")
	}

	playlog, err := getPlaylog(accessCode)
	if err != nil {
		t.Fatal(err)
	}

	n := len(playlog.Playlog)
	if n != 100 {
		t.Errorf("length of playlog: expected %v, got %v", 100, n)
	}

	// check for duplicates
	seenPlaylogApiId := make(map[string]bool)
	seenUserPlayDate := make(map[string]bool)
	for _, item := range playlog.Playlog {
		if seenPlaylogApiId[item.PlaylogApiId] {
			t.Error("duplicate PlaylogApiId: ", item.PlaylogApiId)
		}
		seenPlaylogApiId[item.PlaylogApiId] = true

		if seenUserPlayDate[item.Info.UserPlayDate] {
			t.Error("duplicate UserPlayDate: ", item.Info.UserPlayDate)
		}
		seenUserPlayDate[item.Info.UserPlayDate] = true
	}
}
