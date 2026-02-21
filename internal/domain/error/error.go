package domainerrors

import "errors"

var (
	ErrDuplicateError = errors.New("duplicate error")
	// ErrInvalidCredentials auth
	ErrInvalidCredentials = errors.New("invalid_login_credentials")
	ErrUnauthorized       = errors.New("unauthorized")

	// ErrUserNotFound user
	ErrUserNotFound   = errors.New("user_not_found")
	ErrDuplicateEmail = errors.New("duplicate_email")

	// Client errors
	ErrClientNotFound  = errors.New("client_not_found")
	ErrDuplicateClient = errors.New("duplicate_client")

	// Project errors
	ErrProjectNotFound      = errors.New("project_not_found")
	ErrDuplicateProject     = errors.New("duplicate_project")
	ErrProjectAssetNotFound = errors.New("project_asset_not_found")

	// Ticket errors
	ErrTicketNotFound     = errors.New("ticket_not_found")
	ErrDuplicateTicket    = errors.New("duplicate_ticket")
	ErrTicketTypeNotFound = errors.New("ticket_type_not_found")
	ErrNotProjectMember   = errors.New("not_project_member")
	ErrNotTicketAssignee  = errors.New("not_ticket_assignee")
	ErrHoursExceeded      = errors.New("hours_exceeded")
	ErrTimeEntryNotFound  = errors.New("time_entry_not_found")

	// File errors
	ErrFileNotFound    = errors.New("file_not_found")
	ErrFileNotUploaded = errors.New("file_not_uploaded")

	// Master data errors
	ErrGroupNotFound       = errors.New("group_not_found")
	ErrDesignationNotFound = errors.New("designation_not_found")
	ErrDepartmentNotFound  = errors.New("department_not_found")
	ErrRoleNotFound        = errors.New("role_not_found")
	ErrInvalidRoles        = errors.New("invalid_roles")

	// JWT errors
	ErrInvalidToken         = errors.New("invalid_token")
	ErrTokenExpired         = errors.New("token_expired")
	ErrInvalidSigningMethod = errors.New("invalid_signing_method")
	ErrTokenTypeMismatch    = errors.New("token_type_mismatch")

	//ErrBadRequest common
	ErrBadRequest = errors.New("bad_request")
	ErrInternal   = errors.New("internal_error")

	RoleNotFound = errors.New("role_not_found")

	CantPerformThis = errors.New("cant_perform_this_action")

	InvalidPunchOrder = errors.New("invalid_punch_order")
	DuplicatePunch    = errors.New("duplicate_punch")

	// Clock In
	AlreadyClockedIn = errors.New("already_clocked_in")
	MustClockInFirst = errors.New("must_clock_in_first")

	// Clock Out
	AlreadyClockedOut = errors.New("already_clocked_out")
	NotClockedIn      = errors.New("not_clocked_in")
	MustBreakOutFirst = errors.New("must_break_out_first")

	// Break
	AlreadyOnBreak = errors.New("already_on_break")
	NotOnBreak     = errors.New("not_on_break")

	// Timeout
	AlreadyTimedOut = errors.New("already_timed_out")

	// Data / State
	AttendanceNotFound = errors.New("attendance_not_found")
	SummaryNotFound    = errors.New("attendance_summary_not_found")
)
