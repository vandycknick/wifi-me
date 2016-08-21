package main

import (
	"flag"
	"fmt"
	"os/exec"
	"regexp"

	keyring "github.com/nickvdyck/wifi-me/keyring"
)

// AIRPORTSERVICE write some more docs
const AIRPORTSERVICE string = "AirPort"

// AIRPORTCMD write some more docs
const AIRPORTCMD string = "/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport"

func main() {

	version := keyring.GetVersion()

	fmt.Println(version)

	var list = flag.Bool("list", false, "list all ssid that are saved in your keychain account")

	flag.Parse()

	if *list {
		accounts, err := keyring.GetGenericPasswordAccounts(AIRPORTSERVICE)

		if err != nil {
			//handle error
			return
		}

		for _, account := range accounts {
			fmt.Println(account)
		}
	} else {

		result, _ := exec.Command(AIRPORTCMD, "-I").Output()
		re := regexp.MustCompile(` SSID: ([\w,-]*)`)

		matches := re.FindStringSubmatch(string(result))

		ssid := matches[1]

		// ssid = "ALICUDI"

		fmt.Println("Getting wifi password for ssid: " + ssid)
		password, _ := getPasswordForSSID(ssid)
		fmt.Println(password)
	}
}

func getPasswordForSSID(ssid string) (string, error) {
	query := keyring.NewItem()
	query.SetSecClass(keyring.SecClassGenericPassword)
	query.SetService(AIRPORTSERVICE)
	query.SetAccount(ssid)
	query.SetMatchLimit(keyring.MatchLimitOne)
	query.SetReturnData(true)

	results, err := keyring.QueryItem(query)

	if err != nil {
		fmt.Println(err)
		return "", err
	} else if len(results) != 1 {
		fmt.Println("no results")
		// Not found
		return "", nil
	} else {
		password := string(results[0].Data)

		return password, nil
	}
}
