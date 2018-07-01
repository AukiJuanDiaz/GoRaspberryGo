package main

import (
	"encoding/json"
	"log"
    "fmt"
	"time"
	"strconv"
	"database/sql"

	// Added packages for implementing functionality
	resty "gopkg.in/resty.v1"
	_ "github.com/mattn/go-sqlite3"
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
	Data map[int]int
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

	bikesPerStation := map[int]int{}

	for _, element := range input.Data.Marker {
		ElementStandort_ID := element.Hal2option.Standort_id
		ElementAmountBikes := len(element.Hal2option.Bikelist)

		i, err := strconv.Atoi(ElementStandort_ID)
		if err != nil {
			log.Fatal(err)
		}

		//Add amount of bikes in a map of standort_id
		bikesPerStation[i] = ElementAmountBikes
	}

	result := BikesPerStation{
		Data: bikesPerStation,
		Timestamp: input.Timestamp,
	}

	return result
}

func writeToSQLiteDB(input BikesPerStation) {
     db, err := sql.Open("sqlite3", "./stadtRadTest2.db")
     checkErr(err)

     stmt, err := db.Prepare("INSERT INTO testTableTwoStations(s131881, s198077) values(?,?)")
     checkErr(err)


	 //Parse amount of bikes out of map
     var s131881 int = input.Data[131881]
     var s198077 int = input.Data[198077]

	 /*
	 fmt.Println(s131881)
	 fmt.Println(s198077)
	 */

     res, err := stmt.Exec(s131881, s198077)
     checkErr(err)
     
     fmt.Println(res)
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

func main(){
	rawdata := GetStadtRad()
	stationStruct := GetBikesPerStation(rawdata)
	fmt.Printf("%+v", stationStruct)
	writeToSQLiteDB(stationStruct)
}
