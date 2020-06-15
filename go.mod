module github.com/hashicorp/waypoint

go 1.13

require (
	github.com/armon/circbuf v0.0.0-20190214190532-5111143e8da2
	github.com/aws/aws-sdk-go v1.28.12
	github.com/boltdb/bolt v1.3.1
	github.com/briandowns/spinner v1.8.0
	github.com/buildpacks/pack v0.11.1
	github.com/creack/pty v1.1.9
	github.com/davecgh/go-spew v1.1.1
	github.com/docker/cli v0.0.0-20200312141509-ef2f64abbd37
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/docker v1.4.2-0.20200221181110-62bd5a33f707
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/docker/libtrust v0.0.0-20160708172513-aabc10ec26b7 // indirect
	github.com/dustin/go-humanize v1.0.0
	github.com/fatih/color v1.9.0
	github.com/flynn/noise v0.0.0-20180327030543-2492fe189ae6
	github.com/gliderlabs/ssh v0.2.2
	github.com/go-openapi/runtime v0.19.15
	github.com/go-openapi/strfmt v0.19.5
	github.com/golang/protobuf v1.3.5
	github.com/hashicorp/go-argmapper v0.0.0-20200606213939-1b3495a979bb
	github.com/hashicorp/go-hclog v0.14.0
	github.com/hashicorp/go-memdb v1.2.0
	github.com/hashicorp/go-multierror v1.1.0
	github.com/hashicorp/go-plugin v1.1.0
	github.com/hashicorp/go-version v1.2.0
	github.com/hashicorp/hcl/v2 v2.3.0
	github.com/hashicorp/horizon v0.0.0-20200529210319-3ee3f485d2ca
	github.com/hashicorp/securetunnel v0.0.0-20200213234122-704adcadd8b2
	github.com/imdario/mergo v0.3.8
	github.com/kr/text v0.2.0
	github.com/mattn/go-colorable v0.1.6
	github.com/mattn/go-isatty v0.0.12
	github.com/mitchellh/cli v1.0.0
	github.com/mitchellh/go-grpc-net-conn v0.0.0-20200407005438-c00174eff6c8
	github.com/mitchellh/go-linereader v0.0.0-20190213213312-1b945b3263eb
	github.com/mitchellh/go-testing-interface v1.14.0
	github.com/mitchellh/mapstructure v1.1.2
	github.com/mitchellh/protostructure v0.0.0-20200302233719-00c1118a7e52
	github.com/mr-tron/base58 v1.2.0
	github.com/netlify/open-api v0.15.0
	github.com/oklog/run v1.1.0
	github.com/oklog/ulid v1.3.1
	github.com/olekukonko/tablewriter v0.0.4
	github.com/pkg/errors v0.9.1
	github.com/posener/complete v1.1.1
	github.com/sebdah/goldie v1.0.0
	github.com/skratchdot/open-golang v0.0.0-20200116055534-eef842397966
	github.com/stretchr/testify v1.5.1
	golang.org/x/crypto v0.0.0-20200317142112-1b76d66859c6
	golang.org/x/sys v0.0.0-20200602225109-6fdc65e7d980 // indirect
	google.golang.org/api v0.20.0
	google.golang.org/genproto v0.0.0-20200416231807-8751e049a2a0
	google.golang.org/grpc v1.28.1
	k8s.io/api v0.18.0
	k8s.io/apimachinery v0.18.0
	k8s.io/client-go v0.18.0
)

// NOTE(mitchellh): I'm keeping these commented and in here because during
// development at the moment it is common to be working on these libs too.
// replace github.com/hashicorp/go-argmapper => ../go-argmapper
