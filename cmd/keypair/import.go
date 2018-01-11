package keypair

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/alecthomas/kingpin"
	"github.com/fatih/structs"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
)

type importCommand struct {
	cmd.Command
	raw       *bool
	name      *string
	publicKey *string
}

func (c *importCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("import", "Import a keypair").Action(c.action)
	c.name = command.Arg("name", "The keypair name").Required().String()
	c.publicKey = command.Arg("public key file", "The public key file for the keypair").Required().ExistingFile()
}

func (c *importCommand) action(app *kingpin.Application, element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.SetScopedToken()
	if err != nil {
		return err
	}

	publicKeyBytes, err := ioutil.ReadFile(*c.publicKey)
	if err != nil {
		return err
	}
	publicKey := string(publicKeyBytes)

	keypair, err := c.Application.APIClient.Keypair().Create(*c.name, publicKey)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			keyPairMap := structs.Map(keypair)
			keypairBytes, _ := json.MarshalIndent(keyPairMap, "", "  ")
			fmt.Println(string(keypairBytes))
		} else {
			logrus.Infof("Keypair '%s' imported with an ID of '%s'", keypair.Name, keypair.ID)
		}
	}
	return nil
}
