package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"backend/auth"
)

type Post struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	CreatedAt string   `json:"createdAt"`
	Tags      []string `json:"tags,omitempty"`
}

type PostsHandler struct {
	mu       sync.RWMutex
	filePath string
	sessions *auth.SessionManager
	posts    []Post
}

func NewPostsHandler(filePath string, sessions *auth.SessionManager) *PostsHandler {
	h := &PostsHandler{
		filePath: filePath,
		sessions: sessions,
	}
	h.loadPosts()
	return h
}

func (h *PostsHandler) loadPosts() {
	h.mu.Lock()
	defer h.mu.Unlock()

	data, err := os.ReadFile(h.filePath)
	if err != nil {
		h.posts = []Post{}
		return
	}

	json.Unmarshal(data, &h.posts)
}

func (h *PostsHandler) savePosts() error {
	data, err := json.MarshalIndent(h.posts, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(h.filePath, data, 0644)
}

func (h *PostsHandler) extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	return r.URL.Query().Get("token")
}

func (h *PostsHandler) HandlePosts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.mu.RLock()
		defer h.mu.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(h.posts)

	case http.MethodPost:
		token := h.extractToken(r)
		if !h.sessions.ValidateSession(token) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized"})
			return
		}

		var newPost Post
		if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request body"})
			return
		}

		if strings.TrimSpace(newPost.Title) == "" && strings.TrimSpace(newPost.Content) == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Post title or content required"})
			return
		}

		b := make([]byte, 8)
		rand.Read(b)
		newPost.ID = hex.EncodeToString(b)
		newPost.CreatedAt = time.Now().Format("Jan 02, 2006 15:04")

		h.mu.Lock()
		h.posts = append([]Post{newPost}, h.posts...)
		err := h.savePosts()
		h.mu.Unlock()

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Failed to save post"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newPost)

	case http.MethodDelete:
		token := h.extractToken(r)
		if !h.sessions.ValidateSession(token) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized"})
			return
		}

		id := r.URL.Query().Get("id")
		if id == "" {
			var body struct {
				ID string `json:"id"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err == nil {
				id = body.ID
			}
		}

		if id == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Post ID required"})
			return
		}

		h.mu.Lock()
		found := false
		filtered := make([]Post, 0, len(h.posts))
		for _, p := range h.posts {
			if p.ID == id {
				found = true
			} else {
				filtered = append(filtered, p)
			}
		}
		if found {
			h.posts = filtered
			h.savePosts()
		}
		h.mu.Unlock()

		if !found {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"message": "Post not found"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]bool{"success": true})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
