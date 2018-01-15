package client

import (
	"net/http"

	"context"

	"net"

	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/api/client/auth"
	"github.com/sandwichcloud/deli-cli/api/client/database"
	"github.com/sandwichcloud/deli-cli/api/client/flavor"
	"github.com/sandwichcloud/deli-cli/api/client/image"
	"github.com/sandwichcloud/deli-cli/api/client/instance"
	"github.com/sandwichcloud/deli-cli/api/client/keypair"
	"github.com/sandwichcloud/deli-cli/api/client/network"
	"github.com/sandwichcloud/deli-cli/api/client/policy"
	"github.com/sandwichcloud/deli-cli/api/client/project"
	"github.com/sandwichcloud/deli-cli/api/client/region"
	"github.com/sandwichcloud/deli-cli/api/client/role"
	"github.com/sandwichcloud/deli-cli/api/client/serviceAccount"
	"github.com/sandwichcloud/deli-cli/api/client/zone"
	"golang.org/x/oauth2"
)

type SandwichClient struct {
	APIServer *string
	token     *oauth2.Token
}

type ClientInterface interface {
	createOAuthClient() *http.Client
	Auth() AuthClientInterface
	DatabaseAuth() DatabaseAuthClientInterface
	Project() ProjectClientInterface
	Region() RegionClientInterface
	Zone() ZoneClientInterface
	Image() ImageClientInterface
	Network() NetworkClientInterface
	NetworkPort() NetworkPortClientInterface
	Keypair() KeypairClientInterface
	Flavor() FlavorClientInterface
	Instance() InstanceClientInterface
	Policy() PolicyClientInterface
	Role() RoleClientInterface
	ServiceAccount() ServiceAccountClientInterface
	SetToken(token *oauth2.Token)
}

type AuthClientInterface interface {
	DiscoverAuth() (api.AuthDiscover, error)
	GithubLogin(options api.GithubAuthDriver, username, password, otpCode string) (*oauth2.Token, error)
	DatabaseLogin(options api.DatabaseAuthDriver, username, password string) (*oauth2.Token, error)
	ScopeToken(project *api.Project) (*oauth2.Token, error)
	TokenInfo() (*api.TokenInfo, error)
}

type DatabaseAuthClientInterface interface {
	Create(username, password string) (*api.DatabaseUser, error)
	Get(id string) (*api.DatabaseUser, error)
	Delete(id string) error
	List(limit int, marker string) (*api.DatabaseUserList, error)
	ChangePassword(id, password string) error
	AddRole(id, role string) error
	RemoveRole(id, role string) error
}

type ProjectClientInterface interface {
	Create(name string) (*api.Project, error)
	Get(id string) (*api.Project, error)
	Delete(id string) error
	List(all bool, limit int, marker string) (*api.ProjectList, error)
}

type RegionClientInterface interface {
	Create(name, datacenter, imageDatastore, imageFolder string) (*api.Region, error)
	Get(id string) (*api.Region, error)
	Delete(id string) error
	List(name string, limit int, marker string) (*api.RegionList, error)
	ActionSchedule(id string, schedulable bool) error
}

type ZoneClientInterface interface {
	Create(name, regionID, vmCluster, vmDatastore, vmFolder string, coreProvisionPercent, ramProvisionPercent int) (*api.Zone, error)
	Get(id string) (*api.Zone, error)
	Delete(id string) error
	List(regionID string, limit int, marker string) (*api.ZoneList, error)
	ActionSchedule(id string, schedulable bool) error
}

type ImageClientInterface interface {
	Create(name, regionID, fileName, visibility string) (*api.Image, error)
	Get(id string) (*api.Image, error)
	Delete(id string) error
	List(visibility string, limit int, marker string) (*api.ImageList, error)
	ActionLock(id string) error
	ActionUnlock(id string) error
	MemberAdd(id, projectID string) error
	MemberList(id string) (*api.ImageMemberList, error)
	MemberRemove(id, projectID string) error
}

type KeypairClientInterface interface {
	Create(name, publicKey string) (*api.Keypair, error)
	Get(id string) (*api.Keypair, error)
	Delete(id string) error
	List(limit int, marker string) (*api.KeypairList, error)
}

