package server

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/93percentai/sitemaker/internal/build"
	"github.com/93percentai/sitemaker/internal/config"

	"github.com/fsnotify/fsnotify"
)

type sseClients struct {
	mu      sync.Mutex
	clients map[chan struct{}]struct{}
}

func (s *sseClients) add(ch chan struct{}) {
	s.mu.Lock()
	s.clients[ch] = struct{}{}
	s.mu.Unlock()
}

func (s *sseClients) remove(ch chan struct{}) {
	s.mu.Lock()
	delete(s.clients, ch)
	s.mu.Unlock()
}

func (s *sseClients) notify() {
	s.mu.Lock()
	for ch := range s.clients {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
	s.mu.Unlock()
}

func Serve(cfg config.Config, defaults fs.FS, port int) error {
	log.Printf("Building site...")
	if err := build.Build(cfg, defaults); err != nil {
		return fmt.Errorf("initial build: %w", err)
	}

	clients := &sseClients{clients: make(map[chan struct{}]struct{})}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("creating watcher: %w", err)
	}
	defer watcher.Close()

	watchDirs := []string{cfg.Build.ContentDir, cfg.Build.TemplateDir, cfg.Build.StaticDir}
	for _, dir := range watchDirs {
		if _, err := os.Stat(dir); err == nil {
			addRecursive(watcher, dir)
		}
	}

	go func() {
		var debounce *time.Timer
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove) == 0 {
					continue
				}
				if debounce != nil {
					debounce.Stop()
				}
				debounce = time.AfterFunc(200*time.Millisecond, func() {
					log.Printf("Rebuilding...")
					if err := build.Build(cfg, defaults); err != nil {
						log.Printf("Build error: %v", err)
					} else {
						log.Printf("Build complete.")
						clients.notify()
					}
				})
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Watcher error: %v", err)
			}
		}
	}()

	mux := http.NewServeMux()

	mux.HandleFunc("/__reload", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming not supported", http.StatusInternalServerError)
			return
		}

		ch := make(chan struct{}, 1)
		clients.add(ch)
		defer clients.remove(ch)

		notify := r.Context().Done()
		for {
			select {
			case <-ch:
				fmt.Fprintf(w, "data: reload\n\n")
				flusher.Flush()
			case <-notify:
				return
			}
		}
	})

	fileServer := http.FileServer(http.Dir(cfg.Build.OutputDir))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(cfg.Build.OutputDir, r.URL.Path)

		if info, err := os.Stat(path); err == nil && info.IsDir() {
			index := filepath.Join(path, "index.html")
			if _, err := os.Stat(index); err == nil {
				path = index
			}
		}

		if strings.HasSuffix(path, ".html") || strings.HasSuffix(path, "/index.html") {
			data, err := os.ReadFile(path)
			if err != nil {
				fileServer.ServeHTTP(w, r)
				return
			}
			injected := injectLiveReload(data)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(injected)
			return
		}

		fileServer.ServeHTTP(w, r)
	})

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Serving at http://localhost:%d", port)
	log.Printf("Live reload enabled. Watching for changes...")
	return http.ListenAndServe(addr, mux)
}

func injectLiveReload(html []byte) []byte {
	script := []byte(`<script>
(function() {
  var es = new EventSource('/__reload');
  es.onmessage = function() { location.reload(); };
  es.onerror = function() { setTimeout(function() { location.reload(); }, 1000); };
})();
</script>
</body>`)

	return []byte(strings.Replace(string(html), "</body>", string(script), 1))
}

func addRecursive(watcher *fsnotify.Watcher, root string) {
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			watcher.Add(path)
		}
		return nil
	})
}
