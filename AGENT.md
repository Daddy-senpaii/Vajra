AGENT.md — Vayra
Product: Vayra — a self-improvement platform for challenges, personal progress, analytics, and friends/accountability.
Purpose of this file: This is the primary implementation brief for Codex/AI coding agents. Follow it before making architectural or product decisions.

1. Product vision
Vayra helps people choose or create self-improvement challenges, follow their rules, record daily progress, understand their consistency, and encourage friends.

Product principle: Improve yourself privately, or grow together with your friends.

This is not only a checklist app. The core experience combines:
1. Challenge discovery and enrollment.
2. Challenge-specific daily/weekly rules.
3. Personal progress tracking and analytics.
4. Friends, profiles, and social activity.
5. A meaningful “progress since your last visit” comparison.
Build a reliable, maintainable MVP first. Do not add unrelated features or over-engineer the first release.
2. Required technology stack
Use this stack unless the repository already contains a deliberate, working choice that conflicts with it. If a conflict exists, explain it before changing it.
Mobile application
- React Native with TypeScript for the Android mobile application.
- React function components and hooks.
- React Navigation for navigation.
- A consistent component/theme system; avoid styling each screen independently.
- Use Android Studio for the Android SDK, emulator, device debugging, Gradle/build tooling, and native project inspection.
- The app must run on an Android emulator and a physical Android device.
- Do not build a web-only React application as a substitute for the requested Android app.
Backend
- Go for the HTTP API.
- Use a clear layered structure: handlers/transport → services/use-cases → repositories/data access.
- Use context.Context for request-scoped work.
- Validate inputs at the API boundary and enforce authorization in the service layer.
- Keep business rules out of HTTP handlers.
Database and infrastructure
- Use MongoDB as the database for the current development stage.
- Use the official MongoDB Go Driver (go.mongodb.org/mongo-driver) from the Go backend.
- Keep collection/document definitions and index creation scripts or startup index setup checked into the repository. Do not assume a relational migration system.
- Store timestamps in UTC; convert to the user's local timezone for display and date-based challenge rules.
- Store passwords only as strong password hashes (never plaintext).
- Keep secrets and environment-specific configuration out of source control.
- Provide .env.example with variable names and safe placeholders.
- Use MongoDB document modeling deliberately: embed small, bounded data that is normally read/updated together; reference independently growing entities such as users, friendships, activity events, and progress history.
- Define and create indexes for unique usernames/emails, challenge lookups, progress lookups, friend relationships/requests, activity feeds, notifications, and profile-view records.
- Do not require multi-document transactions for ordinary workflows. Prefer single-document atomic updates and idempotent operations. If a workflow truly needs a multi-document transaction, document that MongoDB transactions require a replica-set-capable deployment and ensure local development is configured accordingly.
- Use MongoDB schema validation where practical, while still validating all incoming data in Go.
Important implementation assumptions
The requested stack says React, Android Studio, and Go. Because the product is an Android mobile app, interpret “React” as React Native + TypeScript. Android Studio is the Android development toolchain, not the UI framework. Do not introduce a second frontend framework.

MongoDB replaces the earlier PostgreSQL choice for now. Do not add PostgreSQL, SQL migrations, GORM, or relational-only assumptions.
Go-specific limitations and engineering rules
Go is the backend language only. It does not replace React Native and should not be used to build the mobile UI.
- Go has no exceptions or class inheritance. Return and handle errors explicitly; use composition and small interfaces rather than trying to imitate class hierarchies.
- Keep interfaces small and define them near the code that consumes them. Avoid creating an interface for every struct or function.
- Use context.Context for request cancellation, deadlines, and MongoDB operations. Pass it through handlers, services, and repository calls; do not store request contexts in long-lived structs.
- Do not start goroutines for simple synchronous work. Every goroutine needs a clear owner, shutdown/cancellation path, and error-handling strategy. Avoid unbounded goroutines and shared mutable state.
- Use the official MongoDB Go Driver APIs and BSON types correctly. Handle mongo.ErrNoDocuments, duplicate-key errors, connection failures, and context timeouts explicitly.
- Reuse a configured MongoDB client; do not create a new client for every request. Configure connection pooling and close the client during graceful shutdown.
- Go's static type system does not automatically validate JSON or BSON business rules. Validate request payloads, enum values, IDs, date ranges, and numeric targets explicitly.
- Avoid heavy reflection or generic abstractions where straightforward typed Go code is clearer.
- Keep handlers thin; put challenge calculations, authorization, and progress-delta logic in testable service functions.
- Do not claim a feature is implemented merely because an endpoint compiles. Add tests for business behavior and MongoDB persistence.
Suggested libraries
Choose stable, compatible versions and document exact versions in the project:
- React Native + TypeScript
- React Navigation
- TanStack Query for server state (or one equivalent; do not duplicate server-state libraries)
- React Hook Form + Zod for forms/validation (optional)
- Go standard net/http or a small router such as Chi
- MongoDB driver and a migration tool
- Structured logging
- Go tests and React Native component tests
Do not add libraries just because they are popular. Keep the dependency list small.
3. Repository and project structure
If starting from an empty repository, use a monorepo layout similar to:
vayra/
├── AGENT.md
├── README.md
├── .gitignore
├── .env.example
├── mobile/
│   ├── android/
│   ├── src/
│   │   ├── app/
│   │   ├── navigation/
│   │   ├── screens/
│   │   ├── components/
│   │   ├── features/
│   │   ├── services/
│   │   ├── hooks/
│   │   ├── theme/
│   │   ├── types/
│   │   └── utils/
│   └── package.json
├── backend/
│   ├── cmd/api/
│   ├── internal/
│   │   ├── config/
│   │   ├── database/       # MongoDB client, collections, indexes
│   │   ├── http/
│   │   ├── middleware/
│   │   ├── auth/
│   │   ├── users/
│   │   ├── challenges/
│   │   ├── progress/
│   │   ├── analytics/
│   │   ├── friends/
│   │   ├── activity/
│   │   ├── notifications/
│   │   └── modes/          # Running, reading, workout, digital wellbeing
│   ├── scripts/            # MongoDB index/bootstrap scripts if needed
│   ├── go.mod
│   └── README.md
└── docs/
    ├── PRODUCT.md
    ├── API.md
    └── DATA_MODEL.md