type NetworkPortClientInterface interface {
	Get(id string) (*api.NetworkPort, error)
	List(limit int, marker string) (*api.NetworkPortList, error)
	Delete(id string) error
}

type FlavorClientInterface interface {
	Create(name string, vcpus, ram, disk int) (*api.Flavor, error)
	Get(id string) (*api.Flavor, error)
	Delete(id string) error
	List(limit int, marker string) (*api.FlavorList, error)
}

type InstanceClientInterface interface {
	Create(name, imageID, regionID, zoneID, networkID, serviceAccountID string, flavorID string, disk int, keypairIDs []string, tags map[string]string) (*api.Instance, error)
	Get(id string) (*api.Instance, error)
	Delete(id string) error
	List(imageID string, limit int, marker string) (*api.InstanceList, error)
	ActionStop(id string, hard bool, timeout int) error
	ActionStart(id string) error
	ActionRestart(id string, hard bool, timeout int) error
	ActionImage(id string, name string, visibility string) (*api.Image, error)
	ActionResetState(id string, active bool) error
}

type NetworkClientInterface interface {
	Create(name, regionID, portGroup, cidr string, gateway, poolStart, poolEnd net.IP, dnsServers []net.IP) (*api.Network, error)
	Get(id string) (*api.Network, error)
	Delete(id string) error
	List(name, region_id string, limit int, marker string) (*api.NetworkList, error)
}

type PolicyClientInterface interface {
	Get(id string) (*api.Policy, error)
	List(limit int, marker string) (*api.PolicyList, error)
}

type RoleClientInterface interface {
	Create(name, roleType, description string) (*api.Role, error)
	Get(id string) (*api.Role, error)
	Delete(id string) error
	List(roleType string, limit int, marker string) (*api.RoleList, error)
}

type ServiceAccountClientInterface interface {
	Create(name, roleId string) (*api.ServiceAccount, error)
	Get(id string) (*api.ServiceAccount, error)
	Delete(id string) error
	List(limit int, marker string) (*api.ServiceAccountList, error)
}

func (client *SandwichClient) createOAuthClient() *http.Client {
	ctx := context.Background()
	return oauth2.NewClient(ctx, oauth2.StaticTokenSource(client.token))
}

func (client *SandwichClient) Auth() AuthClientInterface {
	authClient := &auth.AuthClient{APIServer: client.APIServer}

	if client.token != nil {
		authClient.HttpClient = client.createOAuthClient()
	}

	return authClient
}

func (client *SandwichClient) DatabaseAuth() DatabaseAuthClientInterface {
	return &database.DatabaseAuthClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient()}
}

func (client *SandwichClient) Project() ProjectClientInterface {
	return &project.ProjectClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient()}
}

func (client *SandwichClient) Region() RegionClientInterface {
	return &region.RegionClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient()}
}

func (client *SandwichClient) Zone() ZoneClientInterface {
	return &zone.ZoneClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient()}
}

func (client *SandwichClient) Image() ImageClientInterface {
	return &image.ImageClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient()}
}

func (client *SandwichClient) Network() NetworkClientInterface {
	return &network.NetworkClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient()}
}

func (client *SandwichClient) NetworkPort() NetworkPortClientInterface {
	return &network.NetworkPortClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient()}
}

func (client *SandwichClient) Keypair() KeypairClientInterface {
	return &keypair.KeypairClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient()}
}

func (client *SandwichClient) Flavor() FlavorClientInterface {
	return &flavor.FlavorClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient()}
}

func (client *SandwichClient) Instance() InstanceClientInterface {
	return &instance.InstanceClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient()}
}

func (client *SandwichClient) Policy() PolicyClientInterface {
	return &policy.PolicyClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient()}
}

func (client *SandwichClient) Role() RoleClientInterface {
	return &role.RoleClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient()}
}

func (client *SandwichClient) ServiceAccount() ServiceAccountClientInterface {
	return &serviceAccount.ServiceAccountClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient()}
}

func (client *SandwichClient) SetToken(token *oauth2.Token) {
	client.token = token
}
