module github.com/hashicorp/waypoint

go 1.13

require (
	cloud.google.com/go v0.45.1 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/anmitsu/go-shlex v0.0.0-20161002113705-648efa622239 // indirect
	github.com/armon/circbuf v0.0.0-20190214190532-5111143e8da2
	github.com/aws/aws-sdk-go v1.28.12
	github.com/boltdb/bolt v1.3.1
	github.com/briandowns/spinner v1.8.0
	github.com/containerd/containerd v1.3.2 // indirect
	github.com/creack/pty v1.1.9
	github.com/davecgh/go-spew v1.1.1
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.4.2-0.20190924003213-a8608b5b67c7
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/dustin/go-humanize v1.0.0
	github.com/fatih/color v1.9.0
	github.com/flynn/go-shlex v0.0.0-20150515145356-3f9db97f8568 // indirect
	github.com/flynn/noise v0.0.0-20180327030543-2492fe189ae6
	github.com/gliderlabs/ssh v0.2.2
	github.com/golang/protobuf v1.3.4
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/hashicorp/go-argmapper v0.0.0-20200601005345-c13e5b41aa1f
	github.com/hashicorp/go-hclog v0.14.0
	github.com/hashicorp/go-memdb v1.2.0
	github.com/hashicorp/go-multierror v1.1.0
	github.com/hashicorp/go-plugin v1.1.0
	github.com/hashicorp/go-version v1.2.0
	github.com/hashicorp/hcl/v2 v2.3.0
	github.com/hashicorp/securetunnel v0.0.0-20200213234122-704adcadd8b2
	github.com/hashicorp/yamux v0.0.0-20190923154419-df201c70410d // indirect
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/kr/text v0.1.0
	github.com/mattn/go-colorable v0.1.6
	github.com/mattn/go-isatty v0.0.12
	github.com/mitchellh/cli v1.0.0
	github.com/mitchellh/go-grpc-net-conn v0.0.0-20200407005438-c00174eff6c8
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/go-linereader v0.0.0-20190213213312-1b945b3263eb
	github.com/mitchellh/go-testing-interface v1.14.0
	github.com/mitchellh/go-wordwrap v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.1.2
	github.com/mitchellh/protostructure v0.0.0-20200302233719-00c1118a7e52
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/oklog/run v1.1.0
	github.com/oklog/ulid v1.3.1
	github.com/olekukonko/tablewriter v0.0.4
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/posener/complete v1.1.1
	github.com/sebdah/goldie v1.0.0
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/stretchr/testify v1.5.1
	go.etcd.io/bbolt v1.3.3 // indirect
	go.opencensus.io v0.22.2 // indirect
	golang.org/x/crypto v0.0.0-20200220183623-bac4c82f6975
	golang.org/x/net v0.0.0-20200219183655-46282727080f // indirect
	golang.org/x/sys v0.0.0-20200523222454-059865788121 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	google.golang.org/api v0.20.0
	google.golang.org/genproto v0.0.0-20200218151345-dad8c97a84f5
	google.golang.org/grpc v1.28.1
	gotest.tools v2.2.0+incompatible // indirect
	k8s.io/api v0.18.0
	k8s.io/apimachinery v0.18.0
	k8s.io/client-go v0.18.0
)

// NOTE(mitchellh): I'm keeping these commented and in here because during
// development at the moment it is common to be working on these libs too.
// replace github.com/hashicorp/go-argmapper => ../go-argmapper
