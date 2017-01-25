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

func GetForecastlData(location string) []ForecastInfo {
	var Forecasts []ForecastInfo
	var v interface{}
	query := BuildQuery(location)
	jsr := RunQuery(query)
	json.Unmarshal([]byte(jsr), &v)
	forecasts := dproxy.New(v).M("query").M("results").M("channel").M("item").M("forecast") //.A(0).M("day").String()

	fca, _ := forecasts.Array()
	for i := 0; i < len(fca); i++ {
		//f := forecasts.A(i)
		Forecasts = append(Forecasts, newForecastDataFromRow(forecasts.A(i)))
	}
	return Forecasts

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
	f.High, _ = strconv.ParseFloat(high, 64)

	low, _ := row.M("low").String()
	f.Low, _ = strconv.ParseFloat(low, 64)

	f.Text, _ = row.M("text").String()

	return f
}