Adapt the structure to the actual repository if it already exists. Do not move or delete existing files without a clear reason.
4. Application navigation and screen map
Primary navigation
Use a bottom tab bar with:
App
├── Home
├── Friends
├── Analytics
└── Profile
Other screens should be pushed/nested from these tabs.
Screen inventory
Vayra
├── 0. App bootstrap
│   ├── Loading / splash
│   ├── Check stored authentication session
│   └── Route to Login or Home
│
├── 1. Authentication
│   ├── Login
│   │   ├── Email/username
│   │   ├── Password
│   │   ├── Show/hide password
│   │   └── Forgot password (can be deferred if email delivery is not configured)
│   ├── Sign up
│   │   ├── Display name
│   │   ├── Username
│   │   ├── Email
│   │   ├── Password
│   │   └── Terms/privacy acknowledgement
│   └── Logout
│
├── 2. Home dashboard
│   ├── Greeting and notification entry point
│   ├── Current challenge summary
│   ├── Today's task completion summary
│   ├── Quick actions
│   │   ├── Browse challenges
│   │   ├── Open analytics
│   │   ├── Open friends
│   │   └── Open Modes
│   ├── Modes section (visible on Home)
│   │   ├── Running Mode
│   │   ├── Reading Mode
│   │   ├── Workout Logger
│   │   └── Digital Wellbeing
│   ├── Continue challenge
│   └── Recent activity / encouragement preview
│
├── 3. Challenges
│   ├── Browse challenges
│   │   ├── Featured
│   │   ├── All challenges
│   │   ├── Search/filter
│   │   └── Challenge categories
│   ├── Challenge details
│   │   ├── Name and description
│   │   ├── Duration
│   │   ├── Rules and requirements
│   │   ├── Difficulty / category (if defined)
│   │   └── Join challenge
│   ├── My challenges
│   │   ├── Active
│   │   ├── Completed
│   │   └── Abandoned (if supported)
│   └── Custom challenge builder
│       ├── Challenge name and description
│       ├── Duration or no fixed end date
│       ├── Add/edit/remove rules
│       ├── Rule frequency: daily, weekly, or one-time
│       ├── Rule target and unit (where applicable)
│       └── Save as private custom challenge
│
├── 4. Challenge execution / daily progress
│   ├── Challenge overview
│   ├── Current day and total duration
│   ├── Overall completion
│   ├── Today's checklist
│   ├── Numeric input for measurable rules (e.g. distance)
│   ├── Optional short note
│   ├── Save progress
│   └── History/calendar view
│
├── 5. Modes (Home feature group)
│   ├── Modes overview
│   │   ├── Four mode cards on Home
│   │   ├── Quick access to start/log an activity
│   │   └── Recent mode activity summary
│   ├── Running Mode
│   │   ├── Start / pause / resume / finish a run
│   │   ├── Manual run entry for MVP
│   │   ├── Duration, distance, and optional notes
│   │   ├── Pace calculation when distance and duration exist
│   │   ├── Run history and personal totals
│   │   └── GPS route tracking is a later enhancement; do not imply it works until native permissions and background behavior are implemented and tested
│   ├── Reading Mode
│   │   ├── Add a book or reading item
│   │   ├── Record pages read or reading duration
│   │   ├── Set an optional daily reading target
│   │   ├── Reading session history
│   │   └── Current reading streak / progress
│   ├── Workout Logger
│   │   ├── Log a workout session
│   │   ├── Workout type and date
│   │   ├── Exercises
│   │   │   ├── Exercise name
│   │   │   ├── Sets
│   │   │   ├── Reps
│   │   │   └── Weight/load (optional)
│   │   ├── Duration and optional notes
│   │   └── Workout history
│   └── Digital Wellbeing
│       ├── Manual daily screen-time entry for MVP
│       ├── User-defined daily screen-time goal
│       ├── Track distracting app/site usage manually
│       ├── Daily/weekly usage summaries
│       ├── Focus session timer (optional MVP enhancement)
│       └── Android UsageStats / app blocking integrations are later native Android work; require permissions and must not be presented as available until implemented and tested
│
├── 6. Personal analytics
│   ├── Overall progress
│   ├── Per-challenge progress
│   ├── Rule/habit completion rates
│   ├── Current and best streaks
│   ├── Daily/weekly trends
│   ├── Completed and missed requirements
│   └── Date range selection
│
├── 7. Friends
│   ├── Friends list
│   │   ├── Search friends
│   │   ├── Online/recently active indicator (optional; do not fake presence)
│   │   └── Open friend profile
│   ├── Find users
│   ├── Friend requests
│   │   ├── Incoming
│   │   ├── Outgoing
│   │   ├── Accept
│   │   ├── Decline
│   │   └── Cancel
│   ├── Friend profile
│   │   ├── Avatar, display name, username, bio
│   │   ├── Active challenges visible to the viewer
│   │   ├── Challenge progress
│   │   ├── Selected public statistics
│   │   ├── Recent activity
│   │   └── Remove friend / block / report entry points
│   └── Friendship privacy and visibility rules
│
├── 8. Progress since last visit (key differentiator)
│   ├── Track the last time the current user viewed each friend's profile
│   ├── Compare the friend's progress events after that timestamp
│   ├── Show newly completed tasks
│   ├── Show challenge day/progress changes
│   ├── Show streak changes when available
│   ├── Show newly joined/completed challenges
│   ├── Show a clear “Since your last visit” summary
│   └── Update last-viewed timestamp after successfully loading the profile
│
├── 9. Friends activity feed
│   ├── Friend joined a challenge
│   ├── Friend completed a rule/check-in
│   ├── Friend completed a challenge
│   ├── Friend reached a milestone
│   ├── Activity timestamp
│   ├── Reactions / encouragement
│   └── Feed visibility follows privacy settings
│
├── 10. Notifications
│   ├── Friend request received
│   ├── Friend request accepted
│   ├── Friend milestone / challenge completion (if enabled)
│   └── In-app notification list
│
└── 11. Profile and settings
    ├── View/edit own profile
    ├── Avatar and bio
    ├── Privacy and visibility
    │   ├── Profile visibility
    │   ├── Challenge visibility
    │   ├── Progress/statistics visibility
    │   └── Activity-feed sharing
    ├── Notification preferences
    ├── Account settings
    ├── Logout
    └── About / help
