//go:build linux
// +build linux

package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

func run() {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func child() {
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cgroup()

	syscall.Sethostname([]byte("container"))
	syscall.Chroot("/root/tiny-contain")
	syscall.Chdir("/")
	syscall.Mount("proc", "proc", "proc", 0, "")

	cmd.Run()
}

func cgroup() {
	cgroups := "/sys/fs/cgroup/"
	pids := filepath.Join(cgroups, "pids")
	os.Mkdir(filepath.Join(pids, "tiny"), 0755)
	ioutil.WriteFile(filepath.Join(pids, "tiny/pids.max"), []byte("100"), 0700)
	ioutil.WriteFile(filepath.Join(pids, "tiny/notify_on_release"), []byte("1"), 0700)
	ioutil.WriteFile(filepath.Join(pids, "tiny/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700)
}

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("ERROR: Not valid command")
	}
}
