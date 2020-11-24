package config

// Config is the configuration for the installation template.
type BaseConfig struct {
	ServerImage        string
	ImagePullPolicy    string
	AdvertiseInternal  bool
	ServiceAnnotations map[string]string
	ProviderConfig     struct{}
	Namespace          string

	// K8s-specific
	ServerName      string
	ServiceName     string
	OpenShift       bool
	Replicas        int32
	CPULimit        string
	MemLimit        string
	CPURequest      string
	MemRequest      string
	StorageRequest  string
	SecretFile      string
	ImagePullSecret string

	// Nomad-specific
	RegionF         string
	DatacentersF    []string
	PolicyOverrideF bool
}
