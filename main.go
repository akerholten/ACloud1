// Nikolai Åkerholt, studentid:473184

//Sources:
//


package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"strings"
	"log"
	"time"
	"io/ioutil"
)

type RepoData struct {
	Project 			string 	`json:"name"`
	Owner				Owner	`json:"owner"`
	ContributorsLink	string  `json:"contributors_url"`
	LanguageLink 		string 	`json:"languages_url"`
}

type ProcessedData struct {
	Project 		string 		`json:"project"`
	Owner			string		`json:"owner"`
	Committer		string 		`json:"committer"`
	Commits 		int			`json:"commits"`
	Languages 		[]string 	`json:"language"`
}

type Owner struct {
	Login 		string 	`json:"login"`
}

type Contributer struct {
	Login 			string 	`json:"login"`
	Contributions 	int 	`json:"contributions"`
}


func handlerGitURL(w http.ResponseWriter, r *http.Request){
	http.Header.Add(w.Header(), "content-type", "application/json")
	gitRepo := strings.Split(r.URL.Path, "/")
	if len(gitRepo) >= 5 && strings.Compare(gitRepo[2], "github.com") == 0 {
		url := "https://api." + gitRepo[2] + "/repos/" + gitRepo[3] + "/" + gitRepo[4]

		myClient := http.Client{ Timeout: time.Second * 2 }

		req, err := http.NewRequest(http.MethodGet, url, nil)

		if err!= nil {
			log.Fatal(err)
		}

		req.Header.Set("User-Agent", "akerholten" )


		resp, err := myClient.Do(req)

		if err!= nil {
			log.Fatal(err)
		}


		body, readErr :=ioutil.ReadAll(resp.Body)

		if readErr!= nil {
			log.Fatal(readErr)
		}


		repoData := RepoData{}
		jsonErr := json.Unmarshal(body, &repoData)

		if jsonErr!= nil {
			log.Fatal(jsonErr)
		}

		contData := getContributors(repoData.ContributorsLink, myClient)

		langData := getLanguages(repoData.LanguageLink, myClient)

		json.Marshal(&repoData)

		processedData := ProcessedData{}

		processedData.Project 	= gitRepo[2] + "/" + gitRepo[3] + "/" + gitRepo[4]
		processedData.Owner 	= repoData.Owner.Login
		processedData.Committer = contData.Login
		processedData.Commits 	= contData.Contributions
		processedData.Languages = langData


		json.Marshal(&processedData)
		json.NewEncoder(w).Encode(processedData)


	} else {
		fmt.Fprintf(w, "not valid")
	}
}

func getContributors (url string, myClient http.Client) Contributer {
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err!= nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "akerholten" )


	resp, err := myClient.Do(req)

	if err!= nil {
		log.Fatal(err)
	}


	body, readErr :=ioutil.ReadAll(resp.Body)

	if readErr!= nil {
		log.Fatal(readErr)
	}


	contData := []Contributer{}
	jsonErr := json.Unmarshal(body, &contData)

	if jsonErr!= nil {
		log.Fatal(jsonErr)
	}

	return contData[0]
}

func getLanguages (url string, myClient http.Client) []string {
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err!= nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "akerholten" )


	resp, err := myClient.Do(req)

	if err!= nil {
		log.Fatal(err)
	}


	body, readErr :=ioutil.ReadAll(resp.Body)

	if readErr!= nil {
		log.Fatal(readErr)
	}


	var Languages []string
	LanguageMap := make(map [string] int)
	json.Unmarshal(body, &LanguageMap)
	for key := range LanguageMap{
		Languages = append(Languages, key)
	}

	return Languages
}

func main() {
	http.HandleFunc("/giturl/", handlerGitURL)
	http.ListenAndServe("localhost:8080", nil)
}