(()=>{(()=>{let{useEffect:T,useMemo:V,useRef:U,useState:f,useContext:_e}=React,{BrowserRouter:te,Routes:se,Route:q,NavLink:F,Link:j,useParams:re,useNavigate:Ie}=ReactRouterDOM,s=htm.bind(React.createElement),K=React.createContext({addToast:()=>{}}),M=[{label:"NO_SPORT",value:1},{label:"NO_LRU",value:2},{label:"QUIC_VIP",value:4},{label:"DPORT_HASH",value:8},{label:"SRC_ROUTING",value:16},{label:"LOCAL_VIP",value:32},{label:"GLOBAL_LRU",value:64},{label:"HASH_SRC_DST_PORT",value:128},{label:"UDP_STABLE_ROUTING_VIP",value:256},{label:"UDP_FLOW_MIGRATION",value:512}];function H(){return _e(K)}function Se(e){return String(e).replace(/[^a-z0-9_-]/gi,"_")}function Re(e,a,t){let r=Number(e)||0,c=Number(a)||0;return t?r|c:r&~c}function Ce(e,a){let t=Number(e)||0;return a.filter(r=>(t&r.value)!==0)}function oe({mask:e,options:a,showStatus:t=!1,emptyLabel:r="None"}){let c=Number(e)||0,h=t?a:Ce(c,a),l=t?2:1;return s`
      <table className="table flag-table">
        <thead>
          <tr>
            <th>Flag</th>
            ${t?s`<th>Enabled</th>`:null}
          </tr>
        </thead>
        <tbody>
          ${h.length===0?s`<tr><td colspan=${l} className="muted">${r}</td></tr>`:h.map(o=>{let p=(c&o.value)!==0;return s`
                  <tr>
                    <td>${o.label}</td>
                    ${t?s`<td>${p?"Yes":"No"}</td>`:null}
                  </tr>
                `})}
        </tbody>
      </table>
    `}function ne({options:e,value:a,onChange:t,name:r}){let c=Number(a)||0,h=Se(r||"flags");return s`
      <div className="flag-selector">
        <table className="table flag-table">
          <thead>
            <tr>
              <th>Enabled</th>
              <th>Flag</th>
            </tr>
          </thead>
          <tbody>
            ${e.map(l=>{let o=`${h}-${l.value}`,p=(c&l.value)===l.value;return s`
                <tr>
                  <td>
                    <input
                      id=${o}
                      type="checkbox"
                      checked=${p}
                      onChange=${u=>t(Re(c,l.value,u.target.checked))}
                    />
                  </td>
                  <td>
                    <label className="flag-option" htmlFor=${o}>${l.label}</label>
                  </td>
                </tr>
              `})}
          </tbody>
        </table>
      </div>
    `}let N={base:"/api/v1",async request(e,a={}){let t={method:a.method||"GET",headers:{"Content-Type":"application/json"}},r=`${N.base}${e}`;if(a.body!==void 0&&a.body!==null)if(t.method==="GET"){let l=new URLSearchParams;Object.entries(a.body).forEach(([p,u])=>{if(u!=null){if(Array.isArray(u)){u.forEach(b=>l.append(p,String(b)));return}if(typeof u=="object"){l.set(p,JSON.stringify(u));return}l.set(p,String(u))}});let o=l.toString();o&&(r+=`${r.includes("?")?"&":"?"}${o}`)}else t.body=JSON.stringify(a.body);let c=await fetch(r,t),h;try{h=await c.json()}catch{throw new Error("invalid JSON response")}if(!c.ok)throw new Error(h?.error?.message||`HTTP ${c.status}`);if(!h.success){let l=h.error?.message||"request failed";throw new Error(l)}return h.data},get(e,a){return N.request(e,{method:"GET",body:a})},post(e,a){return N.request(e,{method:"POST",body:a})},put(e,a){return N.request(e,{method:"PUT",body:a})},del(e,a){return N.request(e,{method:"DELETE",body:a})}},le="vatran_target_groups";function ie(e){if(!e||!e.address)return null;let a=String(e.address).trim();if(!a)return null;let t=Number(e.weight),r=Number(e.flags??0);return{address:a,weight:Number.isFinite(t)?t:0,flags:Number.isFinite(r)?r:0}}function ce(e){if(!e||typeof e!="object")return{};let a={};return Object.entries(e).forEach(([t,r])=>{let c=String(t).trim();if(!c)return;let h=Array.isArray(r)?r.map(ie).filter(Boolean):[],l=[],o=new Set;h.forEach(p=>{o.has(p.address)||(o.add(p.address),l.push(p))}),a[c]=l}),a}function X(){if(typeof localStorage>"u")return{};try{let e=localStorage.getItem(le);return e?ce(JSON.parse(e)):{}}catch{return{}}}function de(e){if(!(typeof localStorage>"u"))try{localStorage.setItem(le,JSON.stringify(e))}catch{}}function Te(e,a){let t={...e};return Object.entries(a||{}).forEach(([r,c])=>{t[r]||(t[r]=c)}),t}function ue(){let[e,a]=f(()=>X());return T(()=>{de(e)},[e]),{groups:e,setGroups:a,refreshFromStorage:()=>{a(X())},importFromRunningConfig:async()=>{let c=await N.get("/config/export/json"),h=ce(c?.target_groups||{}),l=Te(X(),h);return a(l),de(l),l}}}function W(e){return`${encodeURIComponent(e.address)}:${e.port}:${e.proto}`}function Y(e){let a=e.split(":"),t=Number(a.pop()||0),r=Number(a.pop()||0);return{address:decodeURIComponent(a.join(":")),port:r,proto:t}}function Pe(e,a,t=[]){let[r,c]=f(null),[h,l]=f(""),[o,p]=f(!0);return T(()=>{let u=!0,b=async()=>{try{let m=await e();u&&(c(m),l(""),p(!1))}catch(m){u&&(l(m.message||"request failed"),p(!1))}};b();let d=setInterval(b,a);return()=>{u=!1,clearInterval(d)}},t),{data:r,error:h,loading:o}}function J({path:e,body:a,intervalMs:t=1e3,limit:r=60}){let[c,h]=f([]),[l,o]=f(""),p=V(()=>JSON.stringify(a||{}),[a]);return T(()=>{if(a===null)return h([]),o(""),()=>{};let u=!0,b=async()=>{try{let m=await N.get(e,a);if(!u)return;let y=new Date().toLocaleTimeString();h(I=>I.concat({label:y,...m}).slice(-r)),o("")}catch(m){u&&o(m.message||"request failed")}};b();let d=setInterval(b,t);return()=>{u=!1,clearInterval(d)}},[e,p,t,r]),{points:c,error:l}}function pe({title:e,points:a,keys:t,diff:r=!1,height:c=120,showTitle:h=!1}){let l=U(null),o=U(null);return T(()=>{if(!l.current)return;o.current||(o.current=new Chart(l.current,{type:"line",data:{labels:[],datasets:[]},options:{responsive:!0,maintainAspectRatio:!1,animation:!1,scales:{x:{grid:{display:!1}},y:{beginAtZero:!r}},plugins:{legend:{display:!0,position:"bottom"},title:{display:h&&!!e,text:e}}}}));let p=o.current,u=a.map(b=>b.label);return p.data.labels=u,p.data.datasets=t.map(b=>{let d=a.map(y=>y[b.field]||0),m=r?d.map((y,I)=>I===0?0:y-d[I-1]):d;return{label:b.label,data:m,borderColor:b.color,backgroundColor:b.fill,borderWidth:2,tension:.3}}),p.options.scales.y.beginAtZero=!r,p.options.plugins.title.display=h&&!!e,p.options.plugins.title.text=e||"",p.update(),()=>{}},[a,t,e,r,h]),T(()=>()=>{o.current&&(o.current.destroy(),o.current=null)},[]),s`<canvas ref=${l} height=${c}></canvas>`}function Z({title:e,points:a,keys:t,diff:r=!1,inlineTitle:c=!0}){let[h,l]=f(!1);return s`
      <div className="chart-wrap">
        <button className="chart-zoom-button" type="button" onClick=${()=>l(!0)}>
          <span className="zoom-icon" aria-hidden="true">+</span>
          Zoom
        </button>
        <div className="chart-click" onClick=${()=>l(!0)}>
          <${pe}
            title=${e}
            points=${a}
            keys=${t}
            diff=${r}
            height=${120}
            showTitle=${c&&!!e}
          />
        </div>
        ${h&&s`
          <div className="chart-overlay" onClick=${()=>l(!1)}>
            <div className="chart-modal" onClick=${o=>o.stopPropagation()}>
              <div className="row chart-modal-header">
                <div>
                  <h3>${e||"Chart"}</h3>
                  ${r?s`<p className="muted">Per-second delta.</p>`:""}
                </div>
                <button className="btn ghost" onClick=${()=>l(!1)}>
                  Close
                </button>
              </div>
              <div className="chart-zoom">
                <${pe}
                  title=${e}
                  points=${a}
                  keys=${t}
                  diff=${r}
                  height=${360}
                  showTitle=${!1}
                />
              </div>
            </div>
          </div>
        `}
      </div>
    `}function me({children:e}){return e}function ke({toasts:e,onDismiss:a}){return s`
      <div className="toast-stack">
        ${e.map(t=>s`
            <div className=${`toast ${t.kind}`}>
              <span>${t.message}</span>
              <button className="toast-close" onClick=${()=>a(t.id)}>
                ×
              </button>
            </div>
          `)}
      </div>
    `}function Fe({status:e}){return s`
      <header>
        <div>
          <div style=${{fontSize:20,fontWeight:700}}>Vatran</div>
          <div className="status-pill">
            <span className=${`dot ${e.ready?"ok":""}`}></span>
            ${e.ready?"Ready":"Not ready"}
          </div>
        </div>
        <nav>
          <${F} to="/" end className=${({isActive:a})=>a?"active":""}>
            Dashboard
          </${F}>
          <${F} to="/stats/global" className=${({isActive:a})=>a?"active":""}>
            Global stats
          </${F}>
          <${F} to="/stats/real" className=${({isActive:a})=>a?"active":""}>
            Per-real stats
          </${F}>
          <${F}
            to="/target-groups"
            className=${({isActive:a})=>a?"active":""}
          >
            Target groups
          </${F}>
          <${F} to="/config" className=${({isActive:a})=>a?"active":""}>
            Config export
          </${F}>
        </nav>
      </header>
    `}function xe(){let{addToast:e}=H(),[a,t]=f({initialized:!1,ready:!1}),[r,c]=f([]),[h,l]=f(""),[o,p]=f(!1),[u,b]=f(!1),[d,m]=f({main_interface:"",balancer_prog_path:"",healthchecking_prog_path:"",default_mac:"",local_mac:"",root_map_path:"",root_map_pos:2,katran_src_v4:"",katran_src_v6:"",use_root_map:!1,max_vips:512,max_reals:4096,hash_func:0}),[y,I]=f({address:"",port:80,proto:6,flags:0}),v=async()=>{try{let n=await N.get("/lb/status"),g=await N.get("/vips"),w=await Promise.all((g||[]).map(async P=>{try{let x=await N.get("/vips/flags",{address:P.address,port:P.port,proto:P.proto});return{...P,flags:x?.flags??0}}catch{return{...P,flags:null}}}));t(n||{initialized:!1,ready:!1}),c(w),l("")}catch(n){l(n.message||"request failed")}};T(()=>{let n=!0;return(async()=>{n&&await v()})(),()=>{n=!1}},[]);let k=async n=>{n.preventDefault();try{let g={...d,root_map_pos:d.root_map_pos===""?void 0:Number(d.root_map_pos),max_vips:Number(d.max_vips),max_reals:Number(d.max_reals),hash_func:Number(d.hash_func)};await N.post("/lb/create",g),l(""),p(!1),e("Load balancer initialized.","success"),await v()}catch(g){l(g.message||"request failed"),e(g.message||"Initialize failed.","error")}},S=async n=>{n.preventDefault();try{await N.post("/vips",{...y,port:Number(y.port),proto:Number(y.proto),flags:Number(y.flags||0)}),I({address:"",port:80,proto:6,flags:0}),l(""),b(!1),e("VIP created.","success"),await v()}catch(g){l(g.message||"request failed"),e(g.message||"VIP create failed.","error")}},R=async()=>{try{await N.post("/lb/load-bpf-progs"),l(""),e("BPF programs loaded.","success"),await v()}catch(n){l(n.message||"request failed"),e(n.message||"Load BPF programs failed.","error")}},D=async()=>{try{await N.post("/lb/attach-bpf-progs"),l(""),e("BPF programs attached.","success"),await v()}catch(n){l(n.message||"request failed"),e(n.message||"Attach BPF programs failed.","error")}};return s`
      <main>
        <section className="card">
          <h2>Load balancer</h2>
          <p>Initialized: ${a.initialized?"yes":"no"}</p>
          <p>Ready: ${a.ready?"yes":"no"}</p>
          <div className="row">
            ${!a.initialized&&s`
              <button className="btn" onClick=${()=>p(n=>!n)}>
                ${o?"Close":"Initialize"}
              </button>
            `}
            <button className="btn secondary" onClick=${()=>b(n=>!n)}>
              ${u?"Close":"Create VIP"}
            </button>
          </div>
          ${!a.ready&&s`
            <div className="row" style=${{marginTop:12}}>
              <button
                className="btn ghost"
                disabled=${!a.initialized}
                onClick=${R}
              >
                Load BPF Programs
              </button>
              <button
                className="btn ghost"
                disabled=${!a.initialized}
                onClick=${D}
              >
                Attach BPF Programs
              </button>
            </div>
          `}
          ${o&&s`
            <form className="form" onSubmit=${k}>
              <div className="form-row">
                <label className="field">
                  <span>Main interface</span>
                  <input
                    value=${d.main_interface}
                    onInput=${n=>m({...d,main_interface:n.target.value})}
                    placeholder="eth0"
                    required
                  />
                </label>
                <label className="field">
                  <span>Balancer prog path</span>
                  <input
                    value=${d.balancer_prog_path}
                    onInput=${n=>m({...d,balancer_prog_path:n.target.value})}
                    placeholder="/path/to/bpf.o"
                    required
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Healthchecking path</span>
                  <input
                    value=${d.healthchecking_prog_path}
                    onInput=${n=>m({...d,healthchecking_prog_path:n.target.value})}
                    placeholder="/path/to/hc.o"
                  />
                </label>
                <label className="field">
                  <span>Default MAC</span>
                  <input
                    value=${d.default_mac}
                    onInput=${n=>m({...d,default_mac:n.target.value})}
                    placeholder="aa:bb:cc:dd:ee:ff"
                    required
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Local MAC</span>
                  <input
                    value=${d.local_mac}
                    onInput=${n=>m({...d,local_mac:n.target.value})}
                    placeholder="11:22:33:44:55:66"
                    required=${d.healthchecking_prog_path?.trim()!==""}
                  />
                </label>
                <label className="field">
                  <span>Hash function</span>
                  <input
                    type="number"
                    min="0"
                    value=${d.hash_func}
                    onInput=${n=>m({...d,hash_func:n.target.value})}
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Root map path</span>
                  <input
                    value=${d.root_map_path}
                    onInput=${n=>m({...d,root_map_path:n.target.value})}
                    placeholder="/sys/fs/bpf/root_map"
                  />
                </label>
                <label className="field">
                  <span>Root map position</span>
                  <input
                    type="number"
                    min="0"
                    value=${d.root_map_pos}
                    onInput=${n=>m({...d,root_map_pos:n.target.value})}
                  />
                </label>
              </div>
              <div className="form-row">
                <label className="field">
                  <span>Katran src v4</span>
                  <input
                    value=${d.katran_src_v4}
                    onInput=${n=>m({...d,katran_src_v4:n.target.value})}
                    placeholder="10.0.0.1"
                  />
                </label>
                <label className="field">
                  <span>Katran src v6</span>
                  <input
                    value=${d.katran_src_v6}
                    onInput=${n=>m({...d,katran_src_v6:n.target.value})}
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
                    value=${d.max_vips}
                    onInput=${n=>m({...d,max_vips:n.target.value})}
                  />
                </label>
                <label className="field">
                  <span>Max reals</span>
                  <input
                    type="number"
                    min="1"
                    value=${d.max_reals}
                    onInput=${n=>m({...d,max_reals:n.target.value})}
                  />
                </label>
              </div>
              <label className="field checkbox">
                <input
                  type="checkbox"
                  checked=${d.use_root_map}
                  onChange=${n=>m({...d,use_root_map:n.target.checked})}
                />
                <span>Use root map</span>
              </label>
              <button className="btn" type="submit">Initialize LB</button>
            </form>
          `}
          ${u&&s`
            <form className="form" onSubmit=${S}>
              <div className="form-row">
                <label className="field">
                  <span>VIP address</span>
                  <input
                    value=${y.address}
                    onInput=${n=>I({...y,address:n.target.value})}
                    placeholder="1.2.3.4"
                    required
                  />
                </label>
                <label className="field">
                  <span>Port</span>
                  <input
                    type="number"
                    value=${y.port}
                    onInput=${n=>I({...y,port:n.target.value})}
                    required
                  />
                </label>
                <label className="field">
                  <span>Protocol</span>
                  <select
                    value=${y.proto}
                    onChange=${n=>I({...y,proto:n.target.value})}
                  >
                    <option value="6">TCP (6)</option>
                    <option value="17">UDP (17)</option>
                  </select>
                </label>
                <label className="field">
                  <span>Flags</span>
                  <${ne}
                    options=${M}
                    value=${y.flags}
                    name="vip-add"
                    onChange=${n=>I({...y,flags:n})}
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
                  ${r.map(n=>s`
                      <div className="card">
                        <div style=${{fontWeight:600}}>
                          ${n.address}:${n.port} / ${n.proto}
                        </div>
                        <div className="muted" style=${{marginTop:6}}>
                          <div style=${{fontWeight:600,marginBottom:6}}>Flags</div>
                          <${oe}
                            mask=${n.flags}
                            options=${M}
                            emptyLabel=${n.flags===null?"Unknown":"No flags"}
                          />
                        </div>
                        <div className="row" style=${{marginTop:12}}>
                          <${j} className="btn" to=${`/vips/${W(n)}`}>
                            Open
                          </${j}>
                          <${j}
                            className="btn secondary"
                            to=${`/vips/${W(n)}/stats`}
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
    `}function Ge(){let{addToast:e}=H(),a=re(),t=Ie(),r=V(()=>Y(a.vipId),[a.vipId]),[c,h]=f([]),[l,o]=f(""),[p,u]=f(""),[b,d]=f(!0),[m,y]=f({address:"",weight:100,flags:0}),[I,v]=f({}),[k,S]=f(null),[R,D]=f({flags:0,set:!0}),[n,g]=f({hash_function:0}),{groups:w,setGroups:P,refreshFromStorage:x,importFromRunningConfig:G}=ue(),[A,he]=f(""),[be,O]=f(""),[ve,$e]=f(!1),[Ne,ye]=f(""),[Q,we]=f({add:0,update:0,remove:0}),B=async()=>{try{let i=await N.get("/vips/reals",r);h(i||[]);let $={};(i||[]).forEach(L=>{$[L.address]=L.weight}),v($),o(""),d(!1)}catch(i){o(i.message||"request failed"),d(!1)}},ee=async()=>{try{let i=await N.get("/vips/flags",r);S(i?.flags??0),u("")}catch(i){u(i.message||"request failed")}};T(()=>{B(),ee()},[a.vipId]),T(()=>{if(!A){we({add:0,update:0,remove:0});return}let i=w[A]||[],$=new Map(c.map(C=>[C.address,C])),L=new Map(i.map(C=>[C.address,C])),E=0,z=0,_=0;i.forEach(C=>{let ae=$.get(C.address);if(!ae){E+=1;return}(Number(ae.weight)!==Number(C.weight)||Number(ae.flags||0)!==Number(C.flags||0))&&(z+=1)}),c.forEach(C=>{L.has(C.address)||(_+=1)}),we({add:E,update:z,remove:_})},[A,c,w]);let ze=async i=>{try{let $=Number(I[i.address]);await N.post("/vips/reals",{vip:r,real:{address:i.address,weight:$,flags:i.flags||0}}),await B(),e("Real weight updated.","success")}catch($){o($.message||"request failed"),e($.message||"Update failed.","error")}},Ue=async i=>{try{await N.del("/vips/reals",{vip:r,real:{address:i.address,weight:i.weight,flags:i.flags||0}}),await B(),e("Real removed.","success")}catch($){o($.message||"request failed"),e($.message||"Remove failed.","error")}},je=async i=>{i.preventDefault();try{await N.post("/vips/reals",{vip:r,real:{address:m.address,weight:Number(m.weight),flags:Number(m.flags||0)}}),y({address:"",weight:100,flags:0}),await B(),e("Real added.","success")}catch($){o($.message||"request failed"),e($.message||"Add failed.","error")}},Me=async()=>{if(!A||!w[A]){O("Select a target group to apply.");return}$e(!0),O("");let i=w[A]||[],$=new Map(c.map(_=>[_.address,_])),L=new Map(i.map(_=>[_.address,_])),E=c.filter(_=>!L.has(_.address)),z=i.filter(_=>{let C=$.get(_.address);return C?Number(C.weight)!==Number(_.weight)||Number(C.flags||0)!==Number(_.flags||0):!0});try{E.length>0&&await N.put("/vips/reals/batch",{vip:r,action:1,reals:E.map(_=>({address:_.address,weight:Number(_.weight),flags:Number(_.flags||0)}))}),z.length>0&&await Promise.all(z.map(_=>N.post("/vips/reals",{vip:r,real:{address:_.address,weight:Number(_.weight),flags:Number(_.flags||0)}}))),await B(),e(`Applied target group "${A}".`,"success")}catch(_){O(_.message||"Failed to apply target group."),e(_.message||"Target group apply failed.","error")}finally{$e(!1)}},He=i=>{i.preventDefault();let $=Ne.trim();if(!$){O("Provide a name for the new target group.");return}if(w[$]){O("A target group with that name already exists.");return}let L={...w,[$]:c.map(E=>({address:E.address,weight:Number(E.weight),flags:Number(E.flags||0)}))};P(L),ye(""),he($),O(""),e(`Target group "${$}" saved.`,"success")},We=async()=>{try{await N.del("/vips",r),e("VIP deleted.","success"),t("/")}catch(i){o(i.message||"request failed"),e(i.message||"Delete failed.","error")}},Je=async i=>{i.preventDefault();try{await N.put("/vips/flags",{...r,flag:Number(R.flags||0),set:!!R.set}),await ee(),e("VIP flags updated.","success")}catch($){u($.message||"request failed"),e($.message||"Flag update failed.","error")}},Ke=async i=>{i.preventDefault();try{await N.put("/vips/hash-function",{...r,hash_function:Number(n.hash_function)}),e("Hash function updated.","success")}catch($){u($.message||"request failed"),e($.message||"Hash update failed.","error")}};return s`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>VIP Detail</h2>
              <p className="muted">${r.address}:${r.port} / ${r.proto}</p>
              ${k===null?s`<p className="muted">Flags: —</p>`:s`
                    <div style=${{marginTop:8}}>
                      <${oe}
                        mask=${k}
                        options=${M}
                        showStatus=${!0}
                        emptyLabel="No flags"
                      />
                    </div>
                  `}
            </div>
            <div className="row">
              <button className="btn ghost" onClick=${ee}>Refresh flags</button>
              <button className="btn danger" onClick=${We}>Delete VIP</button>
            </div>
          </div>
          ${l&&s`<p className="error">${l}</p>`}
          ${p&&s`<p className="error">${p}</p>`}
          ${b?s`<p className="muted">Loading reals…</p>`:s`
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
                    ${c.map(i=>s`
                        <tr>
                          <td>
                            <span
                              className=${`dot ${Number(i.weight)>0?"ok":"bad"}`}
                            ></span>
                          </td>
                          <td>${i.address}</td>
                          <td>
                            <input
                              className="inline-input"
                              type="number"
                              min="0"
                              value=${I[i.address]??i.weight}
                              onInput=${$=>v({...I,[i.address]:$.target.value})}
                            />
                          </td>
                          <td className="row">
                            <button className="btn" onClick=${()=>ze(i)}>
                              Update
                            </button>
                            <button className="btn ghost" onClick=${()=>Ue(i)}>
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
                  <${ne}
                    options=${M}
                    value=${R.flags}
                    name="vip-flag-change"
                    onChange=${i=>D({...R,flags:i})}
                  />
                </label>
                <label className="field">
                  <span>Set</span>
                  <select
                    value=${String(R.set)}
                    onChange=${i=>D({...R,set:i.target.value==="true"})}
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
                    value=${n.hash_function}
                    onInput=${i=>g({...n,hash_function:i.target.value})}
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
                  onInput=${i=>y({...m,address:i.target.value})}
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
                  onInput=${i=>y({...m,weight:i.target.value})}
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
                onClick=${async()=>{try{await G(),e("Imported target groups from running config.","success")}catch(i){O(i.message||"Failed to import target groups."),e(i.message||"Import failed.","error")}}}
              >
                Import from running config
              </button>
            </div>
          </div>
          ${be&&s`<p className="error">${be}</p>`}
          <div className="form-row">
            <label className="field">
              <span>Target group</span>
              <select
                value=${A}
                onChange=${i=>he(i.target.value)}
                disabled=${Object.keys(w).length===0}
              >
                <option value="">Select group</option>
                ${Object.keys(w).map(i=>s`<option value=${i}>${i}</option>`)}
              </select>
            </label>
            <label className="field">
              <span>Preview</span>
              <input
                value=${`add ${Q.add} \xB7 update ${Q.update} \xB7 remove ${Q.remove}`}
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
                  onInput=${i=>ye(i.target.value)}
                  placeholder="edge-backends"
                />
              </label>
            </div>
            <button className="btn secondary" type="submit">Save target group</button>
          </form>
        </section>
      </main>
    `}function Ae(){let e=re(),a=V(()=>Y(e.vipId),[e.vipId]),{points:t,error:r}=J({path:"/stats/vip",body:a}),c=V(()=>[{label:"v1",field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.15)"},{label:"v2",field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[]);return s`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>VIP Stats</h2>
              <p className="muted">${a.address}:${a.port} / ${a.proto}</p>
            </div>
          </div>
          ${r&&s`<p className="error">${r}</p>`}
          <${Z} title="Traffic" points=${t} keys=${c} />
        </section>
      </main>
    `}let ge=[{title:"LRU",path:"/stats/lru"},{title:"LRU Miss",path:"/stats/lru/miss"},{title:"LRU Fallback",path:"/stats/lru/fallback"},{title:"LRU Global",path:"/stats/lru/global"},{title:"XDP Total",path:"/stats/xdp/total"},{title:"XDP Pass",path:"/stats/xdp/pass"},{title:"XDP Drop",path:"/stats/xdp/drop"},{title:"XDP Tx",path:"/stats/xdp/tx"}];function fe(e){return Number.isFinite(e)?`${e>0?"+":""}${e}`:"0"}function Ee({title:e,path:a,diff:t=!1}){let{points:r,error:c}=J({path:a}),h=V(()=>[{label:"v1",field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.2)"},{label:"v2",field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[]);return s`
      <div className="card">
        <h3>${e}</h3>
        ${c&&s`<p className="error">${c}</p>`}
        <${Z} title=${e} points=${r} keys=${h} diff=${t} inlineTitle=${!1} />
      </div>
    `}function De({title:e,path:a}){let{points:t,error:r}=J({path:a}),c=t[t.length-1]||{},h=t[t.length-2]||{},l=Number(c.v1??0),o=Number(c.v2??0),p=l-Number(h.v1??0),u=o-Number(h.v2??0);return s`
      <div className="summary-card">
        <div className="summary-title">${e}</div>
        ${r?s`<p className="error">${r}</p>`:s`
              <div className="summary-row">
                <div className="stat">
                  <span className="muted">v1 absolute</span>
                  <strong>${l}</strong>
                </div>
                <div className="stat">
                  <span className="muted">v1 delta/sec</span>
                  <strong className=${p<0?"delta down":"delta up"}>
                    ${fe(p)}
                  </strong>
                </div>
              </div>
              <div className="summary-row">
                <div className="stat">
                  <span className="muted">v2 absolute</span>
                  <strong>${o}</strong>
                </div>
                <div className="stat">
                  <span className="muted">v2 delta/sec</span>
                  <strong className=${u<0?"delta down":"delta up"}>
                    ${fe(u)}
                  </strong>
                </div>
              </div>
            `}
      </div>
    `}function Le(){let{data:e,error:a}=Pe(()=>N.get("/stats/userspace"),1e3,[]);return s`
      <main>
        <section className="card">
          <h2>Global Stats</h2>
          <p className="muted">Polling every second. Charts show per-second deltas.</p>
        </section>
        <section className="grid">
          ${ge.map(t=>s`<${Ee} title=${t.title} path=${t.path} diff=${!0} />`)}
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
            ${ge.map(t=>s`<${De} title=${t.title} path=${t.path} />`)}
          </div>
        </section>
      </main>
    `}function qe(){let[e,a]=f([]),[t,r]=f(""),[c,h]=f([]),[l,o]=f(""),[p,u]=f(null),[b,d]=f("");T(()=>{let v=!0;return(async()=>{try{let S=await N.get("/vips");if(!v)return;a(S||[]),!t&&S&&S.length>0&&r(W(S[0]))}catch(S){v&&d(S.message||"request failed")}})(),()=>{v=!1}},[]),T(()=>{if(!t)return;let v=Y(t),k=!0;return(async()=>{try{let R=await N.get("/vips/reals",v);if(!k)return;h(R||[]),R&&R.length>0?o(D=>D||R[0].address):o(""),d("")}catch(R){k&&d(R.message||"request failed")}})(),()=>{k=!1}},[t]),T(()=>{if(!l){u(null);return}let v=!0;return(async()=>{try{let S=await N.get("/reals/index",{address:l});if(!v)return;u(S?.index??null),d("")}catch(S){v&&d(S.message||"request failed")}})(),()=>{v=!1}},[l]);let{points:m,error:y}=J({path:"/stats/real",body:p!==null?{index:p}:null}),I=V(()=>[{label:"v1",field:"v1",color:"#2f4858",fill:"rgba(47,72,88,0.2)"},{label:"v2",field:"v2",color:"#d97757",fill:"rgba(217,119,87,0.2)"}],[]);return s`
      <main>
        <section className="card">
          <h2>Per-Real Stats</h2>
          <p className="muted">Select a VIP and real address to chart.</p>
          ${b&&s`<p className="error">${b}</p>`}
          <div className="form-row">
            <label className="field">
              <span>VIP</span>
              <select value=${t} onChange=${v=>r(v.target.value)}>
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
                value=${l}
                onChange=${v=>o(v.target.value)}
                disabled=${c.length===0}
              >
                ${c.map(v=>s`
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
          ${y&&s`<p className="error">${y}</p>`}
          ${p===null?s`<p className="muted">Select a real to start polling.</p>`:s`<${Z} points=${m} keys=${I} />`}
        </section>
      </main>
    `}function Oe(){let{addToast:e}=H(),[a,t]=f(""),[r,c]=f(""),[h,l]=f(!0),[o,p]=f(""),u=U(!0),b=async()=>{if(u.current){l(!0),c("");try{let m=await fetch(`${N.base}/config/export`,{headers:{Accept:"application/x-yaml"}});if(!m.ok){let I=`HTTP ${m.status}`;try{I=(await m.json())?.error?.message||I}catch{}throw new Error(I)}let y=await m.text();if(!u.current)return;t(y||""),p(new Date().toLocaleString())}catch(m){u.current&&c(m.message||"request failed")}finally{u.current&&l(!1)}}},d=async()=>{if(a)try{await navigator.clipboard.writeText(a),e("Config copied to clipboard","info")}catch{e("Failed to copy config","error")}};return T(()=>(u.current=!0,b(),()=>{u.current=!1}),[]),s`
      <main>
        <section className="card">
          <div className="section-header">
            <div>
              <h2>Running config</h2>
              <p className="muted">Exported from /api/v1/config/export</p>
            </div>
            <div className="row">
              <button className="btn ghost" onClick=${d} disabled=${!a}>
                Copy YAML
              </button>
              <button className="btn" onClick=${b} disabled=${h}>
                Refresh
              </button>
            </div>
          </div>
          ${r&&s`<p className="error">${r}</p>`}
          ${h?s`<p className="muted">Loading config...</p>`:a?s`<pre className="yaml-view">${a}</pre>`:s`<p className="muted">No config data returned.</p>`}
          ${o&&s`<p className="muted">Last fetched ${o}</p>`}
        </section>
      </main>
    `}function Ve(){let{addToast:e}=H(),{groups:a,setGroups:t,refreshFromStorage:r,importFromRunningConfig:c}=ue(),[h,l]=f(""),[o,p]=f(""),[u,b]=f({address:"",weight:100,flags:0}),[d,m]=f(""),[y,I]=f(!1);T(()=>{if(o){if(!a[o]){let g=Object.keys(a);p(g[0]||"")}}else{let g=Object.keys(a);g.length>0&&p(g[0])}},[a,o]);let v=g=>{g.preventDefault();let w=h.trim();if(!w){m("Provide a group name.");return}if(a[w]){m("That group already exists.");return}t({...a,[w]:[]}),l(""),p(w),m(""),e(`Target group "${w}" created.`,"success")},k=g=>{let w={...a};delete w[g],t(w),e(`Target group "${g}" removed.`,"success")},S=g=>{if(g.preventDefault(),!o){m("Select a group to add a real.");return}let w=ie(u);if(!w){m("Provide a valid real address.");return}let P=a[o]||[],x=P.some(G=>G.address===w.address)?P.map(G=>G.address===w.address?w:G):P.concat(w);t({...a,[o]:x}),b({address:"",weight:100,flags:0}),m(""),e("Real saved to target group.","success")},R=g=>{if(!o)return;let P=(a[o]||[]).filter(x=>x.address!==g);t({...a,[o]:P})},D=(g,w)=>{if(!o)return;let x=(a[o]||[]).map(G=>G.address===g?{...G,...w}:G);t({...a,[o]:x})};return s`
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
              <button className="btn ghost" type="button" onClick=${async()=>{I(!0);try{await c(),e("Imported target groups from running config.","success"),m("")}catch(g){m(g.message||"Failed to import target groups."),e(g.message||"Import failed.","error")}finally{I(!1)}}} disabled=${y}>
                ${y?"Importing...":"Import from running config"}
              </button>
            </div>
          </div>
          ${d&&s`<p className="error">${d}</p>`}
          <form className="form" onSubmit=${v}>
            <div className="form-row">
              <label className="field">
                <span>New group name</span>
                <input
                  value=${h}
                  onInput=${g=>l(g.target.value)}
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
                  value=${o}
                  onChange=${g=>p(g.target.value)}
                  disabled=${Object.keys(a).length===0}
                >
                  ${Object.keys(a).map(g=>s`<option value=${g}>${g}</option>`)}
                </select>
              </label>
              ${o&&s`<button className="btn danger" type="button" onClick=${()=>k(o)}>
                Delete group
              </button>`}
            </div>
          </div>
          ${o?s`
                <table className="table">
                  <thead>
                    <tr>
                      <th>Address</th>
                      <th>Weight</th>
                      <th></th>
                    </tr>
                  </thead>
                  <tbody>
                    ${(a[o]||[]).map(g=>s`
                        <tr>
                          <td>${g.address}</td>
                          <td>
                            <input
                              className="inline-input"
                              type="number"
                              min="0"
                              value=${g.weight}
                              onInput=${w=>D(g.address,{weight:Number(w.target.value)})}
                            />
                          </td>
                          <td>
                            <button
                              className="btn ghost"
                              type="button"
                              onClick=${()=>R(g.address)}
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
                        value=${u.address}
                        onInput=${g=>b({...u,address:g.target.value})}
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
                        onInput=${g=>b({...u,weight:g.target.value})}
                        required
                      />
                    </label>
                  </div>
                  <button className="btn secondary" type="submit">Save real</button>
                </form>
              `:s`<p className="muted">No groups yet. Create one to add reals.</p>`}
        </section>
      </main>
    `}function Be(){let[e,a]=f({initialized:!1,ready:!1}),[t,r]=f([]),c=U({}),h=(o,p="info")=>{let u=`${Date.now()}-${Math.random().toString(16).slice(2)}`;r(b=>b.concat({id:u,message:o,kind:p})),c.current[u]=setTimeout(()=>{r(b=>b.filter(d=>d.id!==u)),delete c.current[u]},4e3)},l=o=>{c.current[o]&&(clearTimeout(c.current[o]),delete c.current[o]),r(p=>p.filter(u=>u.id!==o))};return T(()=>{let o=!0,p=async()=>{try{let b=await N.get("/lb/status");o&&a(b||{initialized:!1,ready:!1})}catch{o&&a({initialized:!1,ready:!1})}};p();let u=setInterval(p,5e3);return()=>{o=!1,clearInterval(u)}},[]),s`
      <${te}>
        <${me}>
          <${K.Provider} value=${{addToast:h}}>
            <${Fe} status=${e} />
            <${se}>
              <${q} path="/" element=${s`<${xe} />`} />
              <${q} path="/vips/:vipId" element=${s`<${Ge} />`} />
              <${q} path="/vips/:vipId/stats" element=${s`<${Ae} />`} />
              <${q} path="/target-groups" element=${s`<${Ve} />`} />
              <${q} path="/stats/global" element=${s`<${Le} />`} />
              <${q} path="/stats/real" element=${s`<${qe} />`} />
              <${q} path="/config" element=${s`<${Oe} />`} />
            </${se}>
            <${ke} toasts=${t} onDismiss=${l} />
          </${K.Provider}>
        </${me}>
      </${te}>
    `}ReactDOM.createRoot(document.getElementById("root")).render(s`<${Be} />`)})();})();
