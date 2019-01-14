package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	hdhr "github.com/mdlayher/hdhomerun"
)

const (
	tuneTimeSeconds = 60
	lockTime        = 1 * time.Second
	statInterval    = 500 * time.Millisecond
)

func main() {
	dev := discoverTuner()

	client, err := hdhr.Dial(dev.Addr)
	if err != nil {
		log.Fatalf("Error connecting to device %s discovery: %v", dev.ID, err)
	}

	channelStatsData := createChannelStatsData(getLineup(dev.URL))

	for k, v := range channelStatsData {
		log.Printf("Tuning to %s %s", k.GuideNumber, k.GuideName)
		ch := make(chan int)
		go tune(k, tuneTimeSeconds, ch)
		time.Sleep(lockTime) // Give tuner time to lock

		for tuned := true; tuned; {
			select {
			case <-ch:
				tuned = false
			default:
				channelStatsData[k] = accumulateStats(client, v)
				time.Sleep(statInterval) // Don't DOS the thing
			}
		}
	}

	report(channelStatsData)
}

func createChannelStatsData(lineup []lineupChannel) map[lineupChannel]*channelStats {
	channelStatsData := make(map[lineupChannel]*channelStats)

	for _, lc := range lineup {
		channelStatsData[lc] = newChannelStats()
	}

	return channelStatsData
}

func report(channelStatsData map[lineupChannel]*channelStats) {
	channels := make([]lineupChannel, 0)
	for k := range channelStatsData {
		channels = append(channels, k)
	}

	sort.Slice(channels, func(i, j int) bool {
		a, _ := strconv.ParseFloat(channels[i].GuideNumber, 32)
		b, _ := strconv.ParseFloat(channels[j].GuideNumber, 32)
		return a < b
	})

	fmt.Println("Channel,Samples," +
		"SS (min),SS (mean),SS (max)," +
		"SNQ (min),SNQ (mean),SNQ (max)," +
		"SEQ (min),SEQ (mean),SEQ (max)")

	for _, c := range channels {
		fmt.Printf("%s %s,%s\n", c.GuideNumber, c.GuideName, channelStatsData[c])
	}
}
