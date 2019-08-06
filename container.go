package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func cg() error {
	cgroups := "/sys/fs/cgroup/"
	pids := filepath.Join(cgroups, "pids")
	os.Mkdir(filepath.Join(pids, "scratch"), 0755)

	if err := ioutil.WriteFile(filepath.Join(pids, "scratch/pids.max"), []byte("1024"), 0700); err != nil {
		return err
	}

	// Removes the new cgroup in place after the container exits
	if err := ioutil.WriteFile(filepath.Join(pids, "scratch/notify_on_release"), []byte("1"), 0700); err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(pids, "scratch/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700)
}

func importEnv() []string {
	env := os.Environ()

	// TODO: filter environment variables
	return env
}

func child(root string) error {
	if err := cg(); err != nil {
		return err
	}

	args := flag.Args()
	if len(args) < 3 {
		args = append(args, "/bin/bash")
	}
	fmt.Printf("Executing in container (%s): %s\n", root, strings.Join(args[2:], " "))

	// set hostname
	if err := syscall.Sethostname([]byte("container")); err != nil {
		return err
	}

	// bind-mount /dev from outside
	if err := syscall.Mount("/dev", filepath.Join(root, "dev"), "devtmpfs",
		syscall.MS_BIND|syscall.MS_NOSUID, ""); err != nil {
		return err
	}
	defer func() {
		fmt.Println("Unmounting /dev")
		if err := syscall.Unmount("dev", 0); err != nil {
			panic(err)
		}
	}()

	// setup chroot
	if err := syscall.Chroot(root); err != nil {
		return err
	}
	if err := os.Chdir("/"); err != nil {
		return err
	}

	// mount /proc
	if err := syscall.Mount("proc", "proc", "proc", 0, ""); err != nil {
		return err
	}
	defer func() {
		fmt.Println("Unmounting /proc")
		if err := syscall.Unmount("proc", 0); err != nil {
			panic(err)
		}
	}()

	// mount /tmp
	if err := syscall.Mount("tmpfs", "tmp", "tmpfs", 0, ""); err != nil {
		return err
	}
	defer func() {
		fmt.Println("Unmounting /tmp")
		if err := syscall.Unmount("tmp", 0); err != nil {
			panic(err)
		}
	}()

	// execute command in container
	return execWithCredentials(uint32(*uid), uint32(*gid), args[2], args[3:]...)
}

func execWithCredentials(uid, gid uint32, name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = importEnv()
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{Uid: uid, Gid: gid},
	}

	return cmd.Run()
}

func hostCopy(dst string, src string) (err error) {
	dst = filepath.Join(dst, src)
	fmt.Println("Copying", src, "to", dst)

	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	stat, err := os.Stat(src)
	if err != nil {
		return
	}
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, stat.Mode())
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
