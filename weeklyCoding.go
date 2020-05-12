package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Language struct {
	Name    string  `json:"name"`
	Percent float32 `json:"percent"`
}

func main() {
	// fmt.Println(os.Getenv("GIST_TOKEN"))
	var i int = 0
	var err error
	var resp *http.Response
	var result map[string][]Language

	for i < 2 {
		log.Println("Getting Json From WakaTime")
		resp, err = http.Get(os.Getenv("WAKATIME_EMBED_URL"))
		i++
	}
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		// bodyString := string(bodyBytes)
		// log.Println(bodyString)
		err = json.Unmarshal(bodyBytes, &result)
		if err != nil {
			log.Fatal(err)
		}
		prettyResult, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", string(prettyResult))
	}

	langDataSlice := result["data"]
	langDataGraph := langGraphGen(langDataSlice)

	status, err := gistUpdater(langDataGraph)
	for i := 0; status != http.StatusOK && i <= 3; i++ {
		_, err = gistUpdater(langDataGraph)
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Gist Updated Successfully\n")
}

//
//
// This Function Generates the ascii graph

func langGraphGen(langDataSlice []Language) string {
	var buf bytes.Buffer
	for _, item := range langDataSlice {

		pMod := item.Percent / 3
		fmt.Println(len(item.Name), 8-len(fmt.Sprintf("%.2f", item.Percent)))

		fmt.Fprint(&buf,
			item.Name,
			strings.Repeat(" ", 15-len(item.Name)),
			fmt.Sprintf("%.2f", item.Percent),
			strings.Repeat(" ", 8-len(fmt.Sprintf("%.2f", item.Percent))),
			strings.Repeat("█", int(pMod)),
			strings.Repeat("░", int(100/3-pMod)),
		)
		fmt.Fprint(&buf, "\n")

	}
	fmt.Println(buf.String())
	return buf.String()
}

//
//
// This Function Updates Gists
func gistUpdater(content string) (statusCode int, err error) {

	reqBody, err := json.Marshal(map[string]interface{}{
		"description": "Weekly Development Break Down",
		"files": map[string]interface{}{
			"coding_loop": map[string]interface{}{
				"content":  content,
				"filename": "coding_loop",
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("PATCH", "https://api.github.com/gists/b6786ee02c58a21103bb7112be12163c", bytes.NewBuffer(reqBody))
	if err != nil {
		return 1, fmt.Errorf("failed to create request %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Incrypto")
	// fmt.Println(os.Getenv("GIST_TOKEN"))
	req.SetBasicAuth("incrypt0", os.Getenv("GIST_TOKEN"))

	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	resp, err := client.Do(req)
	statusCode = resp.StatusCode
	if err != nil {
		return statusCode, fmt.Errorf("failed to get response %v", err)
	}

	defer resp.Body.Close()

	return statusCode, nil
}
