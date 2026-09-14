package cli

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/chenyme/grok2api/backend/internal/domain/account"
	"github.com/chenyme/grok2api/backend/internal/infra/provider"
	"github.com/chenyme/grok2api/backend/internal/infra/provider/conversation"
)

// Both turns take the normal auto route: Build 403 followed by XAI success.
// The second turn must restore XAI's own cached proof for its tool call.
func TestCoordinatorFallbackRestoresOwnConversationCache(t *testing.T) {
	for _, tc := range []struct {
		name, operation, first, second string
	}{
		{"chat", conversation.OperationChat,
			`{"messages":[{"role":"user","content":"hello"}]}`,
			`{"messages":[{"role":"user","content":"hello"},{"role":"assistant","content":null,"tool_calls":[{"id":"call_one","type":"function","function":{"name":"list_dir","arguments":"{}"}}]},{"role":"tool","tool_call_id":"call_one","content":"ok"}]}`},
		{"messages", conversation.OperationMessages,
			`{"model":"grok-4.6","max_tokens":128,"thinking":{"type":"enabled","budget_tokens":1024},"messages":[{"role":"user","content":"hello"}]}`,
			`{"model":"grok-4.6","max_tokens":128,"thinking":{"type":"enabled","budget_tokens":1024},"messages":[{"role":"user","content":"hello"},{"role":"assistant","content":[{"type":"tool_use","id":"toolu_call_one","name":"list_dir","input":{}}]},{"role":"user","content":[{"type":"tool_result","tool_use_id":"toolu_call_one","content":"ok"}]}]}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			adapter, encrypted := newFallbackTestAdapter(t)
			adapter.SetFallbackMarker(&fallbackMarkerStub{})
			primaryCalls, fallbackCalls := 0, 0
			var fallbackBody string
			adapter.http.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
				if strings.Contains(req.URL.Host, "primary.test") {
					primaryCalls++
					return jsonResponse(http.StatusForbidden, `{"error":{"message":"forbidden"}}`, req), nil
				}
				fallbackCalls++
				body, err := io.ReadAll(req.Body)
				if err != nil {
					return nil, err
				}
				fallbackBody = string(body)
				if fallbackCalls == 1 {
					return jsonResponse(http.StatusOK, `{"id":"resp_one","status":"completed","output":[{"id":"rs_one","type":"reasoning","status":"completed","encrypted_content":"synthetic-xai-only-proof"},{"id":"fc_one","type":"function_call","call_id":"call_one","name":"list_dir","arguments":"{}"}]}`, req), nil
				}
				return jsonResponse(http.StatusOK, `{"id":"resp_two","status":"completed","output":[]}`, req), nil
			})
			request := provider.ResponseResourceRequest{
				Credential: account.Credential{ID: 99, Provider: account.ProviderBuild, EncryptedAccessToken: encrypted, BuildSuperEntitled: true},
				Method:     http.MethodPost, Path: "/responses", Model: "grok-4.6", Operation: tc.operation,
				ReasoningReplayKey: "synthetic-fallback-session", NormalizeBody: true, Body: []byte(tc.first),
			}
			first, err := adapter.ForwardResponse(context.Background(), request)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := io.ReadAll(first.Body); err != nil {
				t.Fatal(err)
			}
			_ = first.Body.Close()
			request.Body = []byte(tc.second)
			// Positive control proves the first response populated XAI's cache.
			cached, _, err := conversation.ConvertRequestWithReasoningReplay(request.Body, request.Model, request.Operation, adapter.conversationReasoningCache, adapter.conversationReasoningScope(request, adapter.fallbackBaseURL()))
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(cached), "synthetic-xai-only-proof") {
				t.Fatal("setup failed: XAI conversation cache was not populated")
			}
			second, err := adapter.ForwardResponse(context.Background(), request)
			if err != nil {
				t.Fatal(err)
			}
			_ = second.Body.Close()
			if primaryCalls != 2 || fallbackCalls != 2 {
				t.Fatalf("unexpected calls: primary=%d fallback=%d", primaryCalls, fallbackCalls)
			}
			if !strings.Contains(fallbackBody, "synthetic-xai-only-proof") {
				t.Fatal("second XAI fallback omitted its own cached tool-call reasoning")
			}
		})
	}
}
