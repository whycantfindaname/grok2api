package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"

	"github.com/chenyme/grok2api/backend/internal/domain/account"
	"github.com/chenyme/grok2api/backend/internal/infra/provider"
	"github.com/chenyme/grok2api/backend/internal/infra/provider/conversation"
)

// Synthetic request data only. The project logging contract prohibits bodies.
func TestCoordinatorToolSchemaBodyIsAbsentFromInfoLogs(t *testing.T) {
	var logs bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&logs, nil)))
	t.Cleanup(func() { slog.SetDefault(previous) })
	const marker = "synthetic-private-tool-description"
	payload := map[string]json.RawMessage{
		"tools": json.RawMessage(`[{"type":"function","name":"automation_update","description":"synthetic-private-tool-description","parameters":{"type":"object","properties":{}}}]`),
	}
	if _, err := normalizeResponsesTools(payload); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(logs.String(), marker) {
		t.Fatal("default INFO logger emitted the synthetic tool-description body")
	}
}

func TestCoordinatorChatReasoningDoesNotCrossFallbackPlane(t *testing.T) {
	adapter, encrypted := newFallbackTestAdapter(t)
	adapter.SetFallbackMarker(&fallbackMarkerStub{})
	primaryCalls, fallbackCalls := 0, 0
	var fallbackBody string
	adapter.http.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if strings.Contains(req.URL.Host, "primary.test") {
			primaryCalls++
			if primaryCalls == 1 {
				return jsonResponse(http.StatusOK, `{"id":"resp_one","status":"completed","output":[{"id":"rs_one","type":"reasoning","status":"completed","encrypted_content":"synthetic-build-only-proof"},{"id":"fc_one","type":"function_call","call_id":"call_one","name":"list_dir","arguments":"{}"}]}`, req), nil
			}
			return jsonResponse(http.StatusForbidden, `{"error":{"message":"forbidden"}}`, req), nil
		}
		fallbackCalls++
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		fallbackBody = string(body)
		return jsonResponse(http.StatusOK, `{"id":"resp_two","status":"completed","output":[]}`, req), nil
	})
	request := provider.ResponseResourceRequest{
		Credential: account.Credential{ID: 99, Provider: account.ProviderBuild, EncryptedAccessToken: encrypted, BuildSuperEntitled: true},
		Method:     http.MethodPost, Path: "/responses", Model: "grok-4.6", Operation: conversation.OperationChat,
		ReasoningReplayKey: "synthetic-client-session", NormalizeBody: true,
		Body: []byte(`{"messages":[{"role":"user","content":"hello"}]}`),
	}
	first, err := adapter.ForwardResponse(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.ReadAll(first.Body); err != nil {
		t.Fatal(err)
	}
	_ = first.Body.Close()
	request.Body = []byte(`{"messages":[{"role":"user","content":"hello"},{"role":"assistant","content":null,"tool_calls":[{"id":"call_one","type":"function","function":{"name":"list_dir","arguments":"{}"}}]},{"role":"tool","tool_call_id":"call_one","content":"ok"}]}`)
	second, err := adapter.ForwardResponse(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	_ = second.Body.Close()
	if primaryCalls != 2 || fallbackCalls != 1 {
		t.Fatalf("unexpected calls: primary=%d fallback=%d", primaryCalls, fallbackCalls)
	}
	if strings.Contains(fallbackBody, "synthetic-build-only-proof") {
		t.Fatal("Build-scoped cached reasoning was sent to the XAI fallback plane")
	}
}

func TestCoordinatorMessagesReasoningDoesNotCrossFallbackPlane(t *testing.T) {
	adapter, encrypted := newFallbackTestAdapter(t)
	adapter.SetFallbackMarker(&fallbackMarkerStub{})
	primaryCalls, fallbackCalls := 0, 0
	var fallbackBody string
	adapter.http.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if strings.Contains(req.URL.Host, "primary.test") {
			primaryCalls++
			if primaryCalls == 1 {
				return jsonResponse(http.StatusOK, `{"id":"resp_one","status":"completed","output":[{"id":"rs_one","type":"reasoning","status":"completed","encrypted_content":"synthetic-build-only-proof"},{"id":"fc_one","type":"function_call","call_id":"call_one","name":"list_dir","arguments":"{}"}]}`, req), nil
			}
			return jsonResponse(http.StatusForbidden, `{"error":{"message":"forbidden"}}`, req), nil
		}
		fallbackCalls++
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		fallbackBody = string(body)
		return jsonResponse(http.StatusOK, `{"id":"resp_two","status":"completed","output":[]}`, req), nil
	})
	request := provider.ResponseResourceRequest{
		Credential: account.Credential{ID: 100, Provider: account.ProviderBuild, EncryptedAccessToken: encrypted, BuildSuperEntitled: true},
		Method:     http.MethodPost, Path: "/responses", Model: "claude-3-7-sonnet", Operation: conversation.OperationMessages,
		ReasoningReplayKey: "synthetic-client-session-messages", NormalizeBody: true,
		Body: []byte(`{"model":"claude-3-7-sonnet","max_tokens":128,"thinking":{"type":"enabled","budget_tokens":1024},"messages":[{"role":"user","content":"hello"}]}`),
	}
	first, err := adapter.ForwardResponse(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.ReadAll(first.Body); err != nil {
		t.Fatal(err)
	}
	_ = first.Body.Close()
	request.Body = []byte(`{"model":"claude-3-7-sonnet","max_tokens":128,"thinking":{"type":"enabled","budget_tokens":1024},"messages":[{"role":"user","content":"hello"},{"role":"assistant","content":[{"type":"tool_use","id":"toolu_call_one","name":"list_dir","input":{}}]},{"role":"user","content":[{"type":"tool_result","tool_use_id":"toolu_call_one","content":"ok"}]}]}`)
	second, err := adapter.ForwardResponse(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	_ = second.Body.Close()
	if primaryCalls != 2 || fallbackCalls != 1 {
		t.Fatalf("unexpected calls: primary=%d fallback=%d", primaryCalls, fallbackCalls)
	}
	if strings.Contains(fallbackBody, "synthetic-build-only-proof") {
		t.Fatal("Build-scoped cached reasoning was sent to the XAI fallback plane for Messages")
	}
}
