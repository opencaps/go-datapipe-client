package datapipe

import (
	"encoding/json"
	"strconv"
)

type DataType int

const (
	// DATA TYPE
	Float64 DataType = iota
	String  DataType = iota
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
	ValveOpen   = "VALVE_OPEN"
	State       = "C_STATE"
)

type Info struct {
	Nature       string `json:"Nature"`
	Zone         string `json:"Zone,omitempty"`
	Group        string `json:"Group"`
	DataSourceID string `json:"dataSourceId,omitempty"`
	BuildingID   string `json:"BuildingId"`
	ProjectID    string `json:"ProjectId"`
}

type Fields struct {
	RELIABILITY int      `json:"RELIABILITY,omitempty"`
	VALUENB     *float64 `json:"VALUENB,omitempty"`
	VALUESTR    string   `json:"VALUESTR,omitempty"`
}

type Data struct {
	Date   int64  `json:"date"`
	Fields Fields `json:"fields,omitempty"`
}

type WriteMessage struct {
	Info Info   `json:"info"`
	Data []Data `json:"data,omitempty"`
}

func (d *Datapipe) WriteData(nature, entityID, projectID, factor, zoneID, itemID string, ts int64, value string, dataType DataType) {
	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		d.log.Error("Error converting string to float64")
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

	var fields Fields

	if dataType == Float64 {
		fields = Fields{
			VALUENB:     &val,
			RELIABILITY: 1,
		}
	} else {
		fields = Fields{
			VALUESTR:    value,
			RELIABILITY: 1,
		}
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
		d.log.Error("Error marshalling data")
		return
	}

	d.log.Debug("Writing Data: ", writeMessage)
	d.sendRequest([]byte("[" + string(jsonData) + "]"))
}
