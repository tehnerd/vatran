(()=>{(()=>{let{useEffect:P,useMemo:C,useRef:q,useState:f,useContext:j}=React,{BrowserRouter:H,Routes:M,Route:k,NavLink:S,Link:T,useParams:O,useNavigate:J}=ReactRouterDOM,r=htm.bind(React.createElement),D=React.createContext({addToast:()=>{}});function L(){return j(D)}let h={base:"/api/v1",async request(e,a={}){let u={method:a.method||"GET",headers:{"Content-Type":"application/json"}},n=`${h.base}${e}`;if(a.body!==void 0&&a.body!==null)if(u.method==="GET"){let d=new URLSearchParams;Object.entries(a.body).forEach(([g,i])=>{if(i!=null){if(Array.isArray(i)){i.forEach(b=>d.append(g,String(b)));return}if(typeof i=="object"){d.set(g,JSON.stringify(i));return}d.set(g,String(i))}});let o=d.toString();o&&(n+=`${n.includes("?")?"&":"?"}${o}`)}else u.body=JSON.stringify(a.body);let p=await fetch(n,u),$;try{$=await p.json()}catch{throw new Error("invalid JSON response")}if(!p.ok)throw new Error($?.error?.message||`HTTP ${p.status}`);if(!$.success){let d=$.error?.message||"request failed";throw new Error(d)}return $.data},get(e,a){return h.request(e,{method:"GET",body:a})},post(e,a){return h.request(e,{method:"POST",body:a})},put(e,a){return h.request(e,{method:"PUT",body:a})},del(e,a){return h.request(e,{method:"DELETE",body:a})}};function V(e){return`${encodeURIComponent(e.address)}:${e.port}:${e.proto}`}function A(e){let a=e.split(":"),u=Number(a.pop()||0),n=Number(a.pop()||0);return{address:decodeURIComponent(a.join(":")),port:n,proto:u}}function X(e,a,u=[]){let[n,p]=f(null),[$,d]=f(""),[o,g]=f(!0);return P(()=>{let i=!0,b=async()=>{try{let c=await e();i&&(p(c),d(""),g(!1))}catch(c){i&&(d(c.message||"request failed"),g(!1))}};b();let s=setInterval(b,a);return()=>{i=!1,clearInterval(s)}},u),{data:n,error:$,loading:o}}function B({path:e,body:a,intervalMs:u=1e3,limit:n=60}){let[p,$]=f([]),[d,o]=f(""),g=C(()=>JSON.stringify(a||{}),[a]);return P(()=>{if(a===null)return $([]),o(""),()=>{};let i=!0,b=async()=>{try{let c=await h.get(e,a);if(!i)return;let v=new Date().toLocaleTimeString();$(y=>y.concat({label:v,...c}).slice(-n)),o("")}catch(c){i&&o(c.message||"request failed")}};b();let s=setInterval(b,u);return()=>{i=!1,clearInterval(s)}},[e,g,u,n]),{points:p,error:d}}function z({title:e,points:a,keys:u}){let n=q(null),p=q(null);return P(()=>{if(!n.current)return;p.current||(p.current=new Chart(n.current,{type:"line",data:{labels:[],datasets:[]},options:{responsive:!0,animation:!1,scales:{x:{grid:{display:!1}},y:{beginAtZero:!0}},plugins:{legend:{display:!0,position:"bottom"},title:{display:!!e,text:e}}}}));let $=p.current,d=a.map(o=>o.label);return $.data.labels=d,$.data.datasets=u.map(o=>({label:o.label,data:a.map(g=>g[o.field]||0),borderColor:o.color,backgroundColor:o.fill,borderWidth:2,tension:.3})),$.update(),()=>{}},[a,u,e]),P(()=>()=>{p.current&&(p.current.destroy(),p.current=null)},[]),r`<canvas ref=${n} height="120"></canvas>`}function W({children:e}){return e}function K({toasts:e,onDismiss:a}){return r`
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
          <${S} to="/" end className=${({isActive:a})=>a?"active":""}>
            Dashboard
          </${S}>
          <${S} to="/stats/global" className=${({isActive:a})=>a?"active":""}>
            Global stats
          </${S}>
          <${S} to="/stats/real" className=${({isActive:a})=>a?"active":""}>
            Per-real stats
          </${S}>
          <${S} to="/config" className=${({isActive:a})=>a?"active":""}>
            Config export
          </${S}>
        </nav>
      </header>
    `}function Z(){let{addToast:e}=L(),[a,u]=f({initialized:!1,ready:!1}),[n,p]=f([]),[$,d]=f(""),[o,g]=f(!1),[i,b]=f(!1),[s,c]=f({main_interface:"",balancer_prog_path:"",healthchecking_prog_path:"",default_mac:"",local_mac:"",root_map_path:"",root_map_pos:2,katran_src_v4:"",katran_src_v6:"",use_root_map:!1,max_vips:1024,max_reals:4096,hash_func:0}),[v,y]=f({address:"",port:80,proto:6,flags:0}),m=async()=>{try{let t=await h.get("/lb/status"),I=await h.get("/vips");u(t||{initialized:!1,ready:!1}),p(I||[]),d("")}catch(t){d(t.message||"request failed")}};P(()=>{let t=!0;return(async()=>{t&&await m()})(),()=>{t=!1}},[]);let R=async t=>{t.preventDefault();try{let I={...s,root_map_pos:s.root_map_pos===""?void 0:Number(s.root_map_pos),max_vips:Number(s.max_vips),max_reals:Number(s.max_reals),hash_func:Number(s.hash_func)};await h.post("/lb/create",I),d(""),g(!1),e("Load balancer initialized.","success"),await m()}catch(I){d(I.message||"request failed"),e(I.message||"Initialize failed.","error")}},w=async t=>{t.preventDefault();try{await h.post("/vips",{...v,port:Number(v.port),proto:Number(v.proto),flags:Number(v.flags||0)}),y({address:"",port:80,proto:6,flags:0}),d(""),b(!1),e("VIP created.","success"),await m()}catch(I){d(I.message||"request failed"),e(I.message||"VIP create failed.","error")}},_=async()=>{try{await h.post("/lb/load-bpf-progs"),d(""),e("BPF programs loaded.","success"),await m()}catch(t){d(t.message||"request failed"),e(t.message||"Load BPF programs failed.","error")}},F=async()=>{try{await h.post("/lb/attach-bpf-progs"),d(""),e("BPF programs attached.","success"),await m()}catch(t){d(t.message||"request failed"),e(t.message||"Attach BPF programs failed.","error")}};return r`
      <main>
        <section className="card">
          <h2>Load balancer</h2>
          <p>Initialized: ${a.initialized?"yes":"no"}</p>
          <p>Ready: ${a.ready?"yes":"no"}</p>
          <div className="row">
            <button className="btn" onClick=${()=>g(t=>!t)}>
              ${o?"Close":"Initialize"}
            </button>
            <button className="btn secondary" onClick=${()=>b(t=>!t)}>
              ${i?"Close":"Create VIP"}
            </button>
          </div>
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
          ${o&&r`
            <form className="form" onSubmit=${R}>
              <div className="form-row">
                <label className="field">
                  <span>Main interface</span>
                  <input
                    value=${s.main_interface}
                    onInput=${t=>c({...s,main_interface:t.target.value})}
                    placeholder="eth0"
                    required
                  />
                </label>
                <label className="field">
                  <span>Balancer prog path</span>
                  <input
                    value=${s.balancer_prog_path}
                    onInput=${t=>c({...s,balancer_prog_path:t.target.value})}
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
                    onInput=${t=>c({...s,healthchecking_prog_path:t.target.value})}
                    placeholder="/path/to/hc.o"
                  />
                </label>
                <label className="field">
                  <span>Default MAC</span>
                  <input
                    value=${s.default_mac}
                    onInput=${t=>c({...s,default_mac:t.target.value})}
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
                    onInput=${t=>c({...s,local_mac:t.target.value})}
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
                    onInput=${t=>c({...s,hash_func:t.target.value})}
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Root map path</span>
                  <input
                    value=${s.root_map_path}
                    onInput=${t=>c({...s,root_map_path:t.target.value})}
                    placeholder="/sys/fs/bpf/root_map"
                  />
                </label>
                <label className="field">
                  <span>Root map position</span>
                  <input
                    type="number"
                    min="0"
                    value=${s.root_map_pos}
                    onInput=${t=>c({...s,root_map_pos:t.target.value})}
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Katran src v4</span>
                  <input
                    value=${s.katran_src_v4}
                    onInput=${t=>c({...s,katran_src_v4:t.target.value})}
                    placeholder="10.0.0.1"
                  />
                </label>
                <label className="field">
                  <span>Katran src v6</span>
                  <input
                    value=${s.katran_src_v6}
                    onInput=${t=>c({...s,katran_src_v6:t.target.value})}
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
                    onInput=${t=>c({...s,max_vips:t.target.value})}
                  />
                </label>
                <label className="field">
                  <span>Max reals</span>
                  <input
                    type="number"
                    min="1"
                    value=${s.max_reals}
                    onInput=${t=>c({...s,max_reals:t.target.value})}
                  />
                </label>
              </div>
              <label className="field checkbox">
                <input
                  type="checkbox"
                  checked=${s.use_root_map}
                  onChange=${t=>c({...s,use_root_map:t.target.checked})}
                />
                <span>Use root map</span>
              </label>
              <button className="btn" type="submit">Initialize LB</button>
            </form>
          `}
          ${i&&r`
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
          ${$&&r`<p className="error">${$}</p>`}
        </section>
        <section className="card">
          <div className="section-header">
            <h2>VIPs</h2>
            <button className="btn ghost" onClick=${m}>Refresh</button>
          </div>
          ${n.length===0?r`<p className="muted">No VIPs configured yet.</p>`:r`
                <div className="grid">
                  ${n.map(t=>r`
                      <div className="card">
                        <div style=${{fontWeight:600}}>
                          ${t.address}:${t.port} / ${t.proto}
                        </div>
                        <div className="muted" style=${{marginTop:6}}>
                          Flags: ${t.flags||0}
                        </div>
                        <div className="row" style=${{marginTop:12}}>
                          <${T} className="btn" to=${`/vips/${V(t)}`}>
                            Open
                          </${T}>
                          <${T}
                            className="btn secondary"
                            to=${`/vips/${V(t)}/stats`}
                          >
                            Stats
                          </${T}>
                        </div>
                      </div>
                    `)}
                </div>
              `}
        </section>
      </main>
    `}function Q(){let{addToast:e}=L(),a=O(),u=J(),n=C(()=>A(a.vipId),[a.vipId]),[p,$]=f([]),[d,o]=f(""),[g,i]=f(""),[b,s]=f(!0),[c,v]=f({address:"",weight:100,flags:0}),[y,m]=f({}),[R,w]=f(null),[_,F]=f({flag:0,set:!0}),[t,I]=f({hash_function:0}),E=async()=>{try{let l=await h.get("/vips/reals",n);$(l||[]);let N={};(l||[]).forEach(G=>{N[G.address]=G.weight}),m(N),o(""),s(!1)}catch(l){o(l.message||"request failed"),s(!1)}},U=async()=>{try{let l=await h.get("/vips/flags",n);w(l?.flags??0),i("")}catch(l){i(l.message||"request failed")}};P(()=>{E(),U()},[a.vipId]);let le=async l=>{try{let N=Number(y[l.address]);await h.post("/vips/reals",{vip:n,real:{address:l.address,weight:N,flags:l.flags||0}}),await E(),e("Real weight updated.","success")}catch(N){o(N.message||"request failed"),e(N.message||"Update failed.","error")}},ne=async l=>{try{await h.del("/vips/reals",{vip:n,real:{address:l.address,weight:l.weight,flags:l.flags||0}}),await E(),e("Real removed.","success")}catch(N){o(N.message||"request failed"),e(N.message||"Remove failed.","error")}},oe=async l=>{l.preventDefault();try{await h.post("/vips/reals",{vip:n,real:{address:c.address,weight:Number(c.weight),flags:Number(c.flags||0)}}),v({address:"",weight:100,flags:0}),await E(),e("Real added.","success")}catch(N){o(N.message||"request failed"),e(N.message||"Add failed.","error")}},ie=async()=>{try{await h.del("/vips",n),e("VIP deleted.","success"),u("/")}catch(l){o(l.message||"request failed"),e(l.message||"Delete failed.","error")}},ce=async l=>{l.preventDefault();try{await h.put("/vips/flags",{...n,flag:Number(_.flag),set:!!_.set}),await U(),e("VIP flags updated.","success")}catch(N){i(N.message||"request failed"),e(N.message||"Flag update failed.","error")}},de=async l=>{l.preventDefault();try{await h.put("/vips/hash-function",{...n,hash_function:Number(t.hash_function)}),e("Hash function updated.","success")}catch(N){i(N.message||"request failed"),e(N.message||"Hash update failed.","error")}};return r`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>VIP Detail</h2>
              <p className="muted">${n.address}:${n.port} / ${n.proto}</p>
              <p className="muted">Flags: ${R??"\u2014"}</p>
            </div>
            <div className="row">
              <button className="btn ghost" onClick=${U}>Refresh flags</button>
              <button className="btn danger" onClick=${ie}>Delete VIP</button>
            </div>
          </div>
          ${d&&r`<p className="error">${d}</p>`}
          ${g&&r`<p className="error">${g}</p>`}
          ${b?r`<p className="muted">Loading reals…</p>`:r`
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
                    ${p.map(l=>r`
                        <tr>
                          <td>
                            <span
                              className=${`dot ${Number(l.weight)>0?"ok":"bad"}`}
                            ></span>
                          </td>
                          <td>${l.address}</td>
                          <td>
                            <input
                              className="inline-input"
                              type="number"
                              min="0"
                              value=${y[l.address]??l.weight}
                              onInput=${N=>m({...y,[l.address]:N.target.value})}
                            />
                          </td>
                          <td>${l.flags||0}</td>
                          <td className="row">
                            <button className="btn" onClick=${()=>le(l)}>
                              Update
                            </button>
                            <button className="btn ghost" onClick=${()=>ne(l)}>
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
            <form className="form" onSubmit=${ce}>
              <div className="form-row">
                <label className="field">
                  <span>Flag</span>
                  <input
                    type="number"
                    value=${_.flag}
                    onInput=${l=>F({..._,flag:l.target.value})}
                  />
                </label>
                <label className="field">
                  <span>Set</span>
                  <select
                    value=${String(_.set)}
                    onChange=${l=>F({..._,set:l.target.value==="true"})}
                  >
                    <option value="true">Set</option>
                    <option value="false">Clear</option>
                  </select>
                </label>
              </div>
              <button className="btn" type="submit">Apply flag</button>
            </form>
            <form className="form" onSubmit=${de}>
              <div className="form-row">
                <label className="field">
                  <span>Hash function</span>
                  <input
                    type="number"
                    min="0"
                    value=${t.hash_function}
                    onInput=${l=>I({...t,hash_function:l.target.value})}
                  />
                </label>
              </div>
              <button className="btn" type="submit">Apply hash</button>
            </form>
          </div>
        </section>
        <section className="card">
          <h3>Add real</h3>
          <form className="form" onSubmit=${oe}>
            <div className="form-row">
              <label className="field">
                <span>Real address</span>
                <input
                  value=${c.address}
                  onInput=${l=>v({...c,address:l.target.value})}
                  placeholder="10.0.0.1"
                  required
                />
              </label>
              <label className="field">
                <span>Weight</span>
                <input
                  type="number"
                  min="0"
                  value=${c.weight}
                  onInput=${l=>v({...c,weight:l.target.value})}
                  required
                />
              </label>
              <label className="field">
                <span>Flags</span>
                <input
                  type="number"
                  value=${c.flags}
                  onInput=${l=>v({...c,flags:l.target.value})}
                />
              </label>
            </div>
            <button className="btn" type="submit">Add real</button>
          </form>
        </section>
      </main>
    `}function ee(){let e=O(),a=C(()=>A(e.vipId),[e.vipId]),{points:u,error:n}=B({path:"/stats/vip",body:a}),p=C(()=>[{label:"v1",field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.15)"},{label:"v2",field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[]);return r`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>VIP Stats</h2>
              <p className="muted">${a.address}:${a.port} / ${a.proto}</p>
            </div>
          </div>
          ${n&&r`<p className="error">${n}</p>`}
          <${z} title="Traffic" points=${u} keys=${p} />
        </section>
      </main>
    `}function x({title:e,path:a}){let{points:u,error:n}=B({path:a}),p=C(()=>[{label:"v1",field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.2)"},{label:"v2",field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[]);return r`
      <div className="card">
        <h3>${e}</h3>
        ${n&&r`<p className="error">${n}</p>`}
        <${z} points=${u} keys=${p} />
      </div>
    `}function ae(){let{data:e,error:a}=X(()=>h.get("/stats/userspace"),1e3,[]);return r`
      <main>
        <section className="card">
          <h2>Global Stats</h2>
          <p className="muted">Polling every second.</p>
        </section>
        <section className="grid">
          <${x} title="LRU" path="/stats/lru" />
          <${x} title="LRU Miss" path="/stats/lru/miss" />
          <${x} title="LRU Fallback" path="/stats/lru/fallback" />
          <${x} title="LRU Global" path="/stats/lru/global" />
          <${x} title="XDP Total" path="/stats/xdp/total" />
          <${x} title="XDP Pass" path="/stats/xdp/pass" />
          <${x} title="XDP Drop" path="/stats/xdp/drop" />
          <${x} title="XDP Tx" path="/stats/xdp/tx" />
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
    `}function te(){let[e,a]=f([]),[u,n]=f(""),[p,$]=f([]),[d,o]=f(""),[g,i]=f(null),[b,s]=f("");P(()=>{let m=!0;return(async()=>{try{let w=await h.get("/vips");if(!m)return;a(w||[]),!u&&w&&w.length>0&&n(V(w[0]))}catch(w){m&&s(w.message||"request failed")}})(),()=>{m=!1}},[]),P(()=>{if(!u)return;let m=A(u),R=!0;return(async()=>{try{let _=await h.get("/vips/reals",m);if(!R)return;$(_||[]),_&&_.length>0?o(F=>F||_[0].address):o(""),s("")}catch(_){R&&s(_.message||"request failed")}})(),()=>{R=!1}},[u]),P(()=>{if(!d){i(null);return}let m=!0;return(async()=>{try{let w=await h.get("/reals/index",{address:d});if(!m)return;i(w?.index??null),s("")}catch(w){m&&s(w.message||"request failed")}})(),()=>{m=!1}},[d]);let{points:c,error:v}=B({path:"/stats/real",body:g!==null?{index:g}:null}),y=C(()=>[{label:"v1",field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.2)"},{label:"v2",field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[]);return r`
      <main>
        <section className="card">
          <h2>Per-Real Stats</h2>
          <p className="muted">Select a VIP and real address to chart.</p>
          ${b&&r`<p className="error">${b}</p>`}
          <div className="form-row">
            <label className="field">
              <span>VIP</span>
              <select value=${u} onChange=${m=>n(m.target.value)}>
                ${e.map(m=>r`
                    <option value=${V(m)}>
                      ${m.address}:${m.port} / ${m.proto}
                    </option>
                  `)}
              </select>
            </label>
            <label className="field">
              <span>Real</span>
              <select
                value=${d}
                onChange=${m=>o(m.target.value)}
                disabled=${p.length===0}
              >
                ${p.map(m=>r`
                    <option value=${m.address}>${m.address}</option>
                  `)}
              </select>
            </label>
            <label className="field">
              <span>Index</span>
              <input value=${g??""} readOnly />
            </label>
          </div>
        </section>
        <section className="card">
          <h3>Real stats</h3>
          ${v&&r`<p className="error">${v}</p>`}
          ${g===null?r`<p className="muted">Select a real to start polling.</p>`:r`<${z} points=${c} keys=${y} />`}
        </section>
      </main>
    `}function se(){let{addToast:e}=L(),[a,u]=f(""),[n,p]=f(""),[$,d]=f(!0),[o,g]=f(""),i=q(!0),b=async()=>{if(i.current){d(!0),p("");try{let c=await fetch(`${h.base}/config/export`,{headers:{Accept:"application/x-yaml"}});if(!c.ok){let y=`HTTP ${c.status}`;try{y=(await c.json())?.error?.message||y}catch{}throw new Error(y)}let v=await c.text();if(!i.current)return;u(v||""),g(new Date().toLocaleString())}catch(c){i.current&&p(c.message||"request failed")}finally{i.current&&d(!1)}}},s=async()=>{if(a)try{await navigator.clipboard.writeText(a),e("Config copied to clipboard","info")}catch{e("Failed to copy config","error")}};return P(()=>(i.current=!0,b(),()=>{i.current=!1}),[]),r`
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
              <button className="btn" onClick=${b} disabled=${$}>
                Refresh
              </button>
            </div>
          </div>
          ${n&&r`<p className="error">${n}</p>`}
          ${$?r`<p className="muted">Loading config...</p>`:a?r`<pre className="yaml-view">${a}</pre>`:r`<p className="muted">No config data returned.</p>`}
          ${o&&r`<p className="muted">Last fetched ${o}</p>`}
        </section>
      </main>
    `}function re(){let[e,a]=f({initialized:!1,ready:!1}),[u,n]=f([]),p=q({}),$=(o,g="info")=>{let i=`${Date.now()}-${Math.random().toString(16).slice(2)}`;n(b=>b.concat({id:i,message:o,kind:g})),p.current[i]=setTimeout(()=>{n(b=>b.filter(s=>s.id!==i)),delete p.current[i]},4e3)},d=o=>{p.current[o]&&(clearTimeout(p.current[o]),delete p.current[o]),n(g=>g.filter(i=>i.id!==o))};return P(()=>{let o=!0,g=async()=>{try{let b=await h.get("/lb/status");o&&a(b||{initialized:!1,ready:!1})}catch{o&&a({initialized:!1,ready:!1})}};g();let i=setInterval(g,5e3);return()=>{o=!1,clearInterval(i)}},[]),r`
      <${H}>
        <${W}>
          <${D.Provider} value=${{addToast:$}}>
            <${Y} status=${e} />
            <${M}>
              <${k} path="/" element=${r`<${Z} />`} />
              <${k} path="/vips/:vipId" element=${r`<${Q} />`} />
              <${k} path="/vips/:vipId/stats" element=${r`<${ee} />`} />
              <${k} path="/stats/global" element=${r`<${ae} />`} />
              <${k} path="/stats/real" element=${r`<${te} />`} />
              <${k} path="/config" element=${r`<${se} />`} />
            </${M}>
            <${K} toasts=${u} onDismiss=${d} />
          </${D.Provider}>
        </${W}>
      </${H}>
    `}ReactDOM.createRoot(document.getElementById("root")).render(r`<${re} />`)})();})();
