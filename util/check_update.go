package util

import (
	"io/ioutil"
	"net/http"
	"strings"
)

const verUrl = "https://raw.githubusercontent.com/abs-lang/abs/master/VERSION"

// Returns latest version, plus "new version available?" bool
func UpdateAvailable(version string) (string, bool) {
	resp, err := http.Get(verUrl)
	if err != nil {
		return version, false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return version, false
	}

	latest := strings.TrimSpace(string(body))
	if version != latest {
		return latest, true
	}

	return version, false
}
