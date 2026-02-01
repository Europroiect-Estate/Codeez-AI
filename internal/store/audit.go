package store

import (
	"encoding/json"
	"time"
)

// ToolAuditEntry records a tool call and its approval.
type ToolAuditEntry struct {
	ID            int64
	SessionID     int64
	Agent         string
	ToolName      string
	ArgsJSON      string
	Approved      bool
	ApprovalScope string
	CreatedAt     time.Time
}

// LogToolAudit inserts an audit entry.
func (s *Store) LogToolAudit(sessionID int64, agent, toolName string, args interface{}, approved bool, approvalScope string) (int64, error) {
	var argsJSON string
	if args != nil {
		b, _ := json.Marshal(args)
		argsJSON = string(b)
	}
	approvedInt := 0
	if approved {
		approvedInt = 1
	}
	res, err := s.DB.Exec(
		`INSERT INTO tool_audit (session_id, agent, tool_name, args_json, approved, approval_scope) VALUES (?, ?, ?, ?, ?, ?)`,
		sessionID, agent, toolName, argsJSON, approvedInt, approvalScope,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetToolAudit returns audit entries for a session.
func (s *Store) GetToolAudit(sessionID int64) ([]ToolAuditEntry, error) {
	rows, err := s.DB.Query(
		`SELECT id, session_id, agent, tool_name, args_json, approved, approval_scope, created_at FROM tool_audit WHERE session_id = ? ORDER BY id`,
		sessionID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ToolAuditEntry
	for rows.Next() {
		var e ToolAuditEntry
		var approvedInt int
		var createdAt string
		if err := rows.Scan(&e.ID, &e.SessionID, &e.Agent, &e.ToolName, &e.ArgsJSON, &approvedInt, &e.ApprovalScope, &createdAt); err != nil {
			return nil, err
		}
		e.Approved = approvedInt != 0
		e.CreatedAt = parseTime(createdAt)
		out = append(out, e)
	}
	return out, rows.Err()
}
