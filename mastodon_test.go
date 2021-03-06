package mastodon

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthenticate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("username") != "valid" || r.FormValue("password") != "user" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w, `{"AccessToken": "zoo"}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
	})
	err := client.Authenticate("invalid", "user")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}

	client = NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
	})
	err = client.Authenticate("valid", "user")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
}

func TestPostStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer zoo" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w, `{"AccessToken": "zoo"}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
	})
	_, err := client.PostStatus(&Toot{
		Status: "foobar",
	})
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}

	client = NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err = client.PostStatus(&Toot{
		Status: "foobar",
	})
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
}

func TestGetTimelineHome(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[{"Content": "foo"}, {"Content": "bar"}]`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
	})
	_, err := client.PostStatus(&Toot{
		Status: "foobar",
	})
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}

	client = NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	tl, err := client.GetTimelineHome()
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(tl) != 2 {
		t.Fatalf("result should be two: %d", len(tl))
	}
	if tl[0].Content != "foo" {
		t.Fatalf("want %q but %q", "foo", tl[0].Content)
	}
	if tl[1].Content != "bar" {
		t.Fatalf("want %q but %q", "bar", tl[1].Content)
	}
}

func TestForTheCoverages(t *testing.T) {
	(*UpdateEvent)(nil).event()
	(*NotificationEvent)(nil).event()
	(*DeleteEvent)(nil).event()
	(*ErrorEvent)(nil).event()
	_ = (&ErrorEvent{io.EOF}).Error()
}

func TestGetAccount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/1234567" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `{"Username": "zzz"}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	a, err := client.GetAccount(1)
	if err == nil {
		t.Fatalf("should not be fail: %v", err)
	}
	a, err = client.GetAccount(1234567)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if a.Username != "zzz" {
		t.Fatalf("want %q but %q", "zzz", a.Username)
	}
}

func TestGetAccountFollowing(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/1234567/following" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `[{"Username": "foo"}, {"Username": "bar"}]`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	fl, err := client.GetAccountFollowing(123)
	if err == nil {
		t.Fatalf("should not be fail: %v", err)
	}
	fl, err = client.GetAccountFollowing(1234567)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(fl) != 2 {
		t.Fatalf("result should be two: %d", len(fl))
	}
	if fl[0].Username != "foo" {
		t.Fatalf("want %q but %q", "foo", fl[0].Username)
	}
	if fl[1].Username != "bar" {
		t.Fatalf("want %q but %q", "bar", fl[0].Username)
	}
}

func TestRegisterApp(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if r.URL.Path != "/api/v1/apps" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if r.FormValue("redirect_uris") != "urn:ietf:wg:oauth:2.0:oob" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w, `{"client_id": "foo", "client_secret": "bar"}`)
		return
	}))
	defer ts.Close()

	app, err := RegisterApp(&AppConfig{
		Server: ts.URL,
		Scopes: "read write follow",
	})
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if app.ClientID != "foo" {
		t.Fatalf("want %q but %q", "foo", app.ClientID)
	}
	if app.ClientSecret != "bar" {
		t.Fatalf("want %q but %q", "bar", app.ClientSecret)
	}
}

func TestStreamingPublic(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/streaming/public" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		f, _ := w.(http.Flusher)
		fmt.Fprintln(w, `
event: update
data: {"Content": "foo"}
		`)
		f.Flush()

		fmt.Fprintln(w, `
event: update
data: {"Content": "bar"}
		`)
		f.Flush()
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	ctx, cancel := context.WithCancel(context.Background())
	q, err := client.StreamingPublic(ctx)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	time.AfterFunc(3*time.Second, func() {
		cancel()
		close(q)
	})
	events := []Event{}
	for e := range q {
		events = append(events, e)
	}
	if len(events) != 2 {
		t.Fatalf("result should be two: %d", len(events))
	}
	if events[0].(*UpdateEvent).Status.Content != "foo" {
		t.Fatalf("want %q but %q", "foo", events[0].(*UpdateEvent).Status.Content)
	}
	if events[1].(*UpdateEvent).Status.Content != "bar" {
		t.Fatalf("want %q but %q", "bar", events[1].(*UpdateEvent).Status.Content)
	}
}
