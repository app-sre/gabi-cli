package config

import (
	"github.com/cristianoveiga/gabi-cli/cmd/gabi/utils"
	"github.com/cristianoveiga/gabi-cli/pkg/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Cmd represents the configure command
var Cmd = &cobra.Command{
	Use:     "configure",
	Short:   "Gets or Sets the gabi-cli configs",
	Long:    `Gets or Sets the gabi-cli configs`,
	Aliases: []string{"config"},
}

// InitCmd represents the init command
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a gabi-cli config by creating a config file under the ~/.config/gabi directory",
	Long:  "Initializes a gabi-cli config by creating a config file under ~/.config/gabi directory",
	Run: func(cmd *cobra.Command, args []string) {
		err := config.Init()
		if err != nil {
			log.Error(err)
		}
	},
}

// currentProfileCmd represents the currentprofile command
var currentProfileCmd = &cobra.Command{
	Use:   "currentprofile",
	Short: "Gets the current profile",
	Long:  `Gets the current profile`,
	Run: func(cmd *cobra.Command, args []string) {
		profile, err := config.CurrentProfile()
		profile.Redact()
		if err != nil {
			log.Error(err)
		}
		utils.PrettyPrint(profile)
	},
}

// allProfilesCmd represents the allprofiles command
var allProfilesCmd = &cobra.Command{
	Use:   "allprofiles",
	Short: "Gets all profiles currently configured for gabi-cli",
	Long:  `Gets all profiles currently configured for gabi-cli`,
	Run: func(cmd *cobra.Command, args []string) {
		allProfiles, err := config.AllProfiles()
		if err != nil {
			log.Error(err)
		}
		var redactedProfiles config.Profiles
		for _, p := range allProfiles {
			p.Redact()
			redactedProfiles = append(redactedProfiles, p)
		}
		utils.PrettyPrint(redactedProfiles)
	},
}

// setTokenCmd represents the settoken command
var setTokenCmd = &cobra.Command{
	Use:   "settoken",
	Short: "Sets the token for current profile",
	Long:  "Sets the token for current profile",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		token := args[0]
		err := config.SetToken(token)
		if err != nil {
			log.Error(err)
		}
	},
}

// setURLCmd represents the settoken command
var setURLCmd = &cobra.Command{
	Use:   "seturl",
	Short: "Sets the url for current profile",
	Long:  "Sets the url for current profile",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		err := config.SetURL(url)
		if err != nil {
			log.Error(err)
		}
	},
}

// setProfileCmd represents the setprofile command
var setProfileCmd = &cobra.Command{
	Use:   "setprofile",
	Short: "Sets the current profile",
	Long:  `Sets the current profile`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		profileAlias := args[0]
		err := config.SetCurrentProfile(profileAlias)
		if err != nil {
			log.Error(err)
		}
	},
}

func init() {
	Cmd.AddCommand(InitCmd)
	Cmd.AddCommand(currentProfileCmd)
	Cmd.AddCommand(allProfilesCmd)
	Cmd.AddCommand(setTokenCmd)
	Cmd.AddCommand(setURLCmd)
	Cmd.AddCommand(setProfileCmd)

}
