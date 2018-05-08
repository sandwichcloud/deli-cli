package client

import (
	"net/http"

	"context"

	"net"

	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/api/client/auth"
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
	"github.com/sandwichcloud/deli-cli/api/client/volume"
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
	Project() ProjectClientInterface
	Region() RegionClientInterface
	Zone() ZoneClientInterface
	Image() ImageClientInterface
	Network() NetworkClientInterface
	NetworkPort() NetworkPortClientInterface
	Keypair() KeypairClientInterface
	Flavor() FlavorClientInterface
	Volume() VolumeClientInterface
	Instance() InstanceClientInterface
	Policy() PolicyClientInterface
	GlobalRole() RoleClientInterface
	ProjectRole() RoleClientInterface
	GlobalServiceAccount() ServiceAccountClientInterface
	ProjectServiceAccount() ServiceAccountClientInterface
	SetToken(token *oauth2.Token)
}

type AuthClientInterface interface {
	DiscoverAuth() (api.AuthDiscover, error)
	GithubLogin(options api.GithubAuthDriver, username, password, otpCode string) (*oauth2.Token, error)
	DatabaseLogin(options api.DatabaseAuthDriver, username, password string) (*oauth2.Token, error)
	ScopeToken(project *api.Project) (*oauth2.Token, error)
	TokenInfo() (*api.TokenInfo, error)
}

type ProjectClientInterface interface {
	Create(name string) (*api.Project, error)
	Get(id string) (*api.Project, error)
	Delete(id string) error
	List(all bool, limit int, marker string) (*api.ProjectList, error)
	GetQuota() (*api.ProjectQuota, error)
	SetQuota(vcpu, ram, disk int) error
	AddMember(username, driver string) (*api.ProjectMember, error)
	GetMember(id string) (*api.ProjectMember, error)
	ListMembers(limit int, marker string) (*api.ProjectMemberList, error)
	UpdateMember(id string, roles []string) error
	RemoveMember(id string) error
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

type VolumeClientInterface interface {
	Create(name, zoneID string, size int) (*api.Volume, error)
	Get(id string) (*api.Volume, error)
	Delete(id string) error
	List(limit int, marker string) (*api.VolumeList, error)
	ActionAttach(id, instanceID string) error
	ActionDetach(id string) error
	ActionGrow(id string, newSize int) error
	ActionClone(id, name string) (*api.Volume, error)
}

type ImageClientInterface interface {
	Create(name, regionID, fileName string) (*api.Image, error)
	Get(id string) (*api.Image, error)
	Delete(id string) error
	List(visibility string, limit int, marker string) (*api.ImageList, error)
	ActionSetVisibility(id string, public bool) error
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
	Create(name, imageID, regionID, zoneID, networkID, serviceAccountID, flavorID string, disk int, keypairIDs []string, initialVolumes []api.InstanceInitialVolume, tags map[string]string, userData string) (*api.Instance, error)
	Get(id string) (*api.Instance, error)
	Delete(id string) error
	List(imageID string, limit int, marker string) (*api.InstanceList, error)
	ActionStop(id string, hard bool, timeout int) error
	ActionStart(id string) error
	ActionRestart(id string, hard bool, timeout int) error
	ActionImage(id string, name string) (*api.Image, error)
}

type NetworkClientInterface interface {
	Create(name, regionID, portGroup, cidr string, gateway, poolStart, poolEnd net.IP, dnsServers []net.IP) (*api.Network, error)
	Get(id string) (*api.Network, error)
	Delete(id string) error
	List(name, region_id string, limit int, marker string) (*api.NetworkList, error)
}

type PolicyClientInterface interface {
	Get(name string) (*api.Policy, error)
	List(limit int, marker string) (*api.PolicyList, error)
}

type RoleClientInterface interface {
	Create(name string, policies []string) (*api.Role, error)
	Get(id string) (*api.Role, error)
	GlobalList(limit int, marker string) (*api.GlobalRoleList, error)
	ProjectList(limit int, marker string) (*api.ProjectRoleList, error)
	Update(id string, policies []string) error
	Delete(id string) error
}

type ServiceAccountClientInterface interface {
	Create(name string) (*api.ServiceAccount, error)
	Get(id string) (*api.ServiceAccount, error)
	Delete(id string) error
	GlobalList(limit int, marker string) (*api.GlobalServiceAccountList, error)
	ProjectList(limit int, marker string) (*api.ProjectServiceAccountList, error)
	Update(id string, roles []string) error
	CreateKey(id, name string) (*oauth2.Token, error)
	DeleteKey(id, name string) error
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

func (client *SandwichClient) Project() ProjectClientInterface {
	return &project.ProjectClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient()}
}

func (client *SandwichClient) Region() RegionClientInterface {
	return &region.RegionClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient()}
}

func (client *SandwichClient) Zone() ZoneClientInterface {
	return &zone.ZoneClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient()}
}

func (client *SandwichClient) Volume() VolumeClientInterface {
	return &volume.VolumeClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient()}
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

func (client *SandwichClient) GlobalRole() RoleClientInterface {
	return &role.RoleClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient(), Type: "global-roles"}
}

func (client *SandwichClient) ProjectRole() RoleClientInterface {
	return &role.RoleClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient(), Type: "project-roles"}
}

func (client *SandwichClient) GlobalServiceAccount() ServiceAccountClientInterface {
	return &serviceAccount.ServiceAccountClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient(), Type: "global-service-accounts"}
}

func (client *SandwichClient) ProjectServiceAccount() ServiceAccountClientInterface {
	return &serviceAccount.ServiceAccountClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient(), Type: "project-service-accounts"}
}

func (client *SandwichClient) SetToken(token *oauth2.Token) {
	client.token = token
}
