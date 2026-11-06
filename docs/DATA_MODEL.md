# Data model

MongoDB is the development database. The initial foundation checks in repeatable index definitions for the collections anticipated by `AGENT.md`: users, challenge templates, enrollments, progress entries, friend requests, friendships, activity events, notifications, and friend profile views.

The first application feature should define and document collection fields alongside its repository/service implementation. Small bounded data read or updated together may be embedded; independently growing records such as progress history, friendships, and activity should be separate documents. Store timestamps in UTC, normalize unique usernames/emails consistently, and enforce unique progress records per enrollment, rule, and local challenge date. Mode collections have per-user date indexes; digital wellbeing entries are unique per user and local date.

Indexes are created idempotently by `backend/internal/database.EnsureIndexes` at API startup. No SQL database or migration system is used.
