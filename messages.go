package datapipe

import (
	"encoding/json"
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

func (d *Datapipe) WriteData(nature, entityID, projectID, factor, zoneID, itemID string, ts int64, value interface{}, reliability int) {

	info := Info{
		Nature:       nature,
		Group:        factor,
		Zone:         zoneID,
		BuildingID:   entityID,
		ProjectID:    projectID,
		DataSourceID: itemID,
	}

	var fields Fields

	switch value.(type) {
	case float64:
		v := value.(float64)
		fields = Fields{
			VALUENB:     &v,
			RELIABILITY: reliability,
		}
	case string:
		fields = Fields{
			VALUESTR:    value.(string),
			RELIABILITY: reliability,
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
