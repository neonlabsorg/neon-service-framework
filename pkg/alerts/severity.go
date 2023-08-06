package alerts

type Severity int

const DEFAULT_SEVERITY = Info

const (
	Debug Severity = iota
	Info
	Warning
	Error
	Critical
	Urgent
)

func (level Severity) String() string {
	names := [...]string{
		"debug",
		"info",
		"warning",
		"error",
		"critical",
		"urgent",
	}

	if level < Debug || level > Urgent {
		return "unknown"
	}

	return names[level]
}
