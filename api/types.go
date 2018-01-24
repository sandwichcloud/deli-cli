package api

import (
	"time"

	"net"

	"github.com/satori/go.uuid"
)

type pageLinks struct {
	HREF string `json:"href"`
	REL  string `json:"rel"`
}

type AuthDiscover struct {
	Default  *string             `json:"default"`
	Github   *GithubAuthDriver   `json:"github"`
	Database *DatabaseAuthDriver `json:"database"`
}

type GithubAuthDriver struct {
	//Github auth driver has no options
}

type DatabaseAuthDriver struct {
	//Database auth driver has no options
}

type Project struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type ProjectList struct {
	Links    []pageLinks `json:"projects_links"`
	Projects []Project   `json:"projects"`
}

type ProjectQuota struct {
	VCPU     int `json:"vcpu"`
	Ram      int `json:"ram"`
	Disk     int `json:"disk"`
	UsedVCPU int `json:"used_vcpu"`
	UsedRam  int `json:"used_ram"`
	UsedDisk int `json:"used_disk"`
}

type ProjectMember struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Driver   string    `json:"driver"`
	Roles    []string  `json:"roles"`
}

type ProjectMemberList struct {
	Links          []pageLinks     `json:"project-members_links"`
	ProjectMembers []ProjectMember `json:"project-members"`
}

type Region struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Datacenter     string    `json:"datacenter"`
	ImageDatastore string    `json:"image_datastore"`
	ImageFolder    string    `json:"image_folder"`
	Schedulable    bool      `json:"schedulable"`
	State          string    `json:"state"`
	ErrorMessage   string    `json:"error_message"`
	CreatedAt      time.Time `json:"created_at"`
}

type RegionList struct {
	Links   []pageLinks `json:"regions_links"`
	Regions []Region    `json:"regions"`
}

type Zone struct {
	ID                   uuid.UUID `json:"id"`
	Name                 string    `json:"name"`
	RegionID             uuid.UUID `json:"region_id"`
	VMCluster            string    `json:"vm_cluster"`
	VMDatastore          string    `json:"vm_datastore"`
	VMFolder             string    `json:"vm_folder"`
	CoreProvisionPercent int       `json:"core_provision_percent"`
	RamProvisionPercent  int       `json:"ram_provision_percent"`
	Schedulable          bool      `json:"schedulable"`
	State                string    `json:"state"`
	ErrorMessage         string    `json:"error_message"`
	CreatedAt            time.Time `json:"created_at"`
}

type ZoneList struct {
	Links []pageLinks `json:"zones_links"`
	Zones []Zone      `json:"zones"`
}

type Volume struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	ZoneID       uuid.UUID `json:"zone_id"`
	Size         int       `json:"size"`
	AttachedTo   uuid.UUID `json:"attached_to"`
	Task         string    `json:"task"`
	State        string    `json:"state"`
	ErrorMessage string    `json:"error_message"`
	CreatedAt    time.Time `json:"created_at"`
}

type VolumeList struct {
	Links   []pageLinks `json:"volumes_links"`
	Volumes []Volume    `json:"volumes"`
}

type Image struct {
	ID           uuid.UUID `json:"id"`
	ProjectID    uuid.UUID `json:"project_id"`
	Name         string    `json:"name"`
	FileName     string    `json:"file_name"`
	Public       bool      `json:"public"`
	Locked       bool      `json:"locked"`
	RegionID     uuid.UUID `json:"region_id"`
	State        string    `json:"state"`
	ErrorMessage string    `json:"error_message"`
	CreatedAt    time.Time `json:"created_at"`
}

type ImageList struct {
	Links  []pageLinks `json:"images_links"`
	Images []Image     `json:"images"`
}

type ImageMemberList struct {
	Links   []pageLinks   `json:"members_links"`
	Members []ImageMember `json:"members"`
}

type ImageMember struct {
	ProjectID uuid.UUID `json:"project_id"`
}

type TokenInfo struct {
	Username           string    `json:"username,omitempty"`
	Driver             string    `json:"driver,omitempty"`
	ServiceAccountID   uuid.UUID `json:"service_account_id,omitempty"`
	ServiceAccountName string    `json:"service_account_name,omitempty"`
	ProjectID          string    `json:"project_id,omitempty"`
	GlobalRoles        []string  `json:"global_roles"`
	ProjectRoles       []string  `json:"project_roles,omitempty"`
}

