package datapipe

import (
	"encoding/json"
	"strconv"
)

const (
	// NATURE
	Sensor   = "SEN"
	Factor   = "FACT"
	Actuator = "ACT"
	Target   = "TAR"
	It       = "IT"

	// GROUP
	Temperature = "TMP"
	State       = "C_STATE"
)

type Info struct {
	Nature       string `json:"Nature"`
	Zone         string `json:"Zone"`
	Group        string `json:"Group"`
	DataSourceID string `json:"dataSourceId,omitempty"`
	BuildingID   string `json:"BuildingId"`
	ProjectID    string `json:"ProjectId"`
}

type Fields struct {
	RELIABILITY int     `json:"RELIABILITY,omitempty"`
	VALUENB     float64 `json:"VALUENB,omitempty"`
	VALUESTR    string  `json:"VALUESTR,omitempty"`
}

type Data struct {
	Date   int64  `json:"date"`
	Fields Fields `json:"fields,omitempty"`
}

type WriteMessage struct {
	Info Info   `json:"info"`
	Data []Data `json:"data,omitempty"`
}

func (d *Datapipe) WriteData(nature, entityID, projectID, factor, zoneID, itemID string, ts int64, value string) {
	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		d.Log.Error("Error converting string to float64")
		return
	}

	info := Info{
		Nature:       nature,
		Group:        factor,
		Zone:         zoneID,
		BuildingID:   entityID,
		ProjectID:    projectID,
		DataSourceID: itemID,
	}

	fields := Fields{
		VALUENB:     val,
		RELIABILITY: 1,
	}

	data := Data{
		Date:   ts,
		Fields: fields,
	}

	writeMessage := WriteMessage{
		Info: info,
		Data: []Data{data},
	}

	jsonData, err := json.Marshal(writeMessage)
	if err != nil {
		d.Log.Error("Error marshalling data")
		return
	}

	d.Log.Debug("Writing Data: ", writeMessage)
	d.sendRequest([]byte("[" + string(jsonData) + "]"))
}
