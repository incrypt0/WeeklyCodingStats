package main

import (
	"bytes"
	"encoding/json"
	"errors"
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

var GIST_ID string = "b6786ee02c58a21103bb7112be12163c"

func main() {

	langDataGraph, _ := GetDataFromEmbedUrl()
	fmt.Println("Begining Gist Update")
	err := gistUpdater(langDataGraph)

	fmt.Println("test 6")

	if err != nil {
		for i := 0; err.Error() == "RESP_ERROR" && i < 2; i++ {
			fmt.Println("test 7")
			err = gistUpdater(langDataGraph)
			fmt.Println("test 8")
			if err != nil {
				log.Fatal("Error while Gist Update", err)
			}
		}
		log.Fatal("Error while Gist Uopdate", err)
	}

	fmt.Println("test 9")
	fmt.Printf("Gist Updated Successfully\n")
}

//
//
//This Function Gets data from embed url
func GetDataFromEmbedUrl() (langDataGraph string, err error) {
	var i int = 0
	var resp *http.Response
	var result map[string][]Language
	var funcName string = "GetDataFromEmbedUrl : "
	for i < 2 {
		log.Println("Getting Json From WakaTime")
		resp, err = http.Get(os.Getenv("WAKATIME_EMBED_URL"))
		i++
	}

	if err != nil {
		return "", fmt.Errorf("%v failed getting json from embed url %v", funcName, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("%v failed to read response body %v", funcName, err)

		}

		// bodyString := string(bodyBytes)
		// log.Println(bodyString)
		err = json.Unmarshal(bodyBytes, &result)
		if err != nil {
			return "", fmt.Errorf("%v%v", funcName, err)
		}
		prettyResult, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			return "", fmt.Errorf("%v%v", funcName, err)
		}
		fmt.Printf("%s\n", string(prettyResult))
	}
	langDataSlice := result["data"]
	langDataGraph = langGraphGen(langDataSlice)
	return langDataGraph, nil
}

//
//
//

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
func gistUpdater(content string) (err error) {

	reqBody, err := json.Marshal(map[string]interface{}{
		"description": "Weekly Coding Stats 📈",
		"files": map[string]interface{}{
			"coding_loop": map[string]interface{}{
				"content":  content,
				"filename": "coding_loop",
			},
		},
	})

	if err != nil {
		log.Fatal("failed to marshal json", err)
	}

	req, err := http.NewRequest("PATCH", "https://api.github.com/gists/"+GIST_ID, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request %v", err)
	}
	fmt.Println("test 1")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Incrypto")
	// fmt.Println(os.Getenv("GIST_TOKEN"))
	req.SetBasicAuth("", os.Getenv("GIST_TOKEN"))
	fmt.Println("test 2")
	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	resp, err := client.Do(req)
	fmt.Println("test 3")
	if resp.StatusCode != http.StatusOK {
		return errors.New("RESP_ERROR")
	}
	if err != nil {
		return fmt.Errorf("failed to get response %v", err)
	}
	fmt.Println("test 4")
	defer resp.Body.Close()
	fmt.Println("test 5")
	return err
}
