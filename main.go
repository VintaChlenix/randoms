package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const RandomURL = "https://api.random.org/json-rpc/2/invoke"

const APIKey = ""

type random struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
	ID      int    `json:"id"`
}
type Params struct {
	APIKey      string `json:"apiKey"`
	N           int    `json:"n"`
	Min         int    `json:"min"`
	Max         int    `json:"max"`
	Replacement bool   `json:"replacement"`
}

type result struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Random struct {
			Data           []int  `json:"data"`
			CompletionTime string `json:"completionTime"`
		} `json:"random"`
		BitsUsed      int `json:"bitsUsed"`
		BitsLeft      int `json:"bitsLeft"`
		RequestsLeft  int `json:"requestsLeft"`
		AdvisoryDelay int `json:"advisoryDelay"`
	} `json:"result"`
	ID int `json:"id"`
}

func getRandomNumber(min, max int) (int, error) {
	myrand := random{
		Jsonrpc: "2.0",
		Method:  "generateIntegers",
		Params: Params{
			APIKey:      APIKey,
			N:           1,
			Min:         min,
			Max:         max,
			Replacement: false,
		},
		ID: 16,
	}

	tmp, err := json.Marshal(myrand)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	r := bytes.NewBuffer(tmp)
	fmt.Println(string(tmp))
	resp, err := http.Post(RandomURL, "application/json", r)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	fmt.Println(string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return 0, fmt.Errorf("Response error: %s", resp.Status)
	}

	var result result
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		resp.Body.Close()
		return 0, err
	}

	resp.Body.Close()

	return result.Result.Random.Data[0], nil
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	res := 0
	tmpl, _ := template.ParseFiles("mainPage.html")
	switch r.Method {
	case "GET":
		tmpl.Execute(w, res)

	case "POST":
		min := r.FormValue("Min number")
		max := r.FormValue("Max number")
		m1, err := strconv.Atoi(min)
		if err != nil {
			log.Println(err)
		}
		m2, err := strconv.Atoi(max)
		if err != nil {
			log.Println(err)
		}
		if m1 >= m2 {
			m2 += m1
		}
		res, err := getRandomNumber(m1, m2)
		if err != nil {
			log.Println(err)
		}
		tmpl.Execute(w, res)
	}
}

func main() {
	http.HandleFunc("/", mainPage)
	http.ListenAndServe(":8080", nil)

}
