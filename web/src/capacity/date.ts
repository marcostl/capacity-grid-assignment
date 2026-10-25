const DATE_PATTERN = /^\d{4}-\d{2}-\d{2}$/

export function shiftDate(value: string, days: number): string {
  const date = new Date(`${value}T00:00:00Z`)
  date.setUTCDate(date.getUTCDate() + days)
  return date.toISOString().slice(0, 10)
}

export function validateRange(from: string, to: string): string | null {
  if (!isDate(from) || !isDate(to)) return 'Choose both dates.'
  if (to < from) return 'The end date must be on or after the start date.'
  return null
}

export function formatWeek(start: string): string {
  return new Intl.DateTimeFormat(undefined, {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    timeZone: 'UTC',
  }).format(new Date(`${start}T00:00:00Z`))
}

function isDate(value: string): boolean {
  if (!DATE_PATTERN.test(value)) return false
  const date = new Date(`${value}T00:00:00Z`)
  return !Number.isNaN(date.getTime()) && date.toISOString().slice(0, 10) === value
}
