package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func GetOffer(api string, offerID string) (offer Offer, err error) {
	response, err := http.Get(api + "/offers/" + offerID)
	if err != nil {
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	var unstructured interface{}
	err = json.Unmarshal(body, &unstructured)
	if err != nil {
		return
	}
	offer, err = ParseOffer(unstructured)
	if err != nil {
		return
	}
	return
}
