package handlers

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"company-chat/internal/domain"
	"company-chat/internal/httputil"
	"company-chat/internal/logging"
	"company-chat/internal/repository"
)

// NewUploadHandler returns an http.HandlerFunc which handles multipart uploads.
// uploadSaver implements Save(ctx, *domain.Upload) error.
// getUserID extracts user id from request (can return error if unauthenticated).
func NewUploadHandler(uploadSaver repository.UploadSaver, maxUploadSize int64, getUserID func(*http.Request) (string, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			httputil.JSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize+512)
		if err := r.ParseMultipartForm(4 << 20); err != nil {
			httputil.JSONError(w, http.StatusBadRequest, "invalid multipart form or file too large")
			return
		}
		file, header, err := r.FormFile("file")
		if err != nil {
			httputil.JSONError(w, http.StatusBadRequest, "no file field")
			return
		}
		defer file.Close()

		sniff := make([]byte, 512)
		n, _ := file.Read(sniff)
		contentType := http.DetectContentType(sniff[:n])
		if !(strings.HasPrefix(contentType, "image/") || strings.HasPrefix(contentType, "audio/") || strings.HasPrefix(contentType, "video/") || contentType == "application/pdf" || contentType == "application/octet-stream") {
			httputil.JSONError(w, http.StatusUnsupportedMediaType, "unsupported media type")
			return
		}

		randb := make([]byte, 12)
		if _, err := rand.Read(randb); err != nil {
			httputil.JSONError(w, http.StatusInternalServerError, "could not generate filename")
			return
		}
		ext := ""
		if exts, _ := mime.ExtensionsByType(contentType); len(exts) > 0 {
			ext = exts[0]
		} else {
			ext = filepath.Ext(header.Filename)
		}
		if !(strings.HasPrefix(ext, ".") && len(ext) > 1) {
			ext = ".bin"
		} else {
			ok := true
			for _, r := range ext[1:] {
				if !(('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') || ('0' <= r && r <= '9')) {
					ok = false
					break
				}
			}
			if !ok {
				ext = ".bin"
			}
		}
		fname := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), hex.EncodeToString(randb), ext)
		fname = filepath.Base(fname)
		outPath := filepath.Join("uploads", fname)

		if err := os.MkdirAll("uploads", 0755); err != nil {
			httputil.JSONError(w, http.StatusInternalServerError, "could not create uploads dir")
			return
		}
		out, err := os.OpenFile(outPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err != nil {
			logging.L.Warnf("could not open upload file %s: %v", outPath, err)
			httputil.JSONError(w, http.StatusInternalServerError, "could not save file")
			return
		}
		defer out.Close()
		reader := io.MultiReader(bytes.NewReader(sniff[:n]), file)
		written, err := io.Copy(out, io.LimitReader(reader, maxUploadSize))
		if err != nil {
			logging.L.Warnf("could not write upload to %s: %v", outPath, err)
			httputil.JSONError(w, http.StatusInternalServerError, "could not write file")
			return
		}
		if err := out.Sync(); err != nil {
			logging.L.Warnf("failed to sync upload file %s: %v", outPath, err)
			httputil.JSONError(w, http.StatusInternalServerError, "could not save file")
			return
		}
		if written == 0 {
			logging.L.Warnf("empty upload received for path %s", outPath)
			httputil.JSONError(w, http.StatusBadRequest, "empty file")
			return
		}
		urlPath := "/uploads/" + fname
		uid, err := getUserID(r)
		if err != nil {
			logging.L.Debugf("upload without auth: %v", err)
		} else {
			u := &domain.Upload{UserID: uid, Path: urlPath, MIME: contentType, Size: written}
			if err := uploadSaver.Save(r.Context(), u); err != nil {
				logging.L.Warnf("could not save upload metadata: %v", err)
			} else {
				logging.L.Infof("upload saved: id=%s user=%s path=%s size=%d", u.ID, u.UserID, u.Path, u.Size)
			}
		}
		resp := map[string]string{"url": urlPath, "type": contentType}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(resp)
	}
}
