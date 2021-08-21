package ap

import (
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/worker"
	"github.com/go-fed/httpsig"
	"github.com/gocraft/work"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func inboxGet(c *fiber.Ctx) error {
	c.Status(fiber.StatusMethodNotAllowed)
	return nil
}

func inboxPost(w http.ResponseWriter, r *http.Request) {
	verifier, err := httpsig.NewVerifier(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userPK, err := internal.GetUserPublicKey(verifier.KeyId())
	if err != nil {
		// Signature missing. ignore.
		return
	}
	err = verifier.Verify(userPK, httpsig.RSA_SHA256)
	if err != nil {
		// Invalid HTTP Signature. ignore.
		return
	}

	_, err = worker.Enqueue.Enqueue(worker.JobInboxProcess, work.Q{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
