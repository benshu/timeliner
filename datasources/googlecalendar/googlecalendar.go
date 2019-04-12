package googlecalendar

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mholt/timeliner"
	"golang.org/x/net/context"
	"google.golang.org/api/calendar/v3"
)

const (
	DataSourceName = "Google Calendar"
	DataSourceID   = "google_calendar"
)

var dataSource = timeliner.DataSource{
	ID:   DataSourceID,
	Name: DataSourceID,
	OAuth2: timeliner.OAuth2{
		ProviderID: "google",
		Scopes: []string{
			"https://www.googleapis.com/auth/calendar.readonly",
		},
	},
	RateLimit: timeliner.RateLimit{
		RequestsPerHour: 10000 / 24, // https://developers.google.com/photos/library/guides/api-limits-quotas
		BurstSize:       3,
	},
	NewClient: func(acc timeliner.Account) (timeliner.Client, error) {
		httpClient, err := acc.NewHTTPClient()
		if err != nil {
			return nil, err
		}
		return &Client{
			HTTPClient: httpClient,
			userID:     acc.UserID,
		}, nil
	},
}

func init() {
	err := timeliner.RegisterDataSource(dataSource)
	if err != nil {
		log.Fatal(err)
	}
}

// Client interacts with the Google Photos
// API. It requires an OAuth2-authorized
// HTTP client in order to work properly.
type Client struct {
	HTTPClient *http.Client

	userID string
}

// ListItems lists items from the data source.
// opt.Timeframe precision is day-level at best.
func (c *Client) ListItems(ctx context.Context, itemChan chan<- *timeliner.ItemGraph, opt timeliner.Options) error {
	defer close(itemChan)

	if opt.Filename != "" {
		return fmt.Errorf("importing data from a file is not supported")
	}

	// get items and collections
	errChan := make(chan error)
	go func() {
		err := c.listItems(ctx, itemChan, opt.Timeframe)
		errChan <- err
	}()

	var errs []string
	for i := 0; i < 1; i++ {
		err := <-errChan
		if err != nil {
			log.Printf("[ERROR][%s/%s] A listing goroutine errored: %v", DataSourceID, c.userID, err)
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("one or more errors: %s", strings.Join(errs, ", "))
	}

	return nil
}

func (c *Client) listItems(ctx context.Context, itemChan chan<- *timeliner.ItemGraph, timeframe timeliner.Timeframe) error {
	srv, err := calendar.New(c.HTTPClient)

	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		return fmt.Errorf("getting items on next page: %v", err)
	}
	for _, item := range events.Items {
		var event eventItem
		log.Printf("[info] %v", item)

		itemChan <- &timeliner.ItemGraph{
			Node: event,
		}
	}

	return nil

}
