import { describe, expect, it } from 'vitest'
import { shiftDate, validateRange } from './date'

describe('capacity dates', () => {
  it('shifts a range across month and year boundaries', () => {
    expect(shiftDate('2025-12-29', 7)).toBe('2026-01-05')
    expect(shiftDate('2026-01-05', -7)).toBe('2025-12-29')
  })

  it('rejects incomplete and reversed ranges', () => {
    expect(validateRange('', '2026-01-01')).toBe('Choose both dates.')
    expect(validateRange('2026-02-01', '2026-01-01')).toBe(
      'The end date must be on or after the start date.',
    )
    expect(validateRange('2026-01-01', '2026-01-01')).toBeNull()
  })
})
