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
	"fmt"
	"io"
	"time"
	"strings"
	"net/url"
	"net/http"
	"net/http/cookiejar"
	_ "encoding/json"
)

const apiUrl = "https://www.solips.app/maimai/profile?_data=routes%2Fmaimai.profile"

// Update uses the Mythos access code to get the most recent 100 songs played
// and makes an api request per new song that's not in the database,
// delaying by apiDelay between requests. It then adds them to the database.
func Update(accessCode string, apiDelay time.Duration) {
	getPlaylog(accessCode)
	//time.Sleep(apiDelay)
}

func getPlaylog(accessCode string) (error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}

	client := &http.Client{
		Jar: jar,
	}

	form := url.Values{
		"accessCode": []string{accessCode},
		"requestType": []string{"getUserApiId"},
	}


	req, err := http.NewRequest("POST", apiUrl, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	fmt.Println(string(data))

	req.Method = "GET"
	req.Body = nil

	resp, err = client.Do(req)
	if err != nil {
		return err
	}

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	fmt.Println(string(data))

	return nil
}
