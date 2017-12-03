package main

import (
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sandwichcloud/deli-cli/cmd/auth"
	"github.com/sandwichcloud/deli-cli/cmd/builtin"
	"github.com/sandwichcloud/deli-cli/cmd/image"
	"github.com/sandwichcloud/deli-cli/cmd/instance"
	"github.com/sandwichcloud/deli-cli/cmd/network"
	"github.com/sandwichcloud/deli-cli/cmd/policy"
	"github.com/sandwichcloud/deli-cli/cmd/project"
	"github.com/sandwichcloud/deli-cli/cmd/region"
	"github.com/sandwichcloud/deli-cli/cmd/role"
	"github.com/sandwichcloud/deli-cli/cmd/serviceAccount"
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

	imageCommand := image.Command{}
	imageCommand.Register(app)

	networkCommand := network.Command{}
	networkCommand.Register(app)

	instanceCommand := instance.Command{}
	instanceCommand.Register(app)

	policyCommand := policy.Command{}
	policyCommand.Register(app)

	roleCommand := role.Command{}
	roleCommand.Register(app)

	serviceAccountCommand := serviceAccount.Command{}
	serviceAccountCommand.Register(app)

	builtinCommand := builtin.Command{}
	builtinCommand.Register(app)

	app.Run()
}
