package versioning

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
)

// Info consist information about the last release
type Info struct {
	IsNewVersion   bool
	CurrentVersion string
	LastVersion    string
	Changes        []string
	Link           string
}

// Transform string from
// "+ 1-st change\r\n + 2-st change\r\n + 3-st change"
// to
// [1-st change, 2-st change, 3-st change]
func getChanges(text string) (changes []string) {
	r := regexp.MustCompile(`(?:\+|\*) ?(?P<change>[\w !@#$%^&*()+/*+"№;:?*=~{}\[\],.<'>|^-]+)(?:\r\n|$)`)
	matches := r.FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		// The first group is <change>
		changes = append(changes, match[1])
	}
	return changes
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
		info.Changes = getChanges(response[0].Body)
		info.Link = response[0].URL
	} else {
		return Info{}, errors.New("Error: list of releases is empty")
	}

	return info, nil
}
