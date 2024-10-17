package datapipe

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/op/go-logging"
)

const (
	tokenRenewal = 24 * time.Hour
)

type Datapipe struct {
	Log           *logging.Logger
	url           string
	tokenEndpoint string
	certPath      string
	keyPath       string
	caPath        string
	token         string
}

// Init datapipe
func (d *Datapipe) Init(conf Conf) {
	d.url = conf.DatapipeURL
	d.certPath = conf.DatapipeCertPath
	d.keyPath = conf.DatapipeKeyPath
	d.caPath = conf.DatapipeCAPath
	d.tokenEndpoint = conf.DatapipeTokenEndPoint

	// Get token
	err := d.getToken()
	if err != nil {
		d.Log.Fatal("Error getting token", err)
	}

	go d.updateToken()
}

// Renew token every 24 hours
func (d *Datapipe) updateToken() {
	ticker := time.NewTicker(tokenRenewal)
	for range ticker.C {
		err := d.getToken()
		if err != nil {
			d.Log.Fatal("Error getting token")
		}
	}
}

func (d *Datapipe) sendRequest(data json.RawMessage) error {
	// Send request to datapipe

	// Create new Http Request
	req, err := http.NewRequest("POST", d.url, bytes.NewBuffer(data))
	if err != nil {
		d.Log.Error("Error creating request")
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+d.token)

	// Create new Http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		d.Log.Error("Error sending request")
		return err
	}

	// Read response
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		d.Log.Error("Error getting response")
		return err
	}

	return nil
}
