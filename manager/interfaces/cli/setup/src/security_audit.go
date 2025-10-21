package setup

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type SecurityAuditLogger struct {
	auditPath string
}

type AuditEvent struct {
	Timestamp time.Time              `json:"timestamp"`
	EventType string                 `json:"event_type"`
	Severity  string                 `json:"severity"`
	Resource  string                 `json:"resource"`
	Action    string                 `json:"action"`
	Result    string                 `json:"result"`
	Details   map[string]interface{} `json:"details"`
	Caller    string                 `json:"caller"`
}

func NewSecurityAuditLogger() *SecurityAuditLogger {
	homeDir, _ := os.UserHomeDir()
	auditPath := filepath.Join(homeDir, ".syntropy", "logs", "security-audit.log")
	os.MkdirAll(filepath.Dir(auditPath), 0700)

	return &SecurityAuditLogger{
		auditPath: auditPath,
	}
}

func (sal *SecurityAuditLogger) LogCriticalOperation(eventType, resource, action, result string, details map[string]interface{}) {
	event := AuditEvent{
		Timestamp: time.Now(),
		EventType: eventType,
		Severity:  "CRITICAL",
		Resource:  resource,
		Action:    action,
		Result:    result,
		Details:   details,
		Caller:    getCaller(),
	}

	// Append to audit log
	f, err := os.OpenFile(sal.auditPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return // Fail silently to avoid disrupting main operations
	}
	defer f.Close()

	data, err := json.Marshal(event)
	if err != nil {
		return // Fail silently
	}

	f.Write(append(data, '\n'))
}

func (sal *SecurityAuditLogger) LogWarning(eventType, resource, action, result string, details map[string]interface{}) {
	event := AuditEvent{
		Timestamp: time.Now(),
		EventType: eventType,
		Severity:  "WARNING",
		Resource:  resource,
		Action:    action,
		Result:    result,
		Details:   details,
		Caller:    getCaller(),
	}

	f, err := os.OpenFile(sal.auditPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	defer f.Close()

	data, err := json.Marshal(event)
	if err != nil {
		return
	}

	f.Write(append(data, '\n'))
}

func getCaller() string {
	pc := make([]uintptr, 3)
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])

	var callers []string
	for {
		frame, more := frames.Next()
		callers = append(callers, frame.Function)
		if !more {
			break
		}
	}

	return strings.Join(callers, " -> ")
}
