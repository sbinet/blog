package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	buildDir = "_build"
	blogDir  = "sbinet.github.io"
)

func main() {
	dir := filepath.Join(buildDir, blogDir)

	os.RemoveAll(dir)

	run("git", "clone", "git@github.com:sbinet/sbinet.github.io", dir)
	run("/bin/sh", "-c",
		fmt.Sprintf(
			"/bin/cp -rf %v/* %v/.",
			filepath.Join(buildDir, "public"),
			dir,
		),
	)

	err := os.Chdir(dir)
	if err != nil {
		log.Fatal(err)
	}

	run("git", "add", "-A", ".")
	run("git", "commit", "-m", "update "+time.Now().UTC().Format("2006-01-02"))
	run("git", "push", "origin", "master")
}

func run(cmd string, args ...string) {
	c := exec.Command(cmd, args...)
	fmt.Printf("+ %s\n", strings.Join(c.Args, " "))
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		log.Fatal(err)
	}
}
