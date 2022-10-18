package exec

import (
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	gabi "github.com/cristianoveiga/gabi-cli/pkg/client"
	"github.com/cristianoveiga/gabi-cli/pkg/config"
)

// Cmd represents the execute command
var Cmd = &cobra.Command{
	Use:     "execute [string] | stdin",
	Short:   "Executes a gabi query",
	Long:    "Executes a gabi query received from a string as argument or from stdin. When using stdin, press Enter to move to the next line and then CTRL+D to execute the query (or CTRL+C to Cancel)",
	Run:     run,
	Aliases: []string{"exec"},
}

var args struct {
	json         bool
	raw          bool
	csv          bool
	showRowCount bool
}

func init() {
	flags := Cmd.Flags()
	flags.BoolVar(
		&args.raw,
		"raw",
		false,
		"Raw output",
	)
	flags.BoolVar(
		&args.csv,
		"csv",
		false,
		"CSV output",
	)
	flags.BoolVar(
		&args.showRowCount,
		"show-row-count",
		false,
		"Prints out the number of rows returned by your query",
	)
}

func run(cmd *cobra.Command, argv []string) {
	c, err := config.CurrentProfile()
	if err != nil {
		logErrAndExit(err.Error())
	}

	if valid, msg := c.IsValid(); !valid {
		logErrAndExit(msg)
	}

	var query string
	if len(argv) > 0 {
		query = argv[0]
	} else {
		input, readErr := ioutil.ReadAll(os.Stdin)
		if readErr != nil {
			logErrAndExit(readErr.Error())
		} else {
			query = string(input)
		}
	}

	query = formatQuery(query)

	// todo: define output types as enums
	output := "json"
	if args.raw {
		output = "raw"
	}
	if args.csv {
		output = "csv"
	}

	gabiCli, err := gabi.NewClient(c)
	if err != nil {
		logErrAndExit("error creating gabi client")
	}

	qs := gabi.NewQueryService(gabiCli)
	QueryErr := qs.Query(query, output, args.showRowCount)
	if QueryErr != nil {
		logErrAndExit(QueryErr.Error())
	}
}

func formatQuery(query string) string {
	q := strings.ReplaceAll(query, "\n", " ")
	q = strings.ReplaceAll(q, "\t", " ")
	return q
}

func logErrAndExit(err string) {
	log.Error(err)
	os.Exit(1)
}
