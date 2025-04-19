package plugins

import (
	"encoding/json"
	"fmt"
	pkgHTTP "github.com/apache/apisix-go-plugin-runner/pkg/http"
	"github.com/apache/apisix-go-plugin-runner/pkg/log"
	"io"
	"net/http"
)

// TimestampInserterGo is the name of the plugin.
type TimestampInserterGo struct {
}

type TimestampResponse struct {
	Time string `json:"time"`
}

// TimestampInserterGoConfig represents the plugin's configuration.
type TimestampInserterGoConfig struct {
	TimestampServiceURI string `json:"timestamp_service_uri"`
}

// Name returns the plugin's name.
func (t *TimestampInserterGo) Name() string {
	return "timestamp-inserter-go"
}

// ParseConf parses the plugin's configuration.
func (t *TimestampInserterGo) ParseConf(in []byte) (interface{}, error) {
	conf := TimestampInserterGoConfig{}
	err := json.Unmarshal(in, &conf)
	if err != nil {
		log.Errorf("failed to unmarshal config: %s", err)
		return nil, err
	}
	return conf, nil
}

// RequestFilter executes during the request phase.
func (t *TimestampInserterGo) RequestFilter(conf interface{}, w http.ResponseWriter, r pkgHTTP.Request) {
	config, ok := conf.(TimestampInserterGoConfig)
	if !ok {
		log.Errorf("invalid configuration type: %T", conf)
		return
	}

	// Call the timestamp service
	timestamp, err := fetchTimestamp(config.TimestampServiceURI)
	if err != nil {
		log.Errorf("failed to fetch timestamp: %s", err)
		return
	}

	// Add the timestamp to the request header
	r.Header().Set("X-Timestamp", timestamp.Time)
}

// ResponseFilter executes during the response phase.
func (t *TimestampInserterGo) ResponseFilter(conf interface{}, w pkgHTTP.Response) {
	// This plugin doesn't need to modify the response.
}

// fetchTimestamp makes an HTTP GET request to the timestamp service.
func fetchTimestamp(uri string) (TimestampResponse, error) {
	timestampResponse := TimestampResponse{}
	resp, err := http.Get(uri)
	if err != nil {
		return timestampResponse, fmt.Errorf("error calling timestamp service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return timestampResponse, fmt.Errorf("timestamp service returned non-200 status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return timestampResponse, fmt.Errorf("error reading timestamp service response: %w", err)
	}

	err = json.Unmarshal(bodyBytes, &timestampResponse)
	if err != nil {
		return timestampResponse, fmt.Errorf("error unmarshalling timestamp service response: %w", err)
	}

	return timestampResponse, nil
}
