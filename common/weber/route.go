package weber

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

type Route struct {
	Parsepath string `json:"parsepath"`
	Funcname  string `json:"funcname"`
}

func ParseRoute(file string, parsepath string) Route {
	jsonFile, _ := os.Open(file)
	defer jsonFile.Close()
	jsonData, _ := ioutil.ReadAll(jsonFile)

	var route []Route
	json.Unmarshal(jsonData, &route)
	// fmt.Println(route)

	for _, v := range route {
		if v.Parsepath == parsepath {

			return v
		}
	}

	return Route{}
}
func GetParsepath(str string) string {
	path1 := strings.Split(str, "?")

	return path1[0][1:]
}
