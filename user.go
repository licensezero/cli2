package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

// ReadAccounts reads all accounts in the configuration directory.
func ReadAccounts(configPath string) (accounts []Account, errors []error, err error) {
	directoryPath := path.Join(configPath, "accounts")
	entries, directoryReadError := ioutil.ReadDir(directoryPath)
	if directoryReadError != nil {
		if os.IsNotExist(directoryReadError) {
			return
		}
		return nil, nil, directoryReadError
	}
	for _, entry := range entries {
		name := entry.Name()
		filePath := path.Join(configPath, "accounts", name)
		receipt, err := readAccount(filePath)
		if err != nil {
			errors = append(errors, err)
		} else {
			accounts = append(accounts, *receipt)
		}
	}
	return
}

func readAccount(filePath string) (account *Account, err error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &account)
	return
}

// ReadReceipts reads all receipts in the configuration directory.
func ReadReceipts(configPath string) (receipts []Receipt, errors []error, err error) {
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
			receipts = append(receipts, receipt)
		}
	}
	return
}

func readReceipt(filePath string) (Receipt, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var unstructured interface{}
	err = json.Unmarshal(data, &unstructured)
	if err != nil {
		return nil, err
	}
	return ParseReceipt(unstructured)
}
