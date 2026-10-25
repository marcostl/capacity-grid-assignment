import { describe, expect, it } from 'vitest'
import type { InfiniteData } from '@tanstack/react-query'
import { replacePersonInCapacityData } from './useUpdateWeeklyHours'
import type { CapacityPage } from './types'

const page = (id: number, weeklyHours: number): CapacityPage => ({
  range: { from: '2026-01-05', to: '2026-01-11' },
  weeks: [{ start: '2026-01-05', end: '2026-01-11' }],
  people: [{ id, name: `Person ${id}`, weeklyHours, allocatedHours: [40], etag: 'old' }],
  pagination: { limit: 1 },
})

describe('replacePersonInCapacityData', () => {
  it('updates the person in every loaded page without changing allocations', () => {
    const data: InfiniteData<CapacityPage> = {
      pages: [page(1, 40), page(2, 40)],
      pageParams: [undefined, '1'],
    }

    const result = replacePersonInCapacityData(data, {
      id: 2,
      name: 'Person 2',
      weeklyHours: 24,
      etag: 'new',
    })

    expect(result?.pages[0].people[0].weeklyHours).toBe(40)
    expect(result?.pages[1].people[0]).toMatchObject({ weeklyHours: 24, allocatedHours: [40], etag: 'new' })
  })
})
