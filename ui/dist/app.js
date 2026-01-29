(()=>{(()=>{let{useEffect:P,useMemo:S,useRef:T,useState:m,useContext:J}=React,{BrowserRouter:H,Routes:M,Route:x,NavLink:R,Link:q,useParams:O,useNavigate:X}=ReactRouterDOM,r=htm.bind(React.createElement),D=React.createContext({addToast:()=>{}});function L(){return J(D)}let b={base:"/api/v1",async request(e,a={}){let u={method:a.method||"GET",headers:{"Content-Type":"application/json"}},l=`${b.base}${e}`;if(a.body!==void 0&&a.body!==null)if(u.method==="GET"){let i=new URLSearchParams;Object.entries(a.body).forEach(([p,c])=>{if(c!=null){if(Array.isArray(c)){c.forEach(h=>i.append(p,String(h)));return}if(typeof c=="object"){i.set(p,JSON.stringify(c));return}i.set(p,String(c))}});let o=i.toString();o&&(l+=`${l.includes("?")?"&":"?"}${o}`)}else u.body=JSON.stringify(a.body);let $=await fetch(l,u),g;try{g=await $.json()}catch{throw new Error("invalid JSON response")}if(!$.ok)throw new Error(g?.error?.message||`HTTP ${$.status}`);if(!g.success){let i=g.error?.message||"request failed";throw new Error(i)}return g.data},get(e,a){return b.request(e,{method:"GET",body:a})},post(e,a){return b.request(e,{method:"POST",body:a})},put(e,a){return b.request(e,{method:"PUT",body:a})},del(e,a){return b.request(e,{method:"DELETE",body:a})}};function V(e){return`${encodeURIComponent(e.address)}:${e.port}:${e.proto}`}function A(e){let a=e.split(":"),u=Number(a.pop()||0),l=Number(a.pop()||0);return{address:decodeURIComponent(a.join(":")),port:l,proto:u}}function K(e,a,u=[]){let[l,$]=m(null),[g,i]=m(""),[o,p]=m(!0);return P(()=>{let c=!0,h=async()=>{try{let d=await e();c&&($(d),i(""),p(!1))}catch(d){c&&(i(d.message||"request failed"),p(!1))}};h();let s=setInterval(h,a);return()=>{c=!1,clearInterval(s)}},u),{data:l,error:g,loading:o}}function z({path:e,body:a,intervalMs:u=1e3,limit:l=60}){let[$,g]=m([]),[i,o]=m(""),p=S(()=>JSON.stringify(a||{}),[a]);return P(()=>{if(a===null)return g([]),o(""),()=>{};let c=!0,h=async()=>{try{let d=await b.get(e,a);if(!c)return;let v=new Date().toLocaleTimeString();g(y=>y.concat({label:v,...d}).slice(-l)),o("")}catch(d){c&&o(d.message||"request failed")}};h();let s=setInterval(h,u);return()=>{c=!1,clearInterval(s)}},[e,p,u,l]),{points:$,error:i}}function W({title:e,points:a,keys:u,diff:l=!1,height:$=120,showTitle:g=!1}){let i=T(null),o=T(null);return P(()=>{if(!i.current)return;o.current||(o.current=new Chart(i.current,{type:"line",data:{labels:[],datasets:[]},options:{responsive:!0,maintainAspectRatio:!1,animation:!1,scales:{x:{grid:{display:!1}},y:{beginAtZero:!l}},plugins:{legend:{display:!0,position:"bottom"},title:{display:g&&!!e,text:e}}}}));let p=o.current,c=a.map(h=>h.label);return p.data.labels=c,p.data.datasets=u.map(h=>{let s=a.map(v=>v[h.field]||0),d=l?s.map((v,y)=>y===0?0:v-s[y-1]):s;return{label:h.label,data:d,borderColor:h.color,backgroundColor:h.fill,borderWidth:2,tension:.3}}),p.options.scales.y.beginAtZero=!l,p.options.plugins.title.display=g&&!!e,p.options.plugins.title.text=e||"",p.update(),()=>{}},[a,u,e,l,g]),P(()=>()=>{o.current&&(o.current.destroy(),o.current=null)},[]),r`<canvas ref=${i} height=${$}></canvas>`}function B({title:e,points:a,keys:u,diff:l=!1,inlineTitle:$=!0}){let[g,i]=m(!1);return r`
      <div className="chart-wrap">
        <div className="chart-click" onClick=${()=>i(!0)}>
          <${W}
            title=${e}
            points=${a}
            keys=${u}
            diff=${l}
            height=${120}
            showTitle=${$&&!!e}
          />
        </div>
        ${g&&r`
          <div className="chart-overlay" onClick=${()=>i(!1)}>
            <div className="chart-modal" onClick=${o=>o.stopPropagation()}>
              <div className="row chart-modal-header">
                <div>
                  <h3>${e||"Chart"}</h3>
                  ${l?r`<p className="muted">Per-second delta.</p>`:""}
                </div>
                <button className="btn ghost" onClick=${()=>i(!1)}>
                  Close
                </button>
              </div>
              <div className="chart-zoom">
                <${W}
                  title=${e}
                  points=${a}
                  keys=${u}
                  diff=${l}
                  height=${360}
                  showTitle=${!1}
                />
              </div>
            </div>
          </div>
        `}
      </div>
    `}function G({children:e}){return e}function Z({toasts:e,onDismiss:a}){return r`
      <div className="toast-stack">
        ${e.map(u=>r`
            <div className=${`toast ${u.kind}`}>
              <span>${u.message}</span>
              <button className="toast-close" onClick=${()=>a(u.id)}>
                ×
              </button>
            </div>
          `)}
      </div>
    `}function Y({status:e}){return r`
      <header>
        <div>
          <div style=${{fontSize:20,fontWeight:700}}>Vatran</div>
          <div className="status-pill">
            <span className=${`dot ${e.ready?"ok":""}`}></span>
            ${e.ready?"Ready":"Not ready"}
          </div>
        </div>
        <nav>
          <${R} to="/" end className=${({isActive:a})=>a?"active":""}>
            Dashboard
          </${R}>
          <${R} to="/stats/global" className=${({isActive:a})=>a?"active":""}>
            Global stats
          </${R}>
          <${R} to="/stats/real" className=${({isActive:a})=>a?"active":""}>
            Per-real stats
          </${R}>
          <${R} to="/config" className=${({isActive:a})=>a?"active":""}>
            Config export
          </${R}>
        </nav>
      </header>
    `}function Q(){let{addToast:e}=L(),[a,u]=m({initialized:!1,ready:!1}),[l,$]=m([]),[g,i]=m(""),[o,p]=m(!1),[c,h]=m(!1),[s,d]=m({main_interface:"",balancer_prog_path:"",healthchecking_prog_path:"",default_mac:"",local_mac:"",root_map_path:"",root_map_pos:2,katran_src_v4:"",katran_src_v6:"",use_root_map:!1,max_vips:512,max_reals:4096,hash_func:0}),[v,y]=m({address:"",port:80,proto:6,flags:0}),f=async()=>{try{let t=await b.get("/lb/status"),I=await b.get("/vips");u(t||{initialized:!1,ready:!1}),$(I||[]),i("")}catch(t){i(t.message||"request failed")}};P(()=>{let t=!0;return(async()=>{t&&await f()})(),()=>{t=!1}},[]);let C=async t=>{t.preventDefault();try{let I={...s,root_map_pos:s.root_map_pos===""?void 0:Number(s.root_map_pos),max_vips:Number(s.max_vips),max_reals:Number(s.max_reals),hash_func:Number(s.hash_func)};await b.post("/lb/create",I),i(""),p(!1),e("Load balancer initialized.","success"),await f()}catch(I){i(I.message||"request failed"),e(I.message||"Initialize failed.","error")}},w=async t=>{t.preventDefault();try{await b.post("/vips",{...v,port:Number(v.port),proto:Number(v.proto),flags:Number(v.flags||0)}),y({address:"",port:80,proto:6,flags:0}),i(""),h(!1),e("VIP created.","success"),await f()}catch(I){i(I.message||"request failed"),e(I.message||"VIP create failed.","error")}},_=async()=>{try{await b.post("/lb/load-bpf-progs"),i(""),e("BPF programs loaded.","success"),await f()}catch(t){i(t.message||"request failed"),e(t.message||"Load BPF programs failed.","error")}},F=async()=>{try{await b.post("/lb/attach-bpf-progs"),i(""),e("BPF programs attached.","success"),await f()}catch(t){i(t.message||"request failed"),e(t.message||"Attach BPF programs failed.","error")}};return r`
      <main>
        <section className="card">
          <h2>Load balancer</h2>
          <p>Initialized: ${a.initialized?"yes":"no"}</p>
          <p>Ready: ${a.ready?"yes":"no"}</p>
          <div className="row">
            ${!a.initialized&&r`
              <button className="btn" onClick=${()=>p(t=>!t)}>
                ${o?"Close":"Initialize"}
              </button>
            `}
            <button className="btn secondary" onClick=${()=>h(t=>!t)}>
              ${c?"Close":"Create VIP"}
            </button>
          </div>
          ${!a.ready&&r`
            <div className="row" style=${{marginTop:12}}>
              <button
                className="btn ghost"
                disabled=${!a.initialized}
                onClick=${_}
              >
                Load BPF Programs
              </button>
              <button
                className="btn ghost"
                disabled=${!a.initialized}
                onClick=${F}
              >
                Attach BPF Programs
              </button>
            </div>
          `}
          ${o&&r`
            <form className="form" onSubmit=${C}>
              <div className="form-row">
                <label className="field">
                  <span>Main interface</span>
                  <input
                    value=${s.main_interface}
                    onInput=${t=>d({...s,main_interface:t.target.value})}
                    placeholder="eth0"
                    required
                  />
                </label>
                <label className="field">
                  <span>Balancer prog path</span>
                  <input
                    value=${s.balancer_prog_path}
                    onInput=${t=>d({...s,balancer_prog_path:t.target.value})}
                    placeholder="/path/to/bpf.o"
                    required
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Healthchecking path</span>
                  <input
                    value=${s.healthchecking_prog_path}
                    onInput=${t=>d({...s,healthchecking_prog_path:t.target.value})}
                    placeholder="/path/to/hc.o"
                  />
                </label>
                <label className="field">
                  <span>Default MAC</span>
                  <input
                    value=${s.default_mac}
                    onInput=${t=>d({...s,default_mac:t.target.value})}
                    placeholder="aa:bb:cc:dd:ee:ff"
                    required
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Local MAC</span>
                  <input
                    value=${s.local_mac}
                    onInput=${t=>d({...s,local_mac:t.target.value})}
                    placeholder="11:22:33:44:55:66"
                    required=${s.healthchecking_prog_path?.trim()!==""}
                  />
                </label>
                <label className="field">
                  <span>Hash function</span>
                  <input
                    type="number"
                    min="0"
                    value=${s.hash_func}
                    onInput=${t=>d({...s,hash_func:t.target.value})}
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Root map path</span>
                  <input
                    value=${s.root_map_path}
                    onInput=${t=>d({...s,root_map_path:t.target.value})}
                    placeholder="/sys/fs/bpf/root_map"
                  />
                </label>
                <label className="field">
                  <span>Root map position</span>
                  <input
                    type="number"
                    min="0"
                    value=${s.root_map_pos}
                    onInput=${t=>d({...s,root_map_pos:t.target.value})}
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Katran src v4</span>
                  <input
                    value=${s.katran_src_v4}
                    onInput=${t=>d({...s,katran_src_v4:t.target.value})}
                    placeholder="10.0.0.1"
                  />
                </label>
                <label className="field">
                  <span>Katran src v6</span>
                  <input
                    value=${s.katran_src_v6}
                    onInput=${t=>d({...s,katran_src_v6:t.target.value})}
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
                    value=${s.max_vips}
                    onInput=${t=>d({...s,max_vips:t.target.value})}
                  />
                </label>
                <label className="field">
                  <span>Max reals</span>
                  <input
                    type="number"
                    min="1"
                    value=${s.max_reals}
                    onInput=${t=>d({...s,max_reals:t.target.value})}
                  />
                </label>
              </div>
              <label className="field checkbox">
                <input
                  type="checkbox"
                  checked=${s.use_root_map}
                  onChange=${t=>d({...s,use_root_map:t.target.checked})}
                />
                <span>Use root map</span>
              </label>
              <button className="btn" type="submit">Initialize LB</button>
            </form>
          `}
          ${c&&r`
            <form className="form" onSubmit=${w}>
              <div className="form-row">
                <label className="field">
                  <span>VIP address</span>
                  <input
                    value=${v.address}
                    onInput=${t=>y({...v,address:t.target.value})}
                    placeholder="1.2.3.4"
                    required
                  />
                </label>
                <label className="field">
                  <span>Port</span>
                  <input
                    type="number"
                    value=${v.port}
                    onInput=${t=>y({...v,port:t.target.value})}
                    required
                  />
                </label>
                <label className="field">
                  <span>Protocol</span>
                  <select
                    value=${v.proto}
                    onChange=${t=>y({...v,proto:t.target.value})}
                  >
                    <option value="6">TCP (6)</option>
                    <option value="17">UDP (17)</option>
                  </select>
                </label>
                <label className="field">
                  <span>Flags</span>
                  <input
                    type="number"
                    value=${v.flags}
                    onInput=${t=>y({...v,flags:t.target.value})}
                  />
                </label>
              </div>
              <button className="btn" type="submit">Create VIP</button>
            </form>
          `}
          ${g&&r`<p className="error">${g}</p>`}
        </section>
        <section className="card">
          <div className="section-header">
            <h2>VIPs</h2>
            <button className="btn ghost" onClick=${f}>Refresh</button>
          </div>
          ${l.length===0?r`<p className="muted">No VIPs configured yet.</p>`:r`
                <div className="grid">
                  ${l.map(t=>r`
                      <div className="card">
                        <div style=${{fontWeight:600}}>
                          ${t.address}:${t.port} / ${t.proto}
                        </div>
                        <div className="muted" style=${{marginTop:6}}>
                          Flags: ${t.flags||0}
                        </div>
                        <div className="row" style=${{marginTop:12}}>
                          <${q} className="btn" to=${`/vips/${V(t)}`}>
                            Open
                          </${q}>
                          <${q}
                            className="btn secondary"
                            to=${`/vips/${V(t)}/stats`}
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
    `}function ee(){let{addToast:e}=L(),a=O(),u=X(),l=S(()=>A(a.vipId),[a.vipId]),[$,g]=m([]),[i,o]=m(""),[p,c]=m(""),[h,s]=m(!0),[d,v]=m({address:"",weight:100,flags:0}),[y,f]=m({}),[C,w]=m(null),[_,F]=m({flag:0,set:!0}),[t,I]=m({hash_function:0}),E=async()=>{try{let n=await b.get("/vips/reals",l);g(n||[]);let N={};(n||[]).forEach(j=>{N[j.address]=j.weight}),f(N),o(""),s(!1)}catch(n){o(n.message||"request failed"),s(!1)}},U=async()=>{try{let n=await b.get("/vips/flags",l);w(n?.flags??0),c("")}catch(n){c(n.message||"request failed")}};P(()=>{E(),U()},[a.vipId]);let ne=async n=>{try{let N=Number(y[n.address]);await b.post("/vips/reals",{vip:l,real:{address:n.address,weight:N,flags:n.flags||0}}),await E(),e("Real weight updated.","success")}catch(N){o(N.message||"request failed"),e(N.message||"Update failed.","error")}},oe=async n=>{try{await b.del("/vips/reals",{vip:l,real:{address:n.address,weight:n.weight,flags:n.flags||0}}),await E(),e("Real removed.","success")}catch(N){o(N.message||"request failed"),e(N.message||"Remove failed.","error")}},ie=async n=>{n.preventDefault();try{await b.post("/vips/reals",{vip:l,real:{address:d.address,weight:Number(d.weight),flags:Number(d.flags||0)}}),v({address:"",weight:100,flags:0}),await E(),e("Real added.","success")}catch(N){o(N.message||"request failed"),e(N.message||"Add failed.","error")}},ce=async()=>{try{await b.del("/vips",l),e("VIP deleted.","success"),u("/")}catch(n){o(n.message||"request failed"),e(n.message||"Delete failed.","error")}},de=async n=>{n.preventDefault();try{await b.put("/vips/flags",{...l,flag:Number(_.flag),set:!!_.set}),await U(),e("VIP flags updated.","success")}catch(N){c(N.message||"request failed"),e(N.message||"Flag update failed.","error")}},ue=async n=>{n.preventDefault();try{await b.put("/vips/hash-function",{...l,hash_function:Number(t.hash_function)}),e("Hash function updated.","success")}catch(N){c(N.message||"request failed"),e(N.message||"Hash update failed.","error")}};return r`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>VIP Detail</h2>
              <p className="muted">${l.address}:${l.port} / ${l.proto}</p>
              <p className="muted">Flags: ${C??"\u2014"}</p>
            </div>
            <div className="row">
              <button className="btn ghost" onClick=${U}>Refresh flags</button>
              <button className="btn danger" onClick=${ce}>Delete VIP</button>
            </div>
          </div>
          ${i&&r`<p className="error">${i}</p>`}
          ${p&&r`<p className="error">${p}</p>`}
          ${h?r`<p className="muted">Loading reals…</p>`:r`
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
                    ${$.map(n=>r`
                        <tr>
                          <td>
                            <span
                              className=${`dot ${Number(n.weight)>0?"ok":"bad"}`}
                            ></span>
                          </td>
                          <td>${n.address}</td>
                          <td>
                            <input
                              className="inline-input"
                              type="number"
                              min="0"
                              value=${y[n.address]??n.weight}
                              onInput=${N=>f({...y,[n.address]:N.target.value})}
                            />
                          </td>
                          <td>${n.flags||0}</td>
                          <td className="row">
                            <button className="btn" onClick=${()=>ne(n)}>
                              Update
                            </button>
                            <button className="btn ghost" onClick=${()=>oe(n)}>
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
            <form className="form" onSubmit=${de}>
              <div className="form-row">
                <label className="field">
                  <span>Flag</span>
                  <input
                    type="number"
                    value=${_.flag}
                    onInput=${n=>F({..._,flag:n.target.value})}
                  />
                </label>
                <label className="field">
                  <span>Set</span>
                  <select
                    value=${String(_.set)}
                    onChange=${n=>F({..._,set:n.target.value==="true"})}
                  >
                    <option value="true">Set</option>
                    <option value="false">Clear</option>
                  </select>
                </label>
              </div>
              <button className="btn" type="submit">Apply flag</button>
            </form>
            <form className="form" onSubmit=${ue}>
              <div className="form-row">
                <label className="field">
                  <span>Hash function</span>
                  <input
                    type="number"
                    min="0"
                    value=${t.hash_function}
                    onInput=${n=>I({...t,hash_function:n.target.value})}
                  />
                </label>
              </div>
              <button className="btn" type="submit">Apply hash</button>
            </form>
          </div>
        </section>
        <section className="card">
          <h3>Add real</h3>
          <form className="form" onSubmit=${ie}>
            <div className="form-row">
              <label className="field">
                <span>Real address</span>
                <input
                  value=${d.address}
                  onInput=${n=>v({...d,address:n.target.value})}
                  placeholder="10.0.0.1"
                  required
                />
              </label>
              <label className="field">
                <span>Weight</span>
                <input
                  type="number"
                  min="0"
                  value=${d.weight}
                  onInput=${n=>v({...d,weight:n.target.value})}
                  required
                />
              </label>
              <label className="field">
                <span>Flags</span>
                <input
                  type="number"
                  value=${d.flags}
                  onInput=${n=>v({...d,flags:n.target.value})}
                />
              </label>
            </div>
            <button className="btn" type="submit">Add real</button>
          </form>
        </section>
      </main>
    `}function ae(){let e=O(),a=S(()=>A(e.vipId),[e.vipId]),{points:u,error:l}=z({path:"/stats/vip",body:a}),$=S(()=>[{label:"v1",field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.15)"},{label:"v2",field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[]);return r`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>VIP Stats</h2>
              <p className="muted">${a.address}:${a.port} / ${a.proto}</p>
            </div>
          </div>
          ${l&&r`<p className="error">${l}</p>`}
          <${B} title="Traffic" points=${u} keys=${$} />
        </section>
      </main>
    `}function k({title:e,path:a,diff:u=!1}){let{points:l,error:$}=z({path:a}),g=S(()=>[{label:"v1",field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.2)"},{label:"v2",field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[]);return r`
      <div className="card">
        <h3>${e}</h3>
        ${$&&r`<p className="error">${$}</p>`}
        <${B} title=${e} points=${l} keys=${g} diff=${u} inlineTitle=${!1} />
      </div>
    `}function te(){let{data:e,error:a}=K(()=>b.get("/stats/userspace"),1e3,[]);return r`
      <main>
        <section className="card">
          <h2>Global Stats</h2>
          <p className="muted">Polling every second. Charts show per-second deltas.</p>
        </section>
        <section className="grid">
          <${k} title="LRU" path="/stats/lru" diff=${!0} />
          <${k} title="LRU Miss" path="/stats/lru/miss" diff=${!0} />
          <${k} title="LRU Fallback" path="/stats/lru/fallback" diff=${!0} />
          <${k} title="LRU Global" path="/stats/lru/global" diff=${!0} />
          <${k} title="XDP Total" path="/stats/xdp/total" diff=${!0} />
          <${k} title="XDP Pass" path="/stats/xdp/pass" diff=${!0} />
          <${k} title="XDP Drop" path="/stats/xdp/drop" diff=${!0} />
          <${k} title="XDP Tx" path="/stats/xdp/tx" diff=${!0} />
        </section>
        <section className="card">
          <h3>Userspace</h3>
          ${a&&r`<p className="error">${a}</p>`}
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
    `}function se(){let[e,a]=m([]),[u,l]=m(""),[$,g]=m([]),[i,o]=m(""),[p,c]=m(null),[h,s]=m("");P(()=>{let f=!0;return(async()=>{try{let w=await b.get("/vips");if(!f)return;a(w||[]),!u&&w&&w.length>0&&l(V(w[0]))}catch(w){f&&s(w.message||"request failed")}})(),()=>{f=!1}},[]),P(()=>{if(!u)return;let f=A(u),C=!0;return(async()=>{try{let _=await b.get("/vips/reals",f);if(!C)return;g(_||[]),_&&_.length>0?o(F=>F||_[0].address):o(""),s("")}catch(_){C&&s(_.message||"request failed")}})(),()=>{C=!1}},[u]),P(()=>{if(!i){c(null);return}let f=!0;return(async()=>{try{let w=await b.get("/reals/index",{address:i});if(!f)return;c(w?.index??null),s("")}catch(w){f&&s(w.message||"request failed")}})(),()=>{f=!1}},[i]);let{points:d,error:v}=z({path:"/stats/real",body:p!==null?{index:p}:null}),y=S(()=>[{label:"v1",field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.2)"},{label:"v2",field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[]);return r`
      <main>
        <section className="card">
          <h2>Per-Real Stats</h2>
          <p className="muted">Select a VIP and real address to chart.</p>
          ${h&&r`<p className="error">${h}</p>`}
          <div className="form-row">
            <label className="field">
              <span>VIP</span>
              <select value=${u} onChange=${f=>l(f.target.value)}>
                ${e.map(f=>r`
                    <option value=${V(f)}>
                      ${f.address}:${f.port} / ${f.proto}
                    </option>
                  `)}
              </select>
            </label>
            <label className="field">
              <span>Real</span>
              <select
                value=${i}
                onChange=${f=>o(f.target.value)}
                disabled=${$.length===0}
              >
                ${$.map(f=>r`
                    <option value=${f.address}>${f.address}</option>
                  `)}
              </select>
            </label>
            <label className="field">
              <span>Index</span>
              <input value=${p??""} readOnly />
            </label>
          </div>
        </section>
        <section className="card">
          <h3>Real stats</h3>
          ${v&&r`<p className="error">${v}</p>`}
          ${p===null?r`<p className="muted">Select a real to start polling.</p>`:r`<${B} points=${d} keys=${y} />`}
        </section>
      </main>
    `}function re(){let{addToast:e}=L(),[a,u]=m(""),[l,$]=m(""),[g,i]=m(!0),[o,p]=m(""),c=T(!0),h=async()=>{if(c.current){i(!0),$("");try{let d=await fetch(`${b.base}/config/export`,{headers:{Accept:"application/x-yaml"}});if(!d.ok){let y=`HTTP ${d.status}`;try{y=(await d.json())?.error?.message||y}catch{}throw new Error(y)}let v=await d.text();if(!c.current)return;u(v||""),p(new Date().toLocaleString())}catch(d){c.current&&$(d.message||"request failed")}finally{c.current&&i(!1)}}},s=async()=>{if(a)try{await navigator.clipboard.writeText(a),e("Config copied to clipboard","info")}catch{e("Failed to copy config","error")}};return P(()=>(c.current=!0,h(),()=>{c.current=!1}),[]),r`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>Running config</h2>
              <p className="muted">Exported from /api/v1/config/export</p>
            </div>
            <div className="row">
              <button className="btn ghost" onClick=${s} disabled=${!a}>
                Copy YAML
              </button>
              <button className="btn" onClick=${h} disabled=${g}>
                Refresh
              </button>
            </div>
          </div>
          ${l&&r`<p className="error">${l}</p>`}
          ${g?r`<p className="muted">Loading config...</p>`:a?r`<pre className="yaml-view">${a}</pre>`:r`<p className="muted">No config data returned.</p>`}
          ${o&&r`<p className="muted">Last fetched ${o}</p>`}
        </section>
      </main>
    `}function le(){let[e,a]=m({initialized:!1,ready:!1}),[u,l]=m([]),$=T({}),g=(o,p="info")=>{let c=`${Date.now()}-${Math.random().toString(16).slice(2)}`;l(h=>h.concat({id:c,message:o,kind:p})),$.current[c]=setTimeout(()=>{l(h=>h.filter(s=>s.id!==c)),delete $.current[c]},4e3)},i=o=>{$.current[o]&&(clearTimeout($.current[o]),delete $.current[o]),l(p=>p.filter(c=>c.id!==o))};return P(()=>{let o=!0,p=async()=>{try{let h=await b.get("/lb/status");o&&a(h||{initialized:!1,ready:!1})}catch{o&&a({initialized:!1,ready:!1})}};p();let c=setInterval(p,5e3);return()=>{o=!1,clearInterval(c)}},[]),r`
      <${H}>
        <${G}>
          <${D.Provider} value=${{addToast:g}}>
            <${Y} status=${e} />
            <${M}>
              <${x} path="/" element=${r`<${Q} />`} />
              <${x} path="/vips/:vipId" element=${r`<${ee} />`} />
              <${x} path="/vips/:vipId/stats" element=${r`<${ae} />`} />
              <${x} path="/stats/global" element=${r`<${te} />`} />
              <${x} path="/stats/real" element=${r`<${se} />`} />
              <${x} path="/config" element=${r`<${re} />`} />
            </${M}>
            <${Z} toasts=${u} onDismiss=${i} />
          </${D.Provider}>
        </${G}>
      </${H}>
    `}ReactDOM.createRoot(document.getElementById("root")).render(r`<${le} />`)})();})();
