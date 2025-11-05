package utils

import "os"

func Env(key string, defaultValue string) string {
	v, found := os.LookupEnv(key)
	if !found {
		v = defaultValue
	}
	return v
}
