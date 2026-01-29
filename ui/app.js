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

  function useToast() {
    return useContext(ToastContext);
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

  function StatChart({ title, points, keys }) {
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
            animation: false,
            scales: {
              x: { grid: { display: false } },
              y: { beginAtZero: true },
            },
            plugins: {
              legend: { display: true, position: "bottom" },
              title: { display: Boolean(title), text: title },
            },
          },
        });
      }
      const chart = chartRef.current;
      const labels = points.map((p) => p.label);
      chart.data.labels = labels;
      chart.data.datasets = keys.map((key) => ({
        label: key.label,
        data: points.map((p) => p[key.field] || 0),
        borderColor: key.color,
        backgroundColor: key.fill,
        borderWidth: 2,
        tension: 0.3,
      }));
      chart.update();
      return () => {};
    }, [points, keys, title]);

    useEffect(() => {
      return () => {
        if (chartRef.current) {
          chartRef.current.destroy();
          chartRef.current = null;
        }
      };
    }, []);

    return html`<canvas ref=${canvasRef} height="120"></canvas>`;
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
      max_vips: 1024,
      max_reals: 4096,
      hash_func: 0,
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
        setStatus(lbStatus || { initialized: false, ready: false });
        setVips(list || []);
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
        const payload = {
          ...initForm,
          root_map_pos: initForm.root_map_pos === "" ? undefined : Number(initForm.root_map_pos),
          max_vips: Number(initForm.max_vips),
          max_reals: Number(initForm.max_reals),
          hash_func: Number(initForm.hash_func),
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

    return html`
      <main>
        <section className="card">
          <h2>Load balancer</h2>
          <p>Initialized: ${status.initialized ? "yes" : "no"}</p>
          <p>Ready: ${status.ready ? "yes" : "no"}</p>
          <div className="row">
            <button className="btn" onClick=${() => setShowInit((s) => !s)}>
              ${showInit ? "Close" : "Initialize"}
            </button>
            <button className="btn secondary" onClick=${() => setShowVipForm((s) => !s)}>
              ${showVipForm ? "Close" : "Create VIP"}
            </button>
          </div>
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
                    required
                  />
                </label>
                <label className="field">
                  <span>Hash function</span>
                  <input
                    type="number"
                    min="0"
                    value=${initForm.hash_func}
                    onInput=${(e) =>
                      setInitForm({ ...initForm, hash_func: e.target.value })}
                  />
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
                  <input
                    type="number"
                    value=${vipForm.flags}
                    onInput=${(e) => setVipForm({ ...vipForm, flags: e.target.value })}
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
                        <div style=${{ fontWeight: 600 }}>
                          ${vip.address}:${vip.port} / ${vip.proto}
                        </div>
                        <div className="muted" style=${{ marginTop: 6 }}>
                          Flags: ${vip.flags || 0}
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
    const [flagForm, setFlagForm] = useState({ flag: 0, set: true });
    const [hashForm, setHashForm] = useState({ hash_function: 0 });

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
          flag: Number(flagForm.flag),
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
              <p className="muted">Flags: ${vipFlags ?? "—"}</p>
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
                      <th>Flags</th>
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
                          <td>${real.flags || 0}</td>
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
                  <span>Flag</span>
                  <input
                    type="number"
                    value=${flagForm.flag}
                    onInput=${(e) => setFlagForm({ ...flagForm, flag: e.target.value })}
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
              <label className="field">
                <span>Flags</span>
                <input
                  type="number"
                  value=${newReal.flags}
                  onInput=${(e) => setNewReal({ ...newReal, flags: e.target.value })}
                />
              </label>
            </div>
            <button className="btn" type="submit">Add real</button>
          </form>
        </section>
      </main>
    `;
  }

  function VipStats() {
    const params = useParams();
    const vip = useMemo(() => parseVipId(params.vipId), [params.vipId]);
    const { points, error } = useStatSeries({ path: "/stats/vip", body: vip });
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
          <${StatChart} title="Traffic" points=${points} keys=${keys} />
        </section>
      </main>
    `;
  }

  function StatPanel({ title, path }) {
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
        <${StatChart} points=${points} keys=${keys} />
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
          <p className="muted">Polling every second.</p>
        </section>
        <section className="grid">
          <${StatPanel} title="LRU" path="/stats/lru" />
          <${StatPanel} title="LRU Miss" path="/stats/lru/miss" />
          <${StatPanel} title="LRU Fallback" path="/stats/lru/fallback" />
          <${StatPanel} title="LRU Global" path="/stats/lru/global" />
          <${StatPanel} title="XDP Total" path="/stats/xdp/total" />
          <${StatPanel} title="XDP Pass" path="/stats/xdp/pass" />
          <${StatPanel} title="XDP Drop" path="/stats/xdp/drop" />
          <${StatPanel} title="XDP Tx" path="/stats/xdp/tx" />
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

    const keys = useMemo(
      () => [
        { label: "v1", field: "v1", color: "#2f4858", fill: "rgba(47,72,88,0.2)" },
        { label: "v2", field: "v2", color: "#d97757", fill: "rgba(217,119,87,0.2)" },
      ],
      []
    );

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
            : html`<${StatChart} points=${points} keys=${keys} />`}
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
