package helm

import(
	"fmt"

	"github.com/hashicorp/go-hclog"
	"helm.sh/helm/v3/pkg/action"
)

func (p *Platform) actionInit(log hclog.Logger) (*action.Configuration,error){
	// Get our K8S API
	cs, ns, rc, err := clientset(p.config.KubeconfigPath, p.config.Context)
	if err != nil {
		return nil, err
	}

	driver := "secret"
	if v := p.config.Driver; v != "" {
		driver = v
	}

	// For logging, we'll debug log to a custom named logger.
	actionlog := log.Named("helm_action")
	debug := func(format string, v ...interface{}) {
		actionlog.Debug(fmt.Sprintf(format, v...))
	}

	// Initialize our action
	var ac action.Configuration
	err = ac.Init(&restClientGetter{
		RestConfig:  rc,
		Kubeconfig:  p.config.KubeconfigPath,
		Kubecontext: p.config.Context,
	}, ns, driver, debug)
	if err != nil {
		return nil, err
	}

	return &ac, nil
}


