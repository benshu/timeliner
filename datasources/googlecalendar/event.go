package googlecalendar

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/mholt/timeliner"
)

type EventPerson struct {
	// DisplayName: The creator's name, if available.
	DisplayName string `json:"displayName,omitempty"`

	// Email: The creator's email address, if available.
	Email string `json:"email,omitempty"`

	// Id: The creator's Profile ID, if available. It corresponds to the id
	// field in the People collection of the Google+ API
	Id string `json:"id,omitempty"`

	// Self: Whether the creator corresponds to the calendar on which this
	// copy of the event appears. Read-only. The default is False.
	Self bool `json:"self,omitempty"`
}

type EventAttendee struct {
	// AdditionalGuests: Number of additional guests. Optional. The default
	// is 0.
	AdditionalGuests int64 `json:"additionalGuests,omitempty"`

	// Comment: The attendee's response comment. Optional.
	Comment string `json:"comment,omitempty"`

	// DisplayName: The attendee's name, if available. Optional.
	DisplayName string `json:"displayName,omitempty"`

	// Email: The attendee's email address, if available. This field must be
	// present when adding an attendee. It must be a valid email address as
	// per RFC5322.
	// Required when adding an attendee.
	Email string `json:"email,omitempty"`

	// Id: The attendee's Profile ID, if available. It corresponds to the id
	// field in the People collection of the Google+ API
	Id string `json:"id,omitempty"`

	// Optional: Whether this is an optional attendee. Optional. The default
	// is False.
	Optional bool `json:"optional,omitempty"`

	// Organizer: Whether the attendee is the organizer of the event.
	// Read-only. The default is False.
	Organizer bool `json:"organizer,omitempty"`

	// Resource: Whether the attendee is a resource. Can only be set when
	// the attendee is added to the event for the first time. Subsequent
	// modifications are ignored. Optional. The default is False.
	Resource bool `json:"resource,omitempty"`

	// ResponseStatus: The attendee's response status. Possible values are:
	//
	// - "needsAction" - The attendee has not responded to the invitation.
	//
	// - "declined" - The attendee has declined the invitation.
	// - "tentative" - The attendee has tentatively accepted the invitation.
	//
	// - "accepted" - The attendee has accepted the invitation.
	ResponseStatus string `json:"responseStatus,omitempty"`

	// Self: Whether this entry represents the calendar on which this copy
	// of the event appears. Read-only. The default is False.
	Self bool `json:"self,omitempty"`
}

type EventDateTime struct {
	// Date: The date, in the format "yyyy-mm-dd", if this is an all-day
	// event.
	Date string `json:"date,omitempty"`

	// DateTime: The time, as a combined date-time value (formatted
	// according to RFC3339). A time zone offset is required unless a time
	// zone is explicitly specified in timeZone.
	DateTime string `json:"dateTime,omitempty"`

	// TimeZone: The time zone in which the time is specified. (Formatted as
	// an IANA Time Zone Database name, e.g. "Europe/Zurich".) For recurring
	// events this field is required and specifies the time zone in which
	// the recurrence is expanded. For single events this field is optional
	// and indicates a custom time zone for the event start/end.
	TimeZone string `json:"timeZone,omitempty"`
}

