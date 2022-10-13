package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

var homeFolder = os.Getenv("HOME")
var configFolder = homeFolder + "/.config/gabi"
var configFileName = "gabi.json"
var configFile = configFolder + "/" + configFileName

func Init() error {
	if validConfig() {
		log.Warning("A valid configuration for gabi already exists. Check %s for details", configFile)
		return nil
	}

	var profiles Profiles
	profile := Profile{
		Name:    "default",
		Alias:   "default",
		URL:     "",
		Token:   "",
		Current: true,
	}
	profiles = append(profiles, &profile)
	err := writeConfigs(profiles)
	if err != nil {
		return err
	}

	fmt.Printf("Gabi init success! Check %s for details and complete the setup. \n", configFile)
	return nil
}

func CurrentProfile() (Profile, error) {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Error when opening file %s. Check your settings or run `gabi config init` to get started. \n", err)
	}

	profiles, err := AllProfiles()
	err = json.Unmarshal(content, &profiles)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	for _, profile := range profiles {
		if profile.Current {
			return *profile, nil
		}
	}

	return Profile{}, errors.New("no active profile found")
}

func AllProfiles() (Profiles, error) {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Error when opening file %s. Check your settings or run `gabi config init` to get started. \n", err)
	}

	var profiles Profiles
	err = json.Unmarshal(content, &profiles)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	if len(profiles) == 0 {
		return Profiles{}, errors.New("no profiles found. Check your settings or run `gabi config init` to get started")
	}

	return profiles, nil
}

func SetCurrentProfile(alias string) error {
	allProfiles, err := AllProfiles()
	if err != nil {
		return err
	}

	found := false
	for _, profile := range allProfiles {
		if profile.Alias == alias {
			profile.Current = true
			found = true
		} else {
			profile.Current = false
		}
	}

	if !found {
		return errors.New(fmt.Sprintf("profile %s not found", alias))
	}

	writeErr := writeConfigs(allProfiles)
	if writeErr != nil {
		return err
	}

	fmt.Printf("current profile set to %s! \n", alias)
	return nil
}

func SetToken(token string) error {
	allProfiles, err := AllProfiles()
	if err != nil {
		return err
	}

	for _, profile := range allProfiles {
		if profile.Current {
			profile.Token = token
			break
		}
	}

	writeErr := writeConfigs(allProfiles)
	if writeErr != nil {
		return errors.New("failed to set token")
	}

	fmt.Println("Token successfully updated!")
	return nil
}

func SetURL(url string) error {
	allProfiles, err := AllProfiles()
	if err != nil {
		return err
	}

	for _, profile := range allProfiles {
		if profile.Current {
			profile.URL = url
			break
		}
	}

	writeErr := writeConfigs(allProfiles)
	if writeErr != nil {
		return errors.New("failed to set url")
	}

	fmt.Println("URL successfully updated!")
	return nil
}

func writeConfigs(profiles Profiles) error {
	_, err := os.Stat(configFolder)
	if os.IsNotExist(err) {
		mkdirErr := os.Mkdir(configFolder, 0700)
		if mkdirErr != nil {
			return mkdirErr
		}
	}
	file, _ := json.MarshalIndent(profiles, "", "  ")
	return ioutil.WriteFile(configFile, file, 0644)
}

func validConfig() bool {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return false
	}

	var profiles Profiles
	err = json.Unmarshal(content, &profiles)
	if err != nil {
		return false
	}

	if len(profiles) == 0 {
		return false
	}

	return true
}
