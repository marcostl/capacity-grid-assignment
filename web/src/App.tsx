import { CapacityGrid } from './CapacityGrid'

const FROM = '2025-12-29'
const TO = '2026-01-16'

export function App() {
  return (
    <main>
      <header className="page-header">
        <p className="eyebrow">Team overview</p>
        <h1>Team capacity</h1>
        <p>Review weekly allocation and spot over-committed team members.</p>
      </header>
      <CapacityGrid from={FROM} to={TO} />
    </main>
  )
}
