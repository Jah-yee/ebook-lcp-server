import React, { useMemo, useState } from "react";
import { Activity, BarChart3, CheckCircle2, KeyRound, Play, Shield } from "lucide-react";
import { createRoot } from "react-dom/client";
import "./styles.css";

type ProcessStatus = {
  id: string;
  status: string;
  publicationId?: string;
  error?: string;
  createdAt: string;
  updatedAt: string;
};

type StatusResponse = {
  status: string;
  uptimeSec: number;
  processes: ProcessStatus[];
};

type MetricsResponse = {
  uptimeSec: number;
  processes: number;
  metrics: {
    requestsTotal: number;
    processesOk: number;
    processesFail: number;
  };
};

const API_BASE = import.meta.env.VITE_API_BASE_URL || "";

function App() {
  const [token, setToken] = useState("");
  const [twoFactor, setTwoFactor] = useState("");
  const [title, setTitle] = useState("Example Publication");
  const [fileText, setFileText] = useState("LCP sample content");
  const [status, setStatus] = useState<StatusResponse | null>(null);
  const [metrics, setMetrics] = useState<MetricsResponse | null>(null);
  const [message, setMessage] = useState("");

  const authHeaders = useMemo(
    () => ({
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json"
    }),
    [token]
  );

  async function refreshStatus() {
    const response = await fetch(`${API_BASE}/api/v1/lcp/status`, { headers: authHeaders });
    const body = await response.json();
    if (!response.ok) throw new Error(body.error || "status request failed");
    setStatus(body);
  }

  async function refreshMetrics() {
    const response = await fetch(`${API_BASE}/api/v1/admin/metrics`, {
      headers: { ...authHeaders, "X-2FA-Code": twoFactor }
    });
    const body = await response.json();
    if (!response.ok) throw new Error(body.error || "metrics request failed");
    setMetrics(body);
  }

  async function processContent() {
    const response = await fetch(`${API_BASE}/api/v1/lcp/process`, {
      method: "POST",
      headers: authHeaders,
      body: JSON.stringify({ title, file: btoa(fileText) })
    });
    const body = await response.json();
    if (!response.ok) throw new Error(body.error || body.error || "process request failed");
    setMessage(`Process ${body.id} completed`);
    await refreshStatus();
  }

  async function run(action: () => Promise<void>) {
    setMessage("");
    try {
      await action();
    } catch (error) {
      setMessage(error instanceof Error ? error.message : "request failed");
    }
  }

  return (
    <main className="shell">
      <header className="topbar">
        <div>
          <h1>LCP Admin</h1>
          <p>Operations dashboard for publications, processing, and runtime health.</p>
        </div>
        <div className="status-pill">
          <Activity size={18} />
          {status?.status || "not loaded"}
        </div>
      </header>

      <section className="grid">
        <div className="panel auth-panel">
          <h2><Shield size={18} /> Access</h2>
          <label>
            JWT
            <textarea value={token} onChange={(event) => setToken(event.target.value)} />
          </label>
          <label>
            Admin 2FA
            <input value={twoFactor} onChange={(event) => setTwoFactor(event.target.value)} />
          </label>
        </div>

        <div className="panel">
          <h2><Play size={18} /> Process</h2>
          <label>
            Title
            <input value={title} onChange={(event) => setTitle(event.target.value)} />
          </label>
          <label>
            Content
            <textarea value={fileText} onChange={(event) => setFileText(event.target.value)} />
          </label>
          <button onClick={() => run(processContent)}>
            <CheckCircle2 size={18} />
            Submit
          </button>
        </div>

        <div className="panel">
          <h2><BarChart3 size={18} /> Metrics</h2>
          <button onClick={() => run(refreshMetrics)}>
            <KeyRound size={18} />
            Load Metrics
          </button>
          <dl className="metrics">
            <dt>Uptime</dt>
            <dd>{metrics?.uptimeSec ?? 0}s</dd>
            <dt>Requests</dt>
            <dd>{metrics?.metrics.requestsTotal ?? 0}</dd>
            <dt>OK / Failed</dt>
            <dd>{metrics ? `${metrics.metrics.processesOk} / ${metrics.metrics.processesFail}` : "0 / 0"}</dd>
          </dl>
        </div>
      </section>

      <section className="panel">
        <div className="section-head">
          <h2><Activity size={18} /> Process Status</h2>
          <button onClick={() => run(refreshStatus)}>Refresh</button>
        </div>
        {message && <p className="message">{message}</p>}
        <div className="table">
          <div className="row header">
            <span>ID</span>
            <span>Status</span>
            <span>Publication</span>
            <span>Updated</span>
          </div>
          {(status?.processes || []).map((item) => (
            <div className="row" key={item.id}>
              <span>{item.id}</span>
              <span>{item.status}</span>
              <span>{item.publicationId || "-"}</span>
              <span>{new Date(item.updatedAt).toLocaleString()}</span>
            </div>
          ))}
        </div>
      </section>
    </main>
  );
}

createRoot(document.getElementById("root")!).render(<App />);
