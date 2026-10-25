import { useState, type FormEvent } from 'react'
import { validateRange } from './date'

type Props = {
  from: string
  to: string
  onApply: (from: string, to: string) => void
  disabled?: boolean
}

export function CapacityToolbar({ from, to, onApply, disabled }: Props) {
  const [draftFrom, setDraftFrom] = useState(from)
  const [draftTo, setDraftTo] = useState(to)
  const [error, setError] = useState<string | null>(null)

  function apply(event: FormEvent) {
    event.preventDefault()
    const validationError = validateRange(draftFrom, draftTo)
    if (validationError) {
      setError(validationError)
      return
    }
    setError(null)
    onApply(draftFrom, draftTo)
  }

  return (
    <div className="capacity-toolbar">
      <form className="range-form" onSubmit={apply}>
        <label>
          From
          <input
            type="date"
            value={draftFrom}
            onChange={(event) => setDraftFrom(event.target.value)}
          />
        </label>
        <label>
          To
          <input
            type="date"
            value={draftTo}
            onChange={(event) => setDraftTo(event.target.value)}
          />
        </label>
        <button type="submit" disabled={disabled || (draftFrom === from && draftTo === to)}>
          Apply range
        </button>
      </form>
      {error && <p className="field-error">{error}</p>}
    </div>
  )
}
