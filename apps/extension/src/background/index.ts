import { getSettings, updateSettings, incrementBlockedCount } from '../storage/storage';
import { generateDynamicRules } from '../rules/rule-builder';

const DEFAULT_RULESET_ID = 'ruleset_default';

/**
 * Syncs the browser declarativeNetRequest rules with the user's settings.
 */
async function syncRules(): Promise<void> {
  const settings = await getSettings();

  if (!settings.enabled) {
    // Disable default ruleset
    await chrome.declarativeNetRequest.updateEnabledRulesets({
      disableRulesetIds: [DEFAULT_RULESET_ID],
    });

    // Clear any existing dynamic rules
    const existingRules = await chrome.declarativeNetRequest.getDynamicRules();
    const existingIds = existingRules.map((r) => r.id);
    if (existingIds.length > 0) {
      await chrome.declarativeNetRequest.updateDynamicRules({
        removeRuleIds: existingIds,
      });
    }
    return;
  }

  // Ensure default ruleset is enabled
  await chrome.declarativeNetRequest.updateEnabledRulesets({
    enableRulesetIds: [DEFAULT_RULESET_ID],
  });

  // Rebuild dynamic rules
  const dynamicRules = generateDynamicRules(
    settings.customBlocklist,
    settings.customAllowlist,
    chrome.runtime.getURL('src/blocked-page/blocked.html')
  );

  const existingRules = await chrome.declarativeNetRequest.getDynamicRules();
  const existingIds = existingRules.map((r) => r.id);

  await chrome.declarativeNetRequest.updateDynamicRules({
    removeRuleIds: existingIds,
    addRules: dynamicRules,
  });
}

// Extension Lifecycle
chrome.runtime.onInstalled.addListener(async (details) => {
  console.log('[Vajra] Extension installed/updated:', details.reason);
  await syncRules();
});

chrome.runtime.onStartup.addListener(async () => {
  console.log('[Vajra] Browser startup - synchronizing rules');
  await syncRules();
});

// React to storage changes
chrome.storage.onChanged.addListener((changes, namespace) => {
  if (namespace === 'local' && changes.vajra_user_settings) {
    syncRules().catch((err) => console.error('[Vajra] Error syncing rules:', err));
  }
});

// Handle messages from popup and blocked page
chrome.runtime.onMessage.addListener((message, _sender, sendResponse) => {
  if (message.type === 'RECORD_BLOCK') {
    incrementBlockedCount(message.domain, message.url).then((newCount) => {
      sendResponse({ success: true, count: newCount });
    });
    return true; // async response
  }

  if (message.type === 'TOGGLE_PROTECTION') {
    getSettings().then((settings) => {
      const nextState = !settings.enabled;
      updateSettings({ enabled: nextState }).then((updated) => {
        sendResponse({ success: true, enabled: updated.enabled });
      });
    });
    return true;
  }

  if (message.type === 'GET_STATUS') {
    getSettings().then((settings) => {
      sendResponse({
        enabled: settings.enabled,
        blockedCount: settings.blockedCount,
        customBlocklistCount: settings.customBlocklist.length,
        customAllowlistCount: settings.customAllowlist.length,
      });
    });
    return true;
  }

  return false;
});
