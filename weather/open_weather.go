package weather

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"telegram-bot/services"
	"time"

	"github.com/labstack/gommon/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var WeatherKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("  Today  ", "weather/today"),
		tgbotapi.NewInlineKeyboardButtonData("  Tomorrow  ", "Tomorrow"),
		tgbotapi.NewInlineKeyboardButtonData("  Week  ", "Week"),
	),
)

type Coord struct {
	Lon int `json:"lon"`
	Lat int `json:"lat"`
}

type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
}

type Wind struct {
	Speed float64 `json:"speed"`
	Deg   float64 `json:"deg"`
}

type Clouds struct {
	All int `json:"all"`
}

type Sys struct {
	Type    int     `json:"type"`
	ID      int     `json:"id"`
	Message float64 `json:"message"`
	Country string  `json:"country"`
	Sunrise int     `json:"sunrise"`
	Sunset  int     `json:"sunset"`
}

type WeatherDataBundle struct {
	Coord    Coord     `json:"coord"`
	Weather  []Weather `json:"weather"`
	Base     string    `json:"base"`
	Main     Main      `json:"main"`
	Wind     Wind      `json:"wind"`
	Clouds   Clouds    `json:"clouds"`
	Dt       int       `json:"dt"`
	Sys      Sys       `json:"sys"`
	Timezone int       `json:"timezone"`
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Cod      int       `json:"cod"`
}

func GetOpenWeather(chatId, lat, lon int) (string, error) {
	log.Info("GetOpenWeather")
	var conf services.Config
	var err error

	if conf, err = services.GetConfig(); err != nil {
		return "", err
	}
	client := http.Client{}
	var req *http.Request

	if req, err = http.NewRequest("GET", "https://api.openweathermap.org/data/2.5/weather", nil); err != nil {
		return "", err
	}
	q := req.URL.Query()
	q.Add("lat", strconv.Itoa(lat))
	q.Add("lon", strconv.Itoa(lon))
	q.Add("units", "metric")
	q.Add("appid", conf.WeatherAppId)
	req.URL.RawQuery = q.Encode()

	var httpRsp *http.Response

	if httpRsp, err = client.Do(req); err != nil {
		return "", err
	}
	var bytesDataFromHttp []byte

	if bytesDataFromHttp, err = ioutil.ReadAll(httpRsp.Body); err != nil {
		return "", err
	}
	var weather WeatherDataBundle

	if err = json.Unmarshal(bytesDataFromHttp, &weather); err != nil {
		return "", err
	}
	return weather.String(), nil
}

func (w WeatherDataBundle) String() string {
	log.Info("WeatherDataBundle  String()")
	sunriseInfo, _ := strconv.ParseInt(strconv.Itoa(w.Sys.Sunrise), 10, 64)
	sunriseTime := time.Unix(sunriseInfo, 0)

	sunsetInfo, _ := strconv.ParseInt(strconv.Itoa(w.Sys.Sunset), 10, 64)
	sunsetTime := time.Unix(sunsetInfo, 0)

	res := "Description: " + w.Weather[0].Description +
		"\nTemp: " + strconv.Itoa(int(w.Main.Temp)) + "°C" +
		"\nFeels like: " + strconv.Itoa(int(w.Main.FeelsLike)) + "°C" +
		"\nHumidity: " + strconv.Itoa(w.Main.Humidity) + "%" +
		"\nWind: " + strconv.Itoa(int(w.Wind.Speed)) + " meter/sec" +
		"\nSunrise: " + strconv.Itoa(sunriseTime.Hour()) + ":" + strconv.Itoa(sunriseTime.Minute()) + " a.m." +
		"\nSunset: " + strconv.Itoa(sunsetTime.Hour()) + ":" + strconv.Itoa(sunsetTime.Minute()) + " p.m." +
		"\nLocation: " + w.Name
	return res
}