5. Feature requirements
5.1 Authentication and user profile
- Users can create an account, log in, and log out.
- Usernames must be unique and normalized consistently.
- Never expose password hashes or authentication secrets through API responses.
- Protect authenticated endpoints with middleware.
- Persist sessions securely on Android using an appropriate secure-storage mechanism.
- A user can edit their display name, avatar, and bio.
- Add sensible loading, empty, validation, and error states.
5.2 Challenge catalogue
Seed the first catalogue with these challenge templates:
Challenges
├── Winter Arc
├── 90-Day Challenge
├── 75 Hard
├── Monk Mode
└── Dopamine Detox
These are challenge templates, not hardcoded application logic. Store challenge descriptions, durations, rules, and rule frequencies as data.

For challenge templates:
- Show an accurate description and duration.
- Display all requirements before enrollment.
- Make clear that users are joining a template and starting their own enrollment.
- Do not imply that a user has completed a requirement without a saved progress record.
- Do not force a restart for missed days unless that specific challenge definition explicitly requires it.
- Avoid presenting unsafe or medically inappropriate requirements as universal advice. Allow challenge rules to be customized and include sensible safety language where needed.
5.3 Challenge rules and tracking
A rule must support:
- Title and optional description.
- Frequency: daily, weekly, or once.
- Input type: checkbox or numeric.
- Optional target and unit, such as 2 km.
- Whether the rule is required.
- Display order.
Examples:
- Run 2 km daily.
- Complete a morning workout.
- Wake up by 5:00 AM.
- Keep social media below a chosen time limit.
- Read 10 pages daily.
Do not treat every rule as a checkbox. Numeric rules need a value field and validation.

