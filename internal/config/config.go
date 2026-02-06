package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const configDir = "tics"
const configFile = "config.json"

type Config struct {
	GLPIURL      string        `json:"glpi_url"`
	UserToken    string        `json:"user_token"`
	AppToken     string        `json:"app_token"`
	FollowupHTML string        `json:"followup_html,omitempty"`
	SyncInterval time.Duration `json:"-"`
}

// GetFollowupHTML returns the custom followup HTML if set, otherwise the default.
func (c *Config) GetFollowupHTML() string {
	if c.FollowupHTML != "" {
		return c.FollowupHTML
	}
	return AutoFollowupHTML
}

var Templates = []string{
	"Bom dia! Estamos analisando seu chamado.",
	"Chamado recebido. Em breve retornaremos com uma solução.",
	"Estamos trabalhando na resolução do seu chamado.",
}

// AutoFollowupHTML is the HTML template sent automatically to tickets
// that are assigned to the user and in "processing" status,
// if no such followup has been sent yet.
const AutoFollowupHTML = `<div style="display: flex; flex-direction: column; align-items: center; justify-content: center; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, sans-serif; max-width: 600px; margin: 0 auto; padding: 24px; background-color: transparent; border-radius: 14px; border: 1px solid rgba(128, 128, 128, 0.2);">
<div style="background: linear-gradient(90deg, #0ea575 0%, #10b981 100%); color: white; font-size: 16px; font-weight: bold; padding: 12px 24px; border-radius: 10px; margin-bottom: 20px; width: 90%; text-align: center; box-shadow: 0 3px 10px rgba(16, 185, 129, 0.25); display: flex; align-items: center; justify-content: center; letter-spacing: 0.5px;">
<div style="display: inline-block; height: 12px; width: 12px; background-color: white; border-radius: 50%; margin-right: 12px; border: 2px solid rgba(255, 255, 255, 0.7);"> </div>
<span style="display: inline-block; transform: translateY(1px);">ATENDIMENTO ABERTO</span></div>
<div style="width: 100%; padding: 20px; background-color: rgba(128, 128, 128, 0.05); border: 1px solid rgba(128, 128, 128, 0.12); display: flex; flex-direction: column; justify-content: center; align-items: center; border-radius: 10px; margin: 0 0 20px 0;">
<div style="display: flex; flex-direction: column; width: 100%;">
<div style="display: flex; align-items: center; margin-bottom: 12px; padding-bottom: 12px; border-bottom: 1px solid rgba(128, 128, 128, 0.15);">
<div style="width: 40%; text-align: right; padding-right: 15px; font-weight: 600; font-size: 14px; color: inherit; opacity: 0.7;">ANALISTA</div>
<div style="width: 60%; font-weight: 500; color: inherit;">Paulo Felipe Borszcz</div>
</div>
<div style="display: flex; align-items: center; margin-bottom: 12px; padding-bottom: 12px; border-bottom: 1px solid rgba(128, 128, 128, 0.15);">
<div style="width: 40%; text-align: right; padding-right: 15px; font-weight: 600; font-size: 14px; color: inherit; opacity: 0.7;">SETOR</div>
<div style="width: 60%; font-weight: 500; color: inherit;">TI - HelpDesk</div>
</div>
<div style="display: flex; align-items: center;">
<div style="width: 40%; text-align: right; padding-right: 15px; font-weight: 600; font-size: 14px; color: inherit; opacity: 0.7;">WHATSAPP</div>
<div style="width: 60%;"><span style="font-weight: 500; color: inherit;">42 3309-7213</span></div>
</div>
</div>
</div>
<div style="padding: 15px; border-radius: 8px; width: 100%; display: flex; justify-content: center; align-items: center; background-color: rgba(128, 128, 128, 0.02); border: 1px solid rgba(128, 128, 128, 0.08);"><img style="height: 70px; max-width: 100%; object-fit: contain;" src="https://i.ibb.co/M1VYxMC/Type-Default.png" alt="Logo do sistema"></div>
</div>`

// AutoFollowupMarker is a substring used to detect if the auto-followup
// was already sent to a ticket.
const AutoFollowupMarker = "ATENDIMENTO ABERTO"

// configPath returns the full path to the config file.
func configPath() (string, error) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("getting config dir: %w", err)
	}
	return filepath.Join(cfgDir, configDir, configFile), nil
}

// IsConfigured returns true if there are saved tokens.
func (c *Config) IsConfigured() bool {
	return c.UserToken != ""
}

// Save writes the config to ~/.config/tics/config.json.
func (c *Config) Save() error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}
	return nil
}

// Load reads config from file, falling back to environment variables.
func Load() *Config {
	intervalSec := 30
	if v := os.Getenv("TICS_SYNC_INTERVAL"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			intervalSec = n
		}
	}

	cfg := &Config{
		GLPIURL:      "https://nexus.lojasmm.com.br/apirest.php",
		SyncInterval: time.Duration(intervalSec) * time.Second,
	}

	// Try loading from file first
	if path, err := configPath(); err == nil {
		if data, err := os.ReadFile(path); err == nil {
			if err := json.Unmarshal(data, cfg); err == nil && cfg.UserToken != "" {
				cfg.SyncInterval = time.Duration(intervalSec) * time.Second
				return cfg
			}
		}
	}

	// Fall back to environment variables
	if v := os.Getenv("GLPI_URL"); v != "" {
		cfg.GLPIURL = v
	}
	if v := os.Getenv("GLPI_USER_TOKEN"); v != "" {
		cfg.UserToken = v
	}
	if v := os.Getenv("GLPI_APP_TOKEN"); v != "" {
		cfg.AppToken = v
	}

	return cfg
}
