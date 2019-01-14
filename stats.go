package main

import (
	"fmt"

	hdhr "github.com/mdlayher/hdhomerun"
)

type channelStats struct {
	samples                 int
	signalStrength          int
	signalStrengthMin       int
	signalStrengthMax       int
	signalToNoiseQuality    int
	signalToNoiseQualityMin int
	signalToNoiseQualityMax int
	symbolErrorQuality      int
	symbolErrorQualityMin   int
	symbolErrorQualityMax   int
}

func newChannelStats() *channelStats {
	cs := new(channelStats)
	cs.samples = 0
	cs.signalStrength = 0
	cs.signalStrengthMin = 999
	cs.signalStrengthMax = -1
	cs.signalToNoiseQuality = 0
	cs.signalToNoiseQualityMin = 999
	cs.signalToNoiseQualityMax = -1
	cs.symbolErrorQuality = 0
	cs.symbolErrorQualityMin = 999
	cs.symbolErrorQualityMax = -1
	return cs
}

func (cs *channelStats) Accumulate(ts *hdhr.TunerStatus) {
	cs.samples = cs.samples + 1
	cs.signalStrength = cs.signalStrength + ts.SignalStrength
	cs.signalStrengthMin = min(cs.signalStrengthMin, ts.SignalStrength)
	cs.signalStrengthMax = max(cs.signalStrengthMax, ts.SignalStrength)
	cs.signalToNoiseQuality = cs.signalToNoiseQuality + ts.SignalToNoiseQuality
	cs.signalToNoiseQualityMin = min(cs.signalToNoiseQualityMin, ts.SignalToNoiseQuality)
	cs.signalToNoiseQualityMax = max(cs.signalToNoiseQualityMax, ts.SignalToNoiseQuality)
	cs.symbolErrorQuality = cs.symbolErrorQuality + ts.SymbolErrorQuality
	cs.symbolErrorQualityMin = min(cs.symbolErrorQualityMin, ts.SymbolErrorQuality)
	cs.symbolErrorQualityMax = max(cs.symbolErrorQualityMax, ts.SymbolErrorQuality)
}

// go's stdlib is completely retarded
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func (cs channelStats) SignalStrength() int {
	if cs.samples == 0 {
		return 0
	}

	return cs.signalStrength / cs.samples
}

func (cs channelStats) SignalToNoiseQuality() int {
	if cs.samples == 0 {
		return 0
	}

	return cs.signalToNoiseQuality / cs.samples
}

func (cs channelStats) SymbolErrorQuality() int {
	if cs.samples == 0 {
		return 0
	}

	return cs.symbolErrorQuality / cs.samples
}

func (cs channelStats) String() string {
	return fmt.Sprintf("%d,%d,%d,%d,%d,%d,%d,%d,%d,%d", cs.samples,
		cs.signalStrengthMin, cs.SignalStrength(), cs.signalStrengthMax,
		cs.signalToNoiseQualityMin, cs.SignalToNoiseQuality(), cs.signalToNoiseQualityMax,
		cs.symbolErrorQualityMin, cs.SymbolErrorQuality(), cs.symbolErrorQualityMax)
}
