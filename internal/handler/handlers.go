package handler

import (
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
)

type URLShortener struct {
	urls map[string]string
}

func NewUrlShortener() *URLShortener {
	return &URLShortener{
		urls: make(map[string]string),
	}
}

func (us URLShortener) generateShortCode() (string, error) {
	bytes := make([]byte, 3)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (us URLShortener) handleShortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "text/plain" {
		http.Error(w, "Invalid content type", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	originalURL := strings.TrimSpace(string(body))
	if originalURL == "" {
		http.Error(w, "URL cannot be empty", http.StatusBadRequest)
		return
	}

	shortCode, err := us.generateShortCode()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	shortURL := fmt.Sprintf("http://localhost:8080/%s", shortCode)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func (us URLShortener) handleRedirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	if path == "/" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	shortID := path[1:]

	originalURL, exists := us.urls[shortID]
	if !exists {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (us URLShortener) MainHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch {
	case path == "/" && r.Method == http.MethodPost:
		us.handleShortenURL(w, r)
	case path == "/" && r.Method != http.MethodPost:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	case path != "/" && r.Method == http.MethodGet:
		us.handleRedirect(w, r)
	default:
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
}
