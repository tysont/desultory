package desultory

import (
	"fmt"
	"github.com/shurcooL/go/osutil"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"syscall"
)

func RunCommand(command string, args []string, variables map[string]string, directory string) error {
	_, err := exec.LookPath(command)
	if err != nil {
		p := path.Join(directory, "/", command)
		if _, err := os.Stat(p); os.IsNotExist(err) {
			return err
		}
		command = p
	}
	cmd := exec.Command(command, args ...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = directory
	env := osutil.Environ(os.Environ())
	if variables != nil {
		for k, v := range variables {
			env.Set(k, v)
		}
	}
	cmd.Env = env
	logrus.Infof("Running command: '%v'", cmd.String())
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("error running command: '%v'", err)
	}
	err = cmd.Wait()
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			if status, ok := e.Sys().(syscall.WaitStatus); ok {
				return fmt.Errorf("command returned non-zero exit code '%v': '%v'", status, err)
			}
		} else {
			return fmt.Errorf("error running command: '%v'", err)
		}
	}
	return nil
}

func RunCommandInTempDirectory(command string, args []string, inputFiles map[string][]byte) (map[string][]byte, error) {
	fs := make(map[string][]byte)
	d, err := ioutil.TempDir("/tmp", "run-")
	defer os.RemoveAll(d)
	if err != nil {
		return fs, err
	}
	err = WriteFilesToDirectory(inputFiles, d)
	if err != nil {
		return fs, err
	}
	err = RunCommand(command, args, nil, d)
	if err != nil {
		return fs, err
	}
	fs, err = GetFilesFromDirectory(d)
	if err != nil {
		return fs, err
	}
	return fs, nil
}
