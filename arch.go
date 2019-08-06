package main

import (
	"fmt"
	"os"
	"os/exec"
)

func archBootstrap(dst string) error {
	fmt.Println("Bootstrapping arch:", dst)

	cmd := exec.Command("pacstrap", "-c", dst, "base")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func archInstall(dst string, packages ...string) error {
	if len(packages) == 0 {
		return nil
	}

	fmt.Printf("Installing in root (%s): %s\n", dst, packages)

	args := []string{"--root", dst, "-S", "--needed", "--noconfirm"}
	args = append(args, packages...)
	cmd := exec.Command("pacman", args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
