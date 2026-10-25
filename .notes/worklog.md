# Worklog

Running notes on how this got built — decisions, assumptions, dead ends, and anything
left unfinished. Append as you go; a line or two per entry is right.

---

[ These notes have been written by both the LLM and me ]

## Implement capacity endpoint

- Capacity uses Monday–Sunday weeks and counts all calendar days, including weekends.
- The endpoint is split into an HTTP adapter, application service, and PostgreSQL repository. Going with a hexagonal architecture here to separate concerns and aid in testing.
- ID-based cursor pagination (100 default, 250 max) bounds work for larger rosters. Not defined in the instructions, but I can't see a reason not to add it, especially if we want the API to be performant and assume it will run on databases bigger than the seeded one from this exercise.
- Internal errors are logged and translated to a safe UI message.
- Lacking better logging infrastructure on purpose, as I believe this is out of scope for now. Making a note to revisit it later if time permits.
- Lacking better testing on purpose, too. I might revisit it later if time permits.

## Implement capacity grid

- Using TanStack Query for cursor-based infinite loading.
- Using TanStack Virtual to keep the mounted row count bounded.
- An Intersection Observer preloads the next page near the end of the virtual content.
- Date changes cancel obsolete requests through React Query's abort signal.
- Initial, incremental, empty, and failed loading states have separate UI so already-loaded rows remain usable when a later page fails.
- Only rows are virtualized; week columns remain horizontally scrollable because the API already caps the range at roughly two years.
- Capacity editing is deferred until the PATCH endpoint is implemented.
- Date selection currently uses explicit from/to inputs and an Apply button.
- Week-step navigation is intentionally deferred while its UX is reconsidered.
- CSS has been completely neglected. This is intended for now. I will revisit it if time permits.
- Same for the UX of the date selection. Very clunky for now and I believe can be greatly improved. Maybe for later.

## Implement weekly-hours updates API

- PATCH returns the authoritative person and an informational weak ETag derived from the representation.
- Capacity rows expose the same token.
- The ETag deliberately does not enforce concurrency yet.
- Preventing stale writes correctly needs a schema-backed version updated and checked atomically, which the fixed schema does not provide. Not implemented, as the instructions explicitly say not to update the schema.
- Weekly hours accept decimals from 0 through 168 (24 hours × 7 days, the potential maximum number of hours a person can work, ignoring basic human rights and, I guess, 100% of labor laws).
- Malformed input, missing people, and internal failures are translated to distinct user-safe responses.

## Implement weekly-hours updates in the UI

- The grid edits one person at a time, disables controls during a save, and keeps the attempted value visible on failure with a row-level error.
- Successful saves patch every loaded capacity page immediately, then invalidate capacity queries for server reconciliation; failed saves leave cached values untouched.
- Editing another person's capacity while one is being saved is disabled.

## Previous / next controls in the table toolbar

- Added simple previous/next controls that shift the entire applied range by seven days.

## Final notes

I chose to shift and refetch the complete selected range for previous/next navigation because it was the simplest reliable way to satisfy the requirement within the available time.

With more time, I would normalize frontend state by person and week:

`Map<PersonID, { person, allocations: Map<WeekStart, AllocatedHours> }>`

That would allow adjacent weeks to be fetched and merged into the existing grid without refetching the entire visible range. It would also require coordinating pagination across two dimensions: people vertically and weeks horizontally.
