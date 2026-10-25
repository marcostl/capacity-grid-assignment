import type { CapacityWeek, PersonCapacity } from './types'

const PERSON_COLUMN_WIDTH = 250
const WEEK_COLUMN_WIDTH = 132

export function gridWidth(weekCount: number) {
  return PERSON_COLUMN_WIDTH + weekCount * WEEK_COLUMN_WIDTH
}

export function gridColumns(weekCount: number) {
  return `${PERSON_COLUMN_WIDTH}px repeat(${weekCount}, ${WEEK_COLUMN_WIDTH}px)`
}

type Props = {
  person: PersonCapacity
  weeks: CapacityWeek[]
}

export function CapacityRow({ person, weeks }: Props) {
  return (
    <div className="capacity-row" role="row" style={{ gridTemplateColumns: gridColumns(weeks.length) }}>
      <div className="person-cell" role="rowheader">
        <strong title={person.name}>{person.name}</strong>
        <span>{formatHours(person.weeklyHours)} weekly hours</span>
      </div>
      {weeks.map((week, index) => {
        const allocated = person.allocatedHours[index] ?? 0
        const overAllocated = allocated > person.weeklyHours
        return (
          <div
            className={`allocation-cell${overAllocated ? ' over-allocated' : ''}`}
            role="cell"
            key={week.start}
            title={`${person.name}: ${formatHours(allocated)} allocated of ${formatHours(person.weeklyHours)} hours`}
          >
            <strong>{formatHours(allocated)}</strong>
            <span>of {formatHours(person.weeklyHours)}h</span>
            {overAllocated && (
              <span className="over-label">
                +{formatHours(allocated - person.weeklyHours)}h
              </span>
            )}
          </div>
        )
      })}
    </div>
  )
}

function formatHours(value: number): string {
  return Number.isInteger(value) ? String(value) : value.toFixed(1)
}
