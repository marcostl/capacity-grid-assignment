# Worklog

Running notes on how this got built — decisions, assumptions, dead ends, and anything
left unfinished. Append as you go; a line or two per entry is right.

---

[ This notes has been written by both the LLM and myself ]

## Implement capacity endpoint

- Capacity uses Monday–Sunday weeks and counts all calendar days, including weekends.
- The endpoint is split into an HTTP adapter, application service, and PostgreSQL repository. Going with an hexagonal architecture here to separate concerns and aid in testing.
- ID-based cursor pagination (100 default, 250 max) bounds work for larger rosters. Not defined in the instructions but I can't see a reason why not to add it. Specially if we want the API to be performant and assumin it will run on DBs bigger than the seeded one from this exercise.
- Internal errors are logged and translated to a safe UI message.
- Lacking better logging infrastructure on purpose as I believe this is out of the scope for now. Making a note to revisit it later if time permits
- Lacking better testing too on purpose. Might revisit it later if time permits

## Implement capacity grid

- Using TanStack Query for cursor-based infinite loading
- Using TanStack Virtual to keep the mounted row count bounded.
- An Intersection Observer preloads the next page near the end of the virtual content.
- Date changes cancel obsolete requests through React Query's abort signal.
- Initial, incremental, empty, and failed loading states have separate UI so already-loaded rows remain usable when a later page fails.
- Only rows are virtualized; week columns remain horizontally scrollable because the API already caps the range at roughly two years.
- Capacity editing is deferred until the PATCH endpoint is implemented.
- Date selection currently uses explicit from/to inputs and an Apply button.
- Week-step navigation is intentionally deferred while its UX is reconsidered.
- CSS has been completely neglected. This is intended for now. Will revisit if time permits
- Same for the UX of the date selection. Very clunky for now and I believe can be greatly improved. Maybe for later.

## Implement weekly-hours updates

- PATCH returns the authoritative person and an informational weak ETag derived from the representation;
- Capacity rows expose the same token.
- The ETag deliberately does not enforce concurrency yet.
- Preventing stale writes correctly needs a schema-backed version updated and checked atomically, which the fixed schema does not provide. Not implemented as the instructions explicitely mention not to update the schema.
- Weekly hours accept decimals from 0 through 168 (24h x 7days, potential max number of hours a person can work, ignoring basic human rights and I guess 100% of labor laws)
- Malformed input, missing people, and internal failures are translated to distinct user-safe responses.
