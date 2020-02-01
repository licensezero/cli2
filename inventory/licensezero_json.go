package inventory

import (
	"encoding/json"
	"errors"
	"github.com/yookoala/realpath"
	"io/ioutil"
	"licensezero.com/cli2/schemas"
	"os"
	"path"
)

func findLicenseZeroFiles(cwd string) (findings []finding, err error) {
	entries, err := readAndStatDir(cwd)
	if err != nil {
		if os.IsNotExist(err) {
			return findings, nil
		}
		return nil, err
	}
	for _, entry := range entries {
		name := entry.Name()
		if name == "licensezero.json" {
			findings, err := ReadLicenseZeroJSON(cwd)
			if err != nil {
				return nil, err
			}
			for _, finding := range findings {
				if alreadyHave(findings, &finding) {
					continue
				}
				packageInfo := findPackageInfo(cwd)
				if packageInfo != nil {
					finding.Type = packageInfo.Type
					finding.Name = packageInfo.Name
					finding.Version = packageInfo.Version
					finding.Scope = packageInfo.Scope
				}
				findings = append(findings, finding)
			}
		} else if entry.IsDir() {
			directory := path.Join(cwd, name)
			below, err := findLicenseZeroFiles(directory)
			if err != nil {
				return nil, err
			}
			findings = append(findings, below...)
		}
	}
	return
}

func findPackageInfo(directoryPath string) *finding {
	approaches := []func(string) *finding{
		findNPMPackageInfo,
		// findPythonPackageInfo,
		// findMavenPackageInfo,
		// findComposerPackageInfo,
	}
	for _, approach := range approaches {
		returned := approach(directoryPath)
		if returned != nil {
			return returned
		}
	}
	return nil
}

// LocalFindings reads project metadata from various files.
func LocalFindings(directoryPath string) (findings []finding, err error) {
	var hadFindings = 0
	var readerFunctions = []func(string) ([]finding, error){
		ReadLicenseZeroJSON,
		// ReadCargoTOML,
	}
	for _, readerFunction := range readerFunctions {
		projects, err := readerFunction(directoryPath)
		if err == nil {
			hadFindings = hadFindings + 1
			findings = projects
		}
	}
	if hadFindings > 1 {
		return nil, errors.New("multiple metadata files")
	}
	return
}

// ReadLicenseZeroJSON reads metadata from licensezero.json.
func ReadLicenseZeroJSON(directoryPath string) (findings []finding, err error) {
	jsonFile := path.Join(directoryPath, "licensezero.json")
	data, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return nil, err
	}
	var unstructured interface{}
	json.Unmarshal(data, &unstructured)
	parsed, err := schemas.ParseArtifactMetadata(unstructured)
	for _, offer := range parsed.Offers {
		item := finding{
			Path:    directoryPath,
			API:     offer.API,
			OfferID: offer.OfferID,
		}
		realDirectory, err := realpath.Realpath(directoryPath)
		if err != nil {
			item.Path = realDirectory
		} else {
			item.Path = directoryPath
		}
		findings = append(findings, item)
	}
	return findings, nil
}
