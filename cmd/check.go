package cmd

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/KixPanganiban/bantay/lib"
	"github.com/KixPanganiban/bantay/log"
	"github.com/spf13/cobra"
)

// checkCmd performs uptime checks once
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Perform uptime checks once",
	Long:  `Performs all uptime checks defined in checks.yml once.`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := os.Getwd()
		checksFilePath := path.Join(dir, "checks.yml")
		checksFileBytes, err := ioutil.ReadFile(checksFilePath)
		if err != nil {
			log.Error("Unable to open checks.yml")
			return
		}
		config, err := lib.ParseYAML(checksFileBytes)
		if err != nil {
			log.Error("Unable to parse checks.yml: " + err.Error())
			return
		}
		downCounter := make(map[string]int)
		failed, successful, total := lib.RunChecks(
			config.Checks,
			&config.ExportedReporters,
			downCounter)
		if failed >= successful {
			log.Warnf("Failed/Successful/Total: %d/%d/%d", failed, successful, total)
		} else {
			log.Infof("Failed/Successful/Total: %d/%d/%d", failed, successful, total)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
