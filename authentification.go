package datapipe

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"net/http"
	"os"
)

func (d *Datapipe) newTLSConfig() *tls.Config {
	cert, err := tls.LoadX509KeyPair(d.certPath, d.keyPath)
	if err != nil {
		d.Log.Fatal("Error loading client certificate", err)
	}

	caCert, err := os.ReadFile(d.caPath)
	if err != nil {
		d.Log.Fatal("Error loading CA certificate", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
	}

	return tlsConfig
}

func (d *Datapipe) getToken() error {
	// Get token from authenticator with tls certificate
	tlsConfig := d.newTLSConfig()

	// create http client
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	req, err := http.NewRequest("GET", d.tokenEndpoint+d.url, nil)
	if err != nil {
		d.Log.Error("Error creating request")
		return err
	}

	d.Log.Debug("Request: ", req)

	// send request
	resp, err := httpClient.Do(req)
	if err != nil {
		d.Log.Error("Error sending request")
		return err
	}

	// read response
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		d.Log.Error("Error getting token")
		return err
	}

	// read token
	type token struct {
		Token string `json:"token"`
	}
	var t token
	err = json.NewDecoder(resp.Body).Decode(&t)
	if err != nil {
		d.Log.Error("Error reading token")
		return err
	}
	d.Lock()
	d.token = t.Token
	d.Unlock()

	return nil
}
