package main

import (
	"encoding/json"
	"net/http"
	"time"
)

// our storage mechanism
// device_id => DeviceData
var storage = map[string]*DeviceData{}

type SensorReading struct {
	ID       string    `json:"id"`
	Readings []Reading `json:"readings"`
}

type Reading struct {
	Timestamp time.Time `json:"timestamp"`
	Count     int64     `json:"count"`
}

type DeviceData struct {
	ID       string
	Readings map[string]int64 // timestamp to reading
	Sum      int64
	Latest   time.Time
}

func main() {
	http.HandleFunc("GET /v1/devices/{id}/cumulative", getCumulative)
	http.HandleFunc("GET /v1/devices/{id}/latest", latestTimestamp)
	http.HandleFunc("POST /v1/devices", storeReadings)
	http.ListenAndServe(":8888", nil)
}

func notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "device not found"})
}

func storeReadings(w http.ResponseWriter, r *http.Request) {
	// read request
	readingRequest := &SensorReading{}
	if err := json.NewDecoder(r.Body).Decode(readingRequest); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}
	if readingRequest.ID == "" {
		// invalid request body
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "device id must be specified."})
		return
	}
	device := storage[readingRequest.ID]
	if device == nil {
		device = &DeviceData{
			ID:       readingRequest.ID,
			Readings: map[string]int64{},
		}
	}
	// parse readings
	for _, r := range readingRequest.Readings {
		_, exists := device.Readings[r.Timestamp.String()]
		if exists {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(map[string]string{"error": "duplicate reading."})
			return
		} else {
			// increment sum
			device.Readings[r.Timestamp.String()] = r.Count
			device.Sum += r.Count
			// update latest
			if r.Timestamp.Compare(device.Latest) > 0 {
				device.Latest = r.Timestamp
			}
		}
		// store
		storage[readingRequest.ID] = device
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func getCumulative(w http.ResponseWriter, r *http.Request) {
	deviceID := r.PathValue("id")
	device := storage[deviceID]
	if device == nil {
		notFound(w)
		return
	}
	json.NewEncoder(w).Encode(map[string]int64{"count": device.Sum})
}

func latestTimestamp(w http.ResponseWriter, r *http.Request) {
	deviceID := r.PathValue("id")
	device := storage[deviceID]
	if device == nil {
		notFound(w)
		return
	}
	json.NewEncoder(w).Encode(map[string]time.Time{
		"latest_timestamp": device.Latest,
	})
}
