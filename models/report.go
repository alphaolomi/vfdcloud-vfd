package models

import (
	"encoding/xml"
	"fmt"
)

type (
	// ZACK is the Z report acknowledge received from vfd server after
	// successfully Z report upload request
	ZACK struct {
		XMLName xml.Name `xml:"ZACK"`
		Text    string   `xml:",chardata"`
		ZNUMBER int64    `xml:"ZNUMBER"`
		DATE    string   `xml:"DATE"`
		TIME    string   `xml:"TIME"`
		ACKCODE int64    `xml:"ACKCODE"`
		ACKMSG  string   `xml:"ACKMSG"`
	}

	ReportAckEFDMS struct {
		XMLName        xml.Name `xml:"EFDMS"`
		Text           string   `xml:",chardata"`
		ZACK           ZACK     `xml:"ZACK"`
		EFDMSSIGNATURE string   `xml:"EFDMSSIGNATURE"`
	}
	ReportEFDMS struct {
		XMLName        xml.Name `xml:"EFDMS"`
		Text           string   `xml:",chardata"`
		ZREPORT        Report   `xml:"ZREPORT"`
		EFDMSSIGNATURE string   `xml:"EFDMSSIGNATURE"`
	}
	// Lines ...
	Lines struct {
		Name    string
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
		lines.Street,
		lines.Mobile,
		fmt.Sprintf("%s,%s", lines.City, lines.Country)}
}
