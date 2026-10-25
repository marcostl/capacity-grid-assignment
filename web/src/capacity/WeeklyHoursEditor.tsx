import { useState, type FormEvent } from 'react'
import { useUpdateWeeklyHours } from './useUpdateWeeklyHours'
import type { PersonCapacity } from './types'

type Props = {
  person: PersonCapacity
  onClose: () => void
}

export function WeeklyHoursEditor({ person, onClose }: Props) {
  const [value, setValue] = useState(String(person.weeklyHours))
  const mutation = useUpdateWeeklyHours(onClose)
  const weeklyHours = Number(value)
  const isValid = value.trim() !== '' && Number.isFinite(weeklyHours) && weeklyHours >= 0 && weeklyHours <= 168
  const isUnchanged = weeklyHours === person.weeklyHours

  function save(event: FormEvent) {
    event.preventDefault()
    if (!isValid || isUnchanged || mutation.isPending) return
    mutation.mutate({ personId: person.id, weeklyHours })
  }

  return (
    <form className="hours-editor" onSubmit={save}>
      <label htmlFor={`weekly-hours-${person.id}`}>Weekly hours</label>
      <div>
        <input
          id={`weekly-hours-${person.id}`}
          type="number"
          min="0"
          max="168"
          step="any"
          value={value}
          onChange={(event) => setValue(event.target.value)}
          disabled={mutation.isPending}
          autoFocus
          aria-invalid={!isValid}
        />
        <button type="submit" disabled={!isValid || isUnchanged || mutation.isPending}>
          {mutation.isPending ? 'Saving…' : 'Save'}
        </button>
        <button type="button" onClick={onClose} disabled={mutation.isPending} aria-label="Cancel editing">
          ×
        </button>
      </div>
      {mutation.isError && <span className="save-error" role="alert">{mutation.error.message}</span>}
    </form>
  )
}
