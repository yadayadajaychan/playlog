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

package update

import (
	"io"
	"encoding/json"
	
	"github.com/yadayadajaychan/playlog/database"
)

type jsonPlaylogDetail struct {
	PlaylogDetail []maimaiPlaylogDetail
}

// used by script to import playlog data from json file to sqlite3 db
func Import(playdb *database.PlayDB, data io.Reader) error {
	detail := &jsonPlaylogDetail{}
	decoder := json.NewDecoder(data)
	err := decoder.Decode(detail)
	if err != nil {
		return err
	}

	for _, v := range detail.PlaylogDetail {
		err = addMaimaiPlaylogDetailToPlayDB(playdb, v)
		if err != nil {
			return err
		}
	}

	return nil
}
