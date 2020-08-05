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

	"github.com/incrypt0/WeeklyCodingStats/hourUnit"
)

var CODING_STAT_GIST_ID string = os.Getenv("GIST_ID")
var WAKATIME_API_KEY string = os.Getenv("WAKATIME_API_KEY")
var GIST_TOKEN string = os.Getenv("GIST_TOKEN")
var WAKATIME_EMBED_URL string = os.Getenv("WAKATIME_EMBED_URL")

func main() {
	var langDataGraph string
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) > 0 && argsWithoutProg[0] == "-api" {
		fmt.Println("Using API")
		if WAKATIME_API_KEY == "" {
			log.Fatal("Error WAKA_TIME_API_KEY env Empty")
		}
		langDataGraph, _ = GetDataFromApi()
	} else {
		fmt.Println("Using Embed Url")
		if WAKATIME_EMBED_URL == "" {
			log.Fatal("Error WAKATIME_EMBED_URL env Empty")
		}
		langDataGraph, _ = GetDataFromEmbedUrl()
	}

	fmt.Println("Begining Gist Update")
	err := gistUpdater(langDataGraph)

	fmt.Println("test 6")

	if err != nil {
		for i := 0; err.Error() == "RESP_ERROR" && i < 2; i++ {
			fmt.Println("test 7")
			err = gistUpdater(langDataGraph)
			fmt.Println("test 8")
			if err != nil {
				log.Fatal("Error while Gist Update ", err)
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
	var result map[string][]hourUnit.Language
	var funcName string = "GetDataFromEmbedUrl : "
	for i < 2 {
		log.Println("Getting Json From WakaTime")
		resp, err = http.Get(WAKATIME_EMBED_URL)
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
	langDataGraph = langGraphGen(langDataSlice, false)
	return langDataGraph, nil
}

//
//
//
func GetDataFromApi() (langDataGraph string, err error) {
	data, err := hourUnit.GetHours(WAKATIME_API_KEY)
	if err != nil {
		return "", err
	}
	langDataGraph = langGraphGen(data.Data.Languages, true)
	return langDataGraph, err
}

//
//
// This Function Generates the ascii graph

func langGraphGen(langDataSlice []hourUnit.Language, isApi bool) string {
	var buf bytes.Buffer
	for _, item := range langDataSlice {
		if isApi {
			pMod := item.Percent / 3
			fmt.Println(len(item.Name), 8-len(fmt.Sprintf("%.2f", item.Percent)))
			if len(item.Name) > 15 {
				item.Name = item.Name[0:14]
			}
			fmt.Println(item.Name)
			fmt.Fprint(&buf,
				item.Name,
				strings.Repeat(" ", 15-len(item.Name)),
				fmt.Sprintf("%v", item.Text),
				strings.Repeat(" ", 15-len(fmt.Sprintf("%v", item.Text))),
				strings.Repeat("█", int(pMod)),
				strings.Repeat("░", int(100/3-pMod)),
			)
			fmt.Fprint(&buf, "\n")
		} else {
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

	}
	fmt.Println(buf.String())
	return buf.String()
}

//
//
// This Function Updates Gists
func gistUpdater(content string) (err error) {

	reqBody, err := json.Marshal(map[string]interface{}{
		"description": "📈 Weekly Coding Stats 📈",
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
	fmt.Println("Gist ID is " + CODING_STAT_GIST_ID)
	req, err := http.NewRequest("PATCH", "https://api.github.com/gists/"+CODING_STAT_GIST_ID, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Incrypto")
	// fmt.Println(os.Getenv("GIST_TOKEN"))
	req.SetBasicAuth("", GIST_TOKEN)

	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	resp, err := client.Do(req)
	fmt.Println(resp.Status)
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
