package utils

import (
	"bytes"
	"os/exec"
)

// 调用shell脚本执行
func ExecShell(s string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", s)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	// 启动调用
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out.String(), nil
}
