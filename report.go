package vfd

import (
	"context"
	"crypto/rsa"
	"github.com/vfdcloud/vfd/models"
)

type (
	ReportSubmitter func(ctx context.Context, url string, headers *RequestHeaders,
		privateKey *rsa.PrivateKey,
		report *models.Report) (*Response, error)
)
