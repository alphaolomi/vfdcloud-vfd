package models

import (
	"encoding/xml"
)

const (
	CASH    = PMTYPE("CASH")
	CHEQUE  = PMTYPE("CHEQUE")
	EMONEY  = PMTYPE("EMONEY")
	CCARD   = PMTYPE("CCARD")
	INVOICE = PMTYPE("INVOICE")
)

type (

	// PMTYPE this represents acceptable mode of payment in API referred to as Mode of Payment
	// or Type of Payment
	// Mode of Payment can either be CASH, CHEQUE, EMONEY or CCARD if receipt is generated. In
	// this case payment is already received
	// Mode of Payment can only be INVOICE if invoice is generated.  In this case payment is not
	// yet received that is why we use Invoice
	PMTYPE string

	// RCTACK is the receipt acknowledge received from
	// vfd server after successfully receipt upload request
	RCTACK struct {
		XMLName xml.Name `xml:"RCTACK"`
		Text    string   `xml:",chardata"`
		RCTNUM  int64    `xml:"RCTNUM"`
		DATE    string   `xml:"DATE"`
		TIME    string   `xml:"TIME"`
		ACKCODE int64    `xml:"ACKCODE"`
		ACKMSG  string   `xml:"ACKMSG"`
	}

	RCTACKEFDMS struct {
		XMLName        xml.Name `xml:"EFDMS"`
		Text           string   `xml:",chardata"`
		RCTACK         RCTACK   `xml:"RCTACK"`
		EFDMSSIGNATURE string   `xml:"EFDMSSIGNATURE"`
	}

	RCT struct {
		XMLName    xml.Name  `xml:"RCT"`
		Text       string    `xml:",chardata"`
		DATE       string    `xml:"DATE"`
		TIME       string    `xml:"TIME"`
		TIN        string    `xml:"TIN"`
		REGID      string    `xml:"REGID"`
		EFDSERIAL  string    `xml:"EFDSERIAL"`
		CUSTIDTYPE int64     `xml:"CUSTIDTYPE"`
		CUSTID     string    `xml:"CUSTID"`
		CUSTNAME   string    `xml:"CUSTNAME"`
		MOBILENUM  string    `xml:"MOBILENUM"`
		RCTNUM     string    `xml:"RCTNUM"`
		DC         int64     `xml:"DC"`
		GC         int64     `xml:"GC"`
		ZNUM       string    `xml:"ZNUM"`
		RCTVNUM    string    `xml:"RCTVNUM"`
		ITEMS      ITEMS     `xml:"ITEMS"`
		TOTALS     TOTALS    `xml:"TOTALS"`
		PAYMENTS   PAYMENTS  `xml:"PAYMENTS"`
		VATTOTALS  VATTOTALS `xml:"VATTOTALS"`
	}

	PAYMENT struct {
		XMLName   xml.Name `xml:"PAYMENT"`
		Text      string   `xml:",chardata"`
		PMTTYPE   string   `xml:"PMTTYPE"`
		PMTAMOUNT float64  `xml:"PMTAMOUNT"`
	}

	ITEMS struct {
		XMLName xml.Name `xml:"ITEMS"`
		Text    string   `xml:",chardata"`
		ITEM    []*ITEM  `xml:"ITEM"`
	}

	ITEM struct {
		XMLName xml.Name `xml:"ITEM"`
		Text    string   `xml:",chardata"`
		ID      string   `xml:"ID"`
		DESC    string   `xml:"DESC"`
		QTY     float64  `xml:"QTY"`
		TAXCODE int64    `xml:"TAXCODE"`
		AMT     float64  `xml:"AMT"`
	}

	TOTALS struct {
		XMLName      xml.Name `xml:"TOTALS"`
		Text         string   `xml:",chardata"`
		TOTALTAXEXCL float64  `xml:"TOTALTAXEXCL"`
		TOTALTAXINCL float64  `xml:"TOTALTAXINCL"`
		DISCOUNT     float64  `xml:"DISCOUNT"`
	}

	VATTOTAL struct {
		XMLName    xml.Name `xml:"VATTOTAL"`
		Text       string   `xml:",chardata"`
		VATRATE    string   `xml:"VATRATE"`
		NETTAMOUNT float64  `xml:"NETTAMOUNT"`
		TAXAMOUNT  float64  `xml:"TAXAMOUNT"`
	}

	VATTOTALS struct {
		XMLName  xml.Name    `xml:"VATTOTALS"`
		Text     string      `xml:",chardata"`
		VATTOTAL []*VATTOTAL `xml:"VATTOTAL"`
	}

	PAYMENTS struct {
		XMLName xml.Name   `xml:"PAYMENTS"`
		Text    string     `xml:",chardata"`
		PAYMENT []*PAYMENT `xml:"PAYMENT"`
	}
)
