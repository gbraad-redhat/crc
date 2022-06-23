package mount

import (
	"net/http"

	"github.com/gbraad/go9p"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{
	}
}

func (h *Handler) Mount(c *context) error {
	go9p.StartServer(":5460", 0, "c:")

	return c.Code(http.StatusOK)
}

func (h *Handler) Umount(c *context) error {


	return c.Code(http.StatusOK)
}
