package utils

import (
	"fmt"
	"strings"
)

const urlFormatErr = "url should be in form of subdomain.domain.com; found: %s"

type Url struct {
	Subdomain string
	ZoneName  string
}

func UnpackSubdomainUrl(fullHostname string) (*Url, error) {
	url := &Url{}

	// getting domain
	domainsArr := strings.Split(fullHostname, ".")
	if len(domainsArr) < 2 {
		return nil, fmt.Errorf(urlFormatErr+"; need a subdomain", fullHostname)
	}
	tld := domainsArr[len(domainsArr)-1]
	zoneName := domainsArr[len(domainsArr)-2]
	url.ZoneName = fmt.Sprintf("%s.%s", zoneName, tld)

	// getting subdomain
	subdomain := ""
	for i := 0; i < len(domainsArr)-2; i++ {
		subdomain += domainsArr[i] + "."
	}
	subdomain, _ = strings.CutSuffix(subdomain, ".")
	url.Subdomain = subdomain

	return url, nil
}

func (u *Url) GetFullHostname() string {
	return fmt.Sprintf("%s.%s", u.Subdomain, u.ZoneName)
}
