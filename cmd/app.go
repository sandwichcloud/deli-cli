package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"time"

	"errors"

	"io/ioutil"
	"os/user"
	"path"

	"encoding/json"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/api/client"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context/ctxhttp"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

type AuthTokens struct {
	Unscoped *oauth2.Token `json:"unscoped"`
	Scoped   *oauth2.Token `json:"scoped,omitempty"`
}

type Application struct {
	CLIApp     *kingpin.Application
	Debug      *bool
	APIClient  client.ClientInterface
	AuthTokens *AuthTokens
}

type Command struct {
	Application *Application
}

func (app *Application) Setup() {
	app.CLIApp = kingpin.New("deli", "Sandwich Cloud CLI")

	app.Debug = app.CLIApp.Flag("debug", "Debug logging.").Short('d').PreAction(app.setupLogging).Bool()
	apiServer := app.CLIApp.Flag("api-server", "Sandwich Cloud API Server [Env: DELI_API_SERVER]").Default("http://localhost:8080").Envar("DELI_API_SERVER").String()

	apiClient := &client.SandwichClient{APIServer: apiServer}
	app.APIClient = apiClient
}

func (app *Application) setupLogging(_ *kingpin.Application, _ *kingpin.ParseElement, _ *kingpin.ParseContext) error {
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
		if os.IsNotExist(err) {
			return app.loadCredsMetadata()
		}
		return errors.New("Cannot read existing sandwich cloud credentials. Have you logged in?")
	}

	tokens := &AuthTokens{}
	err = json.Unmarshal(creds_data, tokens)
	if err != nil {
		return errors.New("Existing sandwich cloud credentials have corrupted, please login again.")
	}

	if tokens.Unscoped == nil {
		return errors.New("No token found. Have you logged in?")
	}

	app.AuthTokens = tokens
	return nil
}

func (app *Application) loadCredsMetadata() error {
	cloudFile := path.Join(string(os.PathSeparator), "etc", "cloud", "cloud.cfg")
	cloudConfigFileData, err := ioutil.ReadFile(cloudFile)
	if err != nil {
		return errors.New("Could not read cloud-init config to find metadata service")
	}

	type cloudConfig struct {
		Datasource struct {
			NoCloud struct {
				SeedFrom string `json:"seedfrom"`
			} `json:"NoCloud"`
		} `json:"datasource"`
	}

	cloudConfigData := &cloudConfig{}
	yaml.Unmarshal(cloudConfigFileData, cloudConfigData)
	metadataURL := cloudConfigData.Datasource.NoCloud.SeedFrom
	if metadataURL == "" {
		return errors.New("Could not find metadata service.")
	}

	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	response, err := ctxhttp.Get(ctx, http.DefaultClient, metadataURL+"/security/creds")
	if err != nil {
		if err == context.DeadlineExceeded {
			return api.ErrTimedOut
		}
		return err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		apiError, err := api.ParseErrors(response.StatusCode, responseData)
		if err != nil {
			return err
		}
		return apiError
	}

	type secCreds struct {
		Token string `json:"token"`
	}

	serviceAccountToken := &secCreds{}
	json.Unmarshal(responseData, serviceAccountToken)
	oauthToken := &oauth2.Token{
		AccessToken: serviceAccountToken.Token,
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(10 * time.Minute),
	}
	tokens := &AuthTokens{
		Unscoped: oauthToken,
		Scoped:   oauthToken,
	}
	app.AuthTokens = tokens

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

	token_json, _ := json.Marshal(app.AuthTokens)
	err = ioutil.WriteFile(path.Join(config_dir, "credentials"), token_json, 0644)
	if err != nil {
		return errors.New("Failed to write token data to file!")
	}

	return nil
}

func (app *Application) isExpired(token *oauth2.Token) bool {
	return token.Expiry.Add(-10 * time.Second).Before(time.Now())
}

func (app *Application) SetUnScopedToken() error {

	if app.isExpired(app.AuthTokens.Unscoped) {
		return errors.New("Token is expired, please login to get a new one.")
	}

	app.APIClient.SetToken(app.AuthTokens.Unscoped)
	return nil
}

func (app *Application) SetScopedToken() error {

	if app.AuthTokens.Scoped == nil {
		return errors.New("No scoped token found, please scope your token to a project.")
	}

	if app.isExpired(app.AuthTokens.Scoped) {
		return errors.New("Scoped Token is expired, please re-scope to get a new one.")
	}

	app.APIClient.SetToken(app.AuthTokens.Scoped)
	return nil
}