(()=>{(()=>{let{useEffect:P,useMemo:C,useRef:D,useState:f,useContext:j}=React,{BrowserRouter:U,Routes:O,Route:q,NavLink:k,Link:F,useParams:H,useNavigate:J}=ReactRouterDOM,o=htm.bind(React.createElement),E=React.createContext({addToast:()=>{}});function M(){return j(E)}let g={base:"/api/v1",async request(e,t={}){let d={method:t.method||"GET",headers:{"Content-Type":"application/json"}},l=`${g.base}${e}`;if(t.body!==void 0&&t.body!==null)if(d.method==="GET"){let i=new URLSearchParams;Object.entries(t.body).forEach(([h,c])=>{if(c!=null){if(Array.isArray(c)){c.forEach($=>i.append(h,String($)));return}if(typeof c=="object"){i.set(h,JSON.stringify(c));return}i.set(h,String(c))}});let n=i.toString();n&&(l+=`${l.includes("?")?"&":"?"}${n}`)}else d.body=JSON.stringify(t.body);let p=await fetch(l,d),b;try{b=await p.json()}catch{throw new Error("invalid JSON response")}if(!p.ok)throw new Error(b?.error?.message||`HTTP ${p.status}`);if(!b.success){let i=b.error?.message||"request failed";throw new Error(i)}return b.data},get(e,t){return g.request(e,{method:"GET",body:t})},post(e,t){return g.request(e,{method:"POST",body:t})},put(e,t){return g.request(e,{method:"PUT",body:t})},del(e,t){return g.request(e,{method:"DELETE",body:t})}};function V(e){return`${encodeURIComponent(e.address)}:${e.port}:${e.proto}`}function A(e){let t=e.split(":"),d=Number(t.pop()||0),l=Number(t.pop()||0);return{address:decodeURIComponent(t.join(":")),port:l,proto:d}}function X(e,t,d=[]){let[l,p]=f(null),[b,i]=f(""),[n,h]=f(!0);return P(()=>{let c=!0,$=async()=>{try{let u=await e();c&&(p(u),i(""),h(!1))}catch(u){c&&(i(u.message||"request failed"),h(!1))}};$();let s=setInterval($,t);return()=>{c=!1,clearInterval(s)}},d),{data:l,error:b,loading:n}}function L({path:e,body:t,intervalMs:d=1e3,limit:l=60}){let[p,b]=f([]),[i,n]=f(""),h=C(()=>JSON.stringify(t||{}),[t]);return P(()=>{if(t===null)return b([]),n(""),()=>{};let c=!0,$=async()=>{try{let u=await g.get(e,t);if(!c)return;let v=new Date().toLocaleTimeString();b(w=>w.concat({label:v,...u}).slice(-l)),n("")}catch(u){c&&n(u.message||"request failed")}};$();let s=setInterval($,d);return()=>{c=!1,clearInterval(s)}},[e,h,d,l]),{points:p,error:i}}function B({title:e,points:t,keys:d}){let l=D(null),p=D(null);return P(()=>{if(!l.current)return;p.current||(p.current=new Chart(l.current,{type:"line",data:{labels:[],datasets:[]},options:{responsive:!0,animation:!1,scales:{x:{grid:{display:!1}},y:{beginAtZero:!0}},plugins:{legend:{display:!0,position:"bottom"},title:{display:!!e,text:e}}}}));let b=p.current,i=t.map(n=>n.label);return b.data.labels=i,b.data.datasets=d.map(n=>({label:n.label,data:t.map(h=>h[n.field]||0),borderColor:n.color,backgroundColor:n.fill,borderWidth:2,tension:.3})),b.update(),()=>{}},[t,d,e]),P(()=>()=>{p.current&&(p.current.destroy(),p.current=null)},[]),o`<canvas ref=${l} height="120"></canvas>`}function W({children:e}){return e}function K({toasts:e,onDismiss:t}){return o`
      <div className="toast-stack">
        ${e.map(d=>o`
            <div className=${`toast ${d.kind}`}>
              <span>${d.message}</span>
              <button className="toast-close" onClick=${()=>t(d.id)}>
                ×
              </button>
            </div>
          `)}
      </div>
    `}function Z({status:e}){return o`
      <header>
        <div>
          <div style=${{fontSize:20,fontWeight:700}}>Vatran</div>
          <div className="status-pill">
            <span className=${`dot ${e.ready?"ok":""}`}></span>
            ${e.ready?"Ready":"Not ready"}
          </div>
        </div>
        <nav>
          <${k} to="/" end className=${({isActive:t})=>t?"active":""}>
            Dashboard
          </${k}>
          <${k} to="/stats/global" className=${({isActive:t})=>t?"active":""}>
            Global stats
          </${k}>
          <${k} to="/stats/real" className=${({isActive:t})=>t?"active":""}>
            Per-real stats
          </${k}>
        </nav>
      </header>
    `}function Q(){let{addToast:e}=M(),[t,d]=f({initialized:!1,ready:!1}),[l,p]=f([]),[b,i]=f(""),[n,h]=f(!1),[c,$]=f(!1),[s,u]=f({main_interface:"",balancer_prog_path:"",healthchecking_prog_path:"",default_mac:"",local_mac:"",root_map_path:"",root_map_pos:2,katran_src_v4:"",katran_src_v6:"",use_root_map:!1,max_vips:1024,max_reals:4096,hash_func:0}),[v,w]=f({address:"",port:80,proto:6,flags:0}),m=async()=>{try{let a=await g.get("/lb/status"),I=await g.get("/vips");d(a||{initialized:!1,ready:!1}),p(I||[]),i("")}catch(a){i(a.message||"request failed")}};P(()=>{let a=!0;return(async()=>{a&&await m()})(),()=>{a=!1}},[]);let R=async a=>{a.preventDefault();try{let I={...s,root_map_pos:s.root_map_pos===""?void 0:Number(s.root_map_pos),max_vips:Number(s.max_vips),max_reals:Number(s.max_reals),hash_func:Number(s.hash_func)};await g.post("/lb/create",I),i(""),h(!1),e("Load balancer initialized.","success"),await m()}catch(I){i(I.message||"request failed"),e(I.message||"Initialize failed.","error")}},y=async a=>{a.preventDefault();try{await g.post("/vips",{...v,port:Number(v.port),proto:Number(v.proto),flags:Number(v.flags||0)}),w({address:"",port:80,proto:6,flags:0}),i(""),$(!1),e("VIP created.","success"),await m()}catch(I){i(I.message||"request failed"),e(I.message||"VIP create failed.","error")}},_=async()=>{try{await g.post("/lb/load-bpf-progs"),i(""),e("BPF programs loaded.","success"),await m()}catch(a){i(a.message||"request failed"),e(a.message||"Load BPF programs failed.","error")}},x=async()=>{try{await g.post("/lb/attach-bpf-progs"),i(""),e("BPF programs attached.","success"),await m()}catch(a){i(a.message||"request failed"),e(a.message||"Attach BPF programs failed.","error")}};return o`
      <main>
        <section className="card">
          <h2>Load balancer</h2>
          <p>Initialized: ${t.initialized?"yes":"no"}</p>
          <p>Ready: ${t.ready?"yes":"no"}</p>
          <div className="row">
            <button className="btn" onClick=${()=>h(a=>!a)}>
              ${n?"Close":"Initialize"}
            </button>
            <button className="btn secondary" onClick=${()=>$(a=>!a)}>
              ${c?"Close":"Create VIP"}
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
          ${n&&o`
            <form className="form" onSubmit=${R}>
              <div className="form-row">
                <label className="field">
                  <span>Main interface</span>
                  <input
                    value=${s.main_interface}
                    onInput=${a=>u({...s,main_interface:a.target.value})}
                    placeholder="eth0"
                    required
                  />
                </label>
                <label className="field">
                  <span>Balancer prog path</span>
                  <input
                    value=${s.balancer_prog_path}
                    onInput=${a=>u({...s,balancer_prog_path:a.target.value})}
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
                    onInput=${a=>u({...s,healthchecking_prog_path:a.target.value})}
                    placeholder="/path/to/hc.o"
                  />
                </label>
                <label className="field">
                  <span>Default MAC</span>
                  <input
                    value=${s.default_mac}
                    onInput=${a=>u({...s,default_mac:a.target.value})}
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
                    onInput=${a=>u({...s,local_mac:a.target.value})}
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
                    onInput=${a=>u({...s,hash_func:a.target.value})}
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Root map path</span>
                  <input
                    value=${s.root_map_path}
                    onInput=${a=>u({...s,root_map_path:a.target.value})}
                    placeholder="/sys/fs/bpf/root_map"
                  />
                </label>
                <label className="field">
                  <span>Root map position</span>
                  <input
                    type="number"
                    min="0"
                    value=${s.root_map_pos}
                    onInput=${a=>u({...s,root_map_pos:a.target.value})}
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Katran src v4</span>
                  <input
                    value=${s.katran_src_v4}
                    onInput=${a=>u({...s,katran_src_v4:a.target.value})}
                    placeholder="10.0.0.1"
                  />
                </label>
                <label className="field">
                  <span>Katran src v6</span>
                  <input
                    value=${s.katran_src_v6}
                    onInput=${a=>u({...s,katran_src_v6:a.target.value})}
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
                    onInput=${a=>u({...s,max_vips:a.target.value})}
                  />
                </label>
                <label className="field">
                  <span>Max reals</span>
                  <input
                    type="number"
                    min="1"
                    value=${s.max_reals}
                    onInput=${a=>u({...s,max_reals:a.target.value})}
                  />
                </label>
              </div>
              <label className="field checkbox">
                <input
                  type="checkbox"
                  checked=${s.use_root_map}
                  onChange=${a=>u({...s,use_root_map:a.target.checked})}
                />
                <span>Use root map</span>
              </label>
              <button className="btn" type="submit">Initialize LB</button>
            </form>
          `}
          ${c&&o`
            <form className="form" onSubmit=${y}>
              <div className="form-row">
                <label className="field">
                  <span>VIP address</span>
                  <input
                    value=${v.address}
                    onInput=${a=>w({...v,address:a.target.value})}
                    placeholder="1.2.3.4"
                    required
                  />
                </label>
                <label className="field">
                  <span>Port</span>
                  <input
                    type="number"
                    value=${v.port}
                    onInput=${a=>w({...v,port:a.target.value})}
                    required
                  />
                </label>
                <label className="field">
                  <span>Protocol</span>
                  <select
                    value=${v.proto}
                    onChange=${a=>w({...v,proto:a.target.value})}
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
                    onInput=${a=>w({...v,flags:a.target.value})}
                  />
                </label>
              </div>
              <button className="btn" type="submit">Create VIP</button>
            </form>
          `}
          ${b&&o`<p className="error">${b}</p>`}
        </section>
        <section className="card">
          <div className="section-header">
            <h2>VIPs</h2>
            <button className="btn ghost" onClick=${m}>Refresh</button>
          </div>
          ${l.length===0?o`<p className="muted">No VIPs configured yet.</p>`:o`
                <div className="grid">
                  ${l.map(a=>o`
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
    `}function Y(){let{addToast:e}=M(),t=H(),d=J(),l=C(()=>A(t.vipId),[t.vipId]),[p,b]=f([]),[i,n]=f(""),[h,c]=f(""),[$,s]=f(!0),[u,v]=f({address:"",weight:100,flags:0}),[w,m]=f({}),[R,y]=f(null),[_,x]=f({flag:0,set:!0}),[a,I]=f({hash_function:0}),T=async()=>{try{let r=await g.get("/vips/reals",l);b(r||[]);let N={};(r||[]).forEach(G=>{N[G.address]=G.weight}),m(N),n(""),s(!1)}catch(r){n(r.message||"request failed"),s(!1)}},z=async()=>{try{let r=await g.get("/vips/flags",l);y(r?.flags??0),c("")}catch(r){c(r.message||"request failed")}};P(()=>{T(),z()},[t.vipId]);let re=async r=>{try{let N=Number(w[r.address]);await g.post("/vips/reals",{vip:l,real:{address:r.address,weight:N,flags:r.flags||0}}),await T(),e("Real weight updated.","success")}catch(N){n(N.message||"request failed"),e(N.message||"Update failed.","error")}},le=async r=>{try{await g.del("/vips/reals",{vip:l,real:{address:r.address,weight:r.weight,flags:r.flags||0}}),await T(),e("Real removed.","success")}catch(N){n(N.message||"request failed"),e(N.message||"Remove failed.","error")}},ne=async r=>{r.preventDefault();try{await g.post("/vips/reals",{vip:l,real:{address:u.address,weight:Number(u.weight),flags:Number(u.flags||0)}}),v({address:"",weight:100,flags:0}),await T(),e("Real added.","success")}catch(N){n(N.message||"request failed"),e(N.message||"Add failed.","error")}},oe=async()=>{try{await g.del("/vips",l),e("VIP deleted.","success"),d("/")}catch(r){n(r.message||"request failed"),e(r.message||"Delete failed.","error")}},ie=async r=>{r.preventDefault();try{await g.put("/vips/flags",{...l,flag:Number(_.flag),set:!!_.set}),await z(),e("VIP flags updated.","success")}catch(N){c(N.message||"request failed"),e(N.message||"Flag update failed.","error")}},ce=async r=>{r.preventDefault();try{await g.put("/vips/hash-function",{...l,hash_function:Number(a.hash_function)}),e("Hash function updated.","success")}catch(N){c(N.message||"request failed"),e(N.message||"Hash update failed.","error")}};return o`
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
          ${i&&o`<p className="error">${i}</p>`}
          ${h&&o`<p className="error">${h}</p>`}
          ${$?o`<p className="muted">Loading reals…</p>`:o`
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
                    ${p.map(r=>o`
                        <tr>
                          <td>
                            <span
                              className=${`dot ${Number(r.weight)>0?"ok":"bad"}`}
                            ></span>
                          </td>
                          <td>${r.address}</td>
                          <td>
                            <input
                              className="inline-input"
                              type="number"
                              min="0"
                              value=${w[r.address]??r.weight}
                              onInput=${N=>m({...w,[r.address]:N.target.value})}
                            />
                          </td>
                          <td>${r.flags||0}</td>
                          <td className="row">
                            <button className="btn" onClick=${()=>re(r)}>
                              Update
                            </button>
                            <button className="btn ghost" onClick=${()=>le(r)}>
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
                    onInput=${r=>x({..._,flag:r.target.value})}
                  />
                </label>
                <label className="field">
                  <span>Set</span>
                  <select
                    value=${String(_.set)}
                    onChange=${r=>x({..._,set:r.target.value==="true"})}
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
                    onInput=${r=>I({...a,hash_function:r.target.value})}
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
                  value=${u.address}
                  onInput=${r=>v({...u,address:r.target.value})}
                  placeholder="10.0.0.1"
                  required
                />
              </label>
              <label className="field">
                <span>Weight</span>
                <input
                  type="number"
                  min="0"
                  value=${u.weight}
                  onInput=${r=>v({...u,weight:r.target.value})}
                  required
                />
              </label>
              <label className="field">
                <span>Flags</span>
                <input
                  type="number"
                  value=${u.flags}
                  onInput=${r=>v({...u,flags:r.target.value})}
                />
              </label>
            </div>
            <button className="btn" type="submit">Add real</button>
          </form>
        </section>
      </main>
    `}function ee(){let e=H(),t=C(()=>A(e.vipId),[e.vipId]),{points:d,error:l}=L({path:"/stats/vip",body:t}),p=C(()=>[{label:"v1",field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.15)"},{label:"v2",field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[]);return o`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>VIP Stats</h2>
              <p className="muted">${t.address}:${t.port} / ${t.proto}</p>
            </div>
          </div>
          ${l&&o`<p className="error">${l}</p>`}
          <${B} title="Traffic" points=${d} keys=${p} />
        </section>
      </main>
    `}function S({title:e,path:t}){let{points:d,error:l}=L({path:t}),p=C(()=>[{label:"v1",field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.2)"},{label:"v2",field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[]);return o`
      <div className="card">
        <h3>${e}</h3>
        ${l&&o`<p className="error">${l}</p>`}
        <${B} points=${d} keys=${p} />
      </div>
    `}function ae(){let{data:e,error:t}=X(()=>g.get("/stats/userspace"),1e3,[]);return o`
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
          ${t&&o`<p className="error">${t}</p>`}
          ${e?o`
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
              `:o`<p className="muted">Waiting for data…</p>`}
        </section>
      </main>
    `}function te(){let[e,t]=f([]),[d,l]=f(""),[p,b]=f([]),[i,n]=f(""),[h,c]=f(null),[$,s]=f("");P(()=>{let m=!0;return(async()=>{try{let y=await g.get("/vips");if(!m)return;t(y||[]),!d&&y&&y.length>0&&l(V(y[0]))}catch(y){m&&s(y.message||"request failed")}})(),()=>{m=!1}},[]),P(()=>{if(!d)return;let m=A(d),R=!0;return(async()=>{try{let _=await g.get("/vips/reals",m);if(!R)return;b(_||[]),_&&_.length>0?n(x=>x||_[0].address):n(""),s("")}catch(_){R&&s(_.message||"request failed")}})(),()=>{R=!1}},[d]),P(()=>{if(!i){c(null);return}let m=!0;return(async()=>{try{let y=await g.get("/reals/index",{address:i});if(!m)return;c(y?.index??null),s("")}catch(y){m&&s(y.message||"request failed")}})(),()=>{m=!1}},[i]);let{points:u,error:v}=L({path:"/stats/real",body:h!==null?{index:h}:null}),w=C(()=>[{label:"v1",field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.2)"},{label:"v2",field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[]);return o`
      <main>
        <section className="card">
          <h2>Per-Real Stats</h2>
          <p className="muted">Select a VIP and real address to chart.</p>
          ${$&&o`<p className="error">${$}</p>`}
          <div className="form-row">
            <label className="field">
              <span>VIP</span>
              <select value=${d} onChange=${m=>l(m.target.value)}>
                ${e.map(m=>o`
                    <option value=${V(m)}>
                      ${m.address}:${m.port} / ${m.proto}
                    </option>
                  `)}
              </select>
            </label>
            <label className="field">
              <span>Real</span>
              <select
                value=${i}
                onChange=${m=>n(m.target.value)}
                disabled=${p.length===0}
              >
                ${p.map(m=>o`
                    <option value=${m.address}>${m.address}</option>
                  `)}
              </select>
            </label>
            <label className="field">
              <span>Index</span>
              <input value=${h??""} readOnly />
            </label>
          </div>
        </section>
        <section className="card">
          <h3>Real stats</h3>
          ${v&&o`<p className="error">${v}</p>`}
          ${h===null?o`<p className="muted">Select a real to start polling.</p>`:o`<${B} points=${u} keys=${w} />`}
        </section>
      </main>
    `}function se(){let[e,t]=f({initialized:!1,ready:!1}),[d,l]=f([]),p=D({}),b=(n,h="info")=>{let c=`${Date.now()}-${Math.random().toString(16).slice(2)}`;l($=>$.concat({id:c,message:n,kind:h})),p.current[c]=setTimeout(()=>{l($=>$.filter(s=>s.id!==c)),delete p.current[c]},4e3)},i=n=>{p.current[n]&&(clearTimeout(p.current[n]),delete p.current[n]),l(h=>h.filter(c=>c.id!==n))};return P(()=>{let n=!0,h=async()=>{try{let $=await g.get("/lb/status");n&&t($||{initialized:!1,ready:!1})}catch{n&&t({initialized:!1,ready:!1})}};h();let c=setInterval(h,5e3);return()=>{n=!1,clearInterval(c)}},[]),o`
      <${U}>
        <${W}>
          <${E.Provider} value=${{addToast:b}}>
            <${Z} status=${e} />
            <${O}>
              <${q} path="/" element=${o`<${Q} />`} />
              <${q} path="/vips/:vipId" element=${o`<${Y} />`} />
              <${q} path="/vips/:vipId/stats" element=${o`<${ee} />`} />
              <${q} path="/stats/global" element=${o`<${ae} />`} />
              <${q} path="/stats/real" element=${o`<${te} />`} />
            </${O}>
            <${K} toasts=${d} onDismiss=${i} />
          </${E.Provider}>
        </${W}>
      </${U}>
    `}ReactDOM.createRoot(document.getElementById("root")).render(o`<${se} />`)})();})();
