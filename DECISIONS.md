# Decisions

Yours to write, not your AI's. Short is good — bullets are fine, and half a page is
plenty. We read this first.

## What did the spec not tell you?

There are things this brief doesn't specify. Which ones did you hit, what did you decide,
and why?

- no definition of what a week or day is.
- no pagination mention whatsoever
- no mention about real scale of the DB. Seeded DB contains 500 people. Is this it?
- I went ahead and implemented pagination straight using person ID as cursor. Not ideal because we might want to do some filtering or sorting based on person but good enough for v0
- no mention about how fresh the data can be in the UI. I'm assuming here the data doesn't change as much but should we build some kind of real time stream that updates values? Out of the scope here.

## What did you notice that looked wrong?

Anything in the output that didn't match what you expected. Whether you fixed it or left
it, we want to know you saw it.

[ Not sure I understood the question correctly but here's my answer :') ]

- There seem to be duplicated records in the DB (would be nice to add a notifications section in the dashboard with these for managers to inspect)
- Some values look like ratios instead of total hours (0.125 or 0.5)

## What did the AI get wrong that you caught?

One concrete example. Every real session has one.

- AI assumed weekends should not count. I asked it to actually count them as there are a lot of people working on weekedns (eg: retail), so this wass not a valid assumption.
- No protection about stale writes when patching capcity. I wasn't convinced about suggested approaches (using previous capacity as a guardrail or using PostgresSQL xmin value) so I decided this had to be fixed at the schema level.
- Asked it to use hexagonal architecture to improve code maintainability and testability.
- The proposed previous / next week navigation was completely off so I got rid of it.

## What would you do differently with a week?

- I'd clarify the requirements a bit better and understand the final expected state
- The instructions look as if they were meant for a quick MVP (which I understand given this is just a technical assignment) but for a bigger project there are a lot of unknowns
- UX would be a huge thing for this dashboard. First step would be to understand how managers do use this Trying to understand how managers use this dashboard.
