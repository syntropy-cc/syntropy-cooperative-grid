package usb

// USBDevice representa um dispositivo USB detectado
type USBDevice struct {
	Path        string `json:"path"`
	Size        string `json:"size"`
	SizeGB      int    `json:"size_gb"`
	Model       string `json:"model"`
	Vendor      string `json:"vendor"`
	Serial      string `json:"serial"`
	Removable   bool   `json:"removable"`
	Platform    string `json:"platform"`
	DiskNumber  int    `json:"disk_number,omitempty"`
	WindowsPath string `json:"windows_path,omitempty"`
}

// Config representa a configuração para criação de USB
type Config struct {
	NodeName        string `json:"node_name"`
	NodeDescription string `json:"node_description"`
	Coordinates     string `json:"coordinates"`
	OwnerKeyFile    string `json:"owner_key_file"`
	Label           string `json:"label"`
	ISOPath         string `json:"iso_path"`
	DiscoveryServer string `json:"discovery_server"`
	SSHPublicKey    string `json:"ssh_public_key"`
	SSHPrivateKey   string `json:"ssh_private_key"`
	CreatedBy       string `json:"created_by"`
}

// CloudInitConfig representa a configuração do cloud-init
type CloudInitConfig struct {
	NodeName         string
	NodeDescription  string
	Coordinates      string
	OwnerKey         string
	DiscoveryServer  string
	SSHPublicKey     string
	CreatedBy        string
	CreatedAt        string
	InstanceID       string
	Interface        string
	Gateway          string
	NodeIPSuffix     string
	PrimaryInterface string
	MeshGateway      string
	MgmtGateway      string
	HTTPProxy        string
	HTTPSProxy       string
	NodeType         string
	HardwareType     string
	CPUCores         int
	MemoryGB         int
	StorageGB        int
	InitialRole      string
	CanBeLeader      bool
	CanBeWorker      bool
	NodeCertPath     string
	NodeKeyPath      string
	CACertPath       string
	// Certificados PEM para write_files
	CACertPEM   string
	NodeCertPEM string
	NodeKeyPEM  string
}

// Certificates representa os certificados TLS gerados
type Certificates struct {
	CAKey    []byte
	CACert   []byte
	NodeKey  []byte
	NodeCert []byte
}

// WindowsDisk estrutura para parse do JSON do PowerShell
type WindowsDisk struct {
	Number       int    `json:"Number"`
	FriendlyName string `json:"FriendlyName"`
	Size         int64  `json:"Size"`
	SerialNumber string `json:"SerialNumber"`
	BusType      string `json:"BusType"`
	Model        string `json:"Model"`
}
