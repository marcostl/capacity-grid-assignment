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
