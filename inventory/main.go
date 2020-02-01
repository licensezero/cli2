package inventory

import (
	"licensezero.com/cli2/abstract"
	"licensezero.com/cli2/api"
	"licensezero.com/cli2/user"
	"os"
	"path"
)

// Inventory describes offers to license artifacts in a working directory.
type Inventory struct {
	Licensable []Item
	Licensed   []Item
	Own        []Item
	Unlicensed []Item
	Ignored    []Item
	Invalid    []Item
}

// Item describes an artifact with an offer.
type Item struct {
	Type    string
	Path    string
	Scope   string
	Name    string
	Version string
	Public  string
	API     string
	OfferID string
	Offer   abstract.Offer
}

type finding struct {
	Type    string
	Path    string
	Scope   string
	Name    string
	Version string
	Public  string
	API     string
	OfferID string
}

// CompileInventory discovers artifacts with offers in a working directory.
func CompileInventory(
	configPath string,
	cwd string,
	ignoreNoncommercial bool,
	ignoreReciprocal bool,
) (inventory *Inventory, err error) {
	// TODO: Don't ignore receipt read errors.
	receipts, _, err := user.ReadReceipts(configPath)
	if err != nil {
		return
	}
	// TODO: Don't ignore account read errors.
	accounts, _, err := user.ReadAccounts(configPath)
	findings, err := find(cwd)
	if err != nil {
		return
	}
	for _, finding := range findings {
		offer, err := api.GetOffer(finding.API, finding.OfferID)
		var item Item
		if err != nil {
			inventory.Invalid = append(inventory.Invalid, Item{
				Type:    finding.Type,
				Path:    finding.Path,
				Scope:   finding.Scope,
				Name:    finding.Name,
				Version: finding.Version,
				Public:  finding.Public,
			})
			continue
		} else {
			item = Item{
				Type:    finding.Type,
				Path:    finding.Path,
				Scope:   finding.Scope,
				Name:    finding.Name,
				Version: finding.Version,
				Public:  finding.Public,
				Offer:   offer,
			}
			inventory.Licensable = append(inventory.Licensable, item)
		}
		if haveReceipt(&item, receipts) {
			inventory.Licensed = append(inventory.Licensed, item)
			continue
		}
		if ownProject(&item, accounts) {
			inventory.Own = append(inventory.Own, item)
			continue
		}
		licenseType := licenseTypeOf(item.Public)
		if (licenseType == noncommercial) && ignoreNoncommercial {
			inventory.Ignored = append(inventory.Ignored, item)
			continue
		}
		if (licenseType == reciprocal) && ignoreReciprocal {
			inventory.Ignored = append(inventory.Ignored, item)
			continue
		}
		inventory.Unlicensed = append(inventory.Unlicensed, item)
	}
	return
}

func find(cwd string) (findings []finding, err error) {
	finders := []func(string) ([]finding, error){
		// findNPMPackages,
		// findRubyGems,
		// findGoDeps,
		// findCargoCrates,
		findLicenseZeroFiles,
	}
	for _, finder := range finders {
		findings, err := finder(cwd)
		if err == nil {
			for _, finding := range findings {
				if alreadyHave(findings, &finding) {
					continue
				}
				findings = append(findings, finding)
			}
		}
	}
	return
}

func alreadyHave(findings []finding, finding *finding) bool {
	api := finding.API
	offerID := finding.OfferID
	for _, other := range findings {
		if other.API == api && other.OfferID == offerID {
			return true
		}
	}
	return false
}

func haveReceipt(item *Item, receipts []abstract.Receipt) bool {
	api := item.API
	offerID := item.OfferID
	for _, account := range receipts {
		if account.API() == api && account.OfferID() == offerID {
			return true
		}
	}
	return false
}

func ownProject(item *Item, accounts []abstract.Account) bool {
	api := item.API
	licensorID := item.Offer.LicensorID()
	for _, account := range accounts {
		if account.API() == api && account.LicensorID() == licensorID {
			return true
		}
	}
	return false
}

// Like ioutil.ReadDir, but don't sort, and read all symlinks.
func readAndStatDir(directoryPath string) ([]os.FileInfo, error) {
	directory, err := os.Open(directoryPath)
	if err != nil {
		return nil, err
	}
	entries, err := directory.Readdir(-1)
	directory.Close()
	if err != nil {
		return nil, err
	}
	returned := make([]os.FileInfo, len(entries))
	for i, entry := range entries {
		if isSymlink(entry) {
			linkPath := path.Join(directoryPath, entry.Name())
			targetPath, err := os.Readlink(linkPath)
			if err != nil {
				return nil, err
			}
			if !path.IsAbs(targetPath) {
				targetPath = path.Join(path.Dir(directoryPath), targetPath)
			}
			newEntry, err := os.Stat(targetPath)
			if err != nil {
				return nil, err
			}
			returned[i] = newEntry
		} else {
			returned[i] = entry
		}
	}
	return returned, nil
}

func isSymlink(entry os.FileInfo) bool {
	return entry.Mode()&os.ModeSymlink != 0
}
