package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func (s *Server) longRunningProcess(w http.ResponseWriter, r *http.Request) {
	_, span := s.tracer.Start(r.Context(), "Long running task - handler")
	s.logger.Debug("Long running task", zap.Any("headers", r.Header))
	defer span.End()

	time.Sleep(time.Millisecond * 50)
	span.AddEvent("halfway done!")
	time.Sleep(time.Millisecond * 50)

	s.logger.Debug("Long running task DONE")
	w.Write([]byte("Done"))
}

func (s *Server) healthzHandler(w http.ResponseWriter, r *http.Request) {
	s.JSONResponse(w, r, map[string]string{"status": "OK"})
}

func (s *Server) readyzHandler(w http.ResponseWriter, r *http.Request) {
	s.JSONResponse(w, r, map[string]string{"status": "OK"})
}

func (s *Server) JSONResponse(w http.ResponseWriter, r *http.Request, result interface{}) {
	body, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error("JSON marshal failed", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	w.Write(prettyJSON(body))
}

func prettyJSON(b []byte) []byte {
	var out bytes.Buffer
	json.Indent(&out, b, "", "  ")
	return out.Bytes()
}
