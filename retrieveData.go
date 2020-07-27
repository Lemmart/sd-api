package main

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/sirupsen/logrus"

	"net/http"
	"net/url"
)

// schema: companyHostName -> offerCodes
var uniqueOffersMap map[string]bool

func scrape(website string, chCompanies chan *Company, chFinished chan bool) {
	// retrieve company information
	parsedUrl, err := url.Parse(website)
	if err != nil {
		logrus.WithError(err).Error("failed to parse website website")
	}
	company := &Company{
		Name:    parsedUrl.Hostname(),
		Website: website,
	}

	// return completed job results
	defer func() {
		chCompanies <- company
		chFinished <- true
	}()

	resp, err := http.Get(website)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"company": company.Name,
			"website": company.Website,
		}).WithError(err).Error("Failed to retrieve website")
		return
	}

	defer resp.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logrus.WithError(err).Error("Error loading HTTP response body")
	}

	document.Find("body").Each(func(_ int, element *goquery.Selection) {
		result := strings.TrimPrefix(element.Text(), "")

		// use lazy match regex to consume the first instance of a match
		pattern, err := regexp.Compile(`.*?(\d{1,2})?% off[^}]*?[cC]ode[ :].*?(\w+).*?[</]?`)
		if err != nil {
			logrus.WithError(err).Error("Failed to compile regex")
		}

		offersFound := pattern.FindAllStringSubmatch(result, -1)

		for _, offer := range offersFound {
			if len(offer) > 1 {
				offerCode := strings.ToUpper(strings.TrimSpace(offer[2]))
				offerAmt, err := strconv.Atoi(offer[1])
				if err != nil {
					logrus.WithError(err).WithFields(logrus.Fields{
						"Amount": offerAmt,
						"Code":   offerCode,
					}).Error("Failed to convert offer amount")
					continue
				}

				offer := &Offer{
					Amount:   offerAmt,
					Code:     offerCode,
					Category: "",
				}

				if _, found := uniqueOffersMap[offerCode]; found {
					logrus.WithFields(logrus.Fields{
						"Amount":  offerAmt,
						"Code":    offerCode,
						"Company": company.Name,
					}).Info("Duplicate offer detected")
					continue
				}

				company.Offers = append(company.Offers, offer)
				uniqueOffersMap[offerCode] = true

				logrus.WithFields(logrus.Fields{
					"Amount":  offerAmt,
					"Code":    offerCode,
					"Company": company.Name,
				}).Info("Saved offer")

			}
		}
	})
}

func GetData(request SdsRequest) (*SdsResponse, error) {
	// "clear cache" -- todo: build this as a cache per company
	uniqueOffersMap = make(map[string]bool)

	var companies []*Company

	// Channels for tracking completed scraping jobs
	chCompany := make(chan *Company)
	chFinished := make(chan bool)

	// Kick off the scraping process (concurrently)
	for _, website := range request.Websites {
		go scrape(website, chCompany, chFinished)
	}

	// Subscribe to both channels
	for c := 0; c < len(request.Websites); {
		select {
		case company := <-chCompany:
			companies = append(companies, company)
		case <-chFinished:
			logrus.WithFields(logrus.Fields{
				"current": c + 1,
				"total":   len(request.Websites),
			}).Info("Requested scrapes")
			c++
		}
	}

	return &SdsResponse{
		Companies: companies,
	}, nil
}
