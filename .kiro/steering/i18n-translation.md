# I18N TRANSLATION RULES

## CORE PRINCIPLE
All user-facing messages must support internationalization using the i18n utility.

## TRANSLATION UTILITIES

### Location
`internal/utils/i18n/` - Translation message definitions
`internal/utils/utils.go` - Translation helper functions

### Available Functions

#### GetMessage
Get simple translated message by key:
```go
utils.GetMessage("key", c)
```

#### GetMessageWithExtra
Get translated message with dynamic placeholders:
```go
utils.GetMessageWithExtra("func_success", c, "logout")
```

## HANDLER RESPONSE PATTERNS

### Success Response with Translation
```go
response.JSONOk(c, response.APIResponse{
    Success: true,
    Message: utils.GetMessageWithExtra("func_success", c, "logout"),
    Data: data,
})
```

### Error Response
Use `response.FromError(c, err)` - automatically translates error messages:
```go
if err != nil {
    response.FromError(c, err)
    return
}
```

### Simple Success without Message
```go
response.JSONOk(c, response.APIResponse{
    Success: true,
    Data: token,
})
```

## ERROR TRANSLATION FLOW

1. Service returns domain error (e.g., `domainerrors.ErrInvalidCredentials`)
2. Handler calls `response.FromError(c, err)`
3. Error mapper:
   - Maps error to HTTP status code
   - Gets language from `Accept-Language` header
   - Looks up translated message using error key
   - Returns translated response

## LANGUAGE DETECTION
Language is detected from `Accept-Language` HTTP header in all translation functions.

## WHEN TO USE TRANSLATIONS

### Always Translate
- Success messages shown to users
- Error messages
- Validation messages
- Status messages

### No Translation Needed
- Debug logs
- Internal error details
- Data values (IDs, emails, etc.)

## EXAMPLE HANDLER

```go
func (h *Handler) Action(c *gin.Context) {
    var req dto.Request
    isValid := bind.AndValidate(c, &req, "key")
    if !isValid {
        return
    }
    
    data, err := h.Service.Action(req)
    if err != nil {
        response.FromError(c, err)
        return
    }
    
    response.JSONOk(c, response.APIResponse{
        Success: true,
        Message: utils.GetMessageWithExtra("func_success", c, "action_name"),
        Data: data,
    })
}
```

## RULES
- Never hardcode user-facing messages in English
- Always pass `*gin.Context` to translation functions
- Use `GetMessage` for simple messages
- Use `GetMessageWithExtra` for messages with placeholders
- Let `response.FromError` handle all error translations
