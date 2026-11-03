import { getSettings, updateSettings } from '../storage/storage';
import { normalizeDomain } from '../domain/normalize';

interface BlocklistResponse {
  version: string;
  count: number;
  domains: string[];
}

const API_URL = import.meta.env.VITE_VAJRA_API_URL?.replace(/\/+$/, '');
const MAX_REMOTE_DOMAINS = 5000;
const MAX_RESPONSE_CHARS = 512 * 1024;

/** Refreshes only the small project-curated delta; the bundled rules still work offline. */
export async function refreshRemoteBlocklist(): Promise<void> {
  if (!API_URL) return;

  const settings = await getSettings();
  const headers = new Headers();
  if (settings.remoteBlocklistETag) {
    headers.set('If-None-Match', settings.remoteBlocklistETag);
  }
  const response = await fetch(`${API_URL}/v1/blocklist`, {
    method: 'GET',
    headers,
    cache: 'no-store',
    signal: AbortSignal.timeout(10_000),
  });
  if (response.status === 304) return;
  if (!response.ok) throw new Error(`Vajra API returned HTTP ${response.status}`);

  const responseText = await response.text();
  if (responseText.length > MAX_RESPONSE_CHARS) {
    throw new Error('Vajra API blocklist exceeded the response size limit');
  }
  let payload: Partial<BlocklistResponse>;
  try {
    payload = JSON.parse(responseText) as Partial<BlocklistResponse>;
  } catch {
    throw new Error('Vajra API returned invalid JSON');
  }
  if (
    !payload ||
    typeof payload.version !== 'string' ||
    payload.version.length === 0 ||
    !Array.isArray(payload.domains) ||
    payload.domains.length > MAX_REMOTE_DOMAINS ||
    payload.count !== payload.domains.length
  ) {
    throw new Error('Vajra API returned an invalid blocklist');
  }

  const normalizedDomains: string[] = [];
  for (const raw of payload.domains) {
    if (typeof raw !== 'string') {
      throw new Error('Vajra API returned a non-string domain');
    }
    const normalized = normalizeDomain(raw);
    if (!normalized.isValid || normalized.domain !== raw) {
      throw new Error('Vajra API returned an invalid or non-normalized domain');
    }
    normalizedDomains.push(normalized.domain);
  }
  const domains = [...new Set(normalizedDomains)];
  if (domains.length !== normalizedDomains.length) {
    throw new Error('Vajra API returned duplicate domains');
  }

  await updateSettings({
    remoteBlocklist: domains,
    remoteBlocklistVersion: payload.version,
    remoteBlocklistETag: response.headers.get('ETag') || undefined,
  });
}
