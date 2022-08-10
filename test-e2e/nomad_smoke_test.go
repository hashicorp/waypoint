package test

import (
	"fmt"
	"strings"
	"testing"
)

var (
	nomadTestDir = fmt.Sprintf("%s/nomad/nodejs", examplesRootDir)
)

func TestWaypointNomadInstall(t *testing.T) {
	wp := NewBinary(t, wpBinary, nomadTestDir)
	stdout, stderr, err := wp.RunRaw("install", "-platform=nomad", "-accept-tos", fmt.Sprintf("-nomad-server-image=%s", wpServerImage), fmt.Sprintf("-nomad-odr-image=%s", wpOdrImage))

	if err != nil {
		t.Errorf("unexpected error installing server to nomad: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output installing server to nomad: %s", stderr)
	}

	if !strings.Contains(stdout, "Waypoint server successfully installed and configured!") {
		t.Errorf("No success message detected after nomad server install:\n%s", stdout)
	}
}

func TestWaypointNomadUp(t *testing.T) {
	wp := NewBinary(t, wpBinary, nomadTestDir)
	stdout, stderr, err := wp.RunRaw("init")

	if err != nil {
		t.Errorf("unexpected error initializing waypoint project: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output initializing waypoint project: %s", stderr)
	}

	if !strings.Contains(stdout, "Project initialized!") {
		t.Errorf("No success message detected after initializing project:\n%s", stdout)
	}

	stdout, stderr, err = wp.RunRaw("up")

	if err != nil {
		t.Errorf("unexpected error deploying waypoint project: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output deploying waypoint project: %s", stderr)
	}

	if !strings.Contains(stdout, "The deploy was successful!") {
		t.Errorf("No success message detected after deploying project:\n%s", stdout)
	}
}

func TestWaypointNomadUpgrade(t *testing.T) {
	wp := NewBinary(t, wpBinary, nomadTestDir)
	stdout, stderr, err := wp.RunRaw("server", "upgrade", "-platform=nomad", "-auto-approve", fmt.Sprintf("-nomad-server-image=%s", wpServerImageUpgrade), fmt.Sprintf("-nomad-odr-image=%s", wpOdrImageUpgrade), "-snapshot=false")

	if err != nil {
		t.Errorf("unexpected error upgrading server in nomad: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output upgrading server in nomad: %s", stderr)
	}

	if !strings.Contains(stdout, "Waypoint has finished upgrading the server") {
		t.Errorf("No success message detected after nomad server install:\n%s", stdout)
	}
}

func TestWaypointNomadUpAfterUpgrade(t *testing.T) {
	wp := NewBinary(t, wpBinary, nomadTestDir)
	stdout, stderr, err := wp.RunRaw("up")

	if err != nil {
		t.Errorf("unexpected error deploying waypoint project: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output deploying waypoint project: %s", stderr)
	}

	if !strings.Contains(stdout, "The deploy was successful!") {
		t.Errorf("No success message detected after deploying project:\n%s", stdout)
	}
}

func TestWaypointNomadDestroy(t *testing.T) {
	wp := NewBinary(t, wpBinary, nomadTestDir)
	stdout, stderr, err := wp.RunRaw("destroy")

	if err != nil {
		t.Errorf("unexpected error destroying waypoint project: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output destroying waypoint project: %s", stderr)
	}

	if !strings.Contains(stdout, "Destroy successful!") {
		t.Errorf("No success message detected after destroying project:\n%s", stdout)
	}
}

func TestWaypointNomadUninstall(t *testing.T) {
	wp := NewBinary(t, wpBinary, nomadTestDir)
	stdout, stderr, err := wp.RunRaw("server", "uninstall", "-platform=nomad", "-auto-approve", "-snapshot=false")

	if err != nil {
		t.Errorf("unexpected error uninstalling waypoint server: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output uninstalling waypoint server: %s", stderr)
	}

	if !strings.Contains(stdout, "Waypoint server successfully uninstalled") {
		t.Errorf("No success message detected after uninstalling server:\n%s", stdout)
	}
}
