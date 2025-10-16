package security

import (
	"fmt"
	"os/exec"
)

func ScanImageWithAnchore(image string) (string, error) {
	cmd := exec.Command("anchore", "image", "vuln", image, "all")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("anchore scan failed: %v\n%s", err, string(out))
	}
	return string(out), nil
}
