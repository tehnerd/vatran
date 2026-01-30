(()=>{(()=>{let{useEffect:T,useMemo:V,useRef:U,useState:h,useContext:_e}=React,{BrowserRouter:se,Routes:re,Route:q,NavLink:F,Link:j,useParams:oe,useNavigate:Ie}=ReactRouterDOM,r=htm.bind(React.createElement),X=React.createContext({addToast:()=>{}}),M=[{label:"NO_SPORT",value:1},{label:"NO_LRU",value:2},{label:"QUIC_VIP",value:4},{label:"DPORT_HASH",value:8},{label:"SRC_ROUTING",value:16},{label:"LOCAL_VIP",value:32},{label:"GLOBAL_LRU",value:64},{label:"HASH_SRC_DST_PORT",value:128},{label:"UDP_STABLE_ROUTING_VIP",value:256},{label:"UDP_FLOW_MIGRATION",value:512}];function H(){return _e(X)}function Se(e){return String(e).replace(/[^a-z0-9_-]/gi,"_")}function Re(e,t,a){let o=Number(e)||0,i=Number(t)||0;return a?o|i:o&~i}function Ce(e,t){let a=Number(e)||0;return t.filter(o=>(a&o.value)!==0)}function ne({mask:e,options:t,showStatus:a=!1,emptyLabel:o="None"}){let i=Number(e)||0,g=a?t:Ce(i,t),n=a?2:1;return r`
      <table className="table flag-table">
        <thead>
          <tr>
            <th>Flag</th>
            ${a?r`<th>Enabled</th>`:null}
          </tr>
        </thead>
        <tbody>
          ${g.length===0?r`<tr><td colspan=${n} className="muted">${o}</td></tr>`:g.map(s=>{let p=(i&s.value)!==0;return r`
                  <tr>
                    <td>${s.label}</td>
                    ${a?r`<td>${p?"Yes":"No"}</td>`:null}
                  </tr>
                `})}
        </tbody>
      </table>
    `}function le({options:e,value:t,onChange:a,name:o}){let i=Number(t)||0,g=Se(o||"flags");return r`
      <div className="flag-selector">
        <table className="table flag-table">
          <thead>
            <tr>
              <th>Enabled</th>
              <th>Flag</th>
            </tr>
          </thead>
          <tbody>
            ${e.map(n=>{let s=`${g}-${n.value}`,p=(i&n.value)===n.value;return r`
                <tr>
                  <td>
                    <input
                      id=${s}
                      type="checkbox"
                      checked=${p}
                      onChange=${d=>a(Re(i,n.value,d.target.checked))}
                    />
                  </td>
                  <td>
                    <label className="flag-option" htmlFor=${s}>${n.label}</label>
                  </td>
                </tr>
              `})}
          </tbody>
        </table>
      </div>
    `}let N={base:"/api/v1",async request(e,t={}){let a={method:t.method||"GET",headers:{"Content-Type":"application/json"}},o=`${N.base}${e}`;if(t.body!==void 0&&t.body!==null)if(a.method==="GET"){let n=new URLSearchParams;Object.entries(t.body).forEach(([p,d])=>{if(d!=null){if(Array.isArray(d)){d.forEach(b=>n.append(p,String(b)));return}if(typeof d=="object"){n.set(p,JSON.stringify(d));return}n.set(p,String(d))}});let s=n.toString();s&&(o+=`${o.includes("?")?"&":"?"}${s}`)}else a.body=JSON.stringify(t.body);let i=await fetch(o,a),g;try{g=await i.json()}catch{throw new Error("invalid JSON response")}if(!i.ok)throw new Error(g?.error?.message||`HTTP ${i.status}`);if(!g.success){let n=g.error?.message||"request failed";throw new Error(n)}return g.data},get(e,t){return N.request(e,{method:"GET",body:t})},post(e,t){return N.request(e,{method:"POST",body:t})},put(e,t){return N.request(e,{method:"PUT",body:t})},del(e,t){return N.request(e,{method:"DELETE",body:t})}},ie="vatran_target_groups";function ce(e){if(!e||!e.address)return null;let t=String(e.address).trim();if(!t)return null;let a=Number(e.weight),o=Number(e.flags??0);return{address:t,weight:Number.isFinite(a)?a:0,flags:Number.isFinite(o)?o:0}}function de(e){if(!e||typeof e!="object")return{};let t={};return Object.entries(e).forEach(([a,o])=>{let i=String(a).trim();if(!i)return;let g=Array.isArray(o)?o.map(ce).filter(Boolean):[],n=[],s=new Set;g.forEach(p=>{s.has(p.address)||(s.add(p.address),n.push(p))}),t[i]=n}),t}function Y(){if(typeof localStorage>"u")return{};try{let e=localStorage.getItem(ie);return e?de(JSON.parse(e)):{}}catch{return{}}}function ue(e){if(!(typeof localStorage>"u"))try{localStorage.setItem(ie,JSON.stringify(e))}catch{}}function Te(e,t){let a={...e};return Object.entries(t||{}).forEach(([o,i])=>{a[o]||(a[o]=i)}),a}function pe(){let[e,t]=h(()=>Y());return T(()=>{ue(e)},[e]),{groups:e,setGroups:t,refreshFromStorage:()=>{t(Y())},importFromRunningConfig:async()=>{let i=await N.get("/config/export/json"),g=de(i?.target_groups||{}),n=Te(Y(),g);return t(n),ue(n),n}}}function W(e){return`${encodeURIComponent(e.address)}:${e.port}:${e.proto}`}function Z(e){let t=e.split(":"),a=Number(t.pop()||0),o=Number(t.pop()||0);return{address:decodeURIComponent(t.join(":")),port:o,proto:a}}function Pe(e,t,a=[]){let[o,i]=h(null),[g,n]=h(""),[s,p]=h(!0);return T(()=>{let d=!0,b=async()=>{try{let m=await e();d&&(i(m),n(""),p(!1))}catch(m){d&&(n(m.message||"request failed"),p(!1))}};b();let u=setInterval(b,t);return()=>{d=!1,clearInterval(u)}},a),{data:o,error:g,loading:s}}function J({path:e,body:t,intervalMs:a=1e3,limit:o=60}){let[i,g]=h([]),[n,s]=h(""),p=V(()=>JSON.stringify(t||{}),[t]);return T(()=>{if(t===null)return g([]),s(""),()=>{};let d=!0,b=async()=>{try{let m=await N.get(e,t);if(!d)return;let y=new Date().toLocaleTimeString();g(I=>I.concat({label:y,...m}).slice(-o)),s("")}catch(m){d&&s(m.message||"request failed")}};b();let u=setInterval(b,a);return()=>{d=!1,clearInterval(u)}},[e,p,a,o]),{points:i,error:n}}function me({title:e,points:t,keys:a,diff:o=!1,height:i=120,showTitle:g=!1}){let n=U(null),s=U(null);return T(()=>{if(!n.current)return;s.current||(s.current=new Chart(n.current,{type:"line",data:{labels:[],datasets:[]},options:{responsive:!0,maintainAspectRatio:!1,animation:!1,scales:{x:{grid:{display:!1}},y:{beginAtZero:!o}},plugins:{legend:{display:!0,position:"bottom"},title:{display:g&&!!e,text:e}}}}));let p=s.current,d=t.map(b=>b.label);return p.data.labels=d,p.data.datasets=a.map(b=>{let u=t.map(y=>y[b.field]||0),m=o?u.map((y,I)=>I===0?0:y-u[I-1]):u;return{label:b.label,data:m,borderColor:b.color,backgroundColor:b.fill,borderWidth:2,tension:.3}}),p.options.scales.y.beginAtZero=!o,p.options.plugins.title.display=g&&!!e,p.options.plugins.title.text=e||"",p.update(),()=>{}},[t,a,e,o,g]),T(()=>()=>{s.current&&(s.current.destroy(),s.current=null)},[]),r`<canvas ref=${n} height=${i}></canvas>`}function Q({title:e,points:t,keys:a,diff:o=!1,inlineTitle:i=!0}){let[g,n]=h(!1);return r`
      <div className="chart-wrap">
        <button className="chart-zoom-button" type="button" onClick=${()=>n(!0)}>
          <span className="zoom-icon" aria-hidden="true">+</span>
          Zoom
        </button>
        <div className="chart-click" onClick=${()=>n(!0)}>
          <${me}
            title=${e}
            points=${t}
            keys=${a}
            diff=${o}
            height=${120}
            showTitle=${i&&!!e}
          />
        </div>
        ${g&&r`
          <div className="chart-overlay" onClick=${()=>n(!1)}>
            <div className="chart-modal" onClick=${s=>s.stopPropagation()}>
              <div className="row chart-modal-header">
                <div>
                  <h3>${e||"Chart"}</h3>
                  ${o?r`<p className="muted">Per-second delta.</p>`:""}
                </div>
                <button className="btn ghost" onClick=${()=>n(!1)}>
                  Close
                </button>
              </div>
              <div className="chart-zoom">
                <${me}
                  title=${e}
                  points=${t}
                  keys=${a}
                  diff=${o}
                  height=${360}
                  showTitle=${!1}
                />
              </div>
            </div>
          </div>
        `}
      </div>
    `}function ge({children:e}){return e}function ke({toasts:e,onDismiss:t}){return r`
      <div className="toast-stack">
        ${e.map(a=>r`
            <div className=${`toast ${a.kind}`}>
              <span>${a.message}</span>
              <button className="toast-close" onClick=${()=>t(a.id)}>
                ×
              </button>
            </div>
          `)}
      </div>
    `}function Fe({status:e}){return r`
      <header>
        <div>
          <div style=${{fontSize:20,fontWeight:700}}>Vatran</div>
          <div className="status-pill">
            <span className=${`dot ${e.ready?"ok":""}`}></span>
            ${e.ready?"Ready":"Not ready"}
          </div>
        </div>
        <nav>
          <${F} to="/" end className=${({isActive:t})=>t?"active":""}>
            Dashboard
          </${F}>
          <${F} to="/stats/global" className=${({isActive:t})=>t?"active":""}>
            Global stats
          </${F}>
          <${F} to="/stats/real" className=${({isActive:t})=>t?"active":""}>
            Per-real stats
          </${F}>
          <${F}
            to="/target-groups"
            className=${({isActive:t})=>t?"active":""}
          >
            Target groups
          </${F}>
          <${F} to="/config" className=${({isActive:t})=>t?"active":""}>
            Config export
          </${F}>
        </nav>
      </header>
    `}function xe(){let{addToast:e}=H(),[t,a]=h({initialized:!1,ready:!1}),[o,i]=h([]),[g,n]=h(""),[s,p]=h(!1),[d,b]=h(!1),[u,m]=h({main_interface:"",balancer_prog_path:"",healthchecking_prog_path:"",default_mac:"",local_mac:"",root_map_path:"",root_map_pos:2,katran_src_v4:"",katran_src_v6:"",use_root_map:!1,max_vips:512,max_reals:4096,hash_func:0}),[y,I]=h({address:"",port:80,proto:6,flags:0}),v=async()=>{try{let l=await N.get("/lb/status"),f=await N.get("/vips"),w=await Promise.all((f||[]).map(async P=>{try{let x=await N.get("/vips/flags",{address:P.address,port:P.port,proto:P.proto});return{...P,flags:x?.flags??0}}catch{return{...P,flags:null}}}));a(l||{initialized:!1,ready:!1}),i(w),n("")}catch(l){n(l.message||"request failed")}};T(()=>{let l=!0;return(async()=>{l&&await v()})(),()=>{l=!1}},[]);let k=async l=>{l.preventDefault();try{let f={...u,root_map_pos:u.root_map_pos===""?void 0:Number(u.root_map_pos),max_vips:Number(u.max_vips),max_reals:Number(u.max_reals),hash_func:Number(u.hash_func)};await N.post("/lb/create",f),n(""),p(!1),e("Load balancer initialized.","success"),await v()}catch(f){n(f.message||"request failed"),e(f.message||"Initialize failed.","error")}},S=async l=>{l.preventDefault();try{await N.post("/vips",{...y,port:Number(y.port),proto:Number(y.proto),flags:Number(y.flags||0)}),I({address:"",port:80,proto:6,flags:0}),n(""),b(!1),e("VIP created.","success"),await v()}catch(f){n(f.message||"request failed"),e(f.message||"VIP create failed.","error")}},R=async()=>{try{await N.post("/lb/load-bpf-progs"),n(""),e("BPF programs loaded.","success"),await v()}catch(l){n(l.message||"request failed"),e(l.message||"Load BPF programs failed.","error")}},D=async()=>{try{await N.post("/lb/attach-bpf-progs"),n(""),e("BPF programs attached.","success"),await v()}catch(l){n(l.message||"request failed"),e(l.message||"Attach BPF programs failed.","error")}};return r`
      <main>
        <section className="card">
          <h2>Load balancer</h2>
          <p>Initialized: ${t.initialized?"yes":"no"}</p>
          <p>Ready: ${t.ready?"yes":"no"}</p>
          <div className="row">
            ${!t.initialized&&r`
              <button className="btn" onClick=${()=>p(l=>!l)}>
                ${s?"Close":"Initialize"}
              </button>
            `}
            <button className="btn secondary" onClick=${()=>b(l=>!l)}>
              ${d?"Close":"Create VIP"}
            </button>
          </div>
          ${!t.ready&&r`
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
                onClick=${D}
              >
                Attach BPF Programs
              </button>
            </div>
          `}
          ${s&&r`
            <form className="form" onSubmit=${k}>
              <div className="form-row">
                <label className="field">
                  <span>Main interface</span>
                  <input
                    value=${u.main_interface}
                    onInput=${l=>m({...u,main_interface:l.target.value})}
                    placeholder="eth0"
                    required
                  />
                </label>
                <label className="field">
                  <span>Balancer prog path</span>
                  <input
                    value=${u.balancer_prog_path}
                    onInput=${l=>m({...u,balancer_prog_path:l.target.value})}
                    placeholder="/path/to/bpf.o"
                    required
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Healthchecking path</span>
                  <input
                    value=${u.healthchecking_prog_path}
                    onInput=${l=>m({...u,healthchecking_prog_path:l.target.value})}
                    placeholder="/path/to/hc.o"
                  />
                </label>
                <label className="field">
                  <span>Default MAC</span>
                  <input
                    value=${u.default_mac}
                    onInput=${l=>m({...u,default_mac:l.target.value})}
                    placeholder="aa:bb:cc:dd:ee:ff"
                    required
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Local MAC</span>
                  <input
                    value=${u.local_mac}
                    onInput=${l=>m({...u,local_mac:l.target.value})}
                    placeholder="11:22:33:44:55:66"
                    required=${u.healthchecking_prog_path?.trim()!==""}
                  />
                </label>
                <label className="field">
                  <span>Hash function</span>
                  <input
                    type="number"
                    min="0"
                    value=${u.hash_func}
                    onInput=${l=>m({...u,hash_func:l.target.value})}
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Root map path</span>
                  <input
                    value=${u.root_map_path}
                    onInput=${l=>m({...u,root_map_path:l.target.value})}
                    placeholder="/sys/fs/bpf/root_map"
                  />
                </label>
                <label className="field">
                  <span>Root map position</span>
                  <input
                    type="number"
                    min="0"
                    value=${u.root_map_pos}
                    onInput=${l=>m({...u,root_map_pos:l.target.value})}
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Katran src v4</span>
                  <input
                    value=${u.katran_src_v4}
                    onInput=${l=>m({...u,katran_src_v4:l.target.value})}
                    placeholder="10.0.0.1"
                  />
                </label>
                <label className="field">
                  <span>Katran src v6</span>
                  <input
                    value=${u.katran_src_v6}
                    onInput=${l=>m({...u,katran_src_v6:l.target.value})}
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
                    value=${u.max_vips}
                    onInput=${l=>m({...u,max_vips:l.target.value})}
                  />
                </label>
                <label className="field">
                  <span>Max reals</span>
                  <input
                    type="number"
                    min="1"
                    value=${u.max_reals}
                    onInput=${l=>m({...u,max_reals:l.target.value})}
                  />
                </label>
              </div>
              <label className="field checkbox">
                <input
                  type="checkbox"
                  checked=${u.use_root_map}
                  onChange=${l=>m({...u,use_root_map:l.target.checked})}
                />
                <span>Use root map</span>
              </label>
              <button className="btn" type="submit">Initialize LB</button>
            </form>
          `}
          ${d&&r`
            <form className="form" onSubmit=${S}>
              <div className="form-row">
                <label className="field">
                  <span>VIP address</span>
                  <input
                    value=${y.address}
                    onInput=${l=>I({...y,address:l.target.value})}
                    placeholder="1.2.3.4"
                    required
                  />
                </label>
                <label className="field">
                  <span>Port</span>
                  <input
                    type="number"
                    value=${y.port}
                    onInput=${l=>I({...y,port:l.target.value})}
                    required
                  />
                </label>
                <label className="field">
                  <span>Protocol</span>
                  <select
                    value=${y.proto}
                    onChange=${l=>I({...y,proto:l.target.value})}
                  >
                    <option value="6">TCP (6)</option>
                    <option value="17">UDP (17)</option>
                  </select>
                </label>
                <label className="field">
                  <span>Flags</span>
                  <${le}
                    options=${M}
                    value=${y.flags}
                    name="vip-add"
                    onChange=${l=>I({...y,flags:l})}
                  />
                </label>
              </div>
              <button className="btn" type="submit">Create VIP</button>
            </form>
          `}
          ${g&&r`<p className="error">${g}</p>`}
        </section>
        <section className="card">
          <div className="section-header">
            <h2>VIPs</h2>
            <button className="btn ghost" onClick=${v}>Refresh</button>
          </div>
          ${o.length===0?r`<p className="muted">No VIPs configured yet.</p>`:r`
                <div className="grid">
                  ${o.map(l=>r`
                      <div className="card">
                        <div style=${{fontWeight:600}}>
                          ${l.address}:${l.port} / ${l.proto}
                        </div>
                        <div className="muted" style=${{marginTop:6}}>
                          <div style=${{fontWeight:600,marginBottom:6}}>Flags</div>
                          <${ne}
                            mask=${l.flags}
                            options=${M}
                            emptyLabel=${l.flags===null?"Unknown":"No flags"}
                          />
                        </div>
                        <div className="row" style=${{marginTop:12}}>
                          <${j} className="btn" to=${`/vips/${W(l)}`}>
                            Open
                          </${j}>
                          <${j}
                            className="btn secondary"
                            to=${`/vips/${W(l)}/stats`}
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
    `}function Ge(){let{addToast:e}=H(),t=oe(),a=Ie(),o=V(()=>Z(t.vipId),[t.vipId]),[i,g]=h([]),[n,s]=h(""),[p,d]=h(""),[b,u]=h(!0),[m,y]=h({address:"",weight:100,flags:0}),[I,v]=h({}),[k,S]=h(null),[R,D]=h({flags:0,set:!0}),[l,f]=h({hash_function:0}),{groups:w,setGroups:P,refreshFromStorage:x,importFromRunningConfig:G}=pe(),[A,he]=h(""),[be,O]=h(""),[ve,$e]=h(!1),[Ne,ye]=h(""),[ee,we]=h({add:0,update:0,remove:0}),B=async()=>{try{let c=await N.get("/vips/reals",o);g(c||[]);let $={};(c||[]).forEach(L=>{$[L.address]=L.weight}),v($),s(""),u(!1)}catch(c){s(c.message||"request failed"),u(!1)}},te=async()=>{try{let c=await N.get("/vips/flags",o);S(c?.flags??0),d("")}catch(c){d(c.message||"request failed")}};T(()=>{B(),te()},[t.vipId]),T(()=>{if(!A){we({add:0,update:0,remove:0});return}let c=w[A]||[],$=new Map(i.map(C=>[C.address,C])),L=new Map(c.map(C=>[C.address,C])),E=0,z=0,_=0;c.forEach(C=>{let ae=$.get(C.address);if(!ae){E+=1;return}(Number(ae.weight)!==Number(C.weight)||Number(ae.flags||0)!==Number(C.flags||0))&&(z+=1)}),i.forEach(C=>{L.has(C.address)||(_+=1)}),we({add:E,update:z,remove:_})},[A,i,w]);let ze=async c=>{try{let $=Number(I[c.address]);await N.post("/vips/reals",{vip:o,real:{address:c.address,weight:$,flags:c.flags||0}}),await B(),e("Real weight updated.","success")}catch($){s($.message||"request failed"),e($.message||"Update failed.","error")}},Ue=async c=>{try{await N.del("/vips/reals",{vip:o,real:{address:c.address,weight:c.weight,flags:c.flags||0}}),await B(),e("Real removed.","success")}catch($){s($.message||"request failed"),e($.message||"Remove failed.","error")}},je=async c=>{c.preventDefault();try{await N.post("/vips/reals",{vip:o,real:{address:m.address,weight:Number(m.weight),flags:Number(m.flags||0)}}),y({address:"",weight:100,flags:0}),await B(),e("Real added.","success")}catch($){s($.message||"request failed"),e($.message||"Add failed.","error")}},Me=async()=>{if(!A||!w[A]){O("Select a target group to apply.");return}$e(!0),O("");let c=w[A]||[],$=new Map(i.map(_=>[_.address,_])),L=new Map(c.map(_=>[_.address,_])),E=i.filter(_=>!L.has(_.address)),z=c.filter(_=>{let C=$.get(_.address);return C?Number(C.weight)!==Number(_.weight)||Number(C.flags||0)!==Number(_.flags||0):!0});try{E.length>0&&await N.put("/vips/reals/batch",{vip:o,action:1,reals:E.map(_=>({address:_.address,weight:Number(_.weight),flags:Number(_.flags||0)}))}),z.length>0&&await Promise.all(z.map(_=>N.post("/vips/reals",{vip:o,real:{address:_.address,weight:Number(_.weight),flags:Number(_.flags||0)}}))),await B(),e(`Applied target group "${A}".`,"success")}catch(_){O(_.message||"Failed to apply target group."),e(_.message||"Target group apply failed.","error")}finally{$e(!1)}},He=c=>{c.preventDefault();let $=Ne.trim();if(!$){O("Provide a name for the new target group.");return}if(w[$]){O("A target group with that name already exists.");return}let L={...w,[$]:i.map(E=>({address:E.address,weight:Number(E.weight),flags:Number(E.flags||0)}))};P(L),ye(""),he($),O(""),e(`Target group "${$}" saved.`,"success")},We=async()=>{try{await N.del("/vips",o),e("VIP deleted.","success"),a("/")}catch(c){s(c.message||"request failed"),e(c.message||"Delete failed.","error")}},Je=async c=>{c.preventDefault();try{await N.put("/vips/flags",{...o,flag:Number(R.flags||0),set:!!R.set}),await te(),e("VIP flags updated.","success")}catch($){d($.message||"request failed"),e($.message||"Flag update failed.","error")}},Ke=async c=>{c.preventDefault();try{await N.put("/vips/hash-function",{...o,hash_function:Number(l.hash_function)}),e("Hash function updated.","success")}catch($){d($.message||"request failed"),e($.message||"Hash update failed.","error")}};return r`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>VIP Detail</h2>
              <p className="muted">${o.address}:${o.port} / ${o.proto}</p>
              ${k===null?r`<p className="muted">Flags: —</p>`:r`
                    <div style=${{marginTop:8}}>
                      <${ne}
                        mask=${k}
                        options=${M}
                        showStatus=${!0}
                        emptyLabel="No flags"
                      />
                    </div>
                  `}
            </div>
            <div className="row">
              <button className="btn ghost" onClick=${te}>Refresh flags</button>
              <button className="btn danger" onClick=${We}>Delete VIP</button>
            </div>
          </div>
          ${n&&r`<p className="error">${n}</p>`}
          ${p&&r`<p className="error">${p}</p>`}
          ${b?r`<p className="muted">Loading reals…</p>`:r`
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
                    ${i.map(c=>r`
                        <tr>
                          <td>
                            <span
                              className=${`dot ${Number(c.weight)>0?"ok":"bad"}`}
                            ></span>
                          </td>
                          <td>${c.address}</td>
                          <td>
                            <input
                              className="inline-input"
                              type="number"
                              min="0"
                              value=${I[c.address]??c.weight}
                              onInput=${$=>v({...I,[c.address]:$.target.value})}
                            />
                          </td>
                          <td className="row">
                            <button className="btn" onClick=${()=>ze(c)}>
                              Update
                            </button>
                            <button className="btn ghost" onClick=${()=>Ue(c)}>
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
            <form className="form" onSubmit=${Je}>
              <div className="form-row">
                <label className="field">
                  <span>Flags</span>
                  <${le}
                    options=${M}
                    value=${R.flags}
                    name="vip-flag-change"
                    onChange=${c=>D({...R,flags:c})}
                  />
                </label>
                <label className="field">
                  <span>Set</span>
                  <select
                    value=${String(R.set)}
                    onChange=${c=>D({...R,set:c.target.value==="true"})}
                  >
                    <option value="true">Set</option>
                    <option value="false">Clear</option>
                  </select>
                </label>
              </div>
              <button className="btn" type="submit">Apply flag</button>
            </form>
            <form className="form" onSubmit=${Ke}>
              <div className="form-row">
                <label className="field">
                  <span>Hash function</span>
                  <input
                    type="number"
                    min="0"
                    value=${l.hash_function}
                    onInput=${c=>f({...l,hash_function:c.target.value})}
                  />
                </label>
              </div>
              <button className="btn" type="submit">Apply hash</button>
            </form>
          </div>
        </section>
        <section className="card">
          <h3>Add real</h3>
          <form className="form" onSubmit=${je}>
            <div className="form-row">
              <label className="field">
                <span>Real address</span>
                <input
                  value=${m.address}
                  onInput=${c=>y({...m,address:c.target.value})}
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
                  onInput=${c=>y({...m,weight:c.target.value})}
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
              <button className="btn ghost" type="button" onClick=${x}>
                Reload groups
              </button>
              <button
                className="btn ghost"
                type="button"
                onClick=${async()=>{try{await G(),e("Imported target groups from running config.","success")}catch(c){O(c.message||"Failed to import target groups."),e(c.message||"Import failed.","error")}}}
              >
                Import from running config
              </button>
            </div>
          </div>
          ${be&&r`<p className="error">${be}</p>`}
          <div className="form-row">
            <label className="field">
              <span>Target group</span>
              <select
                value=${A}
                onChange=${c=>he(c.target.value)}
                disabled=${Object.keys(w).length===0}
              >
                <option value="">Select group</option>
                ${Object.keys(w).map(c=>r`<option value=${c}>${c}</option>`)}
              </select>
            </label>
            <label className="field">
              <span>Preview</span>
              <input
                value=${`add ${ee.add} \xB7 update ${ee.update} \xB7 remove ${ee.remove}`}
                readOnly
              />
            </label>
          </div>
          <div className="row">
            <button
              className="btn"
              type="button"
              onClick=${Me}
              disabled=${ve||!A}
            >
              ${ve?"Applying...":"Apply target group"}
            </button>
          </div>
          <form className="form" onSubmit=${He}>
            <div className="form-row">
              <label className="field">
                <span>Save current reals as new group</span>
                <input
                  value=${Ne}
                  onInput=${c=>ye(c.target.value)}
                  placeholder="edge-backends"
                />
              </label>
            </div>
            <button className="btn secondary" type="submit">Save target group</button>
          </form>
        </section>
      </main>
    `}function Ae(){let e=oe(),t=V(()=>Z(e.vipId),[e.vipId]),{points:a,error:o}=J({path:"/stats/vip",body:t}),i=a[a.length-1]||{},g=a[a.length-2]||{},n=Number(i.v1??0),s=Number(i.v2??0),p=n-Number(g.v1??0),d=s-Number(g.v2??0),b=V(()=>[{label:"v1",field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.15)"},{label:"v2",field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[]);return r`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>VIP Stats</h2>
              <p className="muted">${t.address}:${t.port} / ${t.proto}</p>
            </div>
          </div>
          ${o&&r`<p className="error">${o}</p>`}
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
                <td>v1</td>
                <td>${n}</td>
                <td>
                  <span className=${`delta ${p<0?"down":"up"}`}>
                    ${K(p)}
                  </span>
                </td>
              </tr>
              <tr>
                <td>v2</td>
                <td>${s}</td>
                <td>
                  <span className=${`delta ${d<0?"down":"up"}`}>
                    ${K(d)}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
          <${Q} title="Traffic (delta/sec)" points=${a} keys=${b} diff=${!0} />
        </section>
      </main>
    `}let fe=[{title:"LRU",path:"/stats/lru"},{title:"LRU Miss",path:"/stats/lru/miss"},{title:"LRU Fallback",path:"/stats/lru/fallback"},{title:"LRU Global",path:"/stats/lru/global"},{title:"XDP Total",path:"/stats/xdp/total"},{title:"XDP Pass",path:"/stats/xdp/pass"},{title:"XDP Drop",path:"/stats/xdp/drop"},{title:"XDP Tx",path:"/stats/xdp/tx"}];function K(e){return Number.isFinite(e)?`${e>0?"+":""}${e}`:"0"}function Ee({title:e,path:t,diff:a=!1}){let{points:o,error:i}=J({path:t}),g=V(()=>[{label:"v1",field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.2)"},{label:"v2",field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[]);return r`
      <div className="card">
        <h3>${e}</h3>
        ${i&&r`<p className="error">${i}</p>`}
        <${Q} title=${e} points=${o} keys=${g} diff=${a} inlineTitle=${!1} />
      </div>
    `}function De({title:e,path:t}){let{points:a,error:o}=J({path:t}),i=a[a.length-1]||{},g=a[a.length-2]||{},n=Number(i.v1??0),s=Number(i.v2??0),p=n-Number(g.v1??0),d=s-Number(g.v2??0);return r`
      <div className="summary-card">
        <div className="summary-title">${e}</div>
        ${o?r`<p className="error">${o}</p>`:r`
              <div className="summary-row">
                <div className="stat">
                  <span className="muted">v1 absolute</span>
                  <strong>${n}</strong>
                </div>
                <div className="stat">
                  <span className="muted">v1 delta/sec</span>
                  <strong className=${p<0?"delta down":"delta up"}>
                    ${K(p)}
                  </strong>
                </div>
              </div>
              <div className="summary-row">
                <div className="stat">
                  <span className="muted">v2 absolute</span>
                  <strong>${s}</strong>
                </div>
                <div className="stat">
                  <span className="muted">v2 delta/sec</span>
                  <strong className=${d<0?"delta down":"delta up"}>
                    ${K(d)}
                  </strong>
                </div>
              </div>
            `}
      </div>
    `}function Le(){let{data:e,error:t}=Pe(()=>N.get("/stats/userspace"),1e3,[]);return r`
      <main>
        <section className="card">
          <h2>Global Stats</h2>
          <p className="muted">Polling every second. Charts show per-second deltas.</p>
        </section>
        <section className="grid">
          ${fe.map(a=>r`<${Ee} title=${a.title} path=${a.path} diff=${!0} />`)}
        </section>
        <section className="card">
          <h3>Userspace</h3>
          ${t&&r`<p className="error">${t}</p>`}
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
        <section className="card">
          <div className="section-header">
            <div>
              <h3>Absolute & Rate of Change</h3>
              <p className="muted">Latest value and per-second delta.</p>
            </div>
          </div>
          <div className="summary-grid">
            ${fe.map(a=>r`<${De} title=${a.title} path=${a.path} />`)}
          </div>
        </section>
      </main>
    `}function qe(){let[e,t]=h([]),[a,o]=h(""),[i,g]=h([]),[n,s]=h(""),[p,d]=h(null),[b,u]=h("");T(()=>{let v=!0;return(async()=>{try{let S=await N.get("/vips");if(!v)return;t(S||[]),!a&&S&&S.length>0&&o(W(S[0]))}catch(S){v&&u(S.message||"request failed")}})(),()=>{v=!1}},[]),T(()=>{if(!a)return;let v=Z(a),k=!0;return(async()=>{try{let R=await N.get("/vips/reals",v);if(!k)return;g(R||[]),R&&R.length>0?s(D=>D||R[0].address):s(""),u("")}catch(R){k&&u(R.message||"request failed")}})(),()=>{k=!1}},[a]),T(()=>{if(!n){d(null);return}let v=!0;return(async()=>{try{let S=await N.get("/reals/index",{address:n});if(!v)return;d(S?.index??null),u("")}catch(S){v&&u(S.message||"request failed")}})(),()=>{v=!1}},[n]);let{points:m,error:y}=J({path:"/stats/real",body:p!==null?{index:p}:null}),I=V(()=>[{label:"v1",field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.2)"},{label:"v2",field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[]);return r`
      <main>
        <section className="card">
          <h2>Per-Real Stats</h2>
          <p className="muted">Select a VIP and real address to chart.</p>
          ${b&&r`<p className="error">${b}</p>`}
          <div className="form-row">
            <label className="field">
              <span>VIP</span>
              <select value=${a} onChange=${v=>o(v.target.value)}>
                ${e.map(v=>r`
                    <option value=${W(v)}>
                      ${v.address}:${v.port} / ${v.proto}
                    </option>
                  `)}
              </select>
            </label>
            <label className="field">
              <span>Real</span>
              <select
                value=${n}
                onChange=${v=>s(v.target.value)}
                disabled=${i.length===0}
              >
                ${i.map(v=>r`
                    <option value=${v.address}>${v.address}</option>
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
          ${y&&r`<p className="error">${y}</p>`}
          ${p===null?r`<p className="muted">Select a real to start polling.</p>`:r`<${Q} points=${m} keys=${I} />`}
        </section>
      </main>
    `}function Oe(){let{addToast:e}=H(),[t,a]=h(""),[o,i]=h(""),[g,n]=h(!0),[s,p]=h(""),d=U(!0),b=async()=>{if(d.current){n(!0),i("");try{let m=await fetch(`${N.base}/config/export`,{headers:{Accept:"application/x-yaml"}});if(!m.ok){let I=`HTTP ${m.status}`;try{I=(await m.json())?.error?.message||I}catch{}throw new Error(I)}let y=await m.text();if(!d.current)return;a(y||""),p(new Date().toLocaleString())}catch(m){d.current&&i(m.message||"request failed")}finally{d.current&&n(!1)}}},u=async()=>{if(t)try{await navigator.clipboard.writeText(t),e("Config copied to clipboard","info")}catch{e("Failed to copy config","error")}};return T(()=>(d.current=!0,b(),()=>{d.current=!1}),[]),r`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>Running config</h2>
              <p className="muted">Exported from /api/v1/config/export</p>
            </div>
            <div className="row">
              <button className="btn ghost" onClick=${u} disabled=${!t}>
                Copy YAML
              </button>
              <button className="btn" onClick=${b} disabled=${g}>
                Refresh
              </button>
            </div>
          </div>
          ${o&&r`<p className="error">${o}</p>`}
          ${g?r`<p className="muted">Loading config...</p>`:t?r`<pre className="yaml-view">${t}</pre>`:r`<p className="muted">No config data returned.</p>`}
          ${s&&r`<p className="muted">Last fetched ${s}</p>`}
        </section>
      </main>
    `}function Ve(){let{addToast:e}=H(),{groups:t,setGroups:a,refreshFromStorage:o,importFromRunningConfig:i}=pe(),[g,n]=h(""),[s,p]=h(""),[d,b]=h({address:"",weight:100,flags:0}),[u,m]=h(""),[y,I]=h(!1);T(()=>{if(s){if(!t[s]){let f=Object.keys(t);p(f[0]||"")}}else{let f=Object.keys(t);f.length>0&&p(f[0])}},[t,s]);let v=f=>{f.preventDefault();let w=g.trim();if(!w){m("Provide a group name.");return}if(t[w]){m("That group already exists.");return}a({...t,[w]:[]}),n(""),p(w),m(""),e(`Target group "${w}" created.`,"success")},k=f=>{let w={...t};delete w[f],a(w),e(`Target group "${f}" removed.`,"success")},S=f=>{if(f.preventDefault(),!s){m("Select a group to add a real.");return}let w=ce(d);if(!w){m("Provide a valid real address.");return}let P=t[s]||[],x=P.some(G=>G.address===w.address)?P.map(G=>G.address===w.address?w:G):P.concat(w);a({...t,[s]:x}),b({address:"",weight:100,flags:0}),m(""),e("Real saved to target group.","success")},R=f=>{if(!s)return;let P=(t[s]||[]).filter(x=>x.address!==f);a({...t,[s]:P})},D=(f,w)=>{if(!s)return;let x=(t[s]||[]).map(G=>G.address===f?{...G,...w}:G);a({...t,[s]:x})};return r`
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
              <button className="btn ghost" type="button" onClick=${async()=>{I(!0);try{await i(),e("Imported target groups from running config.","success"),m("")}catch(f){m(f.message||"Failed to import target groups."),e(f.message||"Import failed.","error")}finally{I(!1)}}} disabled=${y}>
                ${y?"Importing...":"Import from running config"}
              </button>
            </div>
          </div>
          ${u&&r`<p className="error">${u}</p>`}
          <form className="form" onSubmit=${v}>
            <div className="form-row">
              <label className="field">
                <span>New group name</span>
                <input
                  value=${g}
                  onInput=${f=>n(f.target.value)}
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
                  onChange=${f=>p(f.target.value)}
                  disabled=${Object.keys(t).length===0}
                >
                  ${Object.keys(t).map(f=>r`<option value=${f}>${f}</option>`)}
                </select>
              </label>
              ${s&&r`<button className="btn danger" type="button" onClick=${()=>k(s)}>
                Delete group
              </button>`}
            </div>
          </div>
          ${s?r`
                <table className="table">
                  <thead>
                    <tr>
                      <th>Address</th>
                      <th>Weight</th>
                      <th></th>
                    </tr>
                  </thead>
                  <tbody>
                    ${(t[s]||[]).map(f=>r`
                        <tr>
                          <td>${f.address}</td>
                          <td>
                            <input
                              className="inline-input"
                              type="number"
                              min="0"
                              value=${f.weight}
                              onInput=${w=>D(f.address,{weight:Number(w.target.value)})}
                            />
                          </td>
                          <td>
                            <button
                              className="btn ghost"
                              type="button"
                              onClick=${()=>R(f.address)}
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
                        onInput=${f=>b({...d,address:f.target.value})}
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
                        onInput=${f=>b({...d,weight:f.target.value})}
                        required
                      />
                    </label>
                  </div>
                  <button className="btn secondary" type="submit">Save real</button>
                </form>
              `:r`<p className="muted">No groups yet. Create one to add reals.</p>`}
        </section>
      </main>
    `}function Be(){let[e,t]=h({initialized:!1,ready:!1}),[a,o]=h([]),i=U({}),g=(s,p="info")=>{let d=`${Date.now()}-${Math.random().toString(16).slice(2)}`;o(b=>b.concat({id:d,message:s,kind:p})),i.current[d]=setTimeout(()=>{o(b=>b.filter(u=>u.id!==d)),delete i.current[d]},4e3)},n=s=>{i.current[s]&&(clearTimeout(i.current[s]),delete i.current[s]),o(p=>p.filter(d=>d.id!==s))};return T(()=>{let s=!0,p=async()=>{try{let b=await N.get("/lb/status");s&&t(b||{initialized:!1,ready:!1})}catch{s&&t({initialized:!1,ready:!1})}};p();let d=setInterval(p,5e3);return()=>{s=!1,clearInterval(d)}},[]),r`
      <${se}>
        <${ge}>
          <${X.Provider} value=${{addToast:g}}>
            <${Fe} status=${e} />
            <${re}>
              <${q} path="/" element=${r`<${xe} />`} />
              <${q} path="/vips/:vipId" element=${r`<${Ge} />`} />
              <${q} path="/vips/:vipId/stats" element=${r`<${Ae} />`} />
              <${q} path="/target-groups" element=${r`<${Ve} />`} />
              <${q} path="/stats/global" element=${r`<${Le} />`} />
              <${q} path="/stats/real" element=${r`<${qe} />`} />
              <${q} path="/config" element=${r`<${Oe} />`} />
            </${re}>
            <${ke} toasts=${a} onDismiss=${n} />
          </${X.Provider}>
        </${ge}>
      </${se}>
    `}ReactDOM.createRoot(document.getElementById("root")).render(r`<${Be} />`)})();})();
