package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

type weatherData struct {
	Name string `json:"name"`
	TimeZone int `json:"timezone"`
	Sys  struct {
		Country     string `json:"country"`
		Flag        string
		Sunrise     int64 `json:"sunrise"`
		Sunset      int64 `json:"sunset"`
		SunriseTime string
		SunsetTime  string
		LocalTime   string 
	} `json:"sys"`
	Main struct {
		Kelvin      float64 `json:"temp"`
		KelvinLike  float64 `json:"feels_like"`
		Humidity    float64 `json:"humidity"`
		Pressure    int     `json:"pressure"`
		Celsius     int
		CelsiusLike int
	} `json:"main"`
	WeatherAsSlice []struct {
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Weather struct {
		Description string
		Icon        string
	}
	Wind struct {
		Speed    float64 `json:"speed"`
		SpeedKmh string
	} `json:"wind"`
}

func loadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return apiConfigData{}, err
	}
	var c apiConfigData
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return apiConfigData{}, err
	}
	return c, nil
}

func getWindDescription(description string, feelsLike int, speed float64) string {
	var wind string
	speedKmh := speed * 3.6
	switch {
	case speedKmh < 1.5:
		wind = "Calm"
	case speedKmh < 5.5:
		wind = "Light Breeze"
	case speedKmh < 11.0:
		wind = "Gentle Breeze"
	case speedKmh < 19.0:
		wind = "Moderate Breeze"
	case speedKmh < 28.0:
		wind = "Fresh Breeze"
	case speedKmh < 38.0:
		wind = "Strong Breeze"
	case speedKmh < 49.0:
		wind = "Near Gale"
	case speedKmh < 61.0:
		wind = "Gale"
	case speedKmh < 74.0:
		wind = "Severe Gale"
	case speedKmh < 88.0:
		wind = "Storm"
	case speedKmh < 102.0:
		wind = "Violent Storm"
	default:
		wind = "Hurricane"
	}
	return fmt.Sprintf("Feels like %d°C. %s. %s", feelsLike, description, wind)
}

func (wd *weatherData) setCurrentTime() {
    location := time.FixedZone("Local Time", wd.TimeZone)
    wd.Sys.LocalTime = time.Now().In(location).Format("Jan 02, 03:04pm")
}

// RenderTemplate renders the specified template with data
func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	err = t.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		fmt.Println("Error executing template: ", err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		city = "oujda"
	}
	data, err := query(city)
	if err != nil {
		if err.Error() == "404" {
			// http.Error(w, "City not found", http.StatusNotFound)
			w.WriteHeader(http.StatusNotFound)
			RenderTemplate(w, "index.html", map[string]interface{}{
				"NotFound": true,
				"City":     city,
			})
		} else {
			http.Error(w, "Error querying weather data", http.StatusInternalServerError)
		}
		return
	}
	RenderTemplate(w, "index.html", map[string]interface{}{
		"WeatherData": data,
		"City":        city,
	})
}

func query(city string) (weatherData, error) {
	apiConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		return weatherData{}, err
	}
	resp, err := http.Get("https://api.openweathermap.org/data/2.5/weather?q=" + city + "&appid=" + apiConfig.OpenWeatherMapApiKey)
	if err != nil {
		return weatherData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return weatherData{}, errors.New("404")
	}

	var data weatherData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return weatherData{}, err
	}

	// Convert temperatures from Kelvin to Celsius
	data.Main.Celsius = int(data.Main.Kelvin - 273.15)
	data.Main.CelsiusLike = int(data.Main.KelvinLike - 273.15)
	data.Weather.Description = getWindDescription(data.WeatherAsSlice[0].Description, data.Main.CelsiusLike, data.Wind.Speed)

	data.Weather.Icon = "https://openweathermap.org/img/wn/" + data.WeatherAsSlice[0].Icon + "@2x.png"
	data.Sys.Flag = "https://flagcdn.com/16x12/" + strings.ToLower(data.Sys.Country) + ".png"
	data.Wind.SpeedKmh = strconv.FormatFloat(data.Wind.Speed*3.6, 'f', 1, 64)
	data.Sys.SunriseTime = time.Unix(data.Sys.Sunrise, 0).Format("03:04pm")
	data.Sys.SunsetTime = time.Unix(data.Sys.Sunset, 0).Format("03:04pm")

	data.setCurrentTime()
	return data, nil
}

func weather(w http.ResponseWriter, r *http.Request) {
	city := strings.SplitN(r.URL.Path, "/", 3)[2]
	data, err := query(city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json;charset=utf-8")
	json.NewEncoder(w).Encode(data)
}

func main() {
	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/", index)
	http.HandleFunc("/weather/", weather)
	http.ListenAndServe(":8080", nil)
}
