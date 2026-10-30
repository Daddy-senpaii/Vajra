/**
 * Normalizes and validates a domain or URL string in parity with the Go core engine.
 */
export interface NormalizeResult {
  domain: string;
  isValid: boolean;
  error?: string;
}

export function normalizeDomain(raw: string): NormalizeResult {
  let s = raw.trim();
  if (!s) {
    return { domain: '', isValid: false, error: 'domain is empty' };
  }

  // If a scheme is present, extract host using browser URL parsing
  if (s.includes('://')) {
    try {
      const parsed = new URL(s);
      s = parsed.hostname;
    } catch {
      return { domain: '', isValid: false, error: 'invalid URL format' };
    }
  } else if (s.includes('/')) {
    // Contains path without scheme (e.g. "example.com/path")
    s = s.split('/')[0];
  }

  // Strip port if present (e.g. "example.com:8080")
  if (s.includes(':')) {
    const parts = s.split(':');
    if (parts.length === 2 && !s.includes('[')) {
      s = parts[0];
    }
  }

  // Strip trailing FQDN dot (e.g. "example.com.")
  if (s.endsWith('.')) {
    s = s.slice(0, -1);
  }

  s = s.trim().toLowerCase();
  if (!s) {
    return { domain: '', isValid: false, error: 'domain is empty' };
  }

  // RFC Validation
  if (s.length > 253) {
    return { domain: '', isValid: false, error: 'domain exceeds maximum length of 253 characters' };
  }

  const labels = s.split('.');
  if (labels.length < 2) {
    return { domain: '', isValid: false, error: 'domain must contain at least two labels separated by a dot' };
  }

  for (let i = 0; i < labels.length; i++) {
    const label = labels[i];
    if (label.length === 0) {
      return { domain: '', isValid: false, error: 'consecutive dots not permitted' };
    }
    if (label.length > 63) {
      return { domain: '', isValid: false, error: 'domain label exceeds maximum length of 63 characters' };
    }
    if (label.startsWith('-') || label.endsWith('-')) {
      return { domain: '', isValid: false, error: 'label cannot start or end with a hyphen' };
    }

    // Only alphanumeric and hyphen allowed
    if (!/^[a-z0-9-]+$/.test(label)) {
      return { domain: '', isValid: false, error: 'domain contains invalid characters' };
    }

    // TLD checks
    if (i === labels.length - 1) {
      if (label.length < 2) {
        return { domain: '', isValid: false, error: 'top-level domain must be at least 2 characters' };
      }
      if (/^\d+$/.test(label)) {
        return { domain: '', isValid: false, error: 'top-level domain cannot be all digits' };
      }
    }
  }

  return { domain: s, isValid: true };
}
