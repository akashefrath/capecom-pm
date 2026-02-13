# Tickets & Time Entries API Examples

Base URL: `{{baseUrl}}{{apiPath}}`  
Example: `http://localhost:8080/api/v1`

All endpoints require Bearer token authentication.

---

## Tickets

### 1. Create Ticket
**POST** `/project/{projectId}/tickets`

**Request:**
```json
{
  "title": "Implement login page",
  "description": "Build the login UI with email/password",
  "ticket_type_id": "t1y2p3e4-5678-90ab-cdef-1234567890ab",
  "assigned_to": "u1s2e3r4-5678-90ab-cdef-1234567890ab",
  "estimated_hours": 8,
  "internal_estimated_hours": 6,
  "priority": "high",
  "start_date": "2026-02-14",
  "end_date": "2026-02-20",
  "due_date": "2026-02-20"
}
```

**Response (201):**
```json
{
  "success": true,
  "data": {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "project_id": "550e8400-e29b-41d4-a716-446655440000",
    "code": "ECOM-42",
    "title": "Implement login page",
    "description": "Build the login UI with email/password",
    "ticket_type_id": "t1y2p3e4-5678-90ab-cdef-1234567890ab",
    "ticket_type_name": "Feature",
    "assigned_to": "u1s2e3r4-5678-90ab-cdef-1234567890ab",
    "assigned_to_name": "John Doe",
    "lifecycle_status": "todo",
    "priority": "high",
    "estimated_hours": 8,
    "created_at": "2026-02-13T10:00:00Z"
  }
}
```

---

### 2. List Tickets (with Pagination)
**GET** `/project/{projectId}/tickets?page=1&limit=20`

**Response (200):**
```json
{
  "success": true,
  "data": {
    "data": [
      {
        "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "code": "ECOM-42",
        "title": "Implement login page",
        "lifecycle_status": "in_progress",
        "priority": "high"
      }
    ],
    "meta": {
      "page": 1,
      "limit": 20,
      "total": 1,
      "has_more": false
    }
  }
}
```

---

### 3. Get Single Ticket
**GET** `/project/{projectId}/tickets/{ticketId}`


**Response (200):**
```json
{
  "success": true,
  "data": {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "project_id": "550e8400-e29b-41d4-a716-446655440000",
    "code": "ECOM-42",
    "title": "Implement login page",
    "description": "Build the login UI with email/password",
    "branch": "feature/login-page",
    "ticket_type_id": "t1y2p3e4-5678-90ab-cdef-1234567890ab",
    "ticket_type_name": "Feature",
    "assigned_to": "u1s2e3r4-5678-90ab-cdef-1234567890ab",
    "assigned_to_name": "John Doe",
    "assigned_by": "m1n2g3r4-5678-90ab-cdef-1234567890ab",
    "assigned_by_name": "Jane Manager",
    "lifecycle_status": "in_progress",
    "priority": "high",
    "estimated_hours": 8,
    "internal_estimated_hours": 6,
    "start_date": "2026-02-14",
    "end_date": "2026-02-20",
    "due_date": "2026-02-20",
    "status": "active",
    "created_at": "2026-02-13T10:00:00Z",
    "updated_at": "2026-02-13T14:30:00Z"
  }
}
```

---

### 4. Update Ticket
**PUT** `/project/{projectId}/tickets/{ticketId}`

**Request:**
```json
{
  "title": "Updated: Implement login page with 2FA",
  "priority": "urgent",
  "estimated_hours": 12,
  "due_date": "2026-03-01"
}
```

**Response (200):** Same as Get Single Ticket

---

### 5. Update Lifecycle Status
**PATCH** `/project/{projectId}/tickets/{ticketId}/lifecycle`

**Request:**
```json
{
  "lifecycle_status": "in_progress"
}
```

**Valid statuses:** `todo`, `in_progress`, `in_review`, `testing`, `done`, `closed`, `reopened`

---

### 6. Reassign Ticket
**PATCH** `/project/{projectId}/tickets/{ticketId}/assignee`

**Request:**
```json
{
  "assigned_to": "u1s2e3r4-5678-90ab-cdef-1234567890ab"
}
```

**Note:** User must be a project member.

---

### 7. Delete Ticket (Soft Delete)
**DELETE** `/project/{projectId}/tickets/{ticketId}`

