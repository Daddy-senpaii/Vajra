import { UserSettings } from '../types';
import { normalizeDomain } from '../domain/normalize';

const STORAGE_KEY = 'vajra_user_settings';

export const DEFAULT_SETTINGS: UserSettings = {
  enabled: true,
  blockedCount: 0,
  customBlocklist: [],
  customAllowlist: [],
};

export async function getSettings(): Promise<UserSettings> {
  return new Promise((resolve) => {
    if (typeof chrome === 'undefined' || !chrome.storage || !chrome.storage.local) {
      resolve({ ...DEFAULT_SETTINGS });
      return;
    }

    chrome.storage.local.get([STORAGE_KEY], (result) => {
      if (chrome.runtime.lastError || !result[STORAGE_KEY]) {
        resolve({ ...DEFAULT_SETTINGS });
      } else {
        resolve({ ...DEFAULT_SETTINGS, ...result[STORAGE_KEY] });
      }
    });
  });
}

export async function updateSettings(partial: Partial<UserSettings>): Promise<UserSettings> {
  const current = await getSettings();
  const updated: UserSettings = { ...current, ...partial };

  return new Promise((resolve) => {
    if (typeof chrome === 'undefined' || !chrome.storage || !chrome.storage.local) {
      resolve(updated);
      return;
    }

    chrome.storage.local.set({ [STORAGE_KEY]: updated }, () => {
      resolve(updated);
    });
  });
}

export async function incrementBlockedCount(domain?: string, url?: string): Promise<number> {
  const current = await getSettings();
  const nextCount = current.blockedCount + 1;

  await updateSettings({
    blockedCount: nextCount,
    lastBlockedDomain: domain || current.lastBlockedDomain,
    lastBlockedUrl: url || current.lastBlockedUrl,
    lastBlockedTimestamp: Date.now(),
  });

  return nextCount;
}

export async function addCustomBlock(rawDomain: string): Promise<boolean> {
  const norm = normalizeDomain(rawDomain);
  if (!norm.isValid) return false;

  const current = await getSettings();
  if (current.customBlocklist.includes(norm.domain)) {
    return true; // Already exists
  }

  // Remove from allowlist if present to avoid conflicting rules
  const updatedAllow = current.customAllowlist.filter((d) => d !== norm.domain);
  const updatedBlock = [...current.customBlocklist, norm.domain];

  await updateSettings({
    customBlocklist: updatedBlock,
    customAllowlist: updatedAllow,
  });

  return true;
}

export async function removeCustomBlock(rawDomain: string): Promise<boolean> {
  const norm = normalizeDomain(rawDomain);
  if (!norm.isValid) return false;

  const current = await getSettings();
  const updatedBlock = current.customBlocklist.filter((d) => d !== norm.domain);

  await updateSettings({
    customBlocklist: updatedBlock,
  });

  return true;
}

export async function addCustomAllow(rawDomain: string): Promise<boolean> {
  const norm = normalizeDomain(rawDomain);
  if (!norm.isValid) return false;

  const current = await getSettings();
  if (current.customAllowlist.includes(norm.domain)) {
    return true;
  }

  // Remove from blocklist if present
  const updatedBlock = current.customBlocklist.filter((d) => d !== norm.domain);
  const updatedAllow = [...current.customAllowlist, norm.domain];

  await updateSettings({
    customBlocklist: updatedBlock,
    customAllowlist: updatedAllow,
  });

  return true;
}

export async function removeCustomAllow(rawDomain: string): Promise<boolean> {
  const norm = normalizeDomain(rawDomain);
  if (!norm.isValid) return false;

  const current = await getSettings();
  const updatedAllow = current.customAllowlist.filter((d) => d !== norm.domain);

  await updateSettings({
    customAllowlist: updatedAllow,
  });

  return true;
}
