/*
 * Copyright 2024, Inspect Element Ltd <https://echo.ac>.
 *
 * See LICENCE <https://github.com/Inspect-Element-Ltd/vm/blob/master/LICENCE>
 *
 * cmd.go
 * ---
 * Last Modified: 11/05/2024 01:15AM (BST)
 * Modified By: Gianluca Oliva <hello@gian.sh>
 */

//go:build !windows

package util

import (
	"io"
	"log"
	"os/exec"
	"strings"
)

func InvokeCMD(cmd string, params ...string) (string, error) {
	command := exec.Command(cmd, params...)
	out, _ := command.StdoutPipe()
	defer func(out io.ReadCloser) {
		if err := out.Close(); err != nil && !strings.Contains(err.Error(), "already closed") {
			log.Println(err)
		}
	}(out)

	if err := command.Start(); err != nil {
		return "", err
	}

	output := make([]byte, 0)
	for {
		tmp := make([]byte, 256)
		n, _ := out.Read(tmp)
		if n == 0 {
			break
		}
		output = append(output, tmp[:n]...)
	}

	if err := command.Wait(); err != nil {
		return "", err
	}

	return string(output), nil
}
