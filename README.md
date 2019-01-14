# HDHomeRun Signal Meter

This isn't, at the moment, a signal meter as anyone thinks of a signal meter.
It simply tunes in to each channel in the lineup for a little while and
records the connection status periodically. Once all channels have been
visited, it prints a report, and exits. I was motivated to write this because
I was planning on moving my HDHomeRun to a different location and wanted a
reasonable way to compare the signal strenght/quality. Also, I'm not exactly
proficient in Go at the moment, so the quality of the code is highly suspect.

A few assumptions were made about the physical setup and conditions. If these
assumptions prove untrue, then the tool won't work as expected. The ones I
think that are most likely to be problems are:

1. There is only one HDHomeRun on the network and it is discoverable
2. There is only one tuner being actively used and it's used by the this tool

To run it yourself clone the repo into your `GOPATH` or just
`go get github.com/jfklingler/hdhr-signal-meter`, then

```bash
cd $GOPATH/github.com/jfklingler/hdhr-signal-meter
make all
./hdhr-signal-meter
```

If you don't already have [Go](https://golang.org/doc/install) and
[dep](https://github.com/golang/dep) installed and configured, you'll need to
do that first.
