package mount

import (
	"net/http"
	"net/url"
	"sync"
	"io/ioutil"

	"github.com/code-ready/crc/pkg/crc/logging"
)

type context struct {
	method      string
	requestBody []byte
	url         *url.URL

	code         int
	headers      map[string]string
	responseBody []byte
}

func (c *context) Code(code int) error {
	c.code = code
	return nil
}

func NewMux() http.Handler {
	handler := NewHandler()

	server := newServerWithMounts(handler)

	return server.Handler()
}

type server struct {
	mounts     map[string]map[string]func(*context) error
	mountsLock sync.RWMutex
}

func newServer() *server {
	return &server{
		mounts: make(map[string]map[string]func(*context) error),
	}
}

func newServerWithMounts(handler *Handler) *server {
	server := newServer()

	server.POST("/mount", handler.Mount)
	server.GET("/mount", handler.Mount)

	server.POST("/umount", handler.Umount)
	server.GET("/umount", handler.Umount)

	return server
}


func (s *server) GET(pattern string, handler func(c *context) error) {
	s.mountsLock.Lock()
	defer s.mountsLock.Unlock()
	if _, ok := s.mounts[pattern]; !ok {
		s.mounts[pattern] = make(map[string]func(*context) error)
	}
	s.mounts[pattern][http.MethodGet] = handler
}

func (s *server) POST(pattern string, handler func(c *context) error) {
	s.mountsLock.Lock()
	defer s.mountsLock.Unlock()
	if _, ok := s.mounts[pattern]; !ok {
		s.mounts[pattern] = make(map[string]func(*context) error)
	}
	s.mounts[pattern][http.MethodPost] = handler
}

func (s *server) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.mountsLock.RLock()
		mount, ok := s.mounts[r.URL.Path]
		if !ok {
			s.mountsLock.RUnlock()
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		handler, ok := mount[r.Method]
		if !ok {
			s.mountsLock.RUnlock()
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		s.mountsLock.RUnlock()

		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		c := &context{
			method:      r.Method,
			requestBody: requestBody,
			headers:     make(map[string]string),
			url:         r.URL,
		}
		if err := handler(c); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(c.code)
		for k, v := range c.headers {
			w.Header().Set(k, v)
		}
		if _, err := w.Write(c.responseBody); err != nil {
			logging.Error("Failed to send response: ", err)
		}
	})
}
