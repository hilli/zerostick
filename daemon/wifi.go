package zerostick

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/pbkdf2"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"
)

// Wifi properties struct
type Wifi struct {
	SSID              string `json:"ssid"`
	Password          string `json:",omitempty" yaml:",omitempty"`
	EncryptedPassword string `json:"encrypted_password"`
	Priority          int    `json:"priority"`
	UseForSync        bool   `json:"use_for_sync"`
}

// Wifis is a slice of Wifi
type Wifis struct {
	List []Wifi `json:"wifis"`
}

// WpaNetwork defines a wifi network to connect to.
type WpaNetwork struct {
	Bssid       string `json:"bssid"`
	Frequency   string `json:"frequency"`
	SignalLevel string `json:"signal_level"`
	Flags       string `json:"flags"`
	Ssid        string `json:"ssid"`
}

// GetWifiConfig returns the Wifi as wpa_supplicant.conf block
func (w Wifi) String() string {
	return fmt.Sprintf("network={\n\tssid\"%s\"\n\tpsk=%s\n\tpriority=%d\n}\n", w.SSID, w.EncryptedPassword, w.Priority)
}

// AddWifiToList appends the given Wifi to the list
func (ws *Wifis) AddWifiToList(w Wifi) {
	ws.DeleteWifiFromList(w.SSID) // Remove the SSID from the list first
	ws.List = append(ws.List, w)
	viper.Set("wifis", ws.List)
	viper.WriteConfig()
}

// DeleteWifiFromList will remove a wifi (given the SSID string) from the list if it's there
func (ws *Wifis) DeleteWifiFromList(SSID string) {
	changed := false
	for k, v := range ws.List {
		if v.SSID == SSID {
			ws.List = append(ws.List[:k], ws.List[k+1:]...)
			changed = true
			// log.Debugf("DELETE $d element leaves %+v", k, ws.List)
		}
	}
	if changed {
		viper.Set("wifis", ws.List)
		viper.WriteConfig()
	}
}

// encryptPassword returns password as a WPA2 formatted hash
func (w Wifi) encryptPassword(ssid string, password string) string {
	dk := pbkdf2.Key([]byte(password), []byte(ssid), 4096, 256, sha1.New)
	WPAKey := hex.EncodeToString(dk)[0:64] // First 64 bytes of hex string
	return WPAKey
}

// EncryptPassword encrypts the password in the Wifi struct and removes the unencrypted password
func (w *Wifi) EncryptPassword() {
	w.EncryptedPassword = w.encryptPassword(w.SSID, w.Password)
	w.Password = ""
}

// GetWpaSupplicantConf generates a wpa_supplicant.conf file
func (ws Wifis) GetWpaSupplicantConf() string {
	config := "ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=netdev\nupdate_config=1\ncountry=US\n"
	for _, w := range ws.List {
		config = config + w.String()
	}
	return config
}

// WriteConfig writes the config to disk
func (ws Wifis) WriteConfig(wpaSupplicantFile string) error {
	if wpaSupplicantFile == "" {
		wpaSupplicantFile = "/etc/wpa_supplicant/wpa_supplicant.conf"
	}
	err := ioutil.WriteFile(wpaSupplicantFile, []byte(ws.GetWpaSupplicantConf()), 0600)
	if err != nil {
		return err
	}
	return nil
}

// ScanNetworks returns a map of WpaNetwork data structures.
func ScanNetworks() (map[string]WpaNetwork, error) {
	wpaNetworks := make(map[string]WpaNetwork, 0)

	scanOut, err := exec.Command("wpa_cli", "-i", "wlan0", "scan").Output()
	if err != nil {
		//log.Fatal(err)
		return wpaNetworks, err
	}
	scanOutClean := strings.TrimSpace(string(scanOut))

	// wait one second for results
	time.Sleep(1 * time.Second)

	if scanOutClean == "OK" {
		log.Debug("OK scan")
		networkListOut, err := exec.Command("wpa_cli", "-i", "wlan0", "scan_results").Output()
		if err != nil {
			//wpa.Log.Fatal(err)
			return wpaNetworks, err
		}

		networkListOutArr := strings.Split(string(networkListOut), "\n")
		for _, netRecord := range networkListOutArr[1:] {
			if !strings.Contains(netRecord, "[WPA2-PSK-CCMP]") {
				continue
			}

			fields := strings.Fields(netRecord)

			if len(fields) > 4 {
				ssid := strings.Join(fields[4:], " ")
				wpaNetworks[ssid] = WpaNetwork{
					Bssid:       fields[0],
					Frequency:   fields[1],
					SignalLevel: fields[2],
					Flags:       fields[3],
					Ssid:        ssid,
				}
			}
		}
	}
	return wpaNetworks, nil
}

// ParseViperWifi parses the wifis from viper
func (ws *Wifis) ParseViperWifi() {
	if viper.IsSet("wifis") {
		err := viper.UnmarshalKey("wifis", &ws.List)
		if err != nil {
			log.Println("Error unmarshalling", err)
		}
	} else {
		ws.List = nil
	}
}