For each user enrollment:
- Store start date, optional end date, status, and timezone.
- Derive the challenge day from the enrollment start date and user's timezone.
- Save progress against a specific rule and local challenge date.
- Prevent duplicate progress records for the same enrollment/rule/date.
- Allow a user to correct their own entry, with updated timestamps.
- Show missed/incomplete requirements honestly; never silently mark them complete.
5.4 Analytics
Analytics must be derived from saved progress records, not hardcoded sample numbers.

Minimum metrics:
- Completion rate = completed eligible rule instances / total eligible rule instances.
- Current streak and best streak, with a clearly documented streak definition.
- Total completed rule instances.
- Challenge progress and status.
- Daily/weekly completion trend.
For numeric requirements, define completion based on the rule's target and saved value. Handle days with no eligible requirement separately from missed requirements.
5.5 Modes on the Home screen
Modes are a separate quick-action feature group on the Home screen. They complement challenges; they do not replace the challenge system. A user may log a run, reading session, workout, or wellbeing entry without joining a challenge. Where appropriate, a user may optionally link a mode entry to an active challenge rule, but avoid double-counting the same activity.
Running Mode
- MVP: allow manual entry of distance, duration, date/time, and optional notes.
- Calculate pace only when distance and duration are valid.
- Validate non-negative distance and positive duration.
- Show run history and simple totals.
- GPS tracking, route maps, background tracking, and live pace are out of the initial MVP unless explicitly requested later. These require Android native location permissions, lifecycle handling, battery consideration, and device testing.
Reading Mode
- Allow users to create reading items (for example, book title and optional author).
- Log pages read and/or reading duration per session.
- Support an optional daily target.
- Show reading history and progress toward the target.
- Do not assume page counts or reading duration when the user has not entered them.
Workout Logger
- Allow users to create workout logs with date, workout type, duration, and optional notes.
- Support exercises with sets, reps, and optional weight/load.
- Allow a simple workout log without requiring every exercise field.
- Do not provide medical or injury-prevention claims.
Digital Wellbeing
- MVP is a manual tracker: users enter daily screen-time totals and optionally categorize distracting apps/sites or activities.
- Allow a user-defined daily screen-time goal and show daily/weekly summaries.
- A focus-session timer may be included if it fits the initial scope.
- Do not claim to read Android-wide screen time, monitor other apps, block websites, or enforce restrictions in the MVP.
- Reading Android app usage requires native Android UsageStatsManager access and explicit user-granted special access; app/site blocking requires separate platform-specific design and permissions. Treat these as later Android-native integrations, not ordinary React Native features.
5.6 Friends and social graph
- A user can send a friend request to another user.
- A recipient can accept or decline an incoming request.
- A user can cancel an outgoing pending request.
- A friendship exists only after acceptance.
- Prevent self-friending and duplicate pending/accepted relationships.
- Support removing a friend.
- Support blocking a user; a blocked relationship must prevent new requests and hide restricted social data.
- Friend search must not reveal private account details.
- Do not expose email addresses or other private account fields to friends.
5.7 Friend profiles and visibility
A friend profile can display only fields allowed by the profile owner's privacy settings.

Suggested privacy controls:
- Profile: public, friends, or private.
- Active challenges: friends, public, or private.
- Progress/statistics: friends, public, or private.
- Activity feed: enabled/disabled, with visibility restricted by relationship and profile settings.
Defaults should be conservative: share with accepted friends, not with everyone. Users must be able to hide individual challenges or their statistics.
5.8 Progress since last visit
This is a core feature, not a cosmetic label.

Expected behavior:
1. The current user opens a friend's profile.
2. The backend checks whether the viewer is allowed to see that profile.
3. The backend reads the last-viewed timestamp for this viewer/friend pair.
4. It returns visible progress/activity events created after that timestamp.
5. The UI summarizes meaningful changes, such as:
   - +3 tasks completed
   - Challenge advanced from Day 12 to Day 13
   - Streak increased from 5 to 6 days
   - Started a new challenge
