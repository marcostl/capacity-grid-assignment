import { useRef } from 'react'
import { useVirtualizer } from '@tanstack/react-virtual'
import { formatWeek } from './date'
import { CapacityRow, gridColumns, gridWidth } from './CapacityRow'
import { useIntersectionObserver } from './useIntersectionObserver'
import type { CapacityWeek, PersonCapacity } from './types'

type Props = {
  weeks: CapacityWeek[]
  people: PersonCapacity[]
  hasNextPage: boolean
  isFetchingNextPage: boolean
  nextPageError: boolean
  onLoadMore: () => void
}

export function CapacityTable({
  weeks,
  people,
  hasNextPage,
  isFetchingNextPage,
  nextPageError,
  onLoadMore,
}: Props) {
  const scrollElement = useRef<HTMLDivElement>(null)
  const sentinel = useRef<HTMLDivElement>(null)
  const virtualizer = useVirtualizer({
    count: people.length,
    getScrollElement: () => scrollElement.current,
    estimateSize: () => 64,
    overscan: 8,
  })
  const virtualRows = virtualizer.getVirtualItems()
  const totalHeight = virtualizer.getTotalSize()

  useIntersectionObserver({
    root: scrollElement,
    target: sentinel,
    enabled: hasNextPage && !isFetchingNextPage && !nextPageError,
    onIntersect: onLoadMore,
  })

  const columns = gridColumns(weeks.length)
  const width = gridWidth(weeks.length)

  return (
    <div className="capacity-scroll" ref={scrollElement}>
      <div className="capacity-table" role="table" aria-rowcount={people.length} style={{ width }}>
        <div className="capacity-header" role="row" style={{ gridTemplateColumns: columns }}>
          <div className="person-heading" role="columnheader">
            Team member
          </div>
          {weeks.map((week) => (
            <div role="columnheader" key={week.start} title={`${week.start} to ${week.end}`}>
              <span>Week of</span>
              <strong>{formatWeek(week.start)}</strong>
            </div>
          ))}
        </div>

        <div className="capacity-body" style={{ height: totalHeight + 64 }}>
          {virtualRows.map((virtualRow) => {
            const person = people[virtualRow.index]
            if (!person) return null
            return (
              <div
                className="virtual-row"
                key={person.id}
                ref={virtualizer.measureElement}
                data-index={virtualRow.index}
                style={{ transform: `translateY(${virtualRow.start}px)` }}
              >
                <CapacityRow person={person} weeks={weeks} />
              </div>
            )
          })}

          <div className="load-more" ref={sentinel} style={{ top: totalHeight }} role="status">
            {isFetchingNextPage && <><span className="spinner" /> Loading more people…</>}
            {nextPageError && (
              <>
                <span>More people could not be loaded.</span>
                <button type="button" onClick={onLoadMore}>Try again</button>
              </>
            )}
            {!hasNextPage && !nextPageError && <span>All {people.length} people loaded</span>}
          </div>
        </div>
      </div>
    </div>
  )
}
