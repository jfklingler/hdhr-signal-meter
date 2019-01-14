package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	hdhr "github.com/mdlayher/hdhomerun"
)

const (
	tuneTranscoding = "internet240"
)

func tune(channel lineupChannel, durationSeconds int, done chan int) {
	httpClient := http.Client{}

	urlWithOptions := channel.URL + "?duration=" + strconv.Itoa(durationSeconds) + "&transcode=" + tuneTranscoding
	req, err := http.NewRequest(http.MethodGet, urlWithOptions, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set(userAgentHeader, userAgent)

	res, getErr := httpClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	io.Copy(ioutil.Discard, res.Body)
	res.Body.Close()
	done <- 0
}

func accumulateStats(client *hdhr.Client, stats *channelStats) *channelStats {
	tunerStatus := findTunerOnChannel(client)
	if tunerStatus != nil {
		stats.Accumulate(tunerStatus)
	}

	return stats
}

// findTunerOnChannel is pretty dumb right now. There doesn't seem to be a way
// to sanely correlate the linueup channel to any information available from
// the tuner.
func findTunerOnChannel(client *hdhr.Client) (tunerFound *hdhr.TunerStatus) {
	tunersTuned := 0

	err := client.ForEachTuner(func(tuner *hdhr.Tuner) error {
		// log.Printf("Checking tuner %d", tuner.Index)

		// query(client, fmt.Sprintf("/tuner%d/channel", tuner.Index))
		// query(client, fmt.Sprintf("/tuner%d/status", tuner.Index))
		// query(client, fmt.Sprintf("/tuner%d/streaminfo", tuner.Index))

		debug, dErr := tuner.Debug()
		if dErr != nil {
			return dErr
		}

		if debug.Tuner.Channel != "none" {
			// log.Printf("Tuner %d is tuned to some channel (%s)", tuner.Index, debug.Tuner.Channel)
			tunersTuned = tunersTuned + 1
			tunerFound = debug.Tuner
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error searching for tuner: %v", err)
	}

	if tunersTuned > 1 {
		log.Fatalf("More than one tuner is active - this is currently not supported.")
	}

	return
}

/*
Supported configuration options:
/lineup/scan
/sys/copyright
/sys/debug
/sys/features
/sys/hwmodel
/sys/model
/sys/restart <resource>
/sys/version
/tuner<n>/channel <modulation>:<freq|ch>
/tuner<n>/channelmap <channelmap>
/tuner<n>/debug
/tuner<n>/filter "0x<nnnn>-0x<nnnn> [...]"
/tuner<n>/lockkey
/tuner<n>/program <program number>
/tuner<n>/streaminfo
/tuner<n>/status
/tuner<n>/target <ip>:<port>
/tuner<n>/vchannel <vchannel>
*/
func query(client *hdhr.Client, q string, t int) string {
	qRes, qErr := client.Query(q)
	if qErr != nil {
		log.Fatalf("Error querying: %v", qErr)
	}
	log.Printf("[%d] %s = %s", t, q, qRes)

	return fmt.Sprintf("%s", qRes)
}
