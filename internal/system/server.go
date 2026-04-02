package system

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// DashboardUpdate representa un hallazgo para la web.
type DashboardUpdate struct {
	Type      string      `json:"type"` // "finding", "log", "status"
	Source    string      `json:"source"`
	Content   interface{} `json:"content"`
	Timestamp string      `json:"timestamp"`
}

type WebDashboard struct {
	clients    map[chan string]bool
	mu         sync.Mutex
	updates    chan DashboardUpdate
	Port       int
}

func NewWebDashboard(port int) *WebDashboard {
	return &WebDashboard{
		clients: make(map[chan string]bool),
		updates: make(chan DashboardUpdate, 100),
		Port:    port,
	}
}

func (d *WebDashboard) Broadcast(update DashboardUpdate) {
	update.Timestamp = time.Now().Format("15:04:05")
	d.updates <- update
}

func (d *WebDashboard) Start() {
	http.HandleFunc("/events", d.sseHandler)
	http.HandleFunc("/", d.indexHandler)

	fmt.Printf("[System] Dashboard running at http://localhost:%d\n", d.Port)
	go http.ListenAndServe(fmt.Sprintf(":%d", d.Port), nil)
	
	// Distribuidor de updates
	go func() {
		for update := range d.updates {
			data, _ := json.Marshal(update)
			message := fmt.Sprintf("data: %s\n\n", string(data))
			
			d.mu.Lock()
			for client := range d.clients {
				client <- message
			}
			d.mu.Unlock()
		}
	}()
}

func (d *WebDashboard) sseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	clientChan := make(chan string)
	d.mu.Lock()
	d.clients[clientChan] = true
	d.mu.Unlock()

	defer func() {
		d.mu.Lock()
		delete(d.clients, clientChan)
		d.mu.Unlock()
		close(clientChan)
	}()

	for msg := range clientChan {
		fmt.Fprint(w, msg)
		w.(http.Flusher).Flush()
	}
}

func (d *WebDashboard) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `
<!DOCTYPE html>
<html>
<head>
    <title>Argus Dashboard</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background: #1a1a1a; color: #eee; margin: 0; padding: 20px; }
        .container { max-width: 1000px; margin: auto; }
        h1 { color: #00ff00; border-bottom: 2px solid #333; padding-bottom: 10px; }
        .card { background: #2a2a2a; border-radius: 8px; padding: 15px; margin-bottom: 15px; border-left: 5px solid #00ff00; animation: fadeIn 0.5s; }
        .timestamp { color: #888; font-size: 0.8em; float: right; }
        .source { font-weight: bold; color: #00d4ff; text-transform: uppercase; font-size: 0.9em; }
        pre { background: #000; padding: 10px; border-radius: 4px; overflow-x: auto; color: #adadad; }
        @keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
        .status-bar { position: fixed; bottom: 0; left: 0; right: 0; background: #333; padding: 5px 20px; font-size: 0.8em; color: #00ff00; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Argus :: Cognitive Dashboard</h1>
        <div id="findings"></div>
    </div>
    <div class="status-bar">Status: Connected to Argus Engine</div>
    <script>
        const evtSource = new EventSource("/events");
        const findingsDiv = document.getElementById('findings');
        
        evtSource.onmessage = function(event) {
            const data = JSON.parse(event.data);
            const card = document.createElement('div');
            card.className = 'card';
            
            let content = '';
            if (typeof data.content === 'string') {
                content = '<pre>' + data.content + '</pre>';
            } else {
                content = '<pre>' + JSON.stringify(data.content, null, 2) + '</pre>';
            }

            card.innerHTML = "<span class=\"timestamp\">" + data.timestamp + "</span>" +
                "<div class=\"source\">" + data.source + "</div>" +
                "<div class=\"content\">" + content + "</div>";
            findingsDiv.prepend(card);
        };
    </script>
</body>
</html>
	`)
}