6. After the profile data is successfully returned/displayed, update the viewer/friend last-viewed timestamp.
7. On the next visit, use the new timestamp as the comparison point.
Implementation notes:
- Use a friendship-specific view record, not one global last_seen_at on the friend.
- Suggested unique key: (viewer_user_id, profile_user_id).
- Do not count events the viewer is not permitted to see.
- Avoid updating the timestamp before the comparison query, or the “since last visit” result will be empty.
- First visit should show a friendly state such as “You're all caught up” or “Start following their progress,” not a fabricated historical delta.
- Use server-generated timestamps and stable event IDs.
- Define how deleted/edited progress affects historical summaries. Prefer event records or an auditable progress history if accurate deltas are required.
5.9 Activity feed and reactions
Activity events may include:
- Challenge started.
- Daily/weekly rule completed.
- Challenge milestone reached.
- Challenge completed.
The feed must:
- Include only activities the viewer is authorized to see.
- Be paginated.
- Have stable ordering by event time and ID.
- Support a small set of reactions, such as 👏, 🔥, and 💪.
- Prevent duplicate identical reactions from the same user on the same activity.
- Avoid public leaderboards or competitive ranking in the MVP unless explicitly requested.
5.10 Notifications
Start with in-app notifications. Push notifications can be added later.
- Friend request received.
- Friend request accepted.
- Friend milestone/challenge completion, if the recipient has enabled this category.
- Notifications must respect privacy and notification preferences.
- Mark notifications as read.
- Do not send notifications for private activities.
6. MongoDB data model — initial proposal
Use MongoDB collections and BSON documents. This is a document database: do not translate the old relational schema mechanically into SQL-style tables or assume SQL joins.

Suggested collections:
users
├── _id (ObjectId)
├── username (normalized, unique index)
├── email (normalized, unique index)
├── password_hash
├── display_name
├── avatar_url (nullable)
├── bio (nullable)
├── timezone
├── privacy
│   ├── profile_visibility
│   ├── challenges_visibility
│   ├── stats_visibility
│   └── activity_sharing_enabled
├── notification_preferences
├── created_at
└── updated_at

challenge_templates
├── _id (ObjectId)
├── name
├── slug (unique index)
├── description
├── duration_days (nullable)
├── category
├── rules[]
│   ├── rule_id
│   ├── title
│   ├── description
│   ├── frequency
│   ├── input_type
│   ├── target_value (nullable)
│   ├── unit (nullable)
│   ├── is_required
│   └── sort_order
├── is_published
├── created_at
└── updated_at

user_challenges
├── _id (ObjectId)
├── user_id (ObjectId reference)
├── template_id (nullable ObjectId reference)
├── name
├── description_snapshot (optional)
├── rules[]                 # snapshot of rules at enrollment/custom rules
│   ├── rule_id
│   ├── source_rule_id (nullable)
│   ├── title
│   ├── description
│   ├── frequency
│   ├── input_type
│   ├── target_value
│   ├── unit
│   ├── is_required
│   └── sort_order
├── start_date
├── end_date (nullable)
├── timezone
├── status
├── created_at
└── updated_at

progress_entries
├── _id (ObjectId)
├── user_id (ObjectId)
├── user_challenge_id (ObjectId)
├── rule_id
├── local_date (YYYY-MM-DD in the enrollment timezone)
├── numeric_value (nullable)
├── is_completed
├── note (private by default; never copy into social event payloads)
├── created_at
└── updated_at
    unique compound index: user_challenge_id + rule_id + local_date

friend_requests
├── _id (ObjectId)
├── sender_user_id (ObjectId)
├── receiver_user_id (ObjectId)
├── status
├── created_at
└── updated_at

friendships
├── _id (ObjectId)
├── user_low_id (ObjectId; canonical lower/sorted ID)
├── user_high_id (ObjectId; canonical higher/sorted ID)
└── created_at
    unique compound index: user_low_id + user_high_id

blocks
├── _id (ObjectId)
├── blocker_user_id (ObjectId)
├── blocked_user_id (ObjectId)
└── created_at
    unique compound index: blocker_user_id + blocked_user_id

friend_profile_views
├── _id (ObjectId)
├── viewer_user_id (ObjectId)
├── profile_user_id (ObjectId)
└── last_viewed_at
    unique compound index: viewer_user_id + profile_user_id

activity_events
├── _id (ObjectId)
├── actor_user_id (ObjectId)
├── event_type
├── entity_type
├── entity_id
├── payload (small, privacy-safe snapshot only)
└── created_at

activity_reactions
├── _id (ObjectId)
├── activity_event_id (ObjectId)
├── user_id (ObjectId)
├── reaction_type
└── created_at
    unique compound index: activity_event_id + user_id + reaction_type

notifications
├── _id (ObjectId)
├── recipient_user_id (ObjectId)
├── actor_user_id (nullable ObjectId)
├── type
├── entity_type (nullable)
├── entity_id (nullable)
├── is_read
└── created_at

