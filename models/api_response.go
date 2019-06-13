package models

const (
	notificationSeverityError   = "ERROR"
	notificationSeverityNote    = "NOTE"
	notificationSeverityWarning = "WARNING" // TODO consider doing something with the response if it's a warning
	notificationSeveritySuccess = "SUCCESS"
)

type Response interface {
	Error() error
}

// Reply has common stuff on all responses from FedEx API
type Reply struct {
	HighestSeverity string
	Notifications   []Notification
	Version         VersionResponse
	JobID           string `xml:"JobId"`
}

func (r Reply) Error() error {
	if r.HighestSeverity == notificationSeveritySuccess ||
		r.HighestSeverity == notificationSeverityNote ||
		r.HighestSeverity == notificationSeverityWarning {
		return nil
	}

	err := ReplyError{Severity: r.HighestSeverity}
	for _, notification := range r.Notifications {
		if notification.Severity == r.HighestSeverity {
			err.Message = notification.Message
			break
		}
	}

	return err
}

type Notification struct {
	Severity         string
	Source           string
	Code             string
	Message          string
	LocalizedMessage string
}

type VersionResponse struct {
	ServiceID    string `xml:"ServiceId"`
	Major        int
	Intermediate int
	Minor        int
}
