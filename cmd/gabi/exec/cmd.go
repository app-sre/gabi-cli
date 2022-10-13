package exec

import (
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/spf13/cobra"

	gabi "github.com/cristianoveiga/gabi-cli/pkg/client"
	"github.com/cristianoveiga/gabi-cli/pkg/config"
)

// Cmd represents the execute command
var Cmd = &cobra.Command{
	Use:     "execute",
	Short:   "Executes a gabi query",
	Long:    "Executes a gabi query",
	Run:     run,
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"exec"},
}

var args struct {
	json bool
	raw  bool
	csv  bool
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
}

func run(cmd *cobra.Command, argv []string) {
	c, err := config.CurrentProfile()
	if err != nil {
		logErrAndExit(err.Error())
	}

	if valid, msg := c.IsValid(); !valid {
		logErrAndExit(msg)
	}

	q := argv[0]
	if len(q) == 0 {
		logErrAndExit("You must provide a query")
	}

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
	QueryErr := qs.Query(q, output)
	if QueryErr != nil {
		logErrAndExit(QueryErr.Error())
	}
}

func logErrAndExit(err string) {
	log.Error(err)
	os.Exit(1)
}
