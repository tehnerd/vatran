(()=>{(()=>{let{useEffect:A,useMemo:E,useRef:U,useState:h,useContext:Se}=React,{BrowserRouter:re,Routes:ne,Route:B,NavLink:L,Link:j,useParams:oe,useNavigate:Ce}=ReactRouterDOM,i=htm.bind(React.createElement),Z=React.createContext({addToast:()=>{}}),H=[{label:"NO_SPORT",value:1},{label:"NO_LRU",value:2},{label:"QUIC_VIP",value:4},{label:"DPORT_HASH",value:8},{label:"SRC_ROUTING",value:16},{label:"LOCAL_VIP",value:32},{label:"GLOBAL_LRU",value:64},{label:"HASH_SRC_DST_PORT",value:128},{label:"UDP_STABLE_ROUTING_VIP",value:256},{label:"UDP_FLOW_MIGRATION",value:512}];function W(){return Se(Z)}function Re(e){return String(e).replace(/[^a-z0-9_-]/gi,"_")}function Te(e,t,a){let o=Number(e)||0,c=Number(t)||0;return a?o|c:o&~c}function Pe(e,t){let a=Number(e)||0;return t.filter(o=>(a&o.value)!==0)}function le(e,t){let a=String(e??"").trim();if(!a)return;let c=a.split(/[\s,]+/).filter(Boolean).map(r=>Number(r));if(c.findIndex(r=>!Number.isFinite(r)||!Number.isInteger(r))!==-1)throw new Error(`${t} must be a comma- or space-separated list of integers.`);return c}function ie({mask:e,options:t,showStatus:a=!1,emptyLabel:o="None"}){let c=Number(e)||0,b=a?t:Pe(c,t),r=a?2:1;return i`
      <table className="table flag-table">
        <thead>
          <tr>
            <th>Flag</th>
            ${a?i`<th>Enabled</th>`:null}
          </tr>
        </thead>
        <tbody>
          ${b.length===0?i`<tr><td colspan=${r} className="muted">${o}</td></tr>`:b.map(l=>{let p=(c&l.value)!==0;return i`
                  <tr>
                    <td>${l.label}</td>
                    ${a?i`<td>${p?"Yes":"No"}</td>`:null}
                  </tr>
                `})}
        </tbody>
      </table>
    `}function ce({options:e,value:t,onChange:a,name:o}){let c=Number(t)||0,b=Re(o||"flags");return i`
      <div className="flag-selector">
        <table className="table flag-table">
          <thead>
            <tr>
              <th>Enabled</th>
              <th>Flag</th>
            </tr>
          </thead>
          <tbody>
            ${e.map(r=>{let l=`${b}-${r.value}`,p=(c&r.value)===r.value;return i`
                <tr>
                  <td>
                    <input
                      id=${l}
                      type="checkbox"
                      checked=${p}
                      onChange=${d=>a(Te(c,r.value,d.target.checked))}
                    />
                  </td>
                  <td>
                    <label className="flag-option" htmlFor=${l}>${r.label}</label>
                  </td>
                </tr>
              `})}
          </tbody>
        </table>
      </div>
    `}let N={base:"/api/v1",async request(e,t={}){let a={method:t.method||"GET",headers:{"Content-Type":"application/json"}},o=`${N.base}${e}`;if(t.body!==void 0&&t.body!==null)if(a.method==="GET"){let r=new URLSearchParams;Object.entries(t.body).forEach(([p,d])=>{if(d!=null){if(Array.isArray(d)){d.forEach(m=>r.append(p,String(m)));return}if(typeof d=="object"){r.set(p,JSON.stringify(d));return}r.set(p,String(d))}});let l=r.toString();l&&(o+=`${o.includes("?")?"&":"?"}${l}`)}else a.body=JSON.stringify(t.body);let c=await fetch(o,a),b;try{b=await c.json()}catch{throw new Error("invalid JSON response")}if(!c.ok)throw new Error(b?.error?.message||`HTTP ${c.status}`);if(!b.success){let r=b.error?.message||"request failed";throw new Error(r)}return b.data},get(e,t){return N.request(e,{method:"GET",body:t})},post(e,t){return N.request(e,{method:"POST",body:t})},put(e,t){return N.request(e,{method:"PUT",body:t})},del(e,t){return N.request(e,{method:"DELETE",body:t})}},de="vatran_target_groups";function ue(e){if(!e||!e.address)return null;let t=String(e.address).trim();if(!t)return null;let a=Number(e.weight),o=Number(e.flags??0);return{address:t,weight:Number.isFinite(a)?a:0,flags:Number.isFinite(o)?o:0}}function pe(e){if(!e||typeof e!="object")return{};let t={};return Object.entries(e).forEach(([a,o])=>{let c=String(a).trim();if(!c)return;let b=Array.isArray(o)?o.map(ue).filter(Boolean):[],r=[],l=new Set;b.forEach(p=>{l.has(p.address)||(l.add(p.address),r.push(p))}),t[c]=r}),t}function X(){if(typeof localStorage>"u")return{};try{let e=localStorage.getItem(de);return e?pe(JSON.parse(e)):{}}catch{return{}}}function me(e){if(!(typeof localStorage>"u"))try{localStorage.setItem(de,JSON.stringify(e))}catch{}}function ke(e,t){let a={...e};return Object.entries(t||{}).forEach(([o,c])=>{a[o]||(a[o]=c)}),a}function ge(){let[e,t]=h(()=>X());return A(()=>{me(e)},[e]),{groups:e,setGroups:t,refreshFromStorage:()=>{t(X())},importFromRunningConfig:async()=>{let c=await N.get("/config/export/json"),b=pe(c?.target_groups||{}),r=ke(X(),b);return t(r),me(r),r}}}function J(e){return`${encodeURIComponent(e.address)}:${e.port}:${e.proto}`}function Q(e){let t=e.split(":"),a=Number(t.pop()||0),o=Number(t.pop()||0);return{address:decodeURIComponent(t.join(":")),port:o,proto:a}}function xe(e,t,a=[]){let[o,c]=h(null),[b,r]=h(""),[l,p]=h(!0);return A(()=>{let d=!0,m=async()=>{try{let g=await e();d&&(c(g),r(""),p(!1))}catch(g){d&&(r(g.message||"request failed"),p(!1))}};m();let s=setInterval(m,t);return()=>{d=!1,clearInterval(s)}},a),{data:o,error:b,loading:l}}function K({path:e,body:t,intervalMs:a=1e3,limit:o=60}){let[c,b]=h([]),[r,l]=h(""),p=E(()=>JSON.stringify(t||{}),[t]);return A(()=>{if(t===null)return b([]),l(""),()=>{};let d=!0,m=async()=>{try{let g=await N.get(e,t);if(!d)return;let y=new Date().toLocaleTimeString();b(w=>w.concat({label:y,...g}).slice(-o)),l("")}catch(g){d&&l(g.message||"request failed")}};m();let s=setInterval(m,a);return()=>{d=!1,clearInterval(s)}},[e,p,a,o]),{points:c,error:r}}function fe({title:e,points:t,keys:a,diff:o=!1,height:c=120,showTitle:b=!1,selectedLabel:r=null,onPointSelect:l=null,onLegendSelect:p=null}){let d=U(null),m=U(null);return A(()=>{if(!d.current)return;m.current||(m.current=new Chart(d.current,{type:"line",data:{labels:[],datasets:[]},options:{responsive:!0,maintainAspectRatio:!1,animation:!1,scales:{x:{grid:{display:!1}},y:{beginAtZero:!o}},plugins:{legend:{display:!0,position:"bottom"},title:{display:b&&!!e,text:e}}}}));let s=m.current,g=new Map((s.data.datasets||[]).filter(_=>typeof _.hidden<"u").map(_=>[_.label,_.hidden])),y=t.map(_=>_.label),w=r?a.filter(_=>_.label===r):a,T=r&&w.length===0?a:w;return s.data.labels=y,s.data.datasets=T.map(_=>{let P=t.map(x=>x[_.field]||0),R=o?P.map((x,n)=>n===0?0:x-P[n-1]):P;return{label:_.label,data:R,borderColor:_.color,backgroundColor:_.fill,borderWidth:2,tension:.3,hidden:g.get(_.label)}}),s.options.onClick=(_,P)=>{if(!l||!P||P.length===0)return;let R=P[0].datasetIndex,x=s.data.datasets?.[R]?.label;x&&l(x)},s.options.plugins&&s.options.plugins.legend&&(s.options.plugins.legend.onClick=(_,P)=>{if(!p)return;let R=P?.text;R&&p(R)}),s.options.scales.y.beginAtZero=!o,s.options.plugins.title.display=b&&!!e,s.options.plugins.title.text=e||"",s.update(),()=>{}},[t,a,e,o,b,r,l,p]),A(()=>()=>{m.current&&(m.current.destroy(),m.current=null)},[]),i`<canvas ref=${d} height=${c}></canvas>`}function ee({title:e,points:t,keys:a,diff:o=!1,inlineTitle:c=!0}){let[b,r]=h(!1),[l,p]=h(null),d=U(!1);return i`
      <div className="chart-wrap">
        <button className="chart-zoom-button" type="button" onClick=${()=>r(!0)}>
          <span className="zoom-icon" aria-hidden="true">+</span>
          Zoom
        </button>
        <div
          className="chart-click"
          onClick=${()=>{if(d.current){d.current=!1;return}r(!0)}}
        >
          <${fe}
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
        ${b&&i`
          <div className="chart-overlay" onClick=${()=>r(!1)}>
            <div className="chart-modal" onClick=${m=>m.stopPropagation()}>
              <div className="row chart-modal-header">
                <div>
                  <h3>${e||"Chart"}</h3>
                  ${o?i`<p className="muted">Per-second delta.</p>`:""}
                </div>
                <button className="btn ghost" onClick=${()=>r(!1)}>
                  Close
                </button>
              </div>
              <div className="chart-zoom">
                <${fe}
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
    `}function Ge(){let{addToast:e}=W(),[t,a]=h({initialized:!1,ready:!1}),[o,c]=h([]),[b,r]=h(""),[l,p]=h(!1),[d,m]=h(!1),[s,g]=h({main_interface:"",balancer_prog_path:"",healthchecking_prog_path:"",default_mac:"",local_mac:"",root_map_path:"",root_map_pos:2,katran_src_v4:"",katran_src_v6:"",use_root_map:!1,max_vips:512,max_reals:4096,hash_function:"maglev_v2",forwarding_cores:"",numa_nodes:""}),[y,w]=h({address:"",port:80,proto:6,flags:0}),T=async()=>{try{let n=await N.get("/lb/status"),v=await N.get("/vips"),f=await Promise.all((v||[]).map(async S=>{try{let C=await N.get("/vips/flags",{address:S.address,port:S.port,proto:S.proto});return{...S,flags:C?.flags??0}}catch{return{...S,flags:null}}}));a(n||{initialized:!1,ready:!1}),c(f),r("")}catch(n){r(n.message||"request failed")}};A(()=>{let n=!0;return(async()=>{n&&await T()})(),()=>{n=!1}},[]);let _=async n=>{n.preventDefault();try{let v=le(s.forwarding_cores,"Forwarding cores"),f=le(s.numa_nodes,"NUMA nodes"),S={...s,forwarding_cores:v,numa_nodes:f,root_map_pos:s.root_map_pos===""?void 0:Number(s.root_map_pos),max_vips:Number(s.max_vips),max_reals:Number(s.max_reals),hash_function:s.hash_function};await N.post("/lb/create",S),r(""),p(!1),e("Load balancer initialized.","success"),await T()}catch(v){r(v.message||"request failed"),e(v.message||"Initialize failed.","error")}},P=async n=>{n.preventDefault();try{await N.post("/vips",{...y,port:Number(y.port),proto:Number(y.proto),flags:Number(y.flags||0)}),w({address:"",port:80,proto:6,flags:0}),r(""),m(!1),e("VIP created.","success"),await T()}catch(v){r(v.message||"request failed"),e(v.message||"VIP create failed.","error")}},R=async()=>{try{await N.post("/lb/load-bpf-progs"),r(""),e("BPF programs loaded.","success"),await T()}catch(n){r(n.message||"request failed"),e(n.message||"Load BPF programs failed.","error")}},x=async()=>{try{await N.post("/lb/attach-bpf-progs"),r(""),e("BPF programs attached.","success"),await T()}catch(n){r(n.message||"request failed"),e(n.message||"Attach BPF programs failed.","error")}};return i`
      <main>
        <section className="card">
          <h2>Load balancer</h2>
          <p>Initialized: ${t.initialized?"yes":"no"}</p>
          <p>Ready: ${t.ready?"yes":"no"}</p>
          <div className="row">
            ${!t.initialized&&i`
              <button className="btn" onClick=${()=>p(n=>!n)}>
                ${l?"Close":"Initialize"}
              </button>
            `}
            <button className="btn secondary" onClick=${()=>m(n=>!n)}>
              ${d?"Close":"Create VIP"}
            </button>
          </div>
          ${!t.ready&&i`
            <div className="row" style=${{marginTop:12}}>
              <button
                className="btn ghost"
                disabled=${!t.initialized}
                onClick=${R}
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
          `}
          ${l&&i`
            <form className="form" onSubmit=${_}>
              <div className="form-row">
                <label className="field">
                  <span>Main interface</span>
                  <input
                    value=${s.main_interface}
                    onInput=${n=>g({...s,main_interface:n.target.value})}
                    placeholder="eth0"
                    required
                  />
                </label>
                <label className="field">
                  <span>Balancer prog path</span>
                  <input
                    value=${s.balancer_prog_path}
                    onInput=${n=>g({...s,balancer_prog_path:n.target.value})}
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
                    onInput=${n=>g({...s,healthchecking_prog_path:n.target.value})}
                    placeholder="/path/to/hc.o"
                  />
                </label>
                <label className="field">
                  <span>Default MAC</span>
                  <input
                    value=${s.default_mac}
                    onInput=${n=>g({...s,default_mac:n.target.value})}
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
                    onInput=${n=>g({...s,local_mac:n.target.value})}
                    placeholder="11:22:33:44:55:66"
                    required=${s.healthchecking_prog_path?.trim()!==""}
                  />
                </label>
                <label className="field">
                  <span>Hash function</span>
                  <select
                    value=${s.hash_function}
                    onInput=${n=>g({...s,hash_function:n.target.value})}
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
                    onInput=${n=>g({...s,root_map_path:n.target.value})}
                    placeholder="/sys/fs/bpf/root_map"
                  />
                </label>
                <label className="field">
                  <span>Root map position</span>
                  <input
                    type="number"
                    min="0"
                    value=${s.root_map_pos}
                    onInput=${n=>g({...s,root_map_pos:n.target.value})}
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Katran src v4</span>
                  <input
                    value=${s.katran_src_v4}
                    onInput=${n=>g({...s,katran_src_v4:n.target.value})}
                    placeholder="10.0.0.1"
                  />
                </label>
                <label className="field">
                  <span>Katran src v6</span>
                  <input
                    value=${s.katran_src_v6}
                    onInput=${n=>g({...s,katran_src_v6:n.target.value})}
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
                    onInput=${n=>g({...s,max_vips:n.target.value})}
                  />
                </label>
                <label className="field">
                  <span>Max reals</span>
                  <input
                    type="number"
                    min="1"
                    value=${s.max_reals}
                    onInput=${n=>g({...s,max_reals:n.target.value})}
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Forwarding cores (optional)</span>
                  <input
                    value=${s.forwarding_cores}
                    onInput=${n=>g({...s,forwarding_cores:n.target.value})}
                    placeholder="0,1,2,3"
                  />
                  <span className="muted">Comma or space separated CPU core IDs.</span>
                </label>
                <label className="field">
                  <span>NUMA nodes (optional)</span>
                  <input
                    value=${s.numa_nodes}
                    onInput=${n=>g({...s,numa_nodes:n.target.value})}
                    placeholder="0,0,1,1"
                  />
                  <span className="muted">Match the forwarding cores length.</span>
                </label>
              </div>
              <label className="field checkbox">
                <input
                  type="checkbox"
                  checked=${s.use_root_map}
                  onChange=${n=>g({...s,use_root_map:n.target.checked})}
                />
                <span>Use root map</span>
              </label>
              <button className="btn" type="submit">Initialize LB</button>
            </form>
          `}
          ${d&&i`
            <form className="form" onSubmit=${P}>
              <div className="form-row">
                <label className="field">
                  <span>VIP address</span>
                  <input
                    value=${y.address}
                    onInput=${n=>w({...y,address:n.target.value})}
                    placeholder="1.2.3.4"
                    required
                  />
                </label>
                <label className="field">
                  <span>Port</span>
                  <input
                    type="number"
                    value=${y.port}
                    onInput=${n=>w({...y,port:n.target.value})}
                    required
                  />
                </label>
                <label className="field">
                  <span>Protocol</span>
                  <select
                    value=${y.proto}
                    onChange=${n=>w({...y,proto:n.target.value})}
                  >
                    <option value="6">TCP (6)</option>
                    <option value="17">UDP (17)</option>
                  </select>
                </label>
                <label className="field">
                  <span>Flags</span>
                  <${ce}
                    options=${H}
                    value=${y.flags}
                    name="vip-add"
                    onChange=${n=>w({...y,flags:n})}
                  />
                </label>
              </div>
              <button className="btn" type="submit">Create VIP</button>
            </form>
          `}
          ${b&&i`<p className="error">${b}</p>`}
        </section>
        <section className="card">
          <div className="section-header">
            <h2>VIPs</h2>
            <button className="btn ghost" onClick=${T}>Refresh</button>
          </div>
          ${o.length===0?i`<p className="muted">No VIPs configured yet.</p>`:i`
                <div className="grid">
                  ${o.map(n=>i`
                      <div className="card">
                        <div style=${{fontWeight:600}}>
                          ${n.address}:${n.port} / ${n.proto}
                        </div>
                        <div className="muted" style=${{marginTop:6}}>
                          <div style=${{fontWeight:600,marginBottom:6}}>Flags</div>
                          <${ie}
                            mask=${n.flags}
                            options=${H}
                            emptyLabel=${n.flags===null?"Unknown":"No flags"}
                          />
                        </div>
                        <div className="row" style=${{marginTop:12}}>
                          <${j} className="btn" to=${`/vips/${J(n)}`}>
                            Open
                          </${j}>
                          <${j}
                            className="btn secondary"
                            to=${`/vips/${J(n)}/stats`}
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
    `}function Ee(){let{addToast:e}=W(),t=oe(),a=Ce(),o=E(()=>Q(t.vipId),[t.vipId]),[c,b]=h([]),[r,l]=h(""),[p,d]=h(""),[m,s]=h(!0),[g,y]=h({address:"",weight:100,flags:0}),[w,T]=h({}),[_,P]=h(null),[R,x]=h({flags:0,set:!0}),[n,v]=h({hash_function:0}),{groups:f,setGroups:S,refreshFromStorage:C,importFromRunningConfig:k}=ge(),[G,he]=h(""),[$e,O]=h(""),[Ne,ye]=h(!1),[we,_e]=h(""),[te,Ie]=h({add:0,update:0,remove:0}),z=async()=>{try{let u=await N.get("/vips/reals",o);b(u||[]);let $={};(u||[]).forEach(q=>{$[q.address]=q.weight}),T($),l(""),s(!1)}catch(u){l(u.message||"request failed"),s(!1)}},ae=async()=>{try{let u=await N.get("/vips/flags",o);P(u?.flags??0),d("")}catch(u){d(u.message||"request failed")}};A(()=>{z(),ae()},[t.vipId]),A(()=>{if(!G){Ie({add:0,update:0,remove:0});return}let u=f[G]||[],$=new Map(c.map(F=>[F.address,F])),q=new Map(u.map(F=>[F.address,F])),D=0,M=0,I=0;u.forEach(F=>{let se=$.get(F.address);if(!se){D+=1;return}(Number(se.weight)!==Number(F.weight)||Number(se.flags||0)!==Number(F.flags||0))&&(M+=1)}),c.forEach(F=>{q.has(F.address)||(I+=1)}),Ie({add:D,update:M,remove:I})},[G,c,f]);let je=async u=>{try{let $=Number(w[u.address]);await N.post("/vips/reals",{vip:o,real:{address:u.address,weight:$,flags:u.flags||0}}),await z(),e("Real weight updated.","success")}catch($){l($.message||"request failed"),e($.message||"Update failed.","error")}},He=async u=>{try{await N.del("/vips/reals",{vip:o,real:{address:u.address,weight:u.weight,flags:u.flags||0}}),await z(),e("Real removed.","success")}catch($){l($.message||"request failed"),e($.message||"Remove failed.","error")}},We=async u=>{u.preventDefault();try{await N.post("/vips/reals",{vip:o,real:{address:g.address,weight:Number(g.weight),flags:Number(g.flags||0)}}),y({address:"",weight:100,flags:0}),await z(),e("Real added.","success")}catch($){l($.message||"request failed"),e($.message||"Add failed.","error")}},Je=async()=>{if(!G||!f[G]){O("Select a target group to apply.");return}ye(!0),O("");let u=f[G]||[],$=new Map(c.map(I=>[I.address,I])),q=new Map(u.map(I=>[I.address,I])),D=c.filter(I=>!q.has(I.address)),M=u.filter(I=>{let F=$.get(I.address);return F?Number(F.weight)!==Number(I.weight)||Number(F.flags||0)!==Number(I.flags||0):!0});try{D.length>0&&await N.put("/vips/reals/batch",{vip:o,action:1,reals:D.map(I=>({address:I.address,weight:Number(I.weight),flags:Number(I.flags||0)}))}),M.length>0&&await Promise.all(M.map(I=>N.post("/vips/reals",{vip:o,real:{address:I.address,weight:Number(I.weight),flags:Number(I.flags||0)}}))),await z(),e(`Applied target group "${G}".`,"success")}catch(I){O(I.message||"Failed to apply target group."),e(I.message||"Target group apply failed.","error")}finally{ye(!1)}},Ke=u=>{u.preventDefault();let $=we.trim();if(!$){O("Provide a name for the new target group.");return}if(f[$]){O("A target group with that name already exists.");return}let q={...f,[$]:c.map(D=>({address:D.address,weight:Number(D.weight),flags:Number(D.flags||0)}))};S(q),_e(""),he($),O(""),e(`Target group "${$}" saved.`,"success")},Ye=async()=>{try{await N.del("/vips",o),e("VIP deleted.","success"),a("/")}catch(u){l(u.message||"request failed"),e(u.message||"Delete failed.","error")}},Ze=async u=>{u.preventDefault();try{await N.put("/vips/flags",{...o,flag:Number(R.flags||0),set:!!R.set}),await ae(),e("VIP flags updated.","success")}catch($){d($.message||"request failed"),e($.message||"Flag update failed.","error")}},Xe=async u=>{u.preventDefault();try{await N.put("/vips/hash-function",{...o,hash_function:Number(n.hash_function)}),e("Hash function updated.","success")}catch($){d($.message||"request failed"),e($.message||"Hash update failed.","error")}};return i`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>VIP Detail</h2>
              <p className="muted">${o.address}:${o.port} / ${o.proto}</p>
              ${_===null?i`<p className="muted">Flags: —</p>`:i`
                    <div style=${{marginTop:8}}>
                      <${ie}
                        mask=${_}
                        options=${H}
                        showStatus=${!0}
                        emptyLabel="No flags"
                      />
                    </div>
                  `}
            </div>
            <div className="row">
              <button className="btn ghost" onClick=${ae}>Refresh flags</button>
              <button className="btn danger" onClick=${Ye}>Delete VIP</button>
            </div>
          </div>
          ${r&&i`<p className="error">${r}</p>`}
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
                              onInput=${$=>T({...w,[u.address]:$.target.value})}
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
                  <${ce}
                    options=${H}
                    value=${R.flags}
                    name="vip-flag-change"
                    onChange=${u=>x({...R,flags:u})}
                  />
                </label>
                <label className="field">
                  <span>Set</span>
                  <select
                    value=${String(R.set)}
                    onChange=${u=>x({...R,set:u.target.value==="true"})}
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
                    value=${n.hash_function}
                    onInput=${u=>v({...n,hash_function:u.target.value})}
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
                onClick=${async()=>{try{await k(),e("Imported target groups from running config.","success")}catch(u){O(u.message||"Failed to import target groups."),e(u.message||"Import failed.","error")}}}
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
                value=${G}
                onChange=${u=>he(u.target.value)}
                disabled=${Object.keys(f).length===0}
              >
                <option value="">Select group</option>
                ${Object.keys(f).map(u=>i`<option value=${u}>${u}</option>`)}
              </select>
            </label>
            <label className="field">
              <span>Preview</span>
              <input
                value=${`add ${te.add} \xB7 update ${te.update} \xB7 remove ${te.remove}`}
                readOnly
              />
            </label>
          </div>
          <div className="row">
            <button
              className="btn"
              type="button"
              onClick=${Je}
              disabled=${Ne||!G}
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
    `}function Le(){let e=oe(),t=E(()=>Q(e.vipId),[e.vipId]),{points:a,error:o}=K({path:"/stats/vip",body:t}),c=E(()=>Y("/stats/vip"),[]),b=a[a.length-1]||{},r=a[a.length-2]||{},l=Number(b.v1??0),p=Number(b.v2??0),d=l-Number(r.v1??0),m=p-Number(r.v2??0),s=E(()=>[{label:c.v1,field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.15)"},{label:c.v2,field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[c]);return i`
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
          <${ee} title="Traffic (delta/sec)" points=${a} keys=${s} diff=${!0} />
        </section>
      </main>
    `}let De={"/stats/vip":{v1:"Packets",v2:"Bytes"},"/stats/real":{v1:"Packets",v2:"Bytes"},"/stats/lru":{v1:"Total packets",v2:"LRU hits"},"/stats/lru/miss":{v1:"TCP SYN misses",v2:"Non-SYN misses"},"/stats/lru/fallback":{v1:"Fallback LRU hits",v2:"Unused"},"/stats/lru/global":{v1:"Map lookup failures",v2:"Global LRU routed"},"/stats/xdp/total":{v1:"Packets",v2:"Bytes"},"/stats/xdp/pass":{v1:"Packets",v2:"Bytes"},"/stats/xdp/drop":{v1:"Packets",v2:"Bytes"},"/stats/xdp/tx":{v1:"Packets",v2:"Bytes"}},ve=[{title:"LRU",path:"/stats/lru"},{title:"LRU Miss",path:"/stats/lru/miss"},{title:"LRU Fallback",path:"/stats/lru/fallback"},{title:"LRU Global",path:"/stats/lru/global"},{title:"XDP Total",path:"/stats/xdp/total"},{title:"XDP Pass",path:"/stats/xdp/pass"},{title:"XDP Drop",path:"/stats/xdp/drop"},{title:"XDP Tx",path:"/stats/xdp/tx"}];function V(e){return Number.isFinite(e)?`${e>0?"+":""}${e}`:"0"}function Y(e){return De[e]||{v1:"v1",v2:"v2"}}function qe({title:e,path:t,diff:a=!1}){let{points:o,error:c}=K({path:t}),b=E(()=>Y(t),[t]),r=E(()=>[{label:b.v1,field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.2)"},{label:b.v2,field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[b]);return i`
      <div className="card">
        <h3>${e}</h3>
        ${c&&i`<p className="error">${c}</p>`}
        <${ee} title=${e} points=${o} keys=${r} diff=${a} inlineTitle=${!1} />
      </div>
    `}function Be({title:e,path:t}){let{points:a,error:o}=K({path:t}),c=E(()=>Y(t),[t]),b=a[a.length-1]||{},r=a[a.length-2]||{},l=Number(b.v1??0),p=Number(b.v2??0),d=l-Number(r.v1??0),m=p-Number(r.v2??0);return i`
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
    `}function Oe(){let{data:e,error:t}=xe(()=>N.get("/stats/userspace"),1e3,[]);return i`
      <main>
        <section className="card">
          <h2>Global Stats</h2>
          <p className="muted">Polling every second. Charts show per-second deltas.</p>
        </section>
        <section className="grid">
          ${ve.map(a=>i`<${qe} title=${a.title} path=${a.path} diff=${!0} />`)}
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
            ${ve.map(a=>i`<${Be} title=${a.title} path=${a.path} />`)}
          </div>
        </section>
      </main>
    `}function Ve(){let[e,t]=h([]),[a,o]=h(""),[c,b]=h([]),[r,l]=h(""),[p,d]=h(null),[m,s]=h("");A(()=>{let f=!0;return(async()=>{try{let C=await N.get("/vips");if(!f)return;t(C||[]),!a&&C&&C.length>0&&o(J(C[0]))}catch(C){f&&s(C.message||"request failed")}})(),()=>{f=!1}},[]),A(()=>{if(!a)return;let f=Q(a),S=!0;return(async()=>{try{let k=await N.get("/vips/reals",f);if(!S)return;b(k||[]),k&&k.length>0?l(G=>G||k[0].address):l(""),s("")}catch(k){S&&s(k.message||"request failed")}})(),()=>{S=!1}},[a]),A(()=>{if(!r){d(null);return}let f=!0;return(async()=>{try{let C=await N.get("/reals/index",{address:r});if(!f)return;d(C?.index??null),s("")}catch(C){f&&s(C.message||"request failed")}})(),()=>{f=!1}},[r]);let{points:g,error:y}=K({path:"/stats/real",body:p!==null?{index:p}:null}),w=E(()=>Y("/stats/real"),[]),T=E(()=>[{label:w.v1,field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.2)"},{label:w.v2,field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[w]),_=g[g.length-1]||{},P=g[g.length-2]||{},R=Number(_.v1??0),x=Number(_.v2??0),n=R-Number(P.v1??0),v=x-Number(P.v2??0);return i`
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
                value=${r}
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
                      <td>${R}</td>
                      <td>
                        <span className=${`delta ${n<0?"down":"up"}`}>
                          ${V(n)}
                        </span>
                      </td>
                    </tr>
                    <tr>
                      <td>${w.v2}</td>
                      <td>${x}</td>
                      <td>
                        <span className=${`delta ${v<0?"down":"up"}`}>
                          ${V(v)}
                        </span>
                      </td>
                    </tr>
                  </tbody>
                </table>
                <${ee}
                  title="Traffic (delta/sec)"
                  points=${g}
                  keys=${T}
                  diff=${!0}
                />
              `}
        </section>
      </main>
    `}function Ue(){let{addToast:e}=W(),[t,a]=h(""),[o,c]=h(""),[b,r]=h(!0),[l,p]=h(""),d=U(!0),m=async()=>{if(d.current){r(!0),c("");try{let g=await fetch(`${N.base}/config/export`,{headers:{Accept:"application/x-yaml"}});if(!g.ok){let w=`HTTP ${g.status}`;try{w=(await g.json())?.error?.message||w}catch{}throw new Error(w)}let y=await g.text();if(!d.current)return;a(y||""),p(new Date().toLocaleString())}catch(g){d.current&&c(g.message||"request failed")}finally{d.current&&r(!1)}}},s=async()=>{if(t)try{await navigator.clipboard.writeText(t),e("Config copied to clipboard","info")}catch{e("Failed to copy config","error")}};return A(()=>(d.current=!0,m(),()=>{d.current=!1}),[]),i`
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
              <button className="btn" onClick=${m} disabled=${b}>
                Refresh
              </button>
            </div>
          </div>
          ${o&&i`<p className="error">${o}</p>`}
          ${b?i`<p className="muted">Loading config...</p>`:t?i`<pre className="yaml-view">${t}</pre>`:i`<p className="muted">No config data returned.</p>`}
          ${l&&i`<p className="muted">Last fetched ${l}</p>`}
        </section>
      </main>
    `}function ze(){let{addToast:e}=W(),{groups:t,setGroups:a,refreshFromStorage:o,importFromRunningConfig:c}=ge(),[b,r]=h(""),[l,p]=h(""),[d,m]=h({address:"",weight:100,flags:0}),[s,g]=h(""),[y,w]=h(!1);A(()=>{if(l){if(!t[l]){let v=Object.keys(t);p(v[0]||"")}}else{let v=Object.keys(t);v.length>0&&p(v[0])}},[t,l]);let T=v=>{v.preventDefault();let f=b.trim();if(!f){g("Provide a group name.");return}if(t[f]){g("That group already exists.");return}a({...t,[f]:[]}),r(""),p(f),g(""),e(`Target group "${f}" created.`,"success")},_=v=>{let f={...t};delete f[v],a(f),e(`Target group "${v}" removed.`,"success")},P=v=>{if(v.preventDefault(),!l){g("Select a group to add a real.");return}let f=ue(d);if(!f){g("Provide a valid real address.");return}let S=t[l]||[],C=S.some(k=>k.address===f.address)?S.map(k=>k.address===f.address?f:k):S.concat(f);a({...t,[l]:C}),m({address:"",weight:100,flags:0}),g(""),e("Real saved to target group.","success")},R=v=>{if(!l)return;let S=(t[l]||[]).filter(C=>C.address!==v);a({...t,[l]:S})},x=(v,f)=>{if(!l)return;let C=(t[l]||[]).map(k=>k.address===v?{...k,...f}:k);a({...t,[l]:C})};return i`
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
              <button className="btn ghost" type="button" onClick=${async()=>{w(!0);try{await c(),e("Imported target groups from running config.","success"),g("")}catch(v){g(v.message||"Failed to import target groups."),e(v.message||"Import failed.","error")}finally{w(!1)}}} disabled=${y}>
                ${y?"Importing...":"Import from running config"}
              </button>
            </div>
          </div>
          ${s&&i`<p className="error">${s}</p>`}
          <form className="form" onSubmit=${T}>
            <div className="form-row">
              <label className="field">
                <span>New group name</span>
                <input
                  value=${b}
                  onInput=${v=>r(v.target.value)}
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
                  onChange=${v=>p(v.target.value)}
                  disabled=${Object.keys(t).length===0}
                >
                  ${Object.keys(t).map(v=>i`<option value=${v}>${v}</option>`)}
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
                    ${(t[l]||[]).map(v=>i`
                        <tr>
                          <td>${v.address}</td>
                          <td>
                            <input
                              className="inline-input"
                              type="number"
                              min="0"
                              value=${v.weight}
                              onInput=${f=>x(v.address,{weight:Number(f.target.value)})}
                            />
                          </td>
                          <td>
                            <button
                              className="btn ghost"
                              type="button"
                              onClick=${()=>R(v.address)}
                            >
                              Remove
                            </button>
                          </td>
                        </tr>
                      `)}
                  </tbody>
                </table>
                <form className="form" onSubmit=${P}>
                  <div className="form-row">
                    <label className="field">
                      <span>Real address</span>
                      <input
                        value=${d.address}
                        onInput=${v=>m({...d,address:v.target.value})}
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
                        onInput=${v=>m({...d,weight:v.target.value})}
                        required
                      />
                    </label>
                  </div>
                  <button className="btn secondary" type="submit">Save real</button>
                </form>
              `:i`<p className="muted">No groups yet. Create one to add reals.</p>`}
        </section>
      </main>
    `}function Me(){let[e,t]=h({initialized:!1,ready:!1}),[a,o]=h([]),c=U({}),b=(l,p="info")=>{let d=`${Date.now()}-${Math.random().toString(16).slice(2)}`;o(m=>m.concat({id:d,message:l,kind:p})),c.current[d]=setTimeout(()=>{o(m=>m.filter(s=>s.id!==d)),delete c.current[d]},4e3)},r=l=>{c.current[l]&&(clearTimeout(c.current[l]),delete c.current[l]),o(p=>p.filter(d=>d.id!==l))};return A(()=>{let l=!0,p=async()=>{try{let m=await N.get("/lb/status");l&&t(m||{initialized:!1,ready:!1})}catch{l&&t({initialized:!1,ready:!1})}};p();let d=setInterval(p,5e3);return()=>{l=!1,clearInterval(d)}},[]),i`
      <${re}>
        <${be}>
          <${Z.Provider} value=${{addToast:b}}>
            <${Ae} status=${e} />
            <${ne}>
              <${B} path="/" element=${i`<${Ge} />`} />
              <${B} path="/vips/:vipId" element=${i`<${Ee} />`} />
              <${B} path="/vips/:vipId/stats" element=${i`<${Le} />`} />
              <${B} path="/target-groups" element=${i`<${ze} />`} />
              <${B} path="/stats/global" element=${i`<${Oe} />`} />
              <${B} path="/stats/real" element=${i`<${Ve} />`} />
              <${B} path="/config" element=${i`<${Ue} />`} />
            </${ne}>
            <${Fe} toasts=${a} onDismiss=${r} />
          </${Z.Provider}>
        </${be}>
      </${re}>
    `}ReactDOM.createRoot(document.getElementById("root")).render(i`<${Me} />`)})();})();
