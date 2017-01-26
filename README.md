# Yahoo Query Language (YQL)

[![GoDoc](https://godoc.org/github.com/sndnvaps/yahoo_weather_api?status.svg)](https://godoc.org/github.com/sndnvaps/yahoo_weather_api)


#how to get it 
	go get github.com/sndnvaps/yahoo_weather_api

```
type ForecastInfo struct {
	Code float64
	Date string //time.Time
	Day  string // Tue,Wed,Thu,Fri,Sat,Sun
	High float64
	Low  float64
	Text string
}
```

# example

```go
package main

import (
	"github.com/sndnvaps/yahoo_weather_api"
	"fmt"
)

func main() {

	f, icon_url := yahoo.GetForecastlData("jiangmen,guangdong,china")
	fmt.Println(f,icon_url)
}
```

output:

```json
	[{30 25 Jan 2017 Wed 71 56 Partly Cloudy} {32 26 Jan 2017 Thu 71 52 Sunny} {34 27 Jan 2017 Fri 71 51 Mostly Sunny} {30 28 Jan 2017 Sat 72 55 Partly Cloudy} {28 29 Jan 2017 Sun 72 62 Mostly Cloudy} {28 30 Jan 2017 Mon 73 64 Mostly Cloudy} {28 31 Jan 2017 Tue 71 59 Mostly Cloudy} {30 01 Feb 2017 Wed 74 60 Partly Cloudy} {30 02 Feb 2017 Thu 73 60 Partly Cloudy} {30 03 Feb 2017 Fri 70 61 Partly Cloudy}]
```
