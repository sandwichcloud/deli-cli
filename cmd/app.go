package cmd

import (
	"fmt"
	"os"

	"errors"

	"io/ioutil"
	"os/user"
	"path"

	"encoding/json"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api/client"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type Application struct {
	CLIApp    *kingpin.Application
	Debug     *bool
	APIClient client.ClientInterface
	AuthToken *oauth2.Token
}

type Command struct {
	Application *Application
}

func (app *Application) Setup() {
	app.CLIApp = kingpin.New("deli", "Sandwich Cloud CLI").PreAction(app.setupLogging)

	app.Debug = app.CLIApp.Flag("debug", "Debug logging.").Short('d').Bool()
	apiServer := app.CLIApp.Flag("api-server", "Sandwich Cloud API Server [Env: DELI_API_SERVER]").Default("http://localhost:8080").Envar("DELI_API_SERVER").String()

	apiClient := &client.SandwichClient{APIServer: apiServer}
	app.APIClient = apiClient
}

func (app *Application) setupLogging(_ *kingpin.ParseElement, _ *kingpin.ParseContext) error {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetLevel(log.InfoLevel)

	if *app.Debug {
		log.SetLevel(log.DebugLevel)
	}

	return nil
}

func (app *Application) Run() {
	_, err := app.CLIApp.Parse(os.Args[1:])
	if err != nil {
		message := err.Error()
		if message[0] == '{' { //This is a json error message so just print that.
			fmt.Println(message)
			os.Exit(2)
		} else {
			log.SetOutput(os.Stderr) //Switch to stderr for readable errors
			log.Fatal(message)
		}
	}
}

func (app *Application) LoadCreds() error {
	u, err := user.Current()
	if err != nil {
		return errors.New("Cannot find the current system user.")
	}

	creds_file := path.Join(u.HomeDir, ".sandwich", "credentials")
	creds_data, err := ioutil.ReadFile(creds_file)
	if err != nil {
		return errors.New("Cannot read existing sandwich cloud credentials. Have you logged in?")
	}

	token := &oauth2.Token{}
	err = json.Unmarshal(creds_data, token)
	if err != nil {
		return errors.New("Existing sandwich cloud credentials have corrupted, please login again.")
	}

	app.AuthToken = token
	app.APIClient.SetToken(app.AuthToken)
	return nil
}

func (app *Application) SaveCreds() error {
	u, err := user.Current()
	if err != nil {
		return errors.New("Cannot find the current system user.")
	}

	config_dir := path.Join(u.HomeDir, ".sandwich")
	if _, err := os.Stat(config_dir); os.IsNotExist(err) {
		os.Mkdir(config_dir, os.ModePerm)
	}

	token_json, _ := json.Marshal(app.AuthToken)
	err = ioutil.WriteFile(path.Join(config_dir, "credentials"), token_json, 0644)
	if err != nil {
		return errors.New("Failed to write token data to file!")
	}

	return nil
}
