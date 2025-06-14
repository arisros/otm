package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"otm/internal/encryption"
	"otm/internal/middleware"
	"otm/internal/storage"
)

type readMessageResponse struct {
	Message string `json:"message"`
}

func handleReadMessage(db *storage.DBHandle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		msg, err := storage.LoadMessage(db.Conn, id)
		if err != nil {
			http.Error(w, "Message not found or expired", http.StatusNotFound)
			return
		}

		if time.Now().After(msg.ExpiresAt) {
			storage.DeleteMessage(db.Conn, id)
			http.Error(w, "Message expired", http.StatusGone)
			return
		}

		// Decrypt the message
		plaintext, err := encryption.DecryptMessageLayer(
			msg.EncryptedText, msg.NonceMsg,
			msg.EncryptedKey, msg.NonceKey,
		)
		if err != nil {
			http.Error(w, "Decryption failed", http.StatusInternalServerError)
			return
		}

		// If read_once, delete after reading
		if msg.ReadOnce {
			_ = storage.DeleteMessage(db.Conn, id)
		}

		// Log the read
		ip := middleware.GetIP(r)
		country := middleware.LookupCountry(ip)
		storage.LogRead(db.Conn, storage.LogEntry{
			MessageID: id,
			IPAddress: ip,
			Country:   country,
			Timestamp: time.Now(),
		})

		// Respond
		res := readMessageResponse{Message: string(plaintext)}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}
