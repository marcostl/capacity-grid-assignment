import type { CapacityPage } from './types'

export class CapacityAPIError extends Error {
  constructor(message: string, readonly status: number) {
    super(message)
    this.name = 'CapacityAPIError'
  }
}

type FetchCapacityPageOptions = {
  from: string
  to: string
  cursor?: string
  signal?: AbortSignal
}

export async function fetchCapacityPage({
  from,
  to,
  cursor,
  signal,
}: FetchCapacityPageOptions): Promise<CapacityPage> {
  const params = new URLSearchParams({ from, to, limit: '100' })
  if (cursor) params.set('cursor', cursor)

  const response = await fetch(`/api/capacity?${params}`, { signal })
  if (!response.ok) {
    let message = 'Capacity could not be loaded. Please try again.'
    try {
      const body: unknown = await response.json()
      if (isErrorResponse(body)) message = body.error
    } catch {
      // Keep the safe fallback when an upstream response is not JSON.
    }
    throw new CapacityAPIError(message, response.status)
  }
  return response.json() as Promise<CapacityPage>
}

function isErrorResponse(value: unknown): value is { error: string } {
  return (
    typeof value === 'object' &&
    value !== null &&
    'error' in value &&
    typeof value.error === 'string'
  )
}
