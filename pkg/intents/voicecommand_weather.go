package intents

import (
	"encoding/json"
	"fmt"
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"image/color"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

/**********************************************************************************************************************/
/*                                EXTENDED WEATHER FORECAST                                                           */
/**********************************************************************************************************************/

// *** OPENWEATHERMAP.ORG ***

type openWeatherMapAPIGeoCodingStruct struct {
	Name       string            `json:"name"`
	LocalNames map[string]string `json:"local_names"`
	Lat        float64           `json:"lat"`
	Lon        float64           `json:"lon"`
	Country    string            `json:"country"`
	State      string            `json:"state"`
}

//2.5 API

type WeatherStruct struct {
	Id          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type openWeatherMapAPIResponseStruct struct {
	Coord struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"coord"`
	Weather []WeatherStruct `json:"weather"`
	Base    string          `json:"base"`
	Main    struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	DT  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		Id      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

type openWeatherMapForecastAPIResponseStruct struct {
	Cod     string                            `json:"cod"`
	Message int                               `json:"message"`
	Cnt     int                               `json:"cnt"`
	List    []openWeatherMapAPIResponseStruct `json:"list"`
}

func Weather_Register(intentList *[]IntentDef) error {
	utterances := make(map[string][]string)
	utterances[LOCALE_ENGLISH] = []string{"weather", "whether", "the other", "the water", "no other", "weather forecast", "weather tomorrow"}
	utterances[LOCALE_ITALIAN] = []string{"che tempo fa", "com'è il tempo", "com'è fuori", "previsioni del tempo", "che tempo farà"}
	utterances[LOCALE_SPANISH] = []string{"qué tiempo hace", "cómo es el tiempo", "cómo está fuera", "pronóstico del tiempo", "qué tiempo hará"}
	utterances[LOCALE_FRENCH] = []string{"quel temps fait-il", "quel temps fera", "prévisions météorologiques", "meteo"}
	utterances[LOCALE_GERMAN] = []string{"Was ist vor langer Zeit", "Wie Zeit ist es", "wie draußen", "Wettervorhersage"}

	var intent = IntentDef{
		IntentName: "extended_intent_weather",
		Utterances: utterances,
		Parameters: []string{PARAMETER_WEATHER},
		Handler:    doWeatherForecast,
	}
	*intentList = append(*intentList, intent)
	return nil
}

func doWeatherForecast(intent IntentDef, params IntentParams) string {
	returnIntent := STANDARD_INTENT_IMPERATIVE_AFFIRMATIVE
	sdk_wrapper.UseVectorEyeColorInImages(true)
	sdk_wrapper.SetBackgroundColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
	intT, _ := strconv.ParseInt(params.Weather.Temperature, 10, 32)
	sdk_wrapper.DisplayTemperature(int(intT), params.Weather.TemperatureUnit, 500, false)
	sdk_wrapper.SayText(params.Weather.Temperature + getText(STR_WEATHER_AND) + params.Weather.Condition)
	sdk_wrapper.DisplayAnimatedGif(params.Weather.Icon, sdk_wrapper.ANIMATED_GIF_SPEED_FAST, 3, true)
	return returnIntent
}

func weatherParser(speechText string, botLocation string, botUnits string) (string, string, string, string, string, string, string) {
	var specificLocation bool
	var apiLocation string
	var speechLocation string
	var hoursFromNow int
	if strings.Contains(speechText, getText(STR_WEATHER_IN)) {
		splitPhrase := strings.SplitAfter(speechText, getText(STR_WEATHER_IN))
		speechLocation = strings.TrimSpace(splitPhrase[1])
		if len(splitPhrase) == 3 {
			speechLocation = speechLocation + " " + strings.TrimSpace(splitPhrase[2])
		} else if len(splitPhrase) == 4 {
			speechLocation = speechLocation + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
		} else if len(splitPhrase) > 4 {
			speechLocation = speechLocation + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
		}
		splitLocation := strings.Split(speechLocation, " ")
		if len(splitLocation) == 2 {
			speechLocation = splitLocation[0] + ", " + splitLocation[1]
		} else if len(splitLocation) == 3 {
			speechLocation = splitLocation[0] + " " + splitLocation[1] + ", " + splitLocation[2]
		}
		println("Location parsed from speech: " + "`" + speechLocation + "`")
		specificLocation = true
	} else {
		println("No location parsed from speech")
		specificLocation = false
	}
	hoursFromNow = 0
	hours, _, _ := time.Now().Clock()
	if strings.Contains(speechText, getText(STR_WEATHER_THIS_AFTERNOON)) {
		if hours < 14 {
			hoursFromNow = 14 - hours
		}
	} else if strings.Contains(speechText, getText(STR_WEATHER_TONIGHT)) {
		if hours < 20 {
			hoursFromNow = 20 - hours
		}
	} else if strings.Contains(speechText, getText(STR_WEATHER_THE_DAY_AFTER_TOMORROW)) {
		hoursFromNow = 24 - hours + 24 + 9
	} else if strings.Contains(speechText, getText(STR_WEATHER_FORECAST)) ||
		strings.Contains(speechText, getText(STR_WEATHER_TOMORROW)) {
		hoursFromNow = 24 - hours + 9
	}
	println("Looking for forecast " + strconv.Itoa(hoursFromNow) + " hours from now...")

	if specificLocation {
		apiLocation = speechLocation
	} else {
		apiLocation = botLocation
	}
	// call to weather API
	condition, is_forecast, local_datetime, speakable_location_string, temperature, temperature_unit, icon := getWeather(apiLocation, botUnits, hoursFromNow)
	return condition, is_forecast, local_datetime, speakable_location_string, temperature, temperature_unit, icon
}

func getWeather(location string, botUnits string, hoursFromNow int) (string, string, string, string, string, string, string) {
	var weatherEnabled bool
	var condition string
	var is_forecast string
	var local_datetime string
	var speakable_location_string string
	var temperature string
	var temperature_unit string
	var icon string = sdk_wrapper.GetDataPath("images/weather/conditions/rain_light.gif")
	weatherAPIEnabled := os.Getenv("WEATHERAPI_ENABLED")
	weatherAPIKey := os.Getenv("WEATHERAPI_KEY")
	weatherAPIUnit := os.Getenv("WEATHERAPI_UNIT")
	weatherAPIProvider := os.Getenv("WEATHERAPI_PROVIDER")
	if weatherAPIEnabled == "true" && weatherAPIKey != "" {
		weatherEnabled = true
		println("Weather API Enabled")
	} else {
		weatherEnabled = false
		println("Weather API not enabled, using placeholder")
		if weatherAPIEnabled == "true" && weatherAPIKey == "" {
			println("Weather API enabled, but Weather API key not set")
		}
	}
	if weatherEnabled {
		if botUnits != "" {
			if botUnits == "F" {
				println("Weather units set to F")
				weatherAPIUnit = "F"
			} else if botUnits == "C" {
				println("Weather units set to C")
				weatherAPIUnit = "C"
			}
		} else if weatherAPIUnit != "F" && weatherAPIUnit != "C" {
			println("Weather API unit not set, using F")
			weatherAPIUnit = "F"
		}
	}

	if weatherEnabled && weatherAPIProvider == "openweathermap.org" {
		// First use geocoding api to convert location into coordinates
		// E.G. http://api.openweathermap.org/geo/1.0/direct?q={city name},{state code},{country code}&limit={limit}&appid={API key}
		url := "http://api.openweathermap.org/geo/1.0/direct?q=" + location + "&limit=1&appid=" + weatherAPIKey
		resp, err := http.Get(url)
		if err != nil {
			println(err.Error())
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		geoCodingResponse := string(body)

		var geoCodingInfoStruct []openWeatherMapAPIGeoCodingStruct

		err = json.Unmarshal([]byte(geoCodingResponse), &geoCodingInfoStruct)
		if err != nil {
			println(err)
		}
		if len(geoCodingInfoStruct) == 0 {
			println("Geo provided no response.")
			condition = "undefined"
			is_forecast = "false"
			local_datetime = "test"              // preferably local time in UTC ISO 8601 format ("2022-06-15 12:21:22.123")
			speakable_location_string = location // preferably the processed location
			temperature = "120"
			temperature_unit = "C"
			return condition, is_forecast, local_datetime, speakable_location_string, temperature, temperature_unit, icon
		}
		Lat := fmt.Sprintf("%f", geoCodingInfoStruct[0].Lat)
		Lon := fmt.Sprintf("%f", geoCodingInfoStruct[0].Lon)

		println("Lat: " + Lat + ", Lon: " + Lon)
		println("Name: " + geoCodingInfoStruct[0].Name)
		println("Country: " + geoCodingInfoStruct[0].Country)

		// Now that we have Lat and Lon, let's query the weather
		units := "metric"
		if weatherAPIUnit == "F" {
			units = "imperial"
		}
		if hoursFromNow == 0 {
			url = "https://api.openweathermap.org/data/2.5/weather?lat=" + Lat + "&lon=" + Lon + "&units=" + units + "&appid=" + weatherAPIKey
		} else {
			url = "https://api.openweathermap.org/data/2.5/forecast?lat=" + Lat + "&lon=" + Lon + "&units=" + units + "&appid=" + weatherAPIKey
		}
		resp, err = http.Get(url)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, _ = io.ReadAll(resp.Body)
		weatherResponse := string(body)
		var openWeatherMapAPIResponse openWeatherMapAPIResponseStruct

		if hoursFromNow > 0 {
			// Forecast request: free API results are returned in 3 hours slots
			var openWeatherMapForecastAPIResponse openWeatherMapForecastAPIResponseStruct
			err = json.Unmarshal([]byte(weatherResponse), &openWeatherMapForecastAPIResponse)
			openWeatherMapAPIResponse = openWeatherMapForecastAPIResponse.List[hoursFromNow/3]
		} else {
			// Current weather request
			err = json.Unmarshal([]byte(weatherResponse), &openWeatherMapAPIResponse)
		}

		if err != nil {
			panic(err)
		}

		conditionCode := openWeatherMapAPIResponse.Weather[0].Id
		println(conditionCode)

		if conditionCode < 300 {
			// Thunderstorm
			condition = "Thunderstorms"
		} else if conditionCode < 400 {
			// Drizzle
			condition = "Rain"
		} else if conditionCode < 600 {
			// Rain
			condition = "Rain"
		} else if conditionCode < 700 {
			// Snow
			condition = "Snow"
		} else if conditionCode < 800 {
			// Athmosphere
			if openWeatherMapAPIResponse.Weather[0].Main == "Mist" ||
				openWeatherMapAPIResponse.Weather[0].Main == "Fog" {
				condition = "Rain"
			} else {
				condition = "Windy"
			}
		} else if conditionCode == 800 {
			// Clear
			if openWeatherMapAPIResponse.DT < openWeatherMapAPIResponse.Sys.Sunset {
				condition = "Sunny"
			} else {
				condition = "Stars"
			}
		} else if conditionCode < 900 {
			// Cloud
			condition = "Cloudy"
		} else {
			condition = openWeatherMapAPIResponse.Weather[0].Main
		}

		is_forecast = "false"
		t := time.Unix(int64(openWeatherMapAPIResponse.DT), 0)
		local_datetime = t.Format(time.RFC850)
		println(local_datetime)
		speakable_location_string = openWeatherMapAPIResponse.Name
		temperature = fmt.Sprintf("%d", int(math.Round(openWeatherMapAPIResponse.Main.Temp)))
		if weatherAPIUnit == "C" {
			temperature_unit = "C"
		} else {
			temperature_unit = "F"
		}
	} else {
		condition = "Snow"
		is_forecast = "false"
		local_datetime = "test"              // preferably local time in UTC ISO 8601 format ("2022-06-15 12:21:22.123")
		speakable_location_string = location // preferably the processed location
		temperature = "120"
		temperature_unit = "C"
	}
	return condition, is_forecast, local_datetime, speakable_location_string, temperature, temperature_unit, icon
}
