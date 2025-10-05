// Content script: sends the current page URL to the background script, gets scan result, and redirects to a warning page if blocked.
(function() {
  const url = window.location.href;

  // Allow user to bypass the warning by adding gw_bypass=1 to the URL
  try {
    const tmp = new URL(url);
    if (tmp.searchParams.has('gw_bypass')) return;
  } catch (e) {
    // ignore invalid URL parsing
  }

  // Avoid redirect loop: if we're already on the extension's warning page, do nothing
  const warningPageUrl = chrome.runtime.getURL('warning.html');
  if (url.startsWith(warningPageUrl)) return;

  // Ask service worker to scan
  chrome.runtime.sendMessage({type: 'scan', url}, (resp) => {
    if (!resp || !resp.ok) return;
    const data = resp.data;
    if (!data) return;
    if (data.allowed === false) {
      // Redirect to the extension warning page and pass the original URL and reason as query params
      const params = new URLSearchParams();
      params.set('url', url);
      if (data.reason) params.set('reason', data.reason);
      const target = warningPageUrl + '?' + params.toString();

      // Use replace so the blocked page isn't left in the back/forward history
      window.location.replace(target);
    }
  });
})();
