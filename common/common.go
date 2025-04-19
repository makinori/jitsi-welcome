package common

import "os"

var (
	_, ConfigInDev = os.LookupEnv("DEV")

	ConfigHTTPPort = getEnv("PORT", "8080")

	ConfigAniListUsername = getEnv("ANILIST_USERNAME", "makinori")

	ConfigCacheJSONPath = getEnv("CACHE_PATH", "./cache.json")
)

func getEnv(key string, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	} else {
		return value
	}
}
