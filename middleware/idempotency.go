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

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

var idempoTracer = otel.Tracer("middleware/idempotency")

// IdempotencyMiddleware obezbeđuje idempotent operacije koristeći Consul kao storage
func IdempotencyMiddleware(consulClient *api.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Start span odmah, da pokrije cijeli middleware flow
			ctx := r.Context()
			ctx, span := idempoTracer.Start(ctx, "IdempotencyMiddleware")
			defer span.End()

			idempotencyKey := r.Header.Get("Idempotency-Key")
			span.SetAttributes(attribute.String("idempotency.key", idempotencyKey))

			// Ako nema ključa, samo nastavi
			if idempotencyKey == "" {
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}

			kv := consulClient.KV()
			keyPath := fmt.Sprintf("idempotency/%s", idempotencyKey)
			span.SetAttributes(attribute.String("consul.key", keyPath))

			// 1) Proveri da li ključ već postoji (consul get)
			var pair *api.KVPair
			{
				_, s := idempoTracer.Start(ctx, "consul.kv.get")
				s.SetAttributes(attribute.String("consul.key", keyPath))
				var err error
				pair, _, err = kv.Get(keyPath, nil)
				if err != nil {
					s.RecordError(err)
					s.SetStatus(codes.Error, "consul get failed")
					s.End()
					http.Error(w, "Failed to connect to Consul", http.StatusInternalServerError)
					return
				}
				s.End()
			}

			if pair != nil {
				var record model.IdempotencyRecord
				if err := json.Unmarshal(pair.Value, &record); err != nil {
					span.RecordError(err)
					span.SetStatus(codes.Error, "failed to unmarshal stored record")
					http.Error(w, "Failed to parse stored record", http.StatusInternalServerError)
					return
				}

				// Ako je završen, vrati keširani odgovor
				if record.Status == model.StatusCompleted {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(record.StatusCode)
					_, _ = w.Write([]byte(record.Body))
					return
				}

				// Ako je u toku, odbaci novi zahtev
				if record.Status == model.StatusInProgress {
					http.Error(w, "Request with this idempotency key is already in progress.", http.StatusConflict)
					return
				}
			}

			// 2) Pročitaj body i resetuj r.Body da se može ponovo čitati
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, "failed to read request body")
				http.Error(w, "Failed to read request body", http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

			// 3) Kreiraj placeholder za zahtev u toku (consul cas)
			placeholder := &model.IdempotencyRecord{Status: model.StatusInProgress}
			placeholderJSON, _ := json.Marshal(placeholder)
			p := &api.KVPair{Key: keyPath, Value: placeholderJSON, CreateIndex: 0}

			{
				_, s := idempoTracer.Start(ctx, "consul.kv.cas")
				s.SetAttributes(attribute.String("consul.key", keyPath))
				success, _, err := kv.CAS(p, nil)
				if err != nil {
					s.RecordError(err)
					s.SetStatus(codes.Error, "consul cas failed")
					s.End()
					http.Error(w, "Failed to write to Consul", http.StatusInternalServerError)
					return
				}
				if !success {
					s.SetStatus(codes.Error, "cas not successful (concurrent request)")
					s.End()
					http.Error(w, "A concurrent request with the same idempotency key is in progress.", http.StatusConflict)
					return
				}
				s.End()
			}

			// Ako handler panikuje, obriši key da se ne zaglavi "in_progress"
			defer func() {
				if rec := recover(); rec != nil {
					{
						_, s := idempoTracer.Start(ctx, "consul.kv.delete")
						s.SetAttributes(attribute.String("consul.key", keyPath))
						_, _ = kv.Delete(keyPath, nil)
						s.End()
					}
					panic(rec)
				}
			}()

			// 4) Pozovi sledeći handler ali sa istim ctx (da tracing nastavi)
			r = r.WithContext(ctx)

			rec := httptest.NewRecorder()
			next.ServeHTTP(rec, r)

			// 5) Sačuvaj finalni odgovor ako je uspešan (consul put) / ili obriši key (delete)
			if rec.Code >= 200 && rec.Code < 300 {
				finalRecord := model.IdempotencyRecord{
					Status:     model.StatusCompleted,
					StatusCode: rec.Code,
					Body:       rec.Body.String(),
				}
				finalJSON, _ := json.Marshal(finalRecord)
				finalPair := &api.KVPair{Key: keyPath, Value: finalJSON}

				{
					_, s := idempoTracer.Start(ctx, "consul.kv.put")
					s.SetAttributes(attribute.String("consul.key", keyPath))
					if _, err := kv.Put(finalPair, nil); err != nil {
						s.RecordError(err)
						s.SetStatus(codes.Error, "consul put failed")
						s.End()
						log.Printf("ERROR: Failed to save final response for key '%s': %v", idempotencyKey, err)
					} else {
						s.End()
						log.Printf("Saved final response for key '%s'", idempotencyKey)
					}
				}
			} else {
				{
					_, s := idempoTracer.Start(ctx, "consul.kv.delete")
					s.SetAttributes(attribute.String("consul.key", keyPath))
					_, _ = kv.Delete(keyPath, nil)
					s.End()
				}
			}

			// 6) Vrati odgovor klijentu
			for k, v := range rec.Header() {
				w.Header()[k] = v
			}
			w.WriteHeader(rec.Code)
			_, _ = io.Copy(w, rec.Body)
		})
	}
}
