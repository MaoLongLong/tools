package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

var (
	cfgFile = flag.StringP("config", "c", "", "config file path")
	global  = flag.BoolP("global", "g", false, "modify global config")
)

type config struct {
	Users map[string]user
}

type user struct {
	Name  string
	Email string
}

func main() {
	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: gsu [options] <user>")
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
	}

	cfg, err := loadCfg()
	if err != nil {
		log.Fatal(err)
	}

	id := flag.Arg(0)
	u, ok := cfg.Users[id]
	if !ok {
		log.Fatal("gsu: user does not exist")
	}

	gitConfig := []string{"git", "config"}
	if *global {
		gitConfig = append(gitConfig, "--global")
	}

	if err := multiExec([][]string{
		append(gitConfig, "user.name", u.Name),
		append(gitConfig, "user.email", u.Email),
	}); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("user.name=%s\nuser.email=%s\n", u.Name, u.Email)
}

func multiExec(commands [][]string) error {
	for _, a := range commands {
		cmd := exec.Command(a[0], a[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func loadCfg() (*config, error) {
	name := *cfgFile
	if name == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		name = filepath.Join(home, ".gsu.yml")
	}
	b, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	var cfg config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