type eventMetadata struct {
	// Attendees: The attendees of the event. See the Events with attendees
	// guide for more information on scheduling events with other calendar
	// users.
	Attendees []*EventAttendee `json:"attendees,omitempty"`

	// Created: Creation time of the event (as a RFC3339 timestamp).
	// Read-only.
	Created string `json:"created,omitempty"`

	// Creator: The creator of the event. Read-only.
	Creator *EventPerson `json:"creator,omitempty"`

	// Description: Description of the event. Optional.
	Description string `json:"description,omitempty"`

	// End: The (exclusive) end time of the event. For a recurring event,
	// this is the end time of the first instance.
	End *EventDateTime `json:"end,omitempty"`

	// EndTimeUnspecified: Whether the end time is actually unspecified. An
	// end time is still provided for compatibility reasons, even if this
	// attribute is set to True. The default is False.
	EndTimeUnspecified bool `json:"endTimeUnspecified,omitempty"`

	// HtmlLink: An absolute link to this event in the Google Calendar Web
	// UI. Read-only.
	HtmlLink string `json:"htmlLink,omitempty"`

	// ICalUID: Event unique identifier as defined in RFC5545. It is used to
	// uniquely identify events accross calendaring systems and must be
	// supplied when importing events via the import method.
	// Note that the icalUID and the id are not identical and only one of
	// them should be supplied at event creation time. One difference in
	// their semantics is that in recurring events, all occurrences of one
	// event have different ids while they all share the same icalUIDs.
	ICalUID string `json:"iCalUID,omitempty"`

	// Id: Opaque identifier of the event. When creating new single or
	// recurring events, you can specify their IDs. Provided IDs must follow
	// these rules:
	// - characters allowed in the ID are those used in base32hex encoding,
	// i.e. lowercase letters a-v and digits 0-9, see section 3.1.2 in
	// RFC2938
	// - the length of the ID must be between 5 and 1024 characters
	// - the ID must be unique per calendar  Due to the globally distributed
	// nature of the system, we cannot guarantee that ID collisions will be
	// detected at event creation time. To minimize the risk of collisions
	// we recommend using an established UUID algorithm such as one
	// described in RFC4122.
	// If you do not specify an ID, it will be automatically generated by
	// the server.
	// Note that the icalUID and the id are not identical and only one of
	// them should be supplied at event creation time. One difference in
	// their semantics is that in recurring events, all occurrences of one
	// event have different ids while they all share the same icalUIDs.
	Id string `json:"id,omitempty"`

	// Kind: Type of the resource ("calendar#event").
	Kind string `json:"kind,omitempty"`

	// Location: Geographic location of the event as free-form text.
	// Optional.
	Location string `json:"location,omitempty"`

	// Organizer: The organizer of the event. If the organizer is also an
	// attendee, this is indicated with a separate entry in attendees with
	// the organizer field set to True. To change the organizer, use the
	// move operation. Read-only, except when importing an event.
	Organizer *EventPerson `json:"organizer,omitempty"`

	// OriginalStartTime: For an instance of a recurring event, this is the
	// time at which this event would start according to the recurrence data
	// in the recurring event identified by recurringEventId. It uniquely
	// identifies the instance within the recurring event series even if the
	// instance was moved to a different time. Immutable.
	OriginalStartTime *EventDateTime `json:"originalStartTime,omitempty"`

	// Recurrence: List of RRULE, EXRULE, RDATE and EXDATE lines for a
	// recurring event, as specified in RFC5545. Note that DTSTART and DTEND
	// lines are not allowed in this field; event start and end times are
	// specified in the start and end fields. This field is omitted for
	// single events or instances of recurring events.
	Recurrence []string `json:"recurrence,omitempty"`

	// RecurringEventId: For an instance of a recurring event, this is the
	// id of the recurring event to which this instance belongs. Immutable.
	RecurringEventId string `json:"recurringEventId,omitempty"`

	// Sequence: Sequence number as per iCalendar.
	Sequence int64 `json:"sequence,omitempty"`

	// Start: The (inclusive) start time of the event. For a recurring
	// event, this is the start time of the first instance.
	Start *EventDateTime `json:"start,omitempty"`

	// Status: Status of the event. Optional. Possible values are:
	// - "confirmed" - The event is confirmed. This is the default status.
	//
	// - "tentative" - The event is tentatively confirmed.
	// - "cancelled" - The event is cancelled (deleted). The list method
	// returns cancelled events only on incremental sync (when syncToken or
	// updatedMin are specified) or if the showDeleted flag is set to true.
	// The get method always returns them.
	// A cancelled status represents two different states depending on the
	// event type:
	// - Cancelled exceptions of an uncancelled recurring event indicate
	// that this instance should no longer be presented to the user. Clients
	// should store these events for the lifetime of the parent recurring
	// event.
	// Cancelled exceptions are only guaranteed to have values for the id,
	// recurringEventId and originalStartTime fields populated. The other
	// fields might be empty.
	// - All other cancelled events represent deleted events. Clients should
	// remove their locally synced copies. Such cancelled events will
	// eventually disappear, so do not rely on them being available
	// indefinitely.
	// Deleted events are only guaranteed to have the id field populated.
	// On the organizer's calendar, cancelled events continue to expose
	// event details (summary, location, etc.) so that they can be restored
	// (undeleted). Similarly, the events to which the user was invited and
	// that they manually removed continue to provide details. However,
	// incremental sync requests with showDeleted set to false will not
	// return these details.
	// If an event changes its organizer (for example via the move
	// operation) and the original organizer is not on the attendee list, it
	// will leave behind a cancelled event where only the id field is
	// guaranteed to be populated.
	Status string `json:"status,omitempty"`

	// Summary: Title of the event.
	Summary string `json:"summary,omitempty"`
	// Updated: Last modification time of the event (as a RFC3339
	// timestamp). Read-only.
	Updated string `json:"updated,omitempty"`
}

type eventItem struct {
	eventID       string        `json:"id"`
	BaseURL       string        `json:"baseUrl"`
	Description   string        `json:"description"`
	EventMetadata eventMetadata `json:"eventMetadata"`
}

func (m eventItem) ID() string {
	return m.eventID
}

func (m eventItem) Timestamp() time.Time {
	fmt.Println(m.EventMetadata.Created)
	created_time, err := time.Parse("", m.EventMetadata.Created)
	if err != nil {
		log.Fatal(err)
	}

	return created_time
}

func (m eventItem) DataText() (*string, error) {
	return &m.Description, nil
}

func (m eventItem) DataFileName() *string {
	return nil
}

func (m eventItem) DataFileHash() []byte {
	return nil
}

func (m eventItem) DataFileMIMEType() *string {
	return nil
}

func (m eventItem) DataFileReader() (io.ReadCloser, error) {
	return nil, nil
}

func (m eventItem) Owner() (*string, *string) {
	// since we only download event owned by the account,
	// we can leave ID nil and assume the display name
	// is the account owner's name
	return nil, &m.EventMetadata.Organizer.DisplayName
}

func (m eventItem) Class() timeliner.ItemClass {
	return timeliner.ClassEvent
}

func (m eventItem) Metadata() (*timeliner.Metadata, error) {
	meta := &timeliner.Metadata{
		// Attendees:   m.EventMetadata.Attendees,
		// Creator:     m.EventMetadata.Creator,
		Description: m.EventMetadata.Description,
		// HtmlLInk:    m.EventMetadata.HtmlLInk,
		// Kind:        m.EventMetadata.Kind,
	}

	return meta, nil
}

func (m eventItem) Location() (*timeliner.Location, error) {
	// See https://issuetracker.google.com/issues/80379228 😭
	location_free_form_text := m.EventMetadata.Location
	fmt.Println(location_free_form_text)
	latitude := float64(30000)
	longitude := float64(30000)
	location := timeliner.Location{
		Latitude:  &latitude,
		Longitude: &longitude,
	}

	return &location, nil
}
