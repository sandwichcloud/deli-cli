package keypair

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"errors"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
)

type createCommand struct {
	cmd.Command
	raw       *bool
	name      *string
	publicKey *string
}

func (c *createCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("create", "Create a keypair").Action(c.action)
	c.name = command.Flag("name", "The image name").Required().String()
	c.publicKey = command.Arg("public key", "The public key for the keypair").Required().String()
}

func (c *createCommand) action(app *kingpin.Application, element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.SetScopedToken()
	if err != nil {
		return err
	}

	publicKey := *c.publicKey
	if strings.HasPrefix(publicKey, "@") {
		publicKeyBytes, err := ioutil.ReadFile(publicKey[1:])
		if err != nil {
			return err
		}
		publicKey = string(publicKeyBytes)
	}

	keypair, err := c.Application.APIClient.Keypair().Create(*c.name, publicKey)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			keypairBytes, _ := json.MarshalIndent(keypair, "", "  ")
			fmt.Println(string(keypairBytes))
		} else {
			logrus.Infof("Keypair '%s' created with an ID of '%s'", keypair.Name, keypair.ID)
		}
	}
	return nil
}
