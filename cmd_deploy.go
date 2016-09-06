package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	blogDir = "sbinet.github.io"
)

func main() {
	os.RemoveAll(blogDir)
	run("git", "clone", "git@github.com:sbinet/sbinet.github.io", blogDir)
	run("/bin/sh", "-c", "/bin/cp -rf public/* "+blogDir+"/.")

	err := os.Chdir(blogDir)
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
