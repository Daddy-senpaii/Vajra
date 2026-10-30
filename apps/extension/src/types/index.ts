export type ProtectionStatus = 'ENABLED' | 'DISABLED';

export type Decision = 'ALLOW' | 'BLOCK' | 'INVALID';

export interface UserSettings {
  enabled: boolean;
  blockedCount: number;
  customBlocklist: string[];
  customAllowlist: string[];
  lastBlockedUrl?: string;
  lastBlockedDomain?: string;
  lastBlockedTimestamp?: number;
}

export interface EvaluationResult {
  decision: Decision;
  normalizedDomain: string;
  reason: string;
}

export interface DNRRule {
  id: number;
  priority: number;
  action: {
    type: 'redirect' | 'allow' | 'block';
    redirect?: {
      extensionPath?: string;
      url?: string;
    };
  };
  condition: {
    urlFilter?: string;
    requestDomains?: string[];
    resourceTypes: chrome.declarativeNetRequest.ResourceType[];
  };
}
