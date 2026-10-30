import { addCustomAllow } from '../storage/storage';

const quotes = [
  '"True freedom is the mastery of oneself."',
  '"Between stimulus and response there is a space. In that space is our power to choose."',
  '"Discipline is choosing between what you want now and what you want most."',
  '"A quiet mind brings unwavering focus and resilience."',
];

function init(): void {
  const urlParams = new URLSearchParams(window.location.search);
  const domain = urlParams.get('domain') || 'Restricted Content';

  const domainEl = document.getElementById('blocked-domain');
  if (domainEl) {
    domainEl.textContent = domain;
  }

  // Display random mindful quote
  const quoteEl = document.getElementById('quote');
  if (quoteEl) {
    const randomQuote = quotes[Math.floor(Math.random() * quotes.length)];
    quoteEl.textContent = randomQuote;
  }

  // Record block event with background service worker
  if (typeof chrome !== 'undefined' && chrome.runtime && chrome.runtime.sendMessage) {
    chrome.runtime.sendMessage({
      type: 'RECORD_BLOCK',
      domain,
      url: window.location.href,
    });
  }

  // Buttons
  const closeBtn = document.getElementById('close-tab-btn');
  closeBtn?.addEventListener('click', () => {
    if (typeof chrome !== 'undefined' && chrome.tabs) {
      chrome.tabs.getCurrent((tab) => {
        if (tab?.id) {
          chrome.tabs.remove(tab.id);
        } else {
          window.close();
        }
      });
    } else {
      window.close();
    }
  });

  const goBackBtn = document.getElementById('go-back-btn');
  goBackBtn?.addEventListener('click', () => {
    if (window.history.length > 1) {
      window.history.back();
    } else {
      window.location.href = 'https://google.com';
    }
  });

  const allowBtn = document.getElementById('allow-domain-btn');
  allowBtn?.addEventListener('click', async () => {
    if (domain === 'Restricted Content') return;
    const confirmed = window.confirm(`Are you sure you want to allowlist "${domain}"?`);
    if (confirmed) {
      await addCustomAllow(domain);
      alert(`"${domain}" has been added to your allowlist.`);
      window.location.href = `https://${domain}`;
    }
  });
}

document.addEventListener('DOMContentLoaded', init);
