package divar

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/MagicalCrawler/RealEstateApp/crawlers"
	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
)

type DivarRealEstateCrawler struct {
	baseURL    string
	userAgents []string
}

// NewDivarRealEstateCrawler creates a new instance of DivarRealEstateCrawler
func NewDivarRealEstateCrawler() *DivarRealEstateCrawler {
	return &DivarRealEstateCrawler{
		baseURL: "https://divar.ir",
		userAgents: []string{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
			"Mozilla/5.0 (Linux; Android 10; SM-G975F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Mobile Safari/537.36",
			"Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15A372 Safari/604.1",
			"Mozilla/5.0 (Linux; Android 11; Pixel 5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Mobile Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:92.0) Gecko/20100101 Firefox/92.0",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Safari/605.1.15",
			"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
			"Mozilla/5.0 (Linux; Android 9; SM-G960F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.89 Mobile Safari/537.36",
			"Mozilla/5.0 (iPad; CPU OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15A5341f Safari/604.1",
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Safari/537.36",
			"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.11; rv:42.0) Gecko/20100101 Firefox/42.0",
			"Mozilla/5.0 (Linux; Android 10; SM-A505F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Mobile Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64; Trident/7.0; AS; rv:11.0) like Gecko",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
			"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:89.0) Gecko/20100101 Firefox/89.0",
			"Mozilla/5.0 (Linux; U; Android 9; en-US; SM-G960U Build/PPR1.180610.011) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.89 Mobile Safari/537.36",
		},
	}
}

// RandomUserAgent selects a random User-Agent string from the list
func (c *DivarRealEstateCrawler) RandomUserAgent() string {
	return c.userAgents[rand.Intn(len(c.userAgents))]
}

// GetCities fetches city URLs from the main page
func (c *DivarRealEstateCrawler) GetCities(ctx context.Context) ([]string, error) {
	var cityURLs []string

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start Playwright: %v", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		return nil, fmt.Errorf("could not create new page: %w", err)
	}
	defer page.Close()

	// Navigate to the base URL and load the page
	if _, err = page.Goto(c.baseURL); err != nil {
		return nil, fmt.Errorf("could not navigate to %s: %w", c.baseURL, err)
	}

	// Adding a short delay to ensure the page fully loads
	time.Sleep(2 * time.Second)

	// Extracting city links from the page
	elements, err := page.QuerySelectorAll("a[class*=cities__item]") // Adjust selector based on the actual class of city links
	if err != nil {
		return nil, fmt.Errorf("could not get city links: %w", err)
	}

	for _, element := range elements {
		href, err := element.GetAttribute("href")
		if err == nil && strings.HasPrefix(href, "/s/") {
			cityURLs = append(cityURLs, fmt.Sprintf("%s%s", c.baseURL, href))
		}
	}

	return cityURLs, nil
}

// CrawlCity performs crawling for each city URL
func (c *DivarRealEstateCrawler) CrawlCity(ctx context.Context, cityURL string, pageLimit int) ([]crawlers.Post, error) {
	paginationURLs := getPaginationURLs(cityURL, pageLimit)
	var allPosts []crawlers.Post

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start Playwright: %v", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	defer browser.Close()

	for _, pageURL := range paginationURLs {
		fmt.Println("Crawling page: ", pageURL)

		page, err := browser.NewPage()
		if err != nil {
			return nil, fmt.Errorf("could not create new page: %w", err)
		}

		page.SetExtraHTTPHeaders(map[string]string{
			"User-Agent":      c.RandomUserAgent(),
			"Accept-Language": "en-US,en;q=0.9",
			"Referer":         "https://google.com",
			"Accept-Encoding": "gzip, deflate, br",
		})

		if _, err = page.Goto(pageURL, playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateNetworkidle,
		}); err != nil {
			return nil, fmt.Errorf("could not navigate to %s: %w", pageURL, err)
		}

		time.Sleep(3 * time.Second) // Adjust the delay as needed

		content, err := page.Content()
		if err != nil {
			return nil, fmt.Errorf("could not get page content: %w", err)
		}
		page.Close()

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
		if err != nil {
			return nil, fmt.Errorf("failed to parse HTML from %s: %w", pageURL, err)
		}

		doc.Find("div.kt-post-card__body").Each(func(i int, s *goquery.Selection) {
			post := c.extractPostFromSelection(s)
			allPosts = append(allPosts, post)
		})

		time.Sleep(time.Duration(rand.Intn(7)+3) * time.Second)
	}

	return allPosts, nil
}

// getPaginationURLs generates URLs for multiple pages based on the pagination limit
func getPaginationURLs(base string, limit int) []string {
	var pages []string
	for i := 1; i <= limit; i++ {
		if i == 1 {
			pages = append(pages, base)
		} else {
			pages = append(pages, fmt.Sprintf("%s?page=%d", base, i))
		}
	}
	return pages
}

// extractPostFromSelection extracts post information from a goquery selection
func (c *DivarRealEstateCrawler) extractPostFromSelection(s *goquery.Selection) crawlers.Post {
	title := strings.TrimSpace(s.Find("h2.kt-post-card__title").Text())
	if title == "" {
		title = "No Title"
	}

	price := strings.TrimSpace(s.Find("div.kt-post-card__description").Text())
	if price == "" {
		price = "No Price"
	}

	link, exists := s.Parent().Attr("href")
	if !exists {
		link = "No Link"
	}

	return crawlers.Post{
		Title: title,
		Price: price,
		Link:  fmt.Sprintf("https://divar.ir%s", link),
	}
}
