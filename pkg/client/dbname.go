package gabi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type SwitchDBNameRequest struct {
	DBName string `json:"db_name"`
}

type DBNameResponse struct {
	DBName string `json:"db_name"`
}

type DBNameService struct {
	client *Client
}

func NewDBNameService(c *Client) DBNameService {
	return DBNameService{client: c}
}

func (s DBNameService) SwitchDBName(dbName string) error {
	marshalledRequest, err := json.Marshal(SwitchDBNameRequest{DBName: dbName})
	if err != nil {
		return err
	}

	url := s.client.baseURL + "/dbname/switch"
	req, err := http.NewRequest("POST", url, bytes.NewReader(marshalledRequest))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+s.client.Token)
	res, err := s.client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("failed to switch DB name: %s", string(body))
	}

	return nil
}

func (s DBNameService) GetDBName() (string, error) {
	url := s.client.baseURL + "/dbname"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+s.client.Token)
	res, err := s.client.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		return "", fmt.Errorf("failed to get DB name: %s", string(body))
	}

	var respData DBNameResponse
	err = json.NewDecoder(res.Body).Decode(&respData)
	if err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return respData.DBName, nil
}
