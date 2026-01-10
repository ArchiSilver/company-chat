package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"company-chat/internal/domain"
	"company-chat/internal/logging"
)

type mockSaver struct {
    saved *domain.Upload
}

func (m *mockSaver) Save(ctx context.Context, u *domain.Upload) error {
    m.saved = u
    // simulate DB assigning ID and created_at
    u.ID = "mock-id"
    u.CreatedAt = u.CreatedAt
    return nil
}

func TestUploadHandler_Success(t *testing.T) {
    logging.Init(false)
    defer logging.Sync()
    // prepare a small PNG-like payload (magic header) so DetectContentType sees image/png
    pngHeader := []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    part, err := writer.CreateFormFile("file", "test.png")
    if err != nil {
        t.Fatalf("create form file: %v", err)
    }
    part.Write(append(pngHeader, []byte("hello")...))
    writer.Close()

    req := httptest.NewRequest(http.MethodPost, "/api/uploads", body)
    req.Header.Set("Content-Type", writer.FormDataContentType())

    ms := &mockSaver{}
    // simple getUserID that returns a fixed user id
    getUserID := func(r *http.Request) (string, error) { return "user-123", nil }

    handler := NewUploadHandler(ms, 1024*1024, getUserID)
    rr := httptest.NewRecorder()
    handler.ServeHTTP(rr, req)

    if rr.Code != http.StatusCreated {
        t.Fatalf("expected status 201, got %d body=%s", rr.Code, rr.Body.String())
    }
    var resp map[string]string
    if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
        t.Fatalf("decode response: %v", err)
    }
    if resp["type"] != "image/png" {
        t.Fatalf("expected type image/png, got %s", resp["type"])
    }
    if resp["url"] == "" {
        t.Fatalf("expected url in response")
    }
    if ms.saved == nil {
        t.Fatalf("expected Save to be called")
    }
    if ms.saved.UserID != "user-123" {
        t.Fatalf("expected saved user id, got %s", ms.saved.UserID)
    }
    // cleanup uploaded file if any
    if ms.saved != nil && ms.saved.Path != "" {
        _ = os.Remove(strings.TrimPrefix(ms.saved.Path, "/uploads/"))
    }
}
