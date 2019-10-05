package cmd

import (
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/KixPanganiban/bantay/lib"
	"github.com/KixPanganiban/bantay/log"
	"github.com/spf13/cobra"
)

// serverCmd spawns a server to run checks periodically.
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Spawn a server and run checks periodically.",
	Long:  `Spawn a server and run checks periodically, with settings defined in the YAML config under server.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Infoln("Server started.")
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
		for true {
			log.Debugln("Running checks...")
			failed, successful, total := lib.RunChecks(
				config.Checks,
				&config.ExportedReporters,
				downCounter)
			if failed >= successful {
				log.Warnf("Failed/Successful/Total: %d/%d/%d", failed, successful, total)
			} else {
				log.Infof("Failed/Successful/Total: %d/%d/%d", failed, successful, total)
			}
			log.Debugln("Sleeping for 10 seconds.")
			time.Sleep(10 * time.Second)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
