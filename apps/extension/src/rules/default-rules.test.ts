import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';
import { describe, expect, it } from 'vitest';

type StaticRule = {
  id: number;
  action: { type: string; redirect?: { extensionPath?: string } };
  condition: { urlFilter?: string; resourceTypes?: string[] };
};

const rules = JSON.parse(
  readFileSync(resolve(process.cwd(), 'public/rules.json'), 'utf8')
) as StaticRule[];

describe('bundled default blocking rules', () => {
  it('stays within Chrome’s guaranteed static rule budget', () => {
    expect(rules).toHaveLength(30_000);
    expect(new Set(rules.map((rule) => rule.id)).size).toBe(rules.length);
  });

  it('prioritizes major adult domains and redirects top-level navigation', () => {
    const filters = new Set(rules.map((rule) => rule.condition.urlFilter));
    for (const domain of ['pornhub.com', 'xvideos.com', 'xnxx.com', 'xhamster.com']) {
      expect(filters.has(`||${domain}^`)).toBe(true);
    }

    const rule = rules.find((candidate) => candidate.condition.urlFilter === '||pornhub.com^');
    expect(rule?.action.type).toBe('redirect');
    expect(rule?.condition.resourceTypes).toEqual(['main_frame', 'sub_frame']);
    expect(rule?.action.redirect?.extensionPath).toContain('/blocked.html?domain=pornhub.com');
  });
});
