package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log <...ID>",
	Short: "Download log of executed Task by ID",
	Long: `
	Usage :

	icli log ID

	will try to download logfile from Infra Server with matching ID
	you may provide only a part of the ID (will find the first match),
	and it works on multiple IDs 
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			if err := cmd.Usage(); err != nil {
				log.Fatalf("Error printing Usage %v\n", err)
			}
			os.Exit(1)
		}
		downloadLogs(args)
	},
}

func init() {
	RootCmd.AddCommand(logCmd)
}

func downloadLogs(args []string) {
	for _, id := range args {
		task, err := queryTask(id)
		if err != nil {
			fmt.Printf("Task not found %v %v\n", id, err)
			continue
		}
		if !strings.Contains(task.Status, "log available") {
			fmt.Printf("Log not available for Task %v status : %v\n", task.ID, task.Status)
			continue
		}
		if _, err := os.Stat("logs/" + task.ID + ".log"); !os.IsNotExist(err) {
			fmt.Printf("Logfile logs/%v.log already exist\n", task.ID)
			continue
		}
		if err := downloadLog(task.ID); err != nil {
			fmt.Printf("Error downloading logfile logs/%v.log %v\n", task.ID, err)
			continue
		}
		fmt.Printf("Successfully downloaded logfile logs/%v.log\n", task.ID)
	}
}

func downloadLog(ID string) error {
	filepath := path.Join("logs/", ID+".log")
	outFile, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.Wrapf(err, "Error opening logfile %v\n", filepath)
	}
	defer func() {
		if closeErr := outFile.Close(); closeErr != nil {
			log.Fatalf("Error closing logfile %v\n%v\n", filepath, closeErr)
		}
	}()

	req, err := http.NewRequest("GET", viper.GetString("apiURL")+"/logs/"+ID+".log", nil)
	if err != nil {
		return errors.Wrapf(err, "Cannot download logfile %v from server\n", filepath)
	}
	req.Header.Set("APITOKEN", viper.GetString("apiToken"))
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "Error requesting logfile %v\n", filepath)
	}
	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil {
			log.Fatalf("Error closing response body %v\n", closeErr)
		}
	}()

	_, err = io.Copy(outFile, res.Body)
	if err != nil {
		return errors.Wrapf(err, "Error writinng logfile %v\n", filepath)
	}

	return nil

}
