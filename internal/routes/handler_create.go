package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"otm/internal/encryption"
	"otm/internal/middleware"
	"otm/internal/storage"
)

type createMessageRequest struct {
	Message   string `json:"message"`
	ReadOnce  bool   `json:"read_once"`
	ExpiresIn int    `json:"expires_in"` // seconds
}

type createMessageResponse struct {
	ID string `json:"id"`
}

func handleCreateMessage(db *storage.DBHandle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		id := uuid.NewString()
		now := time.Now()
		expires := now.Add(time.Duration(req.ExpiresIn) * time.Second)

		// Encrypt the message
		cipherText, nonceMsg, encryptedKey, nonceKey, err := encryption.EncryptMessageLayer([]byte(req.Message))
		if err != nil {
			http.Error(w, "Encryption failed", http.StatusInternalServerError)
			return
		}

		// Save to database
		msg := storage.Message{
			ID:            id,
			EncryptedText: cipherText,
			NonceMsg:      nonceMsg,
			EncryptedKey:  encryptedKey,
			NonceKey:      nonceKey,
			ReadOnce:      req.ReadOnce,
			CreatedAt:     now,
			ExpiresAt:     expires,
		}

		if err := storage.SaveMessage(db.Conn, msg); err != nil {
			http.Error(w, "DB write error", http.StatusInternalServerError)
			return
		}

		// Log the write
		ip := middleware.GetIP(r)
		country := middleware.LookupCountry(ip)
		storage.LogWrite(db.Conn, storage.LogEntry{
			MessageID: id,
			IPAddress: ip,
			Country:   country,
			Timestamp: now,
		})

		res := createMessageResponse{ID: id}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}
