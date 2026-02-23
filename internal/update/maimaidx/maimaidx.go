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
	"log"
	"time"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"errors"
	"strings"
	"strconv"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"


	"github.com/yadayadajaychan/playlog/database"
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
	playlog, err := getPlaylog(ctx)
	if err != nil {return err}

	playlog, err = deleteOldEntries(ctx, playlog)
	if err != nil {return err}

	for _, v := range playlog {
		play, err := getPlaylogDetail(ctx, v)
		if err != nil {return err}

		log.Println(play)
		//log.Println(v.UserPlayDate)
		//log.Println(v.Idx)

		time.Sleep(ctx.ApiInterval)
	}

	return nil
}

type playlogEntry struct {
	UserPlayDate int64
	Idx          string
}

// getPlaylog gets the non-detailed playlog of the most recent 50 plays.
// only the idx and user play date matter in this case.
func getPlaylog(ctx context.PlaylogCtx) ([]playlogEntry, error) {
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

	//test3, _ := io.ReadAll(resp.Body)
	//fmt.Println(string(test3))

	// parse html for play date and idx
	var playlog []playlogEntry
	z := html.NewTokenizer(resp.Body)
	var isDate bool
	var userPlayDate int64

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				return playlog, nil
			} else {
				return playlog, z.Err()
			}

		case html.StartTagToken:
			name, hasAttr := z.TagName()

			if string(name) == "span" && hasAttr {
				key, val, _ := z.TagAttr()
				if string(key) == "class" && string(val) == "v_b" {
					isDate = true
				}
			}

		case html.SelfClosingTagToken:
			name, hasAttr := z.TagName()

			if string(name) == "input" && hasAttr {
				var key, val []byte
				moreAttr := true
				for moreAttr {
					key, val, moreAttr = z.TagAttr()
					if string(key) == "value" {
						p := playlogEntry{
							UserPlayDate : userPlayDate,
							Idx : string(val),
						}
						playlog = append(playlog, p)
						userPlayDate = 0
					}
				}
			}

		case html.TextToken:
			if isDate {
				date, err := time.Parse("2006/01/02 15:04", string(z.Text()))
				if err != nil {return playlog, err}

				// correct for JST timezone
				date = date.Add(-9 * time.Hour)

				userPlayDate = date.Unix()
				isDate = false
			}
		}
	}

	return playlog, nil
}

// deleteOldEntries takes a slice of playlog entries and returns
// a slice of playlog entries that aren't already in the play db
func deleteOldEntries(ctx context.PlaylogCtx, playlog []playlogEntry) ([]playlogEntry, error) {
	playdb := ctx.Playdb
	var output []playlogEntry

	for _, entry := range playlog {
		_, err := playdb.GetPlay(entry.UserPlayDate)
		if _, ok := err.(*database.PlayNotFoundError); ok {
			output = append(output, entry)
		} else if err == nil && ctx.Verbose >= 2 {
			log.Printf("%v already in play db", entry.UserPlayDate)
		} else if err != nil {
			return output, err
		}
	}

	return output, nil
}

// getPlaylogDetail gets the detailed playlog 
func getPlaylogDetail(ctx context.PlaylogCtx, play playlogEntry) (database.PlayInfo, error) {
	output := database.PlayInfo{}
	output.UserPlayDate = play.UserPlayDate

	client := &http.Client{
		Jar: globalCookieJar,
		Transport: &headerTransport{
			base: http.DefaultTransport,
			headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
			},
		},
	}

	// url5
	if ctx.Verbose >= 2 {
		log.Printf("getting idx %s", play.Idx)
	}
	resp, err := client.Get(url5 + "?idx=" + play.Idx)
	if err != nil {return output, err}
	defer resp.Body.Close()

	//test, _ := io.ReadAll(resp.Body)
	//log.Println(string(test))

	doc, err := html.Parse(resp.Body)
	if err != nil {return output, err}

	// find the playlog_top_container div
	var top *html.Node
	top_loop:
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.Div {
			for _, a := range n.Attr {
				if a.Key == "class" && a.Val == "playlog_top_container p_r" {
					top = n
					break top_loop
				}
			}
		}
	}
	if top == nil {
		return output, errors.New("playlogDetail: failed to find playlog_top_container")
	}

	err = parseTop(top, &output)
	if err != nil {return output, err}


	return output, nil
}

func parseTop(top *html.Node, play *database.PlayInfo) error {
	//parseError := errors.New("playlogDetail: failed to parse playlog_top_container")

	// img
	if top.FirstChild == nil {
		return errors.New("playlogDetail: failed to find difficulty image")
	}
	diff := top.FirstChild.NextSibling
	if diff == nil {
		return errors.New("playlogDetail: failed to find difficulty image")
	}

	diffFound := false
	for _, a := range diff.Attr {
		if a.Key == "src" {
			switch a.Val {
			case "https://maimaidx-eng.com/maimai-mobile/img/diff_basic.png":
				play.Difficulty = database.Basic
			case "https://maimaidx-eng.com/maimai-mobile/img/diff_advanced.png":
				play.Difficulty = database.Advanced
			case "https://maimaidx-eng.com/maimai-mobile/img/diff_expert.png":
				play.Difficulty = database.Expert
			case "https://maimaidx-eng.com/maimai-mobile/img/diff_master.png":
				play.Difficulty = database.Master
			case "https://maimaidx-eng.com/maimai-mobile/img/diff_remaster.png":
				play.Difficulty = database.ReMaster
			case "https://maimaidx-eng.com/maimai-mobile/img/diff_utage.png":
				play.Difficulty = database.Utage
			default:
				return errors.New("playlogDetail: invalid difficulty: " + a.Val)
			}
			diffFound = true
		}
	}
	if !diffFound {
		return errors.New("playlogDetail: failed to find difficulty image source")
	}

	// div
	if diff.NextSibling == nil {
		return errors.New("playlogDetail: failed to find sub_title div")
	}
	subtitle := diff.NextSibling.NextSibling
	if subtitle == nil {
		return errors.New("playlogDetail: failed to find sub_title div")
	}

	// span
	if subtitle.FirstChild == nil {
		return errors.New("playlogDetail: failed to find track span")
	}
	track := subtitle.FirstChild.NextSibling
	if track == nil {
		return errors.New("playlogDetail: failed to find track span")
	}

	// text
	t := track.FirstChild
	if t == nil {
		return errors.New("playlogDetail: failed to find track")
	}

	_, tt, ok := strings.Cut(t.Data, " ")
	if !ok {
		return errors.New("playlogDetail: failed to cut track")
	}

	trackNumber, err := strconv.Atoi(tt)
	if err != nil {
		return errors.New("playlogDetail: failed to convert track to int")
	}

	play.Track = trackNumber

	return nil
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
