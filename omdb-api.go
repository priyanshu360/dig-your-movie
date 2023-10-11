package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	API_KEY  = "9496f5e1"
	BASE_URL = "http://www.omdbapi.com/"
)

type Movie struct {
	Title       string `json:"Title"`
	Description string `json:"Plot"`
}

func GetMovieDescription(apiKey string, movieName string) (Movie, error) {
	queryParams := url.Values{
		"apikey": {apiKey},
		"t":      {movieName},
	}
	url := BASE_URL + "?" + queryParams.Encode()

	response, err := http.Get(url)
	if err != nil {
		return Movie{}, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return Movie{}, err
	}

	var movie Movie
	err = json.Unmarshal(body, &movie)
	if err != nil {
		return Movie{}, err
	}

	return movie, nil
}
