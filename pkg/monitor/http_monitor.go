package monitor

import (
	"context"
	"time"
	"net/http"
	"fmt"
)

type CheckResult struct {
	Name		string
	Healthy		bool
	StatusCode	int
	Latency		time.Duration
	Err			error
	CheckedAt	time.Time
}

type HTTPMonitor struct {
	Name string
	URL string
	Client *http.Client
}

func (h *HTTPMonitor) Check(ctx context.Context) CheckResult {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.URL, nil)
	if err != nil {
		return CheckResult{Name: h.Name, Err: fmt.Errorf("building request: %w", err), CheckedAt: time.Now()}
	}

	start := time.Now()
	resp, err := h.Client.Do(req)
	if err != nil {
		return CheckResult{Name: h.Name, Err: fmt.Errorf("request failed: %w", err), Latency: time.Since(start), CheckedAt: time.Now()}
	}
	defer resp.Body.Close()

	return CheckResult{
		Name:		h.Name,
		StatusCode:	resp.StatusCode,
		Latency:	time.Since(start),
		Healthy:	resp.StatusCode >= 200 && resp.StatusCode < 400,
		CheckedAt: 	time.Now(),
	}
}
