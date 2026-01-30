(()=>{(()=>{let{useEffect:T,useMemo:B,useRef:j,useState:f,useContext:Ne}=React,{BrowserRouter:ae,Routes:te,Route:D,NavLink:P,Link:M,useParams:se,useNavigate:ye}=ReactRouterDOM,s=htm.bind(React.createElement),J=React.createContext({addToast:()=>{}});function U(){return Ne(J)}let w={base:"/api/v1",async request(e,a={}){let n={method:a.method||"GET",headers:{"Content-Type":"application/json"}},r=`${w.base}${e}`;if(a.body!==void 0&&a.body!==null)if(n.method==="GET"){let i=new URLSearchParams;Object.entries(a.body).forEach(([g,d])=>{if(d!=null){if(Array.isArray(d)){d.forEach(b=>i.append(g,String(b)));return}if(typeof d=="object"){i.set(g,JSON.stringify(d));return}i.set(g,String(d))}});let t=i.toString();t&&(r+=`${r.includes("?")?"&":"?"}${t}`)}else n.body=JSON.stringify(a.body);let p=await fetch(r,n),h;try{h=await p.json()}catch{throw new Error("invalid JSON response")}if(!p.ok)throw new Error(h?.error?.message||`HTTP ${p.status}`);if(!h.success){let i=h.error?.message||"request failed";throw new Error(i)}return h.data},get(e,a){return w.request(e,{method:"GET",body:a})},post(e,a){return w.request(e,{method:"POST",body:a})},put(e,a){return w.request(e,{method:"PUT",body:a})},del(e,a){return w.request(e,{method:"DELETE",body:a})}},re="vatran_target_groups";function oe(e){if(!e||!e.address)return null;let a=String(e.address).trim();if(!a)return null;let n=Number(e.weight),r=Number(e.flags??0);return{address:a,weight:Number.isFinite(n)?n:0,flags:Number.isFinite(r)?r:0}}function ne(e){if(!e||typeof e!="object")return{};let a={};return Object.entries(e).forEach(([n,r])=>{let p=String(n).trim();if(!p)return;let h=Array.isArray(r)?r.map(oe).filter(Boolean):[],i=[],t=new Set;h.forEach(g=>{t.has(g.address)||(t.add(g.address),i.push(g))}),a[p]=i}),a}function K(){if(typeof localStorage>"u")return{};try{let e=localStorage.getItem(re);return e?ne(JSON.parse(e)):{}}catch{return{}}}function le(e){if(!(typeof localStorage>"u"))try{localStorage.setItem(re,JSON.stringify(e))}catch{}}function we(e,a){let n={...e};return Object.entries(a||{}).forEach(([r,p])=>{n[r]||(n[r]=p)}),n}function ie(){let[e,a]=f(()=>K());return T(()=>{le(e)},[e]),{groups:e,setGroups:a,refreshFromStorage:()=>{a(K())},importFromRunningConfig:async()=>{let p=await w.get("/config/export/json"),h=ne(p?.target_groups||{}),i=we(K(),h);return a(i),le(i),i}}}function W(e){return`${encodeURIComponent(e.address)}:${e.port}:${e.proto}`}function X(e){let a=e.split(":"),n=Number(a.pop()||0),r=Number(a.pop()||0);return{address:decodeURIComponent(a.join(":")),port:r,proto:n}}function _e(e,a,n=[]){let[r,p]=f(null),[h,i]=f(""),[t,g]=f(!0);return T(()=>{let d=!0,b=async()=>{try{let u=await e();d&&(p(u),i(""),g(!1))}catch(u){d&&(i(u.message||"request failed"),g(!1))}};b();let c=setInterval(b,a);return()=>{d=!1,clearInterval(c)}},n),{data:r,error:h,loading:t}}function H({path:e,body:a,intervalMs:n=1e3,limit:r=60}){let[p,h]=f([]),[i,t]=f(""),g=B(()=>JSON.stringify(a||{}),[a]);return T(()=>{if(a===null)return h([]),t(""),()=>{};let d=!0,b=async()=>{try{let u=await w.get(e,a);if(!d)return;let $=new Date().toLocaleTimeString();h(I=>I.concat({label:$,...u}).slice(-r)),t("")}catch(u){d&&t(u.message||"request failed")}};b();let c=setInterval(b,n);return()=>{d=!1,clearInterval(c)}},[e,g,n,r]),{points:p,error:i}}function ce({title:e,points:a,keys:n,diff:r=!1,height:p=120,showTitle:h=!1}){let i=j(null),t=j(null);return T(()=>{if(!i.current)return;t.current||(t.current=new Chart(i.current,{type:"line",data:{labels:[],datasets:[]},options:{responsive:!0,maintainAspectRatio:!1,animation:!1,scales:{x:{grid:{display:!1}},y:{beginAtZero:!r}},plugins:{legend:{display:!0,position:"bottom"},title:{display:h&&!!e,text:e}}}}));let g=t.current,d=a.map(b=>b.label);return g.data.labels=d,g.data.datasets=n.map(b=>{let c=a.map($=>$[b.field]||0),u=r?c.map(($,I)=>I===0?0:$-c[I-1]):c;return{label:b.label,data:u,borderColor:b.color,backgroundColor:b.fill,borderWidth:2,tension:.3}}),g.options.scales.y.beginAtZero=!r,g.options.plugins.title.display=h&&!!e,g.options.plugins.title.text=e||"",g.update(),()=>{}},[a,n,e,r,h]),T(()=>()=>{t.current&&(t.current.destroy(),t.current=null)},[]),s`<canvas ref=${i} height=${p}></canvas>`}function Z({title:e,points:a,keys:n,diff:r=!1,inlineTitle:p=!0}){let[h,i]=f(!1);return s`
      <div className="chart-wrap">
        <button className="chart-zoom-button" type="button" onClick=${()=>i(!0)}>
          <span className="zoom-icon" aria-hidden="true">+</span>
          Zoom
        </button>
        <div className="chart-click" onClick=${()=>i(!0)}>
          <${ce}
            title=${e}
            points=${a}
            keys=${n}
            diff=${r}
            height=${120}
            showTitle=${p&&!!e}
          />
        </div>
        ${h&&s`
          <div className="chart-overlay" onClick=${()=>i(!1)}>
            <div className="chart-modal" onClick=${t=>t.stopPropagation()}>
              <div className="row chart-modal-header">
                <div>
                  <h3>${e||"Chart"}</h3>
                  ${r?s`<p className="muted">Per-second delta.</p>`:""}
                </div>
                <button className="btn ghost" onClick=${()=>i(!1)}>
                  Close
                </button>
              </div>
              <div className="chart-zoom">
                <${ce}
                  title=${e}
                  points=${a}
                  keys=${n}
                  diff=${r}
                  height=${360}
                  showTitle=${!1}
                />
              </div>
            </div>
          </div>
        `}
      </div>
    `}function de({children:e}){return e}function Ie({toasts:e,onDismiss:a}){return s`
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
    `}function Se({status:e}){return s`
      <header>
        <div>
          <div style=${{fontSize:20,fontWeight:700}}>Vatran</div>
          <div className="status-pill">
            <span className=${`dot ${e.ready?"ok":""}`}></span>
            ${e.ready?"Ready":"Not ready"}
          </div>
        </div>
        <nav>
          <${P} to="/" end className=${({isActive:a})=>a?"active":""}>
            Dashboard
          </${P}>
          <${P} to="/stats/global" className=${({isActive:a})=>a?"active":""}>
            Global stats
          </${P}>
          <${P} to="/stats/real" className=${({isActive:a})=>a?"active":""}>
            Per-real stats
          </${P}>
          <${P}
            to="/target-groups"
            className=${({isActive:a})=>a?"active":""}
          >
            Target groups
          </${P}>
          <${P} to="/config" className=${({isActive:a})=>a?"active":""}>
            Config export
          </${P}>
        </nav>
      </header>
    `}function Ce(){let{addToast:e}=U(),[a,n]=f({initialized:!1,ready:!1}),[r,p]=f([]),[h,i]=f(""),[t,g]=f(!1),[d,b]=f(!1),[c,u]=f({main_interface:"",balancer_prog_path:"",healthchecking_prog_path:"",default_mac:"",local_mac:"",root_map_path:"",root_map_pos:2,katran_src_v4:"",katran_src_v6:"",use_root_map:!1,max_vips:512,max_reals:4096,hash_func:0}),[$,I]=f({address:"",port:80,proto:6,flags:0}),v=async()=>{try{let o=await w.get("/lb/status"),m=await w.get("/vips");n(o||{initialized:!1,ready:!1}),p(m||[]),i("")}catch(o){i(o.message||"request failed")}};T(()=>{let o=!0;return(async()=>{o&&await v()})(),()=>{o=!1}},[]);let k=async o=>{o.preventDefault();try{let m={...c,root_map_pos:c.root_map_pos===""?void 0:Number(c.root_map_pos),max_vips:Number(c.max_vips),max_reals:Number(c.max_reals),hash_func:Number(c.hash_func)};await w.post("/lb/create",m),i(""),g(!1),e("Load balancer initialized.","success"),await v()}catch(m){i(m.message||"request failed"),e(m.message||"Initialize failed.","error")}},S=async o=>{o.preventDefault();try{await w.post("/vips",{...$,port:Number($.port),proto:Number($.proto),flags:Number($.flags||0)}),I({address:"",port:80,proto:6,flags:0}),i(""),b(!1),e("VIP created.","success"),await v()}catch(m){i(m.message||"request failed"),e(m.message||"VIP create failed.","error")}},C=async()=>{try{await w.post("/lb/load-bpf-progs"),i(""),e("BPF programs loaded.","success"),await v()}catch(o){i(o.message||"request failed"),e(o.message||"Load BPF programs failed.","error")}},F=async()=>{try{await w.post("/lb/attach-bpf-progs"),i(""),e("BPF programs attached.","success"),await v()}catch(o){i(o.message||"request failed"),e(o.message||"Attach BPF programs failed.","error")}};return s`
      <main>
        <section className="card">
          <h2>Load balancer</h2>
          <p>Initialized: ${a.initialized?"yes":"no"}</p>
          <p>Ready: ${a.ready?"yes":"no"}</p>
          <div className="row">
            ${!a.initialized&&s`
              <button className="btn" onClick=${()=>g(o=>!o)}>
                ${t?"Close":"Initialize"}
              </button>
            `}
            <button className="btn secondary" onClick=${()=>b(o=>!o)}>
              ${d?"Close":"Create VIP"}
            </button>
          </div>
          ${!a.ready&&s`
            <div className="row" style=${{marginTop:12}}>
              <button
                className="btn ghost"
                disabled=${!a.initialized}
                onClick=${C}
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
          ${t&&s`
            <form className="form" onSubmit=${k}>
              <div className="form-row">
                <label className="field">
                  <span>Main interface</span>
                  <input
                    value=${c.main_interface}
                    onInput=${o=>u({...c,main_interface:o.target.value})}
                    placeholder="eth0"
                    required
                  />
                </label>
                <label className="field">
                  <span>Balancer prog path</span>
                  <input
                    value=${c.balancer_prog_path}
                    onInput=${o=>u({...c,balancer_prog_path:o.target.value})}
                    placeholder="/path/to/bpf.o"
                    required
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Healthchecking path</span>
                  <input
                    value=${c.healthchecking_prog_path}
                    onInput=${o=>u({...c,healthchecking_prog_path:o.target.value})}
                    placeholder="/path/to/hc.o"
                  />
                </label>
                <label className="field">
                  <span>Default MAC</span>
                  <input
                    value=${c.default_mac}
                    onInput=${o=>u({...c,default_mac:o.target.value})}
                    placeholder="aa:bb:cc:dd:ee:ff"
                    required
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Local MAC</span>
                  <input
                    value=${c.local_mac}
                    onInput=${o=>u({...c,local_mac:o.target.value})}
                    placeholder="11:22:33:44:55:66"
                    required=${c.healthchecking_prog_path?.trim()!==""}
                  />
                </label>
                <label className="field">
                  <span>Hash function</span>
                  <input
                    type="number"
                    min="0"
                    value=${c.hash_func}
                    onInput=${o=>u({...c,hash_func:o.target.value})}
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Root map path</span>
                  <input
                    value=${c.root_map_path}
                    onInput=${o=>u({...c,root_map_path:o.target.value})}
                    placeholder="/sys/fs/bpf/root_map"
                  />
                </label>
                <label className="field">
                  <span>Root map position</span>
                  <input
                    type="number"
                    min="0"
                    value=${c.root_map_pos}
                    onInput=${o=>u({...c,root_map_pos:o.target.value})}
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Katran src v4</span>
                  <input
                    value=${c.katran_src_v4}
                    onInput=${o=>u({...c,katran_src_v4:o.target.value})}
                    placeholder="10.0.0.1"
                  />
                </label>
                <label className="field">
                  <span>Katran src v6</span>
                  <input
                    value=${c.katran_src_v6}
                    onInput=${o=>u({...c,katran_src_v6:o.target.value})}
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
                    value=${c.max_vips}
                    onInput=${o=>u({...c,max_vips:o.target.value})}
                  />
                </label>
                <label className="field">
                  <span>Max reals</span>
                  <input
                    type="number"
                    min="1"
                    value=${c.max_reals}
                    onInput=${o=>u({...c,max_reals:o.target.value})}
                  />
                </label>
              </div>
              <label className="field checkbox">
                <input
                  type="checkbox"
                  checked=${c.use_root_map}
                  onChange=${o=>u({...c,use_root_map:o.target.checked})}
                />
                <span>Use root map</span>
              </label>
              <button className="btn" type="submit">Initialize LB</button>
            </form>
          `}
          ${d&&s`
            <form className="form" onSubmit=${S}>
              <div className="form-row">
                <label className="field">
                  <span>VIP address</span>
                  <input
                    value=${$.address}
                    onInput=${o=>I({...$,address:o.target.value})}
                    placeholder="1.2.3.4"
                    required
                  />
                </label>
                <label className="field">
                  <span>Port</span>
                  <input
                    type="number"
                    value=${$.port}
                    onInput=${o=>I({...$,port:o.target.value})}
                    required
                  />
                </label>
                <label className="field">
                  <span>Protocol</span>
                  <select
                    value=${$.proto}
                    onChange=${o=>I({...$,proto:o.target.value})}
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
                    onInput=${o=>I({...$,flags:o.target.value})}
                  />
                </label>
              </div>
              <button className="btn" type="submit">Create VIP</button>
            </form>
          `}
          ${h&&s`<p className="error">${h}</p>`}
        </section>
        <section className="card">
          <div className="section-header">
            <h2>VIPs</h2>
            <button className="btn ghost" onClick=${v}>Refresh</button>
          </div>
          ${r.length===0?s`<p className="muted">No VIPs configured yet.</p>`:s`
                <div className="grid">
                  ${r.map(o=>s`
                      <div className="card">
                        <div style=${{fontWeight:600}}>
                          ${o.address}:${o.port} / ${o.proto}
                        </div>
                        <div className="muted" style=${{marginTop:6}}>
                          Flags: ${o.flags||0}
                        </div>
                        <div className="row" style=${{marginTop:12}}>
                          <${M} className="btn" to=${`/vips/${W(o)}`}>
                            Open
                          </${M}>
                          <${M}
                            className="btn secondary"
                            to=${`/vips/${W(o)}/stats`}
                          >
                            Stats
                          </${M}>
                        </div>
                      </div>
                    `)}
                </div>
              `}
        </section>
      </main>
    `}function Re(){let{addToast:e}=U(),a=se(),n=ye(),r=B(()=>X(a.vipId),[a.vipId]),[p,h]=f([]),[i,t]=f(""),[g,d]=f(""),[b,c]=f(!0),[u,$]=f({address:"",weight:100,flags:0}),[I,v]=f({}),[k,S]=f(null),[C,F]=f({flag:0,set:!0}),[o,m]=f({hash_function:0}),{groups:y,setGroups:A,refreshFromStorage:V,importFromRunningConfig:x}=ie(),[G,me]=f(""),[ge,L]=f(""),[fe,he]=f(!1),[be,ve]=f(""),[Y,$e]=f({add:0,update:0,remove:0}),O=async()=>{try{let l=await w.get("/vips/reals",r);h(l||[]);let N={};(l||[]).forEach(q=>{N[q.address]=q.weight}),v(N),t(""),c(!1)}catch(l){t(l.message||"request failed"),c(!1)}},Q=async()=>{try{let l=await w.get("/vips/flags",r);S(l?.flags??0),d("")}catch(l){d(l.message||"request failed")}};T(()=>{O(),Q()},[a.vipId]),T(()=>{if(!G){$e({add:0,update:0,remove:0});return}let l=y[G]||[],N=new Map(p.map(R=>[R.address,R])),q=new Map(l.map(R=>[R.address,R])),E=0,z=0,_=0;l.forEach(R=>{let ee=N.get(R.address);if(!ee){E+=1;return}(Number(ee.weight)!==Number(R.weight)||Number(ee.flags||0)!==Number(R.flags||0))&&(z+=1)}),p.forEach(R=>{q.has(R.address)||(_+=1)}),$e({add:E,update:z,remove:_})},[G,p,y]);let qe=async l=>{try{let N=Number(I[l.address]);await w.post("/vips/reals",{vip:r,real:{address:l.address,weight:N,flags:l.flags||0}}),await O(),e("Real weight updated.","success")}catch(N){t(N.message||"request failed"),e(N.message||"Update failed.","error")}},De=async l=>{try{await w.del("/vips/reals",{vip:r,real:{address:l.address,weight:l.weight,flags:l.flags||0}}),await O(),e("Real removed.","success")}catch(N){t(N.message||"request failed"),e(N.message||"Remove failed.","error")}},Ve=async l=>{l.preventDefault();try{await w.post("/vips/reals",{vip:r,real:{address:u.address,weight:Number(u.weight),flags:Number(u.flags||0)}}),$({address:"",weight:100,flags:0}),await O(),e("Real added.","success")}catch(N){t(N.message||"request failed"),e(N.message||"Add failed.","error")}},Le=async()=>{if(!G||!y[G]){L("Select a target group to apply.");return}he(!0),L("");let l=y[G]||[],N=new Map(p.map(_=>[_.address,_])),q=new Map(l.map(_=>[_.address,_])),E=p.filter(_=>!q.has(_.address)),z=l.filter(_=>{let R=N.get(_.address);return R?Number(R.weight)!==Number(_.weight)||Number(R.flags||0)!==Number(_.flags||0):!0});try{E.length>0&&await w.put("/vips/reals/batch",{vip:r,action:1,reals:E.map(_=>({address:_.address,weight:Number(_.weight),flags:Number(_.flags||0)}))}),z.length>0&&await Promise.all(z.map(_=>w.post("/vips/reals",{vip:r,real:{address:_.address,weight:Number(_.weight),flags:Number(_.flags||0)}}))),await O(),e(`Applied target group "${G}".`,"success")}catch(_){L(_.message||"Failed to apply target group."),e(_.message||"Target group apply failed.","error")}finally{he(!1)}},Be=l=>{l.preventDefault();let N=be.trim();if(!N){L("Provide a name for the new target group.");return}if(y[N]){L("A target group with that name already exists.");return}let q={...y,[N]:p.map(E=>({address:E.address,weight:Number(E.weight),flags:Number(E.flags||0)}))};A(q),ve(""),me(N),L(""),e(`Target group "${N}" saved.`,"success")},Oe=async()=>{try{await w.del("/vips",r),e("VIP deleted.","success"),n("/")}catch(l){t(l.message||"request failed"),e(l.message||"Delete failed.","error")}},ze=async l=>{l.preventDefault();try{await w.put("/vips/flags",{...r,flag:Number(C.flag),set:!!C.set}),await Q(),e("VIP flags updated.","success")}catch(N){d(N.message||"request failed"),e(N.message||"Flag update failed.","error")}},je=async l=>{l.preventDefault();try{await w.put("/vips/hash-function",{...r,hash_function:Number(o.hash_function)}),e("Hash function updated.","success")}catch(N){d(N.message||"request failed"),e(N.message||"Hash update failed.","error")}};return s`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>VIP Detail</h2>
              <p className="muted">${r.address}:${r.port} / ${r.proto}</p>
              <p className="muted">Flags: ${k??"\u2014"}</p>
            </div>
            <div className="row">
              <button className="btn ghost" onClick=${Q}>Refresh flags</button>
              <button className="btn danger" onClick=${Oe}>Delete VIP</button>
            </div>
          </div>
          ${i&&s`<p className="error">${i}</p>`}
          ${g&&s`<p className="error">${g}</p>`}
          ${b?s`<p className="muted">Loading reals…</p>`:s`
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
                    ${p.map(l=>s`
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
                              value=${I[l.address]??l.weight}
                              onInput=${N=>v({...I,[l.address]:N.target.value})}
                            />
                          </td>
                          <td>${l.flags||0}</td>
                          <td className="row">
                            <button className="btn" onClick=${()=>qe(l)}>
                              Update
                            </button>
                            <button className="btn ghost" onClick=${()=>De(l)}>
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
            <form className="form" onSubmit=${ze}>
              <div className="form-row">
                <label className="field">
                  <span>Flag</span>
                  <input
                    type="number"
                    value=${C.flag}
                    onInput=${l=>F({...C,flag:l.target.value})}
                  />
                </label>
                <label className="field">
                  <span>Set</span>
                  <select
                    value=${String(C.set)}
                    onChange=${l=>F({...C,set:l.target.value==="true"})}
                  >
                    <option value="true">Set</option>
                    <option value="false">Clear</option>
                  </select>
                </label>
              </div>
              <button className="btn" type="submit">Apply flag</button>
            </form>
            <form className="form" onSubmit=${je}>
              <div className="form-row">
                <label className="field">
                  <span>Hash function</span>
                  <input
                    type="number"
                    min="0"
                    value=${o.hash_function}
                    onInput=${l=>m({...o,hash_function:l.target.value})}
                  />
                </label>
              </div>
              <button className="btn" type="submit">Apply hash</button>
            </form>
          </div>
        </section>
        <section className="card">
          <h3>Add real</h3>
          <form className="form" onSubmit=${Ve}>
            <div className="form-row">
              <label className="field">
                <span>Real address</span>
                <input
                  value=${u.address}
                  onInput=${l=>$({...u,address:l.target.value})}
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
                  onInput=${l=>$({...u,weight:l.target.value})}
                  required
                />
              </label>
              <label className="field">
                <span>Flags</span>
                <input
                  type="number"
                  value=${u.flags}
                  onInput=${l=>$({...u,flags:l.target.value})}
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
              <button className="btn ghost" type="button" onClick=${V}>
                Reload groups
              </button>
              <button
                className="btn ghost"
                type="button"
                onClick=${async()=>{try{await x(),e("Imported target groups from running config.","success")}catch(l){L(l.message||"Failed to import target groups."),e(l.message||"Import failed.","error")}}}
              >
                Import from running config
              </button>
            </div>
          </div>
          ${ge&&s`<p className="error">${ge}</p>`}
          <div className="form-row">
            <label className="field">
              <span>Target group</span>
              <select
                value=${G}
                onChange=${l=>me(l.target.value)}
                disabled=${Object.keys(y).length===0}
              >
                <option value="">Select group</option>
                ${Object.keys(y).map(l=>s`<option value=${l}>${l}</option>`)}
              </select>
            </label>
            <label className="field">
              <span>Preview</span>
              <input
                value=${`add ${Y.add} \xB7 update ${Y.update} \xB7 remove ${Y.remove}`}
                readOnly
              />
            </label>
          </div>
          <div className="row">
            <button
              className="btn"
              type="button"
              onClick=${Le}
              disabled=${fe||!G}
            >
              ${fe?"Applying...":"Apply target group"}
            </button>
          </div>
          <form className="form" onSubmit=${Be}>
            <div className="form-row">
              <label className="field">
                <span>Save current reals as new group</span>
                <input
                  value=${be}
                  onInput=${l=>ve(l.target.value)}
                  placeholder="edge-backends"
                />
              </label>
            </div>
            <button className="btn secondary" type="submit">Save target group</button>
          </form>
        </section>
      </main>
    `}function Te(){let e=se(),a=B(()=>X(e.vipId),[e.vipId]),{points:n,error:r}=H({path:"/stats/vip",body:a}),p=B(()=>[{label:"v1",field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.15)"},{label:"v2",field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[]);return s`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>VIP Stats</h2>
              <p className="muted">${a.address}:${a.port} / ${a.proto}</p>
            </div>
          </div>
          ${r&&s`<p className="error">${r}</p>`}
          <${Z} title="Traffic" points=${n} keys=${p} />
        </section>
      </main>
    `}let ue=[{title:"LRU",path:"/stats/lru"},{title:"LRU Miss",path:"/stats/lru/miss"},{title:"LRU Fallback",path:"/stats/lru/fallback"},{title:"LRU Global",path:"/stats/lru/global"},{title:"XDP Total",path:"/stats/xdp/total"},{title:"XDP Pass",path:"/stats/xdp/pass"},{title:"XDP Drop",path:"/stats/xdp/drop"},{title:"XDP Tx",path:"/stats/xdp/tx"}];function pe(e){return Number.isFinite(e)?`${e>0?"+":""}${e}`:"0"}function ke({title:e,path:a,diff:n=!1}){let{points:r,error:p}=H({path:a}),h=B(()=>[{label:"v1",field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.2)"},{label:"v2",field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[]);return s`
      <div className="card">
        <h3>${e}</h3>
        ${p&&s`<p className="error">${p}</p>`}
        <${Z} title=${e} points=${r} keys=${h} diff=${n} inlineTitle=${!1} />
      </div>
    `}function Pe({title:e,path:a}){let{points:n,error:r}=H({path:a}),p=n[n.length-1]||{},h=n[n.length-2]||{},i=Number(p.v1??0),t=Number(p.v2??0),g=i-Number(h.v1??0),d=t-Number(h.v2??0);return s`
      <div className="summary-card">
        <div className="summary-title">${e}</div>
        ${r?s`<p className="error">${r}</p>`:s`
              <div className="summary-row">
                <div className="stat">
                  <span className="muted">v1 absolute</span>
                  <strong>${i}</strong>
                </div>
                <div className="stat">
                  <span className="muted">v1 delta/sec</span>
                  <strong className=${g<0?"delta down":"delta up"}>
                    ${pe(g)}
                  </strong>
                </div>
              </div>
              <div className="summary-row">
                <div className="stat">
                  <span className="muted">v2 absolute</span>
                  <strong>${t}</strong>
                </div>
                <div className="stat">
                  <span className="muted">v2 delta/sec</span>
                  <strong className=${d<0?"delta down":"delta up"}>
                    ${pe(d)}
                  </strong>
                </div>
              </div>
            `}
      </div>
    `}function xe(){let{data:e,error:a}=_e(()=>w.get("/stats/userspace"),1e3,[]);return s`
      <main>
        <section className="card">
          <h2>Global Stats</h2>
          <p className="muted">Polling every second. Charts show per-second deltas.</p>
        </section>
        <section className="grid">
          ${ue.map(n=>s`<${ke} title=${n.title} path=${n.path} diff=${!0} />`)}
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
            ${ue.map(n=>s`<${Pe} title=${n.title} path=${n.path} />`)}
          </div>
        </section>
      </main>
    `}function Fe(){let[e,a]=f([]),[n,r]=f(""),[p,h]=f([]),[i,t]=f(""),[g,d]=f(null),[b,c]=f("");T(()=>{let v=!0;return(async()=>{try{let S=await w.get("/vips");if(!v)return;a(S||[]),!n&&S&&S.length>0&&r(W(S[0]))}catch(S){v&&c(S.message||"request failed")}})(),()=>{v=!1}},[]),T(()=>{if(!n)return;let v=X(n),k=!0;return(async()=>{try{let C=await w.get("/vips/reals",v);if(!k)return;h(C||[]),C&&C.length>0?t(F=>F||C[0].address):t(""),c("")}catch(C){k&&c(C.message||"request failed")}})(),()=>{k=!1}},[n]),T(()=>{if(!i){d(null);return}let v=!0;return(async()=>{try{let S=await w.get("/reals/index",{address:i});if(!v)return;d(S?.index??null),c("")}catch(S){v&&c(S.message||"request failed")}})(),()=>{v=!1}},[i]);let{points:u,error:$}=H({path:"/stats/real",body:g!==null?{index:g}:null}),I=B(()=>[{label:"v1",field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.2)"},{label:"v2",field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[]);return s`
      <main>
        <section className="card">
          <h2>Per-Real Stats</h2>
          <p className="muted">Select a VIP and real address to chart.</p>
          ${b&&s`<p className="error">${b}</p>`}
          <div className="form-row">
            <label className="field">
              <span>VIP</span>
              <select value=${n} onChange=${v=>r(v.target.value)}>
                ${e.map(v=>s`
                    <option value=${W(v)}>
                      ${v.address}:${v.port} / ${v.proto}
                    </option>
                  `)}
              </select>
            </label>
            <label className="field">
              <span>Real</span>
              <select
                value=${i}
                onChange=${v=>t(v.target.value)}
                disabled=${p.length===0}
              >
                ${p.map(v=>s`
                    <option value=${v.address}>${v.address}</option>
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
          ${$&&s`<p className="error">${$}</p>`}
          ${g===null?s`<p className="muted">Select a real to start polling.</p>`:s`<${Z} points=${u} keys=${I} />`}
        </section>
      </main>
    `}function Ge(){let{addToast:e}=U(),[a,n]=f(""),[r,p]=f(""),[h,i]=f(!0),[t,g]=f(""),d=j(!0),b=async()=>{if(d.current){i(!0),p("");try{let u=await fetch(`${w.base}/config/export`,{headers:{Accept:"application/x-yaml"}});if(!u.ok){let I=`HTTP ${u.status}`;try{I=(await u.json())?.error?.message||I}catch{}throw new Error(I)}let $=await u.text();if(!d.current)return;n($||""),g(new Date().toLocaleString())}catch(u){d.current&&p(u.message||"request failed")}finally{d.current&&i(!1)}}},c=async()=>{if(a)try{await navigator.clipboard.writeText(a),e("Config copied to clipboard","info")}catch{e("Failed to copy config","error")}};return T(()=>(d.current=!0,b(),()=>{d.current=!1}),[]),s`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>Running config</h2>
              <p className="muted">Exported from /api/v1/config/export</p>
            </div>
            <div className="row">
              <button className="btn ghost" onClick=${c} disabled=${!a}>
                Copy YAML
              </button>
              <button className="btn" onClick=${b} disabled=${h}>
                Refresh
              </button>
            </div>
          </div>
          ${r&&s`<p className="error">${r}</p>`}
          ${h?s`<p className="muted">Loading config...</p>`:a?s`<pre className="yaml-view">${a}</pre>`:s`<p className="muted">No config data returned.</p>`}
          ${t&&s`<p className="muted">Last fetched ${t}</p>`}
        </section>
      </main>
    `}function Ee(){let{addToast:e}=U(),{groups:a,setGroups:n,refreshFromStorage:r,importFromRunningConfig:p}=ie(),[h,i]=f(""),[t,g]=f(""),[d,b]=f({address:"",weight:100,flags:0}),[c,u]=f(""),[$,I]=f(!1);T(()=>{if(t){if(!a[t]){let m=Object.keys(a);g(m[0]||"")}}else{let m=Object.keys(a);m.length>0&&g(m[0])}},[a,t]);let v=m=>{m.preventDefault();let y=h.trim();if(!y){u("Provide a group name.");return}if(a[y]){u("That group already exists.");return}n({...a,[y]:[]}),i(""),g(y),u(""),e(`Target group "${y}" created.`,"success")},k=m=>{let y={...a};delete y[m],n(y),e(`Target group "${m}" removed.`,"success")},S=m=>{if(m.preventDefault(),!t){u("Select a group to add a real.");return}let y=oe(d);if(!y){u("Provide a valid real address.");return}let A=a[t]||[],V=A.some(x=>x.address===y.address)?A.map(x=>x.address===y.address?y:x):A.concat(y);n({...a,[t]:V}),b({address:"",weight:100,flags:0}),u(""),e("Real saved to target group.","success")},C=m=>{if(!t)return;let A=(a[t]||[]).filter(V=>V.address!==m);n({...a,[t]:A})},F=(m,y)=>{if(!t)return;let V=(a[t]||[]).map(x=>x.address===m?{...x,...y}:x);n({...a,[t]:V})};return s`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>Target groups</h2>
              <p className="muted">Define reusable sets of reals (address + weight).</p>
            </div>
            <div className="row">
              <button className="btn ghost" type="button" onClick=${r}>
                Reload groups
              </button>
              <button className="btn ghost" type="button" onClick=${async()=>{I(!0);try{await p(),e("Imported target groups from running config.","success"),u("")}catch(m){u(m.message||"Failed to import target groups."),e(m.message||"Import failed.","error")}finally{I(!1)}}} disabled=${$}>
                ${$?"Importing...":"Import from running config"}
              </button>
            </div>
          </div>
          ${c&&s`<p className="error">${c}</p>`}
          <form className="form" onSubmit=${v}>
            <div className="form-row">
              <label className="field">
                <span>New group name</span>
                <input
                  value=${h}
                  onInput=${m=>i(m.target.value)}
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
                  value=${t}
                  onChange=${m=>g(m.target.value)}
                  disabled=${Object.keys(a).length===0}
                >
                  ${Object.keys(a).map(m=>s`<option value=${m}>${m}</option>`)}
                </select>
              </label>
              ${t&&s`<button className="btn danger" type="button" onClick=${()=>k(t)}>
                Delete group
              </button>`}
            </div>
          </div>
          ${t?s`
                <table className="table">
                  <thead>
                    <tr>
                      <th>Address</th>
                      <th>Weight</th>
                      <th>Flags</th>
                      <th></th>
                    </tr>
                  </thead>
                  <tbody>
                    ${(a[t]||[]).map(m=>s`
                        <tr>
                          <td>${m.address}</td>
                          <td>
                            <input
                              className="inline-input"
                              type="number"
                              min="0"
                              value=${m.weight}
                              onInput=${y=>F(m.address,{weight:Number(y.target.value)})}
                            />
                          </td>
                          <td>
                            <input
                              className="inline-input"
                              type="number"
                              min="0"
                              value=${m.flags||0}
                              onInput=${y=>F(m.address,{flags:Number(y.target.value)})}
                            />
                          </td>
                          <td>
                            <button
                              className="btn ghost"
                              type="button"
                              onClick=${()=>C(m.address)}
                            >
                              Remove
                            </button>
                          </td>
                        </tr>
                      `)}
                  </tbody>
                </table>
                <form className="form" onSubmit=${S}>
                  <div className="form-row">
                    <label className="field">
                      <span>Real address</span>
                      <input
                        value=${d.address}
                        onInput=${m=>b({...d,address:m.target.value})}
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
                        onInput=${m=>b({...d,weight:m.target.value})}
                        required
                      />
                    </label>
                    <label className="field">
                      <span>Flags</span>
                      <input
                        type="number"
                        value=${d.flags}
                        onInput=${m=>b({...d,flags:m.target.value})}
                      />
                    </label>
                  </div>
                  <button className="btn secondary" type="submit">Save real</button>
                </form>
              `:s`<p className="muted">No groups yet. Create one to add reals.</p>`}
        </section>
      </main>
    `}function Ae(){let[e,a]=f({initialized:!1,ready:!1}),[n,r]=f([]),p=j({}),h=(t,g="info")=>{let d=`${Date.now()}-${Math.random().toString(16).slice(2)}`;r(b=>b.concat({id:d,message:t,kind:g})),p.current[d]=setTimeout(()=>{r(b=>b.filter(c=>c.id!==d)),delete p.current[d]},4e3)},i=t=>{p.current[t]&&(clearTimeout(p.current[t]),delete p.current[t]),r(g=>g.filter(d=>d.id!==t))};return T(()=>{let t=!0,g=async()=>{try{let b=await w.get("/lb/status");t&&a(b||{initialized:!1,ready:!1})}catch{t&&a({initialized:!1,ready:!1})}};g();let d=setInterval(g,5e3);return()=>{t=!1,clearInterval(d)}},[]),s`
      <${ae}>
        <${de}>
          <${J.Provider} value=${{addToast:h}}>
            <${Se} status=${e} />
            <${te}>
              <${D} path="/" element=${s`<${Ce} />`} />
              <${D} path="/vips/:vipId" element=${s`<${Re} />`} />
              <${D} path="/vips/:vipId/stats" element=${s`<${Te} />`} />
              <${D} path="/target-groups" element=${s`<${Ee} />`} />
              <${D} path="/stats/global" element=${s`<${xe} />`} />
              <${D} path="/stats/real" element=${s`<${Fe} />`} />
              <${D} path="/config" element=${s`<${Ge} />`} />
            </${te}>
            <${Ie} toasts=${n} onDismiss=${i} />
          </${J.Provider}>
        </${de}>
      </${ae}>
    `}ReactDOM.createRoot(document.getElementById("root")).render(s`<${Ae} />`)})();})();
