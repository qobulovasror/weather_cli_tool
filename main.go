package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

type Weather struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} "json: coord"
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	}
	Main struct {
		Temp     float64 `json:"temp"`
		Pressure float64 `json:"pressure"`
		Humidity float64 `json:"humidity"`
	} "json: main"
	Wind struct {
		Speed float64 `json:"speed"`
		Deg   float64 `json:"deg"`
	}
	Clouds struct {
		All float64 `json:"all"`
	}
	Sys struct {
		Country string `json: country`
		Sunrise int64  `json: sunrise`
		Sunset  int64  `json: sunset`
	}
	Timezone int64  `json: timezone`
	Name     string `json: name`
}

func main() {
	q := "samarkand"
	if len(os.Args) >= 2 {
		q = os.Args[1]
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	token := os.Getenv("API_TOKEN")

	res, err := http.Get("https://api.openweathermap.org/data/2.5/weather?units=metric&q=" + q + "&appid=" + token)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("Weather API not available")
	}

	body, err := io.ReadAll((res.Body))

	if err != nil {
		panic(err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)

	if err != nil {
		panic(err)
	}

	city, country, temp, desc := weather.Name, weather.Sys.Country, weather.Main.Temp, weather.Weather[0]
	fmt.Printf("%s, %s: ", city, country)
	if temp < 8 {
		color.Red("%.0f C 🌡\n", temp)
	} else if temp > 8 && temp <= 30 {
		color.Green("%.0f C 🌡\n", temp)
	} else if temp > 30 {
		color.Yellow("%.0f C 🌡\n", temp)
	}
	fmt.Printf("Description:  %s, %s\n", desc.Main, desc.Description)
	clouds, wind, rise, set := weather.Clouds.All, weather.Wind, weather.Sys.Sunrise, weather.Sys.Sunset
	fmt.Printf("Sun 🌤 : ⬆️ %d: %d, ⬇️ %d: %d\n", time.Unix(rise, 0).Hour(), time.Unix(rise, 0).Minute(), time.Unix(set, 0).Hour(), time.Unix(set, 0).Minute())
	fmt.Printf("Clouds ☁️ :  %.0f\n", clouds)
	fmt.Printf("Wind 🌬 :  %.0f m/s\n", wind.Speed)

}
