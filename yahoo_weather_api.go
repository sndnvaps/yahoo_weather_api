package yahoo

import (
	"encoding/json"
	"fmt"
	"github.com/koron/go-dproxy"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	publicAPIURL string = "http://query.yahooapis.com/v1/public/yql"
)

type ForecastInfo struct {
	Code float64
	Date string //time.Time
	Day  string // Tue,Wed,Thu,Fri,Sat,Sun
	High float64
	Low  float64
	Text string
}

type WindInfo struct {
	Chill     float64
	Direction float64
	Speed     float64
}

type Atmosphere struct {
	Humidity   float64
	Pressure   float64
	Rising     float64
	Visibility float64
}

type Image struct {
	url string
}

type Conditions struct {
	Title string
	Lat   float64
	Long  float64
	Date  time.Time
	temp  float64
	Text  string
}

// runQuery runs the query and retuns the
// results in an io.Reader
func RunQuery(query string) []byte {
	queryURL := BuildURL(query)
	return getJSON(queryURL)
}

//buildQuery creates a YQL sql query
// select * from weather.forecast where woeid in  (select woeid from geo.places where text="jiangmen,guangdong,china")
//location == "jiangmen,guangdong,china or
//location == "jiangmen,guangdong
func BuildQuery(location string) string {
	query := fmt.Sprintf(`select * from weather.forecast where woeid in  (select woeid from geo.places where text="%s")`, location)
	return query
}

// buildURL creates a YQL URL from a query
func BuildURL(query string) string {
	params := url.Values{}
	params.Add("q", query)
	params.Add("format", "json")
	//	params.Add("diagnostics", "true")
	params.Add("callback", "")

	return publicAPIURL + "?" + params.Encode()
}

// getJSON returns the JSON response
// from a request to a given URL
func getJSON(url string) []byte {
	resp, errReq := http.Get(url)

	if errReq != nil {
		log.Fatalf("%s", errReq.Error())
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return body //bytes.NewReader(body)
}

func GetForecastlData(location string) ([]ForecastInfo, string) {
	var Forecasts []ForecastInfo

	ChannelNode := GetChannelNode(location)
	forecasts := ChannelNode.M("item").M("forecast")
	fca, _ := forecasts.Array()
	for i := 0; i < len(fca); i++ {
		//f := forecasts.A(i)
		Forecasts = append(Forecasts, newForecastDataFromRow(forecasts.A(i)))
	}
	icon_url := GetWeatherIcon(ChannelNode)
	return Forecasts, icon_url

}

// newHistoricalPieceFromRow returns a
func newForecastDataFromRow(row dproxy.Proxy) ForecastInfo {
	f := ForecastInfo{}

	code, _ := row.M("code").String()
	f.Code, _ = strconv.ParseFloat(code, 64)

	date, _ := row.M("date").String()
	//for debug
	//fmt.Println(date)
	f.Date = date
	//f.Date, _ =  time.Parse("2006-01-02", date)

	f.Day, _ = row.M("day").String() //

	high, _ := row.M("high").String()
	High, _ := strconv.ParseFloat(high, 64)
	f.High = Fahrenheit2Celsius(High)//High 

	low, _ := row.M("low").String()
	Low, _ := strconv.ParseFloat(low, 64)
	f.Low =  Fahrenheit2Celsius(Low)//Low

	f.Text, _ = row.M("text").String()

	return f
}

func GetChannelNode(location string) dproxy.Proxy {
	var v interface{}
	query := BuildQuery(location)
	jsr := RunQuery(query)
	json.Unmarshal([]byte(jsr), &v)
	channel := dproxy.New(v).M("query").M("results").M("channel")
	return channel
}

//get the icon of weather yahoo
// i := GetChannelNode(location)
func GetWeatherIcon(i dproxy.Proxy) string {
	icon, _ := i.M("image").M("url").String()
	return icon
}

//get the wind info
//w := GetChannelNode(location)
func GetWindInfo(w dproxy.Proxy) WindInfo {
	wind := WindInfo{}
	wi := w.M("wind")

	chill, _ := wi.M("chill").String()
	wind.Chill, _ = strconv.ParseFloat(chill, 64)

	direction, _ := wi.M("direction").String()
	wind.Direction, _ = strconv.ParseFloat(direction, 64)

	speed, _ := wi.M("speed").String()
	wind.Speed, _ = strconv.ParseFloat(speed, 64)

	return wind
}

//℃
func Fahrenheit2Celsius(f float64) float64 {
	return ((f - 32) / 1.8)
}
