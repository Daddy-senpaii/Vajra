import { describe, it, expect } from 'vitest';
import {
  buildBlockRule,
  buildAllowRule,
  generateDynamicRules,
  PRIORITY_ALLOW,
  PRIORITY_BLOCK,
  BASE_BLOCK_ID,
  BASE_ALLOW_ID,
  MAX_DYNAMIC_BLOCK_RULES,
} from './rule-builder';

describe('rule-builder', () => {
  it('builds block rule with redirect action and priority 1', () => {
    const rule = buildBlockRule(1001, 'adult-example.com');
    expect(rule).not.toBeNull();
    expect(rule?.id).toBe(1001);
    expect(rule?.priority).toBe(PRIORITY_BLOCK);
    expect(rule?.action.type).toBe('redirect');
    expect(rule?.action.redirect?.extensionPath).toContain('domain=adult-example.com');
    expect(rule?.condition.urlFilter).toBe('||adult-example.com^');
  });

  it('builds allow rule with allow action and priority 2', () => {
    const rule = buildAllowRule(50001, 'safe.adult-example.com');
    expect(rule).not.toBeNull();
    expect(rule?.id).toBe(50001);
    expect(rule?.priority).toBe(PRIORITY_ALLOW);
    expect(rule?.action.type).toBe('allow');
    expect(rule?.condition.urlFilter).toBe('||safe.adult-example.com^');
  });

  it('rejects invalid domains for rules', () => {
    expect(buildBlockRule(1002, '-invalid..domain')).toBeNull();
    expect(buildAllowRule(50002, '')).toBeNull();
  });

  it('generates dynamic rules list with higher allow priority', () => {
    const rules = generateDynamicRules(['bad1.com', 'bad2.com'], ['good1.com']);
    expect(rules).toHaveLength(3);

    const blockRule = rules.find((r) => r.id === BASE_BLOCK_ID);
    const allowRule = rules.find((r) => r.id === BASE_ALLOW_ID);

    expect(blockRule?.priority).toBe(PRIORITY_BLOCK);
    expect(allowRule?.priority).toBe(PRIORITY_ALLOW);
    expect(allowRule!.priority).toBeGreaterThan(blockRule!.priority!);
  });

  it('keeps combined dynamic rules within Chrome quota and preserves allow rules', () => {
    const blocked = Array.from({ length: MAX_DYNAMIC_BLOCK_RULES }, (_, i) => `blocked-${i}.example`);
    const allowed = Array.from({ length: 10 }, (_, i) => `allowed-${i}.example`);
    const rules = generateDynamicRules(blocked, allowed);

    expect(rules).toHaveLength(MAX_DYNAMIC_BLOCK_RULES);
    expect(rules.filter((rule) => rule.priority === PRIORITY_ALLOW)).toHaveLength(10);
    expect(rules.filter((rule) => rule.priority === PRIORITY_BLOCK)).toHaveLength(MAX_DYNAMIC_BLOCK_RULES - 10);
  });
});
