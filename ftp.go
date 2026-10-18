package main

import (
	"context"
	"net/netip"
	"sync"
	"time"

	"github.com/shizukutanaka/Mimamori/internal/ftp"
	"github.com/shizukutanaka/Mimamori/internal/i18n"
	"github.com/shizukutanaka/Mimamori/internal/probe"
	"github.com/shizukutanaka/Mimamori/internal/report"
)

func anonymousFTPHosts(ctx context.Context, opened []probe.Open) map[netip.Addr]bool {
	if ctx == nil {
		return nil
	}
	hosts := make([]probe.Open, 0)
	for _, candidate := range opened {
		if candidate.Port == 21 {
			hosts = append(hosts, candidate)
		}
	}
	if len(hosts) == 0 {
		return nil
	}

	workers := 16
	if len(hosts) < workers {
		workers = len(hosts)
	}
	jobs := make(chan probe.Open)
	results := make(chan netip.Addr, len(hosts))
	var wait sync.WaitGroup
	wait.Add(workers)
	for range workers {
		go func() {
			defer wait.Done()
			for candidate := range jobs {
				if ftp.Anonymous(ctx, candidate.IP, candidate.Port, 2*time.Second) {
					results <- candidate.IP
				}
			}
		}()
	}
	go func() {
		defer close(jobs)
		for _, candidate := range hosts {
			select {
			case jobs <- candidate:
			case <-ctx.Done():
				return
			}
		}
	}()

	wait.Wait()
	close(results)
	anonymous := make(map[netip.Addr]bool)
	for host := range results {
		anonymous[host] = true
	}
	return anonymous
}

func markAnonymousFTP(findings []report.Finding, anonymous map[netip.Addr]bool, lang string) []report.Finding {
	for i := range findings {
		finding := &findings[i]
		if finding.ID != "ftp_open" {
			continue
		}
		ip, err := netip.ParseAddr(finding.IP)
		if err != nil || !anonymous[ip] {
			continue
		}
		finding.ID = "ftp_anonymous"
		finding.Level = report.Danger
		finding.Points = 25
		finding.Message = i18n.T(lang, "ftp_anonymous", finding.Device)
		finding.Fix = i18n.T(lang, "ftp_anonymous_fix")
	}
	return findings
}
