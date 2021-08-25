package ap

import (
	"bytes"
	"encoding/json"
	"github.com/MatticNote/MatticNote/worker"
	"github.com/go-fed/httpsig"
	"github.com/gocraft/work"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strings"
)

func inboxGet(c *fiber.Ctx) error {
	c.Status(fiber.StatusMethodNotAllowed)
	return nil
}

func inboxPost(w http.ResponseWriter, r *http.Request) {
	var err error
	if !strings.HasPrefix(
		strings.ToLower(r.Header.Get("Content-Type")),
		"application/activity+json") &&
		!strings.HasPrefix(
			strings.ToLower(r.Header.Get("Content-Type")),
			"application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"") {
		// Invalid header. return error.
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bufBody := new(bytes.Buffer)
	_, err = bufBody.ReadFrom(r.Body)
	defer func() {
		_ = r.Body.Close()
	}()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var inboxData map[string]interface{}
	if err := json.Unmarshal(bufBody.Bytes(), &inboxData); err != nil {
		// Not json. return error.
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	verifier, err := httpsig.NewVerifier(r)
	if err != nil {
		// not found HTTP signature. return error.
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, err = worker.Enqueue.Enqueue(worker.JobInboxProcess, work.Q{
		"data":      inboxData,
		"signature": verifier,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Queue Accepted
	w.WriteHeader(http.StatusAccepted)
}
