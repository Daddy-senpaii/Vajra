import { normalizeDomain } from '../domain/normalize';

export const PRIORITY_ALLOW = 2;
export const PRIORITY_BLOCK = 1;

export const BASE_BLOCK_ID = 1000;
export const BASE_ALLOW_ID = 50000;
export const MAX_DYNAMIC_BLOCK_RULES = 5000;

export function buildBlockRule(
  id: number,
  domain: string,
  blockedPagePath: string = '/blocked.html'
): chrome.declarativeNetRequest.Rule | null {
  const norm = normalizeDomain(domain);
  if (!norm.isValid) return null;

  return {
    id,
    priority: PRIORITY_BLOCK,
    action: {
      type: 'redirect' as chrome.declarativeNetRequest.RuleActionType.REDIRECT,
      redirect: {
        extensionPath: `${blockedPagePath}?domain=${encodeURIComponent(norm.domain)}`,
      },
    },
    condition: {
      urlFilter: `||${norm.domain}^`,
      resourceTypes: [
        'main_frame' as chrome.declarativeNetRequest.ResourceType.MAIN_FRAME,
        'sub_frame' as chrome.declarativeNetRequest.ResourceType.SUB_FRAME,
      ],
    },
  };
}

export function buildAllowRule(
  id: number,
  domain: string
): chrome.declarativeNetRequest.Rule | null {
  const norm = normalizeDomain(domain);
  if (!norm.isValid) return null;

  return {
    id,
    priority: PRIORITY_ALLOW,
    action: {
      type: 'allow' as chrome.declarativeNetRequest.RuleActionType.ALLOW,
    },
    condition: {
      urlFilter: `||${norm.domain}^`,
      resourceTypes: [
        'main_frame' as chrome.declarativeNetRequest.ResourceType.MAIN_FRAME,
        'sub_frame' as chrome.declarativeNetRequest.ResourceType.SUB_FRAME,
        'xmlhttprequest' as chrome.declarativeNetRequest.ResourceType.XMLHTTPREQUEST,
        'image' as chrome.declarativeNetRequest.ResourceType.IMAGE,
      ],
    },
  };
}

export function generateDynamicRules(
  blockedDomains: string[],
  allowedDomains: string[],
  blockedPagePath: string = '/blocked.html'
): chrome.declarativeNetRequest.Rule[] {
  const rules: chrome.declarativeNetRequest.Rule[] = [];

  // Chrome's dynamic-rule quota applies to the combined block and allow rules.
  // Keep allow rules first so a large personal allowlist still takes precedence.
  const normalizedAllows = [...new Set(allowedDomains
    .map((domain) => normalizeDomain(domain))
    .filter((result) => result.isValid)
    .map((result) => result.domain))];
  const allowEntries = normalizedAllows.slice(0, MAX_DYNAMIC_BLOCK_RULES);
  const remainingBlockSlots = MAX_DYNAMIC_BLOCK_RULES - allowEntries.length;

  let blockId = BASE_BLOCK_ID;
  for (const d of [...new Set(blockedDomains)].slice(0, remainingBlockSlots)) {
    const rule = buildBlockRule(blockId++, d, blockedPagePath);
    if (rule) {
      rules.push(rule);
    }
  }

  let allowId = BASE_ALLOW_ID;
  for (const domain of allowEntries) {
    const rule = buildAllowRule(allowId++, domain);
    if (rule) {
      rules.push(rule);
    }
  }

  return rules;
}
