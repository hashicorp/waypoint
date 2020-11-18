package serverinstall

// Config is the configuration for the Kubernetes installation template.
type Config struct {
	ServerImage        string
	ImagePullPolicy    string
	AdvertiseInternal  bool
	ServiceAnnotations map[string]string

	// K8s config
	Namespace       string
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

	// NomadConfig
	RegionF         string
	DatacentersF    []string
	NamespaceF      string
	PolicyOverrideF bool
}

