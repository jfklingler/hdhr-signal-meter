package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type lineupChannel struct {
	GuideNumber string
	GuideName   string
	VideoCodec  string
	AudioCodec  string
	HD          int
	URL         string
}

const (
	lineupPath = "/lineup.json"

	userAgentHeader = "User-Agent"
	userAgent       = "hdhr-signal-meter"
)

func getLineup(devURL *url.URL) []lineupChannel {
	lineupURL := devURL.String() + lineupPath

	httpClient := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, lineupURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set(userAgentHeader, userAgent)

	res, getErr := httpClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	defer res.Body.Close()
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	channels := make([]lineupChannel, 0)
	json.Unmarshal(body, &channels)

	return channels
}
