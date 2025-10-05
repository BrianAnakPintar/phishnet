// Read query params (original URL and reason) and wire up buttons
(function(){
  function qs(){
    try{ return new URLSearchParams(location.search); }catch(e){ return new URLSearchParams(); }
  }
  const params = qs();
  const orig = params.get('url') || '';
  const reason = params.get('reason') || '';

  if(orig) {
    const el = document.querySelector('#orig span');
    el.textContent = orig;
    el.title = orig;
  }

  const reasonLogEl = document.getElementById('reasonLog');
  const detailsEl = document.getElementById('reasonDetails');

  function formatReason() {
    if(!reason) return 'No details provided.';
    // If reason looks like JSON, pretty-print it
    try {
      const parsed = JSON.parse(reason);
      return JSON.stringify(parsed, null, 2);
    } catch (e) {
      return reason;
    }
  }

  // Lazy populate the log only when the user expands the details
  if(detailsEl) {
    // If already open (unlikely), populate immediately
    if(detailsEl.open) {
      reasonLogEl.textContent = formatReason();
    }
    detailsEl.addEventListener('toggle', ()=>{
      if(detailsEl.open) {
        reasonLogEl.textContent = formatReason();
      }
    });
  } else {
    // Fallback: populate immediately if no details element
    reasonLogEl.textContent = formatReason();
  }

  // copy-log removed; log is lazy-loaded in the details element

  document.getElementById('goBack').addEventListener('click', ()=>{
    // Try history.back first, else close the tab
    if(history.length>1) history.back(); else window.close();
  });

  document.getElementById('continue').addEventListener('click', ()=>{
    if(!orig) return;
    try{
      const u = new URL(orig);
      // Add bypass param so content script won't redirect again
      u.searchParams.set('gw_bypass','1');
      // Use replace to avoid keeping warning page in history
      window.location.replace(u.toString());
    }catch(e){
      // if invalid URL, do nothing
    }
  });
})();
