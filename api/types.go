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
	Default *string           `json:"default"`
	Github  *GithubAuthDriver `json:"github"`
}

type GithubAuthDriver struct {
	//Github auth driver has no options
}

type Task struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	State        string    `json:"state"`
	ErrorMessage string    `json:"error_message"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	StoppedAt    time.Time `json:"stopped_at"`
}

type Project struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	State         string    `json:"state"`
	CurrentTaskID uuid.UUID `json:"current_task_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ProjectList struct {
	Links    []pageLinks `json:"projects_links"`
	Projects []Project   `json:"projects"`
}

type Region struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Datacenter     string    `json:"datacenter"`
	ImageDatastore string    `json:"image_datastore"`
	ImageFolder    string    `json:"image_folder"`
	Schedulable    bool      `json:"schedulable"`
	State          string    `json:"state"`
	CurrentTaskID  uuid.UUID `json:"current_task_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
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
	CurrentTaskID        uuid.UUID `json:"current_task_id"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type ZoneList struct {
	Links []pageLinks `json:"zones_links"`
	Zones []Zone      `json:"zones"`
}

type Image struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	FileName      string    `json:"file_name"`
	Visibility    string    `json:"visibility"`
	Locked        bool      `json:"locked"`
	RegionID      uuid.UUID `json:"region_id"`
	State         string    `json:"state"`
	CurrentTaskID uuid.UUID `json:"current_task_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ImageList struct {
	Links  []pageLinks `json:"images_links"`
	Images []Image     `json:"images"`
}

type TokenInfo struct {
	UserID             uuid.UUID `json:"user_id,omitempty"`
	Username           string    `json:"username,omitempty"`
	Driver             string    `json:"driver,omitempty"`
	ServiceAccountID   uuid.UUID `json:"service_account_id,omitempty"`
	ServiceAccountName string    `json:"service_account_name,omitempty"`
	ProjectID          string    `json:"project_id,omitempty"`
	GlobalRoles        []string  `json:"global_roles"`
	ProjectRoles       []string  `json:"project_roles,omitempty"`
}

type Network struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	PortGroup     string    `json:"port_group"`
	Cidr          string    `json:"cidr"`
	Gateway       net.IP    `json:"gateway"`
	DNSServers    []net.IP  `json:"dns_servers"`
	PoolStart     net.IP    `json:"pool_start"`
	PoolEnd       net.IP    `json:"pool_end"`
	RegionID      uuid.UUID `json:"region_id"`
	State         string    `json:"state"`
	CurrentTaskID uuid.UUID `json:"current_task_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type NetworkList struct {
	Links    []pageLinks `json:"networks_links"`
	Networks []Network   `json:"networks"`
}

type Instance struct {
	ID            uuid.UUID         `json:"id"`
	Name          string            `json:"name"`
	ImageID       uuid.UUID         `json:"image_id"`
	NetworkPortID uuid.UUID         `json:"network_port_id"`
	RegionID      uuid.UUID         `json:"region_id"`
	ZoneID        string            `json:"zone_id"`
	State         string            `json:"state"`
	CurrentTaskID uuid.UUID         `json:"current_task_id"`
	Tags          map[string]string `json:"tags"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

type InstanceList struct {
	Links     []pageLinks `json:"instances_links"`
	Instances []Instance  `json:"instances"`
}

type Policy struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PolicyList struct {
	Links    []pageLinks `json:"policies_links"`
	Policies []Policy    `json:"policies"`
}

type Role struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ProjectID   string    `json:"project_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RoleList struct {
	Links []pageLinks `json:"roles_links"`
	Roles []Role      `json:"roles"`
}

type ServiceAccount struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	RoleID    string    `json:"role_id"`
	ProjectID string    `json:"project_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ServiceAccountList struct {
	Links     []pageLinks      `json:"service-accounts_links"`
	Instances []ServiceAccount `json:"service-accounts"`
}