**Response (200):**
```json
{
  "success": true,
  "message": "ticket_deleted"
}
```

---

## Time Entries

### 1. Create Time Entry
**POST** `/project/{projectId}/tickets/{ticketId}/time-entries`

**Authorization:** Only the ticket assignee can add time entries.

**Request:**
```json
{
  "work_date": "2026-02-13",
  "hours": 4.5,
  "description": "Implemented login validation and error handling"
}
```

**Response (201):**
```json
{
  "success": true,
  "data": {
    "id": "e1n2t3r4-5678-90ab-cdef-1234567890ab",
    "ticket_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "ticket_code": "ECOM-42",
    "project_id": "550e8400-e29b-41d4-a716-446655440000",
    "user_id": "u1s2e3r4-5678-90ab-cdef-1234567890ab",
    "user_name": "John Doe",
    "work_date": "2026-02-13",
    "hours": 4.5,
    "description": "Implemented login validation and error handling",
    "status": "active",
    "created_at": "2026-02-13T18:00:00Z"
  }
}
```

---

### 2. List Time Entries
**GET** `/project/{projectId}/tickets/{ticketId}/time-entries?page=1&limit=20`

**Response (200):**
```json
{
  "success": true,
  "data": {
    "data": [
      {
        "id": "e1n2t3r4-5678-90ab-cdef-1234567890ab",
        "ticket_code": "ECOM-42",
        "user_name": "John Doe",
        "work_date": "2026-02-13",
        "hours": 4.5,
        "description": "Implemented login validation and error handling"
      },
      {
        "id": "e2n2t3r4-5678-90ab-cdef-1234567890cd",
        "ticket_code": "ECOM-42",
        "user_name": "John Doe",
        "work_date": "2026-02-14",
        "hours": 3.0,
        "description": "Added unit tests"
      }
    ],
    "meta": {
      "page": 1,
      "limit": 20,
      "total": 2,
      "has_more": false
    }
  }
}
```

---

### 3. Get Single Time Entry
**GET** `/project/{projectId}/tickets/{ticketId}/time-entries/{entryId}`

**Response (200):** Same structure as Create response

---

### 4. Update Time Entry
**PUT** `/project/{projectId}/tickets/{ticketId}/time-entries/{entryId}`

**Authorization:** Only the entry owner can update.

**Request:**
```json
{
  "hours": 5.0,
  "description": "Implemented login validation, error handling, and unit tests"
}
```

**Response (200):** Updated time entry

---

### 5. Delete Time Entry
**DELETE** `/project/{projectId}/tickets/{ticketId}/time-entries/{entryId}`

**Authorization:** Only the entry owner can delete.

**Response (200):**
```json
{
  "success": true,
  "message": "time_entry_deleted"
}
```

---

## Ticket History

### Get Ticket History
**GET** `/project/{projectId}/tickets/{ticketId}/history?page=1&limit=20`

**Response (200):**
```json
{
  "success": true,
  "data": {
    "data": [
      {
        "id": "h1s2t3o4-5678-90ab-cdef-1234567890ab",
        "ticket_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "changed_by": "u1s2e3r4-5678-90ab-cdef-1234567890ab",
        "changed_by_name": "John Doe",
        "field_name": "time_entry_added",
        "old_value": null,
        "new_value": "2026-02-13",
        "note": null,
        "created_at": "2026-02-13T18:00:00Z"
      },
      {
        "id": "h2s2t3o4-5678-90ab-cdef-1234567890cd",
        "ticket_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "changed_by": "m1n2g3r4-5678-90ab-cdef-1234567890ab",
        "changed_by_name": "Jane Manager",
        "field_name": "lifecycle_status",
        "old_value": "todo",
        "new_value": "in_progress",
        "note": null,
        "created_at": "2026-02-13T14:30:00Z"
      }
    ],
    "meta": {
      "page": 1,
      "limit": 20,
      "total": 2,
      "has_more": false
    }
  }
}
```

---

## Error Responses

### 403 Forbidden (Not Ticket Assignee)
```json
{
  "success": false,
  "message": "not_ticket_assignee"
}
```

### 404 Not Found
```json
{
  "success": false,
  "message": "ticket_not_found"
}
```

### 400 Bad Request (Validation Error)
```json
{
  "success": false,
  "errors": {
    "hours": ["must be at least 0.01"],
    "work_date": ["required"]
  }
}
```
