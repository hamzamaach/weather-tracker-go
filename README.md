# Weather Tracker

Weather Tracker is a web application that allows users to search for the weather conditions of a specific city. The application fetches data from the OpenWeatherMap API and displays it in a user-friendly format.

## Features

- Search for the current weather of any city.
- Display temperature in Celsius.
- Show additional weather details such as humidity, wind speed, and rain possibility.
- Display weather icons and country flags.

## Installation

### Prerequisites

- Go (version 1.22 or later)
- OpenWeatherMap API key

### Steps

1. Clone the repository:

```sh
git clone https://github.com/hamzamaach/weather-tracker.git
cd weather-tracker
```

2. Rename the `.apiConfig.example` file to `.apiConfig` and add your OpenWeatherMap API key:

```sh
mv .apiConfig.example .apiConfig
```

Edit the `.apiConfig` file to include your OpenWeatherMap API key:

```json
{
    "OpenWeatherMapApiKey": "[put your open weather map api key here]"
}
```
<!-- 
3. Install the dependencies:

```sh
go mod tidy
```
-->
3. Run the application:

```sh
go run main.go
```

The application will be accessible at `http://localhost:8080`.

## Usage

1. Open your web browser and go to `http://localhost:8080`.
2. Enter the name of a city in the search box and click the search button.
3. The weather information for the city will be displayed.

## Project Structure

- `main.go`: The main Go file containing the server and handler functions.
- `assets/`: Directory containing static assets like CSS files.
- `index.html`: The main HTML template for the web application.
- `.apiConfig`: Configuration file for storing the OpenWeatherMap API key.

## API Endpoints

- `/`: The main endpoint for searching and displaying weather information.


## Contributing

1. Fork the repository.
2. Create your feature branch (`git checkout -b feature/AmazingFeature`).
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4. Push to the branch (`git push origin feature/AmazingFeature`).
5. Open a pull request.

## Acknowledgements

- OpenWeatherMap for providing the weather data API.
- Flagcdn for providing country flag images.

---

Feel free to reach out if you have any questions or suggestions!