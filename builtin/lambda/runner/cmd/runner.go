package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/creack/pty"
	"github.com/gliderlabs/ssh"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/securetunnel"
	"github.com/hashicorp/waypoint/builtin/lambda/runner"
	"github.com/mattn/go-isatty"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	fConfig  = flag.String("config", "app.json", "app config file")
	fExtract = flag.String("extract", "", "extract config from this lambda")
	fConnect = flag.Bool("connect", false, "connect to the secure tunnel endpoint and act as the client")
	fApp     = flag.String("app", "", "lambda app to extract config from")
	fServe   = flag.Bool("serve", false, "connect to the secure tunnel endpoint and act as the server")
	fDocker  = flag.Bool("docker", false, "launch the server inside local docker automatically")
)

func main() {
	flag.Parse()

	var r runner.Runner

	if *fServe {
		connectAndRun()
		return
	}

	if *fConnect {
		clientConnect()
		return
	}

	if *fExtract != "" {
		cfg, err := r.ExtractFromLambda(*fExtract)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(os.Stdout).Encode(cfg)
		return
	}

	L := hclog.L().Named("runner")

	f, err := os.Open(*fConfig)
	if err != nil {
		log.Fatal(err)
	}

	var cfg runner.LambdaConfiguration

	err = json.NewDecoder(f).Decode(&cfg)

	L.Info("setting up task", "runtime", cfg.Runtime, "app", cfg.AppUrl)

	err = r.SetupEnv(L, cfg.LayerUrls, cfg.AppUrl)
	if err != nil {
		L.Error("error setting up env", "error", err)
		os.Exit(1)
	}

	args := flag.Args()

	var cmd string

	if len(args) == 0 {
		cmd = "/bin/bash"
	} else {
		cmd = args[0]
		args = args[1:]
	}

	err = r.ExecTask(L, cmd, args...)
	if err != nil {
		L.Error("error executing task", "error", err)
		os.Exit(1)
	}
}

type tunConn struct {
	*securetunnel.Tunnel
}

func (t *tunConn) Network() string {
	return "securetunnel"
}

