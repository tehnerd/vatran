(()=>{(()=>{let{useEffect:E,useMemo:H,useRef:W,useState:h,useContext:he,useCallback:Oe}=React,{BrowserRouter:fe,Routes:ge,Route:B,NavLink:O,Link:K,useParams:be,useNavigate:ze}=ReactRouterDOM,n=htm.bind(React.createElement),ne=React.createContext({addToast:()=>{}}),oe=React.createContext({required:!1,username:"",login:async()=>{},logout:async()=>{}}),Y=[{label:"NO_SPORT",value:1},{label:"NO_LRU",value:2},{label:"QUIC_VIP",value:4},{label:"DPORT_HASH",value:8},{label:"SRC_ROUTING",value:16},{label:"LOCAL_VIP",value:32},{label:"GLOBAL_LRU",value:64},{label:"HASH_SRC_DST_PORT",value:128},{label:"UDP_STABLE_ROUTING_VIP",value:256},{label:"UDP_FLOW_MIGRATION",value:512}];function X(){return he(ne)}function Be(){return he(oe)}class z extends Error{constructor(t,a=0,s=""){super(t||"request failed"),this.name="ApiError",this.status=a,this.code=s,this.unauthorized=a===401||s==="UNAUTHORIZED"}}function Ve(e){return String(e).replace(/[^a-z0-9_-]/gi,"_")}function je(e,t,a){let s=Number(e)||0,c=Number(t)||0;return a?s|c:s&~c}function Me(e,t){let a=Number(e)||0;return t.filter(s=>(a&s.value)!==0)}function ve(e,t){let a=String(e??"").trim();if(!a)return;let c=a.split(/[\s,]+/).filter(Boolean).map(o=>Number(o));if(c.findIndex(o=>!Number.isFinite(o)||!Number.isInteger(o))!==-1)throw new Error(`${t} must be a comma- or space-separated list of integers.`);return c}function $e({mask:e,options:t,showStatus:a=!1,emptyLabel:s="None"}){let c=Number(e)||0,p=a?t:Me(c,t),o=a?2:1;return n`
      <table className="table flag-table">
        <thead>
          <tr>
            <th>Flag</th>
            ${a?n`<th>Enabled</th>`:null}
          </tr>
        </thead>
        <tbody>
          ${p.length===0?n`<tr><td colspan=${o} className="muted">${s}</td></tr>`:p.map(d=>{let g=(c&d.value)!==0;return n`
                  <tr>
                    <td>${d.label}</td>
                    ${a?n`<td>${g?"Yes":"No"}</td>`:null}
                  </tr>
                `})}
        </tbody>
      </table>
    `}function Ne({options:e,value:t,onChange:a,name:s}){let c=Number(t)||0,p=Ve(s||"flags");return n`
      <div className="flag-selector">
        <table className="table flag-table">
          <thead>
            <tr>
              <th>Enabled</th>
              <th>Flag</th>
            </tr>
          </thead>
          <tbody>
            ${e.map(o=>{let d=`${p}-${o.value}`,g=(c&o.value)===o.value;return n`
                <tr>
                  <td>
                    <input
                      id=${d}
                      type="checkbox"
                      checked=${g}
                      onChange=${u=>a(je(c,o.value,u.target.checked))}
                    />
                  </td>
                  <td>
                    <label className="flag-option" htmlFor=${d}>${o.label}</label>
                  </td>
                </tr>
              `})}
          </tbody>
        </table>
      </div>
    `}let N={base:"/api/v1",authHandlers:{onUnauthorized:null},setAuthHandlers(e={}){N.authHandlers={onUnauthorized:e.onUnauthorized||null}},notifyUnauthorized(e){typeof N.authHandlers.onUnauthorized=="function"&&N.authHandlers.onUnauthorized(e||"Authentication required")},async parseJSON(e){try{return await e.json()}catch{return null}},async buildError(e,t){let a=t?.error?.code||"",s=t?.error?.message||`HTTP ${e.status}`;return(e.status===401||a==="UNAUTHORIZED")&&N.notifyUnauthorized(s),new z(s,e.status,a)},async request(e,t={}){let a={method:t.method||"GET",credentials:"same-origin",headers:{Accept:"application/json",...t.headers||{}}},s=`${N.base}${e}`;if(t.body!==void 0&&t.body!==null)if(a.method==="GET"){let o=new URLSearchParams;Object.entries(t.body).forEach(([g,u])=>{if(u!=null){if(Array.isArray(u)){u.forEach(f=>o.append(g,String(f)));return}if(typeof u=="object"){o.set(g,JSON.stringify(u));return}o.set(g,String(u))}});let d=o.toString();d&&(s+=`${s.includes("?")?"&":"?"}${d}`)}else a.headers["Content-Type"]="application/json",a.body=JSON.stringify(t.body);let c=await fetch(s,a),p=await N.parseJSON(c);if(!p)throw new z("invalid JSON response",c.status);if(!c.ok)throw await N.buildError(c,p);if(!p.success)throw new z(p.error?.message||"request failed",c.status,p.error?.code);return p.data},get(e,t){return N.request(e,{method:"GET",body:t})},post(e,t){return N.request(e,{method:"POST",body:t})},put(e,t){return N.request(e,{method:"PUT",body:t})},del(e,t){return N.request(e,{method:"DELETE",body:t})},async login(e,t){let a=await fetch("/login",{method:"POST",credentials:"same-origin",headers:{Accept:"application/json","Content-Type":"application/json"},body:JSON.stringify({username:e,password:t})}),s=await N.parseJSON(a);if(!s)throw new z("invalid JSON response",a.status);if(!a.ok)throw await N.buildError(a,s);if(!s.success)throw new z(s.error?.message||"login failed",a.status,s.error?.code);return s.data||{}},async logout(){let e=await fetch("/logout",{method:"POST",credentials:"same-origin",headers:{Accept:"application/json"}}),t=await N.parseJSON(e);if(!e.ok)throw await N.buildError(e,t||{});if(t&&t.success===!1)throw new z(t.error?.message||"logout failed",e.status,t.error?.code);return t?.data||{}}},ye="vatran_target_groups";function _e(e){if(!e||!e.address)return null;let t=String(e.address).trim();if(!t)return null;let a=Number(e.weight),s=Number(e.flags??0);return{address:t,weight:Number.isFinite(a)?a:0,flags:Number.isFinite(s)?s:0}}function we(e){if(!e||typeof e!="object")return{};let t={};return Object.entries(e).forEach(([a,s])=>{let c=String(a).trim();if(!c)return;let p=Array.isArray(s)?s.map(_e).filter(Boolean):[],o=[],d=new Set;p.forEach(g=>{d.has(g.address)||(d.add(g.address),o.push(g))}),t[c]=o}),t}function le(){if(typeof localStorage>"u")return{};try{let e=localStorage.getItem(ye);return e?we(JSON.parse(e)):{}}catch{return{}}}function ke(e){if(!(typeof localStorage>"u"))try{localStorage.setItem(ye,JSON.stringify(e))}catch{}}function We(e,t){let a={...e};return Object.entries(t||{}).forEach(([s,c])=>{a[s]||(a[s]=c)}),a}function Se(){let[e,t]=h(()=>le());return E(()=>{ke(e)},[e]),{groups:e,setGroups:t,refreshFromStorage:()=>{t(le())},importFromRunningConfig:async()=>{let c=await N.get("/config/export/json"),p=we(c?.target_groups||{}),o=We(le(),p);return t(o),ke(o),o}}}function Q(e){return`${encodeURIComponent(e.address)}:${e.port}:${e.proto}`}function ie(e){let t=e.split(":"),a=Number(t.pop()||0),s=Number(t.pop()||0);return{address:decodeURIComponent(t.join(":")),port:s,proto:a}}function J(){return{type:"dummy",port:"",interval_ms:5e3,timeout_ms:2e3,healthy_threshold:3,unhealthy_threshold:3,http_path:"/healthz",http_expected_status:200,http_host:"",https_path:"/healthz",https_expected_status:200,https_host:"",https_skip_tls_verify:!1}}function Je(e){if(!e||typeof e!="object")return J();let t=J();return{...t,type:e.type||t.type,port:Number.isFinite(Number(e.port))&&Number(e.port)>0?String(e.port):"",interval_ms:Number.isFinite(Number(e.interval_ms))&&Number(e.interval_ms)>0?Number(e.interval_ms):t.interval_ms,timeout_ms:Number.isFinite(Number(e.timeout_ms))&&Number(e.timeout_ms)>0?Number(e.timeout_ms):t.timeout_ms,healthy_threshold:Number.isFinite(Number(e.healthy_threshold))&&Number(e.healthy_threshold)>0?Number(e.healthy_threshold):t.healthy_threshold,unhealthy_threshold:Number.isFinite(Number(e.unhealthy_threshold))&&Number(e.unhealthy_threshold)>0?Number(e.unhealthy_threshold):t.unhealthy_threshold,http_path:e.http?.path||t.http_path,http_expected_status:Number.isFinite(Number(e.http?.expected_status))&&Number(e.http?.expected_status)>0?Number(e.http.expected_status):t.http_expected_status,http_host:e.http?.host||"",https_path:e.https?.path||t.https_path,https_expected_status:Number.isFinite(Number(e.https?.expected_status))&&Number(e.https?.expected_status)>0?Number(e.https.expected_status):t.https_expected_status,https_host:e.https?.host||"",https_skip_tls_verify:!!e.https?.skip_tls_verify}}function Ze(e,t,a=[]){let[s,c]=h(null),[p,o]=h(""),[d,g]=h(!0);return E(()=>{let u=!0,f=async()=>{try{let m=await e();u&&(c(m),o(""),g(!1))}catch(m){u&&(o(m.message||"request failed"),g(!1))}};f();let l=setInterval(f,t);return()=>{u=!1,clearInterval(l)}},a),{data:s,error:p,loading:d}}function ee({path:e,body:t,intervalMs:a=1e3,limit:s=60}){let[c,p]=h([]),[o,d]=h(""),g=H(()=>JSON.stringify(t||{}),[t]);return E(()=>{if(t===null)return p([]),d(""),()=>{};let u=!0,f=async()=>{try{let m=await N.get(e,t);if(!u)return;let v=new Date().toLocaleTimeString();p(_=>_.concat({label:v,...m}).slice(-s)),d("")}catch(m){u&&d(m.message||"request failed")}};f();let l=setInterval(f,a);return()=>{u=!1,clearInterval(l)}},[e,g,a,s]),{points:c,error:o}}function Ie({title:e,points:t,keys:a,diff:s=!1,height:c=120,showTitle:p=!1,selectedLabel:o=null,onPointSelect:d=null,onLegendSelect:g=null}){let u=W(null),f=W(null);return E(()=>{if(!u.current)return;f.current||(f.current=new Chart(u.current,{type:"line",data:{labels:[],datasets:[]},options:{responsive:!0,maintainAspectRatio:!1,animation:!1,scales:{x:{grid:{display:!1}},y:{beginAtZero:!s}},plugins:{legend:{display:!0,position:"bottom"},title:{display:p&&!!e,text:e}}}}));let l=f.current,m=new Map((l.data.datasets||[]).filter(S=>typeof S.hidden<"u").map(S=>[S.label,S.hidden])),v=t.map(S=>S.label),_=o?a.filter(S=>S.label===o):a,I=o&&_.length===0?a:_;return l.data.labels=v,l.data.datasets=I.map(S=>{let F=t.map(L=>L[S.field]||0),A=s?F.map((L,i)=>i===0?0:L-F[i-1]):F;return{label:S.label,data:A,borderColor:S.color,backgroundColor:S.fill,borderWidth:2,tension:.3,hidden:m.get(S.label)}}),l.options.onClick=(S,F)=>{if(!d||!F||F.length===0)return;let A=F[0].datasetIndex,L=l.data.datasets?.[A]?.label;L&&d(L)},l.options.plugins&&l.options.plugins.legend&&(l.options.plugins.legend.onClick=(S,F)=>{if(!g)return;let A=F?.text;A&&g(A)}),l.options.scales.y.beginAtZero=!s,l.options.plugins.title.display=p&&!!e,l.options.plugins.title.text=e||"",l.update(),()=>{}},[t,a,e,s,p,o,d,g]),E(()=>()=>{f.current&&(f.current.destroy(),f.current=null)},[]),n`<canvas ref=${u} height=${c}></canvas>`}function ce({title:e,points:t,keys:a,diff:s=!1,inlineTitle:c=!0}){let[p,o]=h(!1),[d,g]=h(null),u=W(!1);return n`
      <div className="chart-wrap">
        <button className="chart-zoom-button" type="button" onClick=${()=>o(!0)}>
          <span className="zoom-icon" aria-hidden="true">+</span>
          Zoom
        </button>
        <div
          className="chart-click"
          onClick=${()=>{if(u.current){u.current=!1;return}o(!0)}}
        >
          <${Ie}
            title=${e}
            points=${t}
            keys=${a}
            diff=${s}
            height=${120}
            showTitle=${c&&!!e}
            selectedLabel=${d}
            onPointSelect=${f=>{g(l=>l===f?null:f),u.current=!0,setTimeout(()=>{u.current=!1},0)}}
            onLegendSelect=${f=>{g(l=>l===f?null:f),u.current=!0,setTimeout(()=>{u.current=!1},0)}}
          />
        </div>
        ${p&&n`
          <div className="chart-overlay" onClick=${()=>o(!1)}>
            <div className="chart-modal" onClick=${f=>f.stopPropagation()}>
              <div className="row chart-modal-header">
                <div>
                  <h3>${e||"Chart"}</h3>
                  ${s?n`<p className="muted">Per-second delta.</p>`:""}
                </div>
                <button className="btn ghost" onClick=${()=>o(!1)}>
                  Close
                </button>
              </div>
              <div className="chart-zoom">
                <${Ie}
                  title=${e}
                  points=${t}
                  keys=${a}
                  diff=${s}
                  height=${360}
                  showTitle=${!1}
                  selectedLabel=${d}
                  onPointSelect=${f=>g(l=>l===f?null:f)}
                  onLegendSelect=${f=>g(l=>l===f?null:f)}
                />
              </div>
            </div>
          </div>
        `}
      </div>
    `}function Ke(){let{login:e}=Be(),[t,a]=h(""),[s,c]=h(""),[p,o]=h(""),[d,g]=h(!1),u=async f=>{if(f.preventDefault(),!t.trim()||!s){o("Username and password are required.");return}g(!0),o("");try{await e(t.trim(),s)}catch(l){o(l.message||"Login failed.")}finally{g(!1)}};return n`
      <main className="auth-main">
        <section className="auth-card">
          <h1>Sign in</h1>
          <p className="muted">Authentication is required to access Vatran.</p>
          ${p?n`<p className="error">${p}</p>`:null}
          <form className="form" onSubmit=${u}>
            <label className="field">
              <span>Username</span>
              <input
                value=${t}
                onInput=${f=>a(f.target.value)}
                autoComplete="username"
                required
              />
            </label>
            <label className="field">
              <span>Password</span>
              <input
                type="password"
                value=${s}
                onInput=${f=>c(f.target.value)}
                autoComplete="current-password"
                required
              />
            </label>
            <button className="btn" type="submit" disabled=${d}>
              ${d?"Signing in...":"Sign in"}
            </button>
          </form>
        </section>
      </main>
    `}function Te({checking:e,required:t,children:a}){return e?n`
        <main className="auth-main">
          <section className="auth-card">
            <p className="muted">Checking authentication...</p>
          </section>
        </main>
      `:t?n`<${Ke} />`:a}function Ye({toasts:e,onDismiss:t}){return n`
      <div className="toast-stack">
        ${e.map(a=>n`
            <div className=${`toast ${a.kind}`}>
              <span>${a.message}</span>
              <button className="toast-close" onClick=${()=>t(a.id)}>
                ×
              </button>
            </div>
          `)}
      </div>
    `}function Xe({status:e}){return n`
      <header>
        <div>
          <div style=${{fontSize:20,fontWeight:700}}>Vatran</div>
          <div className="status-pill">
            <span className=${`dot ${e.ready?"ok":""}`}></span>
            ${e.ready?"Ready":"Not ready"}
          </div>
        </div>
        <nav>
          <${O} to="/" end className=${({isActive:t})=>t?"active":""}>
            Dashboard
          </${O}>
          <${O} to="/stats/global" className=${({isActive:t})=>t?"active":""}>
            Global stats
          </${O}>
          <${O} to="/stats/real" className=${({isActive:t})=>t?"active":""}>
            Per-real stats
          </${O}>
          <${O}
            to="/target-groups"
            className=${({isActive:t})=>t?"active":""}
          >
            Target groups
          </${O}>
          <${O} to="/config" className=${({isActive:t})=>t?"active":""}>
            Config export
          </${O}>
        </nav>
      </header>
    `}function Qe(){let{addToast:e}=X(),[t,a]=h({initialized:!1,ready:!1}),[s,c]=h([]),[p,o]=h(""),[d,g]=h(!1),[u,f]=h(!1),[l,m]=h({main_interface:"",balancer_prog_path:"",healthchecking_prog_path:"",default_mac:"",local_mac:"",root_map_path:"",root_map_pos:2,katran_src_v4:"",katran_src_v6:"",use_root_map:!1,max_vips:512,max_reals:4096,hash_function:"maglev_v2",forwarding_cores:"",numa_nodes:""}),[v,_]=h({address:"",port:80,proto:6,flags:0}),I=async()=>{try{let i=await N.get("/lb/status"),$=await N.get("/vips"),b=await Promise.all(($||[]).map(async C=>{let P=null,R=!1;try{P=(await N.get("/vips/flags",{address:C.address,port:C.port,proto:C.proto}))?.flags??0}catch{P=null}try{let q=await N.get("/vips/reals",{address:C.address,port:C.port,proto:C.proto});R=Array.isArray(q)&&q.some(ae=>!!ae?.healthy)}catch{R=!1}return{...C,flags:P,healthy:R}}));a(i||{initialized:!1,ready:!1}),c(b),o("")}catch(i){o(i.message||"request failed")}};E(()=>{let i=!0;return(async()=>{i&&await I()})(),()=>{i=!1}},[]);let S=async i=>{i.preventDefault();try{let $=ve(l.forwarding_cores,"Forwarding cores"),b=ve(l.numa_nodes,"NUMA nodes"),C={...l,forwarding_cores:$,numa_nodes:b,root_map_pos:l.root_map_pos===""?void 0:Number(l.root_map_pos),max_vips:Number(l.max_vips),max_reals:Number(l.max_reals),hash_function:l.hash_function};await N.post("/lb/create",C),o(""),g(!1),e("Load balancer initialized.","success"),await I()}catch($){o($.message||"request failed"),e($.message||"Initialize failed.","error")}},F=async i=>{i.preventDefault();try{await N.post("/vips",{...v,port:Number(v.port),proto:Number(v.proto),flags:Number(v.flags||0)}),_({address:"",port:80,proto:6,flags:0}),o(""),f(!1),e("VIP created.","success"),await I()}catch($){o($.message||"request failed"),e($.message||"VIP create failed.","error")}},A=async()=>{try{await N.post("/lb/load-bpf-progs"),o(""),e("BPF programs loaded.","success"),await I()}catch(i){o(i.message||"request failed"),e(i.message||"Load BPF programs failed.","error")}},L=async()=>{try{await N.post("/lb/attach-bpf-progs"),o(""),e("BPF programs attached.","success"),await I()}catch(i){o(i.message||"request failed"),e(i.message||"Attach BPF programs failed.","error")}};return n`
      <main>
        <section className="card">
          <h2>Load balancer</h2>
          <p>Initialized: ${t.initialized?"yes":"no"}</p>
          <p>Ready: ${t.ready?"yes":"no"}</p>
          <div className="row">
            ${!t.initialized&&n`
              <button className="btn" onClick=${()=>g(i=>!i)}>
                ${d?"Close":"Initialize"}
              </button>
            `}
            <button className="btn secondary" onClick=${()=>f(i=>!i)}>
              ${u?"Close":"Create VIP"}
            </button>
          </div>
          ${!t.ready&&n`
            <div className="row" style=${{marginTop:12}}>
              <button
                className="btn ghost"
                disabled=${!t.initialized}
                onClick=${A}
              >
                Load BPF Programs
              </button>
              <button
                className="btn ghost"
                disabled=${!t.initialized}
                onClick=${L}
              >
                Attach BPF Programs
              </button>
            </div>
          `}
          ${d&&n`
            <form className="form" onSubmit=${S}>
              <div className="form-row">
                <label className="field">
                  <span>Main interface</span>
                  <input
                    value=${l.main_interface}
                    onInput=${i=>m({...l,main_interface:i.target.value})}
                    placeholder="eth0"
                    required
                  />
                </label>
                <label className="field">
                  <span>Balancer prog path</span>
                  <input
                    value=${l.balancer_prog_path}
                    onInput=${i=>m({...l,balancer_prog_path:i.target.value})}
                    placeholder="/path/to/bpf.o"
                    required
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Healthchecking path</span>
                  <input
                    value=${l.healthchecking_prog_path}
                    onInput=${i=>m({...l,healthchecking_prog_path:i.target.value})}
                    placeholder="/path/to/hc.o"
                  />
                </label>
                <label className="field">
                  <span>Default MAC</span>
                  <input
                    value=${l.default_mac}
                    onInput=${i=>m({...l,default_mac:i.target.value})}
                    placeholder="aa:bb:cc:dd:ee:ff"
                    required
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Local MAC</span>
                  <input
                    value=${l.local_mac}
                    onInput=${i=>m({...l,local_mac:i.target.value})}
                    placeholder="11:22:33:44:55:66"
                    required=${l.healthchecking_prog_path?.trim()!==""}
                  />
                </label>
                <label className="field">
                  <span>Hash function</span>
                  <select
                    value=${l.hash_function}
                    onInput=${i=>m({...l,hash_function:i.target.value})}
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
                    value=${l.root_map_path}
                    onInput=${i=>m({...l,root_map_path:i.target.value})}
                    placeholder="/sys/fs/bpf/root_map"
                  />
                </label>
                <label className="field">
                  <span>Root map position</span>
                  <input
                    type="number"
                    min="0"
                    value=${l.root_map_pos}
                    onInput=${i=>m({...l,root_map_pos:i.target.value})}
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Katran src v4</span>
                  <input
                    value=${l.katran_src_v4}
                    onInput=${i=>m({...l,katran_src_v4:i.target.value})}
                    placeholder="10.0.0.1"
                  />
                </label>
                <label className="field">
                  <span>Katran src v6</span>
                  <input
                    value=${l.katran_src_v6}
                    onInput=${i=>m({...l,katran_src_v6:i.target.value})}
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
                    value=${l.max_vips}
                    onInput=${i=>m({...l,max_vips:i.target.value})}
                  />
                </label>
                <label className="field">
                  <span>Max reals</span>
                  <input
                    type="number"
                    min="1"
                    value=${l.max_reals}
                    onInput=${i=>m({...l,max_reals:i.target.value})}
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Forwarding cores (optional)</span>
                  <input
                    value=${l.forwarding_cores}
                    onInput=${i=>m({...l,forwarding_cores:i.target.value})}
                    placeholder="0,1,2,3"
                  />
                  <span className="muted">Comma or space separated CPU core IDs.</span>
                </label>
                <label className="field">
                  <span>NUMA nodes (optional)</span>
                  <input
                    value=${l.numa_nodes}
                    onInput=${i=>m({...l,numa_nodes:i.target.value})}
                    placeholder="0,0,1,1"
                  />
                  <span className="muted">Match the forwarding cores length.</span>
                </label>
              </div>
              <label className="field checkbox">
                <input
                  type="checkbox"
                  checked=${l.use_root_map}
                  onChange=${i=>m({...l,use_root_map:i.target.checked})}
                />
                <span>Use root map</span>
              </label>
              <button className="btn" type="submit">Initialize LB</button>
            </form>
          `}
          ${u&&n`
            <form className="form" onSubmit=${F}>
              <div className="form-row">
                <label className="field">
                  <span>VIP address</span>
                  <input
                    value=${v.address}
                    onInput=${i=>_({...v,address:i.target.value})}
                    placeholder="1.2.3.4"
                    required
                  />
                </label>
                <label className="field">
                  <span>Port</span>
                  <input
                    type="number"
                    value=${v.port}
                    onInput=${i=>_({...v,port:i.target.value})}
                    required
                  />
                </label>
                <label className="field">
                  <span>Protocol</span>
                  <select
                    value=${v.proto}
                    onChange=${i=>_({...v,proto:i.target.value})}
                  >
                    <option value="6">TCP (6)</option>
                    <option value="17">UDP (17)</option>
                  </select>
                </label>
                <label className="field">
                  <span>Flags</span>
                  <${Ne}
                    options=${Y}
                    value=${v.flags}
                    name="vip-add"
                    onChange=${i=>_({...v,flags:i})}
                  />
                </label>
              </div>
              <button className="btn" type="submit">Create VIP</button>
            </form>
          `}
          ${p&&n`<p className="error">${p}</p>`}
        </section>
        <section className="card">
          <div className="section-header">
            <h2>VIPs</h2>
            <button className="btn ghost" onClick=${I}>Refresh</button>
          </div>
          ${s.length===0?n`<p className="muted">No VIPs configured yet.</p>`:n`
                <div className="grid">
                  ${s.map(i=>n`
                      <div className="card">
                        <div className="row" style=${{fontWeight:600,gap:8}}>
                          <span className=${`dot ${i.healthy?"ok":"bad"}`}></span>
                          ${i.address}:${i.port} / ${i.proto}
                        </div>
                        <div className="muted" style=${{marginTop:6}}>
                          <div style=${{fontWeight:600,marginBottom:6}}>Flags</div>
                          <${$e}
                            mask=${i.flags}
                            options=${Y}
                            emptyLabel=${i.flags===null?"Unknown":"No flags"}
                          />
                        </div>
                        <div className="row" style=${{marginTop:12}}>
                          <${K} className="btn" to=${`/vips/${Q(i)}`}>
                            Open
                          </${K}>
                          <${K}
                            className="btn secondary"
                            to=${`/vips/${Q(i)}/stats`}
                          >
                            Stats
                          </${K}>
                        </div>
                      </div>
                    `)}
                </div>
              `}
        </section>
      </main>
    `}function et(){let{addToast:e}=X(),t=be(),a=ze(),s=H(()=>ie(t.vipId),[t.vipId]),[c,p]=h([]),[o,d]=h(""),[g,u]=h(""),[f,l]=h(!0),[m,v]=h({address:"",weight:100,flags:0}),[_,I]=h({}),[S,F]=h(null),[A,L]=h({flags:0,set:!0}),[i,$]=h({hash_function:0}),{groups:b,setGroups:C,refreshFromStorage:P,importFromRunningConfig:R}=Se(),[q,ae]=h(""),[Pe,V]=h(""),[xe,Re]=h(!1),[Fe,Ae]=h(""),[de,qe]=h({add:0,update:0,remove:0}),[y,G]=h(()=>J()),[se,ue]=h(!1),[dt,Ee]=h(!0),[Le,j]=h(""),[Ge,Ue]=h(!1),[De,He]=h(!1),Z=async()=>{try{let r=await N.get("/vips/reals",s);p(r||[]);let k={};(r||[]).forEach(U=>{k[U.address]=U.weight}),I(k),d(""),l(!1)}catch(r){d(r.message||"request failed"),l(!1)}},pe=async()=>{try{let r=await N.get("/vips/flags",s);F(r?.flags??0),u("")}catch(r){u(r.message||"request failed")}},re=async()=>{Ee(!0);try{let r=await N.get("/vips/healthcheck",s);r?(ue(!0),G(Je(r))):(ue(!1),G(J())),j("")}catch(r){j(r.message||"request failed")}finally{Ee(!1)}};E(()=>{ue(!1),G(J()),j(""),Z(),pe(),re()},[t.vipId]),E(()=>{if(!q){qe({add:0,update:0,remove:0});return}let r=b[q]||[],k=new Map(c.map(T=>[T.address,T])),U=new Map(r.map(T=>[T.address,T])),x=0,D=0,w=0;r.forEach(T=>{let me=k.get(T.address);if(!me){x+=1;return}(Number(me.weight)!==Number(T.weight)||Number(me.flags||0)!==Number(T.flags||0))&&(D+=1)}),c.forEach(T=>{U.has(T.address)||(w+=1)}),qe({add:x,update:D,remove:w})},[q,c,b]);let ut=async r=>{try{let k=Number(_[r.address]);await N.post("/vips/reals",{vip:s,real:{address:r.address,weight:k,flags:r.flags||0}}),await Z(),e("Real weight updated.","success")}catch(k){d(k.message||"request failed"),e(k.message||"Update failed.","error")}},pt=async r=>{try{await N.del("/vips/reals",{vip:s,real:{address:r.address,weight:r.weight,flags:r.flags||0}}),await Z(),e("Real removed.","success")}catch(k){d(k.message||"request failed"),e(k.message||"Remove failed.","error")}},mt=async r=>{r.preventDefault();try{await N.post("/vips/reals",{vip:s,real:{address:m.address,weight:Number(m.weight),flags:Number(m.flags||0)}}),v({address:"",weight:100,flags:0}),await Z(),e("Real added.","success")}catch(k){d(k.message||"request failed"),e(k.message||"Add failed.","error")}},ht=async()=>{if(!q||!b[q]){V("Select a target group to apply.");return}Re(!0),V("");let r=b[q]||[],k=new Map(c.map(w=>[w.address,w])),U=new Map(r.map(w=>[w.address,w])),x=c.filter(w=>!U.has(w.address)),D=r.filter(w=>{let T=k.get(w.address);return T?Number(T.weight)!==Number(w.weight)||Number(T.flags||0)!==Number(w.flags||0):!0});try{x.length>0&&await N.put("/vips/reals/batch",{vip:s,action:1,reals:x.map(w=>({address:w.address,weight:Number(w.weight),flags:Number(w.flags||0)}))}),D.length>0&&await Promise.all(D.map(w=>N.post("/vips/reals",{vip:s,real:{address:w.address,weight:Number(w.weight),flags:Number(w.flags||0)}}))),await Z(),e(`Applied target group "${q}".`,"success")}catch(w){V(w.message||"Failed to apply target group."),e(w.message||"Target group apply failed.","error")}finally{Re(!1)}},ft=r=>{r.preventDefault();let k=Fe.trim();if(!k){V("Provide a name for the new target group.");return}if(b[k]){V("A target group with that name already exists.");return}let U={...b,[k]:c.map(x=>({address:x.address,weight:Number(x.weight),flags:Number(x.flags||0)}))};C(U),Ae(""),ae(k),V(""),e(`Target group "${k}" saved.`,"success")},gt=async()=>{try{await N.del("/vips",s),e("VIP deleted.","success"),a("/")}catch(r){d(r.message||"request failed"),e(r.message||"Delete failed.","error")}},bt=async r=>{r.preventDefault();try{await N.put("/vips/flags",{...s,flag:Number(A.flags||0),set:!!A.set}),await pe(),e("VIP flags updated.","success")}catch(k){u(k.message||"request failed"),e(k.message||"Flag update failed.","error")}},vt=async r=>{r.preventDefault();try{await N.put("/vips/hash-function",{...s,hash_function:Number(i.hash_function)}),e("Hash function updated.","success")}catch(k){u(k.message||"request failed"),e(k.message||"Hash update failed.","error")}},$t=async r=>{r.preventDefault();let k=String(y.type||"").trim().toLowerCase();if(!["dummy","tcp","http","https"].includes(k)){j("Invalid healthcheck type.");return}let U=(D,w)=>{let T=Number(D);if(!Number.isFinite(T)||!Number.isInteger(T)||T<=0)throw new Error(`${w} must be a positive integer.`);return T},x={type:k};try{let D=String(y.port??"").trim();if(D){let w=Number(D);if(!Number.isFinite(w)||!Number.isInteger(w)||w<0||w>65535)throw new Error("Port must be an integer between 0 and 65535.");x.port=w}if(k!=="dummy"&&(x.interval_ms=U(y.interval_ms,"Interval"),x.timeout_ms=U(y.timeout_ms,"Timeout"),x.healthy_threshold=U(y.healthy_threshold,"Healthy threshold"),x.unhealthy_threshold=U(y.unhealthy_threshold,"Unhealthy threshold"),x.timeout_ms>=x.interval_ms))throw new Error("Timeout must be lower than interval.");if(k==="http"){let w=String(y.http_path||"").trim();if(!w)throw new Error("HTTP path is required.");x.http={path:w,expected_status:U(y.http_expected_status,"HTTP expected status")};let T=String(y.http_host||"").trim();T&&(x.http.host=T)}if(k==="https"){let w=String(y.https_path||"").trim();if(!w)throw new Error("HTTPS path is required.");x.https={path:w,expected_status:U(y.https_expected_status,"HTTPS expected status"),skip_tls_verify:!!y.https_skip_tls_verify};let T=String(y.https_host||"").trim();T&&(x.https.host=T)}}catch(D){j(D.message||"Invalid healthcheck configuration.");return}Ue(!0);try{await N.put("/vips/healthcheck",{vip:s,healthcheck:x}),await re(),e(se?"Healthcheck updated.":"Healthcheck configured.","success")}catch(D){j(D.message||"request failed"),e(D.message||"Healthcheck update failed.","error")}finally{Ue(!1)}},Nt=async()=>{He(!0);try{await N.del("/vips/healthcheck",s),await re(),e("Healthcheck removed.","success")}catch(r){j(r.message||"request failed"),e(r.message||"Healthcheck delete failed.","error")}finally{He(!1)}};return n`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>VIP Detail</h2>
              <p className="muted">${s.address}:${s.port} / ${s.proto}</p>
              ${S===null?n`<p className="muted">Flags: —</p>`:n`
                    <div style=${{marginTop:8}}>
                      <${$e}
                        mask=${S}
                        options=${Y}
                        showStatus=${!0}
                        emptyLabel="No flags"
                      />
                    </div>
                  `}
            </div>
            <div className="row">
              <button className="btn ghost" onClick=${pe}>Refresh flags</button>
              <button className="btn danger" onClick=${gt}>Delete VIP</button>
            </div>
          </div>
          ${o&&n`<p className="error">${o}</p>`}
          ${g&&n`<p className="error">${g}</p>`}
          ${f?n`<p className="muted">Loading reals…</p>`:n`
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
                    ${c.map(r=>n`
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
                              value=${_[r.address]??r.weight}
                              onInput=${k=>I({..._,[r.address]:k.target.value})}
                            />
                          </td>
                          <td className="row">
                            <button className="btn" onClick=${()=>ut(r)}>
                              Update
                            </button>
                            <button className="btn ghost" onClick=${()=>pt(r)}>
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
            <form className="form" onSubmit=${bt}>
              <div className="form-row">
                <label className="field">
                  <span>Flags</span>
                  <${Ne}
                    options=${Y}
                    value=${A.flags}
                    name="vip-flag-change"
                    onChange=${r=>L({...A,flags:r})}
                  />
                </label>
                <label className="field">
                  <span>Set</span>
                  <select
                    value=${String(A.set)}
                    onChange=${r=>L({...A,set:r.target.value==="true"})}
                  >
                    <option value="true">Set</option>
                    <option value="false">Clear</option>
                  </select>
                </label>
              </div>
              <button className="btn" type="submit">Apply flag</button>
            </form>
            <form className="form" onSubmit=${vt}>
              <div className="form-row">
                <label className="field">
                  <span>Hash function</span>
                  <input
                    type="number"
                    min="0"
                    value=${i.hash_function}
                    onInput=${r=>$({...i,hash_function:r.target.value})}
                  />
                </label>
              </div>
              <button className="btn" type="submit">Apply hash</button>
            </form>
          </div>
        </section>
        <section className="card">
          <div className="section-header">
            <div>
              <h3>VIP healthcheck (optional)</h3>
              <p className="muted">
                ${se?"A healthcheck is currently configured for this VIP.":"No healthcheck is configured for this VIP."}
              </p>
            </div>
            <div className="row">
              <button className="btn ghost" type="button" onClick=${re}>
                Refresh healthcheck
              </button>
              ${se?n`
                    <button
                      className="btn danger"
                      type="button"
                      onClick=${Nt}
                      disabled=${De}
                    >
                      ${De?"Removing...":"Remove healthcheck"}
                    </button>
                  `:null}
            </div>
          </div>
          ${Le?n`<p className="error">${Le}</p>`:null}
          ${dt?n`<p className="muted">Loading healthcheck configuration…</p>`:null}
          <form className="form" onSubmit=${$t}>
            <div className="form-row">
              <label className="field">
                <span>Type</span>
                <select
                  value=${y.type}
                  onChange=${r=>G({...y,type:r.target.value})}
                >
                  <option value="dummy">dummy</option>
                  <option value="tcp">tcp</option>
                  <option value="http">http</option>
                  <option value="https">https</option>
                </select>
              </label>
              <label className="field">
                <span>Port (optional)</span>
                <input
                  type="number"
                  min="0"
                  max="65535"
                  value=${y.port}
                  onInput=${r=>G({...y,port:r.target.value})}
                  placeholder=${`VIP port (${s.port})`}
                />
              </label>
            </div>

            ${y.type!=="dummy"?n`
                  <div className="form-row">
                    <label className="field">
                      <span>Interval (ms)</span>
                      <input
                        type="number"
                        min="1"
                        value=${y.interval_ms}
                        onInput=${r=>G({...y,interval_ms:r.target.value})}
                        required
                      />
                    </label>
                    <label className="field">
                      <span>Timeout (ms)</span>
                      <input
                        type="number"
                        min="1"
                        value=${y.timeout_ms}
                        onInput=${r=>G({...y,timeout_ms:r.target.value})}
                        required
                      />
                    </label>
                    <label className="field">
                      <span>Healthy threshold</span>
                      <input
                        type="number"
                        min="1"
                        value=${y.healthy_threshold}
                        onInput=${r=>G({...y,healthy_threshold:r.target.value})}
                        required
                      />
                    </label>
                    <label className="field">
                      <span>Unhealthy threshold</span>
                      <input
                        type="number"
                        min="1"
                        value=${y.unhealthy_threshold}
                        onInput=${r=>G({...y,unhealthy_threshold:r.target.value})}
                        required
                      />
                    </label>
                  </div>
                `:null}

            ${y.type==="http"?n`
                  <div className="form-row">
                    <label className="field">
                      <span>HTTP path</span>
                      <input
                        value=${y.http_path}
                        onInput=${r=>G({...y,http_path:r.target.value})}
                        placeholder="/healthz"
                        required
                      />
                    </label>
                    <label className="field">
                      <span>Expected status</span>
                      <input
                        type="number"
                        min="1"
                        value=${y.http_expected_status}
                        onInput=${r=>G({...y,http_expected_status:r.target.value})}
                        required
                      />
                    </label>
                    <label className="field">
                      <span>Host (optional)</span>
                      <input
                        value=${y.http_host}
                        onInput=${r=>G({...y,http_host:r.target.value})}
                        placeholder="example.com"
                      />
                    </label>
                  </div>
                `:null}

            ${y.type==="https"?n`
                  <div className="form-row">
                    <label className="field">
                      <span>HTTPS path</span>
                      <input
                        value=${y.https_path}
                        onInput=${r=>G({...y,https_path:r.target.value})}
                        placeholder="/healthz"
                        required
                      />
                    </label>
                    <label className="field">
                      <span>Expected status</span>
                      <input
                        type="number"
                        min="1"
                        value=${y.https_expected_status}
                        onInput=${r=>G({...y,https_expected_status:r.target.value})}
                        required
                      />
                    </label>
                    <label className="field">
                      <span>Host (optional)</span>
                      <input
                        value=${y.https_host}
                        onInput=${r=>G({...y,https_host:r.target.value})}
                        placeholder="example.com"
                      />
                    </label>
                    <label className="field">
                      <span>Skip TLS verify</span>
                      <input
                        type="checkbox"
                        checked=${!!y.https_skip_tls_verify}
                        onChange=${r=>G({...y,https_skip_tls_verify:r.target.checked})}
                      />
                    </label>
                  </div>
                `:null}

            <div className="row">
              <button className="btn" type="submit" disabled=${Ge}>
                ${Ge?"Saving...":se?"Update healthcheck":"Set healthcheck"}
              </button>
            </div>
          </form>
        </section>
        <section className="card">
          <h3>Add real</h3>
          <form className="form" onSubmit=${mt}>
            <div className="form-row">
              <label className="field">
                <span>Real address</span>
                <input
                  value=${m.address}
                  onInput=${r=>v({...m,address:r.target.value})}
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
                  onInput=${r=>v({...m,weight:r.target.value})}
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
              <button className="btn ghost" type="button" onClick=${P}>
                Reload groups
              </button>
              <button
                className="btn ghost"
                type="button"
                onClick=${async()=>{try{await R(),e("Imported target groups from running config.","success")}catch(r){V(r.message||"Failed to import target groups."),e(r.message||"Import failed.","error")}}}
              >
                Import from running config
              </button>
            </div>
          </div>
          ${Pe&&n`<p className="error">${Pe}</p>`}
          <div className="form-row">
            <label className="field">
              <span>Target group</span>
              <select
                value=${q}
                onChange=${r=>ae(r.target.value)}
                disabled=${Object.keys(b).length===0}
              >
                <option value="">Select group</option>
                ${Object.keys(b).map(r=>n`<option value=${r}>${r}</option>`)}
              </select>
            </label>
            <label className="field">
              <span>Preview</span>
              <input
                value=${`add ${de.add} \xB7 update ${de.update} \xB7 remove ${de.remove}`}
                readOnly
              />
            </label>
          </div>
          <div className="row">
            <button
              className="btn"
              type="button"
              onClick=${ht}
              disabled=${xe||!q}
            >
              ${xe?"Applying...":"Apply target group"}
            </button>
          </div>
          <form className="form" onSubmit=${ft}>
            <div className="form-row">
              <label className="field">
                <span>Save current reals as new group</span>
                <input
                  value=${Fe}
                  onInput=${r=>Ae(r.target.value)}
                  placeholder="edge-backends"
                />
              </label>
            </div>
            <button className="btn secondary" type="submit">Save target group</button>
          </form>
        </section>
      </main>
    `}function tt(){let e=be(),t=H(()=>ie(e.vipId),[e.vipId]),{points:a,error:s}=ee({path:"/stats/vip",body:t}),c=H(()=>te("/stats/vip"),[]),p=a[a.length-1]||{},o=a[a.length-2]||{},d=Number(p.v1??0),g=Number(p.v2??0),u=d-Number(o.v1??0),f=g-Number(o.v2??0),l=H(()=>[{label:c.v1,field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.15)"},{label:c.v2,field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[c]);return n`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>VIP Stats</h2>
              <p className="muted">${t.address}:${t.port} / ${t.proto}</p>
            </div>
          </div>
          ${s&&n`<p className="error">${s}</p>`}
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
                <td>${d}</td>
                <td>
                  <span className=${`delta ${u<0?"down":"up"}`}>
                    ${M(u)}
                  </span>
                </td>
              </tr>
              <tr>
                <td>${c.v2}</td>
                <td>${g}</td>
                <td>
                  <span className=${`delta ${f<0?"down":"up"}`}>
                    ${M(f)}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
          <${ce} title="Traffic (delta/sec)" points=${a} keys=${l} diff=${!0} />
        </section>
      </main>
    `}let at={"/stats/vip":{v1:"Packets",v2:"Bytes"},"/stats/real":{v1:"Packets",v2:"Bytes"},"/stats/lru":{v1:"Total packets",v2:"LRU hits"},"/stats/lru/miss":{v1:"TCP SYN misses",v2:"Non-SYN misses"},"/stats/lru/fallback":{v1:"Fallback LRU hits",v2:"Unused"},"/stats/lru/global":{v1:"Map lookup failures",v2:"Global LRU routed"},"/stats/xdp/total":{v1:"Packets",v2:"Bytes"},"/stats/xdp/pass":{v1:"Packets",v2:"Bytes"},"/stats/xdp/drop":{v1:"Packets",v2:"Bytes"},"/stats/xdp/tx":{v1:"Packets",v2:"Bytes"}},Ce=[{title:"LRU",path:"/stats/lru"},{title:"LRU Miss",path:"/stats/lru/miss"},{title:"LRU Fallback",path:"/stats/lru/fallback"},{title:"LRU Global",path:"/stats/lru/global"},{title:"XDP Total",path:"/stats/xdp/total"},{title:"XDP Pass",path:"/stats/xdp/pass"},{title:"XDP Drop",path:"/stats/xdp/drop"},{title:"XDP Tx",path:"/stats/xdp/tx"}];function M(e){return Number.isFinite(e)?`${e>0?"+":""}${e}`:"0"}function te(e){return at[e]||{v1:"v1",v2:"v2"}}function st({title:e,path:t,diff:a=!1}){let{points:s,error:c}=ee({path:t}),p=H(()=>te(t),[t]),o=H(()=>[{label:p.v1,field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.2)"},{label:p.v2,field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[p]);return n`
      <div className="card">
        <h3>${e}</h3>
        ${c&&n`<p className="error">${c}</p>`}
        <${ce} title=${e} points=${s} keys=${o} diff=${a} inlineTitle=${!1} />
      </div>
    `}function rt({title:e,path:t}){let{points:a,error:s}=ee({path:t}),c=H(()=>te(t),[t]),p=a[a.length-1]||{},o=a[a.length-2]||{},d=Number(p.v1??0),g=Number(p.v2??0),u=d-Number(o.v1??0),f=g-Number(o.v2??0);return n`
      <div className="summary-card">
        <div className="summary-title">${e}</div>
        ${s?n`<p className="error">${s}</p>`:n`
              <div className="summary-row">
                <div className="stat">
                  <span className="muted">${c.v1}</span>
                  <strong>${d}</strong>
                </div>
                <div className="stat">
                  <span className="muted">${c.v1} delta/sec</span>
                  <strong className=${u<0?"delta down":"delta up"}>
                    ${M(u)}
                  </strong>
                </div>
              </div>
              <div className="summary-row">
                <div className="stat">
                  <span className="muted">${c.v2}</span>
                  <strong>${g}</strong>
                </div>
                <div className="stat">
                  <span className="muted">${c.v2} delta/sec</span>
                  <strong className=${f<0?"delta down":"delta up"}>
                    ${M(f)}
                  </strong>
                </div>
              </div>
            `}
      </div>
    `}function nt(){let{data:e,error:t}=Ze(()=>N.get("/stats/userspace"),1e3,[]);return n`
      <main>
        <section className="card">
          <h2>Global Stats</h2>
          <p className="muted">Polling every second. Charts show per-second deltas.</p>
        </section>
        <section className="grid">
          ${Ce.map(a=>n`<${st} title=${a.title} path=${a.path} diff=${!0} />`)}
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
        <section className="card">
          <div className="section-header">
            <div>
              <h3>Absolute & Rate of Change</h3>
              <p className="muted">Latest value and per-second delta.</p>
            </div>
          </div>
          <div className="summary-grid">
            ${Ce.map(a=>n`<${rt} title=${a.title} path=${a.path} />`)}
          </div>
        </section>
      </main>
    `}function ot(){let[e,t]=h([]),[a,s]=h(""),[c,p]=h([]),[o,d]=h(""),[g,u]=h(null),[f,l]=h("");E(()=>{let b=!0;return(async()=>{try{let P=await N.get("/vips");if(!b)return;t(P||[]),!a&&P&&P.length>0&&s(Q(P[0]))}catch(P){b&&l(P.message||"request failed")}})(),()=>{b=!1}},[]),E(()=>{if(!a)return;let b=ie(a),C=!0;return(async()=>{try{let R=await N.get("/vips/reals",b);if(!C)return;p(R||[]),R&&R.length>0?d(q=>q||R[0].address):d(""),l("")}catch(R){C&&l(R.message||"request failed")}})(),()=>{C=!1}},[a]),E(()=>{if(!o){u(null);return}let b=!0;return(async()=>{try{let P=await N.get("/reals/index",{address:o});if(!b)return;u(P?.index??null),l("")}catch(P){b&&l(P.message||"request failed")}})(),()=>{b=!1}},[o]);let{points:m,error:v}=ee({path:"/stats/real",body:g!==null?{index:g}:null}),_=H(()=>te("/stats/real"),[]),I=H(()=>[{label:_.v1,field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.2)"},{label:_.v2,field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[_]),S=m[m.length-1]||{},F=m[m.length-2]||{},A=Number(S.v1??0),L=Number(S.v2??0),i=A-Number(F.v1??0),$=L-Number(F.v2??0);return n`
      <main>
        <section className="card">
          <h2>Per-Real Stats</h2>
          <p className="muted">Select a VIP and real address to chart.</p>
          ${f&&n`<p className="error">${f}</p>`}
          <div className="form-row">
            <label className="field">
              <span>VIP</span>
              <select value=${a} onChange=${b=>s(b.target.value)}>
                ${e.map(b=>n`
                    <option value=${Q(b)}>
                      ${b.address}:${b.port} / ${b.proto}
                    </option>
                  `)}
              </select>
            </label>
            <label className="field">
              <span>Real</span>
              <select
                value=${o}
                onChange=${b=>d(b.target.value)}
                disabled=${c.length===0}
              >
                ${c.map(b=>n`
                    <option value=${b.address}>${b.address}</option>
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
          ${v&&n`<p className="error">${v}</p>`}
          ${g===null?n`<p className="muted">Select a real to start polling.</p>`:n`
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
                      <td>${_.v1}</td>
                      <td>${A}</td>
                      <td>
                        <span className=${`delta ${i<0?"down":"up"}`}>
                          ${M(i)}
                        </span>
                      </td>
                    </tr>
                    <tr>
                      <td>${_.v2}</td>
                      <td>${L}</td>
                      <td>
                        <span className=${`delta ${$<0?"down":"up"}`}>
                          ${M($)}
                        </span>
                      </td>
                    </tr>
                  </tbody>
                </table>
                <${ce}
                  title="Traffic (delta/sec)"
                  points=${m}
                  keys=${I}
                  diff=${!0}
                />
              `}
        </section>
      </main>
    `}function lt(){let{addToast:e}=X(),[t,a]=h(""),[s,c]=h(""),[p,o]=h(!0),[d,g]=h(""),u=W(!0),f=async()=>{if(u.current){o(!0),c("");try{let m=await fetch(`${N.base}/config/export`,{credentials:"same-origin",headers:{Accept:"application/x-yaml"}});if(!m.ok){let _=`HTTP ${m.status}`,I="";try{let S=await m.json();_=S?.error?.message||_,I=S?.error?.code||""}catch{}throw(m.status===401||I==="UNAUTHORIZED")&&N.notifyUnauthorized(_),new z(_,m.status,I)}let v=await m.text();if(!u.current)return;a(v||""),g(new Date().toLocaleString())}catch(m){u.current&&c(m.message||"request failed")}finally{u.current&&o(!1)}}},l=async()=>{if(t)try{await navigator.clipboard.writeText(t),e("Config copied to clipboard","info")}catch{e("Failed to copy config","error")}};return E(()=>(u.current=!0,f(),()=>{u.current=!1}),[]),n`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>Running config</h2>
              <p className="muted">Exported from /api/v1/config/export</p>
            </div>
            <div className="row">
              <button className="btn ghost" onClick=${l} disabled=${!t}>
                Copy YAML
              </button>
              <button className="btn" onClick=${f} disabled=${p}>
                Refresh
              </button>
            </div>
          </div>
          ${s&&n`<p className="error">${s}</p>`}
          ${p?n`<p className="muted">Loading config...</p>`:t?n`<pre className="yaml-view">${t}</pre>`:n`<p className="muted">No config data returned.</p>`}
          ${d&&n`<p className="muted">Last fetched ${d}</p>`}
        </section>
      </main>
    `}function it(){let{addToast:e}=X(),{groups:t,setGroups:a,refreshFromStorage:s,importFromRunningConfig:c}=Se(),[p,o]=h(""),[d,g]=h(""),[u,f]=h({address:"",weight:100,flags:0}),[l,m]=h(""),[v,_]=h(!1);E(()=>{if(d){if(!t[d]){let $=Object.keys(t);g($[0]||"")}}else{let $=Object.keys(t);$.length>0&&g($[0])}},[t,d]);let I=$=>{$.preventDefault();let b=p.trim();if(!b){m("Provide a group name.");return}if(t[b]){m("That group already exists.");return}a({...t,[b]:[]}),o(""),g(b),m(""),e(`Target group "${b}" created.`,"success")},S=$=>{let b={...t};delete b[$],a(b),e(`Target group "${$}" removed.`,"success")},F=$=>{if($.preventDefault(),!d){m("Select a group to add a real.");return}let b=_e(u);if(!b){m("Provide a valid real address.");return}let C=t[d]||[],P=C.some(R=>R.address===b.address)?C.map(R=>R.address===b.address?b:R):C.concat(b);a({...t,[d]:P}),f({address:"",weight:100,flags:0}),m(""),e("Real saved to target group.","success")},A=$=>{if(!d)return;let C=(t[d]||[]).filter(P=>P.address!==$);a({...t,[d]:C})},L=($,b)=>{if(!d)return;let P=(t[d]||[]).map(R=>R.address===$?{...R,...b}:R);a({...t,[d]:P})};return n`
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
              <button className="btn ghost" type="button" onClick=${async()=>{_(!0);try{await c(),e("Imported target groups from running config.","success"),m("")}catch($){m($.message||"Failed to import target groups."),e($.message||"Import failed.","error")}finally{_(!1)}}} disabled=${v}>
                ${v?"Importing...":"Import from running config"}
              </button>
            </div>
          </div>
          ${l&&n`<p className="error">${l}</p>`}
          <form className="form" onSubmit=${I}>
            <div className="form-row">
              <label className="field">
                <span>New group name</span>
                <input
                  value=${p}
                  onInput=${$=>o($.target.value)}
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
                  value=${d}
                  onChange=${$=>g($.target.value)}
                  disabled=${Object.keys(t).length===0}
                >
                  ${Object.keys(t).map($=>n`<option value=${$}>${$}</option>`)}
                </select>
              </label>
              ${d&&n`<button className="btn danger" type="button" onClick=${()=>S(d)}>
                Delete group
              </button>`}
            </div>
          </div>
          ${d?n`
                <table className="table">
                  <thead>
                    <tr>
                      <th>Address</th>
                      <th>Weight</th>
                      <th></th>
                    </tr>
                  </thead>
                  <tbody>
                    ${(t[d]||[]).map($=>n`
                        <tr>
                          <td>${$.address}</td>
                          <td>
                            <input
                              className="inline-input"
                              type="number"
                              min="0"
                              value=${$.weight}
                              onInput=${b=>L($.address,{weight:Number(b.target.value)})}
                            />
                          </td>
                          <td>
                            <button
                              className="btn ghost"
                              type="button"
                              onClick=${()=>A($.address)}
                            >
                              Remove
                            </button>
                          </td>
                        </tr>
                      `)}
                  </tbody>
                </table>
                <form className="form" onSubmit=${F}>
                  <div className="form-row">
                    <label className="field">
                      <span>Real address</span>
                      <input
                        value=${u.address}
                        onInput=${$=>f({...u,address:$.target.value})}
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
                        onInput=${$=>f({...u,weight:$.target.value})}
                        required
                      />
                    </label>
                  </div>
                  <button className="btn secondary" type="submit">Save real</button>
                </form>
              `:n`<p className="muted">No groups yet. Create one to add reals.</p>`}
        </section>
      </main>
    `}function ct(){let[e,t]=h({initialized:!1,ready:!1}),[a,s]=h([]),[c,p]=h({checking:!0,required:!1,username:""}),o=W({}),d=Oe(async()=>{try{let v=await N.get("/lb/status");t(v||{initialized:!1,ready:!1}),p(_=>({..._,checking:!1,required:!1}))}catch(v){if(v instanceof z&&v.unauthorized){t({initialized:!1,ready:!1}),p(_=>({..._,checking:!1,required:!0}));return}t({initialized:!1,ready:!1}),p(_=>({..._,checking:!1}))}},[]),g=(v,_="info")=>{let I=`${Date.now()}-${Math.random().toString(16).slice(2)}`;s(S=>S.concat({id:I,message:v,kind:_})),o.current[I]=setTimeout(()=>{s(S=>S.filter(F=>F.id!==I)),delete o.current[I]},4e3)},u=v=>{o.current[v]&&(clearTimeout(o.current[v]),delete o.current[v]),s(_=>_.filter(I=>I.id!==v))},f=async(v,_)=>{let I=await N.login(v,_);p({checking:!1,required:!1,username:I?.username||v}),await d(),g("Signed in.","success")},l=async()=>{try{await N.logout()}catch(v){g(v.message||"Sign out failed.","error")}t({initialized:!1,ready:!1}),p({checking:!1,required:!0,username:""})};E(()=>(N.setAuthHandlers({onUnauthorized:()=>{p(v=>({...v,checking:!1,required:!0,username:""}))}}),()=>{N.setAuthHandlers({})}),[]),E(()=>{let v=!0;(async()=>{v&&await d()})();let I=setInterval(()=>{!v||c.required||d()},5e3);return()=>{v=!1,clearInterval(I)}},[d,c.required]);let m=H(()=>({required:c.required,username:c.username,login:f,logout:l}),[c.required,c.username]);return E(()=>()=>{Object.keys(o.current).forEach(v=>clearTimeout(o.current[v]))},[]),n`
      <${fe}>
        <${oe.Provider} value=${m}>
          <${Te} checking=${c.checking} required=${c.required}>
            <${ne.Provider} value=${{addToast:g}}>
              <${Xe} status=${e} />
              <${ge}>
                <${B} path="/" element=${n`<${Qe} />`} />
                <${B} path="/vips/:vipId" element=${n`<${et} />`} />
                <${B} path="/vips/:vipId/stats" element=${n`<${tt} />`} />
                <${B} path="/target-groups" element=${n`<${it} />`} />
                <${B} path="/stats/global" element=${n`<${nt} />`} />
                <${B} path="/stats/real" element=${n`<${ot} />`} />
                <${B} path="/config" element=${n`<${lt} />`} />
              </${ge}>
              <${Ye} toasts=${a} onDismiss=${u} />
            </${ne.Provider}>
          </${Te}>
        </${oe.Provider}>
      </${fe}>
    `}ReactDOM.createRoot(document.getElementById("root")).render(n`<${ct} />`)})();})();
