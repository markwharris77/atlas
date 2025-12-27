package config

import (
	"fmt"
	"os"

	"github.com/markwharris77/atlas/internal/log"
	"gopkg.in/yaml.v3"
)

type Config struct {
	App struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	} `yaml:"app"`

	Deploy struct {
		LocalDir string `yaml:"local-dir"`
		User     string `yaml:"user"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
	}
	Run struct {
		Command []string          `yaml:"command"`
		Env     map[string]string `yaml:"env"`
		User    string            `yaml:"user"`
	}
}

func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %q: %w", path, err)
	}

	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("parse yaml %q: %w", path, err)
	}

	return &c, nil
}

func Print(cfg *Config) error {
	out, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	log.Info(string(out))
	return nil
}
