package stadtrad

import (
	"fmt"
	"time"
	"strconv"
	"database/sql"
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

type ListOfStations []int

type BikesPerStation struct {
	Data map[int]int
	Timestamp time.Time
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

func OpenDatabaseConnection(){
	 var err error 
	 DBConn, err = sql.Open("sqlite3", "./stadtRadDataVault.db")
     checkErr(err)
}

func CreateStationToBikesTable(list ListOfStations){
	for _, elem := range list {
		AddStationToTable(elem)
	} 
}

func AddStationToTable(StationID int){
     var columnName string = "s" + strconv.Itoa(StationID)
     var alterTable string = "ALTER TABLE bikeAmountToStations ADD COLUMN " + columnName + " NUMERIC"
     stmt, err := DBConn.Prepare(alterTable)
     checkErr(err)
     
     res, err := stmt.Exec()
     checkErr(err)
     
     fmt.Println(res)
}

func AddColumnToTable(columnName string, columnType string, tableName string, pathToFile string){
     
     var alterTable string = "ALTER TABLE " + tableName + " ADD COLUMN " + columnName + " " + columnType
     stmt, err := DBConn.Prepare(alterTable)
     checkErr(err)
     
     _, err2 := stmt.Exec()
     checkErr(err2)
} 
