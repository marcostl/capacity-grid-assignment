export type CapacityWeek = {
  start: string
  end: string
}

export type PersonCapacity = {
  id: number
  name: string
  weeklyHours: number
  allocatedHours: number[]
}

export type CapacityPage = {
  range: {
    from: string
    to: string
  }
  weeks: CapacityWeek[]
  people: PersonCapacity[]
  pagination: {
    limit: number
    nextCursor?: string
  }
}
