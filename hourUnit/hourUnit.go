package hourUnit

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Data struct {
	Data struct {
		Languages []Language `json:"languages"`
	} `json:"data"`
}
type Language struct {
	Name    string  `json:"name"`
	Percent float32 `json:"percent"`
	Text    string  `json:"text"`
}

func GetHours(WAKATIME_API_KEY string) (resultBody Data, err error) {

	req, err := http.NewRequest("GET", "https://wakatime.com/api/v1/users/current/stats/last_7_days", nil)

	if err != nil {
		return Data{}, fmt.Errorf("failed to create request %v", err)
	}
	q := req.URL.Query()
	q.Add("apikey", WAKATIME_API_KEY)

	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return Data{}, fmt.Errorf("failed to get response %v", err)
	}
	fmt.Println(req.URL.String())
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&resultBody)
	if err != nil {
		return Data{}, fmt.Errorf("decoding json failed %v", err)
	}

	prettyResult, err := json.MarshalIndent(resultBody, "", "    ")
	if err != nil {
		return Data{}, fmt.Errorf("json marshal failed %v", err)
	}

	fmt.Printf("%s\n", string(prettyResult))

	return resultBody, nil
}
