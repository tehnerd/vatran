(()=>{(()=>{let{useEffect:C,useMemo:S,useRef:T,useState:f,useContext:X}=React,{BrowserRouter:U,Routes:M,Route:k,NavLink:R,Link:F,useParams:O,useNavigate:Z}=ReactRouterDOM,s=htm.bind(React.createElement),D=React.createContext({addToast:()=>{}});function L(){return X(D)}let b={base:"/api/v1",async request(e,a={}){let n={method:a.method||"GET",headers:{"Content-Type":"application/json"}},l=`${b.base}${e}`;if(a.body!==void 0&&a.body!==null)if(n.method==="GET"){let o=new URLSearchParams;Object.entries(a.body).forEach(([p,c])=>{if(c!=null){if(Array.isArray(c)){c.forEach(g=>o.append(p,String(g)));return}if(typeof c=="object"){o.set(p,JSON.stringify(c));return}o.set(p,String(c))}});let i=o.toString();i&&(l+=`${l.includes("?")?"&":"?"}${i}`)}else n.body=JSON.stringify(a.body);let m=await fetch(l,n),$;try{$=await m.json()}catch{throw new Error("invalid JSON response")}if(!m.ok)throw new Error($?.error?.message||`HTTP ${m.status}`);if(!$.success){let o=$.error?.message||"request failed";throw new Error(o)}return $.data},get(e,a){return b.request(e,{method:"GET",body:a})},post(e,a){return b.request(e,{method:"POST",body:a})},put(e,a){return b.request(e,{method:"PUT",body:a})},del(e,a){return b.request(e,{method:"DELETE",body:a})}};function q(e){return`${encodeURIComponent(e.address)}:${e.port}:${e.proto}`}function A(e){let a=e.split(":"),n=Number(a.pop()||0),l=Number(a.pop()||0);return{address:decodeURIComponent(a.join(":")),port:l,proto:n}}function K(e,a,n=[]){let[l,m]=f(null),[$,o]=f(""),[i,p]=f(!0);return C(()=>{let c=!0,g=async()=>{try{let u=await e();c&&(m(u),o(""),p(!1))}catch(u){c&&(o(u.message||"request failed"),p(!1))}};g();let r=setInterval(g,a);return()=>{c=!1,clearInterval(r)}},n),{data:l,error:$,loading:i}}function V({path:e,body:a,intervalMs:n=1e3,limit:l=60}){let[m,$]=f([]),[o,i]=f(""),p=S(()=>JSON.stringify(a||{}),[a]);return C(()=>{if(a===null)return $([]),i(""),()=>{};let c=!0,g=async()=>{try{let u=await b.get(e,a);if(!c)return;let v=new Date().toLocaleTimeString();$(y=>y.concat({label:v,...u}).slice(-l)),i("")}catch(u){c&&i(u.message||"request failed")}};g();let r=setInterval(g,n);return()=>{c=!1,clearInterval(r)}},[e,p,n,l]),{points:m,error:o}}function H({title:e,points:a,keys:n,diff:l=!1,height:m=120,showTitle:$=!1}){let o=T(null),i=T(null);return C(()=>{if(!o.current)return;i.current||(i.current=new Chart(o.current,{type:"line",data:{labels:[],datasets:[]},options:{responsive:!0,maintainAspectRatio:!1,animation:!1,scales:{x:{grid:{display:!1}},y:{beginAtZero:!l}},plugins:{legend:{display:!0,position:"bottom"},title:{display:$&&!!e,text:e}}}}));let p=i.current,c=a.map(g=>g.label);return p.data.labels=c,p.data.datasets=n.map(g=>{let r=a.map(v=>v[g.field]||0),u=l?r.map((v,y)=>y===0?0:v-r[y-1]):r;return{label:g.label,data:u,borderColor:g.color,backgroundColor:g.fill,borderWidth:2,tension:.3}}),p.options.scales.y.beginAtZero=!l,p.options.plugins.title.display=$&&!!e,p.options.plugins.title.text=e||"",p.update(),()=>{}},[a,n,e,l,$]),C(()=>()=>{i.current&&(i.current.destroy(),i.current=null)},[]),s`<canvas ref=${o} height=${m}></canvas>`}function z({title:e,points:a,keys:n,diff:l=!1,inlineTitle:m=!0}){let[$,o]=f(!1);return s`
      <div className="chart-wrap">
        <button className="chart-zoom-button" type="button" onClick=${()=>o(!0)}>
          <span className="zoom-icon" aria-hidden="true">+</span>
          Zoom
        </button>
        <div className="chart-click" onClick=${()=>o(!0)}>
          <${H}
            title=${e}
            points=${a}
            keys=${n}
            diff=${l}
            height=${120}
            showTitle=${m&&!!e}
          />
        </div>
        ${$&&s`
          <div className="chart-overlay" onClick=${()=>o(!1)}>
            <div className="chart-modal" onClick=${i=>i.stopPropagation()}>
              <div className="row chart-modal-header">
                <div>
                  <h3>${e||"Chart"}</h3>
                  ${l?s`<p className="muted">Per-second delta.</p>`:""}
                </div>
                <button className="btn ghost" onClick=${()=>o(!1)}>
                  Close
                </button>
              </div>
              <div className="chart-zoom">
                <${H}
                  title=${e}
                  points=${a}
                  keys=${n}
                  diff=${l}
                  height=${360}
                  showTitle=${!1}
                />
              </div>
            </div>
          </div>
        `}
      </div>
    `}function G({children:e}){return e}function Y({toasts:e,onDismiss:a}){return s`
      <div className="toast-stack">
        ${e.map(n=>s`
            <div className=${`toast ${n.kind}`}>
              <span>${n.message}</span>
              <button className="toast-close" onClick=${()=>a(n.id)}>
                ×
              </button>
            </div>
          `)}
      </div>
    `}function Q({status:e}){return s`
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
    `}function ee(){let{addToast:e}=L(),[a,n]=f({initialized:!1,ready:!1}),[l,m]=f([]),[$,o]=f(""),[i,p]=f(!1),[c,g]=f(!1),[r,u]=f({main_interface:"",balancer_prog_path:"",healthchecking_prog_path:"",default_mac:"",local_mac:"",root_map_path:"",root_map_pos:2,katran_src_v4:"",katran_src_v6:"",use_root_map:!1,max_vips:512,max_reals:4096,hash_func:0}),[v,y]=f({address:"",port:80,proto:6,flags:0}),h=async()=>{try{let t=await b.get("/lb/status"),I=await b.get("/vips");n(t||{initialized:!1,ready:!1}),m(I||[]),o("")}catch(t){o(t.message||"request failed")}};C(()=>{let t=!0;return(async()=>{t&&await h()})(),()=>{t=!1}},[]);let P=async t=>{t.preventDefault();try{let I={...r,root_map_pos:r.root_map_pos===""?void 0:Number(r.root_map_pos),max_vips:Number(r.max_vips),max_reals:Number(r.max_reals),hash_func:Number(r.hash_func)};await b.post("/lb/create",I),o(""),p(!1),e("Load balancer initialized.","success"),await h()}catch(I){o(I.message||"request failed"),e(I.message||"Initialize failed.","error")}},w=async t=>{t.preventDefault();try{await b.post("/vips",{...v,port:Number(v.port),proto:Number(v.proto),flags:Number(v.flags||0)}),y({address:"",port:80,proto:6,flags:0}),o(""),g(!1),e("VIP created.","success"),await h()}catch(I){o(I.message||"request failed"),e(I.message||"VIP create failed.","error")}},_=async()=>{try{await b.post("/lb/load-bpf-progs"),o(""),e("BPF programs loaded.","success"),await h()}catch(t){o(t.message||"request failed"),e(t.message||"Load BPF programs failed.","error")}},x=async()=>{try{await b.post("/lb/attach-bpf-progs"),o(""),e("BPF programs attached.","success"),await h()}catch(t){o(t.message||"request failed"),e(t.message||"Attach BPF programs failed.","error")}};return s`
      <main>
        <section className="card">
          <h2>Load balancer</h2>
          <p>Initialized: ${a.initialized?"yes":"no"}</p>
          <p>Ready: ${a.ready?"yes":"no"}</p>
          <div className="row">
            ${!a.initialized&&s`
              <button className="btn" onClick=${()=>p(t=>!t)}>
                ${i?"Close":"Initialize"}
              </button>
            `}
            <button className="btn secondary" onClick=${()=>g(t=>!t)}>
              ${c?"Close":"Create VIP"}
            </button>
          </div>
          ${!a.ready&&s`
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
                onClick=${x}
              >
                Attach BPF Programs
              </button>
            </div>
          `}
          ${i&&s`
            <form className="form" onSubmit=${P}>
              <div className="form-row">
                <label className="field">
                  <span>Main interface</span>
                  <input
                    value=${r.main_interface}
                    onInput=${t=>u({...r,main_interface:t.target.value})}
                    placeholder="eth0"
                    required
                  />
                </label>
                <label className="field">
                  <span>Balancer prog path</span>
                  <input
                    value=${r.balancer_prog_path}
                    onInput=${t=>u({...r,balancer_prog_path:t.target.value})}
                    placeholder="/path/to/bpf.o"
                    required
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Healthchecking path</span>
                  <input
                    value=${r.healthchecking_prog_path}
                    onInput=${t=>u({...r,healthchecking_prog_path:t.target.value})}
                    placeholder="/path/to/hc.o"
                  />
                </label>
                <label className="field">
                  <span>Default MAC</span>
                  <input
                    value=${r.default_mac}
                    onInput=${t=>u({...r,default_mac:t.target.value})}
                    placeholder="aa:bb:cc:dd:ee:ff"
                    required
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Local MAC</span>
                  <input
                    value=${r.local_mac}
                    onInput=${t=>u({...r,local_mac:t.target.value})}
                    placeholder="11:22:33:44:55:66"
                    required=${r.healthchecking_prog_path?.trim()!==""}
                  />
                </label>
                <label className="field">
                  <span>Hash function</span>
                  <input
                    type="number"
                    min="0"
                    value=${r.hash_func}
                    onInput=${t=>u({...r,hash_func:t.target.value})}
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Root map path</span>
                  <input
                    value=${r.root_map_path}
                    onInput=${t=>u({...r,root_map_path:t.target.value})}
                    placeholder="/sys/fs/bpf/root_map"
                  />
                </label>
                <label className="field">
                  <span>Root map position</span>
                  <input
                    type="number"
                    min="0"
                    value=${r.root_map_pos}
                    onInput=${t=>u({...r,root_map_pos:t.target.value})}
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Katran src v4</span>
                  <input
                    value=${r.katran_src_v4}
                    onInput=${t=>u({...r,katran_src_v4:t.target.value})}
                    placeholder="10.0.0.1"
                  />
                </label>
                <label className="field">
                  <span>Katran src v6</span>
                  <input
                    value=${r.katran_src_v6}
                    onInput=${t=>u({...r,katran_src_v6:t.target.value})}
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
                    value=${r.max_vips}
                    onInput=${t=>u({...r,max_vips:t.target.value})}
                  />
                </label>
                <label className="field">
                  <span>Max reals</span>
                  <input
                    type="number"
                    min="1"
                    value=${r.max_reals}
                    onInput=${t=>u({...r,max_reals:t.target.value})}
                  />
                </label>
              </div>
              <label className="field checkbox">
                <input
                  type="checkbox"
                  checked=${r.use_root_map}
                  onChange=${t=>u({...r,use_root_map:t.target.checked})}
                />
                <span>Use root map</span>
              </label>
              <button className="btn" type="submit">Initialize LB</button>
            </form>
          `}
          ${c&&s`
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
          ${$&&s`<p className="error">${$}</p>`}
        </section>
        <section className="card">
          <div className="section-header">
            <h2>VIPs</h2>
            <button className="btn ghost" onClick=${h}>Refresh</button>
          </div>
          ${l.length===0?s`<p className="muted">No VIPs configured yet.</p>`:s`
                <div className="grid">
                  ${l.map(t=>s`
                      <div className="card">
                        <div style=${{fontWeight:600}}>
                          ${t.address}:${t.port} / ${t.proto}
                        </div>
                        <div className="muted" style=${{marginTop:6}}>
                          Flags: ${t.flags||0}
                        </div>
                        <div className="row" style=${{marginTop:12}}>
                          <${F} className="btn" to=${`/vips/${q(t)}`}>
                            Open
                          </${F}>
                          <${F}
                            className="btn secondary"
                            to=${`/vips/${q(t)}/stats`}
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
    `}function ae(){let{addToast:e}=L(),a=O(),n=Z(),l=S(()=>A(a.vipId),[a.vipId]),[m,$]=f([]),[o,i]=f(""),[p,c]=f(""),[g,r]=f(!0),[u,v]=f({address:"",weight:100,flags:0}),[y,h]=f({}),[P,w]=f(null),[_,x]=f({flag:0,set:!0}),[t,I]=f({hash_function:0}),E=async()=>{try{let d=await b.get("/vips/reals",l);$(d||[]);let N={};(d||[]).forEach(J=>{N[J.address]=J.weight}),h(N),i(""),r(!1)}catch(d){i(d.message||"request failed"),r(!1)}},B=async()=>{try{let d=await b.get("/vips/flags",l);w(d?.flags??0),c("")}catch(d){c(d.message||"request failed")}};C(()=>{E(),B()},[a.vipId]);let ce=async d=>{try{let N=Number(y[d.address]);await b.post("/vips/reals",{vip:l,real:{address:d.address,weight:N,flags:d.flags||0}}),await E(),e("Real weight updated.","success")}catch(N){i(N.message||"request failed"),e(N.message||"Update failed.","error")}},de=async d=>{try{await b.del("/vips/reals",{vip:l,real:{address:d.address,weight:d.weight,flags:d.flags||0}}),await E(),e("Real removed.","success")}catch(N){i(N.message||"request failed"),e(N.message||"Remove failed.","error")}},ue=async d=>{d.preventDefault();try{await b.post("/vips/reals",{vip:l,real:{address:u.address,weight:Number(u.weight),flags:Number(u.flags||0)}}),v({address:"",weight:100,flags:0}),await E(),e("Real added.","success")}catch(N){i(N.message||"request failed"),e(N.message||"Add failed.","error")}},pe=async()=>{try{await b.del("/vips",l),e("VIP deleted.","success"),n("/")}catch(d){i(d.message||"request failed"),e(d.message||"Delete failed.","error")}},me=async d=>{d.preventDefault();try{await b.put("/vips/flags",{...l,flag:Number(_.flag),set:!!_.set}),await B(),e("VIP flags updated.","success")}catch(N){c(N.message||"request failed"),e(N.message||"Flag update failed.","error")}},fe=async d=>{d.preventDefault();try{await b.put("/vips/hash-function",{...l,hash_function:Number(t.hash_function)}),e("Hash function updated.","success")}catch(N){c(N.message||"request failed"),e(N.message||"Hash update failed.","error")}};return s`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>VIP Detail</h2>
              <p className="muted">${l.address}:${l.port} / ${l.proto}</p>
              <p className="muted">Flags: ${P??"\u2014"}</p>
            </div>
            <div className="row">
              <button className="btn ghost" onClick=${B}>Refresh flags</button>
              <button className="btn danger" onClick=${pe}>Delete VIP</button>
            </div>
          </div>
          ${o&&s`<p className="error">${o}</p>`}
          ${p&&s`<p className="error">${p}</p>`}
          ${g?s`<p className="muted">Loading reals…</p>`:s`
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
                    ${m.map(d=>s`
                        <tr>
                          <td>
                            <span
                              className=${`dot ${Number(d.weight)>0?"ok":"bad"}`}
                            ></span>
                          </td>
                          <td>${d.address}</td>
                          <td>
                            <input
                              className="inline-input"
                              type="number"
                              min="0"
                              value=${y[d.address]??d.weight}
                              onInput=${N=>h({...y,[d.address]:N.target.value})}
                            />
                          </td>
                          <td>${d.flags||0}</td>
                          <td className="row">
                            <button className="btn" onClick=${()=>ce(d)}>
                              Update
                            </button>
                            <button className="btn ghost" onClick=${()=>de(d)}>
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
            <form className="form" onSubmit=${me}>
              <div className="form-row">
                <label className="field">
                  <span>Flag</span>
                  <input
                    type="number"
                    value=${_.flag}
                    onInput=${d=>x({..._,flag:d.target.value})}
                  />
                </label>
                <label className="field">
                  <span>Set</span>
                  <select
                    value=${String(_.set)}
                    onChange=${d=>x({..._,set:d.target.value==="true"})}
                  >
                    <option value="true">Set</option>
                    <option value="false">Clear</option>
                  </select>
                </label>
              </div>
              <button className="btn" type="submit">Apply flag</button>
            </form>
            <form className="form" onSubmit=${fe}>
              <div className="form-row">
                <label className="field">
                  <span>Hash function</span>
                  <input
                    type="number"
                    min="0"
                    value=${t.hash_function}
                    onInput=${d=>I({...t,hash_function:d.target.value})}
                  />
                </label>
              </div>
              <button className="btn" type="submit">Apply hash</button>
            </form>
          </div>
        </section>
        <section className="card">
          <h3>Add real</h3>
          <form className="form" onSubmit=${ue}>
            <div className="form-row">
              <label className="field">
                <span>Real address</span>
                <input
                  value=${u.address}
                  onInput=${d=>v({...u,address:d.target.value})}
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
                  onInput=${d=>v({...u,weight:d.target.value})}
                  required
                />
              </label>
              <label className="field">
                <span>Flags</span>
                <input
                  type="number"
                  value=${u.flags}
                  onInput=${d=>v({...u,flags:d.target.value})}
                />
              </label>
            </div>
            <button className="btn" type="submit">Add real</button>
          </form>
        </section>
      </main>
    `}function te(){let e=O(),a=S(()=>A(e.vipId),[e.vipId]),{points:n,error:l}=V({path:"/stats/vip",body:a}),m=S(()=>[{label:"v1",field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.15)"},{label:"v2",field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[]);return s`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>VIP Stats</h2>
              <p className="muted">${a.address}:${a.port} / ${a.proto}</p>
            </div>
          </div>
          ${l&&s`<p className="error">${l}</p>`}
          <${z} title="Traffic" points=${n} keys=${m} />
        </section>
      </main>
    `}let W=[{title:"LRU",path:"/stats/lru"},{title:"LRU Miss",path:"/stats/lru/miss"},{title:"LRU Fallback",path:"/stats/lru/fallback"},{title:"LRU Global",path:"/stats/lru/global"},{title:"XDP Total",path:"/stats/xdp/total"},{title:"XDP Pass",path:"/stats/xdp/pass"},{title:"XDP Drop",path:"/stats/xdp/drop"},{title:"XDP Tx",path:"/stats/xdp/tx"}];function j(e){return Number.isFinite(e)?`${e>0?"+":""}${e}`:"0"}function se({title:e,path:a,diff:n=!1}){let{points:l,error:m}=V({path:a}),$=S(()=>[{label:"v1",field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.2)"},{label:"v2",field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[]);return s`
      <div className="card">
        <h3>${e}</h3>
        ${m&&s`<p className="error">${m}</p>`}
        <${z} title=${e} points=${l} keys=${$} diff=${n} inlineTitle=${!1} />
      </div>
    `}function re({title:e,path:a}){let{points:n,error:l}=V({path:a}),m=n[n.length-1]||{},$=n[n.length-2]||{},o=Number(m.v1??0),i=Number(m.v2??0),p=o-Number($.v1??0),c=i-Number($.v2??0);return s`
      <div className="summary-card">
        <div className="summary-title">${e}</div>
        ${l?s`<p className="error">${l}</p>`:s`
              <div className="summary-row">
                <div className="stat">
                  <span className="muted">v1 absolute</span>
                  <strong>${o}</strong>
                </div>
                <div className="stat">
                  <span className="muted">v1 delta/sec</span>
                  <strong className=${p<0?"delta down":"delta up"}>
                    ${j(p)}
                  </strong>
                </div>
              </div>
              <div className="summary-row">
                <div className="stat">
                  <span className="muted">v2 absolute</span>
                  <strong>${i}</strong>
                </div>
                <div className="stat">
                  <span className="muted">v2 delta/sec</span>
                  <strong className=${c<0?"delta down":"delta up"}>
                    ${j(c)}
                  </strong>
                </div>
              </div>
            `}
      </div>
    `}function le(){let{data:e,error:a}=K(()=>b.get("/stats/userspace"),1e3,[]);return s`
      <main>
        <section className="card">
          <h2>Global Stats</h2>
          <p className="muted">Polling every second. Charts show per-second deltas.</p>
        </section>
        <section className="grid">
          ${W.map(n=>s`<${se} title=${n.title} path=${n.path} diff=${!0} />`)}
        </section>
        <section className="card">
          <h3>Userspace</h3>
          ${a&&s`<p className="error">${a}</p>`}
          ${e?s`
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
              `:s`<p className="muted">Waiting for data…</p>`}
        </section>
        <section className="card">
          <div className="section-header">
            <div>
              <h3>Absolute & Rate of Change</h3>
              <p className="muted">Latest value and per-second delta.</p>
            </div>
          </div>
          <div className="summary-grid">
            ${W.map(n=>s`<${re} title=${n.title} path=${n.path} />`)}
          </div>
        </section>
      </main>
    `}function ne(){let[e,a]=f([]),[n,l]=f(""),[m,$]=f([]),[o,i]=f(""),[p,c]=f(null),[g,r]=f("");C(()=>{let h=!0;return(async()=>{try{let w=await b.get("/vips");if(!h)return;a(w||[]),!n&&w&&w.length>0&&l(q(w[0]))}catch(w){h&&r(w.message||"request failed")}})(),()=>{h=!1}},[]),C(()=>{if(!n)return;let h=A(n),P=!0;return(async()=>{try{let _=await b.get("/vips/reals",h);if(!P)return;$(_||[]),_&&_.length>0?i(x=>x||_[0].address):i(""),r("")}catch(_){P&&r(_.message||"request failed")}})(),()=>{P=!1}},[n]),C(()=>{if(!o){c(null);return}let h=!0;return(async()=>{try{let w=await b.get("/reals/index",{address:o});if(!h)return;c(w?.index??null),r("")}catch(w){h&&r(w.message||"request failed")}})(),()=>{h=!1}},[o]);let{points:u,error:v}=V({path:"/stats/real",body:p!==null?{index:p}:null}),y=S(()=>[{label:"v1",field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.2)"},{label:"v2",field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[]);return s`
      <main>
        <section className="card">
          <h2>Per-Real Stats</h2>
          <p className="muted">Select a VIP and real address to chart.</p>
          ${g&&s`<p className="error">${g}</p>`}
          <div className="form-row">
            <label className="field">
              <span>VIP</span>
              <select value=${n} onChange=${h=>l(h.target.value)}>
                ${e.map(h=>s`
                    <option value=${q(h)}>
                      ${h.address}:${h.port} / ${h.proto}
                    </option>
                  `)}
              </select>
            </label>
            <label className="field">
              <span>Real</span>
              <select
                value=${o}
                onChange=${h=>i(h.target.value)}
                disabled=${m.length===0}
              >
                ${m.map(h=>s`
                    <option value=${h.address}>${h.address}</option>
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
          ${v&&s`<p className="error">${v}</p>`}
          ${p===null?s`<p className="muted">Select a real to start polling.</p>`:s`<${z} points=${u} keys=${y} />`}
        </section>
      </main>
    `}function oe(){let{addToast:e}=L(),[a,n]=f(""),[l,m]=f(""),[$,o]=f(!0),[i,p]=f(""),c=T(!0),g=async()=>{if(c.current){o(!0),m("");try{let u=await fetch(`${b.base}/config/export`,{headers:{Accept:"application/x-yaml"}});if(!u.ok){let y=`HTTP ${u.status}`;try{y=(await u.json())?.error?.message||y}catch{}throw new Error(y)}let v=await u.text();if(!c.current)return;n(v||""),p(new Date().toLocaleString())}catch(u){c.current&&m(u.message||"request failed")}finally{c.current&&o(!1)}}},r=async()=>{if(a)try{await navigator.clipboard.writeText(a),e("Config copied to clipboard","info")}catch{e("Failed to copy config","error")}};return C(()=>(c.current=!0,g(),()=>{c.current=!1}),[]),s`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>Running config</h2>
              <p className="muted">Exported from /api/v1/config/export</p>
            </div>
            <div className="row">
              <button className="btn ghost" onClick=${r} disabled=${!a}>
                Copy YAML
              </button>
              <button className="btn" onClick=${g} disabled=${$}>
                Refresh
              </button>
            </div>
          </div>
          ${l&&s`<p className="error">${l}</p>`}
          ${$?s`<p className="muted">Loading config...</p>`:a?s`<pre className="yaml-view">${a}</pre>`:s`<p className="muted">No config data returned.</p>`}
          ${i&&s`<p className="muted">Last fetched ${i}</p>`}
        </section>
      </main>
    `}function ie(){let[e,a]=f({initialized:!1,ready:!1}),[n,l]=f([]),m=T({}),$=(i,p="info")=>{let c=`${Date.now()}-${Math.random().toString(16).slice(2)}`;l(g=>g.concat({id:c,message:i,kind:p})),m.current[c]=setTimeout(()=>{l(g=>g.filter(r=>r.id!==c)),delete m.current[c]},4e3)},o=i=>{m.current[i]&&(clearTimeout(m.current[i]),delete m.current[i]),l(p=>p.filter(c=>c.id!==i))};return C(()=>{let i=!0,p=async()=>{try{let g=await b.get("/lb/status");i&&a(g||{initialized:!1,ready:!1})}catch{i&&a({initialized:!1,ready:!1})}};p();let c=setInterval(p,5e3);return()=>{i=!1,clearInterval(c)}},[]),s`
      <${U}>
        <${G}>
          <${D.Provider} value=${{addToast:$}}>
            <${Q} status=${e} />
            <${M}>
              <${k} path="/" element=${s`<${ee} />`} />
              <${k} path="/vips/:vipId" element=${s`<${ae} />`} />
              <${k} path="/vips/:vipId/stats" element=${s`<${te} />`} />
              <${k} path="/stats/global" element=${s`<${le} />`} />
              <${k} path="/stats/real" element=${s`<${ne} />`} />
              <${k} path="/config" element=${s`<${oe} />`} />
            </${M}>
            <${Y} toasts=${n} onDismiss=${o} />
          </${D.Provider}>
        </${G}>
      </${U}>
    `}ReactDOM.createRoot(document.getElementById("root")).render(s`<${ie} />`)})();})();
