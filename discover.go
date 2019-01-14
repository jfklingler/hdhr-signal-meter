package main

import (
	"context"
	"io"
	"log"
	"time"

	hdhr "github.com/mdlayher/hdhomerun"
)

func discoverTuner() *hdhr.DiscoveredDevice {
	// Discover tuner devices with any ID
	tunerType := hdhr.DiscoverDeviceType(hdhr.DeviceTypeTuner)

	d, err := hdhr.NewDiscoverer(tunerType)
	if err != nil {
		log.Fatalf("Error starting discovery: %v", err)
	}

	// Discover devices for up to 2 seconds or until canceled.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	for {
		// Note: Discover blocks until a device is found or its context is
		// canceled.  Always pass a context with a cancel function, deadline,
		// or timeout.
		//
		// io.EOF is returned when the context is canceled.
		device, err := d.Discover(ctx)

		switch err {
		case nil:
			// Found a device.
			// If only one device is expected, invoke cancel here.
			log.Printf("Found device with ID %s at %s with %d tuner(s)", device.ID, device.Addr, device.Tuners)
			cancel()
			return device

		case io.EOF:
			// Context canceled; no more devices to be found.
			log.Fatal("No tuner devices found")

		default:
			log.Fatalf("Error during discovery: %v", err)
		}
	}
}
