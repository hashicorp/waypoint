package test

import (
	"fmt"
	"strings"
	"testing"
)

var (
	dockerTestDir = fmt.Sprintf("%s/docker/go", examplesRootDir)
)

func TestWaypointDockerInstall(t *testing.T) {
	wp := NewBinary(wpBinary, dockerTestDir)
	stdout, stderr, err := wp.Run("install", "-platform=docker", "-accept-tos", fmt.Sprintf("-docker-server-image=%s", wpServerImage))

	if err != nil {
		t.Errorf("unexpected error installing server to docker: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output installing server to docker: %s", err)
	}

	if !strings.Contains(stdout, "Waypoint server successfully installed and configured!") {
		t.Errorf("No success message detected after docker server install:\n%s", stdout)
	}
}

func TestWaypointDockerUp(t *testing.T) {
	wp := NewBinary(wpBinary, dockerTestDir)
	stdout, stderr, err := wp.Run("init")

	if err != nil {
		t.Errorf("unexpected error initializing waypoint project: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output initializing waypoint project: %s", err)
	}

	if !strings.Contains(stdout, "Project initialized!") {
		t.Errorf("No success message detected after initializing project:\n%s", stdout)
	}

	stdout, stderr, err = wp.Run("up")

	if err != nil {
		t.Errorf("unexpected error deploying waypoint project: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output deploying waypoint project: %s", err)
	}

	if !strings.Contains(stdout, "The deploy was successful!") {
		t.Errorf("No success message detected after deploying project:\n%s", stdout)
	}
}

func TestWaypointDockerUpgrade(t *testing.T) {
	wp := NewBinary(wpBinary, dockerTestDir)
	stdout, stderr, err := wp.Run("server", "upgrade", "-platform=docker", "-auto-approve", fmt.Sprintf("-docker-server-image=%s", wpServerImageUpgrade), "-snapshot=false")

	if err != nil {
		t.Errorf("unexpected error upgrading server in docker: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output upgrading server in docker: %s", err)
	}

	if !strings.Contains(stdout, "Waypoint has finished upgrading the server") {
		t.Errorf("No success message detected after docker server install:\n%s", stdout)
	}
}

func TestWaypointDockerUpAfterUpgrade(t *testing.T) {
	wp := NewBinary(wpBinary, dockerTestDir)
	stdout, stderr, err := wp.Run("up")

	if err != nil {
		t.Errorf("unexpected error deploying waypoint project: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output deploying waypoint project: %s", err)
	}

	if !strings.Contains(stdout, "The deploy was successful!") {
		t.Errorf("No success message detected after deploying project:\n%s", stdout)
	}
}

func TestWaypointDockerDestroy(t *testing.T) {
	wp := NewBinary(wpBinary, dockerTestDir)
	stdout, stderr, err := wp.Run("destroy")

	if err != nil {
		t.Errorf("unexpected error destroying waypoint project: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output destroying waypoint project: %s", err)
	}

	if !strings.Contains(stdout, "Destroy successful!") {
		t.Errorf("No success message detected after destroying project:\n%s", stdout)
	}
}

func TestWaypointDockerUninstall(t *testing.T) {
	wp := NewBinary(wpBinary, dockerTestDir)
	stdout, stderr, err := wp.Run("server", "uninstall", "-platform=docker", "-auto-approve", "-snapshot=false")

	if err != nil {
		t.Errorf("unexpected error uninstalling waypoint server: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output uninstalling waypoint server: %s", err)
	}

	if !strings.Contains(stdout, "Waypoint server successfully uninstalled") {
		t.Errorf("No success message detected after uninstalling server:\n%s", stdout)
	}
}
