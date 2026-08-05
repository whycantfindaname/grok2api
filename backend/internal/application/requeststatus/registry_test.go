package requeststatus

import (
	"testing"
	"time"
)

func TestRegistryTracksRunningAndTerminalStatePerClientKey(t *testing.T) {
	registry := NewRegistry(10 * time.Minute)
	startedAt := time.Now().UTC()
	registry.Start(1, "request-1", startedAt)

	status, ok := registry.Get(1, "request-1", startedAt.Add(time.Second))
	if !ok || status.State != StateRunning || !status.StartedAt.Equal(startedAt) {
		t.Fatalf("running status = %#v, ok=%v", status, ok)
	}
	if _, ok := registry.Get(2, "request-1", startedAt.Add(time.Second)); ok {
		t.Fatal("request status leaked across client keys")
	}

	finishedAt := startedAt.Add(2 * time.Second)
	registry.Finish(1, "request-1", 502, finishedAt)
	status, ok = registry.Get(1, "request-1", finishedAt)
	if !ok || status.State != StateFailed || status.FinishedAt == nil || status.StatusCode != 502 {
		t.Fatalf("terminal status = %#v, ok=%v", status, ok)
	}
	if _, ok := registry.Get(1, "request-1", finishedAt.Add(11*time.Minute)); ok {
		t.Fatal("expired terminal status was retained")
	}
}

func TestRegistryKeepsDuplicateRequestIDRunningUntilAllRequestsFinish(t *testing.T) {
	registry := NewRegistry(time.Minute)
	now := time.Now().UTC()
	registry.Start(1, "duplicate", now)
	registry.Start(1, "duplicate", now.Add(time.Millisecond))
	registry.Finish(1, "duplicate", 200, now.Add(time.Second))

	status, ok := registry.Get(1, "duplicate", now.Add(2*time.Second))
	if !ok || status.State != StateRunning {
		t.Fatalf("status after first finish = %#v, ok=%v", status, ok)
	}
	registry.Finish(1, "duplicate", 200, now.Add(3*time.Second))
	status, ok = registry.Get(1, "duplicate", now.Add(3*time.Second))
	if !ok || status.State != StateCompleted {
		t.Fatalf("status after final finish = %#v, ok=%v", status, ok)
	}
}
