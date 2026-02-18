---
inclusion: always
---

# PERFORMANCE & SECURITY RULES

## CORE PRINCIPLES
- Speed first
- Security always
- No unnecessary loops
- Raw SQL queries only

## DATABASE QUERIES

### USE RAW SQL
- Always use raw SQL with sqlx
- Never use ORM
- Select only needed columns

Example:
```go
q := `SELECT id, name, email FROM users WHERE id = ?`
err := r.DB.Get(&user, q, id)
```

### AVOID N+1 QUERIES
- Use JOIN instead of loops
- Fetch related data in single query

❌ WRONG:
```go
users := getUsers()
for _, user := range users {
    roles := getRolesByUserID(user.ID)
}
```

✅ CORRECT:
```go
q := `SELECT u.id, u.name, r.name as role_name 
      FROM users u 
      LEFT JOIN user_roles ur ON ur.user_id = u.id
      LEFT JOIN roles r ON r.id = ur.role_id`
```

### SELECT SPECIFIC COLUMNS
- Never use `SELECT *`
- Only select needed columns

❌ WRONG:
```go
q := `SELECT * FROM users WHERE id = ?`
```

✅ CORRECT:
```go
q := `SELECT id, name, email FROM users WHERE id = ?`
```

## LOOPS

### AVOID LOOPS ON QUERIES
- No loops with database calls inside
- Batch operations when needed

❌ WRONG:
```go
for _, id := range ids {
    r.DB.Exec(`UPDATE users SET status = ? WHERE id = ?`, status, id)
}
```

✅ CORRECT:
```go
q := `UPDATE users SET status = ? WHERE id IN (?)`
r.DB.Exec(q, status, ids)
```

### MINIMIZE LOOPS
- Avoid unnecessary iterations
- Use database operations instead

## SECURITY

### PARAMETERIZED QUERIES
- Always use `?` placeholders
- Never concatenate SQL strings

❌ WRONG:
```go
q := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email)
```

✅ CORRECT:
```go
q := `SELECT id, name FROM users WHERE email = ?`
r.DB.Get(&user, q, email)
```

### INPUT VALIDATION
- Validate in handler with `bind.AndValidate`
- Check required fields

## GENERAL RULES
- Avoid nested loops
- Use indexes on database columns
- Cache frequently accessed data in Redis
- Return early on errors
