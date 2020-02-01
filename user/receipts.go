package user

import (
	"encoding/json"
	"io/ioutil"
	"licensezero.com/cli2/abstract"
	"licensezero.com/cli2/schemas"
	"os"
	"path"
)

// ReadReceipts reads all receipts in the configuration directory.
func ReadReceipts(configPath string) (receipts []abstract.Receipt, errors []error, err error) {
	directoryPath := path.Join(configPath, "receipts")
	entries, directoryReadError := ioutil.ReadDir(directoryPath)
	if directoryReadError != nil {
		if os.IsNotExist(directoryReadError) {
			return
		}
		return nil, nil, directoryReadError
	}
	for _, entry := range entries {
		name := entry.Name()
		filePath := path.Join(configPath, "receipts", name)
		receipt, err := readReceipt(filePath)
		if err != nil {
			errors = append(errors, err)
		} else {
			receipts = append(receipts, *receipt)
		}
	}
	return
}

func readReceipt(filePath string) (*abstract.Receipt, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var unstructured interface{}
	err = json.Unmarshal(data, &unstructured)
	if err != nil {
		return nil, err
	}
	return schemas.ParseReceipt(unstructured)
}
