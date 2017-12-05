package keypair

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"strings"

	"errors"

	"github.com/alecthomas/kingpin"
	"github.com/fatih/structs"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
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
	c.publicKey = command.Arg("public key", "The public key for the keypair. If not given one will be generated").Default("").String()
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
	privateKeyString := ""
	if publicKey == "" {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2046)
		if err != nil {
			return err
		}
		privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
		privateKeyString = string(pem.EncodeToMemory(privateKeyPEM))
		pubKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
		publicKey = string(ssh.MarshalAuthorizedKey(pubKey))
	} else if strings.HasPrefix(publicKey, "@") {
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
			keyPairMap := structs.Map(keypair)
			if privateKeyString != "" {
				keyPairMap["private_key"] = privateKeyString
			}
			keypairBytes, _ := json.MarshalIndent(keyPairMap, "", "  ")
			fmt.Println(string(keypairBytes))
		} else {
			logrus.Infof("Keypair '%s' created with an ID of '%s'", keypair.Name, keypair.ID)
			if privateKeyString != "" {
				logrus.Info("Private Key generated bellow:")
				fmt.Println(privateKeyString)
			}
		}
	}
	return nil
}
