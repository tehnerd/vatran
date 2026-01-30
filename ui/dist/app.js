(()=>{(()=>{let{useEffect:G,useMemo:E,useRef:U,useState:h,useContext:Ie}=React,{BrowserRouter:re,Routes:ne,Route:O,NavLink:L,Link:j,useParams:oe,useNavigate:Se}=ReactRouterDOM,n=htm.bind(React.createElement),Z=React.createContext({addToast:()=>{}}),H=[{label:"NO_SPORT",value:1},{label:"NO_LRU",value:2},{label:"QUIC_VIP",value:4},{label:"DPORT_HASH",value:8},{label:"SRC_ROUTING",value:16},{label:"LOCAL_VIP",value:32},{label:"GLOBAL_LRU",value:64},{label:"HASH_SRC_DST_PORT",value:128},{label:"UDP_STABLE_ROUTING_VIP",value:256},{label:"UDP_FLOW_MIGRATION",value:512}];function W(){return Ie(Z)}function Re(e){return String(e).replace(/[^a-z0-9_-]/gi,"_")}function Ce(e,t,a){let o=Number(e)||0,d=Number(t)||0;return a?o|d:o&~d}function Te(e,t){let a=Number(e)||0;return t.filter(o=>(a&o.value)!==0)}function le({mask:e,options:t,showStatus:a=!1,emptyLabel:o="None"}){let d=Number(e)||0,b=a?t:Te(d,t),i=a?2:1;return n`
      <table className="table flag-table">
        <thead>
          <tr>
            <th>Flag</th>
            ${a?n`<th>Enabled</th>`:null}
          </tr>
        </thead>
        <tbody>
          ${b.length===0?n`<tr><td colspan=${i} className="muted">${o}</td></tr>`:b.map(s=>{let p=(d&s.value)!==0;return n`
                  <tr>
                    <td>${s.label}</td>
                    ${a?n`<td>${p?"Yes":"No"}</td>`:null}
                  </tr>
                `})}
        </tbody>
      </table>
    `}function ie({options:e,value:t,onChange:a,name:o}){let d=Number(t)||0,b=Re(o||"flags");return n`
      <div className="flag-selector">
        <table className="table flag-table">
          <thead>
            <tr>
              <th>Enabled</th>
              <th>Flag</th>
            </tr>
          </thead>
          <tbody>
            ${e.map(i=>{let s=`${b}-${i.value}`,p=(d&i.value)===i.value;return n`
                <tr>
                  <td>
                    <input
                      id=${s}
                      type="checkbox"
                      checked=${p}
                      onChange=${c=>a(Ce(d,i.value,c.target.checked))}
                    />
                  </td>
                  <td>
                    <label className="flag-option" htmlFor=${s}>${i.label}</label>
                  </td>
                </tr>
              `})}
          </tbody>
        </table>
      </div>
    `}let N={base:"/api/v1",async request(e,t={}){let a={method:t.method||"GET",headers:{"Content-Type":"application/json"}},o=`${N.base}${e}`;if(t.body!==void 0&&t.body!==null)if(a.method==="GET"){let i=new URLSearchParams;Object.entries(t.body).forEach(([p,c])=>{if(c!=null){if(Array.isArray(c)){c.forEach(m=>i.append(p,String(m)));return}if(typeof c=="object"){i.set(p,JSON.stringify(c));return}i.set(p,String(c))}});let s=i.toString();s&&(o+=`${o.includes("?")?"&":"?"}${s}`)}else a.body=JSON.stringify(t.body);let d=await fetch(o,a),b;try{b=await d.json()}catch{throw new Error("invalid JSON response")}if(!d.ok)throw new Error(b?.error?.message||`HTTP ${d.status}`);if(!b.success){let i=b.error?.message||"request failed";throw new Error(i)}return b.data},get(e,t){return N.request(e,{method:"GET",body:t})},post(e,t){return N.request(e,{method:"POST",body:t})},put(e,t){return N.request(e,{method:"PUT",body:t})},del(e,t){return N.request(e,{method:"DELETE",body:t})}},ce="vatran_target_groups";function de(e){if(!e||!e.address)return null;let t=String(e.address).trim();if(!t)return null;let a=Number(e.weight),o=Number(e.flags??0);return{address:t,weight:Number.isFinite(a)?a:0,flags:Number.isFinite(o)?o:0}}function ue(e){if(!e||typeof e!="object")return{};let t={};return Object.entries(e).forEach(([a,o])=>{let d=String(a).trim();if(!d)return;let b=Array.isArray(o)?o.map(de).filter(Boolean):[],i=[],s=new Set;b.forEach(p=>{s.has(p.address)||(s.add(p.address),i.push(p))}),t[d]=i}),t}function X(){if(typeof localStorage>"u")return{};try{let e=localStorage.getItem(ce);return e?ue(JSON.parse(e)):{}}catch{return{}}}function pe(e){if(!(typeof localStorage>"u"))try{localStorage.setItem(ce,JSON.stringify(e))}catch{}}function Pe(e,t){let a={...e};return Object.entries(t||{}).forEach(([o,d])=>{a[o]||(a[o]=d)}),a}function me(){let[e,t]=h(()=>X());return G(()=>{pe(e)},[e]),{groups:e,setGroups:t,refreshFromStorage:()=>{t(X())},importFromRunningConfig:async()=>{let d=await N.get("/config/export/json"),b=ue(d?.target_groups||{}),i=Pe(X(),b);return t(i),pe(i),i}}}function J(e){return`${encodeURIComponent(e.address)}:${e.port}:${e.proto}`}function Q(e){let t=e.split(":"),a=Number(t.pop()||0),o=Number(t.pop()||0);return{address:decodeURIComponent(t.join(":")),port:o,proto:a}}function ke(e,t,a=[]){let[o,d]=h(null),[b,i]=h(""),[s,p]=h(!0);return G(()=>{let c=!0,m=async()=>{try{let g=await e();c&&(d(g),i(""),p(!1))}catch(g){c&&(i(g.message||"request failed"),p(!1))}};m();let r=setInterval(m,t);return()=>{c=!1,clearInterval(r)}},a),{data:o,error:b,loading:s}}function K({path:e,body:t,intervalMs:a=1e3,limit:o=60}){let[d,b]=h([]),[i,s]=h(""),p=E(()=>JSON.stringify(t||{}),[t]);return G(()=>{if(t===null)return b([]),s(""),()=>{};let c=!0,m=async()=>{try{let g=await N.get(e,t);if(!c)return;let y=new Date().toLocaleTimeString();b(w=>w.concat({label:y,...g}).slice(-o)),s("")}catch(g){c&&s(g.message||"request failed")}};m();let r=setInterval(m,a);return()=>{c=!1,clearInterval(r)}},[e,p,a,o]),{points:d,error:i}}function ge({title:e,points:t,keys:a,diff:o=!1,height:d=120,showTitle:b=!1,selectedLabel:i=null,onPointSelect:s=null,onLegendSelect:p=null}){let c=U(null),m=U(null);return G(()=>{if(!c.current)return;m.current||(m.current=new Chart(c.current,{type:"line",data:{labels:[],datasets:[]},options:{responsive:!0,maintainAspectRatio:!1,animation:!1,scales:{x:{grid:{display:!1}},y:{beginAtZero:!o}},plugins:{legend:{display:!0,position:"bottom"},title:{display:b&&!!e,text:e}}}}));let r=m.current,g=new Map((r.data.datasets||[]).filter(_=>typeof _.hidden<"u").map(_=>[_.label,_.hidden])),y=t.map(_=>_.label),w=i?a.filter(_=>_.label===i):a,T=i&&w.length===0?a:w;return r.data.labels=y,r.data.datasets=T.map(_=>{let P=t.map(x=>x[_.field]||0),R=o?P.map((x,l)=>l===0?0:x-P[l-1]):P;return{label:_.label,data:R,borderColor:_.color,backgroundColor:_.fill,borderWidth:2,tension:.3,hidden:g.get(_.label)}}),r.options.onClick=(_,P)=>{if(!s||!P||P.length===0)return;let R=P[0].datasetIndex,x=r.data.datasets?.[R]?.label;x&&s(x)},r.options.plugins&&r.options.plugins.legend&&(r.options.plugins.legend.onClick=(_,P)=>{if(!p)return;let R=P?.text;R&&p(R)}),r.options.scales.y.beginAtZero=!o,r.options.plugins.title.display=b&&!!e,r.options.plugins.title.text=e||"",r.update(),()=>{}},[t,a,e,o,b,i,s,p]),G(()=>()=>{m.current&&(m.current.destroy(),m.current=null)},[]),n`<canvas ref=${c} height=${d}></canvas>`}function ee({title:e,points:t,keys:a,diff:o=!1,inlineTitle:d=!0}){let[b,i]=h(!1),[s,p]=h(null),c=U(!1);return n`
      <div className="chart-wrap">
        <button className="chart-zoom-button" type="button" onClick=${()=>i(!0)}>
          <span className="zoom-icon" aria-hidden="true">+</span>
          Zoom
        </button>
        <div
          className="chart-click"
          onClick=${()=>{if(c.current){c.current=!1;return}i(!0)}}
        >
          <${ge}
            title=${e}
            points=${t}
            keys=${a}
            diff=${o}
            height=${120}
            showTitle=${d&&!!e}
            selectedLabel=${s}
            onPointSelect=${m=>{p(r=>r===m?null:m),c.current=!0,setTimeout(()=>{c.current=!1},0)}}
            onLegendSelect=${m=>{p(r=>r===m?null:m),c.current=!0,setTimeout(()=>{c.current=!1},0)}}
          />
        </div>
        ${b&&n`
          <div className="chart-overlay" onClick=${()=>i(!1)}>
            <div className="chart-modal" onClick=${m=>m.stopPropagation()}>
              <div className="row chart-modal-header">
                <div>
                  <h3>${e||"Chart"}</h3>
                  ${o?n`<p className="muted">Per-second delta.</p>`:""}
                </div>
                <button className="btn ghost" onClick=${()=>i(!1)}>
                  Close
                </button>
              </div>
              <div className="chart-zoom">
                <${ge}
                  title=${e}
                  points=${t}
                  keys=${a}
                  diff=${o}
                  height=${360}
                  showTitle=${!1}
                  selectedLabel=${s}
                  onPointSelect=${m=>p(r=>r===m?null:m)}
                  onLegendSelect=${m=>p(r=>r===m?null:m)}
                />
              </div>
            </div>
          </div>
        `}
      </div>
    `}function fe({children:e}){return e}function xe({toasts:e,onDismiss:t}){return n`
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
    `}function Fe({status:e}){return n`
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
    `}function Ge(){let{addToast:e}=W(),[t,a]=h({initialized:!1,ready:!1}),[o,d]=h([]),[b,i]=h(""),[s,p]=h(!1),[c,m]=h(!1),[r,g]=h({main_interface:"",balancer_prog_path:"",healthchecking_prog_path:"",default_mac:"",local_mac:"",root_map_path:"",root_map_pos:2,katran_src_v4:"",katran_src_v6:"",use_root_map:!1,max_vips:512,max_reals:4096,hash_func:0}),[y,w]=h({address:"",port:80,proto:6,flags:0}),T=async()=>{try{let l=await N.get("/lb/status"),v=await N.get("/vips"),f=await Promise.all((v||[]).map(async C=>{try{let S=await N.get("/vips/flags",{address:C.address,port:C.port,proto:C.proto});return{...C,flags:S?.flags??0}}catch{return{...C,flags:null}}}));a(l||{initialized:!1,ready:!1}),d(f),i("")}catch(l){i(l.message||"request failed")}};G(()=>{let l=!0;return(async()=>{l&&await T()})(),()=>{l=!1}},[]);let _=async l=>{l.preventDefault();try{let v={...r,root_map_pos:r.root_map_pos===""?void 0:Number(r.root_map_pos),max_vips:Number(r.max_vips),max_reals:Number(r.max_reals),hash_func:Number(r.hash_func)};await N.post("/lb/create",v),i(""),p(!1),e("Load balancer initialized.","success"),await T()}catch(v){i(v.message||"request failed"),e(v.message||"Initialize failed.","error")}},P=async l=>{l.preventDefault();try{await N.post("/vips",{...y,port:Number(y.port),proto:Number(y.proto),flags:Number(y.flags||0)}),w({address:"",port:80,proto:6,flags:0}),i(""),m(!1),e("VIP created.","success"),await T()}catch(v){i(v.message||"request failed"),e(v.message||"VIP create failed.","error")}},R=async()=>{try{await N.post("/lb/load-bpf-progs"),i(""),e("BPF programs loaded.","success"),await T()}catch(l){i(l.message||"request failed"),e(l.message||"Load BPF programs failed.","error")}},x=async()=>{try{await N.post("/lb/attach-bpf-progs"),i(""),e("BPF programs attached.","success"),await T()}catch(l){i(l.message||"request failed"),e(l.message||"Attach BPF programs failed.","error")}};return n`
      <main>
        <section className="card">
          <h2>Load balancer</h2>
          <p>Initialized: ${t.initialized?"yes":"no"}</p>
          <p>Ready: ${t.ready?"yes":"no"}</p>
          <div className="row">
            ${!t.initialized&&n`
              <button className="btn" onClick=${()=>p(l=>!l)}>
                ${s?"Close":"Initialize"}
              </button>
            `}
            <button className="btn secondary" onClick=${()=>m(l=>!l)}>
              ${c?"Close":"Create VIP"}
            </button>
          </div>
          ${!t.ready&&n`
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
          ${s&&n`
            <form className="form" onSubmit=${_}>
              <div className="form-row">
                <label className="field">
                  <span>Main interface</span>
                  <input
                    value=${r.main_interface}
                    onInput=${l=>g({...r,main_interface:l.target.value})}
                    placeholder="eth0"
                    required
                  />
                </label>
                <label className="field">
                  <span>Balancer prog path</span>
                  <input
                    value=${r.balancer_prog_path}
                    onInput=${l=>g({...r,balancer_prog_path:l.target.value})}
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
                    onInput=${l=>g({...r,healthchecking_prog_path:l.target.value})}
                    placeholder="/path/to/hc.o"
                  />
                </label>
                <label className="field">
                  <span>Default MAC</span>
                  <input
                    value=${r.default_mac}
                    onInput=${l=>g({...r,default_mac:l.target.value})}
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
                    onInput=${l=>g({...r,local_mac:l.target.value})}
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
                    onInput=${l=>g({...r,hash_func:l.target.value})}
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Root map path</span>
                  <input
                    value=${r.root_map_path}
                    onInput=${l=>g({...r,root_map_path:l.target.value})}
                    placeholder="/sys/fs/bpf/root_map"
                  />
                </label>
                <label className="field">
                  <span>Root map position</span>
                  <input
                    type="number"
                    min="0"
                    value=${r.root_map_pos}
                    onInput=${l=>g({...r,root_map_pos:l.target.value})}
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Katran src v4</span>
                  <input
                    value=${r.katran_src_v4}
                    onInput=${l=>g({...r,katran_src_v4:l.target.value})}
                    placeholder="10.0.0.1"
                  />
                </label>
                <label className="field">
                  <span>Katran src v6</span>
                  <input
                    value=${r.katran_src_v6}
                    onInput=${l=>g({...r,katran_src_v6:l.target.value})}
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
                    onInput=${l=>g({...r,max_vips:l.target.value})}
                  />
                </label>
                <label className="field">
                  <span>Max reals</span>
                  <input
                    type="number"
                    min="1"
                    value=${r.max_reals}
                    onInput=${l=>g({...r,max_reals:l.target.value})}
                  />
                </label>
              </div>
              <label className="field checkbox">
                <input
                  type="checkbox"
                  checked=${r.use_root_map}
                  onChange=${l=>g({...r,use_root_map:l.target.checked})}
                />
                <span>Use root map</span>
              </label>
              <button className="btn" type="submit">Initialize LB</button>
            </form>
          `}
          ${c&&n`
            <form className="form" onSubmit=${P}>
              <div className="form-row">
                <label className="field">
                  <span>VIP address</span>
                  <input
                    value=${y.address}
                    onInput=${l=>w({...y,address:l.target.value})}
                    placeholder="1.2.3.4"
                    required
                  />
                </label>
                <label className="field">
                  <span>Port</span>
                  <input
                    type="number"
                    value=${y.port}
                    onInput=${l=>w({...y,port:l.target.value})}
                    required
                  />
                </label>
                <label className="field">
                  <span>Protocol</span>
                  <select
                    value=${y.proto}
                    onChange=${l=>w({...y,proto:l.target.value})}
                  >
                    <option value="6">TCP (6)</option>
                    <option value="17">UDP (17)</option>
                  </select>
                </label>
                <label className="field">
                  <span>Flags</span>
                  <${ie}
                    options=${H}
                    value=${y.flags}
                    name="vip-add"
                    onChange=${l=>w({...y,flags:l})}
                  />
                </label>
              </div>
              <button className="btn" type="submit">Create VIP</button>
            </form>
          `}
          ${b&&n`<p className="error">${b}</p>`}
        </section>
        <section className="card">
          <div className="section-header">
            <h2>VIPs</h2>
            <button className="btn ghost" onClick=${T}>Refresh</button>
          </div>
          ${o.length===0?n`<p className="muted">No VIPs configured yet.</p>`:n`
                <div className="grid">
                  ${o.map(l=>n`
                      <div className="card">
                        <div style=${{fontWeight:600}}>
                          ${l.address}:${l.port} / ${l.proto}
                        </div>
                        <div className="muted" style=${{marginTop:6}}>
                          <div style=${{fontWeight:600,marginBottom:6}}>Flags</div>
                          <${le}
                            mask=${l.flags}
                            options=${H}
                            emptyLabel=${l.flags===null?"Unknown":"No flags"}
                          />
                        </div>
                        <div className="row" style=${{marginTop:12}}>
                          <${j} className="btn" to=${`/vips/${J(l)}`}>
                            Open
                          </${j}>
                          <${j}
                            className="btn secondary"
                            to=${`/vips/${J(l)}/stats`}
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
    `}function Ae(){let{addToast:e}=W(),t=oe(),a=Se(),o=E(()=>Q(t.vipId),[t.vipId]),[d,b]=h([]),[i,s]=h(""),[p,c]=h(""),[m,r]=h(!0),[g,y]=h({address:"",weight:100,flags:0}),[w,T]=h({}),[_,P]=h(null),[R,x]=h({flags:0,set:!0}),[l,v]=h({hash_function:0}),{groups:f,setGroups:C,refreshFromStorage:S,importFromRunningConfig:k}=me(),[A,ve]=h(""),[he,B]=h(""),[$e,Ne]=h(!1),[ye,we]=h(""),[te,_e]=h({add:0,update:0,remove:0}),z=async()=>{try{let u=await N.get("/vips/reals",o);b(u||[]);let $={};(u||[]).forEach(q=>{$[q.address]=q.weight}),T($),s(""),r(!1)}catch(u){s(u.message||"request failed"),r(!1)}},ae=async()=>{try{let u=await N.get("/vips/flags",o);P(u?.flags??0),c("")}catch(u){c(u.message||"request failed")}};G(()=>{z(),ae()},[t.vipId]),G(()=>{if(!A){_e({add:0,update:0,remove:0});return}let u=f[A]||[],$=new Map(d.map(F=>[F.address,F])),q=new Map(u.map(F=>[F.address,F])),D=0,M=0,I=0;u.forEach(F=>{let se=$.get(F.address);if(!se){D+=1;return}(Number(se.weight)!==Number(F.weight)||Number(se.flags||0)!==Number(F.flags||0))&&(M+=1)}),d.forEach(F=>{q.has(F.address)||(I+=1)}),_e({add:D,update:M,remove:I})},[A,d,f]);let Me=async u=>{try{let $=Number(w[u.address]);await N.post("/vips/reals",{vip:o,real:{address:u.address,weight:$,flags:u.flags||0}}),await z(),e("Real weight updated.","success")}catch($){s($.message||"request failed"),e($.message||"Update failed.","error")}},je=async u=>{try{await N.del("/vips/reals",{vip:o,real:{address:u.address,weight:u.weight,flags:u.flags||0}}),await z(),e("Real removed.","success")}catch($){s($.message||"request failed"),e($.message||"Remove failed.","error")}},He=async u=>{u.preventDefault();try{await N.post("/vips/reals",{vip:o,real:{address:g.address,weight:Number(g.weight),flags:Number(g.flags||0)}}),y({address:"",weight:100,flags:0}),await z(),e("Real added.","success")}catch($){s($.message||"request failed"),e($.message||"Add failed.","error")}},We=async()=>{if(!A||!f[A]){B("Select a target group to apply.");return}Ne(!0),B("");let u=f[A]||[],$=new Map(d.map(I=>[I.address,I])),q=new Map(u.map(I=>[I.address,I])),D=d.filter(I=>!q.has(I.address)),M=u.filter(I=>{let F=$.get(I.address);return F?Number(F.weight)!==Number(I.weight)||Number(F.flags||0)!==Number(I.flags||0):!0});try{D.length>0&&await N.put("/vips/reals/batch",{vip:o,action:1,reals:D.map(I=>({address:I.address,weight:Number(I.weight),flags:Number(I.flags||0)}))}),M.length>0&&await Promise.all(M.map(I=>N.post("/vips/reals",{vip:o,real:{address:I.address,weight:Number(I.weight),flags:Number(I.flags||0)}}))),await z(),e(`Applied target group "${A}".`,"success")}catch(I){B(I.message||"Failed to apply target group."),e(I.message||"Target group apply failed.","error")}finally{Ne(!1)}},Je=u=>{u.preventDefault();let $=ye.trim();if(!$){B("Provide a name for the new target group.");return}if(f[$]){B("A target group with that name already exists.");return}let q={...f,[$]:d.map(D=>({address:D.address,weight:Number(D.weight),flags:Number(D.flags||0)}))};C(q),we(""),ve($),B(""),e(`Target group "${$}" saved.`,"success")},Ke=async()=>{try{await N.del("/vips",o),e("VIP deleted.","success"),a("/")}catch(u){s(u.message||"request failed"),e(u.message||"Delete failed.","error")}},Ye=async u=>{u.preventDefault();try{await N.put("/vips/flags",{...o,flag:Number(R.flags||0),set:!!R.set}),await ae(),e("VIP flags updated.","success")}catch($){c($.message||"request failed"),e($.message||"Flag update failed.","error")}},Ze=async u=>{u.preventDefault();try{await N.put("/vips/hash-function",{...o,hash_function:Number(l.hash_function)}),e("Hash function updated.","success")}catch($){c($.message||"request failed"),e($.message||"Hash update failed.","error")}};return n`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>VIP Detail</h2>
              <p className="muted">${o.address}:${o.port} / ${o.proto}</p>
              ${_===null?n`<p className="muted">Flags: —</p>`:n`
                    <div style=${{marginTop:8}}>
                      <${le}
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
              <button className="btn danger" onClick=${Ke}>Delete VIP</button>
            </div>
          </div>
          ${i&&n`<p className="error">${i}</p>`}
          ${p&&n`<p className="error">${p}</p>`}
          ${m?n`<p className="muted">Loading reals…</p>`:n`
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
                    ${d.map(u=>n`
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
                            <button className="btn" onClick=${()=>Me(u)}>
                              Update
                            </button>
                            <button className="btn ghost" onClick=${()=>je(u)}>
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
            <form className="form" onSubmit=${Ye}>
              <div className="form-row">
                <label className="field">
                  <span>Flags</span>
                  <${ie}
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
            <form className="form" onSubmit=${Ze}>
              <div className="form-row">
                <label className="field">
                  <span>Hash function</span>
                  <input
                    type="number"
                    min="0"
                    value=${l.hash_function}
                    onInput=${u=>v({...l,hash_function:u.target.value})}
                  />
                </label>
              </div>
              <button className="btn" type="submit">Apply hash</button>
            </form>
          </div>
        </section>
        <section className="card">
          <h3>Add real</h3>
          <form className="form" onSubmit=${He}>
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
              <button className="btn ghost" type="button" onClick=${S}>
                Reload groups
              </button>
              <button
                className="btn ghost"
                type="button"
                onClick=${async()=>{try{await k(),e("Imported target groups from running config.","success")}catch(u){B(u.message||"Failed to import target groups."),e(u.message||"Import failed.","error")}}}
              >
                Import from running config
              </button>
            </div>
          </div>
          ${he&&n`<p className="error">${he}</p>`}
          <div className="form-row">
            <label className="field">
              <span>Target group</span>
              <select
                value=${A}
                onChange=${u=>ve(u.target.value)}
                disabled=${Object.keys(f).length===0}
              >
                <option value="">Select group</option>
                ${Object.keys(f).map(u=>n`<option value=${u}>${u}</option>`)}
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
              onClick=${We}
              disabled=${$e||!A}
            >
              ${$e?"Applying...":"Apply target group"}
            </button>
          </div>
          <form className="form" onSubmit=${Je}>
            <div className="form-row">
              <label className="field">
                <span>Save current reals as new group</span>
                <input
                  value=${ye}
                  onInput=${u=>we(u.target.value)}
                  placeholder="edge-backends"
                />
              </label>
            </div>
            <button className="btn secondary" type="submit">Save target group</button>
          </form>
        </section>
      </main>
    `}function Ee(){let e=oe(),t=E(()=>Q(e.vipId),[e.vipId]),{points:a,error:o}=K({path:"/stats/vip",body:t}),d=E(()=>Y("/stats/vip"),[]),b=a[a.length-1]||{},i=a[a.length-2]||{},s=Number(b.v1??0),p=Number(b.v2??0),c=s-Number(i.v1??0),m=p-Number(i.v2??0),r=E(()=>[{label:d.v1,field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.15)"},{label:d.v2,field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[d]);return n`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>VIP Stats</h2>
              <p className="muted">${t.address}:${t.port} / ${t.proto}</p>
            </div>
          </div>
          ${o&&n`<p className="error">${o}</p>`}
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
                <td>${d.v1}</td>
                <td>${s}</td>
                <td>
                  <span className=${`delta ${c<0?"down":"up"}`}>
                    ${V(c)}
                  </span>
                </td>
              </tr>
              <tr>
                <td>${d.v2}</td>
                <td>${p}</td>
                <td>
                  <span className=${`delta ${m<0?"down":"up"}`}>
                    ${V(m)}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
          <${ee} title="Traffic (delta/sec)" points=${a} keys=${r} diff=${!0} />
        </section>
      </main>
    `}let Le={"/stats/vip":{v1:"Packets",v2:"Bytes"},"/stats/real":{v1:"Packets",v2:"Bytes"},"/stats/lru":{v1:"Total packets",v2:"LRU hits"},"/stats/lru/miss":{v1:"TCP SYN misses",v2:"Non-SYN misses"},"/stats/lru/fallback":{v1:"Fallback LRU hits",v2:"Unused"},"/stats/lru/global":{v1:"Map lookup failures",v2:"Global LRU routed"},"/stats/xdp/total":{v1:"Packets",v2:"Bytes"},"/stats/xdp/pass":{v1:"Packets",v2:"Bytes"},"/stats/xdp/drop":{v1:"Packets",v2:"Bytes"},"/stats/xdp/tx":{v1:"Packets",v2:"Bytes"}},be=[{title:"LRU",path:"/stats/lru"},{title:"LRU Miss",path:"/stats/lru/miss"},{title:"LRU Fallback",path:"/stats/lru/fallback"},{title:"LRU Global",path:"/stats/lru/global"},{title:"XDP Total",path:"/stats/xdp/total"},{title:"XDP Pass",path:"/stats/xdp/pass"},{title:"XDP Drop",path:"/stats/xdp/drop"},{title:"XDP Tx",path:"/stats/xdp/tx"}];function V(e){return Number.isFinite(e)?`${e>0?"+":""}${e}`:"0"}function Y(e){return Le[e]||{v1:"v1",v2:"v2"}}function De({title:e,path:t,diff:a=!1}){let{points:o,error:d}=K({path:t}),b=E(()=>Y(t),[t]),i=E(()=>[{label:b.v1,field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.2)"},{label:b.v2,field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[b]);return n`
      <div className="card">
        <h3>${e}</h3>
        ${d&&n`<p className="error">${d}</p>`}
        <${ee} title=${e} points=${o} keys=${i} diff=${a} inlineTitle=${!1} />
      </div>
    `}function qe({title:e,path:t}){let{points:a,error:o}=K({path:t}),d=E(()=>Y(t),[t]),b=a[a.length-1]||{},i=a[a.length-2]||{},s=Number(b.v1??0),p=Number(b.v2??0),c=s-Number(i.v1??0),m=p-Number(i.v2??0);return n`
      <div className="summary-card">
        <div className="summary-title">${e}</div>
        ${o?n`<p className="error">${o}</p>`:n`
              <div className="summary-row">
                <div className="stat">
                  <span className="muted">${d.v1}</span>
                  <strong>${s}</strong>
                </div>
                <div className="stat">
                  <span className="muted">${d.v1} delta/sec</span>
                  <strong className=${c<0?"delta down":"delta up"}>
                    ${V(c)}
                  </strong>
                </div>
              </div>
              <div className="summary-row">
                <div className="stat">
                  <span className="muted">${d.v2}</span>
                  <strong>${p}</strong>
                </div>
                <div className="stat">
                  <span className="muted">${d.v2} delta/sec</span>
                  <strong className=${m<0?"delta down":"delta up"}>
                    ${V(m)}
                  </strong>
                </div>
              </div>
            `}
      </div>
    `}function Oe(){let{data:e,error:t}=ke(()=>N.get("/stats/userspace"),1e3,[]);return n`
      <main>
        <section className="card">
          <h2>Global Stats</h2>
          <p className="muted">Polling every second. Charts show per-second deltas.</p>
        </section>
        <section className="grid">
          ${be.map(a=>n`<${De} title=${a.title} path=${a.path} diff=${!0} />`)}
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
            ${be.map(a=>n`<${qe} title=${a.title} path=${a.path} />`)}
          </div>
        </section>
      </main>
    `}function Be(){let[e,t]=h([]),[a,o]=h(""),[d,b]=h([]),[i,s]=h(""),[p,c]=h(null),[m,r]=h("");G(()=>{let f=!0;return(async()=>{try{let S=await N.get("/vips");if(!f)return;t(S||[]),!a&&S&&S.length>0&&o(J(S[0]))}catch(S){f&&r(S.message||"request failed")}})(),()=>{f=!1}},[]),G(()=>{if(!a)return;let f=Q(a),C=!0;return(async()=>{try{let k=await N.get("/vips/reals",f);if(!C)return;b(k||[]),k&&k.length>0?s(A=>A||k[0].address):s(""),r("")}catch(k){C&&r(k.message||"request failed")}})(),()=>{C=!1}},[a]),G(()=>{if(!i){c(null);return}let f=!0;return(async()=>{try{let S=await N.get("/reals/index",{address:i});if(!f)return;c(S?.index??null),r("")}catch(S){f&&r(S.message||"request failed")}})(),()=>{f=!1}},[i]);let{points:g,error:y}=K({path:"/stats/real",body:p!==null?{index:p}:null}),w=E(()=>Y("/stats/real"),[]),T=E(()=>[{label:w.v1,field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.2)"},{label:w.v2,field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[w]),_=g[g.length-1]||{},P=g[g.length-2]||{},R=Number(_.v1??0),x=Number(_.v2??0),l=R-Number(P.v1??0),v=x-Number(P.v2??0);return n`
      <main>
        <section className="card">
          <h2>Per-Real Stats</h2>
          <p className="muted">Select a VIP and real address to chart.</p>
          ${m&&n`<p className="error">${m}</p>`}
          <div className="form-row">
            <label className="field">
              <span>VIP</span>
              <select value=${a} onChange=${f=>o(f.target.value)}>
                ${e.map(f=>n`
                    <option value=${J(f)}>
                      ${f.address}:${f.port} / ${f.proto}
                    </option>
                  `)}
              </select>
            </label>
            <label className="field">
              <span>Real</span>
              <select
                value=${i}
                onChange=${f=>s(f.target.value)}
                disabled=${d.length===0}
              >
                ${d.map(f=>n`
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
          ${y&&n`<p className="error">${y}</p>`}
          ${p===null?n`<p className="muted">Select a real to start polling.</p>`:n`
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
                        <span className=${`delta ${l<0?"down":"up"}`}>
                          ${V(l)}
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
    `}function Ve(){let{addToast:e}=W(),[t,a]=h(""),[o,d]=h(""),[b,i]=h(!0),[s,p]=h(""),c=U(!0),m=async()=>{if(c.current){i(!0),d("");try{let g=await fetch(`${N.base}/config/export`,{headers:{Accept:"application/x-yaml"}});if(!g.ok){let w=`HTTP ${g.status}`;try{w=(await g.json())?.error?.message||w}catch{}throw new Error(w)}let y=await g.text();if(!c.current)return;a(y||""),p(new Date().toLocaleString())}catch(g){c.current&&d(g.message||"request failed")}finally{c.current&&i(!1)}}},r=async()=>{if(t)try{await navigator.clipboard.writeText(t),e("Config copied to clipboard","info")}catch{e("Failed to copy config","error")}};return G(()=>(c.current=!0,m(),()=>{c.current=!1}),[]),n`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>Running config</h2>
              <p className="muted">Exported from /api/v1/config/export</p>
            </div>
            <div className="row">
              <button className="btn ghost" onClick=${r} disabled=${!t}>
                Copy YAML
              </button>
              <button className="btn" onClick=${m} disabled=${b}>
                Refresh
              </button>
            </div>
          </div>
          ${o&&n`<p className="error">${o}</p>`}
          ${b?n`<p className="muted">Loading config...</p>`:t?n`<pre className="yaml-view">${t}</pre>`:n`<p className="muted">No config data returned.</p>`}
          ${s&&n`<p className="muted">Last fetched ${s}</p>`}
        </section>
      </main>
    `}function Ue(){let{addToast:e}=W(),{groups:t,setGroups:a,refreshFromStorage:o,importFromRunningConfig:d}=me(),[b,i]=h(""),[s,p]=h(""),[c,m]=h({address:"",weight:100,flags:0}),[r,g]=h(""),[y,w]=h(!1);G(()=>{if(s){if(!t[s]){let v=Object.keys(t);p(v[0]||"")}}else{let v=Object.keys(t);v.length>0&&p(v[0])}},[t,s]);let T=v=>{v.preventDefault();let f=b.trim();if(!f){g("Provide a group name.");return}if(t[f]){g("That group already exists.");return}a({...t,[f]:[]}),i(""),p(f),g(""),e(`Target group "${f}" created.`,"success")},_=v=>{let f={...t};delete f[v],a(f),e(`Target group "${v}" removed.`,"success")},P=v=>{if(v.preventDefault(),!s){g("Select a group to add a real.");return}let f=de(c);if(!f){g("Provide a valid real address.");return}let C=t[s]||[],S=C.some(k=>k.address===f.address)?C.map(k=>k.address===f.address?f:k):C.concat(f);a({...t,[s]:S}),m({address:"",weight:100,flags:0}),g(""),e("Real saved to target group.","success")},R=v=>{if(!s)return;let C=(t[s]||[]).filter(S=>S.address!==v);a({...t,[s]:C})},x=(v,f)=>{if(!s)return;let S=(t[s]||[]).map(k=>k.address===v?{...k,...f}:k);a({...t,[s]:S})};return n`
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
              <button className="btn ghost" type="button" onClick=${async()=>{w(!0);try{await d(),e("Imported target groups from running config.","success"),g("")}catch(v){g(v.message||"Failed to import target groups."),e(v.message||"Import failed.","error")}finally{w(!1)}}} disabled=${y}>
                ${y?"Importing...":"Import from running config"}
              </button>
            </div>
          </div>
          ${r&&n`<p className="error">${r}</p>`}
          <form className="form" onSubmit=${T}>
            <div className="form-row">
              <label className="field">
                <span>New group name</span>
                <input
                  value=${b}
                  onInput=${v=>i(v.target.value)}
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
                  value=${s}
                  onChange=${v=>p(v.target.value)}
                  disabled=${Object.keys(t).length===0}
                >
                  ${Object.keys(t).map(v=>n`<option value=${v}>${v}</option>`)}
                </select>
              </label>
              ${s&&n`<button className="btn danger" type="button" onClick=${()=>_(s)}>
                Delete group
              </button>`}
            </div>
          </div>
          ${s?n`
                <table className="table">
                  <thead>
                    <tr>
                      <th>Address</th>
                      <th>Weight</th>
                      <th></th>
                    </tr>
                  </thead>
                  <tbody>
                    ${(t[s]||[]).map(v=>n`
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
                        value=${c.address}
                        onInput=${v=>m({...c,address:v.target.value})}
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
                        onInput=${v=>m({...c,weight:v.target.value})}
                        required
                      />
                    </label>
                  </div>
                  <button className="btn secondary" type="submit">Save real</button>
                </form>
              `:n`<p className="muted">No groups yet. Create one to add reals.</p>`}
        </section>
      </main>
    `}function ze(){let[e,t]=h({initialized:!1,ready:!1}),[a,o]=h([]),d=U({}),b=(s,p="info")=>{let c=`${Date.now()}-${Math.random().toString(16).slice(2)}`;o(m=>m.concat({id:c,message:s,kind:p})),d.current[c]=setTimeout(()=>{o(m=>m.filter(r=>r.id!==c)),delete d.current[c]},4e3)},i=s=>{d.current[s]&&(clearTimeout(d.current[s]),delete d.current[s]),o(p=>p.filter(c=>c.id!==s))};return G(()=>{let s=!0,p=async()=>{try{let m=await N.get("/lb/status");s&&t(m||{initialized:!1,ready:!1})}catch{s&&t({initialized:!1,ready:!1})}};p();let c=setInterval(p,5e3);return()=>{s=!1,clearInterval(c)}},[]),n`
      <${re}>
        <${fe}>
          <${Z.Provider} value=${{addToast:b}}>
            <${Fe} status=${e} />
            <${ne}>
              <${O} path="/" element=${n`<${Ge} />`} />
              <${O} path="/vips/:vipId" element=${n`<${Ae} />`} />
              <${O} path="/vips/:vipId/stats" element=${n`<${Ee} />`} />
              <${O} path="/target-groups" element=${n`<${Ue} />`} />
              <${O} path="/stats/global" element=${n`<${Oe} />`} />
              <${O} path="/stats/real" element=${n`<${Be} />`} />
              <${O} path="/config" element=${n`<${Ve} />`} />
            </${ne}>
            <${xe} toasts=${a} onDismiss=${i} />
          </${Z.Provider}>
        </${fe}>
      </${re}>
    `}ReactDOM.createRoot(document.getElementById("root")).render(n`<${ze} />`)})();})();
