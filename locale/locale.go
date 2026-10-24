// Package locale owns library-level accessible labels. Applications set their
// language once; the library never registers a language or reads the host
// locale behind the caller's back.
package locale

import "sync/atomic"

// Messages holds user-visible library text that is not supplied by the caller.
type Messages struct {
	Loading       string
	Close         string
	Clear         string
	Search        string
	ShowPassword  string
	HidePassword  string
	Increase      string
	Decrease      string
	Previous      string
	Next          string
	MorePages     string
	Remove        string
	Notifications string
	Required      string
	Invalid       string
	NoResults     string
	Today         string
	Selected      string
	Empty         string
}

var current atomic.Pointer[Messages]

func init() { Set(Messages{}) }

// Set publishes library messages. Empty fields keep the English default.
func Set(m Messages) {
	d := Messages{
		Loading: "Loading", Close: "Close", Clear: "Clear", Search: "Search",
		ShowPassword: "Show characters", HidePassword: "Conceal characters",
		Increase: "Increase", Decrease: "Decrease", Previous: "Previous", Next: "Next",
		MorePages: "More pages", Remove: "Remove", Notifications: "Notifications",
		Required: "Required", Invalid: "Invalid value", NoResults: "No results",
		Today: "Today", Selected: "Selected", Empty: "Nothing here",
	}
	override := func(dst *string, src string) {
		if src != "" {
			*dst = src
		}
	}
	override(&d.Loading, m.Loading)
	override(&d.Close, m.Close)
	override(&d.Clear, m.Clear)
	override(&d.Search, m.Search)
	override(&d.ShowPassword, m.ShowPassword)
	override(&d.HidePassword, m.HidePassword)
	override(&d.Increase, m.Increase)
	override(&d.Decrease, m.Decrease)
	override(&d.Previous, m.Previous)
	override(&d.Next, m.Next)
	override(&d.MorePages, m.MorePages)
	override(&d.Remove, m.Remove)
	override(&d.Notifications, m.Notifications)
	override(&d.Required, m.Required)
	override(&d.Invalid, m.Invalid)
	override(&d.NoResults, m.NoResults)
	override(&d.Today, m.Today)
	override(&d.Selected, m.Selected)
	override(&d.Empty, m.Empty)
	current.Store(&d)
}

// M returns the published messages.
func M() Messages {
	m := current.Load()
	if m == nil {
		Set(Messages{})
		m = current.Load()
	}
	return *m
}
