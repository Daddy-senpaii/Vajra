import {
  getSettings,
  updateSettings,
  addCustomBlock,
  addCustomAllow,
} from '../storage/storage';
import { normalizeDomain } from '../domain/normalize';

let currentTabDomain = '';

async function initPopup(): Promise<void> {
  const settings = await getSettings();
  renderSettings(settings);

  // Setup toggle listener
  const toggle = document.getElementById('protection-toggle') as HTMLInputElement;
  toggle?.addEventListener('change', async () => {
    const isChecked = toggle.checked;
    await updateSettings({ enabled: isChecked });
    const fresh = await getSettings();
    renderSettings(fresh);
  });

  // Query current tab
  if (typeof chrome !== 'undefined' && chrome.tabs && chrome.tabs.query) {
    chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
      const activeTab = tabs[0];
      if (activeTab?.url) {
        const norm = normalizeDomain(activeTab.url);
        if (norm.isValid) {
          currentTabDomain = norm.domain;
          const currentDomainEl = document.getElementById('current-domain');
          if (currentDomainEl) {
            currentDomainEl.textContent = norm.domain;
          }
        }
      }
    });
  }

  // Quick actions on current site
  const blockCurrentBtn = document.getElementById('btn-block-current');
  blockCurrentBtn?.addEventListener('click', async () => {
    if (!currentTabDomain) return;
    const ok = await addCustomBlock(currentTabDomain);
    if (ok) {
      showFeedback(`Blocked ${currentTabDomain}`, 'success');
      refreshStats();
    }
  });

  const allowCurrentBtn = document.getElementById('btn-allow-current');
  allowCurrentBtn?.addEventListener('click', async () => {
    if (!currentTabDomain) return;
    const ok = await addCustomAllow(currentTabDomain);
    if (ok) {
      showFeedback(`Allowed ${currentTabDomain}`, 'success');
      refreshStats();
    }
  });

  // Manual domain actions
  const inputEl = document.getElementById('custom-domain-input') as HTMLInputElement;
  const addBlockBtn = document.getElementById('btn-add-block');
  addBlockBtn?.addEventListener('click', async () => {
    const val = inputEl?.value.trim();
    if (!val) return;
    const norm = normalizeDomain(val);
    if (!norm.isValid) {
      showFeedback(norm.error || 'Invalid domain format', 'error');
      return;
    }
    await addCustomBlock(norm.domain);
    inputEl.value = '';
    showFeedback(`Added ${norm.domain} to blocklist`, 'success');
    refreshStats();
  });

  const addAllowBtn = document.getElementById('btn-add-allow');
  addAllowBtn?.addEventListener('click', async () => {
    const val = inputEl?.value.trim();
    if (!val) return;
    const norm = normalizeDomain(val);
    if (!norm.isValid) {
      showFeedback(norm.error || 'Invalid domain format', 'error');
      return;
    }
    await addCustomAllow(norm.domain);
    inputEl.value = '';
    showFeedback(`Added ${norm.domain} to allowlist`, 'success');
    refreshStats();
  });
}

function renderSettings(settings: { enabled: boolean; blockedCount: number; customBlocklist: string[]; customAllowlist: string[] }): void {
  const toggle = document.getElementById('protection-toggle') as HTMLInputElement;
  if (toggle) toggle.checked = settings.enabled;

  const badge = document.getElementById('status-badge');
  const subtitle = document.getElementById('toggle-subtitle');

  if (settings.enabled) {
    if (badge) {
      badge.textContent = 'Active';
      badge.className = 'badge active';
    }
    if (subtitle) subtitle.textContent = 'Blocking adult content';
  } else {
    if (badge) {
      badge.textContent = 'Disabled';
      badge.className = 'badge inactive';
    }
    if (subtitle) subtitle.textContent = 'Protection paused';
  }

  const blockedCountEl = document.getElementById('blocked-count');
  if (blockedCountEl) {
    blockedCountEl.textContent = settings.blockedCount.toString();
  }

  const rulesCountEl = document.getElementById('custom-rules-count');
  if (rulesCountEl) {
    const totalCustom = settings.customBlocklist.length + settings.customAllowlist.length;
    rulesCountEl.textContent = totalCustom.toString();
  }
}

async function refreshStats(): Promise<void> {
  const settings = await getSettings();
  renderSettings(settings);
}

function showFeedback(msg: string, type: 'success' | 'error'): void {
  const el = document.getElementById('rule-feedback');
  if (el) {
    el.textContent = msg;
    el.className = `feedback-msg ${type}`;
    setTimeout(() => {
      if (el.textContent === msg) {
        el.textContent = '';
        el.className = 'feedback-msg';
      }
    }, 3000);
  }
}

document.addEventListener('DOMContentLoaded', initPopup);
