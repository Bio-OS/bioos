//
// Copyright 2023 Beijing Volcano Engine Technology Ltd.
// Copyright 2023 Guangzhou Laboratory
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build linux

package exec

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"syscall"
	"time"
)

func Exec(ctx context.Context, timeout time.Duration, name string, arg ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, arg...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		//	Credential: &syscall.Credential{Uid: uint32(1), Gid: uint32(1)},
		//	Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET,
	}

	var b bytes.Buffer
	// Set process IO
	cmd.Stdout = &b
	cmd.Stderr = &b

	// Start process
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	// Wait for process to complete or timeout
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		// Command was timeout
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) // Kill process group
		return nil, fmt.Errorf("command timeout:%w", ctx.Err())
	case err := <-done:
		if err != nil {
			return nil, err
		}
	}

	return b.Bytes(), nil
}
