package main

import (
	"flag"
	"log"
	"fmt"
	"time"
	"net/http"
	"context"

	"github.com/markvosstb/homelab_monitor/internal/config"
	"github.com/markvosstb/homelab_monitor/pkg/monitor"
)

func main() {
	configPath := flag.String("config", "cfg/ibn5100.yaml", "Path to server cfg file.")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, svcCfg := range cfg.Services {
		url := fmt.Sprintf("https://%s.%s", svcCfg.Subdomain, cfg.Site)
		client := &http.Client {
			Timeout: 10 * time.Second,
		}
		httpMonitor := &monitor.HTTPMonitor{Name: svcCfg.Name, URL: url, Client: client}
		ctx, cancel := context.WithTimeout(ctx, 10 * time.Second)
		defer cancel()

		result := httpMonitor.Check(ctx)
		log.Println(result)
	}
}
