package config

import "os"

func GetServerPort() string {
	portEnv := os.Getenv("PORT")
	if portEnv == "" {
		portEnv = "5000"
	}
	return portEnv
}
