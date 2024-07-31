package db

import (
	"fmt"
	"os"

	gabi "github.com/cristianoveiga/gabi-cli/pkg/client"
	"github.com/cristianoveiga/gabi-cli/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "db",
	Short: "Database related commands",
	Long:  `Database related commands`,
}

var switchCmd = &cobra.Command{
	Use:   "switch [dbname]",
	Short: "Switches the current database name",
	Long:  `Switches the current database name`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.CurrentProfile()
		if err != nil {
			logErrAndExit(err.Error())
		}

		if valid, msg := c.IsValid(); !valid {
			logErrAndExit(msg)
		}

		client, err := gabi.NewClient(c)
		if err != nil {
			logErrAndExit("error creating gabi client")
		}

		dbNameService := gabi.NewDBNameService(client)
		dbName := args[0]
		err = dbNameService.SwitchDBName(dbName)
		if err != nil {
			logErrAndExit("Failed to switch database name: " + err.Error())
		}
		fmt.Println("Database name switched to:", dbName)
	},
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets the current database name",
	Long:  `Gets the current database name`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.CurrentProfile()
		if err != nil {
			logErrAndExit(err.Error())
		}

		if valid, msg := c.IsValid(); !valid {
			logErrAndExit(msg)
		}

		client, err := gabi.NewClient(c)
		if err != nil {
			logErrAndExit("error creating gabi client")
		}

		dbNameService := gabi.NewDBNameService(client)
		dbName, err := dbNameService.GetDBName()
		if err != nil {
			logErrAndExit("Failed to get database name: " + err.Error())
		}
		fmt.Println("Current database name:", dbName)
	},
}

func init() {
	Cmd.AddCommand(switchCmd)
	Cmd.AddCommand(getCmd)
}

func logErrAndExit(err string) {
	log.Error(err)
	os.Exit(1)
}
