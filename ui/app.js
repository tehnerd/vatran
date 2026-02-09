(() => {
  const { useEffect, useMemo, useRef, useState, useContext } = React;
  const {
    BrowserRouter,
    Routes,
    Route,
    NavLink,
    Link,
    useParams,
    useNavigate,
  } = ReactRouterDOM;

  const html = htm.bind(React.createElement);
  const ToastContext = React.createContext({ addToast: () => {} });

  const VIP_FLAG_OPTIONS = [
    { label: "NO_SPORT", value: 1 },
    { label: "NO_LRU", value: 2 },
    { label: "QUIC_VIP", value: 4 },
    { label: "DPORT_HASH", value: 8 },
    { label: "SRC_ROUTING", value: 16 },
    { label: "LOCAL_VIP", value: 32 },
    { label: "GLOBAL_LRU", value: 64 },
    { label: "HASH_SRC_DST_PORT", value: 128 },
    { label: "UDP_STABLE_ROUTING_VIP", value: 256 },
    { label: "UDP_FLOW_MIGRATION", value: 512 },
  ];

  function useToast() {
    return useContext(ToastContext);
  }

  function sanitizeFlagId(value) {
    return String(value).replace(/[^a-z0-9_-]/gi, "_");
  }

  function toggleFlag(mask, flagValue, enabled) {
    const current = Number(mask) || 0;
    const flag = Number(flagValue) || 0;
    if (enabled) {
      return current | flag;
    }
    return current & ~flag;
  }

  function getEnabledFlags(mask, options) {
    const value = Number(mask) || 0;
    return options.filter((opt) => (value & opt.value) !== 0);
  }

  function parseNumberList(raw, label) {
    const text = String(raw ?? "").trim();
    if (!text) return undefined;
    const parts = text.split(/[\s,]+/).filter(Boolean);
    const numbers = parts.map((item) => Number(item));
    const invalid = numbers.findIndex((value) => !Number.isFinite(value) || !Number.isInteger(value));
    if (invalid !== -1) {
      throw new Error(`${label} must be a comma- or space-separated list of integers.`);
    }
    return numbers;
  }

  function FlagTable({ mask, options, showStatus = false, emptyLabel = "None" }) {
    const value = Number(mask) || 0;
    const rows = showStatus ? options : getEnabledFlags(value, options);
    const colSpan = showStatus ? 2 : 1;
    return html`
      <table className="table flag-table">
        <thead>
          <tr>
            <th>Flag</th>
            ${showStatus ? html`<th>Enabled</th>` : null}
          </tr>
        </thead>
        <tbody>
          ${rows.length === 0
            ? html`<tr><td colspan=${colSpan} className="muted">${emptyLabel}</td></tr>`
            : rows.map((opt) => {
                const enabled = (value & opt.value) !== 0;
                return html`
                  <tr>
                    <td>${opt.label}</td>
                    ${showStatus ? html`<td>${enabled ? "Yes" : "No"}</td>` : null}
                  </tr>
                `;
              })}
        </tbody>
      </table>
    `;
  }

  function FlagSelector({ options, value, onChange, name }) {
    const current = Number(value) || 0;
    const prefix = sanitizeFlagId(name || "flags");
    return html`
      <div className="flag-selector">
        <table className="table flag-table">
          <thead>
            <tr>
              <th>Enabled</th>
              <th>Flag</th>
            </tr>
          </thead>
          <tbody>
            ${options.map((opt) => {
              const id = `${prefix}-${opt.value}`;
              const checked = (current & opt.value) === opt.value;
              return html`
                <tr>
                  <td>
                    <input
                      id=${id}
                      type="checkbox"
                      checked=${checked}
                      onChange=${(e) => onChange(toggleFlag(current, opt.value, e.target.checked))}
                    />
                  </td>
                  <td>
                    <label className="flag-option" htmlFor=${id}>${opt.label}</label>
                  </td>
                </tr>
              `;
            })}
          </tbody>
        </table>
      </div>
    `;
  }

  const api = {
    base: "/api/v1",
    async request(path, options = {}) {
      const init = {
        method: options.method || "GET",
        headers: { "Content-Type": "application/json" },
      };
      let url = `${api.base}${path}`;
      if (options.body !== undefined && options.body !== null) {
        if (init.method === "GET") {
          const params = new URLSearchParams();
          Object.entries(options.body).forEach(([key, value]) => {
            if (value === undefined || value === null) return;
            if (Array.isArray(value)) {
              value.forEach((item) => params.append(key, String(item)));
              return;
            }
            if (typeof value === "object") {
              params.set(key, JSON.stringify(value));
              return;
            }
            params.set(key, String(value));
          });
          const query = params.toString();
          if (query) {
            url += `${url.includes("?") ? "&" : "?"}${query}`;
          }
        } else {
          init.body = JSON.stringify(options.body);
        }
      }
      const res = await fetch(url, init);
      let payload;
      try {
        payload = await res.json();
      } catch (err) {
        throw new Error("invalid JSON response");
      }
      if (!res.ok) {
        throw new Error(payload?.error?.message || `HTTP ${res.status}`);
      }
      if (!payload.success) {
        const msg = payload.error?.message || "request failed";
        throw new Error(msg);
      }
      return payload.data;
    },
    get(path, body) {
      return api.request(path, { method: "GET", body });
    },
    post(path, body) {
      return api.request(path, { method: "POST", body });
    },
    put(path, body) {
      return api.request(path, { method: "PUT", body });
    },
    del(path, body) {
      return api.request(path, { method: "DELETE", body });
    },
  };

  const TARGET_GROUPS_STORAGE_KEY = "vatran_target_groups";

  function normalizeReal(raw) {
    if (!raw || !raw.address) return null;
    const address = String(raw.address).trim();
    if (!address) return null;
    const weight = Number(raw.weight);
    const flags = Number(raw.flags ?? 0);
    return {
      address,
      weight: Number.isFinite(weight) ? weight : 0,
      flags: Number.isFinite(flags) ? flags : 0,
    };
  }

  function normalizeTargetGroups(groups) {
    if (!groups || typeof groups !== "object") return {};
    const result = {};
    Object.entries(groups).forEach(([name, list]) => {
      const groupName = String(name).trim();
      if (!groupName) return;
      const reals = Array.isArray(list) ? list.map(normalizeReal).filter(Boolean) : [];
      const unique = [];
      const seen = new Set();
      reals.forEach((real) => {
        if (seen.has(real.address)) return;
        seen.add(real.address);
        unique.push(real);
      });
      result[groupName] = unique;
    });
    return result;
  }

  function loadStoredTargetGroups() {
    if (typeof localStorage === "undefined") return {};
    try {
      const raw = localStorage.getItem(TARGET_GROUPS_STORAGE_KEY);
      if (!raw) return {};
      return normalizeTargetGroups(JSON.parse(raw));
    } catch (err) {
      return {};
    }
  }

  function saveStoredTargetGroups(groups) {
    if (typeof localStorage === "undefined") return;
    try {
      localStorage.setItem(TARGET_GROUPS_STORAGE_KEY, JSON.stringify(groups));
    } catch (err) {
      // ignore storage errors
    }
  }

  function mergeTargetGroups(base, incoming) {
    const next = { ...base };
    Object.entries(incoming || {}).forEach(([name, list]) => {
      if (!next[name]) {
        next[name] = list;
      }
    });
    return next;
  }

  function useTargetGroups() {
    const [groups, setGroups] = useState(() => loadStoredTargetGroups());

    useEffect(() => {
      saveStoredTargetGroups(groups);
    }, [groups]);

    const refreshFromStorage = () => {
      setGroups(loadStoredTargetGroups());
    };

    const importFromRunningConfig = async () => {
      const data = await api.get("/config/export/json");
      const imported = normalizeTargetGroups(data?.target_groups || {});
      const next = mergeTargetGroups(loadStoredTargetGroups(), imported);
      setGroups(next);
      saveStoredTargetGroups(next);
      return next;
    };

    return { groups, setGroups, refreshFromStorage, importFromRunningConfig };
  }

  function vipIdFromVip(vip) {
    return `${encodeURIComponent(vip.address)}:${vip.port}:${vip.proto}`;
  }

  function parseVipId(vipId) {
    const parts = vipId.split(":");
    const proto = Number(parts.pop() || 0);
    const port = Number(parts.pop() || 0);
    const address = decodeURIComponent(parts.join(":"));
    return { address, port, proto };
  }

  function usePolling(fetcher, intervalMs, deps = []) {
    const [data, setData] = useState(null);
    const [error, setError] = useState("");
    const [loading, setLoading] = useState(true);

    useEffect(() => {
      let mounted = true;
      const load = async () => {
        try {
          const result = await fetcher();
          if (mounted) {
            setData(result);
            setError("");
            setLoading(false);
          }
        } catch (err) {
          if (mounted) {
            setError(err.message || "request failed");
            setLoading(false);
          }
        }
      };
      load();
      const interval = setInterval(load, intervalMs);
      return () => {
        mounted = false;
        clearInterval(interval);
      };
    }, deps);

    return { data, error, loading };
  }

  function useStatSeries({ path, body, intervalMs = 1000, limit = 60 }) {
    const [points, setPoints] = useState([]);
    const [error, setError] = useState("");
    const bodyKey = useMemo(() => JSON.stringify(body || {}), [body]);

    useEffect(() => {
      if (body === null) {
        setPoints([]);
        setError("");
        return () => {};
      }
      let mounted = true;
      const load = async () => {
        try {
          const result = await api.get(path, body);
          if (!mounted) return;
          const label = new Date().toLocaleTimeString();
          setPoints((prev) => {
            const next = prev.concat({ label, ...result });
            return next.slice(-limit);
          });
          setError("");
        } catch (err) {
          if (mounted) {
            setError(err.message || "request failed");
          }
        }
      };
      load();
      const interval = setInterval(load, intervalMs);
      return () => {
        mounted = false;
        clearInterval(interval);
      };
    }, [path, bodyKey, intervalMs, limit]);

    return { points, error };
  }

  function ChartCanvas({
    title,
    points,
    keys,
    diff = false,
    height = 120,
    showTitle = false,
    selectedLabel = null,
    onPointSelect = null,
    onLegendSelect = null,
  }) {
    const canvasRef = useRef(null);
    const chartRef = useRef(null);

    useEffect(() => {
      if (!canvasRef.current) return;
      if (!chartRef.current) {
        chartRef.current = new Chart(canvasRef.current, {
          type: "line",
          data: {
            labels: [],
            datasets: [],
          },
          options: {
            responsive: true,
            maintainAspectRatio: false,
            animation: false,
            scales: {
              x: { grid: { display: false } },
              y: { beginAtZero: !diff },
            },
            plugins: {
              legend: { display: true, position: "bottom" },
              title: { display: showTitle && Boolean(title), text: title },
            },
          },
        });
      }
      const chart = chartRef.current;
      const hiddenByLabel = new Map(
        (chart.data.datasets || [])
          .filter((dataset) => typeof dataset.hidden !== "undefined")
          .map((dataset) => [dataset.label, dataset.hidden])
      );
      const labels = points.map((p) => p.label);
      const selectedKeys = selectedLabel
        ? keys.filter((key) => key.label === selectedLabel)
        : keys;
      const visibleKeys = selectedLabel && selectedKeys.length === 0 ? keys : selectedKeys;
      chart.data.labels = labels;
      chart.data.datasets = visibleKeys.map((key) => {
        const values = points.map((p) => p[key.field] || 0);
        const data = diff
          ? values.map((value, index) => (index === 0 ? 0 : value - values[index - 1]))
          : values;
        return {
          label: key.label,
          data,
          borderColor: key.color,
          backgroundColor: key.fill,
          borderWidth: 2,
          tension: 0.3,
          hidden: hiddenByLabel.get(key.label),
        };
      });
      chart.options.onClick = (event, elements) => {
        if (!onPointSelect || !elements || elements.length === 0) return;
        const index = elements[0].datasetIndex;
        const label = chart.data.datasets?.[index]?.label;
        if (label) onPointSelect(label);
      };
      if (chart.options.plugins && chart.options.plugins.legend) {
        chart.options.plugins.legend.onClick = (event, legendItem) => {
          if (!onLegendSelect) return;
          const label = legendItem?.text;
          if (label) onLegendSelect(label);
        };
      }
      chart.options.scales.y.beginAtZero = !diff;
      chart.options.plugins.title.display = showTitle && Boolean(title);
      chart.options.plugins.title.text = title || "";
      chart.update();
      return () => {};
    }, [points, keys, title, diff, showTitle, selectedLabel, onPointSelect, onLegendSelect]);

    useEffect(() => {
      return () => {
        if (chartRef.current) {
          chartRef.current.destroy();
          chartRef.current = null;
        }
      };
    }, []);

    return html`<canvas ref=${canvasRef} height=${height}></canvas>`;
  }

  function StatChart({ title, points, keys, diff = false, inlineTitle = true }) {
    const [zoomed, setZoomed] = useState(false);
    const [selectedLabel, setSelectedLabel] = useState(null);
    const suppressZoomRef = useRef(false);

    return html`
      <div className="chart-wrap">
        <button className="chart-zoom-button" type="button" onClick=${() => setZoomed(true)}>
          <span className="zoom-icon" aria-hidden="true">+</span>
          Zoom
        </button>
        <div
          className="chart-click"
          onClick=${() => {
            if (suppressZoomRef.current) {
              suppressZoomRef.current = false;
              return;
            }
            setZoomed(true);
          }}
        >
          <${ChartCanvas}
            title=${title}
            points=${points}
            keys=${keys}
            diff=${diff}
            height=${120}
            showTitle=${inlineTitle && Boolean(title)}
            selectedLabel=${selectedLabel}
            onPointSelect=${(label) => {
              setSelectedLabel((prev) => (prev === label ? null : label));
              suppressZoomRef.current = true;
              setTimeout(() => {
                suppressZoomRef.current = false;
              }, 0);
            }}
            onLegendSelect=${(label) => {
              setSelectedLabel((prev) => (prev === label ? null : label));
              suppressZoomRef.current = true;
              setTimeout(() => {
                suppressZoomRef.current = false;
              }, 0);
            }}
          />
        </div>
        ${zoomed &&
        html`
          <div className="chart-overlay" onClick=${() => setZoomed(false)}>
            <div className="chart-modal" onClick=${(e) => e.stopPropagation()}>
              <div className="row chart-modal-header">
                <div>
                  <h3>${title || "Chart"}</h3>
                  ${diff ? html`<p className="muted">Per-second delta.</p>` : ""}
                </div>
                <button className="btn ghost" onClick=${() => setZoomed(false)}>
                  Close
                </button>
              </div>
              <div className="chart-zoom">
                <${ChartCanvas}
                  title=${title}
                  points=${points}
                  keys=${keys}
                  diff=${diff}
                  height=${360}
                  showTitle=${false}
                  selectedLabel=${selectedLabel}
                  onPointSelect=${(label) =>
                    setSelectedLabel((prev) => (prev === label ? null : label))
                  }
                  onLegendSelect=${(label) =>
                    setSelectedLabel((prev) => (prev === label ? null : label))
                  }
                />
              </div>
            </div>
          </div>
        `}
      </div>
    `;
  }

  function AuthGate({ children }) {
    return children;
  }

  function Toasts({ toasts, onDismiss }) {
    return html`
      <div className="toast-stack">
        ${toasts.map(
          (toast) => html`
            <div className=${`toast ${toast.kind}`}>
              <span>${toast.message}</span>
              <button className="toast-close" onClick=${() => onDismiss(toast.id)}>
                ×
              </button>
            </div>
          `
        )}
      </div>
    `;
  }

  function Header({ status }) {
    return html`
      <header>
        <div>
          <div style=${{ fontSize: 20, fontWeight: 700 }}>Vatran</div>
          <div className="status-pill">
            <span className=${`dot ${status.ready ? "ok" : ""}`}></span>
            ${status.ready ? "Ready" : "Not ready"}
          </div>
        </div>
        <nav>
          <${NavLink} to="/" end className=${({ isActive }) => (isActive ? "active" : "")}>
            Dashboard
          </${NavLink}>
          <${NavLink} to="/stats/global" className=${({ isActive }) => (isActive ? "active" : "")}>
            Global stats
          </${NavLink}>
          <${NavLink} to="/stats/real" className=${({ isActive }) => (isActive ? "active" : "")}>
            Per-real stats
          </${NavLink}>
          <${NavLink}
            to="/target-groups"
            className=${({ isActive }) => (isActive ? "active" : "")}
          >
            Target groups
          </${NavLink}>
          <${NavLink} to="/config" className=${({ isActive }) => (isActive ? "active" : "")}>
            Config export
          </${NavLink}>
        </nav>
      </header>
    `;
  }

  function Dashboard() {
    const { addToast } = useToast();
    const [status, setStatus] = useState({ initialized: false, ready: false });
    const [vips, setVips] = useState([]);
    const [error, setError] = useState("");
    const [showInit, setShowInit] = useState(false);
    const [showVipForm, setShowVipForm] = useState(false);
    const [initForm, setInitForm] = useState({
      main_interface: "",
      balancer_prog_path: "",
      healthchecking_prog_path: "",
      default_mac: "",
      local_mac: "",
      root_map_path: "",
      root_map_pos: 2,
      katran_src_v4: "",
      katran_src_v6: "",
      use_root_map: false,
      max_vips: 512,
      max_reals: 4096,
      hash_function: "maglev_v2",
      forwarding_cores: "",
      numa_nodes: "",
    });
    const [vipForm, setVipForm] = useState({
      address: "",
      port: 80,
      proto: 6,
      flags: 0,
    });

    const load = async () => {
      try {
        const lbStatus = await api.get("/lb/status");
        const list = await api.get("/vips");
        const vipsWithFlags = await Promise.all(
          (list || []).map(async (vip) => {
            let flags = null;
            let healthy = false;
            try {
              const flagData = await api.get("/vips/flags", {
                address: vip.address,
                port: vip.port,
                proto: vip.proto,
              });
              flags = flagData?.flags ?? 0;
            } catch (err) {
              flags = null;
            }
            try {
              const reals = await api.get("/vips/reals", {
                address: vip.address,
                port: vip.port,
                proto: vip.proto,
              });
              healthy = Array.isArray(reals) && reals.some((real) => Boolean(real?.healthy));
            } catch (err) {
              healthy = false;
            }
            return { ...vip, flags, healthy };
          })
        );
        setStatus(lbStatus || { initialized: false, ready: false });
        setVips(vipsWithFlags);
        setError("");
      } catch (err) {
        setError(err.message || "request failed");
      }
    };

    useEffect(() => {
      let mounted = true;
      const loadSafe = async () => {
        if (!mounted) return;
        await load();
      };
      loadSafe();
      return () => {
        mounted = false;
      };
    }, []);

    const submitInit = async (event) => {
      event.preventDefault();
      try {
        const forwardingCores = parseNumberList(initForm.forwarding_cores, "Forwarding cores");
        const numaNodes = parseNumberList(initForm.numa_nodes, "NUMA nodes");
        const payload = {
          ...initForm,
          forwarding_cores: forwardingCores,
          numa_nodes: numaNodes,
          root_map_pos: initForm.root_map_pos === "" ? undefined : Number(initForm.root_map_pos),
          max_vips: Number(initForm.max_vips),
          max_reals: Number(initForm.max_reals),
          hash_function: initForm.hash_function,
        };
        await api.post("/lb/create", payload);
        setError("");
        setShowInit(false);
        addToast("Load balancer initialized.", "success");
        await load();
      } catch (err) {
        setError(err.message || "request failed");
        addToast(err.message || "Initialize failed.", "error");
      }
    };

    const submitVip = async (event) => {
      event.preventDefault();
      try {
        await api.post("/vips", {
          ...vipForm,
          port: Number(vipForm.port),
          proto: Number(vipForm.proto),
          flags: Number(vipForm.flags || 0),
        });
        setVipForm({ address: "", port: 80, proto: 6, flags: 0 });
        setError("");
        setShowVipForm(false);
        addToast("VIP created.", "success");
        await load();
      } catch (err) {
        setError(err.message || "request failed");
        addToast(err.message || "VIP create failed.", "error");
      }
    };

    const loadBpfProgs = async () => {
      try {
        await api.post("/lb/load-bpf-progs");
        setError("");
        addToast("BPF programs loaded.", "success");
        await load();
      } catch (err) {
        setError(err.message || "request failed");
        addToast(err.message || "Load BPF programs failed.", "error");
      }
    };

    const attachBpfProgs = async () => {
      try {
        await api.post("/lb/attach-bpf-progs");
        setError("");
        addToast("BPF programs attached.", "success");
        await load();
      } catch (err) {
        setError(err.message || "request failed");
        addToast(err.message || "Attach BPF programs failed.", "error");
      }
    };

    return html`
      <main>
        <section className="card">
          <h2>Load balancer</h2>
          <p>Initialized: ${status.initialized ? "yes" : "no"}</p>
          <p>Ready: ${status.ready ? "yes" : "no"}</p>
          <div className="row">
            ${!status.initialized &&
            html`
              <button className="btn" onClick=${() => setShowInit((s) => !s)}>
                ${showInit ? "Close" : "Initialize"}
              </button>
            `}
            <button className="btn secondary" onClick=${() => setShowVipForm((s) => !s)}>
              ${showVipForm ? "Close" : "Create VIP"}
            </button>
          </div>
          ${!status.ready &&
          html`
            <div className="row" style=${{ marginTop: 12 }}>
              <button
                className="btn ghost"
                disabled=${!status.initialized}
                onClick=${loadBpfProgs}
              >
                Load BPF Programs
              </button>
              <button
                className="btn ghost"
                disabled=${!status.initialized}
                onClick=${attachBpfProgs}
              >
                Attach BPF Programs
              </button>
            </div>
          `}
          ${showInit &&
          html`
            <form className="form" onSubmit=${submitInit}>
              <div className="form-row">
                <label className="field">
                  <span>Main interface</span>
                  <input
                    value=${initForm.main_interface}
                    onInput=${(e) =>
                      setInitForm({ ...initForm, main_interface: e.target.value })}
                    placeholder="eth0"
                    required
                  />
                </label>
                <label className="field">
                  <span>Balancer prog path</span>
                  <input
                    value=${initForm.balancer_prog_path}
                    onInput=${(e) =>
                      setInitForm({ ...initForm, balancer_prog_path: e.target.value })}
                    placeholder="/path/to/bpf.o"
                    required
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Healthchecking path</span>
                  <input
                    value=${initForm.healthchecking_prog_path}
                    onInput=${(e) =>
                      setInitForm({
                        ...initForm,
                        healthchecking_prog_path: e.target.value,
                      })}
                    placeholder="/path/to/hc.o"
                  />
                </label>
                <label className="field">
                  <span>Default MAC</span>
                  <input
                    value=${initForm.default_mac}
                    onInput=${(e) =>
                      setInitForm({ ...initForm, default_mac: e.target.value })}
                    placeholder="aa:bb:cc:dd:ee:ff"
                    required
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Local MAC</span>
                  <input
                    value=${initForm.local_mac}
                    onInput=${(e) =>
                      setInitForm({ ...initForm, local_mac: e.target.value })}
                    placeholder="11:22:33:44:55:66"
                    required=${initForm.healthchecking_prog_path?.trim() !== ""}
                  />
                </label>
                <label className="field">
                  <span>Hash function</span>
                  <select
                    value=${initForm.hash_function}
                    onInput=${(e) =>
                      setInitForm({ ...initForm, hash_function: e.target.value })}
                  >
                    <option value="maglev">maglev</option>
                    <option value="maglev_v2">maglev_v2</option>
                  </select>
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Root map path</span>
                  <input
                    value=${initForm.root_map_path}
                    onInput=${(e) =>
                      setInitForm({ ...initForm, root_map_path: e.target.value })}
                    placeholder="/sys/fs/bpf/root_map"
                  />
                </label>
                <label className="field">
                  <span>Root map position</span>
                  <input
                    type="number"
                    min="0"
                    value=${initForm.root_map_pos}
                    onInput=${(e) =>
                      setInitForm({ ...initForm, root_map_pos: e.target.value })}
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Katran src v4</span>
                  <input
                    value=${initForm.katran_src_v4}
                    onInput=${(e) =>
                      setInitForm({ ...initForm, katran_src_v4: e.target.value })}
                    placeholder="10.0.0.1"
                  />
                </label>
                <label className="field">
                  <span>Katran src v6</span>
                  <input
                    value=${initForm.katran_src_v6}
                    onInput=${(e) =>
                      setInitForm({ ...initForm, katran_src_v6: e.target.value })}
                    placeholder="fc00::1"
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Max VIPs</span>
                  <input
                    type="number"
                    min="1"
                    value=${initForm.max_vips}
                    onInput=${(e) =>
                      setInitForm({ ...initForm, max_vips: e.target.value })}
                  />
                </label>
                <label className="field">
                  <span>Max reals</span>
                  <input
                    type="number"
                    min="1"
                    value=${initForm.max_reals}
                    onInput=${(e) =>
                      setInitForm({ ...initForm, max_reals: e.target.value })}
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Forwarding cores (optional)</span>
                  <input
                    value=${initForm.forwarding_cores}
                    onInput=${(e) =>
                      setInitForm({ ...initForm, forwarding_cores: e.target.value })}
                    placeholder="0,1,2,3"
                  />
                  <span className="muted">Comma or space separated CPU core IDs.</span>
                </label>
                <label className="field">
                  <span>NUMA nodes (optional)</span>
                  <input
                    value=${initForm.numa_nodes}
                    onInput=${(e) => setInitForm({ ...initForm, numa_nodes: e.target.value })}
                    placeholder="0,0,1,1"
                  />
                  <span className="muted">Match the forwarding cores length.</span>
                </label>
              </div>
              <label className="field checkbox">
                <input
                  type="checkbox"
                  checked=${initForm.use_root_map}
                  onChange=${(e) =>
                    setInitForm({ ...initForm, use_root_map: e.target.checked })}
                />
                <span>Use root map</span>
              </label>
              <button className="btn" type="submit">Initialize LB</button>
            </form>
          `}
          ${showVipForm &&
          html`
            <form className="form" onSubmit=${submitVip}>
              <div className="form-row">
                <label className="field">
                  <span>VIP address</span>
                  <input
                    value=${vipForm.address}
                    onInput=${(e) => setVipForm({ ...vipForm, address: e.target.value })}
                    placeholder="1.2.3.4"
                    required
                  />
                </label>
                <label className="field">
                  <span>Port</span>
                  <input
                    type="number"
                    value=${vipForm.port}
                    onInput=${(e) => setVipForm({ ...vipForm, port: e.target.value })}
                    required
                  />
                </label>
                <label className="field">
                  <span>Protocol</span>
                  <select
                    value=${vipForm.proto}
                    onChange=${(e) => setVipForm({ ...vipForm, proto: e.target.value })}
                  >
                    <option value="6">TCP (6)</option>
                    <option value="17">UDP (17)</option>
                  </select>
                </label>
                <label className="field">
                  <span>Flags</span>
                  <${FlagSelector}
                    options=${VIP_FLAG_OPTIONS}
                    value=${vipForm.flags}
                    name="vip-add"
                    onChange=${(next) => setVipForm({ ...vipForm, flags: next })}
                  />
                </label>
              </div>
              <button className="btn" type="submit">Create VIP</button>
            </form>
          `}
          ${error && html`<p className="error">${error}</p>`}
        </section>
        <section className="card">
          <div className="section-header">
            <h2>VIPs</h2>
            <button className="btn ghost" onClick=${load}>Refresh</button>
          </div>
          ${vips.length === 0
            ? html`<p className="muted">No VIPs configured yet.</p>`
            : html`
                <div className="grid">
                  ${vips.map(
                    (vip) => html`
                      <div className="card">
                        <div className="row" style=${{ fontWeight: 600, gap: 8 }}>
                          <span className=${`dot ${vip.healthy ? "ok" : "bad"}`}></span>
                          ${vip.address}:${vip.port} / ${vip.proto}
                        </div>
                        <div className="muted" style=${{ marginTop: 6 }}>
                          <div style=${{ fontWeight: 600, marginBottom: 6 }}>Flags</div>
                          <${FlagTable}
                            mask=${vip.flags}
                            options=${VIP_FLAG_OPTIONS}
                            emptyLabel=${vip.flags === null ? "Unknown" : "No flags"}
                          />
                        </div>
                        <div className="row" style=${{ marginTop: 12 }}>
                          <${Link} className="btn" to=${`/vips/${vipIdFromVip(vip)}`}>
                            Open
                          </${Link}>
                          <${Link}
                            className="btn secondary"
                            to=${`/vips/${vipIdFromVip(vip)}/stats`}
                          >
                            Stats
                          </${Link}>
                        </div>
                      </div>
                    `
                  )}
                </div>
              `}
        </section>
      </main>
    `;
  }

  function VipDetail() {
    const { addToast } = useToast();
    const params = useParams();
    const navigate = useNavigate();
    const vip = useMemo(() => parseVipId(params.vipId), [params.vipId]);
    const [reals, setReals] = useState([]);
    const [error, setError] = useState("");
    const [flagError, setFlagError] = useState("");
    const [loading, setLoading] = useState(true);
    const [newReal, setNewReal] = useState({ address: "", weight: 100, flags: 0 });
    const [weights, setWeights] = useState({});
    const [vipFlags, setVipFlags] = useState(null);
    const [flagForm, setFlagForm] = useState({ flags: 0, set: true });
    const [hashForm, setHashForm] = useState({ hash_function: 0 });
    const { groups, setGroups, refreshFromStorage, importFromRunningConfig } = useTargetGroups();
    const [targetGroupName, setTargetGroupName] = useState("");
    const [targetGroupError, setTargetGroupError] = useState("");
    const [targetGroupBusy, setTargetGroupBusy] = useState(false);
    const [saveGroupName, setSaveGroupName] = useState("");
    const [groupDiff, setGroupDiff] = useState({ add: 0, update: 0, remove: 0 });

    const loadReals = async () => {
      try {
        const list = await api.get("/vips/reals", vip);
        setReals(list || []);
        const nextWeights = {};
        (list || []).forEach((real) => {
          nextWeights[real.address] = real.weight;
        });
        setWeights(nextWeights);
        setError("");
        setLoading(false);
      } catch (err) {
        setError(err.message || "request failed");
        setLoading(false);
      }
    };

    const loadFlags = async () => {
      try {
        const data = await api.get("/vips/flags", vip);
        setVipFlags(data?.flags ?? 0);
        setFlagError("");
      } catch (err) {
        setFlagError(err.message || "request failed");
      }
    };

    useEffect(() => {
      loadReals();
      loadFlags();
    }, [params.vipId]);

    useEffect(() => {
      if (!targetGroupName) {
        setGroupDiff({ add: 0, update: 0, remove: 0 });
        return;
      }
      const group = groups[targetGroupName] || [];
      const currentByAddress = new Map(reals.map((real) => [real.address, real]));
      const targetByAddress = new Map(group.map((real) => [real.address, real]));
      let add = 0;
      let update = 0;
      let remove = 0;
      group.forEach((real) => {
        const existing = currentByAddress.get(real.address);
        if (!existing) {
          add += 1;
          return;
        }
        if (
          Number(existing.weight) !== Number(real.weight) ||
          Number(existing.flags || 0) !== Number(real.flags || 0)
        ) {
          update += 1;
        }
      });
      reals.forEach((real) => {
        if (!targetByAddress.has(real.address)) {
          remove += 1;
        }
      });
      setGroupDiff({ add, update, remove });
    }, [targetGroupName, reals, groups]);

    const updateWeight = async (real) => {
      try {
        const weight = Number(weights[real.address]);
        await api.post("/vips/reals", {
          vip,
          real: { address: real.address, weight, flags: real.flags || 0 },
        });
        await loadReals();
        addToast("Real weight updated.", "success");
      } catch (err) {
        setError(err.message || "request failed");
        addToast(err.message || "Update failed.", "error");
      }
    };

    const deleteReal = async (real) => {
      try {
        await api.del("/vips/reals", {
          vip,
          real: { address: real.address, weight: real.weight, flags: real.flags || 0 },
        });
        await loadReals();
        addToast("Real removed.", "success");
      } catch (err) {
        setError(err.message || "request failed");
        addToast(err.message || "Remove failed.", "error");
      }
    };

    const addReal = async (event) => {
      event.preventDefault();
      try {
        await api.post("/vips/reals", {
          vip,
          real: {
            address: newReal.address,
            weight: Number(newReal.weight),
            flags: Number(newReal.flags || 0),
          },
        });
        setNewReal({ address: "", weight: 100, flags: 0 });
        await loadReals();
        addToast("Real added.", "success");
      } catch (err) {
        setError(err.message || "request failed");
        addToast(err.message || "Add failed.", "error");
      }
    };

    const applyTargetGroup = async () => {
      if (!targetGroupName || !groups[targetGroupName]) {
        setTargetGroupError("Select a target group to apply.");
        return;
      }
      setTargetGroupBusy(true);
      setTargetGroupError("");
      const group = groups[targetGroupName] || [];
      const currentByAddress = new Map(reals.map((real) => [real.address, real]));
      const targetByAddress = new Map(group.map((real) => [real.address, real]));
      const toRemove = reals.filter((real) => !targetByAddress.has(real.address));
      const toUpsert = group.filter((real) => {
        const current = currentByAddress.get(real.address);
        if (!current) return true;
        return (
          Number(current.weight) !== Number(real.weight) ||
          Number(current.flags || 0) !== Number(real.flags || 0)
        );
      });

      try {
        if (toRemove.length > 0) {
          await api.put("/vips/reals/batch", {
            vip,
            action: 1,
            reals: toRemove.map((real) => ({
              address: real.address,
              weight: Number(real.weight),
              flags: Number(real.flags || 0),
            })),
          });
        }
        if (toUpsert.length > 0) {
          await Promise.all(
            toUpsert.map((real) =>
              api.post("/vips/reals", {
                vip,
                real: {
                  address: real.address,
                  weight: Number(real.weight),
                  flags: Number(real.flags || 0),
                },
              })
            )
          );
        }
        await loadReals();
        addToast(`Applied target group "${targetGroupName}".`, "success");
      } catch (err) {
        setTargetGroupError(err.message || "Failed to apply target group.");
        addToast(err.message || "Target group apply failed.", "error");
      } finally {
        setTargetGroupBusy(false);
      }
    };

    const saveCurrentAsGroup = (event) => {
      event.preventDefault();
      const name = saveGroupName.trim();
      if (!name) {
        setTargetGroupError("Provide a name for the new target group.");
        return;
      }
      if (groups[name]) {
        setTargetGroupError("A target group with that name already exists.");
        return;
      }
      const next = {
        ...groups,
        [name]: reals.map((real) => ({
          address: real.address,
          weight: Number(real.weight),
          flags: Number(real.flags || 0),
        })),
      };
      setGroups(next);
      setSaveGroupName("");
      setTargetGroupName(name);
      setTargetGroupError("");
      addToast(`Target group "${name}" saved.`, "success");
    };

    const deleteVip = async () => {
      try {
        await api.del("/vips", vip);
        addToast("VIP deleted.", "success");
        navigate("/");
      } catch (err) {
        setError(err.message || "request failed");
        addToast(err.message || "Delete failed.", "error");
      }
    };

    const submitFlagChange = async (event) => {
      event.preventDefault();
      try {
        await api.put("/vips/flags", {
          ...vip,
          flag: Number(flagForm.flags || 0),
          set: Boolean(flagForm.set),
        });
        await loadFlags();
        addToast("VIP flags updated.", "success");
      } catch (err) {
        setFlagError(err.message || "request failed");
        addToast(err.message || "Flag update failed.", "error");
      }
    };

    const submitHashChange = async (event) => {
      event.preventDefault();
      try {
        await api.put("/vips/hash-function", {
          ...vip,
          hash_function: Number(hashForm.hash_function),
        });
        addToast("Hash function updated.", "success");
      } catch (err) {
        setFlagError(err.message || "request failed");
        addToast(err.message || "Hash update failed.", "error");
      }
    };

    return html`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>VIP Detail</h2>
              <p className="muted">${vip.address}:${vip.port} / ${vip.proto}</p>
              ${vipFlags === null
                ? html`<p className="muted">Flags: —</p>`
                : html`
                    <div style=${{ marginTop: 8 }}>
                      <${FlagTable}
                        mask=${vipFlags}
                        options=${VIP_FLAG_OPTIONS}
                        showStatus=${true}
                        emptyLabel="No flags"
                      />
                    </div>
                  `}
            </div>
            <div className="row">
              <button className="btn ghost" onClick=${loadFlags}>Refresh flags</button>
              <button className="btn danger" onClick=${deleteVip}>Delete VIP</button>
            </div>
          </div>
          ${error && html`<p className="error">${error}</p>`}
          ${flagError && html`<p className="error">${flagError}</p>`}
          ${loading
            ? html`<p className="muted">Loading reals…</p>`
            : html`
                <table className="table">
                  <thead>
                    <tr>
                      <th>Status</th>
                      <th>Address</th>
                      <th>Weight</th>
                      <th></th>
                    </tr>
                  </thead>
                  <tbody>
                    ${reals.map(
                      (real) => html`
                        <tr>
                          <td>
                            <span
                              className=${`dot ${Number(real.weight) > 0 ? "ok" : "bad"}`}
                            ></span>
                          </td>
                          <td>${real.address}</td>
                          <td>
                            <input
                              className="inline-input"
                              type="number"
                              min="0"
                              value=${weights[real.address] ?? real.weight}
                              onInput=${(e) =>
                                setWeights({
                                  ...weights,
                                  [real.address]: e.target.value,
                                })}
                            />
                          </td>
                          <td className="row">
                            <button className="btn" onClick=${() => updateWeight(real)}>
                              Update
                            </button>
                            <button className="btn ghost" onClick=${() => deleteReal(real)}>
                              Remove
                            </button>
                          </td>
                        </tr>
                      `
                    )}
                  </tbody>
                </table>
              `}
        </section>
        <section className="card">
          <h3>VIP flags & hash</h3>
          <div className="grid">
            <form className="form" onSubmit=${submitFlagChange}>
              <div className="form-row">
                <label className="field">
                  <span>Flags</span>
                  <${FlagSelector}
                    options=${VIP_FLAG_OPTIONS}
                    value=${flagForm.flags}
                    name="vip-flag-change"
                    onChange=${(next) => setFlagForm({ ...flagForm, flags: next })}
                  />
                </label>
                <label className="field">
                  <span>Set</span>
                  <select
                    value=${String(flagForm.set)}
                    onChange=${(e) =>
                      setFlagForm({ ...flagForm, set: e.target.value === "true" })}
                  >
                    <option value="true">Set</option>
                    <option value="false">Clear</option>
                  </select>
                </label>
              </div>
              <button className="btn" type="submit">Apply flag</button>
            </form>
            <form className="form" onSubmit=${submitHashChange}>
              <div className="form-row">
                <label className="field">
                  <span>Hash function</span>
                  <input
                    type="number"
                    min="0"
                    value=${hashForm.hash_function}
                    onInput=${(e) =>
                      setHashForm({ ...hashForm, hash_function: e.target.value })}
                  />
                </label>
              </div>
              <button className="btn" type="submit">Apply hash</button>
            </form>
          </div>
        </section>
        <section className="card">
          <h3>Add real</h3>
          <form className="form" onSubmit=${addReal}>
            <div className="form-row">
              <label className="field">
                <span>Real address</span>
                <input
                  value=${newReal.address}
                  onInput=${(e) => setNewReal({ ...newReal, address: e.target.value })}
                  placeholder="10.0.0.1"
                  required
                />
              </label>
              <label className="field">
                <span>Weight</span>
                <input
                  type="number"
                  min="0"
                  value=${newReal.weight}
                  onInput=${(e) => setNewReal({ ...newReal, weight: e.target.value })}
                  required
                />
              </label>
            </div>
            <button className="btn" type="submit">Add real</button>
          </form>
        </section>
        <section className="card">
          <div className="section-header">
            <div>
              <h3>Target group</h3>
              <p className="muted">Sync this VIP with a saved target group of reals.</p>
            </div>
            <div className="row">
              <button className="btn ghost" type="button" onClick=${refreshFromStorage}>
                Reload groups
              </button>
              <button
                className="btn ghost"
                type="button"
                onClick=${async () => {
                  try {
                    await importFromRunningConfig();
                    addToast("Imported target groups from running config.", "success");
                  } catch (err) {
                    setTargetGroupError(err.message || "Failed to import target groups.");
                    addToast(err.message || "Import failed.", "error");
                  }
                }}
              >
                Import from running config
              </button>
            </div>
          </div>
          ${targetGroupError && html`<p className="error">${targetGroupError}</p>`}
          <div className="form-row">
            <label className="field">
              <span>Target group</span>
              <select
                value=${targetGroupName}
                onChange=${(e) => setTargetGroupName(e.target.value)}
                disabled=${Object.keys(groups).length === 0}
              >
                <option value="">Select group</option>
                ${Object.keys(groups).map(
                  (name) => html`<option value=${name}>${name}</option>`
                )}
              </select>
            </label>
            <label className="field">
              <span>Preview</span>
              <input
                value=${`add ${groupDiff.add} · update ${groupDiff.update} · remove ${groupDiff.remove}`}
                readOnly
              />
            </label>
          </div>
          <div className="row">
            <button
              className="btn"
              type="button"
              onClick=${applyTargetGroup}
              disabled=${targetGroupBusy || !targetGroupName}
            >
              ${targetGroupBusy ? "Applying..." : "Apply target group"}
            </button>
          </div>
          <form className="form" onSubmit=${saveCurrentAsGroup}>
            <div className="form-row">
              <label className="field">
                <span>Save current reals as new group</span>
                <input
                  value=${saveGroupName}
                  onInput=${(e) => setSaveGroupName(e.target.value)}
                  placeholder="edge-backends"
                />
              </label>
            </div>
            <button className="btn secondary" type="submit">Save target group</button>
          </form>
        </section>
      </main>
    `;
  }

  function VipStats() {
    const params = useParams();
    const vip = useMemo(() => parseVipId(params.vipId), [params.vipId]);
    const { points, error } = useStatSeries({ path: "/stats/vip", body: vip });
    const latest = points[points.length - 1] || {};
    const prev = points[points.length - 2] || {};
    const v1 = Number(latest.v1 ?? 0);
    const v2 = Number(latest.v2 ?? 0);
    const d1 = v1 - Number(prev.v1 ?? 0);
    const d2 = v2 - Number(prev.v2 ?? 0);
    const keys = useMemo(
      () => [
        { label: "v1", field: "v1", color: "#2f4858", fill: "rgba(47,72,88,0.15)" },
        { label: "v2", field: "v2", color: "#d97757", fill: "rgba(217,119,87,0.2)" },
      ],
      []
    );

    return html`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>VIP Stats</h2>
              <p className="muted">${vip.address}:${vip.port} / ${vip.proto}</p>
            </div>
          </div>
          ${error && html`<p className="error">${error}</p>`}
          <table className="table">
            <thead>
              <tr>
                <th>Counter</th>
                <th>Absolute</th>
                <th>Delta/sec</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>v1</td>
                <td>${v1}</td>
                <td>
                  <span className=${`delta ${d1 < 0 ? "down" : "up"}`}>
                    ${formatDelta(d1)}
                  </span>
                </td>
              </tr>
              <tr>
                <td>v2</td>
                <td>${v2}</td>
                <td>
                  <span className=${`delta ${d2 < 0 ? "down" : "up"}`}>
                    ${formatDelta(d2)}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
          <${StatChart} title="Traffic (delta/sec)" points=${points} keys=${keys} diff=${true} />
        </section>
      </main>
    `;
  }

  const GLOBAL_STAT_ITEMS = [
    { title: "LRU", path: "/stats/lru" },
    { title: "LRU Miss", path: "/stats/lru/miss" },
    { title: "LRU Fallback", path: "/stats/lru/fallback" },
    { title: "LRU Global", path: "/stats/lru/global" },
    { title: "XDP Total", path: "/stats/xdp/total" },
    { title: "XDP Pass", path: "/stats/xdp/pass" },
    { title: "XDP Drop", path: "/stats/xdp/drop" },
    { title: "XDP Tx", path: "/stats/xdp/tx" },
  ];

  function formatDelta(value) {
    if (!Number.isFinite(value)) return "0";
    const sign = value > 0 ? "+" : "";
    return `${sign}${value}`;
  }

  function StatPanel({ title, path, diff = false }) {
    const { points, error } = useStatSeries({ path });
    const keys = useMemo(
      () => [
        { label: "v1", field: "v1", color: "#2f4858", fill: "rgba(47,72,88,0.2)" },
        { label: "v2", field: "v2", color: "#d97757", fill: "rgba(217,119,87,0.2)" },
      ],
      []
    );

    return html`
      <div className="card">
        <h3>${title}</h3>
        ${error && html`<p className="error">${error}</p>`}
        <${StatChart} title=${title} points=${points} keys=${keys} diff=${diff} inlineTitle=${false} />
      </div>
    `;
  }

  function StatSummaryCard({ title, path }) {
    const { points, error } = useStatSeries({ path });
    const latest = points[points.length - 1] || {};
    const prev = points[points.length - 2] || {};
    const v1 = Number(latest.v1 ?? 0);
    const v2 = Number(latest.v2 ?? 0);
    const d1 = v1 - Number(prev.v1 ?? 0);
    const d2 = v2 - Number(prev.v2 ?? 0);

    return html`
      <div className="summary-card">
        <div className="summary-title">${title}</div>
        ${error
          ? html`<p className="error">${error}</p>`
          : html`
              <div className="summary-row">
                <div className="stat">
                  <span className="muted">v1 absolute</span>
                  <strong>${v1}</strong>
                </div>
                <div className="stat">
                  <span className="muted">v1 delta/sec</span>
                  <strong className=${d1 < 0 ? "delta down" : "delta up"}>
                    ${formatDelta(d1)}
                  </strong>
                </div>
              </div>
              <div className="summary-row">
                <div className="stat">
                  <span className="muted">v2 absolute</span>
                  <strong>${v2}</strong>
                </div>
                <div className="stat">
                  <span className="muted">v2 delta/sec</span>
                  <strong className=${d2 < 0 ? "delta down" : "delta up"}>
                    ${formatDelta(d2)}
                  </strong>
                </div>
              </div>
            `}
      </div>
    `;
  }

  function GlobalStats() {
    const { data: userspace, error: userspaceError } = usePolling(
      () => api.get("/stats/userspace"),
      1000,
      []
    );

    return html`
      <main>
        <section className="card">
          <h2>Global Stats</h2>
          <p className="muted">Polling every second. Charts show per-second deltas.</p>
        </section>
        <section className="grid">
          ${GLOBAL_STAT_ITEMS.map(
            (item) => html`<${StatPanel} title=${item.title} path=${item.path} diff=${true} />`
          )}
        </section>
        <section className="card">
          <h3>Userspace</h3>
          ${userspaceError && html`<p className="error">${userspaceError}</p>`}
          ${userspace
            ? html`
                <div className="row">
                  <div className="stat">
                    <span className="muted">BPF failed calls</span>
                    <strong>${userspace.bpf_failed_calls ?? 0}</strong>
                  </div>
                  <div className="stat">
                    <span className="muted">Addr validation failed</span>
                    <strong>${userspace.addr_validation_failed ?? 0}</strong>
                  </div>
                </div>
              `
            : html`<p className="muted">Waiting for data…</p>`}
        </section>
        <section className="card">
          <div className="section-header">
            <div>
              <h3>Absolute & Rate of Change</h3>
              <p className="muted">Latest value and per-second delta.</p>
            </div>
          </div>
          <div className="summary-grid">
            ${GLOBAL_STAT_ITEMS.map(
              (item) => html`<${StatSummaryCard} title=${item.title} path=${item.path} />`
            )}
          </div>
        </section>
      </main>
    `;
  }

  function RealStats() {
    const [vips, setVips] = useState([]);
    const [vipId, setVipId] = useState("");
    const [reals, setReals] = useState([]);
    const [realAddress, setRealAddress] = useState("");
    const [realIndex, setRealIndex] = useState(null);
    const [error, setError] = useState("");

    useEffect(() => {
      let mounted = true;
      const load = async () => {
        try {
          const list = await api.get("/vips");
          if (!mounted) return;
          setVips(list || []);
          if (!vipId && list && list.length > 0) {
            setVipId(vipIdFromVip(list[0]));
          }
        } catch (err) {
          if (mounted) {
            setError(err.message || "request failed");
          }
        }
      };
      load();
      return () => {
        mounted = false;
      };
    }, []);

    useEffect(() => {
      if (!vipId) return;
      const vip = parseVipId(vipId);
      let mounted = true;
      const load = async () => {
        try {
          const list = await api.get("/vips/reals", vip);
          if (!mounted) return;
          setReals(list || []);
          if (list && list.length > 0) {
            setRealAddress((prev) => prev || list[0].address);
          } else {
            setRealAddress("");
          }
          setError("");
        } catch (err) {
          if (mounted) {
            setError(err.message || "request failed");
          }
        }
      };
      load();
      return () => {
        mounted = false;
      };
    }, [vipId]);

    useEffect(() => {
      if (!realAddress) {
        setRealIndex(null);
        return;
      }
      let mounted = true;
      const load = async () => {
        try {
          const data = await api.get("/reals/index", { address: realAddress });
          if (!mounted) return;
          setRealIndex(data?.index ?? null);
          setError("");
        } catch (err) {
          if (mounted) {
            setError(err.message || "request failed");
          }
        }
      };
      load();
      return () => {
        mounted = false;
      };
    }, [realAddress]);

    const { points, error: statsError } = useStatSeries({
      path: "/stats/real",
      body: realIndex !== null ? { index: realIndex } : null,
    });

    const labels = useMemo(() => getStatLabels("/stats/real"), []);
    const keys = useMemo(
      () => [
        { label: labels.v1, field: "v1", color: "#2f4858", fill: "rgba(47,72,88,0.2)" },
        { label: labels.v2, field: "v2", color: "#d97757", fill: "rgba(217,119,87,0.2)" },
      ],
      [labels]
    );
    const latest = points[points.length - 1] || {};
    const prev = points[points.length - 2] || {};
    const v1 = Number(latest.v1 ?? 0);
    const v2 = Number(latest.v2 ?? 0);
    const d1 = v1 - Number(prev.v1 ?? 0);
    const d2 = v2 - Number(prev.v2 ?? 0);

    return html`
      <main>
        <section className="card">
          <h2>Per-Real Stats</h2>
          <p className="muted">Select a VIP and real address to chart.</p>
          ${error && html`<p className="error">${error}</p>`}
          <div className="form-row">
            <label className="field">
              <span>VIP</span>
              <select value=${vipId} onChange=${(e) => setVipId(e.target.value)}>
                ${vips.map(
                  (vip) => html`
                    <option value=${vipIdFromVip(vip)}>
                      ${vip.address}:${vip.port} / ${vip.proto}
                    </option>
                  `
                )}
              </select>
            </label>
            <label className="field">
              <span>Real</span>
              <select
                value=${realAddress}
                onChange=${(e) => setRealAddress(e.target.value)}
                disabled=${reals.length === 0}
              >
                ${reals.map(
                  (real) => html`
                    <option value=${real.address}>${real.address}</option>
                  `
                )}
              </select>
            </label>
            <label className="field">
              <span>Index</span>
              <input value=${realIndex ?? ""} readOnly />
            </label>
          </div>
        </section>
        <section className="card">
          <h3>Real stats</h3>
          ${statsError && html`<p className="error">${statsError}</p>`}
          ${realIndex === null
            ? html`<p className="muted">Select a real to start polling.</p>`
            : html`
                <table className="table">
                  <thead>
                    <tr>
                      <th>Counter</th>
                      <th>Absolute</th>
                      <th>Delta/sec</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr>
                      <td>${labels.v1}</td>
                      <td>${v1}</td>
                      <td>
                        <span className=${`delta ${d1 < 0 ? "down" : "up"}`}>
                          ${formatDelta(d1)}
                        </span>
                      </td>
                    </tr>
                    <tr>
                      <td>${labels.v2}</td>
                      <td>${v2}</td>
                      <td>
                        <span className=${`delta ${d2 < 0 ? "down" : "up"}`}>
                          ${formatDelta(d2)}
                        </span>
                      </td>
                    </tr>
                  </tbody>
                </table>
                <${StatChart}
                  title="Traffic (delta/sec)"
                  points=${points}
                  keys=${keys}
                  diff=${true}
                />
              `}
        </section>
      </main>
    `;
  }

  function ConfigExport() {
    const { addToast } = useToast();
    const [yaml, setYaml] = useState("");
    const [error, setError] = useState("");
    const [loading, setLoading] = useState(true);
    const [lastFetched, setLastFetched] = useState("");
    const mountedRef = useRef(true);

    const load = async () => {
      if (!mountedRef.current) return;
      setLoading(true);
      setError("");
      try {
        const res = await fetch(`${api.base}/config/export`, {
          headers: { Accept: "application/x-yaml" },
        });
        if (!res.ok) {
          let message = `HTTP ${res.status}`;
          try {
            const payload = await res.json();
            message = payload?.error?.message || message;
          } catch (err) {
            // ignore JSON parsing errors for YAML responses
          }
          throw new Error(message);
        }
        const text = await res.text();
        if (!mountedRef.current) return;
        setYaml(text || "");
        setLastFetched(new Date().toLocaleString());
      } catch (err) {
        if (mountedRef.current) {
          setError(err.message || "request failed");
        }
      } finally {
        if (mountedRef.current) {
          setLoading(false);
        }
      }
    };

    const copyToClipboard = async () => {
      if (!yaml) return;
      try {
        await navigator.clipboard.writeText(yaml);
        addToast("Config copied to clipboard", "info");
      } catch (err) {
        addToast("Failed to copy config", "error");
      }
    };

    useEffect(() => {
      mountedRef.current = true;
      load();
      return () => {
        mountedRef.current = false;
      };
    }, []);

    return html`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>Running config</h2>
              <p className="muted">Exported from /api/v1/config/export</p>
            </div>
            <div className="row">
              <button className="btn ghost" onClick=${copyToClipboard} disabled=${!yaml}>
                Copy YAML
              </button>
              <button className="btn" onClick=${load} disabled=${loading}>
                Refresh
              </button>
            </div>
          </div>
          ${error && html`<p className="error">${error}</p>`}
          ${loading
            ? html`<p className="muted">Loading config...</p>`
            : yaml
              ? html`<pre className="yaml-view">${yaml}</pre>`
              : html`<p className="muted">No config data returned.</p>`}
          ${lastFetched && html`<p className="muted">Last fetched ${lastFetched}</p>`}
        </section>
      </main>
    `;
  }

  function TargetGroups() {
    const { addToast } = useToast();
    const { groups, setGroups, refreshFromStorage, importFromRunningConfig } = useTargetGroups();
    const [groupName, setGroupName] = useState("");
    const [selectedGroup, setSelectedGroup] = useState("");
    const [newReal, setNewReal] = useState({ address: "", weight: 100, flags: 0 });
    const [error, setError] = useState("");
    const [importing, setImporting] = useState(false);

    useEffect(() => {
      if (!selectedGroup) {
        const names = Object.keys(groups);
        if (names.length > 0) {
          setSelectedGroup(names[0]);
        }
      } else if (!groups[selectedGroup]) {
        const names = Object.keys(groups);
        setSelectedGroup(names[0] || "");
      }
    }, [groups, selectedGroup]);

    const createGroup = (event) => {
      event.preventDefault();
      const name = groupName.trim();
      if (!name) {
        setError("Provide a group name.");
        return;
      }
      if (groups[name]) {
        setError("That group already exists.");
        return;
      }
      setGroups({ ...groups, [name]: [] });
      setGroupName("");
      setSelectedGroup(name);
      setError("");
      addToast(`Target group "${name}" created.`, "success");
    };

    const deleteGroup = (name) => {
      const next = { ...groups };
      delete next[name];
      setGroups(next);
      addToast(`Target group "${name}" removed.`, "success");
    };

    const addRealToGroup = (event) => {
      event.preventDefault();
      if (!selectedGroup) {
        setError("Select a group to add a real.");
        return;
      }
      const normalized = normalizeReal(newReal);
      if (!normalized) {
        setError("Provide a valid real address.");
        return;
      }
      const group = groups[selectedGroup] || [];
      const nextGroup = group.some((real) => real.address === normalized.address)
        ? group.map((real) =>
            real.address === normalized.address ? normalized : real
          )
        : group.concat(normalized);
      setGroups({ ...groups, [selectedGroup]: nextGroup });
      setNewReal({ address: "", weight: 100, flags: 0 });
      setError("");
      addToast("Real saved to target group.", "success");
    };

    const removeRealFromGroup = (address) => {
      if (!selectedGroup) return;
      const group = groups[selectedGroup] || [];
      const nextGroup = group.filter((real) => real.address !== address);
      setGroups({ ...groups, [selectedGroup]: nextGroup });
    };

    const updateGroupReal = (address, updates) => {
      if (!selectedGroup) return;
      const group = groups[selectedGroup] || [];
      const nextGroup = group.map((real) =>
        real.address === address ? { ...real, ...updates } : real
      );
      setGroups({ ...groups, [selectedGroup]: nextGroup });
    };

    const importGroups = async () => {
      setImporting(true);
      try {
        await importFromRunningConfig();
        addToast("Imported target groups from running config.", "success");
        setError("");
      } catch (err) {
        setError(err.message || "Failed to import target groups.");
        addToast(err.message || "Import failed.", "error");
      } finally {
        setImporting(false);
      }
    };

    return html`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>Target groups</h2>
              <p className="muted">Define reusable sets of reals (address + weight).</p>
            </div>
            <div className="row">
              <button className="btn ghost" type="button" onClick=${refreshFromStorage}>
                Reload groups
              </button>
              <button className="btn ghost" type="button" onClick=${importGroups} disabled=${importing}>
                ${importing ? "Importing..." : "Import from running config"}
              </button>
            </div>
          </div>
          ${error && html`<p className="error">${error}</p>`}
          <form className="form" onSubmit=${createGroup}>
            <div className="form-row">
              <label className="field">
                <span>New group name</span>
                <input
                  value=${groupName}
                  onInput=${(e) => setGroupName(e.target.value)}
                  placeholder="edge-backends"
                />
              </label>
            </div>
            <button className="btn" type="submit">Create group</button>
          </form>
        </section>
        <section className="card">
          <div className="section-header">
            <div>
              <h3>Group contents</h3>
              <p className="muted">
                Add or update real entries. Existing addresses are updated in-place.
              </p>
            </div>
            <div className="row">
              <label className="field">
                <span>Selected group</span>
                <select
                  value=${selectedGroup}
                  onChange=${(e) => setSelectedGroup(e.target.value)}
                  disabled=${Object.keys(groups).length === 0}
                >
                  ${Object.keys(groups).map(
                    (name) => html`<option value=${name}>${name}</option>`
                  )}
                </select>
              </label>
              ${selectedGroup &&
              html`<button className="btn danger" type="button" onClick=${() => deleteGroup(selectedGroup)}>
                Delete group
              </button>`}
            </div>
          </div>
          ${!selectedGroup
            ? html`<p className="muted">No groups yet. Create one to add reals.</p>`
            : html`
                <table className="table">
                  <thead>
                    <tr>
                      <th>Address</th>
                      <th>Weight</th>
                      <th></th>
                    </tr>
                  </thead>
                  <tbody>
                    ${(groups[selectedGroup] || []).map(
                      (real) => html`
                        <tr>
                          <td>${real.address}</td>
                          <td>
                            <input
                              className="inline-input"
                              type="number"
                              min="0"
                              value=${real.weight}
                              onInput=${(e) =>
                                updateGroupReal(real.address, {
                                  weight: Number(e.target.value),
                                })}
                            />
                          </td>
                          <td>
                            <button
                              className="btn ghost"
                              type="button"
                              onClick=${() => removeRealFromGroup(real.address)}
                            >
                              Remove
                            </button>
                          </td>
                        </tr>
                      `
                    )}
                  </tbody>
                </table>
                <form className="form" onSubmit=${addRealToGroup}>
                  <div className="form-row">
                    <label className="field">
                      <span>Real address</span>
                      <input
                        value=${newReal.address}
                        onInput=${(e) => setNewReal({ ...newReal, address: e.target.value })}
                        placeholder="10.0.0.1"
                        required
                      />
                    </label>
                    <label className="field">
                      <span>Weight</span>
                      <input
                        type="number"
                        min="0"
                        value=${newReal.weight}
                        onInput=${(e) => setNewReal({ ...newReal, weight: e.target.value })}
                        required
                      />
                    </label>
                  </div>
                  <button className="btn secondary" type="submit">Save real</button>
                </form>
              `}
        </section>
      </main>
    `;
  }

  function App() {
    const [status, setStatus] = useState({ initialized: false, ready: false });
    const [toasts, setToasts] = useState([]);
    const toastTimers = useRef({});

    const addToast = (message, kind = "info") => {
      const id = `${Date.now()}-${Math.random().toString(16).slice(2)}`;
      setToasts((prev) => prev.concat({ id, message, kind }));
      toastTimers.current[id] = setTimeout(() => {
        setToasts((prev) => prev.filter((toast) => toast.id !== id));
        delete toastTimers.current[id];
      }, 4000);
    };

    const dismissToast = (id) => {
      if (toastTimers.current[id]) {
        clearTimeout(toastTimers.current[id]);
        delete toastTimers.current[id];
      }
      setToasts((prev) => prev.filter((toast) => toast.id !== id));
    };

    useEffect(() => {
      let mounted = true;
      const load = async () => {
        try {
          const lbStatus = await api.get("/lb/status");
          if (mounted) {
            setStatus(lbStatus || { initialized: false, ready: false });
          }
        } catch (err) {
          if (mounted) {
            setStatus({ initialized: false, ready: false });
          }
        }
      };
      load();
      const interval = setInterval(load, 5000);
      return () => {
        mounted = false;
        clearInterval(interval);
      };
    }, []);

    return html`
      <${BrowserRouter}>
        <${AuthGate}>
          <${ToastContext.Provider} value=${{ addToast }}>
            <${Header} status=${status} />
            <${Routes}>
              <${Route} path="/" element=${html`<${Dashboard} />`} />
              <${Route} path="/vips/:vipId" element=${html`<${VipDetail} />`} />
              <${Route} path="/vips/:vipId/stats" element=${html`<${VipStats} />`} />
              <${Route} path="/target-groups" element=${html`<${TargetGroups} />`} />
              <${Route} path="/stats/global" element=${html`<${GlobalStats} />`} />
              <${Route} path="/stats/real" element=${html`<${RealStats} />`} />
              <${Route} path="/config" element=${html`<${ConfigExport} />`} />
            </${Routes}>
            <${Toasts} toasts=${toasts} onDismiss=${dismissToast} />
          </${ToastContext.Provider}>
        </${AuthGate}>
      </${BrowserRouter}>
    `;
  }

  ReactDOM.createRoot(document.getElementById("root")).render(html`<${App} />`);
})();