mode_entries
├── _id (ObjectId)
├── user_id (ObjectId)
├── mode_type (running/reading/workout/digital_wellbeing)
├── occurred_at
├── local_date
├── data (typed per mode; validate by mode_type)
├── note (private by default)
├── linked_user_challenge_id (nullable ObjectId)
└── created_at / updated_at
MongoDB modeling and correctness rules:
- Embed bounded challenge-template rules inside challenge_templates; they are read with the template and are small.
- Snapshot the rules into each user_challenges document at enrollment. Existing enrollments must not change when a template is edited later.
- Keep progress entries in their own collection because they grow over time and need independent date-range queries.
- Keep friendships, friend requests, activity events, reactions, notifications, profile views, and mode entries in separate collections because they grow independently and have different query patterns.
- Use ObjectId consistently for MongoDB document identifiers and references in the backend. Convert IDs to validated strings at the HTTP boundary; never accept arbitrary malformed IDs.
- MongoDB does not provide relational foreign-key enforcement. Validate referenced documents in services and handle missing/deleted references.
- Create unique and compound indexes during application bootstrap or through a repeatable index setup command. Index creation must be idempotent.
- For progress, enforce uniqueness with a compound index. Use an atomic upsert where appropriate; handle duplicate-key races.
- For friendships, canonicalize the pair of IDs before insert so A-B and B-A cannot create duplicate friendships.
- Use projections to avoid fetching private fields such as password_hash when reading user profiles.
- Avoid unbounded arrays in documents. Do not embed all progress, all activity, all notifications, or all friend IDs inside a user document.
- MongoDB documents have a 16 MB size limit; keep documents bounded and paginate growing collections.
- Use atomic single-document updates where possible. Multi-document transactions require a replica set (including local replica-set setup) and should not be assumed in a standalone local MongoDB instance.
- Design multi-step operations to be idempotent and recoverable. For example, accepting a friend request should tolerate retries and avoid duplicate friendship records.
- Use schema validation at the database level where useful, plus explicit Go request/domain validation.
- Define deletion behavior for referenced documents; do not leave the API returning private or orphaned social data.
Modes data payload examples
Keep mode_entries.data type-specific and validated in Go. Suggested shapes:
running:
  distance_meters
  duration_seconds
  pace_seconds_per_km (derived; optional stored snapshot)
  entry_method (manual initially; gps later)

reading:
  title
  author (optional)
  pages_read (optional)
  duration_seconds (optional)
  daily_target_pages (optional)

workout:
  workout_type
  duration_seconds (optional)
  exercises[]:
    name
    sets
    reps
    weight_value (optional)
    weight_unit (optional)

digital_wellbeing:
  total_screen_time_seconds
  daily_goal_seconds (optional)
  categories[]:
    label
    duration_seconds (optional)
Do not accept arbitrary unvalidated JSON for these payloads. Decode into a typed Go request structure based on mode_type, validate it, then persist the normalized document.
7. Backend API outline
Use versioned routes under /api/v1. Return consistent JSON errors and appropriate HTTP status codes.
Auth
├── POST   /auth/register
├── POST   /auth/login
├── POST   /auth/refresh        (if refresh-token flow is implemented)
└── POST   /auth/logout

Me / Profile
├── GET    /me
├── PATCH  /me
├── GET    /users/:username
└── PATCH  /me/privacy

Challenges
├── GET    /challenges
├── GET    /challenges/:id
├── POST   /challenges/custom
├── POST   /challenges/:id/join
├── GET    /me/challenges
├── GET    /me/challenges/:id
├── PATCH  /me/challenges/:id/status
└── DELETE /me/challenges/:id       (only if deletion semantics are defined)

Progress
├── GET    /me/challenges/:id/progress
├── PUT    /me/challenges/:id/rules/:ruleId/progress/:date
├── GET    /me/progress/summary
└── GET    /me/analytics

Friends
├── GET    /users/search?q=
├── GET    /me/friends
├── GET    /me/friend-requests/incoming
├── GET    /me/friend-requests/outgoing
├── POST   /users/:userId/friend-requests
├── POST   /friend-requests/:requestId/accept
├── POST   /friend-requests/:requestId/decline
├── DELETE /friend-requests/:requestId
├── DELETE /me/friends/:userId
└── POST   /users/:userId/block

Friend profile and progress comparison
├── GET    /users/:userId/profile
├── GET    /users/:userId/progress-since-last-visit
└── POST   /users/:userId/profile-viewed

Modes
├── GET    /me/modes/summary
├── GET    /me/modes/:modeType/entries
├── POST   /me/modes/:modeType/entries
├── GET    /me/modes/:modeType/entries/:entryId
├── PATCH  /me/modes/:modeType/entries/:entryId
└── DELETE /me/modes/:modeType/entries/:entryId

