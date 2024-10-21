package datapipe

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/op/go-logging"
)

const (
	tokenRenewal = 24 * time.Hour
	logName      = "datapipe"
)

type Datapipe struct {
	log           *logging.Logger
	url           string
	tokenEndpoint string
	certPath      string
	keyPath       string
	caPath        string
	token         string
	sync.Mutex
}

// Init datapipe
func (d *Datapipe) Init(conf Conf) {
	d.log = logging.MustGetLogger(logName)
	logging.SetLevel(logging.INFO, logName)
	d.url = conf.DatapipeURL
	d.certPath = conf.DatapipeCertPath
	d.keyPath = conf.DatapipeKeyPath
	d.tokenEndpoint = conf.DatapipeTokenEndPoint
	loglevel, err := logging.LogLevel(conf.LogLevel)
	if err == nil {
		logging.SetLevel(loglevel, logName)
	}

	// Get token
	err = d.getToken()
	if err != nil {
		d.log.Fatal("Error getting token", err)
	}

	go d.updateToken()
}

// Renew token every 24 hours
func (d *Datapipe) updateToken() {
	ticker := time.NewTicker(tokenRenewal)
	for range ticker.C {
		err := d.getToken()
		if err != nil {
			d.log.Fatal("Error getting token")
		}
	}
}

func (d *Datapipe) sendRequest(data json.RawMessage) error {
	// Send request to datapipe

	// Create new Http Request
	req, err := http.NewRequest("POST", d.url+"/data/service", bytes.NewBuffer(data))
	if err != nil {
		d.log.Error("Error creating request")
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	d.Lock()
	req.Header.Set("Authorization", "Bearer "+d.token)
	d.Unlock()

	d.log.Debug("Request: ", req)

	// Create new Http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		d.log.Error("Error sending request")
		return err
	}

	// Read response
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		d.log.Error("Error getting response", resp.Status)
		return err
	}

	return nil
}
