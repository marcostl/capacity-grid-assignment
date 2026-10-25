import { describe, expect, it } from 'vitest'
import { validateRange } from './date'

describe('capacity dates', () => {
  it('rejects incomplete and reversed ranges', () => {
    expect(validateRange('', '2026-01-01')).toBe('Choose both dates.')
    expect(validateRange('2026-02-01', '2026-01-01')).toBe(
      'The end date must be on or after the start date.',
    )
    expect(validateRange('2026-01-01', '2026-01-01')).toBeNull()
  })
})
