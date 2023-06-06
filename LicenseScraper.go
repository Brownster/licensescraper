package collectors

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"gopkg.in/alecthomas/kingpin.v2"
	"vcs-maas.maas.services.sabio.co.uk/exporters/maas"
)

// NewLicenseCollector creates a new ScheduledScraper instance for collecting license metrics.
func NewLicenseCollector(app *kingpin.Application) *maas.ScheduledScraper {
	const timeout = time.Second * 5
	const frequency = time.Hour * 12

	return maas.NewScheduledScraper("license", NewLicenseScraper(app),
		maas.WithSchedule(maas.NewSchedule(maas.WithTimeout(timeout), maas.WithFrequency(frequency))),
		maas.WithDescription(app, "used", "Peak used licenses", []string{"name"}),
		maas.WithDescription(app, "capacity", "Total available licenses", []string{"name"}),
		maas.WithDescription(app, "parent_info", "License parent information", []string{"name", "folder"}),
		maas.WithDescription(app, "name", "Feature name", []string{"name"}),
		maas.WithDescription(app, "display_name", "Feature display name", []string{"display_name"}),
		maas.WithDescription(app, "value", "Feature value", []string{"value"}),
		maas.WithDescription(app, "expiration_date", "Feature expiration date", []string{"expiration_date"}),
		maas.WithDescription(app, "expiration_days_remaining", "Remaining time until the license expires in days", []string{"display_name"}),
	)
}

// LicenseScraper contains the path to the data folder containing UsageHistory.xml files.
type LicenseScraper struct {
	dataPath string
}

// FeatureUsageHistory represents the XML structure of the UsageHistory.xml files.
type FeatureUsageHistory struct {
	XMLName xml.Name `xml:"FeatureUsageHistory"`
	Text    string   `xml:",chardata"`
	Feature []struct {
		Text          string `xml:",chardata"`
		ID            string `xml:"id,attr"`
		Name          string `xml:"Name"`
		DisplayName   string `xml:"DisplayName"`
		Value         string `xml:"Value"`
		ExpirationDate string `xml:"ExpirationDate"`
		Usage         struct {
			Text      string `xml:",chardata"`
			UpdatedOn string `xml:"updatedOn,attr"`
		} `xml:"Usage"`
		Capacity string `xml:"Capacity"`
	} `xml:"Feature"`
	Signature string `xml:"Signature"`
}

// NewLicenseScraper initializes a new LicenseScraper instance.
func NewLicenseScraper(app *kingpin.Application) *LicenseScraper {
	ls := &LicenseScraper{}

	app.Flag("license.datapath", "Location of the data folder containing UsageHistory.xml files").
		Required().
		StringVar(&ls.dataPath)

	return ls
}

// Scrape performs the actual scraping of the license data from the XML files.
func (s *LicenseScraper) Scrape(c maas.Connector) ([]maas.Metric, error) {
	// Find all XML files in the specified directory.
	out, err := c.Execute(fmt.Sprintf(`find /opt/Avaya/JBoss/wildfly-10.1.0.Final/avmgmt/configuration/weblm/licenses -iname '*.xml' 2>/dev/null`))
	if err != nil {
		return nil, err
	}

	metrics := make([]maas.Metric, 0)

	// Loop through the found XML files.
	for _, path := range strings.Split(strings.ReplaceAll(out.(string), "\r\n", "\n"), "\n") {
		if path == "" {
			continue
		}
		// Read the content of the XML file.
		out, err := c.Execute(fmt.Sprintf(`cat %s`, path))
		if err != nil {
			log.Warn(err)
			continue
		}

		// Split the path to the file to retrieve the folder name later.
		s := strings.Split(path, string(os.PathSeparator))

		// Unmarshal the XML file content into a FeatureUsageHistory struct.
		var usage FeatureUsageHistory
		err = xml.Unmarshal([]byte(out.(string)), &usage)
		if err != nil {
			log.Warn(err)
			continue
		}

		// Loop through the features in the FeatureUsageHistory struct.
		for _, f := range usage.Feature {
			// Extract the peak used licenses.
			used := strings.Split(f.Usage.Text, ", ")
			if len(used) > 2 {
				metrics = append(metrics, maas.NewMetric("used", prometheus.GaugeValue, cast.ToFloat64(used[1]), []string{f.ID}))
			}

			// Extract the total available licenses.
			capacity := strings.Split(f.Capacity, ", ")
			if len(capacity) > 2 {
				metrics = append(metrics, maas.NewMetric("capacity", prometheus.GaugeValue, cast.ToFloat64(capacity[1]), []string{f.ID}))
			}

			// Add license parent information if available.
			if len(s) >= 2 {
				metrics = append(metrics, maas.NewMetric("parent_info", prometheus.GaugeValue, 1, []string{f.ID, s[len(s)-2]}))
			}
            // Calculate the remaining time until the license expires in days.
            expirationDate, err := time.Parse("2006-01-02", f.ExpirationDate)
            if err == nil {
                now := time.Now()
                remainingTime := expirationDate.Sub(now).Seconds()
                remainingDays := remainingTime / 86400 // Convert remaining time to days
                metrics = append(metrics, maas.NewMetric("expiration_days_remaining", prometheus.GaugeValue, remainingDays, []string{f.DisplayName}))
        }
		}
	}

	// Return the collected metrics.
	return metrics, nil
}
