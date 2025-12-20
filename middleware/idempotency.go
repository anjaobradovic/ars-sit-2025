package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/anjaobradovic/ars-sit-2025/model"
	"github.com/hashicorp/consul/api"
)

// IdempotencyMiddleware obezbeđuje idempotent operacije koristeći Consul kao storage
func IdempotencyMiddleware(consulClient *api.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			idempotencyKey := r.Header.Get("Idempotency-Key")

			// Ako nema ključa, samo nastavi sa sledećim handler-om
			if idempotencyKey == "" {
				next.ServeHTTP(w, r)
				return
			}

			kv := consulClient.KV()
			keyPath := fmt.Sprintf("idempotency/%s", idempotencyKey)

			// Proveri da li ključ već postoji
			pair, _, err := kv.Get(keyPath, nil)
			if err != nil {
				http.Error(w, "Failed to connect to Consul", http.StatusInternalServerError)
				return
			}

			if pair != nil {
				var record model.IdempotencyRecord
				if err := json.Unmarshal(pair.Value, &record); err != nil {
					http.Error(w, "Failed to parse stored record", http.StatusInternalServerError)
					return
				}

				// Ako je zahtev završen, vrati keširani odgovor
				if record.Status == model.StatusCompleted {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(record.StatusCode)
					w.Write([]byte(record.Body))
					return
				}

				// Ako je zahtev u toku, odbaci novi zahtev
				if record.Status == model.StatusInProgress {
					http.Error(w, "Request with this idempotency key is already in progress.", http.StatusConflict)
					return
				}
			}

			// Pročitaj body i resetuj r.Body da se može ponovo čitati
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read request body", http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

			// Kreiraj placeholder za zahtev u toku
			placeholder := &model.IdempotencyRecord{Status: model.StatusInProgress}
			placeholderJSON, _ := json.Marshal(placeholder)
			p := &api.KVPair{Key: keyPath, Value: placeholderJSON, CreateIndex: 0}
			success, _, err := kv.CAS(p, nil)
			if err != nil {
				http.Error(w, "Failed to write to Consul", http.StatusInternalServerError)
				return
			}
			if !success {
				http.Error(w, "A concurrent request with the same idempotency key is in progress.", http.StatusConflict)
				return
			}

			// Obradi stvarni zahtev
			defer func() {
				if r := recover(); r != nil {
					kv.Delete(keyPath, nil)
					panic(r)
				}
			}()

			rec := httptest.NewRecorder()
			next.ServeHTTP(rec, r)

			// Sačuvaj finalni odgovor ako je uspešan
			if rec.Code >= 200 && rec.Code < 300 {
				finalRecord := model.IdempotencyRecord{
					Status:     model.StatusCompleted,
					StatusCode: rec.Code,
					Body:       rec.Body.String(),
				}
				finalJSON, _ := json.Marshal(finalRecord)

				finalPair := &api.KVPair{Key: keyPath, Value: finalJSON}
				if _, err := kv.Put(finalPair, nil); err != nil {
					log.Printf("ERROR: Failed to save final response for key '%s': %v", idempotencyKey, err)
				} else {
					log.Printf("Saved final response for key '%s'", idempotencyKey)
				}
			} else {
				// Ako nije uspešno, ukloni ključ da dozvoli retry
				kv.Delete(keyPath, nil)
			}

			// Vrati odgovor klijentu
			for k, v := range rec.Header() {
				w.Header()[k] = v
			}
			w.WriteHeader(rec.Code)
			io.Copy(w, rec.Body)
		})
	}
}
