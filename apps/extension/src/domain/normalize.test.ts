import { describe, it, expect } from 'vitest';
import { normalizeDomain } from './normalize';

describe('normalizeDomain', () => {
  it('normalizes valid domains correctly', () => {
    expect(normalizeDomain('  Example.COM  ')).toEqual({
      domain: 'example.com',
      isValid: true,
    });

    expect(normalizeDomain('https://WWW.Adult-Site.COM/videos?id=1#tag')).toEqual({
      domain: 'www.adult-site.com',
      isValid: true,
    });

    expect(normalizeDomain('example.com:8080')).toEqual({
      domain: 'example.com',
      isValid: true,
    });

    expect(normalizeDomain('example.com.')).toEqual({
      domain: 'example.com',
      isValid: true,
    });

    expect(normalizeDomain('sub.deep.example.com')).toEqual({
      domain: 'sub.deep.example.com',
      isValid: true,
    });
  });

  it('rejects invalid domains', () => {
    expect(normalizeDomain('').isValid).toBe(false);
    expect(normalizeDomain('   ').isValid).toBe(false);
    expect(normalizeDomain('invalid').isValid).toBe(false);
    expect(normalizeDomain('invalid domain').isValid).toBe(false);
    expect(normalizeDomain('-example.com').isValid).toBe(false);
    expect(normalizeDomain('example-.com').isValid).toBe(false);
    expect(normalizeDomain('example..com').isValid).toBe(false);
    expect(normalizeDomain('exam@ple.com').isValid).toBe(false);
    expect(normalizeDomain('example.123').isValid).toBe(false);
    expect(normalizeDomain('example.c').isValid).toBe(false);
  });
});
