package togo

import (
	"encoding/json"
	"io/ioutil"
)

// Config holds info required to configure a togo.server.
type Config struct {
	HTTPAddr     string `json:"http_addr"`
	HTTPPort     int    `json:"http_port"`
	IdleTimeout  int    `json:"idle_timeout"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
	LogFilename  string `json:"log_filename"`
}

// LoadJSONFile attempts to read a specified JSON file by provided filename.
// It then attempts to unmarshal the JSON data into a Config type object.
// It returns a populated Config or any errors it encountered during the JSON file load and parse.
func LoadJSONFile(filename string) (config Config) {
	var (
		content []byte
		err     error
	)
	config = Config{
		HTTPAddr:     "0.0.0.0",
		HTTPPort:     3000,
		IdleTimeout:  30,
		ReadTimeout:  5,
		WriteTimeout: 10,
	}
	if content, err = ioutil.ReadFile(filename); err != nil {
		Log.Printf("Unable to read file %q: %s", filename, err)
		return
	}
	if err = json.Unmarshal(content, &config); err != nil {
		Log.Printf("Unable to parse JSON from file %q: %s", filename, err)
	}
	return
}
