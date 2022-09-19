package models

import (
	"encoding/xml"
	"fmt"
	"strings"
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

	ZREPORT struct {
		XMLName xml.Name `xml:"ZREPORT"`
		Text    string   `xml:",chardata"`
		DATE    string   `xml:"DATE"`
		TIME    string   `xml:"TIME"`
		HEADER  struct {
			Text string   `xml:",chardata"`
			LINE []string `xml:"LINE"`
		} `xml:"HEADER"`
		VRN              string `xml:"VRN"`
		TIN              string `xml:"TIN"`
		TAXOFFICE        string `xml:"TAXOFFICE"`
		REGID            string `xml:"REGID"`
		ZNUMBER          string `xml:"ZNUMBER"`
		EFDSERIAL        string `xml:"EFDSERIAL"`
		REGISTRATIONDATE string `xml:"REGISTRATIONDATE"`
		USER             string `xml:"USER"`
		SIMIMSI          string `xml:"SIMIMSI"`
		TOTALS           struct {
			Text             string `xml:",chardata"`
			DAILYTOTALAMOUNT string `xml:"DAILYTOTALAMOUNT"`
			GROSS            string `xml:"GROSS"`
			CORRECTIONS      string `xml:"CORRECTIONS"`
			DISCOUNTS        string `xml:"DISCOUNTS"`
			SURCHARGES       string `xml:"SURCHARGES"`
			TICKETSVOID      string `xml:"TICKETSVOID"`
			TICKETSVOIDTOTAL string `xml:"TICKETSVOIDTOTAL"`
			TICKETSFISCAL    string `xml:"TICKETSFISCAL"`
			TICKETSNONFISCAL string `xml:"TICKETSNONFISCAL"`
		} `xml:"TOTALS"`
		VATTOTALS struct {
			Text       string   `xml:",chardata"`
			VATRATE    []string `xml:"VATRATE"`
			NETTAMOUNT []string `xml:"NETTAMOUNT"`
			TAXAMOUNT  []string `xml:"TAXAMOUNT"`
		} `xml:"VATTOTALS"`
		PAYMENTS struct {
			Text      string   `xml:",chardata"`
			PMTTYPE   []string `xml:"PMTTYPE"`
			PMTAMOUNT []string `xml:"PMTAMOUNT"`
		} `xml:"PAYMENTS"`
		CHANGES struct {
			Text          string `xml:",chardata"`
			VATCHANGENUM  string `xml:"VATCHANGENUM"`
			HEADCHANGENUM string `xml:"HEADCHANGENUM"`
		} `xml:"CHANGES"`
		ERRORS     string `xml:"ERRORS"`
		FWVERSION  string `xml:"FWVERSION"`
		FWCHECKSUM string `xml:"FWCHECKSUM"`
	}
)

func (lines *Lines) List() []string {
	return []string{
		strings.ToUpper(lines.Name),
		strings.ToUpper(lines.Street),
		fmt.Sprintf("MOBILE: %s", lines.Mobile),
		strings.ToUpper(fmt.Sprintf("%s,%s", lines.City, lines.Country))}
}
