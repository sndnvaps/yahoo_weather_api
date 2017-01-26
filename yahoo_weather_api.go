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

/*
   "astronomy": {
     "sunrise": "7:4 am",
     "sunset": "4:55 pm"
    }
*/
type Astronomy struct {
	Sunrise string
	Sunset  string
}

/*
   "distance": "mi",
   "pressure": "in",
   "speed": "mph",
   "temperature": "F"
*/
type Units struct {
	Distance    string
	Pressure    string
	Speed       string
	Temperature string
}

type Image struct {
	url string
}

type Conditions struct {
	Title   string
	Lat     float64
	Long    float64
	PubDate time.Time
	Temp    float64
	Text    string
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
	f.High = Fahrenheit2Celsius(High) //High

	low, _ := row.M("low").String()
	Low, _ := strconv.ParseFloat(low, 64)
	f.Low = Fahrenheit2Celsius(Low) //Low

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

//F -> ℃
func Fahrenheit2Celsius(f float64) float64 {
	return ((f - 32) / 1.8)
}

//get the units info
//u := GetChannelNode(location)
func GetUnits(u dproxy.Proxy) Units {
	units := Units{}
	ui := u.M("units")
	units.Distance, _ = ui.M("distance").String()
	units.Pressure, _ = ui.M("pressure").String()
	units.Speed, _ = ui.M("speed").String()
	units.Temperature, _ = ui.M("temperature").String()

	return units
}

//get the astronomy info
//a := GetChannelNode(location)
func GetAstronomy(a dproxy.Proxy) Astronomy {
	astronomy := Astronomy{}
	ai := a.M("astronomy")
	astronomy.Sunrise, _ = ai.M("sunrise").String()
	astronomy.Sunset, _ = ai.M("sunset").String()

	return astronomy
}

//get the Conditions info
//c := GetChannelNode(location)
func GetConditions(c dproxy.Proxy) Conditions {
	con := Conditions{}
	c_item := c.M("item")
	c_con := c_item.M("condition")
	con.Title, _ = c_item.M("title").String()

	lat, _ := c_item.M("lat").String()
	con.Lat, _ = strconv.ParseFloat(lat, 64)

	long, _ := c_item.M("long").String()
	con.Long, _ = strconv.ParseFloat(long, 64)

	pubDate, _ := c_item.M("pubDate").String()
	con.PubDate, _ = time.Parse("2006-01-02 15:04:05", pubDate)

	temp, _ := c_con.M("temp").String()
	con.Temp, _ = strconv.ParseFloat(temp, 64)

	con.Text, _ = c_con.M("text").String()

	return con
}
