# Synchronization Strategy Documentation

This document outlines the complete synchronization strategy for the Blocks Tracker application ecosystem (Electron, React Native, and Backend). It is designed to ensure data consistency, handle offline-first scenarios gracefully, and resolve conflicts predictably.

## 1. Guiding Principles

1.  **Offline-First**: The application must be fully functional without a network connection. All user actions are performed on the local database first for an optimistic and responsive UI.
2.  **Server is the Single Source of Truth**: While clients operate independently, the backend server is the ultimate arbiter of data conflicts.
3.  **Last Write Wins (LWW)**: The entity (Task, Space, etc.) with the most recent `modifiedAt` timestamp is considered the "winner" in any conflict. This is the cornerstone of our conflict resolution.
4.  **Atomic Operations**: Client-side database changes that require a sync operation must be performed within a single transaction (Outbox Pattern). If creating a task fails, the corresponding pending operation must also fail to be created.

## 2. The Sync Cycle

The sync process is initiated automatically when the app starts (if logged in), when network connectivity is restored, or after a local data-changing operation. It consists of two distinct phases:

1.  **PUSH Phase**: The client sends its local, unsynced changes to the server.
2.  **PULL Phase**: After the PUSH is complete, the client fetches the latest changes from the server.

### PUSH Phase

- The `SyncService` queries the `pending_operations` table for the oldest operation with a `pending` status.
- It processes operations one by one, in the order they were created.
- For each operation, it makes the corresponding API call (e.g., `POST /tasks`, `PUT /spaces/:id`).
- The client then handles the server's response according to the **Client-Side Error Handling Strategy** detailed below.

### PULL Phase

- The client fetches its `last_change_id` from local settings.
- It makes a `GET /changes/sync?last_change_id=<id>` request to the server.
- The server returns a list of all entities that have changed since that ID, along with the new `latestChangeId`.
- The client performs a local `upsertMany` operation for each entity type (Tasks, Spaces, etc.) within a single database transaction.
- If the transaction is successful, the client updates its local `last_change_id` to the new `latestChangeId` from the server.

## 3. Backend Conflict Resolution Strategy

The server is responsible for intelligently resolving conflicts based on the "Last Write Wins" principle. It **must not** blindly reject requests.

### Scenario A: Stale Update (e.g., `PUT /tasks/:id`)

This occurs when a client tries to update an entity with data that is older than what's already on the server.

1.  The handler receives the `PUT` request.
2.  It fetches the existing entity from the database.
3.  It compares the `modifiedAt` timestamp from the incoming request with the `modifiedAt` timestamp of the existing entity.
4.  **If the incoming `modifiedAt` is older**, the server **rejects** the request with an `HTTP 409 Conflict` and a specific JSON body: `{"code": "STALE_DATA"}`.

### Scenario B: Duplicate Creation (e.g., `POST /tasks`)

This scenario is specific to entities with server-side unique constraints. Currently, this only applies to the `tasks` table, which enforces uniqueness on `(repetitive_task_template_id, due_date)`.

Entities like `Space`, `Tag`, and `RepetitiveTaskTemplate` do **not** have unique name constraints on the backend, so creating entities with duplicate names will result in two distinct entities being created. This is an intentional design choice to prioritize user experience and simplify sync logic.

1.  The handler receives the `POST` request and attempts to `INSERT` the new entity.
2.  The database throws a **unique constraint violation** error.
3.  The handler **catches this specific error** and does not immediately fail.
4.  It then **queries for the existing entity** that caused the violation (using the unique keys from the request, like `repetitive_task_template_id` and `dueDate`).
5.  It compares the `modifiedAt` timestamp from the incoming request with the `modifiedAt` of the existing entity.
    - **If the incoming request is NEWER**: The server treats the `POST` as an `UPDATE`. It updates the existing record with the data from the incoming request and responds with `HTTP 200 OK`.
    - **If the incoming request is OLDER or the same**: The server's version is correct. It responds with an `HTTP 409 Conflict` and the JSON body: `{"code": "DUPLICATE_ENTITY"}`.

This "create-or-merge" logic is critical to prevent data loss when a device with newer changes syncs second.

## 4. Client-Side Error Handling Strategy

The client's `SyncService` must intelligently handle API responses during the PUSH phase.

| HTTP Status        | Response Body Code             | Reason                                                                                                   | Client Action                                                                                                                                       |
| :----------------- | :----------------------------- | :------------------------------------------------------------------------------------------------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- |
| `2xx`              | -                              | Success                                                                                                  | Delete the `pending_operation`.                                                                                                                     |
| `409 Conflict`     | `{"code": "DUPLICATE_ENTITY"}` | The server has a canonical, newer version of this entity. The local one is a redundant, older duplicate. | **Delete the local entity** (e.g., `taskRepo.deleteTaskById`) AND **delete the `pending_operation`**.                                               |
| `409 Conflict`     | `{"code": "STALE_DATA"}`       | The client tried to push an update that was older than the server's version.                             | **Delete the `pending_operation`**. The PULL phase will fetch the newer version.                                                                    |
| `404 Not Found`    | -                              | The client tried to update an entity that was already deleted on another device.                         | **Delete the `pending_operation`**. The conflict is resolved.                                                                                       |
| `401 Unauthorized` | -                              | The user's session is invalid and could not be refreshed.                                                | **Transient Error**. Keep the operation in the queue. Record a failed attempt (`recordFailedAttempt`). The user will be prompted to sign in again.  |
| `5xx` Server Error | -                              | The server is temporarily unavailable.                                                                   | **Transient Error**. Keep the operation. Record a failed attempt.                                                                                   |
| Network Error      | -                              | The device is offline.                                                                                   | **Transient Error**. Keep the operation. Record a failed attempt.                                                                                   |
| `400`, `422`, etc. | -                              | The data in the operation's payload is malformed or violates a business rule.                            | **Permanent Failure**. Mark the operation's status as `'failed'`. This unblocks the queue and indicates a client-side bug that needs investigation. |

## 5. Client-Side Data Merging (PULL Phase)

To protect unsynced local changes from being overwritten by stale data from the server, the client's `upsertMany` methods must also follow the "Last Write Wins" principle.

- When processing entities from the PULL phase, the repository method must check if a local version of the entity already exists.
- **If it exists**, the `UPDATE` part of the `UPSERT` should only execute if the incoming record's `modifiedAt` timestamp is strictly greater than the local record's `modifiedAt`.

**Example (SQLite):**

```sql
INSERT INTO tasks (...)
VALUES (...)
ON CONFLICT(id) DO UPDATE SET
  title = excluded.title,
  ...
WHERE excluded.modified_at > tasks.modified_at;
```

This ensures that if a user makes a local change while a PULL is in progress, their change will not be overwritten by the slightly older data just fetched from the server. The user's local change will be correctly pushed in the next sync cycle.
