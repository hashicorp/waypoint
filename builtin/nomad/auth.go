package nomad

import "os"

func ConsulAuth() (string, error) {
	return os.Getenv("CONSUL_HTTP_TOKEN"), nil
}

func VaultAuth() (string, error) {
	return os.Getenv("VAULT_TOKEN"), nil
}
