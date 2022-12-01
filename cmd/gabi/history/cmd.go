package history

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/cristianoveiga/gabi-cli/pkg/history"
)

var args struct {
	maxRows int
}

func init() {
	showHistoryCmd.Flags().IntVar(&args.maxRows, "max-rows", 100, "Maximum number of rows returned in the show command")
	Cmd.AddCommand(showHistoryCmd)
	Cmd.AddCommand(clearHistoryCmd)
}

// Cmd represents the history command
var Cmd = &cobra.Command{
	Use:   "history",
	Short: "Executes history operations",
	Long:  `Executes history operations`,
}

// showHistoryCmd represents the history show command
var showHistoryCmd = &cobra.Command{
	Use:   "show",
	Short: "Shows gabi-cli query history",
	Long:  `Shows gabi-cli query history`,
	Run: func(cmd *cobra.Command, argv []string) {
		rows, err := history.Read()
		if err != nil {
			log.Error(err)
		}

		if len(rows) == 0 {
			log.Warning("your query history is empty")
			return
		}

		totalRows := len(rows)
		rowCount := 0
		for i := totalRows - 1; i > -1 && rowCount < args.maxRows; i-- {
			rowCount += 1
			fmt.Printf("%d %s \n", i+1, rows[i])
		}

		if rowCount < totalRows {
			fmt.Println("[truncated content...]")
		}
	},
}

// clearHistoryCmd represents the history clear command
var clearHistoryCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clears gabi-cli query history",
	Long:  `Clears gabi-cli query history`,
	Run: func(cmd *cobra.Command, args []string) {
		err := history.Clear()
		if err != nil {
			log.Error(err)
		}
		fmt.Println("History successfully cleared!")
	},
}
