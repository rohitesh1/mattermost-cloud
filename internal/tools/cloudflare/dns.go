// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
//

package cloudflare

import (
	"context"
	"strings"
	"time"

	cf "github.com/cloudflare/cloudflare-go"
	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
)

const defaultTimeout = 30 * time.Second

func (c *Client) getZoneID(zoneName string) (zoneID string, err error) {
	zoneID, err = c.cfClient.ZoneIDByName(zoneName)
	if err != nil {
		return "", err
	}

	return zoneID, nil
}

func (c *Client) getZoneName(zoneNameList []string, customerDNSName string) (zoneName string, found bool) {
	for _, zoneName := range zoneNameList {
		if zoneName == "" {
			return "", false
		}
		if strings.HasSuffix(customerDNSName, zoneName) {
			return zoneName, true
		}
	}
	return "", false
}

func (c *Client) getRecordIDs(zoneID, customerDNSName string, logger logrus.FieldLogger) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	dnsRecords, err := c.cfClient.DNSRecords(ctx, zoneID, cf.DNSRecord{Name: customerDNSName})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get DNS Record ID from Cloudflare")
	}
	if len(dnsRecords) == 0 {
		logger.Infof("No DNS records for %q domain found in Cloudflare %q zone", customerDNSName, zoneID)
		return nil, nil
	}

	ids := make([]string, 0, len(dnsRecords))
	for _, rec := range dnsRecords {
		ids = append(ids, rec.ID)
	}

	return ids, nil

}

// CreateDNSRecords creates a DNS records in the first given Cloudflare zone name of the list
func (c *Client) CreateDNSRecords(dnsNames []string, dnsEndpoints []string, logger logrus.FieldLogger) error {
	if len(dnsNames) == 0 {
		return errors.New("no domain names provided")
	}
	if len(dnsEndpoints) == 0 {
		return errors.New("no DNS endpoints provided for Cloudflare creation request")
	}
	if len(dnsEndpoints) > 1 {
		return errors.New("creating record for more than one endpoint not supported")
	}
	dnsEndpoint := dnsEndpoints[0]
	if dnsEndpoint == "" {
		return errors.New("DNS endpoint was an empty string")
	}

	zoneNameList := c.aws.GetPublicHostedZoneNames()
	if len(zoneNameList) == 0 {
		return errors.New("no public hosted zones names found from AWS")
	}

	for _, dnsName := range dnsNames {
		// Fetch the zone name for that customer DNS name
		zoneName, found := c.getZoneName(zoneNameList, dnsName)
		if !found {
			return errors.Errorf("hosted zone for %q domain name not found", dnsName)
		}

		// Fetch the zone ID
		zoneID, err := c.getZoneID(zoneName)
		if err != nil {
			return errors.Wrap(err, "failed to fetch Zone ID from Cloudflare")
		}

		proxied := true

		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()
		recordResp, err := c.cfClient.CreateDNSRecord(ctx, zoneID, cf.DNSRecord{
			Name:    dnsName,
			Type:    "CNAME",
			Content: dnsEndpoint,
			TTL:     1,
			Proxied: &proxied,
		})
		if err != nil {
			return errors.Wrap(err, "failed to create DNS Record at Cloudflare")
		}

		logger.WithFields(logrus.Fields{
			"cloudflare-dns-value":    dnsName,
			"cloudflare-dns-endpoint": dnsEndpoint,
			"cloudflare-zone-id":      zoneID,
		}).Debugf("Cloudflare create DNS record response: %v", recordResp)
	}

	return nil
}

// DeleteDNSRecords gets DNS name and zone name which uses to delete that DNS record from Cloudflare
func (c *Client) DeleteDNSRecords(dnsNames []string, logger logrus.FieldLogger) error {
	zoneNameList := c.aws.GetPublicHostedZoneNames()
	if len(zoneNameList) == 0 {
		return errors.New("no public hosted zones names found from AWS")
	}

	for _, dnsName := range dnsNames {
		// Fetch the zone name for that customer DNS name
		zoneName, found := c.getZoneName(zoneNameList, dnsName)
		if !found {
			return errors.Errorf("hosted zone for %q domain name not found", dnsName)
		}

		// Fetch the zone ID
		zoneID, err := c.getZoneID(zoneName)
		if err != nil {
			return errors.Wrap(err, "failed to fetch Zone ID from Cloudflare")
		}

		recordIDs, err := c.getRecordIDs(zoneID, dnsName, logger)
		if err != nil {
			return errors.Wrapf(err, "Failed to get record ID from Cloudflare for DNS: %s", dnsName)
		}

		for _, recID := range recordIDs {
			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			defer cancel()

			err = c.cfClient.DeleteDNSRecord(ctx, zoneID, recID)
			if err != nil {
				return errors.Wrap(err, "Failed to delete DNS Record at Cloudflare")
			}
		}
	}
	return nil
}
