---
inclusion: auto
description: Security rule - UUID for external, auto-increment ID for internal use only
---

# ID vs UUID Security Rule

This is a strict security convention enforced across the entire application.

---

## Rule

- Every database table has two identifiers: `id` (bigint, auto-increment) and `uuid` (varchar(36), unique).
- `id` is internal only. It MUST NEVER be exposed in API requests, responses, URLs, or any external-facing surface.
- `uuid` is the external identifier. All API endpoints, request params, route params, JSON responses, and DTOs MUST use `uuid`.

---

## Why

Exposing auto-increment IDs leaks information about record count, creation order, and makes enumeration attacks trivial. UUIDs are non-sequential and unpredictable, making them safe for external use.

---

## Where to Use What

| Context                        | Use     |
|--------------------------------|---------|
| API route params (`:uuid`)     | `uuid`  |
| JSON request bodies            | `uuid`  |
| JSON response bodies           | `uuid`  |
| Foreign keys in DB             | `id`    |
| Internal joins and queries     | `id`    |
| Service-to-repository calls    | `id`    |
| Repository lookups from API    | `uuid` → resolve to `id` internally |
| Logs (internal/debug)          | `id` is acceptable |

---

## Enforcement

- DTOs and response structs MUST NOT include an `ID` or `id` field. Only `UUID` / `uuid`.
- Handler layer receives `uuid` from the client, passes it to the service layer.
- Service layer resolves `uuid` to the internal `id` via repository (e.g., `FindByUUID`), then uses `id` for all further internal operations (updates, deletes, relation lookups).
- Repository layer accepts `id` (uint64) for internal operations and `uuid` (string) for external-facing lookups.
- Never pass raw `id` values back to the client under any field name.

---

## Example Flow

```
Client → POST /api/v1/users/:uuid/files
  Handler: extracts `uuid` from route param
  Service: calls repo.FindByUUID(uuid) → gets user with internal id
  Service: uses user.ID internally for foreign key (e.g., file.CreatedBy = user.ID)
  Service: returns response DTO with file.UUID, user.UUID — never IDs
Client ← { "uuid": "abc-123", "uploaded_by": "def-456" }
```

---

## Violations to Watch For

- Returning `id` in any JSON response
- Accepting `id` in any request body or query param
- Using `id` in URL route params
- Logging `uuid` where `id` would suffice internally (minor, but prefer `id` in logs)
