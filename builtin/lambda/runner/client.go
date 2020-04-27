package runner

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
	"github.com/mattn/go-isatty"
	"golang.org/x/crypto/ssh"

	"github.com/hashicorp/waypoint/sdk/terminal"
)

type ConsoleClient struct {
	Tunnel    *Tunnel
	AppConfig *LambdaConfiguration
}

func NewConsoleClient(host string) (*ConsoleClient, error) {
	tun, err := CreateTunnel(host)
	if err != nil {
		return nil, err
	}

	return &ConsoleClient{Tunnel: tun}, nil
}

func (c *ConsoleClient) UseApp(cfg *LambdaConfiguration) {
	c.AppConfig = cfg
}

func (c *ConsoleClient) Exec(ui terminal.UI, name, cmd string) error {
	S := ui.Status()
	defer S.Close()

	S.Update("connecting to tunnel broker")
	conn, err := c.Tunnel.Connect()
	if err != nil {
		return err
	}

	cfg := &ssh.ClientConfig{}
	cfg.SetDefaults()
	cfg.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	S.Update("establishing console session")
	sshc, chans, reqs, err := ssh.NewClientConn(conn, "tunnel", cfg)
	if err != nil {
		return err
	}

	client := ssh.NewClient(sshc, chans, reqs)

	sess, err := client.NewSession()
	if err != nil {
		return err
	}

	if c.AppConfig != nil {
		data, err := json.Marshal(c.AppConfig)
		if err != nil {
			return err
		}

		sess.Setenv("DEVFLOW_CONFIG", base64.RawURLEncoding.EncodeToString(data))
	}

	S.Update("configuring session parameters")
	sess.Stdin = os.Stdin
	sess.Stdout = os.Stdout
	sess.Stderr = os.Stderr

	sess.Setenv("TERM", os.Getenv("TERM"))
	sess.Setenv("APP_NAME", name)

	var cw chan os.Signal

	if isatty.IsTerminal(os.Stdout.Fd()) {
		h, w, err := pty.Getsize(os.Stdout)
		if err != nil {
			return err
		}

		term := os.Getenv("TERM")

		sess.RequestPty(term, h, w, nil)

		cw = make(chan os.Signal, 1)
		signal.Notify(cw, syscall.SIGWINCH)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM,
		syscall.SIGALRM, syscall.SIGUSR1, syscall.SIGUSR2)

	S.Close()

	if cmd == "" {
		sess.Shell()
	} else {
		sess.Start(cmd)
	}

	go func() {
		for {
			select {
			case <-cw:
				h, w, err := pty.Getsize(os.Stdout)

				if err == nil {
					sess.WindowChange(h, w)
				}
			case sig := <-sigs:
				switch sig {
				case syscall.SIGINT:
					sess.Signal(ssh.SIGINT)
				case syscall.SIGQUIT:
					sess.Signal(ssh.SIGQUIT)
				case syscall.SIGTERM:
					sess.Signal(ssh.SIGTERM)
				case syscall.SIGALRM:
					sess.Signal(ssh.SIGALRM)
				case syscall.SIGUSR1:
					sess.Signal(ssh.SIGUSR1)
				case syscall.SIGUSR2:
					sess.Signal(ssh.SIGUSR2)
				}
			}
		}
	}()

	return sess.Wait()
}
