/*
 * Copyright (c) 2022  The AdvantEDGE Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package meepdaimgr

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

func cmdExec(cli string) (appExecEntry AppExecEntry, err error) {
	parts := strings.Fields(cli)
	head := parts[0]
	parts = parts[1:]
	// fmt.Println("cmdExec: head: ", head)
	// fmt.Println("cmdExec: parts: ", parts)

	appExecEntry.cmd = exec.Command(head, parts...)
	appExecEntry.stdout = new(bytes.Buffer)
	appExecEntry.stderr = new(bytes.Buffer)
	appExecEntry.cmd.Stdout = appExecEntry.stdout
	appExecEntry.cmd.Stderr = appExecEntry.stderr
	err = appExecEntry.cmd.Start()
	if err != nil {
		fmt.Println("error in exec command: ", err, " for command: ", cli)
		fmt.Println("detailed output: ", appExecEntry.stderr.String(), "---", appExecEntry.stdout.String())
		return appExecEntry, err
	}
	fmt.Println("cmdExec: process started: ", appExecEntry.cmd.Process.Pid)

	go func() {
		err = appExecEntry.cmd.Wait()
		// Flush process logs
		fmt.Println("////////////////////////// stdout for process: ", appExecEntry.cmd.Process.Pid)
		fmt.Println("out:", appExecEntry.stdout.String())
		fmt.Println("////////////////////////// sterr for process: ", appExecEntry.cmd.Process.Pid)
		fmt.Println("err:", appExecEntry.stderr.String())
		fmt.Println("////////////////////////// Terminated process: ", appExecEntry.cmd.Process.Pid)
		if err == nil {
			fmt.Println("process terminated normally for pid ", appExecEntry.cmd.Process.Pid)
		} else {
			fmt.Println("process terminated abnormally for pid :", appExecEntry.cmd.Process.Pid, ", ", err)
		}
	}()

	//fmt.Println("cmdExec: appExecEntry: ", appExecEntry)
	return appExecEntry, nil
}

func terminatePidProcess(pid int) {
	if pid <= 0 {
		return
	}

	str := "kill " + strconv.Itoa(pid) // SIGTERM
	_, _ = cmdExec(str)
}

func deletePidProcess(pid int) {
	if pid <= 0 {
		return
	}

	str := "kill -9 " + strconv.Itoa(pid) // SIGKILL
	_, _ = cmdExec(str)
}

func pidExists(pid int) (bool, error) {
	if pid <= 0 {
		return false, fmt.Errorf("invalid pid %v", pid)
	}

	proc, err := os.FindProcess(int(pid))
	if err != nil {
		return false, err
	}

	err = proc.Signal(syscall.Signal(0))
	if err == nil {
		return true, nil
	}
	if err.Error() == "os: process already finished" {
		return false, nil
	}

	errno, ok := err.(syscall.Errno)
	if !ok {
		return false, err
	}
	switch errno {
	case syscall.ESRCH:
		return false, nil
	case syscall.EPERM:
		return true, nil
	}
	return false, err
}

func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

// ExecuteCmd wraps exec.Command.CombinedOutput with verbosity
func ExecuteCmd(cmd *exec.Cmd, cobraCmd *cobra.Command) (output string, err error) {
	err = nil
	verbose, _ := cobraCmd.Flags().GetBool("verbose")

	if verbose {
		fmt.Println("Cmd:", cmd.Args)
	}
	start := time.Now()
	out, err := cmd.CombinedOutput()
	elapsed := time.Since(start)
	output = string(out)
	if verbose {
		var r string
		if len(output) > 0 {
			r = "Result: " + output
		} else {
			r = "Result: OK"
		}
		r = FormatResult(r, elapsed, cobraCmd)
		fmt.Println(r)
	}
	return output, err
}

// GetProcess get a running process ID by name
func GetProcess(name string, cobraCmd *cobra.Command) (string, error) {
	cmd := exec.Command("pidof", name)
	pid, err := ExecuteCmd(cmd, cobraCmd)
	if err != nil {
		return pid, err
	}
	pid = strings.TrimSuffix(pid, "\n")
	return pid, nil
}

// WaitProcess get a running process ID by name
func WaitProcess(pid string, timeout string, cobraCmd *cobra.Command) error {
	cmd := exec.Command("timeout", timeout, "tail", "--pid="+pid, "-f", "/dev/null")
	_, err := ExecuteCmd(cmd, cobraCmd)
	return err
}

// KillProcess kill a process by PID
func KillProcess(pid string, cobraCmd *cobra.Command) {
	cmd := exec.Command("kill", "-9", pid)
	_, err := ExecuteCmd(cmd, cobraCmd)
	if err != nil {
		fmt.Println("Error terminating " + pid)
		fmt.Println(err)
	}
}

// InterruptProcess SIGINT a process by PID
func InterruptProcess(pid string, cobraCmd *cobra.Command) {
	cmd := exec.Command("kill", "-2", pid)
	_, err := ExecuteCmd(cmd, cobraCmd)
	if err != nil {
		fmt.Println("Error interrupting " + pid)
		fmt.Println(err)
	}
}