Activity and notifications
├── GET    /me/feed
├── POST   /activities/:activityId/reactions
├── DELETE /activities/:activityId/reactions/:reactionType
├── GET    /me/notifications
└── PATCH  /me/notifications/:notificationId/read
API rules:
- Never trust a user ID supplied by the client as the authenticated identity. Derive the acting user from the authenticated session/token.
- Validate UUIDs, dates, pagination, enums, numeric targets, and request bodies.
- Use pagination for friends, search, feed, and notifications.
- Return only fields the caller is allowed to see.
- Prefer atomic MongoDB document updates; use a multi-document transaction only when required and only when the deployment supports replica sets.
- Add request IDs and structured logs.
- Document request/response examples in docs/API.md.
- Do not return raw database errors to clients.
8. Mobile UI/UX direction
Visual direction
- Modern, calm, motivating, and uncluttered.
- Dark theme as the initial design direction, with accessible contrast.
- Use restrained purple/lavender accents and green for positive progress; avoid excessive gradients and visual noise.
- Cards should have consistent spacing, radius, typography, and states.
- Use clear labels and readable text; do not rely on color alone to communicate status.
- Use reusable components rather than duplicating screen layouts.
Reusable component tree
components/
├── common/
│   ├── AppButton
│   ├── AppInput
│   ├── Avatar
│   ├── Badge
│   ├── EmptyState
│   ├── ErrorState
│   ├── LoadingState
│   ├── ProgressBar
│   └── ScreenContainer
├── challenge/
│   ├── ChallengeCard
│   ├── ChallengeRuleRow
│   ├── ChallengeProgressHeader
│   ├── DailyChecklist
│   └── NumericRuleInput
├── analytics/
│   ├── MetricCard
│   ├── CompletionChart
│   └── StreakCard
├── modes/
│   ├── ModeCard
│   ├── RunningEntryForm
│   ├── ReadingSessionForm
│   ├── WorkoutLogForm
│   └── DigitalWellbeingForm
└── social/
    ├── FriendCard
    ├── FriendRequestCard
    ├── ActivityCard
    ├── ProgressDeltaCard
    └── ReactionBar
Required screen states
Every network-backed screen must account for:
- Initial loading.
- Refreshing where appropriate.
- Empty state.
- Error state with retry.
- Success state.
- Offline/network failure where practical.
- Permission/visibility-denied state.
Do not fill the UI with fake permanent sample data. Mock data may be used during early UI development, but clearly isolate it and replace it with API integration.
9. Security, privacy, and reliability
- Enforce authentication and authorization on the backend.
- Hash passwords using a suitable password-hashing algorithm; never encrypt or store plaintext passwords.
- Use TLS in deployed environments.
- Keep access tokens and secrets out of logs.
- Use Android secure storage for tokens.
- Validate and limit request sizes.
- Add rate limits to login, registration, user search, and friend-request endpoints.
- Prevent IDOR: a user must not access another user's private challenge/progress by changing a URL ID.
- Enforce blocks and privacy settings in every profile, feed, analytics, and progress-delta query.
- Keep notes private by default; never include them in activity payloads.
- Do not expose email, password hashes, private notes, or private challenge details in social APIs.
- Avoid fabricated online status. If presence is not implemented, omit it.
- Provide a way to delete an account and associated data before public release; document retention behavior.
- Do not make health, fitness, or sleep claims beyond what the recorded data supports.
10. Development phases and implementation order
Work in small, testable increments. Do not attempt to implement every feature in one giant change.
Phase 0 — Repository inspection and plan
- Inspect the current repository before editing.
- Identify existing framework, files, scripts, and conventions.
- Report any conflicts with this document.
- Create or update the project structure without deleting existing work.
- Add a concise README and .env.example.
Phase 1 — Foundation
- Set up React Native TypeScript Android app.
- Confirm Android Studio emulator/device build works.
- Set up Go API server.
- Set up MongoDB and repeatable collection/index initialization.
- Add health-check endpoint.
- Add configuration and structured logging.
- Add basic tests and run instructions.
Acceptance: App launches on Android; backend starts; MongoDB connects; required indexes initialize idempotently; health endpoint responds.
Phase 2 — Authentication and profile
- Register, login, logout.
- Secure token/session handling.
- Get and edit own profile.
- Basic navigation and protected screens.
Acceptance: A user can register, log in, close/reopen the app, and remain authenticated safely.
Phase 3 — Challenges and progress
- Seed the five initial challenge templates.
- Browse and view challenge details.
- Join a challenge.
- Snapshot rules into the user's enrollment.
- Display today's rules.
- Save checkbox and numeric progress.
- Show challenge history and basic summary.
Acceptance: A user can join a challenge and record/retrieve progress without duplicate daily entries.
Phase 4 — Home Modes
- Add the four mode cards to the Home screen.
- Implement manual running entries and history.
- Implement reading items/sessions and targets.
- Implement workout logs with exercises, sets, reps, and optional load.
- Implement manual digital wellbeing entries and goals.
- Add typed Go request models, MongoDB persistence, and validation for each mode.
- Keep GPS tracking and Android-wide usage monitoring out of scope for this phase.
Acceptance: Users can create, edit, view, and delete their own mode entries; summaries are calculated from saved entries; users cannot read or modify another user's entries.
Phase 5 — Analytics
- Implement completion calculations.
- Implement streak rules and document their definitions.
- Add challenge and overall analytics.
- Add basic charts/trends.
Acceptance: Analytics match independently calculated results from test data.
Phase 6 — Friends
- User search.
- Friend request lifecycle.
- Friends list.
- Friend profile.
- Privacy settings and authorization.
Acceptance: Two test users can send/accept a request and view only the data allowed by privacy settings.
Phase 7 — Progress since last visit
- Add friend_profile_views.
- Query visible events after the previous visit.
- Calculate progress deltas.
- Update last-viewed timestamp after the comparison has been read.
- Handle first visit and no-new-progress states.
Acceptance: Reopening a friend's profile shows only changes since the previous successful visit; a subsequent visit does not repeat old changes.
Phase 8 — Activity feed and notifications
- Activity event creation for supported actions.
- Paginated friends feed.
- Reactions.
- In-app notifications and read state.
Acceptance: A friend sees permitted activity, can react once per reaction type, and can mark notifications read.
Phase 9 — Hardening and release preparation
- Loading, empty, error, and offline states.
- Authorization and privacy tests.
- Android device testing.
- Performance checks and database indexes.
- Account deletion and privacy review.
- Build/release documentation.
11. Testing requirements
Backend
- Unit tests for challenge day calculations, rule eligibility, completion rate, streaks, and progress deltas.
- Integration tests for MongoDB repositories, indexes, and persistence behavior.
- API tests for authentication, validation, pagination, and authorization.
- Privacy tests proving users cannot access private friend data.
- Friendship tests for duplicate requests, self-requests, blocks, accept/decline, and removal.
- Progress-since-last-visit tests:
  - first visit;
  - no new activity;
  - one or multiple progress events;
  - hidden/private events excluded;
  - timestamp updated in the correct order;
  - two different viewers have independent last-viewed timestamps.
