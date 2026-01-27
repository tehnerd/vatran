(()=>{(()=>{let{useEffect:_,useMemo:S,useRef:T,useState:u,useContext:G}=React,{BrowserRouter:U,Routes:H,Route:C,NavLink:x,Link:q,useParams:M,useNavigate:X}=ReactRouterDOM,r=htm.bind(React.createElement),F=React.createContext({addToast:()=>{}});function W(){return G(F)}let h={base:"/api/v1",async request(e,t={}){let c={method:t.method||"GET",headers:{"Content-Type":"application/json"}};t.body!==void 0&&(c.body=JSON.stringify(t.body));let l=await fetch(`${h.base}${e}`,c),i;try{i=await l.json()}catch{throw new Error("invalid JSON response")}if(!l.ok)throw new Error(i?.error?.message||`HTTP ${l.status}`);if(!i.success){let b=i.error?.message||"request failed";throw new Error(b)}return i.data},get(e,t){return h.request(e,{method:"GET",body:t})},post(e,t){return h.request(e,{method:"POST",body:t})},put(e,t){return h.request(e,{method:"PUT",body:t})},del(e,t){return h.request(e,{method:"DELETE",body:t})}};function V(e){return`${encodeURIComponent(e.address)}:${e.port}:${e.proto}`}function D(e){let t=e.split(":"),c=Number(t.pop()||0),l=Number(t.pop()||0);return{address:decodeURIComponent(t.join(":")),port:l,proto:c}}function j(e,t,c=[]){let[l,i]=u(null),[b,$]=u(""),[n,f]=u(!0);return _(()=>{let m=!0,N=async()=>{try{let p=await e();m&&(i(p),$(""),f(!1))}catch(p){m&&($(p.message||"request failed"),f(!1))}};N();let o=setInterval(N,t);return()=>{m=!1,clearInterval(o)}},c),{data:l,error:b,loading:n}}function E({path:e,body:t,intervalMs:c=1e3,limit:l=60}){let[i,b]=u([]),[$,n]=u(""),f=S(()=>JSON.stringify(t||{}),[t]);return _(()=>{if(t===null)return b([]),n(""),()=>{};let m=!0,N=async()=>{try{let p=await h.get(e,t);if(!m)return;let g=new Date().toLocaleTimeString();b(w=>w.concat({label:g,...p}).slice(-l)),n("")}catch(p){m&&n(p.message||"request failed")}};N();let o=setInterval(N,c);return()=>{m=!1,clearInterval(o)}},[e,f,c,l]),{points:i,error:$}}function L({title:e,points:t,keys:c}){let l=T(null),i=T(null);return _(()=>{if(!l.current)return;i.current||(i.current=new Chart(l.current,{type:"line",data:{labels:[],datasets:[]},options:{responsive:!0,animation:!1,scales:{x:{grid:{display:!1}},y:{beginAtZero:!0}},plugins:{legend:{display:!0,position:"bottom"},title:{display:!!e,text:e}}}}));let b=i.current,$=t.map(n=>n.label);return b.data.labels=$,b.data.datasets=c.map(n=>({label:n.label,data:t.map(f=>f[n.field]||0),borderColor:n.color,backgroundColor:n.fill,borderWidth:2,tension:.3})),b.update(),()=>{}},[t,c,e]),_(()=>()=>{i.current&&(i.current.destroy(),i.current=null)},[]),r`<canvas ref=${l} height="120"></canvas>`}function O({children:e}){return e}function J({toasts:e,onDismiss:t}){return r`
      <div className="toast-stack">
        ${e.map(c=>r`
            <div className=${`toast ${c.kind}`}>
              <span>${c.message}</span>
              <button className="toast-close" onClick=${()=>t(c.id)}>
                ×
              </button>
            </div>
          `)}
      </div>
    `}function K({status:e}){return r`
      <header>
        <div>
          <div style=${{fontSize:20,fontWeight:700}}>Vatran</div>
          <div className="status-pill">
            <span className=${`dot ${e.ready?"ok":""}`}></span>
            ${e.ready?"Ready":"Not ready"}
          </div>
        </div>
        <nav>
          <${x} to="/" end className=${({isActive:t})=>t?"active":""}>
            Dashboard
          </${x}>
          <${x} to="/stats/global" className=${({isActive:t})=>t?"active":""}>
            Global stats
          </${x}>
          <${x} to="/stats/real" className=${({isActive:t})=>t?"active":""}>
            Per-real stats
          </${x}>
        </nav>
      </header>
    `}function Z(){let{addToast:e}=W(),[t,c]=u({initialized:!1,ready:!1}),[l,i]=u([]),[b,$]=u(""),[n,f]=u(!1),[m,N]=u(!1),[o,p]=u({main_interface:"",balancer_prog_path:"",healthchecking_prog_path:"",default_mac:"",local_mac:"",use_root_map:!1,max_vips:1024,max_reals:4096,hash_func:0}),[g,w]=u({address:"",port:80,proto:6,flags:0}),d=async()=>{try{let a=await h.get("/lb/status"),I=await h.get("/vips");c(a||{initialized:!1,ready:!1}),i(I||[]),$("")}catch(a){$(a.message||"request failed")}};_(()=>{let a=!0;return(async()=>{a&&await d()})(),()=>{a=!1}},[]);let R=async a=>{a.preventDefault();try{await h.post("/lb/create",{...o,max_vips:Number(o.max_vips),max_reals:Number(o.max_reals),hash_func:Number(o.hash_func)}),$(""),f(!1),e("Load balancer initialized.","success"),await d()}catch(I){$(I.message||"request failed"),e(I.message||"Initialize failed.","error")}},y=async a=>{a.preventDefault();try{await h.post("/vips",{...g,port:Number(g.port),proto:Number(g.proto),flags:Number(g.flags||0)}),w({address:"",port:80,proto:6,flags:0}),$(""),N(!1),e("VIP created.","success"),await d()}catch(I){$(I.message||"request failed"),e(I.message||"VIP create failed.","error")}};return r`
      <main>
        <section className="card">
          <h2>Load balancer</h2>
          <p>Initialized: ${t.initialized?"yes":"no"}</p>
          <p>Ready: ${t.ready?"yes":"no"}</p>
          <div className="row">
            <button className="btn" onClick=${()=>f(a=>!a)}>
              ${n?"Close":"Initialize"}
            </button>
            <button className="btn secondary" onClick=${()=>N(a=>!a)}>
              ${m?"Close":"Create VIP"}
            </button>
          </div>
          ${n&&r`
            <form className="form" onSubmit=${R}>
              <div className="form-row">
                <label className="field">
                  <span>Main interface</span>
                  <input
                    value=${o.main_interface}
                    onInput=${a=>p({...o,main_interface:a.target.value})}
                    placeholder="eth0"
                    required
                  />
                </label>
                <label className="field">
                  <span>Balancer prog path</span>
                  <input
                    value=${o.balancer_prog_path}
                    onInput=${a=>p({...o,balancer_prog_path:a.target.value})}
                    placeholder="/path/to/bpf.o"
                    required
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Healthchecking path</span>
                  <input
                    value=${o.healthchecking_prog_path}
                    onInput=${a=>p({...o,healthchecking_prog_path:a.target.value})}
                    placeholder="/path/to/hc.o"
                  />
                </label>
                <label className="field">
                  <span>Default MAC</span>
                  <input
                    value=${o.default_mac}
                    onInput=${a=>p({...o,default_mac:a.target.value})}
                    placeholder="aa:bb:cc:dd:ee:ff"
                    required
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Local MAC</span>
                  <input
                    value=${o.local_mac}
                    onInput=${a=>p({...o,local_mac:a.target.value})}
                    placeholder="11:22:33:44:55:66"
                    required
                  />
                </label>
                <label className="field">
                  <span>Hash function</span>
                  <input
                    type="number"
                    min="0"
                    value=${o.hash_func}
                    onInput=${a=>p({...o,hash_func:a.target.value})}
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Max VIPs</span>
                  <input
                    type="number"
                    min="1"
                    value=${o.max_vips}
                    onInput=${a=>p({...o,max_vips:a.target.value})}
                  />
                </label>
                <label className="field">
                  <span>Max reals</span>
                  <input
                    type="number"
                    min="1"
                    value=${o.max_reals}
                    onInput=${a=>p({...o,max_reals:a.target.value})}
                  />
                </label>
              </div>
              <label className="field checkbox">
                <input
                  type="checkbox"
                  checked=${o.use_root_map}
                  onChange=${a=>p({...o,use_root_map:a.target.checked})}
                />
                <span>Use root map</span>
              </label>
              <button className="btn" type="submit">Initialize LB</button>
            </form>
          `}
          ${m&&r`
            <form className="form" onSubmit=${y}>
              <div className="form-row">
                <label className="field">
                  <span>VIP address</span>
                  <input
                    value=${g.address}
                    onInput=${a=>w({...g,address:a.target.value})}
                    placeholder="1.2.3.4"
                    required
                  />
                </label>
                <label className="field">
                  <span>Port</span>
                  <input
                    type="number"
                    value=${g.port}
                    onInput=${a=>w({...g,port:a.target.value})}
                    required
                  />
                </label>
                <label className="field">
                  <span>Protocol</span>
                  <select
                    value=${g.proto}
                    onChange=${a=>w({...g,proto:a.target.value})}
                  >
                    <option value="6">TCP (6)</option>
                    <option value="17">UDP (17)</option>
                  </select>
                </label>
                <label className="field">
                  <span>Flags</span>
                  <input
                    type="number"
                    value=${g.flags}
                    onInput=${a=>w({...g,flags:a.target.value})}
                  />
                </label>
              </div>
              <button className="btn" type="submit">Create VIP</button>
            </form>
          `}
          ${b&&r`<p className="error">${b}</p>`}
        </section>
        <section className="card">
          <div className="section-header">
            <h2>VIPs</h2>
            <button className="btn ghost" onClick=${d}>Refresh</button>
          </div>
          ${l.length===0?r`<p className="muted">No VIPs configured yet.</p>`:r`
                <div className="grid">
                  ${l.map(a=>r`
                      <div className="card">
                        <div style=${{fontWeight:600}}>
                          ${a.address}:${a.port} / ${a.proto}
                        </div>
                        <div className="muted" style=${{marginTop:6}}>
                          Flags: ${a.flags||0}
                        </div>
                        <div className="row" style=${{marginTop:12}}>
                          <${q} className="btn" to=${`/vips/${V(a)}`}>
                            Open
                          </${q}>
                          <${q}
                            className="btn secondary"
                            to=${`/vips/${V(a)}/stats`}
                          >
                            Stats
                          </${q}>
                        </div>
                      </div>
                    `)}
                </div>
              `}
        </section>
      </main>
    `}function Q(){let{addToast:e}=W(),t=M(),c=X(),l=S(()=>D(t.vipId),[t.vipId]),[i,b]=u([]),[$,n]=u(""),[f,m]=u(""),[N,o]=u(!0),[p,g]=u({address:"",weight:100,flags:0}),[w,d]=u({}),[R,y]=u(null),[a,I]=u({flag:0,set:!0}),[A,se]=u({hash_function:0}),k=async()=>{try{let s=await h.get("/vips/reals",l);b(s||[]);let v={};(s||[]).forEach(B=>{v[B.address]=B.weight}),d(v),n(""),o(!1)}catch(s){n(s.message||"request failed"),o(!1)}},z=async()=>{try{let s=await h.get("/vips/flags",l);y(s?.flags??0),m("")}catch(s){m(s.message||"request failed")}};_(()=>{k(),z()},[t.vipId]);let le=async s=>{try{let v=Number(w[s.address]);await h.post("/vips/reals",{vip:l,real:{address:s.address,weight:v,flags:s.flags||0}}),await k(),e("Real weight updated.","success")}catch(v){n(v.message||"request failed"),e(v.message||"Update failed.","error")}},re=async s=>{try{await h.del("/vips/reals",{vip:l,real:{address:s.address,weight:s.weight,flags:s.flags||0}}),await k(),e("Real removed.","success")}catch(v){n(v.message||"request failed"),e(v.message||"Remove failed.","error")}},ne=async s=>{s.preventDefault();try{await h.post("/vips/reals",{vip:l,real:{address:p.address,weight:Number(p.weight),flags:Number(p.flags||0)}}),g({address:"",weight:100,flags:0}),await k(),e("Real added.","success")}catch(v){n(v.message||"request failed"),e(v.message||"Add failed.","error")}},oe=async()=>{try{await h.del("/vips",l),e("VIP deleted.","success"),c("/")}catch(s){n(s.message||"request failed"),e(s.message||"Delete failed.","error")}},ie=async s=>{s.preventDefault();try{await h.put("/vips/flags",{...l,flag:Number(a.flag),set:!!a.set}),await z(),e("VIP flags updated.","success")}catch(v){m(v.message||"request failed"),e(v.message||"Flag update failed.","error")}},ce=async s=>{s.preventDefault();try{await h.put("/vips/hash-function",{...l,hash_function:Number(A.hash_function)}),e("Hash function updated.","success")}catch(v){m(v.message||"request failed"),e(v.message||"Hash update failed.","error")}};return r`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>VIP Detail</h2>
              <p className="muted">${l.address}:${l.port} / ${l.proto}</p>
              <p className="muted">Flags: ${R??"\u2014"}</p>
            </div>
            <div className="row">
              <button className="btn ghost" onClick=${z}>Refresh flags</button>
              <button className="btn danger" onClick=${oe}>Delete VIP</button>
            </div>
          </div>
          ${$&&r`<p className="error">${$}</p>`}
          ${f&&r`<p className="error">${f}</p>`}
          ${N?r`<p className="muted">Loading reals…</p>`:r`
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
                    ${i.map(s=>r`
                        <tr>
                          <td>
                            <span
                              className=${`dot ${Number(s.weight)>0?"ok":"bad"}`}
                            ></span>
                          </td>
                          <td>${s.address}</td>
                          <td>
                            <input
                              className="inline-input"
                              type="number"
                              min="0"
                              value=${w[s.address]??s.weight}
                              onInput=${v=>d({...w,[s.address]:v.target.value})}
                            />
                          </td>
                          <td>${s.flags||0}</td>
                          <td className="row">
                            <button className="btn" onClick=${()=>le(s)}>
                              Update
                            </button>
                            <button className="btn ghost" onClick=${()=>re(s)}>
                              Remove
                            </button>
                          </td>
                        </tr>
                      `)}
                  </tbody>
                </table>
              `}
        </section>
        <section className="card">
          <h3>VIP flags & hash</h3>
          <div className="grid">
            <form className="form" onSubmit=${ie}>
              <div className="form-row">
                <label className="field">
                  <span>Flag</span>
                  <input
                    type="number"
                    value=${a.flag}
                    onInput=${s=>I({...a,flag:s.target.value})}
                  />
                </label>
                <label className="field">
                  <span>Set</span>
                  <select
                    value=${String(a.set)}
                    onChange=${s=>I({...a,set:s.target.value==="true"})}
                  >
                    <option value="true">Set</option>
                    <option value="false">Clear</option>
                  </select>
                </label>
              </div>
              <button className="btn" type="submit">Apply flag</button>
            </form>
            <form className="form" onSubmit=${ce}>
              <div className="form-row">
                <label className="field">
                  <span>Hash function</span>
                  <input
                    type="number"
                    min="0"
                    value=${A.hash_function}
                    onInput=${s=>se({...A,hash_function:s.target.value})}
                  />
                </label>
              </div>
              <button className="btn" type="submit">Apply hash</button>
            </form>
          </div>
        </section>
        <section className="card">
          <h3>Add real</h3>
          <form className="form" onSubmit=${ne}>
            <div className="form-row">
              <label className="field">
                <span>Real address</span>
                <input
                  value=${p.address}
                  onInput=${s=>g({...p,address:s.target.value})}
                  placeholder="10.0.0.1"
                  required
                />
              </label>
              <label className="field">
                <span>Weight</span>
                <input
                  type="number"
                  min="0"
                  value=${p.weight}
                  onInput=${s=>g({...p,weight:s.target.value})}
                  required
                />
              </label>
              <label className="field">
                <span>Flags</span>
                <input
                  type="number"
                  value=${p.flags}
                  onInput=${s=>g({...p,flags:s.target.value})}
                />
              </label>
            </div>
            <button className="btn" type="submit">Add real</button>
          </form>
        </section>
      </main>
    `}function Y(){let e=M(),t=S(()=>D(e.vipId),[e.vipId]),{points:c,error:l}=E({path:"/stats/vip",body:t}),i=S(()=>[{label:"v1",field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.15)"},{label:"v2",field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[]);return r`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>VIP Stats</h2>
              <p className="muted">${t.address}:${t.port} / ${t.proto}</p>
            </div>
          </div>
          ${l&&r`<p className="error">${l}</p>`}
          <${L} title="Traffic" points=${c} keys=${i} />
        </section>
      </main>
    `}function P({title:e,path:t}){let{points:c,error:l}=E({path:t}),i=S(()=>[{label:"v1",field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.2)"},{label:"v2",field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[]);return r`
      <div className="card">
        <h3>${e}</h3>
        ${l&&r`<p className="error">${l}</p>`}
        <${L} points=${c} keys=${i} />
      </div>
    `}function ee(){let{data:e,error:t}=j(()=>h.get("/stats/userspace"),1e3,[]);return r`
      <main>
        <section className="card">
          <h2>Global Stats</h2>
          <p className="muted">Polling every second.</p>
        </section>
        <section className="grid">
          <${P} title="LRU" path="/stats/lru" />
          <${P} title="LRU Miss" path="/stats/lru/miss" />
          <${P} title="LRU Fallback" path="/stats/lru/fallback" />
          <${P} title="LRU Global" path="/stats/lru/global" />
          <${P} title="XDP Total" path="/stats/xdp/total" />
          <${P} title="XDP Pass" path="/stats/xdp/pass" />
          <${P} title="XDP Drop" path="/stats/xdp/drop" />
          <${P} title="XDP Tx" path="/stats/xdp/tx" />
        </section>
        <section className="card">
          <h3>Userspace</h3>
          ${t&&r`<p className="error">${t}</p>`}
          ${e?r`
                <div className="row">
                  <div className="stat">
                    <span className="muted">BPF failed calls</span>
                    <strong>${e.bpf_failed_calls??0}</strong>
                  </div>
                  <div className="stat">
                    <span className="muted">Addr validation failed</span>
                    <strong>${e.addr_validation_failed??0}</strong>
                  </div>
                </div>
              `:r`<p className="muted">Waiting for data…</p>`}
        </section>
      </main>
    `}function ae(){let[e,t]=u([]),[c,l]=u(""),[i,b]=u([]),[$,n]=u(""),[f,m]=u(null),[N,o]=u("");_(()=>{let d=!0;return(async()=>{try{let y=await h.get("/vips");if(!d)return;t(y||[]),!c&&y&&y.length>0&&l(V(y[0]))}catch(y){d&&o(y.message||"request failed")}})(),()=>{d=!1}},[]),_(()=>{if(!c)return;let d=D(c),R=!0;return(async()=>{try{let a=await h.get("/vips/reals",d);if(!R)return;b(a||[]),a&&a.length>0?n(I=>I||a[0].address):n(""),o("")}catch(a){R&&o(a.message||"request failed")}})(),()=>{R=!1}},[c]),_(()=>{if(!$){m(null);return}let d=!0;return(async()=>{try{let y=await h.get("/reals/index",{address:$});if(!d)return;m(y?.index??null),o("")}catch(y){d&&o(y.message||"request failed")}})(),()=>{d=!1}},[$]);let{points:p,error:g}=E({path:"/stats/real",body:f!==null?{index:f}:null}),w=S(()=>[{label:"v1",field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.2)"},{label:"v2",field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[]);return r`
      <main>
        <section className="card">
          <h2>Per-Real Stats</h2>
          <p className="muted">Select a VIP and real address to chart.</p>
          ${N&&r`<p className="error">${N}</p>`}
          <div className="form-row">
            <label className="field">
              <span>VIP</span>
              <select value=${c} onChange=${d=>l(d.target.value)}>
                ${e.map(d=>r`
                    <option value=${V(d)}>
                      ${d.address}:${d.port} / ${d.proto}
                    </option>
                  `)}
              </select>
            </label>
            <label className="field">
              <span>Real</span>
              <select
                value=${$}
                onChange=${d=>n(d.target.value)}
                disabled=${i.length===0}
              >
                ${i.map(d=>r`
                    <option value=${d.address}>${d.address}</option>
                  `)}
              </select>
            </label>
            <label className="field">
              <span>Index</span>
              <input value=${f??""} readOnly />
            </label>
          </div>
        </section>
        <section className="card">
          <h3>Real stats</h3>
          ${g&&r`<p className="error">${g}</p>`}
          ${f===null?r`<p className="muted">Select a real to start polling.</p>`:r`<${L} points=${p} keys=${w} />`}
        </section>
      </main>
    `}function te(){let[e,t]=u({initialized:!1,ready:!1}),[c,l]=u([]),i=T({}),b=(n,f="info")=>{let m=`${Date.now()}-${Math.random().toString(16).slice(2)}`;l(N=>N.concat({id:m,message:n,kind:f})),i.current[m]=setTimeout(()=>{l(N=>N.filter(o=>o.id!==m)),delete i.current[m]},4e3)},$=n=>{i.current[n]&&(clearTimeout(i.current[n]),delete i.current[n]),l(f=>f.filter(m=>m.id!==n))};return _(()=>{let n=!0,f=async()=>{try{let N=await h.get("/lb/status");n&&t(N||{initialized:!1,ready:!1})}catch{n&&t({initialized:!1,ready:!1})}};f();let m=setInterval(f,5e3);return()=>{n=!1,clearInterval(m)}},[]),r`
      <${U}>
        <${O}>
          <${F.Provider} value=${{addToast:b}}>
            <${K} status=${e} />
            <${H}>
              <${C} path="/" element=${r`<${Z} />`} />
              <${C} path="/vips/:vipId" element=${r`<${Q} />`} />
              <${C} path="/vips/:vipId/stats" element=${r`<${Y} />`} />
              <${C} path="/stats/global" element=${r`<${ee} />`} />
              <${C} path="/stats/real" element=${r`<${ae} />`} />
            </${H}>
            <${J} toasts=${c} onDismiss=${$} />
          </${F.Provider}>
        </${O}>
      </${U}>
    `}ReactDOM.createRoot(document.getElementById("root")).render(r`<${te} />`)})();})();
