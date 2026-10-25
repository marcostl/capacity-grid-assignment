import type { CapacityPage } from './types'

export class CapacityAPIError extends Error {
  constructor(message: string, readonly status: number) {
    super(message)
    this.name = 'CapacityAPIError'
  }
}

export type UpdatedPerson = {
  id: number
  name: string
  weeklyHours: number
  etag: string
}

type FetchCapacityPageOptions = {
  from: string
  to: string
  cursor?: string
  signal?: AbortSignal
}

export async function updateWeeklyHours(
  personId: number,
  weeklyHours: number,
): Promise<UpdatedPerson> {
  await developmentDelay()
  const response = await fetch(`/api/people/${personId}`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ weeklyHours }),
  })
  if (!response.ok) throw await responseError(response, 'Weekly hours could not be updated. Please try again.')
  return response.json() as Promise<UpdatedPerson>
}

export async function fetchCapacityPage({
  from,
  to,
  cursor,
  signal,
}: FetchCapacityPageOptions): Promise<CapacityPage> {
  await developmentDelay(signal)
  const params = new URLSearchParams({ from, to, limit: '100' })
  if (cursor) params.set('cursor', cursor)

  const response = await fetch(`/api/capacity?${params}`, { signal })
  if (!response.ok) throw await responseError(response, 'Capacity could not be loaded. Please try again.')
  return response.json() as Promise<CapacityPage>
}

async function responseError(response: Response, fallback: string): Promise<CapacityAPIError> {
  let message = fallback
  try {
    const body: unknown = await response.json()
    if (isErrorResponse(body)) message = body.error
  } catch {
    // Keep the safe fallback when an upstream response is not JSON.
  }
  return new CapacityAPIError(message, response.status)
}

function developmentDelay(signal?: AbortSignal): Promise<void> {
  if (!import.meta.env.DEV) return Promise.resolve()
  if (signal?.aborted) return Promise.reject(signal.reason)

  return new Promise((resolve, reject) => {
    const onAbort = () => {
      window.clearTimeout(timeout)
      reject(signal?.reason ?? new DOMException('Request aborted', 'AbortError'))
    }
    const timeout = window.setTimeout(() => {
      signal?.removeEventListener('abort', onAbort)
      resolve()
    }, 1_000)
    signal?.addEventListener('abort', onAbort, { once: true })
  })
}

function isErrorResponse(value: unknown): value is { error: string } {
  return (
    typeof value === 'object' &&
    value !== null &&
    'error' in value &&
    typeof value.error === 'string'
  )
}
