package main

import (
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sandwichcloud/deli-cli/cmd/auth"
	"github.com/sandwichcloud/deli-cli/cmd/flavor"
	"github.com/sandwichcloud/deli-cli/cmd/image"
	"github.com/sandwichcloud/deli-cli/cmd/instance"
	"github.com/sandwichcloud/deli-cli/cmd/keypair"
	"github.com/sandwichcloud/deli-cli/cmd/metadata"
	"github.com/sandwichcloud/deli-cli/cmd/network"
	"github.com/sandwichcloud/deli-cli/cmd/permission"
	"github.com/sandwichcloud/deli-cli/cmd/policy"
	"github.com/sandwichcloud/deli-cli/cmd/project"
	"github.com/sandwichcloud/deli-cli/cmd/region"
	"github.com/sandwichcloud/deli-cli/cmd/role"
	"github.com/sandwichcloud/deli-cli/cmd/serviceAccount"
	"github.com/sandwichcloud/deli-cli/cmd/volume"
	"github.com/sandwichcloud/deli-cli/cmd/zone"
)

func main() {

	app := &cmd.Application{}
	app.Setup()

	authCommand := auth.Command{}
	authCommand.Register(app)

	projectCommand := project.Command{}
	projectCommand.Register(app)

	regionCommand := region.Command{}
	regionCommand.Register(app)

	zoneCommand := zone.Command{}
	zoneCommand.Register(app)

	flavorCommand := flavor.Command{}
	flavorCommand.Register(app)

	volumeCommand := volume.Command{}
	volumeCommand.Register(app)

	imageCommand := image.Command{}
	imageCommand.Register(app)

	networkCommand := network.Command{}
	networkCommand.Register(app)

	instanceCommand := instance.Command{}
	instanceCommand.Register(app)

	permissionCommand := permission.Command{}
	permissionCommand.Register(app)

	roleCommand := role.Command{}
	roleCommand.Register(app)

	policyCommand := policy.Command{}
	policyCommand.Register(app)

	serviceAccountCommand := serviceAccount.Command{}
	serviceAccountCommand.Register(app)

	keypairCommand := keypair.Command{}
	keypairCommand.Register(app)

	metaDataCommand := metadata.Command{}
	metaDataCommand.Register(app)

	app.Run()
}
