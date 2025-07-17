package gabi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/cristianoveiga/gabi-cli/cmd/gabi/utils"
)

type row map[string]string

// OutputQuery represents the structure which Gabi expects when receiving a query as a JSON file.
type OutputQuery struct {
	Query string `json:"query"`
}

type QueryService struct {
	client *Client
}

func NewQueryService(c *Client) QueryService {
	return QueryService{client: c}
}

func (s QueryService) Query(q string, output string, showRowCount bool) error {
	// Marshal the given query to JSON. This way, any specified double quotes will be escaped for us.
	marshalledQuery, err := json.Marshal(OutputQuery{Query: q})
	if err != nil {
		return err
	}

	log.Debugf("Marshalled query to be sent to Gabi: %s", marshalledQuery)

	url := s.client.baseURL + "/query"
	req, err := http.NewRequest("POST", url, bytes.NewReader(marshalledQuery))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+s.client.Token)
	res, err := s.client.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 403 {
		return errors.New("unable to execute query. Error: 403 Forbidden. Please review your token and try again")
	}

	if res.StatusCode == 401 {
		return errors.New("unable to execute query. Error: 401 Unauthorized. Please review your permissions and try again")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// in some cases, gabi returns the error message in the response body
	if strings.Contains(string(body), "ERROR") {
		return errors.New(string(body))
	}

	if strings.Contains(string(body), "The service instance has expired") {
		return errors.New(strings.TrimSpace(string(body)) +
			" -- Please extend the `expirationDate` in the deployment config.")
	}

	dic := map[string]interface{}{}
	if err = json.Unmarshal(body, &dic); err != nil {
		return err
	}

	// the error message can also be embedded in the json response
	resultErr := dic["error"]
	if resultErr != "" {
		return errors.New(fmt.Sprintf("%v", resultErr))
	}

	result := dic["result"]
	switch output {
	case "raw":
		printRaw(result)
	case "json":
		printJson(result)
	case "csv":
		printCSV(result)
	}

	if showRowCount {
		printRowCount(result)
	}

	return nil
}

// printRowCount prints the number of rows returned by the query
func printRowCount(result interface{}) {
	fmt.Printf("your query returned %d row(s)\n", getRowsCount(result))
}

// printRaw returns the raw gabi's response
func printRaw(result interface{}) {
	utils.PrettyPrint(result)
}

// printJson parses the gabi's output and formats it nicely in
// an array of objects in the `key:value` format
func printJson(result interface{}) {
	rows, _ := getRowsAndAttrs(result)
	if len(rows) == 0 {
		fmt.Println("your query didn't return any results")
		return
	}
	utils.PrettyPrint(rows)
}

// printCSV parses gabi's output and formats it as CSV
func printCSV(result interface{}) {
	rows, attrs := getRowsAndAttrs(result)
	if len(rows) == 0 {
		fmt.Println("your query didn't return any results")
		return
	}
	header := strings.Join(attrs, ",")
	fmt.Println(header)
	for _, item := range rows {
		var vals []string
		for _, a := range attrs {
			vals = append(vals, item[a])
		}
		fmt.Println(strings.Join(vals, ","))
	}
}

func getRowsAndAttrs(result interface{}) ([]row, []string) {
	var attrs []string
	var rows []row
	// gabi returns an array of arrays and we need some logic to
	// combine the values from the first array item, which contains the attributes
	// with the remaining items (actual database rows)
	if reflect.TypeOf(result).Kind() == reflect.Slice {
		resultVal := reflect.ValueOf(result)
		// handles the case when the query didn't return any rows (only the header)
		if resultVal.Len() == 1 {
			return []row{}, []string{}
		}
		for i := 0; i < resultVal.Len(); i++ {
			internalVals := reflect.ValueOf(resultVal.Index(i).Interface())
			if i == 0 {
				attrs = getAttrs(internalVals)
			} else {
				rows = append(rows, getRow(internalVals, attrs))
			}
		}
	}
	return rows, attrs
}

func getAttrs(collection reflect.Value) []string {
	var attrs []string
	for i := 0; i < collection.Len(); i++ {
		v := fmt.Sprintf("%s", collection.Index(i))
		attrs = append(attrs, v)
	}
	return attrs
}

func getRow(collection reflect.Value, attrs []string) row {
	r := row{}
	for i := 0; i < collection.Len(); i++ {
		v := fmt.Sprintf("%s", collection.Index(i))
		r[attrs[i]] = v
	}
	return r
}

func getRowsCount(result interface{}) int {
	rowCount := 0
	if reflect.TypeOf(result).Kind() == reflect.Slice {
		resultVal := reflect.ValueOf(result)
		// removes the header from the rowcount
		if resultVal.Len() > 1 {
			rowCount = resultVal.Len() - 1
		}
	}
	return rowCount
}
