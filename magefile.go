//go:build mage
// +build mage

package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/magefile/mage/mg"
)

var WD = must(os.Getwd())

var Default = Run

func must[T any](value T, err error) T {
	if err != nil {
		log.Fatal(err)
	}

	return value
}

func must0(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func run(name string, args ...string) {
	cmd(name, args...)()
}

func cmd(name string, args ...string) func(...string) {
	cmd := exec.Command(name, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return func(args ...string) {
		cmd.Args = append(cmd.Args, args...)
		must0(cmd.Run())
	}
}

// Start the server
func Run() {
	run("go", "run", filepath.Join(WD, "cmd", "app"))
}

type Docker mg.Namespace

// Rebuilds the Dockerfile
func (Docker) Rebuild() {
	compose := cmd("docker", "compose")

	compose("up", "-d", "--no-deps", "--build", "server")
	compose("down")
}

// Start all containers
func (Docker) Compose() {
	compose := cmd("docker", "compose")

	compose("up")
}

// Start only aux container (db + webUI)
func (Docker) ComposeAux() {
	compose := cmd("docker", "compose")

	compose("down")
	compose("-f", "docker-compose-aux.yml", "up")
}

// Generates Dockerfile with ./cmd/app/main.go as an entry point
func (Docker) Generate() {
	if _, err := os.Stat("Dockerfile"); errors.Is(err, os.ErrNotExist) {
		cmd("go", "install", "github.com/zeromicro/go-zero/tools/goctl@latest")

		goctl := cmd("goctl")
		goctl("docker", "-go", filepath.Join(WD, "cmd", "app", "main.go"), "--tz", "Europe/Moscow")

		// applyGoctlPatch()
	}
}
