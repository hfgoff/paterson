// go build -o bus
package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

//go:embed .env
var dotenv string

type BustimeResponse struct {
	BusTimeResponse struct {
		Predictions []Prd `json:"prd"`
	} `json:"bustime-response"`
}

type Prd struct {
	DistanceFeet       int    `json:"dstp"`
	RouteDisplay       string `json:"rtdd"`
	Direction          string `json:"rtdir"`
	Destination        string `json:"des"`
	PredictionDateTime string `json:"prdtm"`
	IsDelayed          bool   `json:"dly"`
	PredictionMinutes  string `json:"prdctdn"`
}

func main() {
	env := make(map[string]string)

	lines := strings.SplitSeq(dotenv, "\n")
	for line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			val := strings.Trim(value, `"'`)
			env[key] = val
		}
	}

	if env["API_KEY"] != "" {
		fmt.Printf("[info] loaded BusTime API key\n")
	}

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		client := &http.Client{}

		endpoint, err := url.Parse("https://metromap.cityofmadison.com/bustime/api/v3/getpredictions")
		if err != nil {
			panic(err)
		}

		requestParams := url.Values{}
		requestParams.Set("key", env["API_KEY"])
		// Paterson SOUTHBOUND
		// See: https://metromap.cityofmadison.com/bustime/wireless/html/home.jsp
		requestParams.Set("stpid", "10089")
		requestParams.Set("top", "3")   // Get the next three predicitions
		requestParams.Set("tmres", "s") // Set prediciton to seconds precision. Defaults to m (minutes)
		requestParams.Set("format", "json")

		endpoint.RawQuery = requestParams.Encode()

		req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
		if err != nil {
			panic(err)
		}

		fmt.Printf("[info] making request to: %s\n", req.URL)

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			panic(err)
		}

		var bustime BustimeResponse
		err = json.Unmarshal(body, &bustime)
		if err != nil {
			panic(err)
		}

		type Out struct {
			Bustime     BustimeResponse
			GeneratedAt time.Time
		}

		var out Out
		out.Bustime = bustime
		out.GeneratedAt = time.Now()

		data, err := json.MarshalIndent(out, "", "  ")
		if err != nil {
			panic(err)
		}

		fmt.Printf("[info] writing %d bytes\n", len(data))

		os.WriteFile("next.json", data, 0644)
		<-ticker.C
	}
}
