package yahoo

import (
	"bou.ke/monkey"
	"fmt"
	"testing"
	"unsafe"
)

func TestYahooApiGetForecastData(t *testing.T) {

	f, _ := GetForecastlData("jiangmen,guangdong,china")

	if f == nil {
		t.Fatal("forecast is nil\n")
	}

	if len(f) == 0 {
		t.Fatal("forecast is empty\n")
	}

	for _, forecast := range f {
		if unsafe.Sizeof(forecast) == 0 {
			t.Fatalf("The forecast representing  is nil.\n")
		}
	}
}

func TestYahooApiGetWrongForecastData(t *testing.T) {
	monkey.Patch(getPublicAPIURL, func() string {
		return "http://wrong_url"
	})

	defer func(){
		if r := recover(); r != nil {
			monkey.UnpatchAll()
			f, _ := GetForecastlData("jiangmen,guangdong,china")

			if f == nil {
				t.Fatal("forecast is nil\n")
			}

		}
	}()

	GetForecastlData("jiangmen,guangdong,china")
	t.Fatal("No panic")
}

func TestYahooApi_(t *testing.T) {

	c := GetChannelNode("foshan,guangdong,china")
	wind := GetWindInfo(c)
	units := GetUnits(c)
	astronomy := GetAstronomy(c)
	conditions := GetConditions(c)
	atmosphere := GetAtmosphere(c)
	icon := GetWeatherIcon(c)
	if unsafe.Sizeof(wind) == 0 {
		t.Fatalf("The wind representing is nil.\n")
	}

	if unsafe.Sizeof(units) == 0 {
		t.Fatalf("The units representing is nil.\n")
	}

	if unsafe.Sizeof(astronomy) == 0 {
		t.Fatalf("The astronomy representing is nil.\n")
	}

	if unsafe.Sizeof(conditions) == 0 {
		t.Fatalf("The conditions representing is nil.\n")
	}

	if unsafe.Sizeof(atmosphere) == 0 {
		t.Fatal("The Atmosphere representing is nil. \n")
	}

	if unsafe.Sizeof(icon) == 0 {
		t.Fatal("Cannot get the url of icon, it was nil.\n")
	}

	fmt.Println(atmosphere)
	fmt.Println(conditions)
	fmt.Println(astronomy)
	fmt.Println(units)
	fmt.Println(wind)
	fmt.Println(icon)
}
