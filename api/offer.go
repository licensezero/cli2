package api

import (
	"encoding/json"
	"io/ioutil"
	"licensezero.com/cli2/abstract"
	"net/http"
)

func GetOffer(api string, offerID string) (offer abstract.Offer, err error) {
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
	offer, err = abstract.ParseOffer(unstructured)
	if err != nil {
		return
	}
	return
}
