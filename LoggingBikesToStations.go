package stadtrad

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

// Database Pointer and RowCounter
var DBConn *sql.DB
var IDcurrRowBikesToStations int

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

func ListAllStations(input GetStadtRadData) ListOfStations{
	var result []int
	for _, element := range input.Data.Marker {
		ElementStandort_ID := element.Hal2option.Standort_id
		i, err := strconv.Atoi(ElementStandort_ID)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, i)
	}
	return result
}

func logBikesPerStation(input BikesPerStation, list ListOfStations, rowid int) {
	
	// create the row by entering the first element
	rowid = rowid + 1
	var firstColumnName string = "s" + strconv.Itoa(list[0])
	stmt, err := DBConn.Prepare("INSERT INTO bikeAmountToStations(id, "+ firstColumnName +") values(?, ?)")
	checkErr(err)
	var bikesAtStation int = input.Data[list[0]]
	_, err2 := stmt.Exec(rowid, bikesAtStation)
	checkErr(err2)
	list = list[1:]
	
	// update all the other columns
	for _,elem := range list{
		var columnName string = "s" + strconv.Itoa(elem)
		stmt, err := DBConn.Prepare("UPDATE bikeAmountToStations SET " + columnName + " = ? WHERE id = ?")
		checkErr(err)
		
		var bikesAtStation int = input.Data[elem]
		
		_, err3 := stmt.Exec(bikesAtStation, rowid)
		checkErr(err3)
	}
	
	IDcurrRowBikesToStations = IDcurrRowBikesToStations + 1
}

func GetHighestIDInBikesToStations() int{
	row := DBConn.QueryRow("SELECT MAX(id) FROM bikeAmountToStations")
	
	var result int
	err := row.Scan(&result)
	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		// Everything alright
	default:
		checkErr(err)
	}
	
	return result
}

func main(){
	OpenDatabaseConnection()
	IDcurrRowBikesToStations = GetHighestIDInBikesToStations()
	
	// Creating new columns in table  
	//CreateStationToBikesTable(stationList)
	
	// Repeat function calls every 10 seconds
	for t := range time.NewTicker(10 * time.Second).C {
			rawdata := GetStadtRad()
			bikesPerStationStruct := GetBikesPerStation(rawdata)
			stationList := ListAllStations(rawdata)
			logBikesPerStation(bikesPerStationStruct, stationList, IDcurrRowBikesToStations) 
			fmt.Println("Logging at time " + t.String() + "completed.")
			// Print the rawdata to console
			//fmt.Printf("%+v", bikesPerStationStruct)
	}
}