type Network struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	PortGroup    string    `json:"port_group"`
	Cidr         string    `json:"cidr"`
	Gateway      net.IP    `json:"gateway"`
	DNSServers   []net.IP  `json:"dns_servers"`
	PoolStart    net.IP    `json:"pool_start"`
	PoolEnd      net.IP    `json:"pool_end"`
	RegionID     uuid.UUID `json:"region_id"`
	State        string    `json:"state"`
	ErrorMessage string    `json:"error_message"`
	CreatedAt    time.Time `json:"created_at"`
}

type NetworkList struct {
	Links    []pageLinks `json:"networks_links"`
	Networks []Network   `json:"networks"`
}

type Keypair struct {
	ID        uuid.UUID `json:"id" structs:"id"`
	Name      string    `json:"name" structs:"name"`
	PublicKey string    `json:"public_key" structs:"public_key"`
}

type KeypairList struct {
	Links    []pageLinks `json:"keypairs_links"`
	KeyPairs []Keypair   `json:"keypairs"`
}

type NetworkPort struct {
	ID        uuid.UUID `json:"id"`
	NetworkID uuid.UUID `json:"network_id"`
	IPAddress net.IP    `json:"ip_address"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
}

type NetworkPortList struct {
	Links        []pageLinks   `json:"network-ports_links"`
	NetworkPorts []NetworkPort `json:"network-ports"`
}

type Flavor struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	VCPUS     int       `json:"vcpus"`
	Ram       int       `json:"ram"`
	Disk      int       `json:"disk"`
	CreatedAt time.Time `json:"created_at"`
}

type FlavorList struct {
	Links   []pageLinks `json:"flavors_links"`
	Flavors []Flavor    `json:"flavors"`
}

type Instance struct {
	ID               uuid.UUID         `json:"id"`
	Name             string            `json:"name"`
	ImageID          uuid.UUID         `json:"image_id"`
	NetworkPortID    uuid.UUID         `json:"network_port_id"`
	RegionID         uuid.UUID         `json:"region_id"`
	ZoneID           uuid.UUID         `json:"zone_id"`
	ServiceAccountID uuid.UUID         `json:"service_account_id"`
	Tags             map[string]string `json:"tags"`
	KeypairIDs       []uuid.UUID       `json:"keypair_ids"`
	FlavorID         uuid.UUID         `json:"flavor_id"`
	VCPUS            int               `json:"vcpus"`
	Ram              int               `json:"ram"`
	Disk             int               `json:"disk"`
	State            string            `json:"state"`
	PowerState       string            `json:"power_state"`
	Task             string            `json:"task"`
	ErrorMessage     string            `json:"error_message"`
	CreatedAt        time.Time         `json:"created_at"`
}

type InstanceList struct {
	Links     []pageLinks `json:"instances_links"`
	Instances []Instance  `json:"instances"`
}

type Policy struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type PolicyList struct {
	Links    []pageLinks `json:"policies_links"`
	Policies []Policy    `json:"policies"`
}

type Role struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Policies  []string  `json:"policies"`
	CreatedAt time.Time `json:"created_at"`
}

type GlobalRoleList struct {
	Links []pageLinks `json:"global-roles_links"`
	Roles []Role      `json:"global-roles"`
}

type ProjectRoleList struct {
	Links []pageLinks `json:"project-roles_links"`
	Roles []Role      `json:"project-roles"`
}

type ServiceAccount struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	RoleID    string    `json:"role_id"`
	ProjectID string    `json:"project_id"`
	CreatedAt time.Time `json:"created_at"`
}

type ServiceAccountList struct {
	Links           []pageLinks      `json:"service-accounts_links"`
	ServiceAccounts []ServiceAccount `json:"service-accounts"`
}

type DatabaseUser struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Roles     []string  `json:"roles"`
	CreatedAt time.Time `json:"created_at"`
}

type DatabaseUserList struct {
	Links []pageLinks    `json:"users_links"`
	Users []DatabaseUser `json:"users"`
}
