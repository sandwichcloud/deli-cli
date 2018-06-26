package policy

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
)

type setCommand struct {
	cmd.Command
	policyFile *string
	project    *string
	raw        *bool
}

func (c *setCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("set", "Set a policy").Action(c.action)
	c.policyFile = command.Arg("policy file", "The policy json document").Required().ExistingFile()
}

func (c *setCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	policy := api.Policy{}

	policyBytes, err := ioutil.ReadFile(*c.policyFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(policyBytes, &policy)
	if err != nil {
		return err
	}

	if c.project != nil {
		err = c.Application.APIClient.ProjectPolicy(*c.project).Set(policy)
	} else {
		err = c.Application.APIClient.SystemPolicy().Set(policy)
	}
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			fmt.Println("{}")
		} else {
			logrus.Info("Policy has been set")
		}
	}
	return nil
}
