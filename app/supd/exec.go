package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"syscall"
	"time"

	"github.com/gwaylib/errors"
	log "github.com/sirupsen/logrus"
)

type ExecCommand struct {
	// api是否等待返回
	Wait bool
	// 超时时间，此值为0时自动更换为10秒
	Timeout int

	// Args holds command line arguments, including the command as Args[0].
	// If the Args field is empty or nil, Run uses {Path}.
	//
	// In typical use, both Path and Args are set by calling Command.
	Args []string

	// Env specifies the environment of the process.
	// Each entry is of the form "key=value".
	// If Env is nil, the new process uses the current process's
	// environment.
	// If Env contains duplicate environment keys, only the last
	// value in the slice for each duplicate key is used.
	// As a special case on Windows, SYSTEMROOT is always added if
	// missing and not explicitly set to the empty string.
	Env []string

	// Dir specifies the working directory of the command.
	// If Dir is the empty string, Run runs the command in the
	// calling process's current directory.
	Dir string

	// SysProcAttr holds optional, operating system-specific attributes.
	// Run passes it to os.StartProcess as the os.ProcAttr's Sys field.
	SysProcAttr *syscall.SysProcAttr
}

type ExecStdout struct {
	data []byte
}

func (c *ExecStdout) Write(data []byte) (int, error) {
	c.data = append(c.data, data...)
	return len(data), nil
}

type ExecStderr struct {
	data []byte
}

func (c *ExecStderr) Write(data []byte) (int, error) {
	c.data = append(c.data, data...)
	return len(data), nil
}

// example:
// curl -v -X PUT --data '{"Args":["/sbin/ping","baidu.com"],"Wait":true}' "http://127.0.0.1:9001/exec"
func (sr *SupervisorRestful) ExecCommand(w http.ResponseWriter, req *http.Request) {
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.Write([]byte(errors.As(err).Error()))
		return
	}
	log.Info("Exec:" + string(data))
	inputCmd := &ExecCommand{}
	if err := json.Unmarshal(data, inputCmd); err != nil {
		w.WriteHeader(403)
		w.Write([]byte(errors.As(err, string(data)).Error()))
		return
	}

	timeout := time.Duration(inputCmd.Timeout * 1e6)
	if timeout == 0 {
		timeout = 10 * 1e9
	}
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	cmd := exec.CommandContext(ctx, inputCmd.Args[0])
	cmd.Args = inputCmd.Args
	cmd.Env = inputCmd.Env
	cmd.Dir = inputCmd.Dir
	cmd.SysProcAttr = inputCmd.SysProcAttr
	cmd.Stdout = &ExecStdout{}
	cmd.Stderr = &ExecStderr{}

	if err := cmd.Start(); err != nil {
		w.WriteHeader(403)
		w.Write([]byte(errors.As(err).Error()))
		return
	}

	done := make(chan string, 1)
	go func() {
		if err := cmd.Wait(); err != nil {
			log.Info(err.Error())
		}
		stdout := cmd.Stdout.(*ExecStdout).data
		stderr := cmd.Stderr.(*ExecStderr).data
		output := fmt.Sprintf("exec:\n%+v\nstdout:\n%s\nstderr:\n%s\n", *inputCmd, string(stdout), string(stderr))
		log.Println(output)
		if !inputCmd.Wait {
			close(done)
			return
		}
		done <- output
	}()

	if !inputCmd.Wait {
		w.Write([]byte("done"))
		return
	}
	output := <-done
	close(done)
	w.Write([]byte(output))
}
