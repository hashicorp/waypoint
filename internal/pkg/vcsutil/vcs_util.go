package vcsutil

import (
	"fmt"
	"os/exec"
)

func IsDirty(path string) (bool, error) {
	cmd := exec.Command("git", "-C", path, "status", "-s")
	stdout, err := cmd.Output()
	if err != nil {
		return false, err
	}
	if len(stdout) > 0 {
		return true, nil
	}

	diff, err := remoteHasDiff(path)
	if err != nil {
		return false, err
	}

	return diff, nil
}

func remoteHasDiff(path string) (bool, error) {

	// check that local dir has a default git remote
	cmd := exec.Command("git", "-C", path, "rev-parse", "--abbrev-ref", "HEAD")
	fmt.Println(cmd)
	sdout, err := cmd.Output()
	if err != nil {
		return false, err
	}

	cmd = exec.Command("git", "-C", path, "config", "branch."+string(sdout)+".remote")
	remote, err := cmd.Output()
	fmt.Println(string(remote))
	if err != nil {
		return false, err
	}

	// cmd = exec.Command("git", "-C", path, "diff", string(remote))
	// stdout, err := cmd.Output()
	// if err != nil {
	// 	return false, err
	// }
	// if len(stdout) > 0 {
	// 	return true, nil
	// }

	return false, nil
}
