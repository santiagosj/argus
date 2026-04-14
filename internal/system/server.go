package system

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

// DashboardUpdate representa un hallazgo para la web.
type DashboardUpdate struct {
	Type      string      `json:"type"` // "finding", "log", "status", "sdd_progress"
	Source    string      `json:"source"`
	Phase     string      `json:"phase,omitempty"`
	Content   interface{} `json:"content"`
	Timestamp string      `json:"timestamp"`
}

type WebDashboard struct {
	clients map[chan string]bool
	mu      sync.Mutex
	updates chan DashboardUpdate
	Port    int
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
	mux := http.NewServeMux()
	mux.HandleFunc("/events", d.sseHandler)
	mux.HandleFunc("/", d.indexHandler)

	fmt.Printf("[System] Dashboard running at http://localhost:%d\n", d.Port)

	// Detect WSL and show external IP
	if ip := getLocalIP(); ip != "" {
		fmt.Printf("[System] Also accessible from Windows at: http://%s:%d\n", ip, d.Port)
	}

	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", d.Port), mux)
		if err != nil {
			fmt.Printf("[System] Dashboard failed: %v\n", err)
		}
	}()

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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `
<!DOCTYPE html>
<html>
<head>
    <title>Argus Dashboard</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background: #121212; color: #f0f0f0; margin: 0; padding: 0; }
        .dashboard { max-width: 1300px; margin: auto; padding: 20px 24px 36px; }
        .hero { display: flex; flex-wrap: wrap; justify-content: space-between; align-items: center; gap: 16px; margin-bottom: 24px; }
        .hero h1 { margin: 0; font-size: 2.4rem; letter-spacing: -0.03em; color: #7df9ff; }
        .hero p { margin: 6px 0 0; color: #cfd8dc; }
        .status-card { background: #1f1f25; border: 1px solid #2c2c33; border-radius: 16px; padding: 18px 22px; min-width: 260px; box-shadow: 0 16px 40px rgba(0,0,0,0.18); }
        .status-card strong { display: block; margin-bottom: 10px; color: #a5ffb2; }
        .status-pill { display: inline-flex; align-items: center; gap: 8px; background: #111; border: 1px solid #333; border-radius: 999px; padding: 8px 12px; font-size: 0.9rem; }
        .status-pill.online { border-color: #39d98a; color: #39d98a; }
        .status-pill.offline { border-color: #ff5f5f; color: #ff5f5f; }
        .status-pill.connecting { border-color: #ffd86f; color: #ffd86f; }
        .summary-row { display: grid; grid-template-columns: repeat(3, minmax(180px, 1fr)); gap: 14px; margin-top: 12px; }
        .summary-box { background: #19191f; border-radius: 12px; padding: 14px 16px; border: 1px solid #2d2d34; }
        .summary-box span { display: block; font-size: 0.75rem; text-transform: uppercase; color: #7a7f8a; margin-bottom: 6px; }
        .summary-box strong { display: block; font-size: 1.4rem; color: #ffffff; }
        .layout { display: grid; grid-template-columns: 320px 1fr; gap: 20px; }
        .panel { background: #18181f; border: 1px solid #2c2c33; border-radius: 18px; padding: 20px; box-shadow: 0 16px 40px rgba(0,0,0,0.18); }
        .panel h2 { margin: 0 0 18px; font-size: 1.1rem; color: #d7eaff; }
        .phase-list { display: grid; gap: 12px; }
        .phase-chip { display: flex; justify-content: space-between; align-items: center; gap: 10px; padding: 12px 14px; border-radius: 14px; border: 1px solid #2c2c33; cursor: pointer; transition: transform 0.15s ease, border-color 0.15s ease; }
        .phase-chip:hover { transform: translateY(-1px); border-color: #5b8cff; }
        .phase-chip.active { background: #0d2944; border-color: #5b8cff; }
        .phase-chip span:first-child { text-transform: uppercase; font-size: 0.85rem; letter-spacing: 0.05em; color: #8bc2ff; }
        .phase-chip strong { font-size: 1rem; color: #ffffff; }
        .phase-chip .badge { background: #252933; border-radius: 999px; padding: 6px 10px; font-size: 0.8rem; color: #c0c5d8; }
        .filter-group { display: grid; gap: 10px; margin-top: 18px; }
        .filter-group input { width: 100%; box-sizing: border-box; border-radius: 14px; border: 1px solid #2c2c33; background: #11131a; color: #f2f2f2; padding: 12px 14px; font-family: inherit; }
        .filter-group button { width: 100%; box-sizing: border-box; border: none; border-radius: 14px; background: #3a6eff; color: white; padding: 12px; cursor: pointer; font-weight: 600; font-family: inherit; }
        .filter-group button:hover { background: #4b7bff; }
        .phase-tabs { display: flex; flex-wrap: wrap; gap: 10px; margin-bottom: 18px; }
        .tab { padding: 10px 16px; border-radius: 999px; border: 1px solid #2c2c33; background: #11131a; color: #cbd4e6; cursor: pointer; transition: background 0.15s ease, border-color 0.15s ease; }
        .tab.active { background: #0f2b5b; border-color: #3c6eff; color: #ffffff; }
        .stage-section { display: none; gap: 16px; }
        .stage-section.active { display: grid; }
        .finding-card { background: #16161e; border: 1px solid #2b2b33; border-radius: 16px; padding: 18px; transition: box-shadow 0.2s ease; }
        .finding-card:hover { box-shadow: 0 16px 30px rgba(0,0,0,0.25); }
        .finding-title { display: flex; justify-content: space-between; align-items: center; gap: 12px; margin-bottom: 12px; }
        .finding-title h3 { margin: 0; font-size: 1rem; color: #e6f7ff; }
        .finding-title .tag { text-transform: uppercase; font-size: 0.75rem; color: #9cd6ff; letter-spacing: 0.06em; }
        .finding-detail { display: grid; gap: 8px; font-size: 0.95rem; line-height: 1.5; color: #d3d7dd; }
        .finding-detail pre { margin: 0; background: #10101a; padding: 12px; border-radius: 12px; overflow-x: auto; color: #d8d8d8; }
        .footer { margin-top: 28px; color: #7a7f8a; font-size: 0.9rem; text-align: center; }
    </style>
</head>
<body>
    <div class="dashboard">
        <header class="hero">
            <div>
                <h1>Argus Cognitive Dashboard</h1>
                <p>Demostración guiada por fases con hallazgos organizados y priorización visual.</p>
            </div>
            <aside class="status-card">
                <strong>Estado de la conexión</strong>
                <div id="connection-status" class="status-pill connecting">Conectando...</div>
                <div class="summary-row">
                    <div class="summary-box"><span>Total de eventos</span><strong id="total-events">0</strong></div>
                    <div class="summary-box"><span>Hallazgos</span><strong id="total-findings">0</strong></div>
                    <div class="summary-box"><span>Último evento</span><strong id="last-event">-</strong></div>
                    <div class="summary-box"><span>Fase SDD actual</span><strong id="current-phase">-</strong></div>
                    <div class="summary-box"><span>Próximo paso</span><strong id="next-phase">-</strong></div>
                    <div class="summary-box"><span>Artefactos SDD</span><strong id="artifact-count">0</strong></div>
                </div>
            </aside>
        </header>
        <div class="layout">
            <aside class="panel">
                <h2>Resumen de fases</h2>
                <div id="phase-summary" class="phase-list"></div>
                <div class="filter-group">
                    <input id="filter-input" placeholder="Buscar hallazgos, herramientas, frases..." />
                    <button id="clear-filter">Limpiar filtro</button>
                </div>
            </aside>
            <main class="panel">
                <div class="phase-tabs" id="phase-tabs"></div>
                <section id="phase-all" class="stage-section active"></section>
                <section id="phase-identify" class="stage-section"></section>
                <section id="phase-protect" class="stage-section"></section>
                <section id="phase-detect" class="stage-section"></section>
                <section id="phase-respond" class="stage-section"></section>
                <section id="phase-recover" class="stage-section"></section>
                <section id="phase-evolve" class="stage-section"></section>
            </main>
        </div>
        <footer class="footer">Navegue entre fases, filtre hallazgos y prepare un demo claro para OWASP Juice Shop u objetivos similares.</footer>
    </div>
    <script>
        const phases = ['all', 'identify', 'protect', 'detect', 'respond', 'recover', 'evolve'];
        const phaseLabels = {
            all: 'Todos',
            identify: 'Identify',
            protect: 'Protect',
            detect: 'Detect',
            respond: 'Respond',
            recover: 'Recover',
            evolve: 'Evolve'
        };
        const counts = { all: 0, identify: 0, protect: 0, detect: 0, respond: 0, recover: 0, evolve: 0 };
        const phaseColor = { all: '#7df9ff', identify: '#8bc2ff', protect: '#73ff94', detect: '#ffd35b', respond: '#ff9376', recover: '#a37cff', evolve: '#ffb3d9' };

        const phaseSummary = document.getElementById('phase-summary');
        const phaseTabs = document.getElementById('phase-tabs');
        const connectionStatus = document.getElementById('connection-status');
        const totalEvents = document.getElementById('total-events');
        const totalFindings = document.getElementById('total-findings');
        const lastEvent = document.getElementById('last-event');
        const currentPhaseDisplay = document.getElementById('current-phase');
        const nextPhaseDisplay = document.getElementById('next-phase');
        const artifactCountDisplay = document.getElementById('artifact-count');
        const filterInput = document.getElementById('filter-input');
        const clearFilter = document.getElementById('clear-filter');

        const sections = {};
        phases.forEach(phase => { sections[phase] = document.getElementById('phase-' + phase); });

        function createPhaseSummary() {
            phaseSummary.innerHTML = '';
            phases.forEach(phase => {
                const card = document.createElement('div');
                card.className = 'phase-chip';
                card.id = 'summary-' + phase;
                card.onclick = () => showTab(phase);
                if (phase === 'all') card.classList.add('active');
                card.innerHTML = '<span>' + phaseLabels[phase] + '</span><strong>' + counts[phase] + '</strong>';
                card.style.borderColor = phaseColor[phase];
                phaseSummary.appendChild(card);
            });
        }

        function createPhaseTabs() {
            phaseTabs.innerHTML = '';
            phases.forEach(phase => {
                const tab = document.createElement('div');
                tab.className = 'tab' + (phase === 'all' ? ' active' : '');
                tab.id = 'tab-' + phase;
                tab.textContent = phaseLabels[phase];
                tab.onclick = () => showTab(phase);
                phaseTabs.appendChild(tab);
            });
        }

        function showTab(phase) {
            phases.forEach(p => {
                document.getElementById('tab-' + p).classList.toggle('active', p === phase);
                document.getElementById('summary-' + p).classList.toggle('active', p === phase);
                sections[p].classList.toggle('active', p === phase);
            });
        }

        function getPhase(source) {
            const lower = source.toLowerCase();
            if (lower.includes('identify') || lower.includes('recon') || lower.includes('nmap') || lower.includes('nuclei')) return 'identify';
            if (lower.includes('protect')) return 'protect';
            if (lower.includes('detect')) return 'detect';
            if (lower.includes('respond')) return 'respond';
            if (lower.includes('recover')) return 'recover';
            if (lower.includes('evolve') || lower.includes('sdd')) return 'evolve';
            return 'all';
        }

        function createCard(data, phase) {
            const card = document.createElement('div');
            card.className = 'finding-card';
            const source = data.source || 'Unknown';
            const type = data.type || 'log';
            const content = typeof data.content === 'string' ? data.content : JSON.stringify(data.content, null, 2);
            const time = data.timestamp || '-';
            const displayPhase = phaseLabels[phase] || data.phase || phase;
            const displayColor = phaseColor[phase] || '#ffb3d9';
            card.innerHTML =
                '<div class="finding-title">' +
                '<h3>' + source + '</h3>' +
                '<div class="tag" style="color:' + displayColor + '">' + displayPhase + '</div>' +
                '</div>' +
                '<div class="finding-detail">' +
                '<div><strong>Tipo:</strong> ' + type + '</div>' +
                '<div><strong>Hora:</strong> ' + time + '</div>' +
                '<pre>' + content + '</pre>' +
                '</div>';
            return card;
        }

        function addFinding(data) {
            const phase = data.phase ? 'evolve' : getPhase(data.source);
            const content = typeof data.content === 'string' ? data.content : JSON.stringify(data.content, null, 2);
            if (filterInput.value.trim()) {
                const filter = filterInput.value.toLowerCase();
                if (!data.source.toLowerCase().includes(filter) && !content.toLowerCase().includes(filter)) {
                    return;
                }
            }
            const item = createCard(data, phase);
            item.setAttribute('data-ts', data.timestamp || '');
            const duplicate = sections[phase].querySelector('.finding-card[data-ts="' + data.timestamp + '"]');
            if (duplicate) return;
            counts.all++;
            counts[phase]++;
            totalEvents.textContent = counts.all;
            totalFindings.textContent = counts.all;
            lastEvent.textContent = data.timestamp;
            phaseSummary.querySelectorAll('.phase-chip').forEach(card => {
                const key = card.id.replace('summary-', '');
                card.querySelector('strong').textContent = counts[key];
            });
            sections[phase].prepend(item.cloneNode(true));
            sections.all.prepend(item);
        }

        function updateSDDStatus(data) {
            if (!data.content || typeof data.content !== 'object') {
                return;
            }
            const phase = data.phase || data.source || '-';
            currentPhaseDisplay.textContent = phase;
            nextPhaseDisplay.textContent = data.content.next || nextPhaseDisplay.textContent;
            if (Array.isArray(data.content.artifacts)) {
                artifactCountDisplay.textContent = data.content.artifacts.length;
            }
        }

        createPhaseSummary();
        createPhaseTabs();

        const evtSource = new EventSource('/events');
        evtSource.onopen = function() {
            connectionStatus.textContent = 'Conectado';
            connectionStatus.className = 'status-pill online';
        };

        evtSource.onerror = function() {
            connectionStatus.textContent = 'Conexión perdida';
            connectionStatus.className = 'status-pill offline';
        };

        evtSource.onmessage = function(event) {
            const data = JSON.parse(event.data);
            if (data.type === 'sdd_progress') {
                updateSDDStatus(data);
            }
            addFinding(data);
        };

        clearFilter.onclick = function() {
            filterInput.value = '';
            location.reload();
        };
    </script>
</body>
</html>
	`)
}

func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
