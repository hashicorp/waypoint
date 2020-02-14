package runner

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/hashicorp/go-hclog"
)

type Runner struct {
	Runtime string
}

var sess = session.New(aws.NewConfig().WithRegion("us-west-2"))

type LayerConfiguration struct {
	Runtime   string   `json:"runtime"`
	BuildID   string   `json:"build_id"`
	AppUrl    string   `json:"app_url"`
	LayerUrls []string `json:"layer_urls"`
}

func (r *Runner) ExtractFromLambda(name string) (*LayerConfiguration, error) {
	lamSvc := lambda.New(sess)

	var cfg LayerConfiguration

	fnInfo, err := lamSvc.GetFunction(&lambda.GetFunctionInput{
		FunctionName: aws.String(name),
	})

	if err != nil {
		return nil, err
	}

	cfg.Runtime = *fnInfo.Configuration.Runtime
	cfg.BuildID = *fnInfo.Tags["devflow.app.id"]
	cfg.AppUrl = *fnInfo.Code.Location

	for _, layer := range fnInfo.Configuration.Layers {
		ver, err := lamSvc.GetLayerVersionByArn(&lambda.GetLayerVersionByArnInput{
			Arn: layer.Arn,
		})

		if err != nil {
			return nil, err
		}

		cfg.LayerUrls = append(cfg.LayerUrls, *ver.Content.Location)
	}

	return &cfg, nil
}

func (r *Runner) extractZip(L hclog.Logger, base, path string) error {
	if strings.HasPrefix(path, "https://") {
		L.Info("downloading zip", "url", path)

		resp, err := http.Get(path)
		if err != nil {
			return nil
		}

		defer resp.Body.Close()

		f, err := ioutil.TempFile("", "lambda-runner")
		if err != nil {
			return err
		}

		io.Copy(f, resp.Body)

		f.Close()

		defer os.Remove(f.Name())

		path = f.Name()
	}

	zr, err := zip.OpenReader(path)
	if err != nil {
		return err
	}

	defer zr.Close()

	createdDirs := map[string]struct{}{}

	for _, file := range zr.File {
		wPath := filepath.Join(base, file.Name)
		dir := filepath.Dir(wPath)

		if _, ok := createdDirs[dir]; !ok {
			os.MkdirAll(dir, 0755)
			createdDirs[dir] = struct{}{}
		}

		switch file.Mode() & os.ModeType {
		case os.ModeDir:
			os.MkdirAll(wPath, 0755)
		case os.ModeSymlink:
			r, err := file.Open()
			if err != nil {
				return err
			}

			data, err := ioutil.ReadAll(r)
			r.Close()

			if err != nil {
				return err
			}

			err = os.Symlink(string(data), wPath)
			if err != nil {
				return err
			}
		case 0: // reg shows up as 0 here
			err := func() error {
				r, err := file.Open()
				if err != nil {
					return err
				}

				defer r.Close()

				w, err := os.Create(wPath)
				if err != nil {
					return err
				}

				defer w.Close()

				io.Copy(w, r)

				return nil
			}()

			if err != nil {
				return err
			}
		default:
			// skip it for now
		}
	}

	return nil
}

func (r *Runner) SetupEnv(L hclog.Logger, layers []string, app string) error {
	os.MkdirAll("/opt", 0755)

	for _, layer := range layers {
		L.Info("extracting layer", "path", layer)

		err := r.extractZip(L, "/opt", layer)
		if err != nil {
			return err
		}
	}

	os.MkdirAll("/var/task", 0755)

	L.Info("extracting app", "path", app)
	return r.extractZip(L, "/var/task", app)
}

func (r *Runner) ExecTask(L hclog.Logger, str string, args ...string) error {
	L.Info("executing command", "command", str, "args", args)

	out := append([]string{str}, args...)

	cmd := exec.Command("/var/runtime/bootstrap-raw", out...)
	cmd.Dir = "/var/task"
	cmd.Env = []string{
		"PATH=/var/lang/bin:/usr/local/bin:/usr/bin/:/bin:/opt/bin",
		"LD_LIBRARY_PATH=/var/lang/lib:/lib64:/usr/lib64:/var/runtime:/var/runtime/lib:/var/task:/var/task/lib:/opt/lib",
		"LANG=en_US.UTF-8",
		"TZ=:UTC",
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func (r *Runner) Command(L hclog.Logger, args ...string) *exec.Cmd {
	cmd := exec.Command("/var/runtime/bootstrap-raw", args...)
	cmd.Dir = "/var/task"
	cmd.Env = []string{
		"PATH=/var/lang/bin:/usr/local/bin:/usr/bin/:/bin:/opt/bin",
		"LD_LIBRARY_PATH=/var/lang/lib:/lib64:/usr/lib64:/var/runtime:/var/runtime/lib:/var/task:/var/task/lib:/opt/lib",
		"LANG=en_US.UTF-8",
		"TZ=:UTC",
	}

	return cmd
}
