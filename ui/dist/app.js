(()=>{(()=>{let{useEffect:P,useMemo:k,useRef:D,useState:p,useContext:X}=React,{BrowserRouter:U,Routes:H,Route:q,NavLink:C,Link:F,useParams:M,useNavigate:j}=ReactRouterDOM,n=htm.bind(React.createElement),E=React.createContext({addToast:()=>{}});function W(){return X(E)}let h={base:"/api/v1",async request(e,t={}){let d={method:t.method||"GET",headers:{"Content-Type":"application/json"}};t.body!==void 0&&(d.body=JSON.stringify(t.body));let r=await fetch(`${h.base}${e}`,d),i;try{i=await r.json()}catch{throw new Error("invalid JSON response")}if(!r.ok)throw new Error(i?.error?.message||`HTTP ${r.status}`);if(!i.success){let v=i.error?.message||"request failed";throw new Error(v)}return i.data},get(e,t){return h.request(e,{method:"GET",body:t})},post(e,t){return h.request(e,{method:"POST",body:t})},put(e,t){return h.request(e,{method:"PUT",body:t})},del(e,t){return h.request(e,{method:"DELETE",body:t})}};function V(e){return`${encodeURIComponent(e.address)}:${e.port}:${e.proto}`}function L(e){let t=e.split(":"),d=Number(t.pop()||0),r=Number(t.pop()||0);return{address:decodeURIComponent(t.join(":")),port:r,proto:d}}function J(e,t,d=[]){let[r,i]=p(null),[v,m]=p(""),[o,g]=p(!0);return P(()=>{let f=!0,N=async()=>{try{let c=await e();f&&(i(c),m(""),g(!1))}catch(c){f&&(m(c.message||"request failed"),g(!1))}};N();let s=setInterval(N,t);return()=>{f=!1,clearInterval(s)}},d),{data:r,error:v,loading:o}}function A({path:e,body:t,intervalMs:d=1e3,limit:r=60}){let[i,v]=p([]),[m,o]=p(""),g=k(()=>JSON.stringify(t||{}),[t]);return P(()=>{if(t===null)return v([]),o(""),()=>{};let f=!0,N=async()=>{try{let c=await h.get(e,t);if(!f)return;let $=new Date().toLocaleTimeString();v(w=>w.concat({label:$,...c}).slice(-r)),o("")}catch(c){f&&o(c.message||"request failed")}};N();let s=setInterval(N,d);return()=>{f=!1,clearInterval(s)}},[e,g,d,r]),{points:i,error:m}}function B({title:e,points:t,keys:d}){let r=D(null),i=D(null);return P(()=>{if(!r.current)return;i.current||(i.current=new Chart(r.current,{type:"line",data:{labels:[],datasets:[]},options:{responsive:!0,animation:!1,scales:{x:{grid:{display:!1}},y:{beginAtZero:!0}},plugins:{legend:{display:!0,position:"bottom"},title:{display:!!e,text:e}}}}));let v=i.current,m=t.map(o=>o.label);return v.data.labels=m,v.data.datasets=d.map(o=>({label:o.label,data:t.map(g=>g[o.field]||0),borderColor:o.color,backgroundColor:o.fill,borderWidth:2,tension:.3})),v.update(),()=>{}},[t,d,e]),P(()=>()=>{i.current&&(i.current.destroy(),i.current=null)},[]),n`<canvas ref=${r} height="120"></canvas>`}function O({children:e}){return e}function K({toasts:e,onDismiss:t}){return n`
      <div className="toast-stack">
        ${e.map(d=>n`
            <div className=${`toast ${d.kind}`}>
              <span>${d.message}</span>
              <button className="toast-close" onClick=${()=>t(d.id)}>
                ×
              </button>
            </div>
          `)}
      </div>
    `}function Z({status:e}){return n`
      <header>
        <div>
          <div style=${{fontSize:20,fontWeight:700}}>Vatran</div>
          <div className="status-pill">
            <span className=${`dot ${e.ready?"ok":""}`}></span>
            ${e.ready?"Ready":"Not ready"}
          </div>
        </div>
        <nav>
          <${C} to="/" end className=${({isActive:t})=>t?"active":""}>
            Dashboard
          </${C}>
          <${C} to="/stats/global" className=${({isActive:t})=>t?"active":""}>
            Global stats
          </${C}>
          <${C} to="/stats/real" className=${({isActive:t})=>t?"active":""}>
            Per-real stats
          </${C}>
        </nav>
      </header>
    `}function Q(){let{addToast:e}=W(),[t,d]=p({initialized:!1,ready:!1}),[r,i]=p([]),[v,m]=p(""),[o,g]=p(!1),[f,N]=p(!1),[s,c]=p({main_interface:"",balancer_prog_path:"",healthchecking_prog_path:"",default_mac:"",local_mac:"",root_map_path:"",root_map_pos:2,katran_src_v4:"",katran_src_v6:"",use_root_map:!1,max_vips:1024,max_reals:4096,hash_func:0}),[$,w]=p({address:"",port:80,proto:6,flags:0}),u=async()=>{try{let a=await h.get("/lb/status"),I=await h.get("/vips");d(a||{initialized:!1,ready:!1}),i(I||[]),m("")}catch(a){m(a.message||"request failed")}};P(()=>{let a=!0;return(async()=>{a&&await u()})(),()=>{a=!1}},[]);let R=async a=>{a.preventDefault();try{let I={...s,root_map_pos:s.root_map_pos===""?void 0:Number(s.root_map_pos),max_vips:Number(s.max_vips),max_reals:Number(s.max_reals),hash_func:Number(s.hash_func)};await h.post("/lb/create",I),m(""),g(!1),e("Load balancer initialized.","success"),await u()}catch(I){m(I.message||"request failed"),e(I.message||"Initialize failed.","error")}},y=async a=>{a.preventDefault();try{await h.post("/vips",{...$,port:Number($.port),proto:Number($.proto),flags:Number($.flags||0)}),w({address:"",port:80,proto:6,flags:0}),m(""),N(!1),e("VIP created.","success"),await u()}catch(I){m(I.message||"request failed"),e(I.message||"VIP create failed.","error")}},_=async()=>{try{await h.post("/lb/load-bpf-progs"),m(""),e("BPF programs loaded.","success"),await u()}catch(a){m(a.message||"request failed"),e(a.message||"Load BPF programs failed.","error")}},x=async()=>{try{await h.post("/lb/attach-bpf-progs"),m(""),e("BPF programs attached.","success"),await u()}catch(a){m(a.message||"request failed"),e(a.message||"Attach BPF programs failed.","error")}};return n`
      <main>
        <section className="card">
          <h2>Load balancer</h2>
          <p>Initialized: ${t.initialized?"yes":"no"}</p>
          <p>Ready: ${t.ready?"yes":"no"}</p>
          <div className="row">
            <button className="btn" onClick=${()=>g(a=>!a)}>
              ${o?"Close":"Initialize"}
            </button>
            <button className="btn secondary" onClick=${()=>N(a=>!a)}>
              ${f?"Close":"Create VIP"}
            </button>
          </div>
          <div className="row" style=${{marginTop:12}}>
            <button
              className="btn ghost"
              disabled=${!t.initialized}
              onClick=${_}
            >
              Load BPF Programs
            </button>
            <button
              className="btn ghost"
              disabled=${!t.initialized}
              onClick=${x}
            >
              Attach BPF Programs
            </button>
          </div>
          ${o&&n`
            <form className="form" onSubmit=${R}>
              <div className="form-row">
                <label className="field">
                  <span>Main interface</span>
                  <input
                    value=${s.main_interface}
                    onInput=${a=>c({...s,main_interface:a.target.value})}
                    placeholder="eth0"
                    required
                  />
                </label>
                <label className="field">
                  <span>Balancer prog path</span>
                  <input
                    value=${s.balancer_prog_path}
                    onInput=${a=>c({...s,balancer_prog_path:a.target.value})}
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
                    onInput=${a=>c({...s,healthchecking_prog_path:a.target.value})}
                    placeholder="/path/to/hc.o"
                  />
                </label>
                <label className="field">
                  <span>Default MAC</span>
                  <input
                    value=${s.default_mac}
                    onInput=${a=>c({...s,default_mac:a.target.value})}
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
                    onInput=${a=>c({...s,local_mac:a.target.value})}
                    placeholder="11:22:33:44:55:66"
                    required
                  />
                </label>
                <label className="field">
                  <span>Hash function</span>
                  <input
                    type="number"
                    min="0"
                    value=${s.hash_func}
                    onInput=${a=>c({...s,hash_func:a.target.value})}
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Root map path</span>
                  <input
                    value=${s.root_map_path}
                    onInput=${a=>c({...s,root_map_path:a.target.value})}
                    placeholder="/sys/fs/bpf/root_map"
                  />
                </label>
                <label className="field">
                  <span>Root map position</span>
                  <input
                    type="number"
                    min="0"
                    value=${s.root_map_pos}
                    onInput=${a=>c({...s,root_map_pos:a.target.value})}
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Katran src v4</span>
                  <input
                    value=${s.katran_src_v4}
                    onInput=${a=>c({...s,katran_src_v4:a.target.value})}
                    placeholder="10.0.0.1"
                  />
                </label>
                <label className="field">
                  <span>Katran src v6</span>
                  <input
                    value=${s.katran_src_v6}
                    onInput=${a=>c({...s,katran_src_v6:a.target.value})}
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
                    onInput=${a=>c({...s,max_vips:a.target.value})}
                  />
                </label>
                <label className="field">
                  <span>Max reals</span>
                  <input
                    type="number"
                    min="1"
                    value=${s.max_reals}
                    onInput=${a=>c({...s,max_reals:a.target.value})}
                  />
                </label>
              </div>
              <label className="field checkbox">
                <input
                  type="checkbox"
                  checked=${s.use_root_map}
                  onChange=${a=>c({...s,use_root_map:a.target.checked})}
                />
                <span>Use root map</span>
              </label>
              <button className="btn" type="submit">Initialize LB</button>
            </form>
          `}
          ${f&&n`
            <form className="form" onSubmit=${y}>
              <div className="form-row">
                <label className="field">
                  <span>VIP address</span>
                  <input
                    value=${$.address}
                    onInput=${a=>w({...$,address:a.target.value})}
                    placeholder="1.2.3.4"
                    required
                  />
                </label>
                <label className="field">
                  <span>Port</span>
                  <input
                    type="number"
                    value=${$.port}
                    onInput=${a=>w({...$,port:a.target.value})}
                    required
                  />
                </label>
                <label className="field">
                  <span>Protocol</span>
                  <select
                    value=${$.proto}
                    onChange=${a=>w({...$,proto:a.target.value})}
                  >
                    <option value="6">TCP (6)</option>
                    <option value="17">UDP (17)</option>
                  </select>
                </label>
                <label className="field">
                  <span>Flags</span>
                  <input
                    type="number"
                    value=${$.flags}
                    onInput=${a=>w({...$,flags:a.target.value})}
                  />
                </label>
              </div>
              <button className="btn" type="submit">Create VIP</button>
            </form>
          `}
          ${v&&n`<p className="error">${v}</p>`}
        </section>
        <section className="card">
          <div className="section-header">
            <h2>VIPs</h2>
            <button className="btn ghost" onClick=${u}>Refresh</button>
          </div>
          ${r.length===0?n`<p className="muted">No VIPs configured yet.</p>`:n`
                <div className="grid">
                  ${r.map(a=>n`
                      <div className="card">
                        <div style=${{fontWeight:600}}>
                          ${a.address}:${a.port} / ${a.proto}
                        </div>
                        <div className="muted" style=${{marginTop:6}}>
                          Flags: ${a.flags||0}
                        </div>
                        <div className="row" style=${{marginTop:12}}>
                          <${F} className="btn" to=${`/vips/${V(a)}`}>
                            Open
                          </${F}>
                          <${F}
                            className="btn secondary"
                            to=${`/vips/${V(a)}/stats`}
                          >
                            Stats
                          </${F}>
                        </div>
                      </div>
                    `)}
                </div>
              `}
        </section>
      </main>
    `}function Y(){let{addToast:e}=W(),t=M(),d=j(),r=k(()=>L(t.vipId),[t.vipId]),[i,v]=p([]),[m,o]=p(""),[g,f]=p(""),[N,s]=p(!0),[c,$]=p({address:"",weight:100,flags:0}),[w,u]=p({}),[R,y]=p(null),[_,x]=p({flag:0,set:!0}),[a,I]=p({hash_function:0}),T=async()=>{try{let l=await h.get("/vips/reals",r);v(l||[]);let b={};(l||[]).forEach(G=>{b[G.address]=G.weight}),u(b),o(""),s(!1)}catch(l){o(l.message||"request failed"),s(!1)}},z=async()=>{try{let l=await h.get("/vips/flags",r);y(l?.flags??0),f("")}catch(l){f(l.message||"request failed")}};P(()=>{T(),z()},[t.vipId]);let le=async l=>{try{let b=Number(w[l.address]);await h.post("/vips/reals",{vip:r,real:{address:l.address,weight:b,flags:l.flags||0}}),await T(),e("Real weight updated.","success")}catch(b){o(b.message||"request failed"),e(b.message||"Update failed.","error")}},re=async l=>{try{await h.del("/vips/reals",{vip:r,real:{address:l.address,weight:l.weight,flags:l.flags||0}}),await T(),e("Real removed.","success")}catch(b){o(b.message||"request failed"),e(b.message||"Remove failed.","error")}},ne=async l=>{l.preventDefault();try{await h.post("/vips/reals",{vip:r,real:{address:c.address,weight:Number(c.weight),flags:Number(c.flags||0)}}),$({address:"",weight:100,flags:0}),await T(),e("Real added.","success")}catch(b){o(b.message||"request failed"),e(b.message||"Add failed.","error")}},oe=async()=>{try{await h.del("/vips",r),e("VIP deleted.","success"),d("/")}catch(l){o(l.message||"request failed"),e(l.message||"Delete failed.","error")}},ie=async l=>{l.preventDefault();try{await h.put("/vips/flags",{...r,flag:Number(_.flag),set:!!_.set}),await z(),e("VIP flags updated.","success")}catch(b){f(b.message||"request failed"),e(b.message||"Flag update failed.","error")}},ce=async l=>{l.preventDefault();try{await h.put("/vips/hash-function",{...r,hash_function:Number(a.hash_function)}),e("Hash function updated.","success")}catch(b){f(b.message||"request failed"),e(b.message||"Hash update failed.","error")}};return n`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>VIP Detail</h2>
              <p className="muted">${r.address}:${r.port} / ${r.proto}</p>
              <p className="muted">Flags: ${R??"\u2014"}</p>
            </div>
            <div className="row">
              <button className="btn ghost" onClick=${z}>Refresh flags</button>
              <button className="btn danger" onClick=${oe}>Delete VIP</button>
            </div>
          </div>
          ${m&&n`<p className="error">${m}</p>`}
          ${g&&n`<p className="error">${g}</p>`}
          ${N?n`<p className="muted">Loading reals…</p>`:n`
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
                    ${i.map(l=>n`
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
                              value=${w[l.address]??l.weight}
                              onInput=${b=>u({...w,[l.address]:b.target.value})}
                            />
                          </td>
                          <td>${l.flags||0}</td>
                          <td className="row">
                            <button className="btn" onClick=${()=>le(l)}>
                              Update
                            </button>
                            <button className="btn ghost" onClick=${()=>re(l)}>
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
                    value=${_.flag}
                    onInput=${l=>x({..._,flag:l.target.value})}
                  />
                </label>
                <label className="field">
                  <span>Set</span>
                  <select
                    value=${String(_.set)}
                    onChange=${l=>x({..._,set:l.target.value==="true"})}
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
                    value=${a.hash_function}
                    onInput=${l=>I({...a,hash_function:l.target.value})}
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
                  value=${c.address}
                  onInput=${l=>$({...c,address:l.target.value})}
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
                  onInput=${l=>$({...c,weight:l.target.value})}
                  required
                />
              </label>
              <label className="field">
                <span>Flags</span>
                <input
                  type="number"
                  value=${c.flags}
                  onInput=${l=>$({...c,flags:l.target.value})}
                />
              </label>
            </div>
            <button className="btn" type="submit">Add real</button>
          </form>
        </section>
      </main>
    `}function ee(){let e=M(),t=k(()=>L(e.vipId),[e.vipId]),{points:d,error:r}=A({path:"/stats/vip",body:t}),i=k(()=>[{label:"v1",field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.15)"},{label:"v2",field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[]);return n`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>VIP Stats</h2>
              <p className="muted">${t.address}:${t.port} / ${t.proto}</p>
            </div>
          </div>
          ${r&&n`<p className="error">${r}</p>`}
          <${B} title="Traffic" points=${d} keys=${i} />
        </section>
      </main>
    `}function S({title:e,path:t}){let{points:d,error:r}=A({path:t}),i=k(()=>[{label:"v1",field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.2)"},{label:"v2",field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[]);return n`
      <div className="card">
        <h3>${e}</h3>
        ${r&&n`<p className="error">${r}</p>`}
        <${B} points=${d} keys=${i} />
      </div>
    `}function ae(){let{data:e,error:t}=J(()=>h.get("/stats/userspace"),1e3,[]);return n`
      <main>
        <section className="card">
          <h2>Global Stats</h2>
          <p className="muted">Polling every second.</p>
        </section>
        <section className="grid">
          <${S} title="LRU" path="/stats/lru" />
          <${S} title="LRU Miss" path="/stats/lru/miss" />
          <${S} title="LRU Fallback" path="/stats/lru/fallback" />
          <${S} title="LRU Global" path="/stats/lru/global" />
          <${S} title="XDP Total" path="/stats/xdp/total" />
          <${S} title="XDP Pass" path="/stats/xdp/pass" />
          <${S} title="XDP Drop" path="/stats/xdp/drop" />
          <${S} title="XDP Tx" path="/stats/xdp/tx" />
        </section>
        <section className="card">
          <h3>Userspace</h3>
          ${t&&n`<p className="error">${t}</p>`}
          ${e?n`
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
              `:n`<p className="muted">Waiting for data…</p>`}
        </section>
      </main>
    `}function te(){let[e,t]=p([]),[d,r]=p(""),[i,v]=p([]),[m,o]=p(""),[g,f]=p(null),[N,s]=p("");P(()=>{let u=!0;return(async()=>{try{let y=await h.get("/vips");if(!u)return;t(y||[]),!d&&y&&y.length>0&&r(V(y[0]))}catch(y){u&&s(y.message||"request failed")}})(),()=>{u=!1}},[]),P(()=>{if(!d)return;let u=L(d),R=!0;return(async()=>{try{let _=await h.get("/vips/reals",u);if(!R)return;v(_||[]),_&&_.length>0?o(x=>x||_[0].address):o(""),s("")}catch(_){R&&s(_.message||"request failed")}})(),()=>{R=!1}},[d]),P(()=>{if(!m){f(null);return}let u=!0;return(async()=>{try{let y=await h.get("/reals/index",{address:m});if(!u)return;f(y?.index??null),s("")}catch(y){u&&s(y.message||"request failed")}})(),()=>{u=!1}},[m]);let{points:c,error:$}=A({path:"/stats/real",body:g!==null?{index:g}:null}),w=k(()=>[{label:"v1",field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.2)"},{label:"v2",field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[]);return n`
      <main>
        <section className="card">
          <h2>Per-Real Stats</h2>
          <p className="muted">Select a VIP and real address to chart.</p>
          ${N&&n`<p className="error">${N}</p>`}
          <div className="form-row">
            <label className="field">
              <span>VIP</span>
              <select value=${d} onChange=${u=>r(u.target.value)}>
                ${e.map(u=>n`
                    <option value=${V(u)}>
                      ${u.address}:${u.port} / ${u.proto}
                    </option>
                  `)}
              </select>
            </label>
            <label className="field">
              <span>Real</span>
              <select
                value=${m}
                onChange=${u=>o(u.target.value)}
                disabled=${i.length===0}
              >
                ${i.map(u=>n`
                    <option value=${u.address}>${u.address}</option>
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
          ${$&&n`<p className="error">${$}</p>`}
          ${g===null?n`<p className="muted">Select a real to start polling.</p>`:n`<${B} points=${c} keys=${w} />`}
        </section>
      </main>
    `}function se(){let[e,t]=p({initialized:!1,ready:!1}),[d,r]=p([]),i=D({}),v=(o,g="info")=>{let f=`${Date.now()}-${Math.random().toString(16).slice(2)}`;r(N=>N.concat({id:f,message:o,kind:g})),i.current[f]=setTimeout(()=>{r(N=>N.filter(s=>s.id!==f)),delete i.current[f]},4e3)},m=o=>{i.current[o]&&(clearTimeout(i.current[o]),delete i.current[o]),r(g=>g.filter(f=>f.id!==o))};return P(()=>{let o=!0,g=async()=>{try{let N=await h.get("/lb/status");o&&t(N||{initialized:!1,ready:!1})}catch{o&&t({initialized:!1,ready:!1})}};g();let f=setInterval(g,5e3);return()=>{o=!1,clearInterval(f)}},[]),n`
      <${U}>
        <${O}>
          <${E.Provider} value=${{addToast:v}}>
            <${Z} status=${e} />
            <${H}>
              <${q} path="/" element=${n`<${Q} />`} />
              <${q} path="/vips/:vipId" element=${n`<${Y} />`} />
              <${q} path="/vips/:vipId/stats" element=${n`<${ee} />`} />
              <${q} path="/stats/global" element=${n`<${ae} />`} />
              <${q} path="/stats/real" element=${n`<${te} />`} />
            </${H}>
            <${K} toasts=${d} onDismiss=${m} />
          </${E.Provider}>
        </${O}>
      </${U}>
    `}ReactDOM.createRoot(document.getElementById("root")).render(n`<${se} />`)})();})();
