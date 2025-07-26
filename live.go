package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"golang.org/x/net/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var reloadMu sync.Mutex
var reloadTimer *time.Timer

func live(path string, port int, openb bool) error {
	http.Handle("/ws", websocket.Handler(func(ws *websocket.Conn) {
		clients[ws] = true
		defer ws.Close()
		for {
			time.Sleep(1 * time.Hour)
		}
	}))

	dist := filepath.Join(path, "dist")
	fs := http.FileServer(http.Dir(dist))
	http.Handle("/", injectReloadScript(fs, dist))

	go watchFiles(path)

	portStr := strconv.Itoa(port)
	url := "http://localhost:" + portStr
	if openb {
		go openBrowser(url)
	}

	fmt.Println("Goblo is live", path, "on", url)
	if err := http.ListenAndServe(":"+portStr, nil); err != nil {
		return err
	}

	return nil
}

func injectReloadScript(next http.Handler, dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cleanPath := filepath.Clean(r.URL.Path)
		fullPath := filepath.Join(dir, cleanPath)

		info, err := os.Stat(fullPath)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		if info.IsDir() {
			indexPath := filepath.Join(fullPath, "index.html")
			if indexInfo, err := os.Stat(indexPath); err == nil && !indexInfo.IsDir() {
				data, err := os.ReadFile(indexPath)
				if err != nil {
					http.Error(w, "Failed to read index.html", 500)
					return
				}

				html := string(data)
				script := `
<script>
const ws = new WebSocket("ws://" + location.host + "/ws");
ws.onmessage = (msg) => { if (msg.data === "reload") location.reload(); };
</script>`

				modified := ""
				if idx := strings.Index(html, "</head>"); idx != -1 {
					modified = html[:idx] + script + html[idx:]
				} else {
					modified = html + script
				}

				w.Header().Set("Content-Type", "text/html")
				w.Write([]byte(modified))
				return
			} else {
				listDir(w, fullPath, cleanPath)
				return
			}
		}

		ext := strings.ToLower(filepath.Ext(fullPath))
		if ext == ".html" {
			data, err := os.ReadFile(fullPath)
			if err != nil {
				http.NotFound(w, r)
				return
			}

			html := string(data)
			script := `
<script>
const ws = new WebSocket("ws://" + location.host + "/ws");
ws.onmessage = (msg) => { if (msg.data === "reload") location.reload(); };
</script>`

			modified := ""
			if idx := strings.Index(html, "</head>"); idx != -1 {
				modified = html[:idx] + script + html[idx:]
			} else {
				modified = html + script
			}

			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(modified))
			return
		}

		next.ServeHTTP(w, r)
	}
}

func listDir(w http.ResponseWriter, dirPath, urlPath string) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		http.Error(w, "Failed to read directory", 500)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h1>Index of %s</h1><ul>", urlPath)
	for _, file := range files {
		name := file.Name()
		displayName := name
		if file.IsDir() {
			displayName += "/"
		}
		link := urlPath
		if !strings.HasSuffix(link, "/") {
			link += "/"
		}
		link += name
		fmt.Fprintf(w, `<li><a href="%s">%s</a></li>`, link, displayName)
	}
	fmt.Fprint(w, "</ul>")
}

func watchFiles(dir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if filepath.Base(path) == "dist" {
				return filepath.SkipDir
			}
			return watcher.Add(path)
		}
		return nil
	})

	var scheduleReload func()
	scheduleReload = func() {
		reloadMu.Lock()
		defer reloadMu.Unlock()
		if reloadTimer != nil {
			reloadTimer.Stop()
		}
		reloadTimer = time.AfterFunc(100*time.Millisecond, func() {
			if err := live_build(); err != nil {
				werror("Failed to build", err)
				return
			}
			fmt.Println("Reload at", time.Now().Format("15:04:05"))
			broadcastReload()
		})
	}

	for {
		select {
		case event := <-watcher.Events:
			if strings.Contains(event.Name, "/dist/") || strings.HasSuffix(event.Name, "/dist") {
				continue
			}
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove) != 0 {
				scheduleReload()
			}
		case err := <-watcher.Errors:
			log.Println("watch error:", err)
		}
	}

}

func broadcastReload() {
	for ws := range clients {
		_ = websocket.Message.Send(ws, "reload")
	}
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default: // linux
		cmd = exec.Command("xdg-open", url)
	}
	_ = cmd.Start()
}
