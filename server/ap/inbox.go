package ap

import (
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

	// todo: keyIdで共有鍵を探す
	verifier.KeyId()
	// todo: keyIdが見つからない・署名エラーが出たらinbox処理をしない

	_, err = worker.Enqueue.Enqueue("inbox_worker", work.Q{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