func (t *tunConn) String() string {
	return "tunnel"
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (t *tunConn) Close() error {
	return nil
}

// LocalAddr returns the local network address.
func (t *tunConn) LocalAddr() net.Addr { return t }

// RemoteAddr returns the remote network address.
func (t *tunConn) RemoteAddr() net.Addr { return t }

func (tc *tunConn) SetDeadline(t time.Time) error      { return nil }
func (tc *tunConn) SetReadDeadline(t time.Time) error  { return nil }
func (tc *tunConn) SetWriteDeadline(t time.Time) error { return nil }

type oneShot struct {
	tun *securetunnel.Tunnel
}

func (o *oneShot) Accept() (net.Conn, error) {
	t := o.tun
	if t == nil {
		return nil, io.EOF
	}

	hclog.L().Info("returning tunnel as connection")

	o.tun = nil
	return &tunConn{t}, nil
}

func (o *oneShot) Addr() net.Addr {
	return o
}

func (o *oneShot) String() string {
	return "tunnel-oneshot"
}

func (o *oneShot) Network() string {
	return "securetunnel"
}

func (o *oneShot) Close() error {
	return nil
}

func connectAndRun() {
	L := hclog.L()

	token := os.Getenv("DEVFLOW_TUNNEL_TOKEN")
	key := os.Getenv("DEVFLOW_TUNNEL_KEY")

	tun, err := securetunnel.Open(token, key)
	if err != nil {
		log.Fatal(err)
	}

	L.Info("tunnel openned, starting ssh")

	handled := make(chan struct{})

	serv := &ssh.Server{
		Handler: func(sess ssh.Session) {
			defer close(handled)
			sessionHandler(sess)
		},
	}

	serv.Serve(&oneShot{tun})
	L.Info("serve ended, sleeping")
	<-handled
}

func sessionHandler(sess ssh.Session) {
	L := hclog.L()

	var (
		outEnv []string
		cfgStr []byte
	)

	L.Info("configuring app environment")

	for _, str := range sess.Environ() {
		if strings.HasPrefix(str, "DEVFLOW_CONFIG=") {
			b, err := base64.RawURLEncoding.DecodeString(str[len("DEVFLOW_CONFIG="):])
			if err != nil {
				L.Error("error decoding config", "error", err)
				fmt.Fprintf(sess, "error decoding config: %s\n", err)
				return
			}

			cfgStr = b
		} else {
			outEnv = append(outEnv, str)
		}
	}

	// Setup the lambda-esque env

	var cfg runner.LambdaConfiguration

	err := json.Unmarshal(cfgStr, &cfg)
	if err != nil {
		fmt.Fprintf(sess, "error decoding configuration: %s\n", err)
		L.Error("error decoding configuration", "error", err)
		return
	}

	L.Info("setting up task", "runtime", cfg.Runtime, "app", cfg.AppUrl)

	var r runner.Runner

	err = r.SetupEnv(L, cfg.LayerUrls, cfg.AppUrl)
	if err != nil {
		fmt.Fprintf(sess, "error setting up env: %s\n", err)
		L.Error("error setting up env", "error", err)
		return
	}

	cs := sess.Command()

	cmd := r.Command(L, cs...)
	cmd.Env = append(cmd.Env, outEnv...)

	for k, v := range cfg.Variables {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	if req, win, ok := sess.Pty(); ok {
		cmd.Env = append(cmd.Env, "TERM="+req.Term)

		ws := &pty.Winsize{
			Cols: uint16(req.Window.Width),
			Rows: uint16(req.Window.Height),
		}

		f, err := pty.StartWithSize(cmd, ws)
		if err != nil {
			L.Error("error starting pty", "error", err)
			return
		}

		go func() {
			for update := range win {
				pty.Setsize(f, &pty.Winsize{
					Cols: uint16(update.Width),
					Rows: uint16(update.Height),
				})
			}
		}()

		go io.Copy(f, sess)

		io.Copy(sess, f)
		sess.CloseWrite()
	} else {
		cmd.Stdout = sess
		cmd.Stderr = sess.Stderr()
		cmd.Stdin = sess

		err := cmd.Start()
		if err != nil {
			L.Error("error starting command", "error", err)
			return
		}
	}

	state, err := cmd.Process.Wait()
	if err != nil {
		L.Error("error waiting for command to exit", "error", err)
		return
	}

	sess.Exit(state.ExitCode())
}

func clientConnect() {
	L := hclog.L()

	cc, err := runner.NewConsoleClient(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	if *fDocker {
		L.Info("starting server with docker")
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		exec.Command("docker", "run", "-d",
			"-v", pwd+"/tmp:/input", "-w", "/input",
			"-e", "DEVFLOW_TUNNEL_TOKEN="+cc.Tunnel.ServerToken(),
			"-e", "DEVFLOW_TUNNEL_KEY="+cc.Tunnel.ServerKey(),
			"devflow/lambda:ruby2.5", "./lambda-runner", "-serve").Run()
	} else {
		fmt.Printf("DEVFLOW_TUNNEL_TOKEN=%s DEVFLOW_TUNNEL_KEY=%s ./bin/runner -serve\n", cc.Tunnel.ServerToken(), cc.Tunnel.ServerKey())
	}

	L.Info("executing shell")

	if *fApp != "" {
		var r runner.Runner

		cfg, err := r.ExtractFromLambda(*fApp)
		if err != nil {
			log.Fatal(err)
		}

		cc.UseApp(cfg)
	}

	var (
		fd = os.Stdin.Fd()
		st *terminal.State
	)

	isterm := isatty.IsTerminal(fd)

	if isterm {
		st, err = terminal.MakeRaw(int(fd))
		if err == nil {
			defer terminal.Restore(int(fd), st)
		}
	}

	cc.Exec(nil, "app", "/bin/bash")

	terminal.Restore(int(fd), st)
}
