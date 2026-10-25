export function CapacityLoadingGrid() {
  return (
    <div className="loading-grid" aria-label="Loading team capacity" aria-busy="true">
      <div className="loading-header">
        {Array.from({ length: 6 }, (_, index) => (
          <span className="skeleton" key={index} />
        ))}
      </div>
      {Array.from({ length: 8 }, (_, row) => (
        <div className="loading-row" key={row}>
          {Array.from({ length: 6 }, (_, column) => (
            <span className="skeleton" key={column} />
          ))}
        </div>
      ))}
      <span className="sr-only">Loading capacity…</span>
    </div>
  )
}
