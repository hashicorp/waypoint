module github.com/hashicorp/waypoint

go 1.13

require (
	github.com/Azure/azure-sdk-for-go v42.3.0+incompatible
	github.com/Azure/go-autorest/autorest v0.10.2
	github.com/Azure/go-autorest/autorest/azure/auth v0.4.2
	github.com/Azure/go-autorest/autorest/to v0.3.0
	github.com/adrg/xdg v0.2.1
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/armon/circbuf v0.0.0-20190214190532-5111143e8da2
	github.com/aws/aws-sdk-go v1.36.31
	github.com/bmatcuk/doublestar v1.1.5
	github.com/buildpacks/pack v0.18.1
	github.com/cenkalti/backoff/v4 v4.0.2
	github.com/certifi/gocertifi v0.0.0-20200922220541-2c3bb06c6054
	github.com/containerd/console v1.0.1
	github.com/creack/pty v1.1.11
	github.com/davecgh/go-spew v1.1.1
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/docker/cli v20.10.0-beta1.0.20201029214301-1d20b15adc38+incompatible
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/docker v20.10.7+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/docker/go-units v0.4.0
	github.com/docker/libtrust v0.0.0-20160708172513-aabc10ec26b7 // indirect
	github.com/dustin/go-humanize v1.0.0
	github.com/elazarl/go-bindata-assetfs v1.0.1
	github.com/fatih/color v1.12.0
	github.com/ghodss/yaml v1.0.0
	github.com/gliderlabs/ssh v0.3.1
	github.com/go-git/go-git/v5 v5.2.0
	github.com/go-ozzo/ozzo-validation/v4 v4.2.1
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/gofrs/flock v0.8.0
	github.com/golang/protobuf v1.5.2
	github.com/google/go-cmp v0.5.5
	github.com/google/go-jsonnet v0.17.0
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/google/uuid v1.1.2
	github.com/gorilla/handlers v1.4.2
	github.com/hashicorp/aws-sdk-go-base v0.7.0
	github.com/hashicorp/cap v0.1.1
	github.com/hashicorp/go-argmapper v0.2.0
	github.com/hashicorp/go-bexpr v0.1.7
	github.com/hashicorp/go-cleanhttp v0.5.2
	github.com/hashicorp/go-gcp-common v0.6.0
	github.com/hashicorp/go-getter v1.4.1
	github.com/hashicorp/go-hclog v0.16.1
	github.com/hashicorp/go-memdb v1.3.2
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/go-plugin v1.4.2
	github.com/hashicorp/go-uuid v1.0.2
	github.com/hashicorp/go-version v1.2.0
	github.com/hashicorp/hcl/v2 v2.10.1-0.20210621220818-327f3ce2570e
	github.com/hashicorp/horizon v0.0.0-20210317214650-d2053943be04
	github.com/hashicorp/nomad/api v0.0.0-20200814140818-42de70466a9d
	github.com/hashicorp/vault/api v1.0.5-0.20200519221902-385fac77e20f
	github.com/hashicorp/vault/sdk v0.1.14-0.20201202172114-ee5ebeb30fef
	github.com/hashicorp/waypoint-hzn v0.0.0-20201008221232-97cd4d9120b9
	github.com/hashicorp/waypoint-plugin-sdk v0.0.0-20210714195507-d64bfe0a3ffa
	github.com/imdario/mergo v0.3.11
	github.com/improbable-eng/grpc-web v0.13.0
	github.com/kevinburke/go-bindata v3.22.0+incompatible
	github.com/kr/text v0.2.0
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/mitchellh/cli v1.1.2
	github.com/mitchellh/copystructure v1.0.0
	github.com/mitchellh/go-glint v0.0.0-20201015034436-f80573c636de
	github.com/mitchellh/go-grpc-net-conn v0.0.0-20200427190222-eb030e4876f0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/go-testing-interface v1.14.1
	github.com/mitchellh/go-wordwrap v1.0.1
	github.com/mitchellh/hashstructure/v2 v2.0.1
	github.com/mitchellh/mapstructure v1.4.1
	github.com/mitchellh/pointerstructure v1.2.0
	github.com/mitchellh/protoc-gen-go-json v0.0.0-20200917194518-364b693410ae
	github.com/mitchellh/reflectwalk v1.0.1 // indirect
	github.com/moby/buildkit v0.8.3 // indirect
	github.com/moby/sys/mountinfo v0.4.1 // indirect
	github.com/moby/sys/symlink v0.1.0 // indirect
	github.com/mr-tron/base58 v1.2.0
	github.com/natefinch/atomic v0.0.0-20200526193002-18c0533a5b09
	github.com/oklog/ulid v1.3.1
	github.com/oklog/ulid/v2 v2.0.2
	github.com/olekukonko/tablewriter v0.0.4
	github.com/opencontainers/image-spec v1.0.1
	github.com/pkg/errors v0.9.1
	github.com/posener/complete v1.2.3
	github.com/r3labs/diff v1.1.0
	github.com/rs/cors v1.7.0 // indirect
	github.com/sabhiram/go-gitignore v0.0.0-20201211210132-54b8a0bf510f // indirect
	github.com/sebdah/goldie/v2 v2.5.3
	github.com/skratchdot/open-golang v0.0.0-20200116055534-eef842397966
	github.com/slack-go/slack v0.6.5
	github.com/stretchr/testify v1.7.0
	github.com/vektra/mockery v1.1.2
	github.com/zclconf/go-cty v1.8.4
	github.com/zclconf/go-cty-yaml v1.0.2
	go.etcd.io/bbolt v1.3.5
	go.uber.org/goleak v1.1.10
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c
	google.golang.org/api v0.25.0
	google.golang.org/genproto v0.0.0-20201022181438-0ff5f38871d5
	google.golang.org/grpc v1.33.1
	google.golang.org/protobuf v1.26.0
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	k8s.io/api v0.19.4
	k8s.io/apimachinery v0.19.4
	k8s.io/client-go v0.19.4
)

// NOTE(mitchellh): I'm keeping these commented and in here because during
// development at the moment it is common to be working on these libs too.
// replace github.com/hashicorp/go-argmapper => ../go-argmapper
// replace github.com/hashicorp/horizon => ../horizon
// replace github.com/hashicorp/waypoint-plugin-sdk => ../waypoint-plugin-sdk

// replace github.com/hashicorp/go-plugin => ../go-plugin

// v0.3.11 panics for some reason on our tests
replace github.com/imdario/mergo => github.com/imdario/mergo v0.3.9
