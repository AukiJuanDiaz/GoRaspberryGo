package main

import (
	"encoding/json"
	"log"
    "fmt"
	"time"


	// Added packages for implementing functionality
	resty "gopkg.in/resty.v1"
)

// Minimal small types just for the status updates
type Marker []MarkerItem

type Bikelist []BikelistItem

type BikelistItem struct {
	Number string `json:"number"`
	CanBeRented bool `json:"canBeRented"`
}

type Hal2option struct {
	Standort_id string `json:"standort_id"`
	Bikelist Bikelist `json:"bikelist"`
}

type MarkerItem struct {
	Hal2option Hal2option `json:"hal2option"`
}

type GetStadtRadJSON struct {
	Marker Marker `json:"marker"`
}

type GetStadtRadData struct {
	Data GetStadtRadJSON
	Timestamp time.Time
}

type BikesPerStationData struct {
	Standort_id string
	AmountBikes int
}

type BikesPerStationList []BikesPerStationData

type BikesPerStation struct {
	Data BikesPerStationList
	Timestamp time.Time
}


func GetStadtRad() GetStadtRadData {

	// url from StadtradMap
	url := "https://stadtrad.hamburg.de/kundenbuchung/hal2ajax_process.php?zoom=10&lng1=&lat1=&lng2=&lat2=&stadtCache=&mapstation_id=&mapstadt_id=75&verwaltungfirma=&centerLng=9.986872299999959&centerLat=53.56661530000001&searchmode=default&with_staedte=N&buchungsanfrage=N&bereich=2&stoinput=&before=&after=&ajxmod=hal2map&callee=getMarker&requester=bikesuche&key=&webfirma_id=510"

	resp, err := resty.R().Get(url)
	timestamp := time.Now()

	// check, if the http get call suceeded
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nError: %v", err)
	fmt.Printf("\nResponse Status Code: %v", resp.StatusCode())
	fmt.Printf("\nResponse Status: %v", resp.Status())
	fmt.Printf("\nResponse Time: %v", resp.Time())
	fmt.Printf("\nResponse Received At: %v\n\n", resp.ReceivedAt())
	// fmt.Printf("\nResponse Body: %v", resp)

	var body GetStadtRadJSON

	err = json.Unmarshal(resp.Body(), &body)


	result := GetStadtRadData{
		Data: body,
		Timestamp: timestamp,
	}
	return result
}

func GetBikesPerStation(input GetStadtRadData) BikesPerStation{
	var bikesPerStation []BikesPerStationData

	for _, element := range input.Data.Marker {
		ElementStandort_ID := element.Hal2option.Standort_id
		ElementAmountBikes := len(element.Hal2option.Bikelist)
		ElementBikesPerStation := BikesPerStationData{
			Standort_id: ElementStandort_ID,
			AmountBikes: ElementAmountBikes,
		}
		bikesPerStation = append(bikesPerStation, ElementBikesPerStation)
	}

	result := BikesPerStation{
		Data: bikesPerStation,
		Timestamp: input.Timestamp,
	}

	return result
}

func main(){
	rawdata := GetStadtRad()
	stationStruct := GetBikesPerStation(rawdata)
	fmt.Printf("%+v", stationStruct)
}
