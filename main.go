package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	logFile, err := os.OpenFile("jos-api.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logFile.Close()
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	result, err := doRequest("https://joshome.jd.com/api/index")
	if err != nil {
		log.Fatalln(err)
	}
	data := result["data"].([]interface{})
	for _, item := range data {
		itemMap := item.(map[string]interface{})
		id := int(itemMap["id"].(float64))
		log.Printf("id: %d, groupName: %s\n", id, itemMap["groupName"])
		getApiList(id)
	}
}

func doRequest(url string) (map[string]interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getApiList(groupId int) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in getApiList", r)
		}
	}()
	result, err := doRequest("https://joshome.jd.com/api/list?id=" + strconv.Itoa(groupId))
	if err != nil {
		return
	}
	data := result["data"].(map[string]interface{})
	if len(data) == 0 {
		return
	}
	cmsApis := data["cmsApis"]
	if cmsApis == nil {
		return
	}
	cmsApiArray := cmsApis.([]interface{})
	for _, api := range cmsApiArray {
		apiMap := api.(map[string]interface{})
		id := int(apiMap["id"].(float64))
		apiName := apiMap["apiName"].(string)
		log.Printf("id: %d, apiName: %s, dese: %s\n", id, apiName, apiMap["apiDesc"])
		getApiDetail(id, apiName)
	}
	time.Sleep(1 * time.Second)
}

func getApiDetail(id int, name string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in getApiList", r)
		}
	}()
	result, err := doRequest("https://joshome.jd.com/api/detail?id=" + strconv.Itoa(id) + "&apiName=" + name)
	if err != nil {
		return
	}
	success := result["success"].(bool)
	if !success {
		return
	}
	data := result["data"].(map[string]interface{})
	method := data["method"].(map[string]interface{})
	josResult := method["josResult"].(map[string]interface{})
	printElements(1, josResult["elements"])
	time.Sleep(1 * time.Second)
}

func printElements(level int, elements interface{}) {
	if elements == nil {
		return
	}
	elementArray := elements.([]interface{})
	for _, elem := range elementArray {
		elemMap := elem.(map[string]interface{})
		log.Printf("%sname: %s, type: %s, value: %s, desc: %s\n", strings.Repeat("- ", level), elemMap["paramName"], elemMap["type"], elemMap["value"], elemMap["desc"])
		newLevel := level + 1
		printElements(newLevel, elemMap["elements"])
	}
}
