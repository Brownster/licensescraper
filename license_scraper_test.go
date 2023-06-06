package collectors

import (
    "fmt"
    "strings"
    "testing"

    "github.com/spf13/cast"
    "gopkg.in/alecthomas/kingpin.v2"
)

type MockConnector struct {
    files []string
}

func (mc *MockConnector) Execute(cmd string) (interface{}, error) {
    if strings.HasPrefix(cmd, "find") {
        return strings.Join(mc.files, "\n"), nil
    }
    if strings.HasPrefix(cmd, "cat") {
        return `<Your sample XML content here>`, nil
    }
    return nil, fmt.Errorf("unsupported command")
}

func TestLicenseScraper(t *testing.T) {
    mockConnector := &MockConnector{
        files: []string{
            "/path/to/your/sample/XML/file1.xml",
            "/path/to/your/sample/XML/file2.xml",
        },
    }

    app := kingpin.New("test", "Test application")
    scraper := NewLicenseScraper(app)

    metrics, err := scraper.Scrape(mockConnector)
    if err != nil {
        t.Errorf("Scrape() returned an error: %v", err)
    }

    for _, metric := range metrics {
        fmt.Printf("Metric: %s, Value: %f, Labels: %+v\n", metric.Name, metric.Value, metric.Labels)
    }
}

func TestScrapeExpirationDaysRemaining(t *testing.T) {
    mockConnector := &MockConnector{
        files: []string{
            "/path/to/your/sample/XML/file1.xml",
            "/path/to/your/sample/XML/file2.xml",
        },
    }

    app := kingpin.New("test", "Test application")
    scraper := NewLicenseScraper(app)

    metrics, err := scraper.ScrapeExpirationDaysRemaining(mockConnector)
    if err != nil {
        t.Errorf("ScrapeExpirationDaysRemaining() returned an error: %v", err)
    }

    for _, metric := range metrics {
        fmt.Printf("Metric: %s, Value: %f, Labels: %+v\n", metric.Name, metric.Value, metric.Labels)
    }
}
