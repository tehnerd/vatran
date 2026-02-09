(()=>{(()=>{let{useEffect:x,useMemo:E,useRef:V,useState:v,useContext:le,useCallback:Re}=React,{BrowserRouter:ie,Routes:ce,Route:D,NavLink:G,Link:M,useParams:de,useNavigate:Te}=ReactRouterDOM,o=htm.bind(React.createElement),Q=React.createContext({addToast:()=>{}}),ee=React.createContext({required:!1,username:"",login:async()=>{},logout:async()=>{}}),W=[{label:"NO_SPORT",value:1},{label:"NO_LRU",value:2},{label:"QUIC_VIP",value:4},{label:"DPORT_HASH",value:8},{label:"SRC_ROUTING",value:16},{label:"LOCAL_VIP",value:32},{label:"GLOBAL_LRU",value:64},{label:"HASH_SRC_DST_PORT",value:128},{label:"UDP_STABLE_ROUTING_VIP",value:256},{label:"UDP_FLOW_MIGRATION",value:512}];function J(){return le(Q)}function Pe(){return le(ee)}class O extends Error{constructor(t,a=0,s=""){super(t||"request failed"),this.name="ApiError",this.status=a,this.code=s,this.unauthorized=a===401||s==="UNAUTHORIZED"}}function Ae(e){return String(e).replace(/[^a-z0-9_-]/gi,"_")}function xe(e,t,a){let s=Number(e)||0,i=Number(t)||0;return a?s|i:s&~i}function Fe(e,t){let a=Number(e)||0;return t.filter(s=>(a&s.value)!==0)}function ue(e,t){let a=String(e??"").trim();if(!a)return;let i=a.split(/[\s,]+/).filter(Boolean).map(r=>Number(r));if(i.findIndex(r=>!Number.isFinite(r)||!Number.isInteger(r))!==-1)throw new Error(`${t} must be a comma- or space-separated list of integers.`);return i}function pe({mask:e,options:t,showStatus:a=!1,emptyLabel:s="None"}){let i=Number(e)||0,p=a?t:Fe(i,t),r=a?2:1;return o`
      <table className="table flag-table">
        <thead>
          <tr>
            <th>Flag</th>
            ${a?o`<th>Enabled</th>`:null}
          </tr>
        </thead>
        <tbody>
          ${p.length===0?o`<tr><td colspan=${r} className="muted">${s}</td></tr>`:p.map(c=>{let f=(i&c.value)!==0;return o`
                  <tr>
                    <td>${c.label}</td>
                    ${a?o`<td>${f?"Yes":"No"}</td>`:null}
                  </tr>
                `})}
        </tbody>
      </table>
    `}function me({options:e,value:t,onChange:a,name:s}){let i=Number(t)||0,p=Ae(s||"flags");return o`
      <div className="flag-selector">
        <table className="table flag-table">
          <thead>
            <tr>
              <th>Enabled</th>
              <th>Flag</th>
            </tr>
          </thead>
          <tbody>
            ${e.map(r=>{let c=`${p}-${r.value}`,f=(i&r.value)===r.value;return o`
                <tr>
                  <td>
                    <input
                      id=${c}
                      type="checkbox"
                      checked=${f}
                      onChange=${d=>a(xe(i,r.value,d.target.checked))}
                    />
                  </td>
                  <td>
                    <label className="flag-option" htmlFor=${c}>${r.label}</label>
                  </td>
                </tr>
              `})}
          </tbody>
        </table>
      </div>
    `}let N={base:"/api/v1",authHandlers:{onUnauthorized:null},setAuthHandlers(e={}){N.authHandlers={onUnauthorized:e.onUnauthorized||null}},notifyUnauthorized(e){typeof N.authHandlers.onUnauthorized=="function"&&N.authHandlers.onUnauthorized(e||"Authentication required")},async parseJSON(e){try{return await e.json()}catch{return null}},async buildError(e,t){let a=t?.error?.code||"",s=t?.error?.message||`HTTP ${e.status}`;return(e.status===401||a==="UNAUTHORIZED")&&N.notifyUnauthorized(s),new O(s,e.status,a)},async request(e,t={}){let a={method:t.method||"GET",credentials:"same-origin",headers:{Accept:"application/json",...t.headers||{}}},s=`${N.base}${e}`;if(t.body!==void 0&&t.body!==null)if(a.method==="GET"){let r=new URLSearchParams;Object.entries(t.body).forEach(([f,d])=>{if(d!=null){if(Array.isArray(d)){d.forEach(g=>r.append(f,String(g)));return}if(typeof d=="object"){r.set(f,JSON.stringify(d));return}r.set(f,String(d))}});let c=r.toString();c&&(s+=`${s.includes("?")?"&":"?"}${c}`)}else a.headers["Content-Type"]="application/json",a.body=JSON.stringify(t.body);let i=await fetch(s,a),p=await N.parseJSON(i);if(!p)throw new O("invalid JSON response",i.status);if(!i.ok)throw await N.buildError(i,p);if(!p.success)throw new O(p.error?.message||"request failed",i.status,p.error?.code);return p.data},get(e,t){return N.request(e,{method:"GET",body:t})},post(e,t){return N.request(e,{method:"POST",body:t})},put(e,t){return N.request(e,{method:"PUT",body:t})},del(e,t){return N.request(e,{method:"DELETE",body:t})},async login(e,t){let a=await fetch("/login",{method:"POST",credentials:"same-origin",headers:{Accept:"application/json","Content-Type":"application/json"},body:JSON.stringify({username:e,password:t})}),s=await N.parseJSON(a);if(!s)throw new O("invalid JSON response",a.status);if(!a.ok)throw await N.buildError(a,s);if(!s.success)throw new O(s.error?.message||"login failed",a.status,s.error?.code);return s.data||{}},async logout(){let e=await fetch("/logout",{method:"POST",credentials:"same-origin",headers:{Accept:"application/json"}}),t=await N.parseJSON(e);if(!e.ok)throw await N.buildError(e,t||{});if(t&&t.success===!1)throw new O(t.error?.message||"logout failed",e.status,t.error?.code);return t?.data||{}}},ge="vatran_target_groups";function fe(e){if(!e||!e.address)return null;let t=String(e.address).trim();if(!t)return null;let a=Number(e.weight),s=Number(e.flags??0);return{address:t,weight:Number.isFinite(a)?a:0,flags:Number.isFinite(s)?s:0}}function he(e){if(!e||typeof e!="object")return{};let t={};return Object.entries(e).forEach(([a,s])=>{let i=String(a).trim();if(!i)return;let p=Array.isArray(s)?s.map(fe).filter(Boolean):[],r=[],c=new Set;p.forEach(f=>{c.has(f.address)||(c.add(f.address),r.push(f))}),t[i]=r}),t}function te(){if(typeof localStorage>"u")return{};try{let e=localStorage.getItem(ge);return e?he(JSON.parse(e)):{}}catch{return{}}}function be(e){if(!(typeof localStorage>"u"))try{localStorage.setItem(ge,JSON.stringify(e))}catch{}}function qe(e,t){let a={...e};return Object.entries(t||{}).forEach(([s,i])=>{a[s]||(a[s]=i)}),a}function ve(){let[e,t]=v(()=>te());return x(()=>{be(e)},[e]),{groups:e,setGroups:t,refreshFromStorage:()=>{t(te())},importFromRunningConfig:async()=>{let i=await N.get("/config/export/json"),p=he(i?.target_groups||{}),r=qe(te(),p);return t(r),be(r),r}}}function Z(e){return`${encodeURIComponent(e.address)}:${e.port}:${e.proto}`}function ae(e){let t=e.split(":"),a=Number(t.pop()||0),s=Number(t.pop()||0);return{address:decodeURIComponent(t.join(":")),port:s,proto:a}}function Ee(e,t,a=[]){let[s,i]=v(null),[p,r]=v(""),[c,f]=v(!0);return x(()=>{let d=!0,g=async()=>{try{let m=await e();d&&(i(m),r(""),f(!1))}catch(m){d&&(r(m.message||"request failed"),f(!1))}};g();let n=setInterval(g,t);return()=>{d=!1,clearInterval(n)}},a),{data:s,error:p,loading:c}}function K({path:e,body:t,intervalMs:a=1e3,limit:s=60}){let[i,p]=v([]),[r,c]=v(""),f=E(()=>JSON.stringify(t||{}),[t]);return x(()=>{if(t===null)return p([]),c(""),()=>{};let d=!0,g=async()=>{try{let m=await N.get(e,t);if(!d)return;let b=new Date().toLocaleTimeString();p(y=>y.concat({label:b,...m}).slice(-s)),c("")}catch(m){d&&c(m.message||"request failed")}};g();let n=setInterval(g,a);return()=>{d=!1,clearInterval(n)}},[e,f,a,s]),{points:i,error:r}}function $e({title:e,points:t,keys:a,diff:s=!1,height:i=120,showTitle:p=!1,selectedLabel:r=null,onPointSelect:c=null,onLegendSelect:f=null}){let d=V(null),g=V(null);return x(()=>{if(!d.current)return;g.current||(g.current=new Chart(d.current,{type:"line",data:{labels:[],datasets:[]},options:{responsive:!0,maintainAspectRatio:!1,animation:!1,scales:{x:{grid:{display:!1}},y:{beginAtZero:!s}},plugins:{legend:{display:!0,position:"bottom"},title:{display:p&&!!e,text:e}}}}));let n=g.current,m=new Map((n.data.datasets||[]).filter(w=>typeof w.hidden<"u").map(w=>[w.label,w.hidden])),b=t.map(w=>w.label),y=r?a.filter(w=>w.label===r):a,_=r&&y.length===0?a:y;return n.data.labels=b,n.data.datasets=_.map(w=>{let T=t.map(F=>F[w.field]||0),P=s?T.map((F,l)=>l===0?0:F-T[l-1]):T;return{label:w.label,data:P,borderColor:w.color,backgroundColor:w.fill,borderWidth:2,tension:.3,hidden:m.get(w.label)}}),n.options.onClick=(w,T)=>{if(!c||!T||T.length===0)return;let P=T[0].datasetIndex,F=n.data.datasets?.[P]?.label;F&&c(F)},n.options.plugins&&n.options.plugins.legend&&(n.options.plugins.legend.onClick=(w,T)=>{if(!f)return;let P=T?.text;P&&f(P)}),n.options.scales.y.beginAtZero=!s,n.options.plugins.title.display=p&&!!e,n.options.plugins.title.text=e||"",n.update(),()=>{}},[t,a,e,s,p,r,c,f]),x(()=>()=>{g.current&&(g.current.destroy(),g.current=null)},[]),o`<canvas ref=${d} height=${i}></canvas>`}function se({title:e,points:t,keys:a,diff:s=!1,inlineTitle:i=!0}){let[p,r]=v(!1),[c,f]=v(null),d=V(!1);return o`
      <div className="chart-wrap">
        <button className="chart-zoom-button" type="button" onClick=${()=>r(!0)}>
          <span className="zoom-icon" aria-hidden="true">+</span>
          Zoom
        </button>
        <div
          className="chart-click"
          onClick=${()=>{if(d.current){d.current=!1;return}r(!0)}}
        >
          <${$e}
            title=${e}
            points=${t}
            keys=${a}
            diff=${s}
            height=${120}
            showTitle=${i&&!!e}
            selectedLabel=${c}
            onPointSelect=${g=>{f(n=>n===g?null:g),d.current=!0,setTimeout(()=>{d.current=!1},0)}}
            onLegendSelect=${g=>{f(n=>n===g?null:g),d.current=!0,setTimeout(()=>{d.current=!1},0)}}
          />
        </div>
        ${p&&o`
          <div className="chart-overlay" onClick=${()=>r(!1)}>
            <div className="chart-modal" onClick=${g=>g.stopPropagation()}>
              <div className="row chart-modal-header">
                <div>
                  <h3>${e||"Chart"}</h3>
                  ${s?o`<p className="muted">Per-second delta.</p>`:""}
                </div>
                <button className="btn ghost" onClick=${()=>r(!1)}>
                  Close
                </button>
              </div>
              <div className="chart-zoom">
                <${$e}
                  title=${e}
                  points=${t}
                  keys=${a}
                  diff=${s}
                  height=${360}
                  showTitle=${!1}
                  selectedLabel=${c}
                  onPointSelect=${g=>f(n=>n===g?null:g)}
                  onLegendSelect=${g=>f(n=>n===g?null:g)}
                />
              </div>
            </div>
          </div>
        `}
      </div>
    `}function Ge(){let{login:e}=Pe(),[t,a]=v(""),[s,i]=v(""),[p,r]=v(""),[c,f]=v(!1),d=async g=>{if(g.preventDefault(),!t.trim()||!s){r("Username and password are required.");return}f(!0),r("");try{await e(t.trim(),s)}catch(n){r(n.message||"Login failed.")}finally{f(!1)}};return o`
      <main className="auth-main">
        <section className="auth-card">
          <h1>Sign in</h1>
          <p className="muted">Authentication is required to access Vatran.</p>
          ${p?o`<p className="error">${p}</p>`:null}
          <form className="form" onSubmit=${d}>
            <label className="field">
              <span>Username</span>
              <input
                value=${t}
                onInput=${g=>a(g.target.value)}
                autoComplete="username"
                required
              />
            </label>
            <label className="field">
              <span>Password</span>
              <input
                type="password"
                value=${s}
                onInput=${g=>i(g.target.value)}
                autoComplete="current-password"
                required
              />
            </label>
            <button className="btn" type="submit" disabled=${c}>
              ${c?"Signing in...":"Sign in"}
            </button>
          </form>
        </section>
      </main>
    `}function Ne({checking:e,required:t,children:a}){return e?o`
        <main className="auth-main">
          <section className="auth-card">
            <p className="muted">Checking authentication...</p>
          </section>
        </main>
      `:t?o`<${Ge} />`:a}function Le({toasts:e,onDismiss:t}){return o`
      <div className="toast-stack">
        ${e.map(a=>o`
            <div className=${`toast ${a.kind}`}>
              <span>${a.message}</span>
              <button className="toast-close" onClick=${()=>t(a.id)}>
                ×
              </button>
            </div>
          `)}
      </div>
    `}function Oe({status:e}){return o`
      <header>
        <div>
          <div style=${{fontSize:20,fontWeight:700}}>Vatran</div>
          <div className="status-pill">
            <span className=${`dot ${e.ready?"ok":""}`}></span>
            ${e.ready?"Ready":"Not ready"}
          </div>
        </div>
        <nav>
          <${G} to="/" end className=${({isActive:t})=>t?"active":""}>
            Dashboard
          </${G}>
          <${G} to="/stats/global" className=${({isActive:t})=>t?"active":""}>
            Global stats
          </${G}>
          <${G} to="/stats/real" className=${({isActive:t})=>t?"active":""}>
            Per-real stats
          </${G}>
          <${G}
            to="/target-groups"
            className=${({isActive:t})=>t?"active":""}
          >
            Target groups
          </${G}>
          <${G} to="/config" className=${({isActive:t})=>t?"active":""}>
            Config export
          </${G}>
        </nav>
      </header>
    `}function Ue(){let{addToast:e}=J(),[t,a]=v({initialized:!1,ready:!1}),[s,i]=v([]),[p,r]=v(""),[c,f]=v(!1),[d,g]=v(!1),[n,m]=v({main_interface:"",balancer_prog_path:"",healthchecking_prog_path:"",default_mac:"",local_mac:"",root_map_path:"",root_map_pos:2,katran_src_v4:"",katran_src_v6:"",use_root_map:!1,max_vips:512,max_reals:4096,hash_function:"maglev_v2",forwarding_cores:"",numa_nodes:""}),[b,y]=v({address:"",port:80,proto:6,flags:0}),_=async()=>{try{let l=await N.get("/lb/status"),$=await N.get("/vips"),h=await Promise.all(($||[]).map(async C=>{let k=null,R=!1;try{k=(await N.get("/vips/flags",{address:C.address,port:C.port,proto:C.proto}))?.flags??0}catch{k=null}try{let A=await N.get("/vips/reals",{address:C.address,port:C.port,proto:C.proto});R=Array.isArray(A)&&A.some(X=>!!X?.healthy)}catch{R=!1}return{...C,flags:k,healthy:R}}));a(l||{initialized:!1,ready:!1}),i(h),r("")}catch(l){r(l.message||"request failed")}};x(()=>{let l=!0;return(async()=>{l&&await _()})(),()=>{l=!1}},[]);let w=async l=>{l.preventDefault();try{let $=ue(n.forwarding_cores,"Forwarding cores"),h=ue(n.numa_nodes,"NUMA nodes"),C={...n,forwarding_cores:$,numa_nodes:h,root_map_pos:n.root_map_pos===""?void 0:Number(n.root_map_pos),max_vips:Number(n.max_vips),max_reals:Number(n.max_reals),hash_function:n.hash_function};await N.post("/lb/create",C),r(""),f(!1),e("Load balancer initialized.","success"),await _()}catch($){r($.message||"request failed"),e($.message||"Initialize failed.","error")}},T=async l=>{l.preventDefault();try{await N.post("/vips",{...b,port:Number(b.port),proto:Number(b.proto),flags:Number(b.flags||0)}),y({address:"",port:80,proto:6,flags:0}),r(""),g(!1),e("VIP created.","success"),await _()}catch($){r($.message||"request failed"),e($.message||"VIP create failed.","error")}},P=async()=>{try{await N.post("/lb/load-bpf-progs"),r(""),e("BPF programs loaded.","success"),await _()}catch(l){r(l.message||"request failed"),e(l.message||"Load BPF programs failed.","error")}},F=async()=>{try{await N.post("/lb/attach-bpf-progs"),r(""),e("BPF programs attached.","success"),await _()}catch(l){r(l.message||"request failed"),e(l.message||"Attach BPF programs failed.","error")}};return o`
      <main>
        <section className="card">
          <h2>Load balancer</h2>
          <p>Initialized: ${t.initialized?"yes":"no"}</p>
          <p>Ready: ${t.ready?"yes":"no"}</p>
          <div className="row">
            ${!t.initialized&&o`
              <button className="btn" onClick=${()=>f(l=>!l)}>
                ${c?"Close":"Initialize"}
              </button>
            `}
            <button className="btn secondary" onClick=${()=>g(l=>!l)}>
              ${d?"Close":"Create VIP"}
            </button>
          </div>
          ${!t.ready&&o`
            <div className="row" style=${{marginTop:12}}>
              <button
                className="btn ghost"
                disabled=${!t.initialized}
                onClick=${P}
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
          ${c&&o`
            <form className="form" onSubmit=${w}>
              <div className="form-row">
                <label className="field">
                  <span>Main interface</span>
                  <input
                    value=${n.main_interface}
                    onInput=${l=>m({...n,main_interface:l.target.value})}
                    placeholder="eth0"
                    required
                  />
                </label>
                <label className="field">
                  <span>Balancer prog path</span>
                  <input
                    value=${n.balancer_prog_path}
                    onInput=${l=>m({...n,balancer_prog_path:l.target.value})}
                    placeholder="/path/to/bpf.o"
                    required
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Healthchecking path</span>
                  <input
                    value=${n.healthchecking_prog_path}
                    onInput=${l=>m({...n,healthchecking_prog_path:l.target.value})}
                    placeholder="/path/to/hc.o"
                  />
                </label>
                <label className="field">
                  <span>Default MAC</span>
                  <input
                    value=${n.default_mac}
                    onInput=${l=>m({...n,default_mac:l.target.value})}
                    placeholder="aa:bb:cc:dd:ee:ff"
                    required
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Local MAC</span>
                  <input
                    value=${n.local_mac}
                    onInput=${l=>m({...n,local_mac:l.target.value})}
                    placeholder="11:22:33:44:55:66"
                    required=${n.healthchecking_prog_path?.trim()!==""}
                  />
                </label>
                <label className="field">
                  <span>Hash function</span>
                  <select
                    value=${n.hash_function}
                    onInput=${l=>m({...n,hash_function:l.target.value})}
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
                    value=${n.root_map_path}
                    onInput=${l=>m({...n,root_map_path:l.target.value})}
                    placeholder="/sys/fs/bpf/root_map"
                  />
                </label>
                <label className="field">
                  <span>Root map position</span>
                  <input
                    type="number"
                    min="0"
                    value=${n.root_map_pos}
                    onInput=${l=>m({...n,root_map_pos:l.target.value})}
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Katran src v4</span>
                  <input
                    value=${n.katran_src_v4}
                    onInput=${l=>m({...n,katran_src_v4:l.target.value})}
                    placeholder="10.0.0.1"
                  />
                </label>
                <label className="field">
                  <span>Katran src v6</span>
                  <input
                    value=${n.katran_src_v6}
                    onInput=${l=>m({...n,katran_src_v6:l.target.value})}
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
                    value=${n.max_vips}
                    onInput=${l=>m({...n,max_vips:l.target.value})}
                  />
                </label>
                <label className="field">
                  <span>Max reals</span>
                  <input
                    type="number"
                    min="1"
                    value=${n.max_reals}
                    onInput=${l=>m({...n,max_reals:l.target.value})}
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Forwarding cores (optional)</span>
                  <input
                    value=${n.forwarding_cores}
                    onInput=${l=>m({...n,forwarding_cores:l.target.value})}
                    placeholder="0,1,2,3"
                  />
                  <span className="muted">Comma or space separated CPU core IDs.</span>
                </label>
                <label className="field">
                  <span>NUMA nodes (optional)</span>
                  <input
                    value=${n.numa_nodes}
                    onInput=${l=>m({...n,numa_nodes:l.target.value})}
                    placeholder="0,0,1,1"
                  />
                  <span className="muted">Match the forwarding cores length.</span>
                </label>
              </div>
              <label className="field checkbox">
                <input
                  type="checkbox"
                  checked=${n.use_root_map}
                  onChange=${l=>m({...n,use_root_map:l.target.checked})}
                />
                <span>Use root map</span>
              </label>
              <button className="btn" type="submit">Initialize LB</button>
            </form>
          `}
          ${d&&o`
            <form className="form" onSubmit=${T}>
              <div className="form-row">
                <label className="field">
                  <span>VIP address</span>
                  <input
                    value=${b.address}
                    onInput=${l=>y({...b,address:l.target.value})}
                    placeholder="1.2.3.4"
                    required
                  />
                </label>
                <label className="field">
                  <span>Port</span>
                  <input
                    type="number"
                    value=${b.port}
                    onInput=${l=>y({...b,port:l.target.value})}
                    required
                  />
                </label>
                <label className="field">
                  <span>Protocol</span>
                  <select
                    value=${b.proto}
                    onChange=${l=>y({...b,proto:l.target.value})}
                  >
                    <option value="6">TCP (6)</option>
                    <option value="17">UDP (17)</option>
                  </select>
                </label>
                <label className="field">
                  <span>Flags</span>
                  <${me}
                    options=${W}
                    value=${b.flags}
                    name="vip-add"
                    onChange=${l=>y({...b,flags:l})}
                  />
                </label>
              </div>
              <button className="btn" type="submit">Create VIP</button>
            </form>
          `}
          ${p&&o`<p className="error">${p}</p>`}
        </section>
        <section className="card">
          <div className="section-header">
            <h2>VIPs</h2>
            <button className="btn ghost" onClick=${_}>Refresh</button>
          </div>
          ${s.length===0?o`<p className="muted">No VIPs configured yet.</p>`:o`
                <div className="grid">
                  ${s.map(l=>o`
                      <div className="card">
                        <div className="row" style=${{fontWeight:600,gap:8}}>
                          <span className=${`dot ${l.healthy?"ok":"bad"}`}></span>
                          ${l.address}:${l.port} / ${l.proto}
                        </div>
                        <div className="muted" style=${{marginTop:6}}>
                          <div style=${{fontWeight:600,marginBottom:6}}>Flags</div>
                          <${pe}
                            mask=${l.flags}
                            options=${W}
                            emptyLabel=${l.flags===null?"Unknown":"No flags"}
                          />
                        </div>
                        <div className="row" style=${{marginTop:12}}>
                          <${M} className="btn" to=${`/vips/${Z(l)}`}>
                            Open
                          </${M}>
                          <${M}
                            className="btn secondary"
                            to=${`/vips/${Z(l)}/stats`}
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
    `}function De(){let{addToast:e}=J(),t=de(),a=Te(),s=E(()=>ae(t.vipId),[t.vipId]),[i,p]=v([]),[r,c]=v(""),[f,d]=v(""),[g,n]=v(!0),[m,b]=v({address:"",weight:100,flags:0}),[y,_]=v({}),[w,T]=v(null),[P,F]=v({flags:0,set:!0}),[l,$]=v({hash_function:0}),{groups:h,setGroups:C,refreshFromStorage:k,importFromRunningConfig:R}=ve(),[A,X]=v(""),[we,B]=v(""),[_e,Se]=v(!1),[Ie,Ce]=v(""),[re,ke]=v({add:0,update:0,remove:0}),H=async()=>{try{let u=await N.get("/vips/reals",s);p(u||[]);let S={};(u||[]).forEach(U=>{S[U.address]=U.weight}),_(S),c(""),n(!1)}catch(u){c(u.message||"request failed"),n(!1)}},ne=async()=>{try{let u=await N.get("/vips/flags",s);T(u?.flags??0),d("")}catch(u){d(u.message||"request failed")}};x(()=>{H(),ne()},[t.vipId]),x(()=>{if(!A){ke({add:0,update:0,remove:0});return}let u=h[A]||[],S=new Map(i.map(q=>[q.address,q])),U=new Map(u.map(q=>[q.address,q])),L=0,j=0,I=0;u.forEach(q=>{let oe=S.get(q.address);if(!oe){L+=1;return}(Number(oe.weight)!==Number(q.weight)||Number(oe.flags||0)!==Number(q.flags||0))&&(j+=1)}),i.forEach(q=>{U.has(q.address)||(I+=1)}),ke({add:L,update:j,remove:I})},[A,i,h]);let Ke=async u=>{try{let S=Number(y[u.address]);await N.post("/vips/reals",{vip:s,real:{address:u.address,weight:S,flags:u.flags||0}}),await H(),e("Real weight updated.","success")}catch(S){c(S.message||"request failed"),e(S.message||"Update failed.","error")}},Ye=async u=>{try{await N.del("/vips/reals",{vip:s,real:{address:u.address,weight:u.weight,flags:u.flags||0}}),await H(),e("Real removed.","success")}catch(S){c(S.message||"request failed"),e(S.message||"Remove failed.","error")}},Xe=async u=>{u.preventDefault();try{await N.post("/vips/reals",{vip:s,real:{address:m.address,weight:Number(m.weight),flags:Number(m.flags||0)}}),b({address:"",weight:100,flags:0}),await H(),e("Real added.","success")}catch(S){c(S.message||"request failed"),e(S.message||"Add failed.","error")}},Qe=async()=>{if(!A||!h[A]){B("Select a target group to apply.");return}Se(!0),B("");let u=h[A]||[],S=new Map(i.map(I=>[I.address,I])),U=new Map(u.map(I=>[I.address,I])),L=i.filter(I=>!U.has(I.address)),j=u.filter(I=>{let q=S.get(I.address);return q?Number(q.weight)!==Number(I.weight)||Number(q.flags||0)!==Number(I.flags||0):!0});try{L.length>0&&await N.put("/vips/reals/batch",{vip:s,action:1,reals:L.map(I=>({address:I.address,weight:Number(I.weight),flags:Number(I.flags||0)}))}),j.length>0&&await Promise.all(j.map(I=>N.post("/vips/reals",{vip:s,real:{address:I.address,weight:Number(I.weight),flags:Number(I.flags||0)}}))),await H(),e(`Applied target group "${A}".`,"success")}catch(I){B(I.message||"Failed to apply target group."),e(I.message||"Target group apply failed.","error")}finally{Se(!1)}},et=u=>{u.preventDefault();let S=Ie.trim();if(!S){B("Provide a name for the new target group.");return}if(h[S]){B("A target group with that name already exists.");return}let U={...h,[S]:i.map(L=>({address:L.address,weight:Number(L.weight),flags:Number(L.flags||0)}))};C(U),Ce(""),X(S),B(""),e(`Target group "${S}" saved.`,"success")},tt=async()=>{try{await N.del("/vips",s),e("VIP deleted.","success"),a("/")}catch(u){c(u.message||"request failed"),e(u.message||"Delete failed.","error")}},at=async u=>{u.preventDefault();try{await N.put("/vips/flags",{...s,flag:Number(P.flags||0),set:!!P.set}),await ne(),e("VIP flags updated.","success")}catch(S){d(S.message||"request failed"),e(S.message||"Flag update failed.","error")}},st=async u=>{u.preventDefault();try{await N.put("/vips/hash-function",{...s,hash_function:Number(l.hash_function)}),e("Hash function updated.","success")}catch(S){d(S.message||"request failed"),e(S.message||"Hash update failed.","error")}};return o`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>VIP Detail</h2>
              <p className="muted">${s.address}:${s.port} / ${s.proto}</p>
              ${w===null?o`<p className="muted">Flags: —</p>`:o`
                    <div style=${{marginTop:8}}>
                      <${pe}
                        mask=${w}
                        options=${W}
                        showStatus=${!0}
                        emptyLabel="No flags"
                      />
                    </div>
                  `}
            </div>
            <div className="row">
              <button className="btn ghost" onClick=${ne}>Refresh flags</button>
              <button className="btn danger" onClick=${tt}>Delete VIP</button>
            </div>
          </div>
          ${r&&o`<p className="error">${r}</p>`}
          ${f&&o`<p className="error">${f}</p>`}
          ${g?o`<p className="muted">Loading reals…</p>`:o`
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
                    ${i.map(u=>o`
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
                              value=${y[u.address]??u.weight}
                              onInput=${S=>_({...y,[u.address]:S.target.value})}
                            />
                          </td>
                          <td className="row">
                            <button className="btn" onClick=${()=>Ke(u)}>
                              Update
                            </button>
                            <button className="btn ghost" onClick=${()=>Ye(u)}>
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
            <form className="form" onSubmit=${at}>
              <div className="form-row">
                <label className="field">
                  <span>Flags</span>
                  <${me}
                    options=${W}
                    value=${P.flags}
                    name="vip-flag-change"
                    onChange=${u=>F({...P,flags:u})}
                  />
                </label>
                <label className="field">
                  <span>Set</span>
                  <select
                    value=${String(P.set)}
                    onChange=${u=>F({...P,set:u.target.value==="true"})}
                  >
                    <option value="true">Set</option>
                    <option value="false">Clear</option>
                  </select>
                </label>
              </div>
              <button className="btn" type="submit">Apply flag</button>
            </form>
            <form className="form" onSubmit=${st}>
              <div className="form-row">
                <label className="field">
                  <span>Hash function</span>
                  <input
                    type="number"
                    min="0"
                    value=${l.hash_function}
                    onInput=${u=>$({...l,hash_function:u.target.value})}
                  />
                </label>
              </div>
              <button className="btn" type="submit">Apply hash</button>
            </form>
          </div>
        </section>
        <section className="card">
          <h3>Add real</h3>
          <form className="form" onSubmit=${Xe}>
            <div className="form-row">
              <label className="field">
                <span>Real address</span>
                <input
                  value=${m.address}
                  onInput=${u=>b({...m,address:u.target.value})}
                  placeholder="10.0.0.1"
                  required
                />
              </label>
              <label className="field">
                <span>Weight</span>
                <input
                  type="number"
                  min="0"
                  value=${m.weight}
                  onInput=${u=>b({...m,weight:u.target.value})}
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
              <button className="btn ghost" type="button" onClick=${k}>
                Reload groups
              </button>
              <button
                className="btn ghost"
                type="button"
                onClick=${async()=>{try{await R(),e("Imported target groups from running config.","success")}catch(u){B(u.message||"Failed to import target groups."),e(u.message||"Import failed.","error")}}}
              >
                Import from running config
              </button>
            </div>
          </div>
          ${we&&o`<p className="error">${we}</p>`}
          <div className="form-row">
            <label className="field">
              <span>Target group</span>
              <select
                value=${A}
                onChange=${u=>X(u.target.value)}
                disabled=${Object.keys(h).length===0}
              >
                <option value="">Select group</option>
                ${Object.keys(h).map(u=>o`<option value=${u}>${u}</option>`)}
              </select>
            </label>
            <label className="field">
              <span>Preview</span>
              <input
                value=${`add ${re.add} \xB7 update ${re.update} \xB7 remove ${re.remove}`}
                readOnly
              />
            </label>
          </div>
          <div className="row">
            <button
              className="btn"
              type="button"
              onClick=${Qe}
              disabled=${_e||!A}
            >
              ${_e?"Applying...":"Apply target group"}
            </button>
          </div>
          <form className="form" onSubmit=${et}>
            <div className="form-row">
              <label className="field">
                <span>Save current reals as new group</span>
                <input
                  value=${Ie}
                  onInput=${u=>Ce(u.target.value)}
                  placeholder="edge-backends"
                />
              </label>
            </div>
            <button className="btn secondary" type="submit">Save target group</button>
          </form>
        </section>
      </main>
    `}function Be(){let e=de(),t=E(()=>ae(e.vipId),[e.vipId]),{points:a,error:s}=K({path:"/stats/vip",body:t}),i=E(()=>Y("/stats/vip"),[]),p=a[a.length-1]||{},r=a[a.length-2]||{},c=Number(p.v1??0),f=Number(p.v2??0),d=c-Number(r.v1??0),g=f-Number(r.v2??0),n=E(()=>[{label:i.v1,field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.15)"},{label:i.v2,field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[i]);return o`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>VIP Stats</h2>
              <p className="muted">${t.address}:${t.port} / ${t.proto}</p>
            </div>
          </div>
          ${s&&o`<p className="error">${s}</p>`}
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
                <td>${i.v1}</td>
                <td>${c}</td>
                <td>
                  <span className=${`delta ${d<0?"down":"up"}`}>
                    ${z(d)}
                  </span>
                </td>
              </tr>
              <tr>
                <td>${i.v2}</td>
                <td>${f}</td>
                <td>
                  <span className=${`delta ${g<0?"down":"up"}`}>
                    ${z(g)}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
          <${se} title="Traffic (delta/sec)" points=${a} keys=${n} diff=${!0} />
        </section>
      </main>
    `}let ze={"/stats/vip":{v1:"Packets",v2:"Bytes"},"/stats/real":{v1:"Packets",v2:"Bytes"},"/stats/lru":{v1:"Total packets",v2:"LRU hits"},"/stats/lru/miss":{v1:"TCP SYN misses",v2:"Non-SYN misses"},"/stats/lru/fallback":{v1:"Fallback LRU hits",v2:"Unused"},"/stats/lru/global":{v1:"Map lookup failures",v2:"Global LRU routed"},"/stats/xdp/total":{v1:"Packets",v2:"Bytes"},"/stats/xdp/pass":{v1:"Packets",v2:"Bytes"},"/stats/xdp/drop":{v1:"Packets",v2:"Bytes"},"/stats/xdp/tx":{v1:"Packets",v2:"Bytes"}},ye=[{title:"LRU",path:"/stats/lru"},{title:"LRU Miss",path:"/stats/lru/miss"},{title:"LRU Fallback",path:"/stats/lru/fallback"},{title:"LRU Global",path:"/stats/lru/global"},{title:"XDP Total",path:"/stats/xdp/total"},{title:"XDP Pass",path:"/stats/xdp/pass"},{title:"XDP Drop",path:"/stats/xdp/drop"},{title:"XDP Tx",path:"/stats/xdp/tx"}];function z(e){return Number.isFinite(e)?`${e>0?"+":""}${e}`:"0"}function Y(e){return ze[e]||{v1:"v1",v2:"v2"}}function Ve({title:e,path:t,diff:a=!1}){let{points:s,error:i}=K({path:t}),p=E(()=>Y(t),[t]),r=E(()=>[{label:p.v1,field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.2)"},{label:p.v2,field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[p]);return o`
      <div className="card">
        <h3>${e}</h3>
        ${i&&o`<p className="error">${i}</p>`}
        <${se} title=${e} points=${s} keys=${r} diff=${a} inlineTitle=${!1} />
      </div>
    `}function He({title:e,path:t}){let{points:a,error:s}=K({path:t}),i=E(()=>Y(t),[t]),p=a[a.length-1]||{},r=a[a.length-2]||{},c=Number(p.v1??0),f=Number(p.v2??0),d=c-Number(r.v1??0),g=f-Number(r.v2??0);return o`
      <div className="summary-card">
        <div className="summary-title">${e}</div>
        ${s?o`<p className="error">${s}</p>`:o`
              <div className="summary-row">
                <div className="stat">
                  <span className="muted">${i.v1}</span>
                  <strong>${c}</strong>
                </div>
                <div className="stat">
                  <span className="muted">${i.v1} delta/sec</span>
                  <strong className=${d<0?"delta down":"delta up"}>
                    ${z(d)}
                  </strong>
                </div>
              </div>
              <div className="summary-row">
                <div className="stat">
                  <span className="muted">${i.v2}</span>
                  <strong>${f}</strong>
                </div>
                <div className="stat">
                  <span className="muted">${i.v2} delta/sec</span>
                  <strong className=${g<0?"delta down":"delta up"}>
                    ${z(g)}
                  </strong>
                </div>
              </div>
            `}
      </div>
    `}function je(){let{data:e,error:t}=Ee(()=>N.get("/stats/userspace"),1e3,[]);return o`
      <main>
        <section className="card">
          <h2>Global Stats</h2>
          <p className="muted">Polling every second. Charts show per-second deltas.</p>
        </section>
        <section className="grid">
          ${ye.map(a=>o`<${Ve} title=${a.title} path=${a.path} diff=${!0} />`)}
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
        <section className="card">
          <div className="section-header">
            <div>
              <h3>Absolute & Rate of Change</h3>
              <p className="muted">Latest value and per-second delta.</p>
            </div>
          </div>
          <div className="summary-grid">
            ${ye.map(a=>o`<${He} title=${a.title} path=${a.path} />`)}
          </div>
        </section>
      </main>
    `}function Me(){let[e,t]=v([]),[a,s]=v(""),[i,p]=v([]),[r,c]=v(""),[f,d]=v(null),[g,n]=v("");x(()=>{let h=!0;return(async()=>{try{let k=await N.get("/vips");if(!h)return;t(k||[]),!a&&k&&k.length>0&&s(Z(k[0]))}catch(k){h&&n(k.message||"request failed")}})(),()=>{h=!1}},[]),x(()=>{if(!a)return;let h=ae(a),C=!0;return(async()=>{try{let R=await N.get("/vips/reals",h);if(!C)return;p(R||[]),R&&R.length>0?c(A=>A||R[0].address):c(""),n("")}catch(R){C&&n(R.message||"request failed")}})(),()=>{C=!1}},[a]),x(()=>{if(!r){d(null);return}let h=!0;return(async()=>{try{let k=await N.get("/reals/index",{address:r});if(!h)return;d(k?.index??null),n("")}catch(k){h&&n(k.message||"request failed")}})(),()=>{h=!1}},[r]);let{points:m,error:b}=K({path:"/stats/real",body:f!==null?{index:f}:null}),y=E(()=>Y("/stats/real"),[]),_=E(()=>[{label:y.v1,field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.2)"},{label:y.v2,field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[y]),w=m[m.length-1]||{},T=m[m.length-2]||{},P=Number(w.v1??0),F=Number(w.v2??0),l=P-Number(T.v1??0),$=F-Number(T.v2??0);return o`
      <main>
        <section className="card">
          <h2>Per-Real Stats</h2>
          <p className="muted">Select a VIP and real address to chart.</p>
          ${g&&o`<p className="error">${g}</p>`}
          <div className="form-row">
            <label className="field">
              <span>VIP</span>
              <select value=${a} onChange=${h=>s(h.target.value)}>
                ${e.map(h=>o`
                    <option value=${Z(h)}>
                      ${h.address}:${h.port} / ${h.proto}
                    </option>
                  `)}
              </select>
            </label>
            <label className="field">
              <span>Real</span>
              <select
                value=${r}
                onChange=${h=>c(h.target.value)}
                disabled=${i.length===0}
              >
                ${i.map(h=>o`
                    <option value=${h.address}>${h.address}</option>
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
          ${b&&o`<p className="error">${b}</p>`}
          ${f===null?o`<p className="muted">Select a real to start polling.</p>`:o`
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
                      <td>${y.v1}</td>
                      <td>${P}</td>
                      <td>
                        <span className=${`delta ${l<0?"down":"up"}`}>
                          ${z(l)}
                        </span>
                      </td>
                    </tr>
                    <tr>
                      <td>${y.v2}</td>
                      <td>${F}</td>
                      <td>
                        <span className=${`delta ${$<0?"down":"up"}`}>
                          ${z($)}
                        </span>
                      </td>
                    </tr>
                  </tbody>
                </table>
                <${se}
                  title="Traffic (delta/sec)"
                  points=${m}
                  keys=${_}
                  diff=${!0}
                />
              `}
        </section>
      </main>
    `}function We(){let{addToast:e}=J(),[t,a]=v(""),[s,i]=v(""),[p,r]=v(!0),[c,f]=v(""),d=V(!0),g=async()=>{if(d.current){r(!0),i("");try{let m=await fetch(`${N.base}/config/export`,{credentials:"same-origin",headers:{Accept:"application/x-yaml"}});if(!m.ok){let y=`HTTP ${m.status}`,_="";try{let w=await m.json();y=w?.error?.message||y,_=w?.error?.code||""}catch{}throw(m.status===401||_==="UNAUTHORIZED")&&N.notifyUnauthorized(y),new O(y,m.status,_)}let b=await m.text();if(!d.current)return;a(b||""),f(new Date().toLocaleString())}catch(m){d.current&&i(m.message||"request failed")}finally{d.current&&r(!1)}}},n=async()=>{if(t)try{await navigator.clipboard.writeText(t),e("Config copied to clipboard","info")}catch{e("Failed to copy config","error")}};return x(()=>(d.current=!0,g(),()=>{d.current=!1}),[]),o`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>Running config</h2>
              <p className="muted">Exported from /api/v1/config/export</p>
            </div>
            <div className="row">
              <button className="btn ghost" onClick=${n} disabled=${!t}>
                Copy YAML
              </button>
              <button className="btn" onClick=${g} disabled=${p}>
                Refresh
              </button>
            </div>
          </div>
          ${s&&o`<p className="error">${s}</p>`}
          ${p?o`<p className="muted">Loading config...</p>`:t?o`<pre className="yaml-view">${t}</pre>`:o`<p className="muted">No config data returned.</p>`}
          ${c&&o`<p className="muted">Last fetched ${c}</p>`}
        </section>
      </main>
    `}function Je(){let{addToast:e}=J(),{groups:t,setGroups:a,refreshFromStorage:s,importFromRunningConfig:i}=ve(),[p,r]=v(""),[c,f]=v(""),[d,g]=v({address:"",weight:100,flags:0}),[n,m]=v(""),[b,y]=v(!1);x(()=>{if(c){if(!t[c]){let $=Object.keys(t);f($[0]||"")}}else{let $=Object.keys(t);$.length>0&&f($[0])}},[t,c]);let _=$=>{$.preventDefault();let h=p.trim();if(!h){m("Provide a group name.");return}if(t[h]){m("That group already exists.");return}a({...t,[h]:[]}),r(""),f(h),m(""),e(`Target group "${h}" created.`,"success")},w=$=>{let h={...t};delete h[$],a(h),e(`Target group "${$}" removed.`,"success")},T=$=>{if($.preventDefault(),!c){m("Select a group to add a real.");return}let h=fe(d);if(!h){m("Provide a valid real address.");return}let C=t[c]||[],k=C.some(R=>R.address===h.address)?C.map(R=>R.address===h.address?h:R):C.concat(h);a({...t,[c]:k}),g({address:"",weight:100,flags:0}),m(""),e("Real saved to target group.","success")},P=$=>{if(!c)return;let C=(t[c]||[]).filter(k=>k.address!==$);a({...t,[c]:C})},F=($,h)=>{if(!c)return;let k=(t[c]||[]).map(R=>R.address===$?{...R,...h}:R);a({...t,[c]:k})};return o`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>Target groups</h2>
              <p className="muted">Define reusable sets of reals (address + weight).</p>
            </div>
            <div className="row">
              <button className="btn ghost" type="button" onClick=${s}>
                Reload groups
              </button>
              <button className="btn ghost" type="button" onClick=${async()=>{y(!0);try{await i(),e("Imported target groups from running config.","success"),m("")}catch($){m($.message||"Failed to import target groups."),e($.message||"Import failed.","error")}finally{y(!1)}}} disabled=${b}>
                ${b?"Importing...":"Import from running config"}
              </button>
            </div>
          </div>
          ${n&&o`<p className="error">${n}</p>`}
          <form className="form" onSubmit=${_}>
            <div className="form-row">
              <label className="field">
                <span>New group name</span>
                <input
                  value=${p}
                  onInput=${$=>r($.target.value)}
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
                  value=${c}
                  onChange=${$=>f($.target.value)}
                  disabled=${Object.keys(t).length===0}
                >
                  ${Object.keys(t).map($=>o`<option value=${$}>${$}</option>`)}
                </select>
              </label>
              ${c&&o`<button className="btn danger" type="button" onClick=${()=>w(c)}>
                Delete group
              </button>`}
            </div>
          </div>
          ${c?o`
                <table className="table">
                  <thead>
                    <tr>
                      <th>Address</th>
                      <th>Weight</th>
                      <th></th>
                    </tr>
                  </thead>
                  <tbody>
                    ${(t[c]||[]).map($=>o`
                        <tr>
                          <td>${$.address}</td>
                          <td>
                            <input
                              className="inline-input"
                              type="number"
                              min="0"
                              value=${$.weight}
                              onInput=${h=>F($.address,{weight:Number(h.target.value)})}
                            />
                          </td>
                          <td>
                            <button
                              className="btn ghost"
                              type="button"
                              onClick=${()=>P($.address)}
                            >
                              Remove
                            </button>
                          </td>
                        </tr>
                      `)}
                  </tbody>
                </table>
                <form className="form" onSubmit=${T}>
                  <div className="form-row">
                    <label className="field">
                      <span>Real address</span>
                      <input
                        value=${d.address}
                        onInput=${$=>g({...d,address:$.target.value})}
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
                        onInput=${$=>g({...d,weight:$.target.value})}
                        required
                      />
                    </label>
                  </div>
                  <button className="btn secondary" type="submit">Save real</button>
                </form>
              `:o`<p className="muted">No groups yet. Create one to add reals.</p>`}
        </section>
      </main>
    `}function Ze(){let[e,t]=v({initialized:!1,ready:!1}),[a,s]=v([]),[i,p]=v({checking:!0,required:!1,username:""}),r=V({}),c=Re(async()=>{try{let b=await N.get("/lb/status");t(b||{initialized:!1,ready:!1}),p(y=>({...y,checking:!1,required:!1}))}catch(b){if(b instanceof O&&b.unauthorized){t({initialized:!1,ready:!1}),p(y=>({...y,checking:!1,required:!0}));return}t({initialized:!1,ready:!1}),p(y=>({...y,checking:!1}))}},[]),f=(b,y="info")=>{let _=`${Date.now()}-${Math.random().toString(16).slice(2)}`;s(w=>w.concat({id:_,message:b,kind:y})),r.current[_]=setTimeout(()=>{s(w=>w.filter(T=>T.id!==_)),delete r.current[_]},4e3)},d=b=>{r.current[b]&&(clearTimeout(r.current[b]),delete r.current[b]),s(y=>y.filter(_=>_.id!==b))},g=async(b,y)=>{let _=await N.login(b,y);p({checking:!1,required:!1,username:_?.username||b}),await c(),f("Signed in.","success")},n=async()=>{try{await N.logout()}catch(b){f(b.message||"Sign out failed.","error")}t({initialized:!1,ready:!1}),p({checking:!1,required:!0,username:""})};x(()=>(N.setAuthHandlers({onUnauthorized:()=>{p(b=>({...b,checking:!1,required:!0,username:""}))}}),()=>{N.setAuthHandlers({})}),[]),x(()=>{let b=!0;(async()=>{b&&await c()})();let _=setInterval(()=>{!b||i.required||c()},5e3);return()=>{b=!1,clearInterval(_)}},[c,i.required]);let m=E(()=>({required:i.required,username:i.username,login:g,logout:n}),[i.required,i.username]);return x(()=>()=>{Object.keys(r.current).forEach(b=>clearTimeout(r.current[b]))},[]),o`
      <${ie}>
        <${ee.Provider} value=${m}>
          <${Ne} checking=${i.checking} required=${i.required}>
            <${Q.Provider} value=${{addToast:f}}>
              <${Oe} status=${e} />
              <${ce}>
                <${D} path="/" element=${o`<${Ue} />`} />
                <${D} path="/vips/:vipId" element=${o`<${De} />`} />
                <${D} path="/vips/:vipId/stats" element=${o`<${Be} />`} />
                <${D} path="/target-groups" element=${o`<${Je} />`} />
                <${D} path="/stats/global" element=${o`<${je} />`} />
                <${D} path="/stats/real" element=${o`<${Me} />`} />
                <${D} path="/config" element=${o`<${We} />`} />
              </${ce}>
              <${Le} toasts=${a} onDismiss=${d} />
            </${Q.Provider}>
          </${Ne}>
        </${ee.Provider}>
      </${ie}>
    `}ReactDOM.createRoot(document.getElementById("root")).render(o`<${Ze} />`)})();})();
