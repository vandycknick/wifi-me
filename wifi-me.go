package main

import (
	"flag"
	"fmt"
	"os/exec"
	"regexp"
)

// AIRPORTKEYRINGREF write some more docs
const AIRPORTKEYRINGREF = "AirPort"

// AIRPORTCMD write some more docs
const AIRPORTCMD = "/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport"

func main() {

	version := GetVersion()

	var list = flag.Bool("list", false, "list all ssid that are saved in your keychain account")

	var ssid = flag.String("ssid", getCurrentSSID(), "Get the password for the given ssid, defaults to your currently selected wifi hotspot.")

	flag.Parse()

	fmt.Printf("Currently using the following version %s \n", version)

	if *list {
		accounts, err := GetGenericPasswordAccounts(AIRPORTKEYRINGREF)

		if err != nil {
			// properly handle error
			fmt.Println(err)
			return
		}

		for _, account := range accounts {
			fmt.Println(account)
		}
	} else {
		fmt.Println("Getting wifi password for ssid: " + *ssid)
		password, _ := getPasswordForSSID(ssid)
		fmt.Println(password)
	}
}

func getCurrentSSID() string {
	result, _ := exec.Command(AIRPORTCMD, "-I").Output()
	re := regexp.MustCompile(` SSID: ([\w,-]*)`)

	matches := re.FindStringSubmatch(string(result))

	ssid := matches[1]

	return ssid
}

func getPasswordForSSID(ssid *string) (string, error) {
	result, err := GetMacKeyringPassword(AIRPORTKEYRINGREF, *ssid)

	if err != nil {
		return "", err
	}

	password := string(result)
	return password, nil
}
