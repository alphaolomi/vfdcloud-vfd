package models

import (
	"encoding/xml"
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
	// Address ...

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
		VRN              string       `xml:"VRN"`
		TIN              string       `xml:"TIN"`
		TAXOFFICE        string       `xml:"TAXOFFICE"`
		REGID            string       `xml:"REGID"`
		ZNUMBER          string       `xml:"ZNUMBER"`
		EFDSERIAL        string       `xml:"EFDSERIAL"`
		REGISTRATIONDATE string       `xml:"REGISTRATIONDATE"`
		USER             string       `xml:"USER"`
		SIMIMSI          string       `xml:"SIMIMSI"`
		TOTALS           REPORTTOTALS `xml:"TOTALS"`
		VATTOTALS        VATTOTALS    `xml:"VATTOTALS"`
		PAYMENTS         PAYMENTS     `xml:"PAYMENTS"`
		CHANGES          struct {
			Text          string `xml:",chardata"`
			VATCHANGENUM  string `xml:"VATCHANGENUM"`
			HEADCHANGENUM string `xml:"HEADCHANGENUM"`
		} `xml:"CHANGES"`
		ERRORS     string `xml:"ERRORS"`
		FWVERSION  string `xml:"FWVERSION"`
		FWCHECKSUM string `xml:"FWCHECKSUM"`
	}

	REPORTTOTALS struct {
		XMLName          xml.Name `xml:"TOTALS"`
		Text             string   `xml:",chardata"`
		DAILYTOTALAMOUNT float64  `xml:"DAILYTOTALAMOUNT"`
		GROSS            float64  `xml:"GROSS"`
		CORRECTIONS      float64  `xml:"CORRECTIONS"`
		DISCOUNTS        float64  `xml:"DISCOUNTS"`
		SURCHARGES       float64  `xml:"SURCHARGES"`
		TICKETSVOID      int64    `xml:"TICKETSVOID"`
		TICKETSVOIDTOTAL float64  `xml:"TICKETSVOIDTOTAL"`
		TICKETSFISCAL    int64    `xml:"TICKETSFISCAL"`
		TICKETSNONFISCAL int64    `xml:"TICKETSNONFISCAL"`
	}
)
