package user

import (
	"encoding/json"
	"io/ioutil"
	"licensezero.com/cli2/abstract"
	"os"
	"path"
)

// ReadAccounts reads all accounts in the configuration directory.
func ReadAccounts(configPath string) (accounts []abstract.Account, errors []error, err error) {
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

func readAccount(filePath string) (account *abstract.Account, err error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &account)
	return
}
