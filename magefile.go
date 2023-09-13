//go:build mage
// +build mage

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/metafates/smartway-test/config"
	"github.com/metafates/smartway-test/pkg/postgres"
	"golang.org/x/sync/errgroup"
)

var WD = must(os.Getwd())

var Default = Full

func must[T any](value T, err error) T {
	if err != nil {
		log.Fatal(err)
	}

	return value
}

func run(name string, args ...string) error {
	return cmd(name, args...)()
}

func cmd(name string, args ...string) func(...string) error {
	cmd := exec.Command(name, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return func(args ...string) error {
		cmd.Args = append(cmd.Args, args...)

		return cmd.Run()
	}
}

// Start the server
func Run() error {
	return run("go", "run", filepath.Join(WD, "cmd", "app"))
}

type Docker mg.Namespace

// Rebuilds the Dockerfile
func (Docker) Rebuild() error {
	compose := cmd("docker", "compose")

	if err := compose("up", "-d", "--no-deps", "--build", "server"); err != nil {
		return err
	}

	if err := compose("down"); err != nil {
		return err
	}

	return nil
}

// Spin up docker containers, apply migrations and load example data
func Full(ctx context.Context) error {
	if err := cmd("docker", "compose", "down")(); err != nil {
		return err
	}

	var group *errgroup.Group
	group, ctx = errgroup.WithContext(ctx)

	group.Go(func() error {
		return cmd("docker", "compose", "up")()
	})

	group.Go(func() error {
		time.Sleep(5 * time.Second)
		fmt.Println("Loading migrations with example data")

		fmt.Println("Applying migrations")
		if err := migrate("up"); err != nil {
			return err
		}
		fmt.Println("Migrations applied")

		fmt.Println("Loading example data")
		if err := loadExample(ctx); err != nil {
			return err
		}
		fmt.Println("Example data is loaded")

		return nil
	})

	return group.Wait()
}

// Start all containers
func (Docker) Compose() error {
	compose := cmd("docker", "compose")

	return compose("up")
}

// Start only aux container (db + webUI)
func (Docker) ComposeAux() error {
	compose := cmd("docker", "compose")

	return compose("-f", "docker-compose-aux.yml", "up")
}

// Generates Dockerfile with ./cmd/app/main.go as an entry point
func (Docker) Generate() error {
	if _, err := os.Stat("Dockerfile"); errors.Is(err, os.ErrNotExist) {
		mg.Deps(BinDeps)

		goctl := cmd("goctl")
		return goctl("docker", "-go", filepath.Join(WD, "cmd", "app", "main.go"), "--tz", "Europe/Moscow")
	}

	return nil
}

type DB mg.Namespace

// Loads example data from example-data.sql.
func (DB) LoadExample(ctx context.Context) error {
	return loadExample(ctx)
}

func loadExample(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	pg, err := postgres.New(cfg.Postgres.URL, postgres.MaxPoolSize(cfg.Postgres.PoolMax))
	if err != nil {
		return err
	}
	defer pg.Close()

	sql, err := os.ReadFile("example-data.sql")
	if err != nil {
		return err
	}

	_, err = pg.Pool.Exec(ctx, string(sql))
	return err
}

func migrate(subcommand string) error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	goose := cmd(
		"goose",
		"-dir",
		filepath.Join(WD, "migrations"),
		"postgres",
		cfg.Postgres.URL,
		subcommand,
	)

	goose("status")
	return nil
}

// Migrate the DB to the most recent version available
func (DB) MigrateUp() error {
	return migrate("up")
}

// Dump the migration status for the current DB
func (DB) MigrateStatus() error {
	return migrate("status")
}

// Roll back the version by 1
func (DB) MigrateDown() error {
	return migrate("down")
}

// Install binary dependencies using go install
func BinDeps() {
	install := cmd("go", "install")

	install("github.com/zeromicro/go-zero/tools/goctl@latest")
}
