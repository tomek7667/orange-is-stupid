package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"orange-is-stupid/pkg/cloudflare"
	"orange-is-stupid/pkg/inet"
	"orange-is-stupid/pkg/utils"

	"github.com/joho/godotenv"
)

var (
	cfClient *cloudflare.Client
	cfUrls   []CfUrl
)

type CfUrl struct {
	Url     string
	Proxied bool
}

func init() {
	godotenv.Load(".env")
	if _, found := os.LookupEnv("CF_API_TOKEN"); !found {
		panic("set CF_API_TOKEN and URLS in .env; Example URLS='subdomain.example.com|true,this-will-be-set-to-your-ip.without-proxy.com|false'")
	}

	cfClient = cloudflare.New(
		utils.Env("CF_API_TOKEN", "<here your token or set env>"),
	)
	urlsStr := utils.Env("URLS", "subdomain.example.com|true,this-will-be-set-to-your-ip.without-proxy.com|false")
	selfUrls := strings.Split(urlsStr, ",")
	for _, selfUrl := range selfUrls {
		selfUrlArr := strings.Split(selfUrl, "|")
		if len(selfUrlArr) != 2 {
			panic(fmt.Errorf("invalid env: %s", urlsStr))
		}
		cfUrls = append(cfUrls, CfUrl{
			Url:     selfUrlArr[0],
			Proxied: selfUrlArr[1] == "true",
		})
	}
}

func main() {
	currentIp, err := inet.GetCurrentAddr()
	if err != nil {
		panic(err)
	}
	for _, cfUrl := range cfUrls {
		err = cfClient.UpsertUrl(cfUrl.Url, currentIp.String(), cfUrl.Proxied)
		slog.Info(
			"upsert url status",
			"url", cfUrl.Url,
			"err", err,
		)
	}
}
