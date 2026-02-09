(()=>{(()=>{let{useEffect:G,useMemo:E,useRef:U,useState:h,useContext:Se}=React,{BrowserRouter:ne,Routes:oe,Route:q,NavLink:L,Link:j,useParams:le,useNavigate:Ce}=ReactRouterDOM,i=htm.bind(React.createElement),X=React.createContext({addToast:()=>{}}),H=[{label:"NO_SPORT",value:1},{label:"NO_LRU",value:2},{label:"QUIC_VIP",value:4},{label:"DPORT_HASH",value:8},{label:"SRC_ROUTING",value:16},{label:"LOCAL_VIP",value:32},{label:"GLOBAL_LRU",value:64},{label:"HASH_SRC_DST_PORT",value:128},{label:"UDP_STABLE_ROUTING_VIP",value:256},{label:"UDP_FLOW_MIGRATION",value:512}];function W(){return Se(X)}function Re(e){return String(e).replace(/[^a-z0-9_-]/gi,"_")}function Te(e,t,a){let o=Number(e)||0,c=Number(t)||0;return a?o|c:o&~c}function Pe(e,t){let a=Number(e)||0;return t.filter(o=>(a&o.value)!==0)}function ie(e,t){let a=String(e??"").trim();if(!a)return;let c=a.split(/[\s,]+/).filter(Boolean).map(n=>Number(n));if(c.findIndex(n=>!Number.isFinite(n)||!Number.isInteger(n))!==-1)throw new Error(`${t} must be a comma- or space-separated list of integers.`);return c}function ce({mask:e,options:t,showStatus:a=!1,emptyLabel:o="None"}){let c=Number(e)||0,v=a?t:Pe(c,t),n=a?2:1;return i`
      <table className="table flag-table">
        <thead>
          <tr>
            <th>Flag</th>
            ${a?i`<th>Enabled</th>`:null}
          </tr>
        </thead>
        <tbody>
          ${v.length===0?i`<tr><td colspan=${n} className="muted">${o}</td></tr>`:v.map(l=>{let p=(c&l.value)!==0;return i`
                  <tr>
                    <td>${l.label}</td>
                    ${a?i`<td>${p?"Yes":"No"}</td>`:null}
                  </tr>
                `})}
        </tbody>
      </table>
    `}function de({options:e,value:t,onChange:a,name:o}){let c=Number(t)||0,v=Re(o||"flags");return i`
      <div className="flag-selector">
        <table className="table flag-table">
          <thead>
            <tr>
              <th>Enabled</th>
              <th>Flag</th>
            </tr>
          </thead>
          <tbody>
            ${e.map(n=>{let l=`${v}-${n.value}`,p=(c&n.value)===n.value;return i`
                <tr>
                  <td>
                    <input
                      id=${l}
                      type="checkbox"
                      checked=${p}
                      onChange=${d=>a(Te(c,n.value,d.target.checked))}
                    />
                  </td>
                  <td>
                    <label className="flag-option" htmlFor=${l}>${n.label}</label>
                  </td>
                </tr>
              `})}
          </tbody>
        </table>
      </div>
    `}let $={base:"/api/v1",async request(e,t={}){let a={method:t.method||"GET",headers:{"Content-Type":"application/json"}},o=`${$.base}${e}`;if(t.body!==void 0&&t.body!==null)if(a.method==="GET"){let n=new URLSearchParams;Object.entries(t.body).forEach(([p,d])=>{if(d!=null){if(Array.isArray(d)){d.forEach(m=>n.append(p,String(m)));return}if(typeof d=="object"){n.set(p,JSON.stringify(d));return}n.set(p,String(d))}});let l=n.toString();l&&(o+=`${o.includes("?")?"&":"?"}${l}`)}else a.body=JSON.stringify(t.body);let c=await fetch(o,a),v;try{v=await c.json()}catch{throw new Error("invalid JSON response")}if(!c.ok)throw new Error(v?.error?.message||`HTTP ${c.status}`);if(!v.success){let n=v.error?.message||"request failed";throw new Error(n)}return v.data},get(e,t){return $.request(e,{method:"GET",body:t})},post(e,t){return $.request(e,{method:"POST",body:t})},put(e,t){return $.request(e,{method:"PUT",body:t})},del(e,t){return $.request(e,{method:"DELETE",body:t})}},ue="vatran_target_groups";function pe(e){if(!e||!e.address)return null;let t=String(e.address).trim();if(!t)return null;let a=Number(e.weight),o=Number(e.flags??0);return{address:t,weight:Number.isFinite(a)?a:0,flags:Number.isFinite(o)?o:0}}function me(e){if(!e||typeof e!="object")return{};let t={};return Object.entries(e).forEach(([a,o])=>{let c=String(a).trim();if(!c)return;let v=Array.isArray(o)?o.map(pe).filter(Boolean):[],n=[],l=new Set;v.forEach(p=>{l.has(p.address)||(l.add(p.address),n.push(p))}),t[c]=n}),t}function Q(){if(typeof localStorage>"u")return{};try{let e=localStorage.getItem(ue);return e?me(JSON.parse(e)):{}}catch{return{}}}function ge(e){if(!(typeof localStorage>"u"))try{localStorage.setItem(ue,JSON.stringify(e))}catch{}}function ke(e,t){let a={...e};return Object.entries(t||{}).forEach(([o,c])=>{a[o]||(a[o]=c)}),a}function fe(){let[e,t]=h(()=>Q());return G(()=>{ge(e)},[e]),{groups:e,setGroups:t,refreshFromStorage:()=>{t(Q())},importFromRunningConfig:async()=>{let c=await $.get("/config/export/json"),v=me(c?.target_groups||{}),n=ke(Q(),v);return t(n),ge(n),n}}}function J(e){return`${encodeURIComponent(e.address)}:${e.port}:${e.proto}`}function ee(e){let t=e.split(":"),a=Number(t.pop()||0),o=Number(t.pop()||0);return{address:decodeURIComponent(t.join(":")),port:o,proto:a}}function xe(e,t,a=[]){let[o,c]=h(null),[v,n]=h(""),[l,p]=h(!0);return G(()=>{let d=!0,m=async()=>{try{let g=await e();d&&(c(g),n(""),p(!1))}catch(g){d&&(n(g.message||"request failed"),p(!1))}};m();let s=setInterval(m,t);return()=>{d=!1,clearInterval(s)}},a),{data:o,error:v,loading:l}}function K({path:e,body:t,intervalMs:a=1e3,limit:o=60}){let[c,v]=h([]),[n,l]=h(""),p=E(()=>JSON.stringify(t||{}),[t]);return G(()=>{if(t===null)return v([]),l(""),()=>{};let d=!0,m=async()=>{try{let g=await $.get(e,t);if(!d)return;let y=new Date().toLocaleTimeString();v(w=>w.concat({label:y,...g}).slice(-o)),l("")}catch(g){d&&l(g.message||"request failed")}};m();let s=setInterval(m,a);return()=>{d=!1,clearInterval(s)}},[e,p,a,o]),{points:c,error:n}}function ve({title:e,points:t,keys:a,diff:o=!1,height:c=120,showTitle:v=!1,selectedLabel:n=null,onPointSelect:l=null,onLegendSelect:p=null}){let d=U(null),m=U(null);return G(()=>{if(!d.current)return;m.current||(m.current=new Chart(d.current,{type:"line",data:{labels:[],datasets:[]},options:{responsive:!0,maintainAspectRatio:!1,animation:!1,scales:{x:{grid:{display:!1}},y:{beginAtZero:!o}},plugins:{legend:{display:!0,position:"bottom"},title:{display:v&&!!e,text:e}}}}));let s=m.current,g=new Map((s.data.datasets||[]).filter(_=>typeof _.hidden<"u").map(_=>[_.label,_.hidden])),y=t.map(_=>_.label),w=n?a.filter(_=>_.label===n):a,P=n&&w.length===0?a:w;return s.data.labels=y,s.data.datasets=P.map(_=>{let k=t.map(F=>F[_.field]||0),T=o?k.map((F,r)=>r===0?0:F-k[r-1]):k;return{label:_.label,data:T,borderColor:_.color,backgroundColor:_.fill,borderWidth:2,tension:.3,hidden:g.get(_.label)}}),s.options.onClick=(_,k)=>{if(!l||!k||k.length===0)return;let T=k[0].datasetIndex,F=s.data.datasets?.[T]?.label;F&&l(F)},s.options.plugins&&s.options.plugins.legend&&(s.options.plugins.legend.onClick=(_,k)=>{if(!p)return;let T=k?.text;T&&p(T)}),s.options.scales.y.beginAtZero=!o,s.options.plugins.title.display=v&&!!e,s.options.plugins.title.text=e||"",s.update(),()=>{}},[t,a,e,o,v,n,l,p]),G(()=>()=>{m.current&&(m.current.destroy(),m.current=null)},[]),i`<canvas ref=${d} height=${c}></canvas>`}function te({title:e,points:t,keys:a,diff:o=!1,inlineTitle:c=!0}){let[v,n]=h(!1),[l,p]=h(null),d=U(!1);return i`
      <div className="chart-wrap">
        <button className="chart-zoom-button" type="button" onClick=${()=>n(!0)}>
          <span className="zoom-icon" aria-hidden="true">+</span>
          Zoom
        </button>
        <div
          className="chart-click"
          onClick=${()=>{if(d.current){d.current=!1;return}n(!0)}}
        >
          <${ve}
            title=${e}
            points=${t}
            keys=${a}
            diff=${o}
            height=${120}
            showTitle=${c&&!!e}
            selectedLabel=${l}
            onPointSelect=${m=>{p(s=>s===m?null:m),d.current=!0,setTimeout(()=>{d.current=!1},0)}}
            onLegendSelect=${m=>{p(s=>s===m?null:m),d.current=!0,setTimeout(()=>{d.current=!1},0)}}
          />
        </div>
        ${v&&i`
          <div className="chart-overlay" onClick=${()=>n(!1)}>
            <div className="chart-modal" onClick=${m=>m.stopPropagation()}>
              <div className="row chart-modal-header">
                <div>
                  <h3>${e||"Chart"}</h3>
                  ${o?i`<p className="muted">Per-second delta.</p>`:""}
                </div>
                <button className="btn ghost" onClick=${()=>n(!1)}>
                  Close
                </button>
              </div>
              <div className="chart-zoom">
                <${ve}
                  title=${e}
                  points=${t}
                  keys=${a}
                  diff=${o}
                  height=${360}
                  showTitle=${!1}
                  selectedLabel=${l}
                  onPointSelect=${m=>p(s=>s===m?null:m)}
                  onLegendSelect=${m=>p(s=>s===m?null:m)}
                />
              </div>
            </div>
          </div>
        `}
      </div>
    `}function be({children:e}){return e}function Fe({toasts:e,onDismiss:t}){return i`
      <div className="toast-stack">
        ${e.map(a=>i`
            <div className=${`toast ${a.kind}`}>
              <span>${a.message}</span>
              <button className="toast-close" onClick=${()=>t(a.id)}>
                ×
              </button>
            </div>
          `)}
      </div>
    `}function Ae({status:e}){return i`
      <header>
        <div>
          <div style=${{fontSize:20,fontWeight:700}}>Vatran</div>
          <div className="status-pill">
            <span className=${`dot ${e.ready?"ok":""}`}></span>
            ${e.ready?"Ready":"Not ready"}
          </div>
        </div>
        <nav>
          <${L} to="/" end className=${({isActive:t})=>t?"active":""}>
            Dashboard
          </${L}>
          <${L} to="/stats/global" className=${({isActive:t})=>t?"active":""}>
            Global stats
          </${L}>
          <${L} to="/stats/real" className=${({isActive:t})=>t?"active":""}>
            Per-real stats
          </${L}>
          <${L}
            to="/target-groups"
            className=${({isActive:t})=>t?"active":""}
          >
            Target groups
          </${L}>
          <${L} to="/config" className=${({isActive:t})=>t?"active":""}>
            Config export
          </${L}>
        </nav>
      </header>
    `}function Ge(){let{addToast:e}=W(),[t,a]=h({initialized:!1,ready:!1}),[o,c]=h([]),[v,n]=h(""),[l,p]=h(!1),[d,m]=h(!1),[s,g]=h({main_interface:"",balancer_prog_path:"",healthchecking_prog_path:"",default_mac:"",local_mac:"",root_map_path:"",root_map_pos:2,katran_src_v4:"",katran_src_v6:"",use_root_map:!1,max_vips:512,max_reals:4096,hash_function:"maglev_v2",forwarding_cores:"",numa_nodes:""}),[y,w]=h({address:"",port:80,proto:6,flags:0}),P=async()=>{try{let r=await $.get("/lb/status"),b=await $.get("/vips"),f=await Promise.all((b||[]).map(async S=>{let C=null,R=!1;try{C=(await $.get("/vips/flags",{address:S.address,port:S.port,proto:S.proto}))?.flags??0}catch{C=null}try{let x=await $.get("/vips/reals",{address:S.address,port:S.port,proto:S.proto});R=Array.isArray(x)&&x.some(Z=>!!Z?.healthy)}catch{R=!1}return{...S,flags:C,healthy:R}}));a(r||{initialized:!1,ready:!1}),c(f),n("")}catch(r){n(r.message||"request failed")}};G(()=>{let r=!0;return(async()=>{r&&await P()})(),()=>{r=!1}},[]);let _=async r=>{r.preventDefault();try{let b=ie(s.forwarding_cores,"Forwarding cores"),f=ie(s.numa_nodes,"NUMA nodes"),S={...s,forwarding_cores:b,numa_nodes:f,root_map_pos:s.root_map_pos===""?void 0:Number(s.root_map_pos),max_vips:Number(s.max_vips),max_reals:Number(s.max_reals),hash_function:s.hash_function};await $.post("/lb/create",S),n(""),p(!1),e("Load balancer initialized.","success"),await P()}catch(b){n(b.message||"request failed"),e(b.message||"Initialize failed.","error")}},k=async r=>{r.preventDefault();try{await $.post("/vips",{...y,port:Number(y.port),proto:Number(y.proto),flags:Number(y.flags||0)}),w({address:"",port:80,proto:6,flags:0}),n(""),m(!1),e("VIP created.","success"),await P()}catch(b){n(b.message||"request failed"),e(b.message||"VIP create failed.","error")}},T=async()=>{try{await $.post("/lb/load-bpf-progs"),n(""),e("BPF programs loaded.","success"),await P()}catch(r){n(r.message||"request failed"),e(r.message||"Load BPF programs failed.","error")}},F=async()=>{try{await $.post("/lb/attach-bpf-progs"),n(""),e("BPF programs attached.","success"),await P()}catch(r){n(r.message||"request failed"),e(r.message||"Attach BPF programs failed.","error")}};return i`
      <main>
        <section className="card">
          <h2>Load balancer</h2>
          <p>Initialized: ${t.initialized?"yes":"no"}</p>
          <p>Ready: ${t.ready?"yes":"no"}</p>
          <div className="row">
            ${!t.initialized&&i`
              <button className="btn" onClick=${()=>p(r=>!r)}>
                ${l?"Close":"Initialize"}
              </button>
            `}
            <button className="btn secondary" onClick=${()=>m(r=>!r)}>
              ${d?"Close":"Create VIP"}
            </button>
          </div>
          ${!t.ready&&i`
            <div className="row" style=${{marginTop:12}}>
              <button
                className="btn ghost"
                disabled=${!t.initialized}
                onClick=${T}
              >
                Load BPF Programs
              </button>
              <button
                className="btn ghost"
                disabled=${!t.initialized}
                onClick=${F}
              >
                Attach BPF Programs
              </button>
            </div>
          `}
          ${l&&i`
            <form className="form" onSubmit=${_}>
              <div className="form-row">
                <label className="field">
                  <span>Main interface</span>
                  <input
                    value=${s.main_interface}
                    onInput=${r=>g({...s,main_interface:r.target.value})}
                    placeholder="eth0"
                    required
                  />
                </label>
                <label className="field">
                  <span>Balancer prog path</span>
                  <input
                    value=${s.balancer_prog_path}
                    onInput=${r=>g({...s,balancer_prog_path:r.target.value})}
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
                    onInput=${r=>g({...s,healthchecking_prog_path:r.target.value})}
                    placeholder="/path/to/hc.o"
                  />
                </label>
                <label className="field">
                  <span>Default MAC</span>
                  <input
                    value=${s.default_mac}
                    onInput=${r=>g({...s,default_mac:r.target.value})}
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
                    onInput=${r=>g({...s,local_mac:r.target.value})}
                    placeholder="11:22:33:44:55:66"
                    required=${s.healthchecking_prog_path?.trim()!==""}
                  />
                </label>
                <label className="field">
                  <span>Hash function</span>
                  <select
                    value=${s.hash_function}
                    onInput=${r=>g({...s,hash_function:r.target.value})}
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
                    value=${s.root_map_path}
                    onInput=${r=>g({...s,root_map_path:r.target.value})}
                    placeholder="/sys/fs/bpf/root_map"
                  />
                </label>
                <label className="field">
                  <span>Root map position</span>
                  <input
                    type="number"
                    min="0"
                    value=${s.root_map_pos}
                    onInput=${r=>g({...s,root_map_pos:r.target.value})}
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Katran src v4</span>
                  <input
                    value=${s.katran_src_v4}
                    onInput=${r=>g({...s,katran_src_v4:r.target.value})}
                    placeholder="10.0.0.1"
                  />
                </label>
                <label className="field">
                  <span>Katran src v6</span>
                  <input
                    value=${s.katran_src_v6}
                    onInput=${r=>g({...s,katran_src_v6:r.target.value})}
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
                    onInput=${r=>g({...s,max_vips:r.target.value})}
                  />
                </label>
                <label className="field">
                  <span>Max reals</span>
                  <input
                    type="number"
                    min="1"
                    value=${s.max_reals}
                    onInput=${r=>g({...s,max_reals:r.target.value})}
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Forwarding cores (optional)</span>
                  <input
                    value=${s.forwarding_cores}
                    onInput=${r=>g({...s,forwarding_cores:r.target.value})}
                    placeholder="0,1,2,3"
                  />
                  <span className="muted">Comma or space separated CPU core IDs.</span>
                </label>
                <label className="field">
                  <span>NUMA nodes (optional)</span>
                  <input
                    value=${s.numa_nodes}
                    onInput=${r=>g({...s,numa_nodes:r.target.value})}
                    placeholder="0,0,1,1"
                  />
                  <span className="muted">Match the forwarding cores length.</span>
                </label>
              </div>
              <label className="field checkbox">
                <input
                  type="checkbox"
                  checked=${s.use_root_map}
                  onChange=${r=>g({...s,use_root_map:r.target.checked})}
                />
                <span>Use root map</span>
              </label>
              <button className="btn" type="submit">Initialize LB</button>
            </form>
          `}
          ${d&&i`
            <form className="form" onSubmit=${k}>
              <div className="form-row">
                <label className="field">
                  <span>VIP address</span>
                  <input
                    value=${y.address}
                    onInput=${r=>w({...y,address:r.target.value})}
                    placeholder="1.2.3.4"
                    required
                  />
                </label>
                <label className="field">
                  <span>Port</span>
                  <input
                    type="number"
                    value=${y.port}
                    onInput=${r=>w({...y,port:r.target.value})}
                    required
                  />
                </label>
                <label className="field">
                  <span>Protocol</span>
                  <select
                    value=${y.proto}
                    onChange=${r=>w({...y,proto:r.target.value})}
                  >
                    <option value="6">TCP (6)</option>
                    <option value="17">UDP (17)</option>
                  </select>
                </label>
                <label className="field">
                  <span>Flags</span>
                  <${de}
                    options=${H}
                    value=${y.flags}
                    name="vip-add"
                    onChange=${r=>w({...y,flags:r})}
                  />
                </label>
              </div>
              <button className="btn" type="submit">Create VIP</button>
            </form>
          `}
          ${v&&i`<p className="error">${v}</p>`}
        </section>
        <section className="card">
          <div className="section-header">
            <h2>VIPs</h2>
            <button className="btn ghost" onClick=${P}>Refresh</button>
          </div>
          ${o.length===0?i`<p className="muted">No VIPs configured yet.</p>`:i`
                <div className="grid">
                  ${o.map(r=>i`
                      <div className="card">
                        <div className="row" style=${{fontWeight:600,gap:8}}>
                          <span className=${`dot ${r.healthy?"ok":"bad"}`}></span>
                          ${r.address}:${r.port} / ${r.proto}
                        </div>
                        <div className="muted" style=${{marginTop:6}}>
                          <div style=${{fontWeight:600,marginBottom:6}}>Flags</div>
                          <${ce}
                            mask=${r.flags}
                            options=${H}
                            emptyLabel=${r.flags===null?"Unknown":"No flags"}
                          />
                        </div>
                        <div className="row" style=${{marginTop:12}}>
                          <${j} className="btn" to=${`/vips/${J(r)}`}>
                            Open
                          </${j}>
                          <${j}
                            className="btn secondary"
                            to=${`/vips/${J(r)}/stats`}
                          >
                            Stats
                          </${j}>
                        </div>
                      </div>
                    `)}
                </div>
              `}
        </section>
      </main>
    `}function Ee(){let{addToast:e}=W(),t=le(),a=Ce(),o=E(()=>ee(t.vipId),[t.vipId]),[c,v]=h([]),[n,l]=h(""),[p,d]=h(""),[m,s]=h(!0),[g,y]=h({address:"",weight:100,flags:0}),[w,P]=h({}),[_,k]=h(null),[T,F]=h({flags:0,set:!0}),[r,b]=h({hash_function:0}),{groups:f,setGroups:S,refreshFromStorage:C,importFromRunningConfig:R}=fe(),[x,Z]=h(""),[$e,O]=h(""),[Ne,ye]=h(!1),[we,_e]=h(""),[ae,Ie]=h({add:0,update:0,remove:0}),z=async()=>{try{let u=await $.get("/vips/reals",o);v(u||[]);let N={};(u||[]).forEach(B=>{N[B.address]=B.weight}),P(N),l(""),s(!1)}catch(u){l(u.message||"request failed"),s(!1)}},se=async()=>{try{let u=await $.get("/vips/flags",o);k(u?.flags??0),d("")}catch(u){d(u.message||"request failed")}};G(()=>{z(),se()},[t.vipId]),G(()=>{if(!x){Ie({add:0,update:0,remove:0});return}let u=f[x]||[],N=new Map(c.map(A=>[A.address,A])),B=new Map(u.map(A=>[A.address,A])),D=0,M=0,I=0;u.forEach(A=>{let re=N.get(A.address);if(!re){D+=1;return}(Number(re.weight)!==Number(A.weight)||Number(re.flags||0)!==Number(A.flags||0))&&(M+=1)}),c.forEach(A=>{B.has(A.address)||(I+=1)}),Ie({add:D,update:M,remove:I})},[x,c,f]);let je=async u=>{try{let N=Number(w[u.address]);await $.post("/vips/reals",{vip:o,real:{address:u.address,weight:N,flags:u.flags||0}}),await z(),e("Real weight updated.","success")}catch(N){l(N.message||"request failed"),e(N.message||"Update failed.","error")}},He=async u=>{try{await $.del("/vips/reals",{vip:o,real:{address:u.address,weight:u.weight,flags:u.flags||0}}),await z(),e("Real removed.","success")}catch(N){l(N.message||"request failed"),e(N.message||"Remove failed.","error")}},We=async u=>{u.preventDefault();try{await $.post("/vips/reals",{vip:o,real:{address:g.address,weight:Number(g.weight),flags:Number(g.flags||0)}}),y({address:"",weight:100,flags:0}),await z(),e("Real added.","success")}catch(N){l(N.message||"request failed"),e(N.message||"Add failed.","error")}},Je=async()=>{if(!x||!f[x]){O("Select a target group to apply.");return}ye(!0),O("");let u=f[x]||[],N=new Map(c.map(I=>[I.address,I])),B=new Map(u.map(I=>[I.address,I])),D=c.filter(I=>!B.has(I.address)),M=u.filter(I=>{let A=N.get(I.address);return A?Number(A.weight)!==Number(I.weight)||Number(A.flags||0)!==Number(I.flags||0):!0});try{D.length>0&&await $.put("/vips/reals/batch",{vip:o,action:1,reals:D.map(I=>({address:I.address,weight:Number(I.weight),flags:Number(I.flags||0)}))}),M.length>0&&await Promise.all(M.map(I=>$.post("/vips/reals",{vip:o,real:{address:I.address,weight:Number(I.weight),flags:Number(I.flags||0)}}))),await z(),e(`Applied target group "${x}".`,"success")}catch(I){O(I.message||"Failed to apply target group."),e(I.message||"Target group apply failed.","error")}finally{ye(!1)}},Ke=u=>{u.preventDefault();let N=we.trim();if(!N){O("Provide a name for the new target group.");return}if(f[N]){O("A target group with that name already exists.");return}let B={...f,[N]:c.map(D=>({address:D.address,weight:Number(D.weight),flags:Number(D.flags||0)}))};S(B),_e(""),Z(N),O(""),e(`Target group "${N}" saved.`,"success")},Ye=async()=>{try{await $.del("/vips",o),e("VIP deleted.","success"),a("/")}catch(u){l(u.message||"request failed"),e(u.message||"Delete failed.","error")}},Ze=async u=>{u.preventDefault();try{await $.put("/vips/flags",{...o,flag:Number(T.flags||0),set:!!T.set}),await se(),e("VIP flags updated.","success")}catch(N){d(N.message||"request failed"),e(N.message||"Flag update failed.","error")}},Xe=async u=>{u.preventDefault();try{await $.put("/vips/hash-function",{...o,hash_function:Number(r.hash_function)}),e("Hash function updated.","success")}catch(N){d(N.message||"request failed"),e(N.message||"Hash update failed.","error")}};return i`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>VIP Detail</h2>
              <p className="muted">${o.address}:${o.port} / ${o.proto}</p>
              ${_===null?i`<p className="muted">Flags: —</p>`:i`
                    <div style=${{marginTop:8}}>
                      <${ce}
                        mask=${_}
                        options=${H}
                        showStatus=${!0}
                        emptyLabel="No flags"
                      />
                    </div>
                  `}
            </div>
            <div className="row">
              <button className="btn ghost" onClick=${se}>Refresh flags</button>
              <button className="btn danger" onClick=${Ye}>Delete VIP</button>
            </div>
          </div>
          ${n&&i`<p className="error">${n}</p>`}
          ${p&&i`<p className="error">${p}</p>`}
          ${m?i`<p className="muted">Loading reals…</p>`:i`
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
                    ${c.map(u=>i`
                        <tr>
                          <td>
                            <span
                              className=${`dot ${Number(u.weight)>0?"ok":"bad"}`}
                            ></span>
                          </td>
                          <td>${u.address}</td>
                          <td>
                            <input
                              className="inline-input"
                              type="number"
                              min="0"
                              value=${w[u.address]??u.weight}
                              onInput=${N=>P({...w,[u.address]:N.target.value})}
                            />
                          </td>
                          <td className="row">
                            <button className="btn" onClick=${()=>je(u)}>
                              Update
                            </button>
                            <button className="btn ghost" onClick=${()=>He(u)}>
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
            <form className="form" onSubmit=${Ze}>
              <div className="form-row">
                <label className="field">
                  <span>Flags</span>
                  <${de}
                    options=${H}
                    value=${T.flags}
                    name="vip-flag-change"
                    onChange=${u=>F({...T,flags:u})}
                  />
                </label>
                <label className="field">
                  <span>Set</span>
                  <select
                    value=${String(T.set)}
                    onChange=${u=>F({...T,set:u.target.value==="true"})}
                  >
                    <option value="true">Set</option>
                    <option value="false">Clear</option>
                  </select>
                </label>
              </div>
              <button className="btn" type="submit">Apply flag</button>
            </form>
            <form className="form" onSubmit=${Xe}>
              <div className="form-row">
                <label className="field">
                  <span>Hash function</span>
                  <input
                    type="number"
                    min="0"
                    value=${r.hash_function}
                    onInput=${u=>b({...r,hash_function:u.target.value})}
                  />
                </label>
              </div>
              <button className="btn" type="submit">Apply hash</button>
            </form>
          </div>
        </section>
        <section className="card">
          <h3>Add real</h3>
          <form className="form" onSubmit=${We}>
            <div className="form-row">
              <label className="field">
                <span>Real address</span>
                <input
                  value=${g.address}
                  onInput=${u=>y({...g,address:u.target.value})}
                  placeholder="10.0.0.1"
                  required
                />
              </label>
              <label className="field">
                <span>Weight</span>
                <input
                  type="number"
                  min="0"
                  value=${g.weight}
                  onInput=${u=>y({...g,weight:u.target.value})}
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
              <button className="btn ghost" type="button" onClick=${C}>
                Reload groups
              </button>
              <button
                className="btn ghost"
                type="button"
                onClick=${async()=>{try{await R(),e("Imported target groups from running config.","success")}catch(u){O(u.message||"Failed to import target groups."),e(u.message||"Import failed.","error")}}}
              >
                Import from running config
              </button>
            </div>
          </div>
          ${$e&&i`<p className="error">${$e}</p>`}
          <div className="form-row">
            <label className="field">
              <span>Target group</span>
              <select
                value=${x}
                onChange=${u=>Z(u.target.value)}
                disabled=${Object.keys(f).length===0}
              >
                <option value="">Select group</option>
                ${Object.keys(f).map(u=>i`<option value=${u}>${u}</option>`)}
              </select>
            </label>
            <label className="field">
              <span>Preview</span>
              <input
                value=${`add ${ae.add} \xB7 update ${ae.update} \xB7 remove ${ae.remove}`}
                readOnly
              />
            </label>
          </div>
          <div className="row">
            <button
              className="btn"
              type="button"
              onClick=${Je}
              disabled=${Ne||!x}
            >
              ${Ne?"Applying...":"Apply target group"}
            </button>
          </div>
          <form className="form" onSubmit=${Ke}>
            <div className="form-row">
              <label className="field">
                <span>Save current reals as new group</span>
                <input
                  value=${we}
                  onInput=${u=>_e(u.target.value)}
                  placeholder="edge-backends"
                />
              </label>
            </div>
            <button className="btn secondary" type="submit">Save target group</button>
          </form>
        </section>
      </main>
    `}function Le(){let e=le(),t=E(()=>ee(e.vipId),[e.vipId]),{points:a,error:o}=K({path:"/stats/vip",body:t}),c=E(()=>Y("/stats/vip"),[]),v=a[a.length-1]||{},n=a[a.length-2]||{},l=Number(v.v1??0),p=Number(v.v2??0),d=l-Number(n.v1??0),m=p-Number(n.v2??0),s=E(()=>[{label:c.v1,field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.15)"},{label:c.v2,field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[c]);return i`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>VIP Stats</h2>
              <p className="muted">${t.address}:${t.port} / ${t.proto}</p>
            </div>
          </div>
          ${o&&i`<p className="error">${o}</p>`}
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
                <td>${c.v1}</td>
                <td>${l}</td>
                <td>
                  <span className=${`delta ${d<0?"down":"up"}`}>
                    ${V(d)}
                  </span>
                </td>
              </tr>
              <tr>
                <td>${c.v2}</td>
                <td>${p}</td>
                <td>
                  <span className=${`delta ${m<0?"down":"up"}`}>
                    ${V(m)}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
          <${te} title="Traffic (delta/sec)" points=${a} keys=${s} diff=${!0} />
        </section>
      </main>
    `}let De={"/stats/vip":{v1:"Packets",v2:"Bytes"},"/stats/real":{v1:"Packets",v2:"Bytes"},"/stats/lru":{v1:"Total packets",v2:"LRU hits"},"/stats/lru/miss":{v1:"TCP SYN misses",v2:"Non-SYN misses"},"/stats/lru/fallback":{v1:"Fallback LRU hits",v2:"Unused"},"/stats/lru/global":{v1:"Map lookup failures",v2:"Global LRU routed"},"/stats/xdp/total":{v1:"Packets",v2:"Bytes"},"/stats/xdp/pass":{v1:"Packets",v2:"Bytes"},"/stats/xdp/drop":{v1:"Packets",v2:"Bytes"},"/stats/xdp/tx":{v1:"Packets",v2:"Bytes"}},he=[{title:"LRU",path:"/stats/lru"},{title:"LRU Miss",path:"/stats/lru/miss"},{title:"LRU Fallback",path:"/stats/lru/fallback"},{title:"LRU Global",path:"/stats/lru/global"},{title:"XDP Total",path:"/stats/xdp/total"},{title:"XDP Pass",path:"/stats/xdp/pass"},{title:"XDP Drop",path:"/stats/xdp/drop"},{title:"XDP Tx",path:"/stats/xdp/tx"}];function V(e){return Number.isFinite(e)?`${e>0?"+":""}${e}`:"0"}function Y(e){return De[e]||{v1:"v1",v2:"v2"}}function Be({title:e,path:t,diff:a=!1}){let{points:o,error:c}=K({path:t}),v=E(()=>Y(t),[t]),n=E(()=>[{label:v.v1,field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.2)"},{label:v.v2,field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[v]);return i`
      <div className="card">
        <h3>${e}</h3>
        ${c&&i`<p className="error">${c}</p>`}
        <${te} title=${e} points=${o} keys=${n} diff=${a} inlineTitle=${!1} />
      </div>
    `}function qe({title:e,path:t}){let{points:a,error:o}=K({path:t}),c=E(()=>Y(t),[t]),v=a[a.length-1]||{},n=a[a.length-2]||{},l=Number(v.v1??0),p=Number(v.v2??0),d=l-Number(n.v1??0),m=p-Number(n.v2??0);return i`
      <div className="summary-card">
        <div className="summary-title">${e}</div>
        ${o?i`<p className="error">${o}</p>`:i`
              <div className="summary-row">
                <div className="stat">
                  <span className="muted">${c.v1}</span>
                  <strong>${l}</strong>
                </div>
                <div className="stat">
                  <span className="muted">${c.v1} delta/sec</span>
                  <strong className=${d<0?"delta down":"delta up"}>
                    ${V(d)}
                  </strong>
                </div>
              </div>
              <div className="summary-row">
                <div className="stat">
                  <span className="muted">${c.v2}</span>
                  <strong>${p}</strong>
                </div>
                <div className="stat">
                  <span className="muted">${c.v2} delta/sec</span>
                  <strong className=${m<0?"delta down":"delta up"}>
                    ${V(m)}
                  </strong>
                </div>
              </div>
            `}
      </div>
    `}function Oe(){let{data:e,error:t}=xe(()=>$.get("/stats/userspace"),1e3,[]);return i`
      <main>
        <section className="card">
          <h2>Global Stats</h2>
          <p className="muted">Polling every second. Charts show per-second deltas.</p>
        </section>
        <section className="grid">
          ${he.map(a=>i`<${Be} title=${a.title} path=${a.path} diff=${!0} />`)}
        </section>
        <section className="card">
          <h3>Userspace</h3>
          ${t&&i`<p className="error">${t}</p>`}
          ${e?i`
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
              `:i`<p className="muted">Waiting for data…</p>`}
        </section>
        <section className="card">
          <div className="section-header">
            <div>
              <h3>Absolute & Rate of Change</h3>
              <p className="muted">Latest value and per-second delta.</p>
            </div>
          </div>
          <div className="summary-grid">
            ${he.map(a=>i`<${qe} title=${a.title} path=${a.path} />`)}
          </div>
        </section>
      </main>
    `}function Ve(){let[e,t]=h([]),[a,o]=h(""),[c,v]=h([]),[n,l]=h(""),[p,d]=h(null),[m,s]=h("");G(()=>{let f=!0;return(async()=>{try{let C=await $.get("/vips");if(!f)return;t(C||[]),!a&&C&&C.length>0&&o(J(C[0]))}catch(C){f&&s(C.message||"request failed")}})(),()=>{f=!1}},[]),G(()=>{if(!a)return;let f=ee(a),S=!0;return(async()=>{try{let R=await $.get("/vips/reals",f);if(!S)return;v(R||[]),R&&R.length>0?l(x=>x||R[0].address):l(""),s("")}catch(R){S&&s(R.message||"request failed")}})(),()=>{S=!1}},[a]),G(()=>{if(!n){d(null);return}let f=!0;return(async()=>{try{let C=await $.get("/reals/index",{address:n});if(!f)return;d(C?.index??null),s("")}catch(C){f&&s(C.message||"request failed")}})(),()=>{f=!1}},[n]);let{points:g,error:y}=K({path:"/stats/real",body:p!==null?{index:p}:null}),w=E(()=>Y("/stats/real"),[]),P=E(()=>[{label:w.v1,field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.2)"},{label:w.v2,field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[w]),_=g[g.length-1]||{},k=g[g.length-2]||{},T=Number(_.v1??0),F=Number(_.v2??0),r=T-Number(k.v1??0),b=F-Number(k.v2??0);return i`
      <main>
        <section className="card">
          <h2>Per-Real Stats</h2>
          <p className="muted">Select a VIP and real address to chart.</p>
          ${m&&i`<p className="error">${m}</p>`}
          <div className="form-row">
            <label className="field">
              <span>VIP</span>
              <select value=${a} onChange=${f=>o(f.target.value)}>
                ${e.map(f=>i`
                    <option value=${J(f)}>
                      ${f.address}:${f.port} / ${f.proto}
                    </option>
                  `)}
              </select>
            </label>
            <label className="field">
              <span>Real</span>
              <select
                value=${n}
                onChange=${f=>l(f.target.value)}
                disabled=${c.length===0}
              >
                ${c.map(f=>i`
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
          ${y&&i`<p className="error">${y}</p>`}
          ${p===null?i`<p className="muted">Select a real to start polling.</p>`:i`
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
                      <td>${w.v1}</td>
                      <td>${T}</td>
                      <td>
                        <span className=${`delta ${r<0?"down":"up"}`}>
                          ${V(r)}
                        </span>
                      </td>
                    </tr>
                    <tr>
                      <td>${w.v2}</td>
                      <td>${F}</td>
                      <td>
                        <span className=${`delta ${b<0?"down":"up"}`}>
                          ${V(b)}
                        </span>
                      </td>
                    </tr>
                  </tbody>
                </table>
                <${te}
                  title="Traffic (delta/sec)"
                  points=${g}
                  keys=${P}
                  diff=${!0}
                />
              `}
        </section>
      </main>
    `}function Ue(){let{addToast:e}=W(),[t,a]=h(""),[o,c]=h(""),[v,n]=h(!0),[l,p]=h(""),d=U(!0),m=async()=>{if(d.current){n(!0),c("");try{let g=await fetch(`${$.base}/config/export`,{headers:{Accept:"application/x-yaml"}});if(!g.ok){let w=`HTTP ${g.status}`;try{w=(await g.json())?.error?.message||w}catch{}throw new Error(w)}let y=await g.text();if(!d.current)return;a(y||""),p(new Date().toLocaleString())}catch(g){d.current&&c(g.message||"request failed")}finally{d.current&&n(!1)}}},s=async()=>{if(t)try{await navigator.clipboard.writeText(t),e("Config copied to clipboard","info")}catch{e("Failed to copy config","error")}};return G(()=>(d.current=!0,m(),()=>{d.current=!1}),[]),i`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>Running config</h2>
              <p className="muted">Exported from /api/v1/config/export</p>
            </div>
            <div className="row">
              <button className="btn ghost" onClick=${s} disabled=${!t}>
                Copy YAML
              </button>
              <button className="btn" onClick=${m} disabled=${v}>
                Refresh
              </button>
            </div>
          </div>
          ${o&&i`<p className="error">${o}</p>`}
          ${v?i`<p className="muted">Loading config...</p>`:t?i`<pre className="yaml-view">${t}</pre>`:i`<p className="muted">No config data returned.</p>`}
          ${l&&i`<p className="muted">Last fetched ${l}</p>`}
        </section>
      </main>
    `}function ze(){let{addToast:e}=W(),{groups:t,setGroups:a,refreshFromStorage:o,importFromRunningConfig:c}=fe(),[v,n]=h(""),[l,p]=h(""),[d,m]=h({address:"",weight:100,flags:0}),[s,g]=h(""),[y,w]=h(!1);G(()=>{if(l){if(!t[l]){let b=Object.keys(t);p(b[0]||"")}}else{let b=Object.keys(t);b.length>0&&p(b[0])}},[t,l]);let P=b=>{b.preventDefault();let f=v.trim();if(!f){g("Provide a group name.");return}if(t[f]){g("That group already exists.");return}a({...t,[f]:[]}),n(""),p(f),g(""),e(`Target group "${f}" created.`,"success")},_=b=>{let f={...t};delete f[b],a(f),e(`Target group "${b}" removed.`,"success")},k=b=>{if(b.preventDefault(),!l){g("Select a group to add a real.");return}let f=pe(d);if(!f){g("Provide a valid real address.");return}let S=t[l]||[],C=S.some(R=>R.address===f.address)?S.map(R=>R.address===f.address?f:R):S.concat(f);a({...t,[l]:C}),m({address:"",weight:100,flags:0}),g(""),e("Real saved to target group.","success")},T=b=>{if(!l)return;let S=(t[l]||[]).filter(C=>C.address!==b);a({...t,[l]:S})},F=(b,f)=>{if(!l)return;let C=(t[l]||[]).map(R=>R.address===b?{...R,...f}:R);a({...t,[l]:C})};return i`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>Target groups</h2>
              <p className="muted">Define reusable sets of reals (address + weight).</p>
            </div>
            <div className="row">
              <button className="btn ghost" type="button" onClick=${o}>
                Reload groups
              </button>
              <button className="btn ghost" type="button" onClick=${async()=>{w(!0);try{await c(),e("Imported target groups from running config.","success"),g("")}catch(b){g(b.message||"Failed to import target groups."),e(b.message||"Import failed.","error")}finally{w(!1)}}} disabled=${y}>
                ${y?"Importing...":"Import from running config"}
              </button>
            </div>
          </div>
          ${s&&i`<p className="error">${s}</p>`}
          <form className="form" onSubmit=${P}>
            <div className="form-row">
              <label className="field">
                <span>New group name</span>
                <input
                  value=${v}
                  onInput=${b=>n(b.target.value)}
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
                  value=${l}
                  onChange=${b=>p(b.target.value)}
                  disabled=${Object.keys(t).length===0}
                >
                  ${Object.keys(t).map(b=>i`<option value=${b}>${b}</option>`)}
                </select>
              </label>
              ${l&&i`<button className="btn danger" type="button" onClick=${()=>_(l)}>
                Delete group
              </button>`}
            </div>
          </div>
          ${l?i`
                <table className="table">
                  <thead>
                    <tr>
                      <th>Address</th>
                      <th>Weight</th>
                      <th></th>
                    </tr>
                  </thead>
                  <tbody>
                    ${(t[l]||[]).map(b=>i`
                        <tr>
                          <td>${b.address}</td>
                          <td>
                            <input
                              className="inline-input"
                              type="number"
                              min="0"
                              value=${b.weight}
                              onInput=${f=>F(b.address,{weight:Number(f.target.value)})}
                            />
                          </td>
                          <td>
                            <button
                              className="btn ghost"
                              type="button"
                              onClick=${()=>T(b.address)}
                            >
                              Remove
                            </button>
                          </td>
                        </tr>
                      `)}
                  </tbody>
                </table>
                <form className="form" onSubmit=${k}>
                  <div className="form-row">
                    <label className="field">
                      <span>Real address</span>
                      <input
                        value=${d.address}
                        onInput=${b=>m({...d,address:b.target.value})}
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
                        onInput=${b=>m({...d,weight:b.target.value})}
                        required
                      />
                    </label>
                  </div>
                  <button className="btn secondary" type="submit">Save real</button>
                </form>
              `:i`<p className="muted">No groups yet. Create one to add reals.</p>`}
        </section>
      </main>
    `}function Me(){let[e,t]=h({initialized:!1,ready:!1}),[a,o]=h([]),c=U({}),v=(l,p="info")=>{let d=`${Date.now()}-${Math.random().toString(16).slice(2)}`;o(m=>m.concat({id:d,message:l,kind:p})),c.current[d]=setTimeout(()=>{o(m=>m.filter(s=>s.id!==d)),delete c.current[d]},4e3)},n=l=>{c.current[l]&&(clearTimeout(c.current[l]),delete c.current[l]),o(p=>p.filter(d=>d.id!==l))};return G(()=>{let l=!0,p=async()=>{try{let m=await $.get("/lb/status");l&&t(m||{initialized:!1,ready:!1})}catch{l&&t({initialized:!1,ready:!1})}};p();let d=setInterval(p,5e3);return()=>{l=!1,clearInterval(d)}},[]),i`
      <${ne}>
        <${be}>
          <${X.Provider} value=${{addToast:v}}>
            <${Ae} status=${e} />
            <${oe}>
              <${q} path="/" element=${i`<${Ge} />`} />
              <${q} path="/vips/:vipId" element=${i`<${Ee} />`} />
              <${q} path="/vips/:vipId/stats" element=${i`<${Le} />`} />
              <${q} path="/target-groups" element=${i`<${ze} />`} />
              <${q} path="/stats/global" element=${i`<${Oe} />`} />
              <${q} path="/stats/real" element=${i`<${Ve} />`} />
              <${q} path="/config" element=${i`<${Ue} />`} />
            </${oe}>
            <${Fe} toasts=${a} onDismiss=${n} />
          </${X.Provider}>
        </${be}>
      </${ne}>
    `}ReactDOM.createRoot(document.getElementById("root")).render(i`<${Me} />`)})();})();
