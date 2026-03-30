package git

import (
	"fmt"
	"os/exec"
)

func CloneRepo(url, branch, dest string) error {
	cmd := exec.Command("git", "clone", "--depth", "1", "--branch", branch, url, dest)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone failed: %s: %w", string(output), err)
	}
	return nil
}
