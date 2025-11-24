# 1. Consistent Timestamp Formatting with Milliseconds

- **Status**: Accepted
- **Date**: 2025-11-21

## Context and Problem Statement

The client applications (React Native, Electron) serialize dates using JavaScript's `Date.toISOString()`, which produces timestamps with mandatory millisecond precision (e.g., `2025-11-20T18:00:00.000Z`).

The Go backend, using the default `time.Time` JSON marshaling, omits zero-value milliseconds (e.g., `2025-11-20T18:00:00Z`).

This discrepancy caused a critical bug where the `UNIQUE` constraint on the `tasks` table (`repetitive_task_template_id`, `due_date`) would not trigger correctly during synchronization. The client would attempt to create a task that already existed on the server, but because the `dueDate` string was different, the server did not recognize it as a duplicate, leading to data integrity issues.

## Decision Drivers

- The need for a single, canonical timestamp format across the entire stack (clients and backend).
- The need to fix the duplicate entity sync bug.
- The desire to centralize formatting logic on the backend, which is the source of truth.

## Considered Options

1.  **Modify client-side serialization**: Change both React Native and Electron apps to format timestamps without milliseconds if they are zero. This would require changes in multiple codebases and is prone to future error.
2.  **Do nothing**: Leave the bug in place. Unacceptable.
3.  **Introduce a custom time type in Go**: Create a wrapper around `time.Time` that overrides the default JSON marshaling to always include millisecond precision.

## Decision

We will adopt **Option 3**. A new type, `models.JSONTime`, will be created. This type will implement the `json.Marshaler` interface to always format timestamps using the `2006-01-02T15:04:05.000Z` layout. All `time.Time` fields in API-facing models (`Task`, `Space`, `Tag`, etc.) will be converted to use `JSONTime`.

### Consequences

- **Positive**: Guarantees timestamp format consistency. Fixes the sync bug permanently. Centralizes timestamp formatting logic in one place (`models/jsontime.go`).
- **Negative**: Adds a small amount of boilerplate code. Developers must remember to convert `JSONTime` to `time.Time` (e.g., `time.Time(myJSONTime)`) when using methods like `.Before()` or `.After()`. This is documented in the code.
