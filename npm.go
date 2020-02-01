package main

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"strings"
)

type packageJSONFile struct {
	Name        string      `json:"name"`
	Version     string      `json:"version"`
	LicenseZero interface{} `json:"licensezero"`
}

func readPackageJSON(directory string) (*packageJSONFile, error) {
	packageJSON := path.Join(directory, "package.json")
	data, err := ioutil.ReadFile(packageJSON)
	if err != nil {
		return nil, err
	}
	var parsed packageJSONFile
	json.Unmarshal(data, &parsed)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func findNPMPackageInfo(directoryPath string) *finding {
	packageJSON := path.Join(directoryPath, "package.json")
	data, err := ioutil.ReadFile(packageJSON)
	if err != nil {
		return nil
	}
	var parsed struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		return nil
	}
	rawName := parsed.Name
	var name, scope string
	// If the name looks like @scope/name, parse it.
	if strings.HasPrefix(rawName, "@") && strings.Index(rawName, "/") != -1 {
		index := strings.Index(rawName, "/")
		scope = rawName[1 : index-1]
		name = rawName[index:]
	} else {
		name = rawName
	}
	return &finding{
		Type:    "npm",
		Name:    name,
		Scope:   scope,
		Version: parsed.Version,
	}
}
