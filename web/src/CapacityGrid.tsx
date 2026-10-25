import { useMemo, useState } from 'react'
import { CapacityLoadingGrid } from './capacity/CapacityLoadingGrid'
import { CapacityTable } from './capacity/CapacityTable'
import { CapacityToolbar } from './capacity/CapacityToolbar'
import { useCapacity } from './capacity/useCapacity'

type Props = {
  from: string
  to: string
}

export function CapacityGrid({ from: initialFrom, to: initialTo }: Props) {
  const [range, setRange] = useState({ from: initialFrom, to: initialTo })
  const capacity = useCapacity(range.from, range.to)
  const people = useMemo(
    () => capacity.data?.pages.flatMap((page) => page.people) ?? [],
    [capacity.data],
  )
  const firstPage = capacity.data?.pages[0]

  return (
    <section className="capacity-section" aria-labelledby="capacity-heading">
      <h2 className="sr-only" id="capacity-heading">Capacity by person and week</h2>
      <CapacityToolbar
        from={range.from}
        to={range.to}
        onApply={(from, to) => setRange({ from, to })}
      />

      {firstPage && (
        <div className="grid-summary">
          <p>
            <strong>{firstPage.range.from}</strong> to <strong>{firstPage.range.to}</strong>
          </p>
          <p aria-live="polite">
            {people.length} people loaded
            {capacity.isFetchingNextPage ? ' — loading more' : ''}
          </p>
        </div>
      )}

      {capacity.isPending && <CapacityLoadingGrid />}

      {capacity.isError && !capacity.data && (
        <div className="error-panel" role="alert">
          <h2>Capacity could not be loaded</h2>
          <p>{capacity.error.message}</p>
          <button type="button" onClick={() => void capacity.refetch()}>
            Try again
          </button>
        </div>
      )}

      {firstPage && people.length === 0 && !capacity.hasNextPage && (
        <div className="empty-state">
          <h2>No team members found</h2>
          <p>There is no capacity to show for this team.</p>
        </div>
      )}

      {firstPage && people.length > 0 && (
        <CapacityTable
          weeks={firstPage.weeks}
          people={people}
          hasNextPage={capacity.hasNextPage}
          isFetchingNextPage={capacity.isFetchingNextPage}
          nextPageError={capacity.isFetchNextPageError}
          onLoadMore={() => void capacity.fetchNextPage()}
        />
      )}
    </section>
  )
}
