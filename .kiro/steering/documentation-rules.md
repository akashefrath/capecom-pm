---
inclusion: auto
description: Documentation rules - strict rule to never create .md files in codebase, use .docs/ folder instead
---

# Documentation Rules

## CRITICAL RULE: NO MARKDOWN FILES IN CODEBASE

**NEVER create markdown (.md) files in the codebase for documentation purposes.**

### ❌ FORBIDDEN

- Creating README.md files in feature directories
- Creating CHANGELOG.md files
- Creating TODO.md files
- Creating any .md files for documentation
- Creating documentation files in code directories

### ✅ ALLOWED

- Code comments (inline documentation)
- JSDoc/GoDoc style comments
- Package-level documentation in code files
- Steering files in `.kiro/steering/` (architecture patterns only)

### 📁 Documentation Location

If documentation is needed for future reference or AI context:

**Use `.docs/` folder in project root**

```
project/
├── .docs/              # All documentation goes here
│   ├── api/           # API documentation
│   ├── architecture/  # Architecture docs
│   ├── guides/        # How-to guides
│   └── decisions/     # Architecture decision records
├── .kiro/
│   └── steering/      # Only architecture patterns
└── internal/          # NO .md files here
```

### Examples

**❌ WRONG:**
```
internal/
├── services/
│   ├── user.go
│   └── README.md          # NEVER DO THIS
├── handlers/
│   ├── user_handler.go
│   └── DOCUMENTATION.md   # NEVER DO THIS
```

**✅ CORRECT:**
```
.docs/
├── services/
│   └── user-service.md    # Documentation here
├── api/
│   └── endpoints.md       # API docs here

internal/
├── services/
│   └── user.go            # Only code, with comments
├── handlers/
│   └── user_handler.go    # Only code, with comments
```

### Code Documentation

Use inline comments and package documentation:

```go
// Package services provides business logic layer.
// All services follow the same pattern:
// - Accept DTOs from handlers
// - Return domain errors
// - Call repositories for data access
package services

// UserService handles user-related business logic.
type UserService struct {
    repo *repositories.UserRepo
}

// Create creates a new user with validation.
// Returns ErrDuplicateEmail if email already exists.
func (s *UserService) Create(req CreateUserRequest) (*UserResponse, error) {
    // Implementation with inline comments
}
```

### When to Create Documentation

**Only create documentation files in `.docs/` when:**
- Documenting complex business logic
- Creating API reference
- Recording architecture decisions
- Writing integration guides
- Documenting deployment processes

**Never create documentation files for:**
- Simple CRUD operations (use code comments)
- Obvious code structure
- Temporary notes (use TODO comments in code)
- Feature summaries (use commit messages)

### Summary

- ✅ Code comments: YES
- ✅ `.docs/` folder: YES for documentation
- ✅ `.kiro/steering/`: YES for architecture patterns only
- ❌ `.md` files in code directories: NEVER
- ❌ README.md in feature folders: NEVER
- ❌ Documentation files mixed with code: NEVER

**Keep code clean. Keep documentation separate.**