Mobile
- Component tests for forms, rule inputs, progress cards, and empty/error states.
- Navigation tests for protected routes.
- Test on at least one Android emulator and, before release, a physical Android device.
- Verify screen layouts at common Android screen sizes and font scaling.
Definition of done for each feature
A feature is not complete merely because its screen exists. It must include:
1. UI and navigation.
2. Backend endpoint and business logic, when applicable.
3. Database migration/model, when applicable.
4. Loading, empty, error, and success states.
5. Input validation and authorization.
6. Tests.
7. Updated documentation.
12. Agent working instructions
Codex/AI agent must follow these rules:
1. Inspect before coding. Read the repository and existing instructions first.
2. Do not overwrite working code blindly. Preserve existing conventions and functionality.
3. Plan first for large tasks. State the files/modules to be changed and the implementation sequence.
4. Work phase by phase. Start with Phase 0 and Phase 1 unless the user explicitly requests a later phase.
5. Do not silently change the stack. React Native + TypeScript, Android toolchain via Android Studio, Go backend, and MongoDB are the intended choices.
6. Do not invent product requirements. If a requirement is unclear and materially affects data or behavior, ask before committing to it.
7. Keep business logic testable. Put calculations in pure functions/services where possible.
8. Keep frontend and backend contracts aligned. Update API documentation and types when endpoints change.
9. Never trust client-side privacy checks. Backend authorization is mandatory.
10. No fake success. Show success only after the server confirms the operation.
11. No hardcoded analytics. Calculate metrics from persisted records.
12. No secrets in source control. Use environment variables and .env.example.
13. Run checks after changes. Run relevant formatters, linters, tests, MongoDB index/bootstrap checks, and Android build checks. Report exactly what passed and what could not be run.
14. Avoid scope creep. Do not add chat, public leaderboards, payments, AI coaching, or other unrequested features in the MVP.
15. Keep the user informed. At the end of each task, summarize files changed, behavior implemented, commands/tests run, known limitations, and the next logical step.
First task for the agent
Start with Phase 0 — Repository inspection and plan.

Then implement only Phase 1 — Foundation, unless the user explicitly asks for more. Do not begin the Modes feature until the user asks for the next phase or approves it. Before implementation, report:
- What already exists in the repository.
- Whether React Native/Android and Go backend are already initialized.
- Any blockers or stack conflicts.
- The proposed first set of files to create/change.
Do not implement all phases at once.
