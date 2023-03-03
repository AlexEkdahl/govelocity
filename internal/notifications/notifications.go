package notifications

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/AlexEkdahl/govelocity/internal/process"
)

// Notifier represents a notifier that sends notifications.
type Notifier interface {
	SendNotification(string) error
}

// Notification represents a notification.
type Notification struct {
	Timestamp  time.Time
	Process    *process.Process
	Message    string
	Recipients []string
}

// Notifications represents a collection of notifications.
type Notifications struct {
	notifications []*Notification
	notifiers     []Notifier
}

// NewNotifications creates a new Notifications instance.
func NewNotifications() *Notifications {
	return &Notifications{}
}

// AddNotifier adds a notifier to the collection.
func (n *Notifications) AddNotifier(notifier Notifier) {
	n.notifiers = append(n.notifiers, notifier)
}

// Notify sends a notification to all the notifiers.
func (n *Notifications) Notify(p *process.Process, message string) error {
	if len(n.notifiers) == 0 {
		return errors.New("no notifiers registered")
	}

	notification := &Notification{
		Timestamp:  time.Now(),
		Process:    p,
		Message:    message,
		Recipients: []string{},
	}

	// Send the notification to all the notifiers
	var errors []string
	for _, notifier := range n.notifiers {
		err := notifier.SendNotification(n.formatNotification(notification))
		if err != nil {
			errors = append(errors, err.Error())
		}
	}

	// Return an error if any of the notifications failed to send
	if len(errors) > 0 {
		return fmt.Errorf("failed to send notifications: %s", strings.Join(errors, "; "))
	}

	// Add the notification to the collection
	n.notifications = append(n.notifications, notification)

	return nil
}

// formatNotification formats the notification as a string.
func (n *Notifications) formatNotification(notification *Notification) string {
	var recipients string
	if len(notification.Recipients) > 0 {
		recipients = fmt.Sprintf(" to %s", strings.Join(notification.Recipients, ", "))
	}
	return fmt.Sprintf("[%s] Process '%s' (%d) %s%s", notification.Timestamp.Format("2006-01-02 15:04:05"), notification.Process.Name, notification.Process.Pid, notification.Message, recipients)
}

// LatestNotifications returns the latest n notifications.
func (n *Notifications) LatestNotifications(count int) []*Notification {
	if count >= len(n.notifications) {
		return n.notifications
	}
	return n.notifications[len(n.notifications)-count:]
}
