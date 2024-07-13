package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin     float64 `json:"temp"`
		KelvinLike float64 `json:"feels_like"`
		KelvinMax  float64 `json:"temp_max"`
		KelvinMin  float64 `json:"temp_min"`
		Humidity   float64 `json:"humidity"`
		Celsius    int
		CelsiusMax int
		CelsiusMin int
	} `json:"main"`
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
	city := "oujda"
	data, err := query(city)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	RenderTemplate(w, "index.html", data)
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

	var data weatherData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return weatherData{}, err
	}

	// Convert temperatures from Kelvin to Celsius
	data.Main.Celsius = int(data.Main.Kelvin - 273.15)
	data.Main.CelsiusMax = int(data.Main.KelvinMax - 273.15)
	data.Main.CelsiusMin = int(data.Main.KelvinMin - 273.15)

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
