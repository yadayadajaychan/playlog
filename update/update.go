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
	"io"
	"time"
	"strings"
	"net/url"
	"net/http"
	"net/http/cookiejar"
	"encoding/json"
)

const apiUrl = "https://www.solips.app/maimai/profile?_data=routes%2Fmaimai.profile"

type apiPlaylog struct {
	Playlog	[]apiPlaylogItem
}

type apiPlaylogItem struct {
	PlaylogApiId	string
	Info		apiPlaylogItemInfo
}

type apiPlaylogItemInfo struct {
	UserPlayDate	string
}

// Update uses the Mythos access code to get the most recent 100 songs played
// and makes an api request per new song that's not in the database,
// delaying by apiDelay between requests. It then adds them to the database.
func Update(accessCode string, apiDelay time.Duration) error {
	playlog, err := getPlaylog(accessCode)
	if err != nil {
		return err
	}

	_ = playlog

	//time.Sleep(apiDelay)

	return nil
}

// getPlaylog gets the non-detailed playlog of the most recent 100 plays.
// only the playlogApiId and userPlayDate values matter in this case.
func getPlaylog(accessCode string) (*apiPlaylog, error) {
	// POST to API and save cookie
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar: jar,
	}

	form := url.Values{
		"accessCode": {accessCode},
		"requestType": {"getUserApiId"},
	}

	req, err := http.NewRequest("POST", apiUrl, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	// GET with authentication cookie
	req.Method = "GET"
	req.Body = nil

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	// Unmarshal JSON
	playlog := apiPlaylog{
		Playlog: make([]apiPlaylogItem, 0, 100),
	}
	err = json.Unmarshal(data, &playlog)
	if err != nil {
		return nil, err
	}

	return &playlog, nil
}
