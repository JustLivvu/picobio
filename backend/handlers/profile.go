package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"backend/auth"
)

type Profile struct {
	DisplayName string `json:"displayName"`
	Pronouns    string `json:"pronouns"`
	Description string `json:"description"`
	JoinedDate  string `json:"joinedDate"`
	GithubURL   string `json:"githubUrl"`
	AvatarURL   string `json:"avatarUrl"`
}

type ProfileHandler struct {
	mu       sync.RWMutex
	filePath string
	mediaDir string
	sessions *auth.SessionManager
	profile  Profile
}

func NewProfileHandler(filePath string, mediaDir string, sessions *auth.SessionManager) *ProfileHandler {
	h := &ProfileHandler{
		filePath: filePath,
		mediaDir: mediaDir,
		sessions: sessions,
	}
	h.loadProfile()
	return h
}

func (h *ProfileHandler) loadProfile() {
	h.mu.Lock()
	defer h.mu.Unlock()

	data, err := os.ReadFile(h.filePath)
	if err != nil {
		h.profile = Profile{
			DisplayName: "Livvya",
			Pronouns:    "livvya · she/her",
			Description: "Hey! I'm Livvya, Developer and maintainer of Picobio",
			JoinedDate:  "Joined on Jun 14, 2026",
			GithubURL:   "https://github.com/justlivvu",
			AvatarURL:   "/media/avatar.jpg",
		}
		return
	}

	json.Unmarshal(data, &h.profile)
}

func (h *ProfileHandler) saveProfile() error {
	data, err := json.MarshalIndent(h.profile, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(h.filePath, data, 0644)
}

func (h *ProfileHandler) extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	return r.URL.Query().Get("token")
}

func (h *ProfileHandler) HandleProfile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.mu.RLock()
		defer h.mu.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(h.profile)

	case http.MethodPost, http.MethodPut:
		token := h.extractToken(r)
		if !h.sessions.ValidateSession(token) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized"})
			return
		}

		var updated Profile
		if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request body"})
			return
		}

		h.mu.Lock()
		if updated.DisplayName != "" {
			h.profile.DisplayName = updated.DisplayName
		}
		h.profile.Pronouns = updated.Pronouns
		h.profile.Description = updated.Description
		if updated.JoinedDate != "" {
			h.profile.JoinedDate = updated.JoinedDate
		}
		h.profile.GithubURL = updated.GithubURL
		if updated.AvatarURL != "" {
			h.profile.AvatarURL = updated.AvatarURL
		}
		err := h.saveProfile()
		h.mu.Unlock()

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Failed to save profile"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(h.profile)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ProfileHandler) HandleAvatarUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := h.extractToken(r)
	if !h.sessions.ValidateSession(token) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized"})
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "File too large or invalid form"})
		return
	}

	file, handler, err := r.FormFile("avatar")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Avatar file required"})
		return
	}
	defer file.Close()

	ext := filepath.Ext(handler.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	filename := fmt.Sprintf("avatar%s", ext)
	dstPath := filepath.Join(h.mediaDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to save file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to write file"})
		return
	}

	avatarURL := fmt.Sprintf("/media/%s?t=%d", filename, time.Now().Unix())

	h.mu.Lock()
	h.profile.AvatarURL = avatarURL
	h.saveProfile()
	h.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"avatarUrl": avatarURL,
	})
}
