module github.com/hashicorp/waypoint

go 1.13

require (
	github.com/Azure/azure-sdk-for-go v42.3.0+incompatible
	github.com/Azure/go-autorest/autorest v0.10.2 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.8.3 // indirect
	github.com/Azure/go-autorest/autorest/azure/auth v0.4.2
	github.com/Azure/go-autorest/autorest/to v0.3.0
	github.com/Azure/go-autorest/autorest/validation v0.2.0 // indirect
	github.com/adrg/xdg v0.2.1
	github.com/armon/circbuf v0.0.0-20190214190532-5111143e8da2
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/aws/aws-sdk-go v1.33.6
	github.com/bgentry/speakeasy v0.1.0
	github.com/boltdb/bolt v1.3.1
	github.com/briandowns/spinner v1.8.0
	github.com/buildpacks/pack v0.11.1
	github.com/creack/pty v1.1.11
	github.com/davecgh/go-spew v1.1.1
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/docker/cli v0.0.0-20200312141509-ef2f64abbd37
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/docker v1.4.2-0.20200221181110-62bd5a33f707
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/docker/libtrust v0.0.0-20160708172513-aabc10ec26b7 // indirect
	github.com/dustin/go-humanize v1.0.0
	github.com/elazarl/go-bindata-assetfs v1.0.1
	github.com/fatih/color v1.9.0
	github.com/flynn/noise v0.0.0-20180327030543-2492fe189ae6
	github.com/gliderlabs/ssh v0.2.2
	github.com/go-git/go-git/v5 v5.1.0
	github.com/go-openapi/runtime v0.19.15
	github.com/go-openapi/strfmt v0.19.5
	github.com/go-ozzo/ozzo-validation/v4 v4.2.1
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/golang/protobuf v1.4.2
	github.com/google/renameio v0.1.0
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/gorilla/handlers v1.4.2
	github.com/hashicorp/go-argmapper v0.0.0-20200617185818-77662eabd58a
	github.com/hashicorp/go-getter v1.4.1
	github.com/hashicorp/go-hclog v0.14.1
	github.com/hashicorp/go-memdb v1.2.0
	github.com/hashicorp/go-multierror v1.1.0
	github.com/hashicorp/go-plugin v1.1.0
	github.com/hashicorp/go-version v1.2.0
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/hcl/v2 v2.3.0
	github.com/hashicorp/horizon v0.0.0-20200730033532-a0cec81f4e7a
	github.com/hashicorp/nomad/api v0.0.0-20200814140818-42de70466a9d
	github.com/hashicorp/securetunnel v0.0.0-20200213234122-704adcadd8b2
	github.com/hashicorp/waypoint-hzn v0.0.0-20200716172850-719df2bafd65
	github.com/hashicorp/yamux v0.0.0-20200609203250-aecfd211c9ce // indirect
	github.com/imdario/mergo v0.3.9
	github.com/improbable-eng/grpc-web v0.13.0
	github.com/kr/text v0.2.0
	github.com/lab47/vterm v0.0.0-20200730001001-14f91f37349f
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mattn/go-colorable v0.1.7
	github.com/mattn/go-isatty v0.0.12
	github.com/mitchellh/cli v1.1.0
	github.com/mitchellh/go-grpc-net-conn v0.0.0-20200407005438-c00174eff6c8
	github.com/mitchellh/go-linereader v0.0.0-20190213213312-1b945b3263eb
	github.com/mitchellh/go-testing-interface v1.14.0
	github.com/mitchellh/mapstructure v1.3.3
	github.com/mitchellh/pointerstructure v1.0.0
	github.com/mitchellh/protostructure v0.0.0-20200814180458-3cfccdb015ce
	github.com/morikuni/aec v1.0.0
	github.com/mr-tron/base58 v1.2.0
	github.com/netlify/open-api v0.15.0
	github.com/oklog/run v1.1.0
	github.com/oklog/ulid v1.3.1
	github.com/olekukonko/tablewriter v0.0.4
	github.com/pkg/errors v0.9.1
	github.com/posener/complete v1.2.3
	github.com/rs/cors v1.7.0 // indirect
	github.com/sebdah/goldie v1.0.0
	github.com/skratchdot/open-golang v0.0.0-20200116055534-eef842397966
	github.com/slack-go/slack v0.6.5
	github.com/stretchr/testify v1.5.1
	github.com/zclconf/go-cty v1.5.1
	github.com/zclconf/go-cty-yaml v1.0.2
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899
	google.golang.org/api v0.20.0
	google.golang.org/genproto v0.0.0-20200715011427-11fb19a81f2c
	google.golang.org/grpc v1.30.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	k8s.io/api v0.18.0
	k8s.io/apimachinery v0.18.0
	k8s.io/client-go v0.18.0
)

// NOTE(mitchellh): I'm keeping these commented and in here because during
// development at the moment it is common to be working on these libs too.
// replace github.com/hashicorp/go-argmapper => ../go-argmapper
// replace github.com/hashicorp/horizon => ../horizon
