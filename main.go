package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

var (
	uid = flag.Uint("uid", 0, "uid to run as (inside container)")
	gid = flag.Uint("gid", 0, "gid to run as (inside container)")
)

func printUsage() {
	fmt.Println(`Usage:
scratch bootstrap <root>
scratch run <root> <cmd> <args>`)
	os.Exit(1)
}

func main() {
	flag.Parse()
	if len(flag.Args()) < 2 {
		printUsage()
	}

	root, err := filepath.Abs(flag.Args()[1])
	if err != nil {
		fmt.Println("Invalid root:", err)
		os.Exit(1)
	}

	switch flag.Args()[0] {
	case "bootstrap":
		err := os.MkdirAll(root, 0755)
		if err != nil {
			fmt.Println("Cannot create root:", err)
			os.Exit(1)
		}
		if err := archBootstrap(root); err != nil {
			fmt.Println("Bootstrap failed:", err)
			os.Exit(1)
		}
		if err := archInstall(root, flag.Args()[2:]...); err != nil {
			fmt.Println("Installing packages failed:", err)
			os.Exit(1)
		}

		fmt.Println("Successfully bootstrapped root:", root)

	case "run":
		if err := hostCopy(root, "/etc/resolv.conf"); err != nil {
			fmt.Println("Cannot copy resolv.conf:", err)
			os.Exit(1)
		}

		if err := run(root); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	case "child":
		if err := child(root); err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				// exited with exit code != 0
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					// fmt.Println(err)
					os.Exit(status.ExitStatus())
				}
			}

			os.Exit(1)
		}

	default:
		printUsage()
	}
}

func run(root string) error {
	// fmt.Printf("Running %v\n", flag.Args()[1:])

	// bind-mount root so we have a proper mtab in the container
	if err := syscall.Mount(root, root, "", syscall.MS_BIND, ""); err != nil {
		return err
	}
	defer func() {
		fmt.Println("Unmounting root:", root)
		if err := syscall.Unmount(root, 0); err != nil {
			panic(err)
		}
	}()

	// collect params we pass to the child
	params := []string{
		"-uid", strconv.FormatUint(uint64(*uid), 10),
		"-gid", strconv.FormatUint(uint64(*gid), 10),
		"child",
	}
	params = append(params, flag.Args()[1:]...)

	// restart ourselves within our own namespace
	cmd := exec.Command("/proc/self/exe", params...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}
	return cmd.Run()
}
