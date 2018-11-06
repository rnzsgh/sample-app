package main

import (
	"compress/gzip"
	"encoding/json"
	"flag"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	log "github.com/golang/glog"
)

func serveStatic(file string, w http.ResponseWriter, r *http.Request) {

}

func main() {

	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")

	http.HandleFunc("/health", createGzipHandler(func(w http.ResponseWriter, r *http.Request) {
		log.Info("/health called")
		res := &response{Message: "ok"}
		out, _ := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, string(out))
	}))

	http.HandleFunc("/favicon.ico", createGzipHandler(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/vnd.microsoft.icon")
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Length", strconv.Itoa(len(faviconIco)))
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		_, _ = w.Write(faviconIco)
	}))

	http.HandleFunc("/", createGzipHandler(func(w http.ResponseWriter, r *http.Request) {

		f := fib()

		res := &response{Message: "Hello World"}

		for _, e := range os.Environ() {
			pair := strings.Split(e, "=")
			res.EnvVars = append(res.EnvVars, pair[0]+"="+pair[1])
		}

		for i := 1; i <= 90; i++ {
			res.Fib = append(res.Fib, f())
		}

		// Beautify the JSON output
		out, _ := json.MarshalIndent(res, "", "  ")

		// Normally this would be application/json, but we don't want to prompt downloads
		w.Header().Set("Content-Type", "text/plain")

		io.WriteString(w, string(out))

		log.Info("Hello world - the log message")
	}))
	http.ListenAndServe(":3000", nil)
	log.Flush()
}

func createGzipHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			f(w, r)
		} else {
			w.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(w)
			defer gz.Close()
			f(gzipResponseWriter{Writer: gz, ResponseWriter: w}, r)
		}
	}
}

type gzipResponseWriter struct {
	http.ResponseWriter
	io.Writer
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

type response struct {
	Message string   `json:"message"`
	EnvVars []string `json:"env"`
	Fib     []int    `json:"fib"`
}

var faviconIco = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xdc\xd2\x31\x4a\x03\x41\x14\x06\xe0\x3f\x98\xce\xc2\x05\xc1\xd6\x47\x46\x59\x4b\x6f\x90\x1c\xc1\x1b\x4c\xc4\xc6\x26\x18\x50\x11\x64\x14\xcb\x80\x4d\x3a\xb7\xcc\x51\x3c\x42\xf0\x04\xdb\x8e\xd5\xc8\x2c\xb2\xcd\x30\x32\xce\x06\x26\x21\x93\xdd\x69\xf3\xe0\xe7\xed\x3e\xf6\xfb\xab\x05\x7a\xe8\x21\xcb\xdc\x26\xdc\xf6\x81\x13\x00\x17\x00\x32\x00\x04\x7f\x77\x33\xea\x03\xc7\x87\x3e\xbb\x86\x09\x7d\xc5\x84\x5e\x0c\x1e\x54\x45\x5c\x2a\xe2\x72\x41\x5c\x8e\x5a\xcc\x25\x13\x7a\xce\x84\x2e\x99\xd0\xd6\x65\xf0\xa8\x2c\x71\x19\xa6\x24\x2e\xe7\xc4\x25\x05\x6e\xca\x84\x5e\xae\x4c\x98\x2d\x3e\xcc\x92\xb8\x9c\x6e\x73\x1d\xfd\x7f\xf6\xda\x3f\xfd\x58\xba\xf9\x4e\xf7\xaf\x95\x3d\x7f\xaf\x55\xfe\x61\x7e\x5d\xce\x5e\x2a\x15\xeb\x89\xba\xc2\xd8\xb5\x44\x7a\x5a\x5d\x4b\x4f\x67\x17\xe9\xc9\x0b\x33\xce\x0b\x53\x26\x7b\x6f\xc6\xab\xff\x38\xa1\x67\xcd\x6d\x4e\xd3\xa3\x52\xdd\x46\x47\x96\x17\xe6\xad\xe9\xd9\xe9\xea\xd3\xaf\xe7\xcf\x03\xf7\x34\x99\xf9\xcb\x5d\xb3\xaf\x23\xef\xc3\x7b\xf7\xbd\x3a\x9a\xcc\xea\x21\xf0\x17\x00\x00\xff\xff\x3d\x3b\x07\x48\x7e\x04\x00\x00")

func fib() func() int {
	a, b := 0, 1
	return func() int {
		a, b = b, a+b
		return a
	}

}
