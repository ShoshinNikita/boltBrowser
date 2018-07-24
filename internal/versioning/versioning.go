package versioning

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// Info consist information about the last release
type Info struct {
	IsNewVersion   bool
	CurrentVersion string
	LastVersion    string
	Changes        string
	Link           string
}

// CheckVersion check the last release on GitHub
func CheckVersion(version string) (info Info, err error) {
	response := []struct {
		URL     string `json:"html_url"`
		TagName string `json:"tag_name"`
		Body    string `json:"body"`
	}{}

	r, err := http.Get("https://api.github.com/repos/ShoshinNikita/boltbrowser/releases")
	if err != nil {
		return Info{}, err
	}
	defer r.Body.Close()
	json.NewDecoder(r.Body).Decode(&response)

	if len(response) > 0 {
		// Checking of the last release
		info.IsNewVersion = (response[0].TagName != version)
		info.CurrentVersion = version
		info.LastVersion = response[0].TagName
		info.Changes = strings.Replace(response[0].Body, "\r\n", "\n", -1)
		info.Link = response[0].URL
	} else {
		return Info{}, errors.New("Error: list of releases is empty")
	}

	return info, nil
}
