package routes

import (
    "fmt"
    "net/http"
)

// writeSSE writes a Server-Sent Event to the provided http.ResponseWriter.
// The event and data fields are escaped as per the SSE spec.
func writeSSE(w http.ResponseWriter, event, data string) {
    fmt.Fprintf(w, "event:%s\ndata:%s\n\n", event, data)
    if flusher, ok := w.(http.Flusher); ok {
        flusher.Flush()
    }
}
