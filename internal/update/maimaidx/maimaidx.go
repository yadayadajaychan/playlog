// Copyright (C) 2026 Ethan Cheng <ethan@nijika.org>
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

// Package maimaidx handles scraping maimaidx-eng.com for play data and
// updating the play database
package maimaidx

import (
	"io"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"


	//"github.com/yadayadajaychan/playlog/database"
	"github.com/yadayadajaychan/playlog/internal/context"
)

const (
	// GET JSESSIONID
	url1 = `https://lng-tgk-aime-gw.am-all.net/common_auth/login?site_id=maimaidxex&redirect_url=https://maimaidx-eng.com/maimai-mobile/&back_url=https://maimai.sega.com/`

	// POST credentials
	url2 = `https://lng-tgk-aime-gw.am-all.net/common_auth/login/sid/`

	// GET homepage
	url3 = `https://maimaidx-eng.com/maimai-mobile/home/`

	// GET records
	url4 = `https://maimaidx-eng.com/maimai-mobile/record/`

	// GET detailed record
	url5 = `https://maimaidx-eng.com/maimai-mobile/record/playlogDetail/`
)

var globalCookieJar, _ = cookiejar.New(nil)

// Update uses SegaID and SegaPassword to get the most recent 50 songs played
// and makes an http request per new song that's not in the database,
// delaying by ctx.ApiInterval between requests.
// It then adds them to the database.
// ctx requires Playdb, Songdb, AccessCode, ApiInterval, Verbose
func Update(ctx context.PlaylogCtx) error {
	_, err := getPlaylog(ctx)
	if err != nil {return err}

	return nil
}

// getPlaylog gets the non-detailed playlog of the most recent 50 plays.
// only the idx and user play date matter in this case.
// It returns a slice of the idx's not already in the play database.
func getPlaylog(ctx context.PlaylogCtx) ([]string, error) {
	client := &http.Client{
		Jar: globalCookieJar,
		Transport: &headerTransport{
			base: http.DefaultTransport,
			headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
			},
		},
	}

	// url1
	if ctx.Verbose >= 2 {
		log.Printf("getting JSESSIONID")
	}
	resp, err := client.Get(url1)
	if err != nil {return nil, err}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	// url2
	if ctx.Verbose >= 1 {
		log.Printf("logging into maimaidx")
	}
	loginData := url.Values{}
	loginData.Set("retention", "1")
	loginData.Set("sid", ctx.SegaID)
	loginData.Set("password", ctx.SegaPassword)

	resp, err = client.PostForm(url2, loginData)
	if err != nil {return nil, err}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	// url3
	if ctx.Verbose >= 2 {
		log.Printf("getting homepage")
	}
	resp, err = client.Get(url3)
	if err != nil {return nil, err}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	// url4
	if ctx.Verbose >= 1 {
		log.Printf("getting playlog")
	}
	resp, err = client.Get(url4)
	if err != nil {return nil, err}
	defer resp.Body.Close()

	test3, _ := io.ReadAll(resp.Body)
	fmt.Println(string(test3))

	return nil, nil
}

type headerTransport struct {
	base    http.RoundTripper
	headers map[string]string
}

func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range t.headers {
		req.Header.Set(k, v)
	}
	return t.base.RoundTrip(req)
}
