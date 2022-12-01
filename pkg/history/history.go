package history

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"

	. "github.com/cristianoveiga/gabi-cli/pkg/common"
)

var historyFileName = "history.json"
var HistoryFile = ConfigFolder + "/" + historyFileName

type History []string

// PushQuery tries to write the query to the log history
func PushQuery(query string) {
	h, err := Read()
	if err != nil {
		log.Warning("unable to read from the log history: %s ", err)
	}
	// appends the query to the log
	h = append(h, query)

	err = writeHistory(h)
	if err != nil {
		log.Warning("unable to write to the log history: %s ", err)
	}
}

// Read reads the query history from the history file. If it doesn't exist, it attempts to create one.
func Read() (History, error) {
	var history History
	content, err := ioutil.ReadFile(HistoryFile)
	if err != nil {
		initErr := initHistory()
		if initErr != nil {
			return nil, err
		}
	} else {
		err = json.Unmarshal(content, &history)
		if err != nil {
			log.Fatal("Error during Unmarshal: ", err)
		}
	}
	return history, nil
}

func Clear() error {
	err := writeHistory(History{})
	if err != nil {
		return errors.New("failed to clear the history")
	}
	return nil
}

func initHistory() error {
	_, err := os.Stat(ConfigFolder)
	if os.IsNotExist(err) {
		mkdirErr := os.Mkdir(ConfigFolder, 0700)
		if mkdirErr != nil {
			return mkdirErr
		}
	}
	log.Warning("history file doesn't exist. Will try to create it.")
	return writeHistory(History{})
}

func writeHistory(h History) error {
	data, _ := json.MarshalIndent(h, "", "  ")
	return ioutil.WriteFile(HistoryFile, data, 0644)
}
