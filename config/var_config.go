package config

import "os"

var TEMPORAL_HOST_PORT = "temporal:7233"

func init() {
	if val := os.Getenv("TEMPORAL_HOST_PORT"); val != "" {
		TEMPORAL_HOST_PORT = val
	}

}
