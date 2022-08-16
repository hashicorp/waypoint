package test

import (
	"fmt"
	"strings"
	"testing"
)

var (
	dockerTestDir         = fmt.Sprintf("%s/docker/go", examplesRootDir)
	dockerMultiAppTestDir = fmt.Sprintf("%s/docker/go-multiapp", examplesRootDir)
)

func TestWaypointDockerInstall(t *testing.T) {
	wp := NewBinary(t, wpBinary, dockerTestDir)
	stdout, stderr, err := wp.RunRaw("install", "-platform=docker", "-accept-tos", fmt.Sprintf("-docker-server-image=%s", wpServerImage), fmt.Sprintf("-docker-odr-image=%s", wpOdrImage))

	if err != nil {
		t.Errorf("unexpected error installing server to docker: %s\nstderr: %s", err, stderr)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output installing server to docker: %s", stderr)
	}

	if !strings.Contains(stdout, "Waypoint server successfully installed and configured!") {
		t.Errorf("No success message detected after docker server install:\n%s", stdout)
	}
}

func TestWaypointDockerUp(t *testing.T) {
	wp := NewBinary(t, wpBinary, dockerTestDir)
	stdout, stderr, err := wp.RunRaw("init")

	if err != nil {
		t.Errorf("unexpected error initializing waypoint project: %s\nstderr: %s", err, stderr)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output initializing waypoint project: %s", stderr)
	}

	if !strings.Contains(stdout, "Project initialized!") {
		t.Errorf("No success message detected after initializing project:\n%s", stdout)
	}

	stdout, stderr, err = wp.RunRaw("up")

	if err != nil {
		t.Errorf("unexpected error deploying waypoint project: %s\nstderr: %s", err, stderr)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output deploying waypoint project: %s", stderr)
	}

	if !strings.Contains(stdout, "The deploy was successful!") {
		t.Errorf("No success message detected after deploying project:\n%s", stdout)
	}
}

func TestWaypointDockerMultiAppUp(t *testing.T) {
	wp := NewBinary(t, wpBinary, dockerMultiAppTestDir)
	stdout, stderr, err := wp.RunRaw("init")

	if err != nil {
		t.Errorf("unexpected error initializing waypoint project: %s\nstderr: %s", err, stderr)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output initializing waypoint project: %s", stderr)
	}

	if !strings.Contains(stdout, "Project initialized!") {
		t.Errorf("No success message detected after initializing project:\n%s", stdout)
	}

	stdout, stderr, err = wp.RunRaw("up")

	if err != nil {
		t.Errorf("unexpected error deploying waypoint project: %s\nstderr: %s", err, stderr)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output deploying waypoint project: %s", stderr)
	}

	if !strings.Contains(stdout, "The deploy was successful!") {
		t.Errorf("No success message detected after deploying project:\n%s", stdout)
	}
}

func TestWaypointDockerUpgrade(t *testing.T) {
	wp := NewBinary(t, wpBinary, dockerTestDir)
	stdout, stderr, err := wp.RunRaw("server", "upgrade", "-platform=docker", "-auto-approve", fmt.Sprintf("-docker-server-image=%s", wpServerImageUpgrade), fmt.Sprintf("-docker-odr-image=%s", wpOdrImageUpgrade), "-snapshot=false")

	if err != nil {
		t.Errorf("unexpected error upgrading server in docker: %s\nstderr: %s", err, stderr)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output upgrading server in docker: %s", stderr)
	}

	if !strings.Contains(stdout, "Waypoint has finished upgrading the server") {
		t.Errorf("No success message detected after docker server install:\n%s", stdout)
	}
}

func TestWaypointDockerUpAfterUpgrade(t *testing.T) {
	wp := NewBinary(t, wpBinary, dockerTestDir)
	stdout, stderr, err := wp.RunRaw("up")

	if err != nil {
		t.Errorf("unexpected error deploying waypoint project: %s\nstderr: %s", err, stderr)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output deploying waypoint project: %s", stderr)
	}

	if !strings.Contains(stdout, "The deploy was successful!") {
		t.Errorf("No success message detected after deploying project:\n%s", stdout)
	}
}

func TestWaypointDockerMultiAppUpAfterUpgrade(t *testing.T) {
	wp := NewBinary(t, wpBinary, dockerMultiAppTestDir)
	stdout, stderr, err := wp.RunRaw("init")

	if err != nil {
		t.Errorf("unexpected error initializing waypoint project: %s\nstderr: %s", err, stderr)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output initializing waypoint project: %s", stderr)
	}

	if !strings.Contains(stdout, "Project initialized!") {
		t.Errorf("No success message detected after initializing project:\n%s", stdout)
	}

	stdout, stderr, err = wp.RunRaw("up")

	if err != nil {
		t.Errorf("unexpected error deploying waypoint project: %s\nstderr: %s", err, stderr)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output deploying waypoint project: %s", stderr)
	}

	if !strings.Contains(stdout, "The deploy was successful!") {
		t.Errorf("No success message detected after deploying project:\n%s", stdout)
	}
}

func TestWaypointDockerDestroy(t *testing.T) {
	wp := NewBinary(t, wpBinary, dockerTestDir)
	stdout, stderr, err := wp.RunRaw("destroy", "-auto-approve")

	if err != nil {
		t.Errorf("unexpected error destroying waypoint project: %s\nstderr: %s", err, stderr)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output destroying waypoint project: %v", stderr)
	}

	if !strings.Contains(stdout, "Destroy successful!") {
		t.Errorf("No success message detected after destroying project:\n%s", stdout)
	}
}

func TestWaypointDockerDestroyMultiApp(t *testing.T) {
	wp := NewBinary(t, wpBinary, dockerMultiAppTestDir)
	stdout, stderr, err := wp.RunRaw("destroy", "-auto-approve")

	if err != nil {
		t.Errorf("unexpected error destroying waypoint project: %s\nstderr: %s", err, stderr)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output destroying waypoint project: %v", stderr)
	}

	if !strings.Contains(stdout, "Destroy successful!") {
		t.Errorf("No success message detected after destroying project:\n%s", stdout)
	}
}

func TestWaypointDockerUninstall(t *testing.T) {
	wp := NewBinary(t, wpBinary, dockerTestDir)
	stdout, stderr, err := wp.RunRaw("server", "uninstall", "-platform=docker", "-auto-approve", "-snapshot=false")

	if err != nil {
		t.Errorf("unexpected error uninstalling waypoint server: %s\nstderr: %s", err, stderr)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output uninstalling waypoint server: %s", stderr)
	}

	if !strings.Contains(stdout, "Waypoint server successfully uninstalled") {
		t.Errorf("No success message detected after uninstalling server:\n%s", stdout)
	}
}
