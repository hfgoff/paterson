package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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
	DistanceFeet      int    `json:"dstp"`
	RouteDisplay      string `json:"rtdd"`
	Direction         string `json:"rtdir"`
	Destination       string `json:"des"`
	PredictionMinutes string `json:"prdctdn"`
}

func main() {
	var outDir, top string
	flag.StringVar(&outDir, "out", ".", "where to output the predicitons")
	flag.StringVar(&top, "top", "5", "number of predicitons to save")
	flag.Parse()

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

	var bustimeApiKey string
	if env["BUSTIME_API_KEY"] != "" {
		fmt.Printf("[info] loaded BusTime API key\n")
		bustimeApiKey = env["BUSTIME_API_KEY"]
	} else {
		fmt.Printf("[error] did not load BusTime API key\n")
		return
	}

	// Align to the top of the minute
	now := time.Now()
	sleep := now.Truncate(time.Minute).Add(time.Minute).Sub(now)
	fmt.Printf("[info] waiting %s to align to :00\n", sleep)
	time.Sleep(sleep)

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		client := &http.Client{
			Timeout: 15 * time.Second,
		}

		endpoint, err := url.Parse("https://metromap.cityofmadison.com/bustime/api/v3/getpredictions")
		if err != nil {
			fmt.Printf("[error] request failed: %v\n", err)
			<-ticker.C
			continue
		}

		requestParams := url.Values{}
		requestParams.Set("key", bustimeApiKey)
		requestParams.Set("stpid", "10089") // See: https://metromap.cityofmadison.com/bustime/wireless/html/home.jsp
		requestParams.Set("top", top)
		requestParams.Set("tmres", "s")     // Set prediciton to seconds precision. Defaults to m (minutes)
		requestParams.Set("format", "json") // Get response in JSON becuase XML is evil

		endpoint.RawQuery = requestParams.Encode()

		req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
		if err != nil {
			fmt.Printf("[error] building request: %v\n", err)
			<-ticker.C
			continue
		}

		fmt.Printf("[info] making request to: %s\n", req.URL)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("[error] doing: %v\n", err)
			<-ticker.C
			continue
		}

		const maxbody = 1 << 20 // 1 MB

		limited := io.LimitReader(resp.Body, maxbody)

		body, err := io.ReadAll(limited)
		resp.Body.Close()
		if err != nil {
			fmt.Printf("[error] read failed: %v\n", err)
			<-ticker.C
			continue
		}

		if len(body) > maxbody {
			fmt.Printf("[error] response too large (%d bytes)\n", len(body))
			continue
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("[error] bad status: %s\n", resp.Status)
			resp.Body.Close()
			<-ticker.C
			continue
		}

		var bustime BustimeResponse
		err = json.Unmarshal(body, &bustime)
		if err != nil {
			fmt.Printf("[error] invalid JSON: %v\n", err)
			<-ticker.C
			continue
		}

		type Out struct {
			Bustime     BustimeResponse
			GeneratedAt time.Time
		}

		var out Out
		out.Bustime = bustime
		out.GeneratedAt = time.Now().UTC()

		data, err := json.MarshalIndent(out, "", "  ")
		if err != nil {
			fmt.Printf("[error] indenting failed: %v\n", err)
			<-ticker.C
			continue
		}

		file := filepath.Join(outDir, "next.json")
		tmp := file + ".tmp"

		if err := os.WriteFile(tmp, data, 0644); err != nil {
			fmt.Printf("[error] write failed: %v\n", err)
			<-ticker.C
			continue
		}

		if err := os.Rename(tmp, file); err != nil {
			fmt.Printf("[error] rename failed: %v\n", err)
			<-ticker.C
			continue
		}

		fmt.Printf("[info] wrote %d bytes at %s\n", len(data), out.GeneratedAt.Format("2006-01-02 15:04:05"))

		<-ticker.C
	}
}
