// Background service worker: receives messages from content script with current URL, queries local PhishNet server, and forwards result back.
self.addEventListener('message', (event) => {
  // Not used; content script uses fetch directly via extension origin.
});

// Provide a helper via chrome.runtime.onMessage to allow content scripts to request scans through the service worker (in case of CORS)
chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message && message.type === 'scan') {
    const url = message.url;
    fetch(`http://localhost:8080/scan?url=${encodeURIComponent(url)}`)
      .then(resp => resp.json())
      .then(data => sendResponse({ok: true, data}))
      .catch(err => sendResponse({ok: false, error: err.toString()}));
    return true; // indicate async response
  }
});
