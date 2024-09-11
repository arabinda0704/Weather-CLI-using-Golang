package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch    int64   `json:"time_epoch"`
				TempC        float64 `json:"temp_c"`
				ChanceOfRain float64 `json:"chance_of_rain"`
				Condition    struct {
					Text string `json:"text"`
				} `json:"condition"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func readAPIKey(filePath string) (string, error) {
	keyBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(keyBytes), nil
}

func main() {
	q := "Kolkata"
	if len(os.Args) >= 2 {
		q = os.Args[1]
	}

	apiKey, err := readAPIKey("apikey.txt")
	if err != nil {
		panic("Error reading API key file: " + err.Error())
	}

	url := fmt.Sprintf("http://api.weatherapi.com/v1/forecast.json?key=%s&q=%s&days=1&aqi=no&alerts=no", apiKey, q)
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("Weather API Not Available!")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	location, current, forecastHours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour
	fmt.Printf("\n%s, %s: current %.0f°C %s\n", location.Name, location.Country, current.TempC, current.Condition.Text)

	for _, hour := range forecastHours {
		date := time.Unix(hour.TimeEpoch, 0)
		if date.Before(time.Now()) {
			continue
		}

		message := fmt.Sprintf("%s - %.0f°C, %.0f%% chance of rain, %s\n", date.Format("15:04"), hour.TempC, hour.ChanceOfRain, hour.Condition.Text)
		if hour.ChanceOfRain < 40 {
			fmt.Print(message)
		} else {
			color.Red(message) // This will print the message in red
		}
	}
}
