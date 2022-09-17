package models

import "fmt"

type (
	// Lines ...
	Lines struct {
		Name    string
		Address string
		Street  string
		Mobile  string
		City    string
		Country string
	}
	Report struct{}
)

// List takes a list of lines and returns array of lines
func (lines *Lines) List() []string {
	return []string{
		lines.Name,
		fmt.Sprintf("%s,%s", lines.Address, lines.Street),
		lines.Mobile,
		fmt.Sprintf("%s,%s", lines.City, lines.Country)}
}
