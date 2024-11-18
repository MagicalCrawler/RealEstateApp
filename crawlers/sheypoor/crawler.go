package sheypoor

import (
	"context"
	"fmt"
	crawlerModels "github.com/MagicalCrawler/RealEstateApp/models/crawler"
	"github.com/MagicalCrawler/RealEstateApp/types"
	"github.com/MagicalCrawler/RealEstateApp/utils"
	"log"
	"log/slog"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
)

const (
	defaultMaxRetries        = 3
	defaultRetryDelay        = 5 * time.Second
	defaultPageLimit         = 1
	defaultGotoTimeout       = 30000
	defaultMaxScrollAttempts = 5
	scrollWaitDuration       = 2 * time.Second
)

// SheypoorCrawler implements the Crawler interface for the Sheypoor website
type SheypoorCrawler struct {
	baseURL    string
	userAgents []string
	logger     *slog.Logger
}

// NewSheypoorCrawler creates a new instance of SheypoorCrawler
func NewSheypoorCrawler() *SheypoorCrawler {
	return &SheypoorCrawler{
		baseURL: utils.GetConfig("SHEYPOOR_BASE_URL"),
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
		logger: utils.NewLogger("Sheypoor_Crawler"),
	}
}

// Crawl fetches posts for a given city
func (c *SheypoorCrawler) Crawl(ctx context.Context, city crawlerModels.City) ([]crawlerModels.Post, error) {
	pageLimit, err := strconv.Atoi(utils.GetConfig("CRAWLER_PAGE_LIMIT"))
	if err != nil || pageLimit <= 0 {
		pageLimit = defaultPageLimit
	}
	pageURL := fmt.Sprintf("%s/s/%s/real-estate", c.baseURL, city.Slug)
	var allPosts []crawlerModels.Post

	// Initialize Playwright
	pw, err := playwright.Run()
	if err != nil {
		c.logger.Error("could not start Playwright: ", err)
		return nil, fmt.Errorf("could not start Playwright: %w", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		c.logger.Error("could not launch browser: ", err)
		return nil, fmt.Errorf("could not launch browser: %w", err)
	}
	defer browser.Close()

	select {
	case <-ctx.Done():
		return allPosts, ctx.Err()
	default:
	}

	maxPageRetries, err := strconv.Atoi(utils.GetConfig("CRAWLER_MAX_RETRIES"))
	if err != nil || maxPageRetries <= 0 {
		maxPageRetries = defaultMaxRetries
	}

	retryDelaySeconds, err := strconv.Atoi(utils.GetConfig("CRAWLER_RETRY_DELAY"))
	if err != nil || retryDelaySeconds <= 0 {
		retryDelaySeconds = int(defaultRetryDelay.Seconds())
	}
	retryDelay := time.Duration(retryDelaySeconds) * time.Second

	for attempt := 1; attempt <= maxPageRetries; attempt++ {
		c.logger.Info("Crawling page: ", pageURL, "Attempt ", attempt)

		page, err := browser.NewPage()
		if err != nil {
			c.logger.Error("could not create new page: ", err, " | Attempt: ", attempt)
			time.Sleep(retryDelay)
			continue
		}

		err = page.SetExtraHTTPHeaders(map[string]string{
			"User-Agent": c.randomUserAgent(),
			"Referer":    "https://google.com",
		})
		if err != nil {
			c.logger.Error("could not set extra headers: ", err, " | Attempt: ", attempt)
			page.Close()
			time.Sleep(retryDelay)
			continue
		}

		playwrightTimeout, err := strconv.Atoi(utils.GetConfig("PLAYWRIGHT_GOTO_TIMEOUT"))
		if err != nil || playwrightTimeout <= 0 {
			playwrightTimeout = defaultGotoTimeout
		}
		_, err = page.Goto(pageURL, playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateNetworkidle,
			Timeout:   playwright.Float(float64(playwrightTimeout)),
		})
		if err != nil {
			c.logger.Error("Error navigating to: ", pageURL, " | Attempt: ", attempt, " error: ", err)
			page.Close()
			time.Sleep(retryDelay)
			continue
		}

		// Scroll to load all content
		postLinks, err := c.autoScroll(page)
		if err != nil {
			c.logger.Error("Error during auto-scroll: ", err, " | Attempt: ", attempt)
		}

		var mu sync.Mutex
		chunkSize := 15
		chunks := splitIntoChunks(postLinks, chunkSize)

		for _, chunk := range chunks {
			var wg sync.WaitGroup

			for _, link := range chunk {
				select {
				case <-ctx.Done():
					return allPosts, ctx.Err()
				default:
				}

				wg.Add(1)

				go func(link string) {
					defer wg.Done()

					post, err := c.CrawlPostDetails(ctx, link)
					if err != nil {
						c.logger.Error("Error crawling post: ", link, " | error: ", err)
						return
					}
					post.City = city

					// Append to allPosts in a thread-safe manner
					mu.Lock()
					allPosts = append(allPosts, post)
					mu.Unlock()
				}(link)
			}

			wg.Wait()
		}
	}
	return allPosts, nil
}

// CrawlPostDetails fetches details for a single post
func (c *SheypoorCrawler) CrawlPostDetails(ctx context.Context, postURL string) (crawlerModels.Post, error) {
	var post crawlerModels.Post

	maxRetries, err := strconv.Atoi(utils.GetConfig("CRAWLER_MAX_RETRIES"))
	if err != nil || maxRetries <= 0 {
		maxRetries = defaultMaxRetries
	}

	retryDelaySeconds, err := strconv.Atoi(utils.GetConfig("CRAWLER_RETRY_DELAY"))
	if err != nil || retryDelaySeconds <= 0 {
		retryDelaySeconds = int(defaultRetryDelay.Seconds())
	}
	retryDelay := time.Duration(retryDelaySeconds) * time.Second

	// Initialize Playwright
	pw, err := playwright.Run()
	if err != nil {
		c.logger.Error("could not start Playwright: ", err)
		return post, fmt.Errorf("could not start Playwright: %w", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		c.logger.Error("could not launch browser: ", err)
		return post, fmt.Errorf("could not launch browser: %w", err)
	}
	defer browser.Close()

	for attempt := 1; attempt <= maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return post, ctx.Err()
		default:
		}

		page, err := browser.NewPage()
		if err != nil {
			c.logger.Error("could not create new page: ", err, " | Attempt: ", attempt)
			log.Printf("Attempt %d: could not create new page: %v", attempt, err)
			time.Sleep(retryDelay)
			continue
		}

		err = page.SetExtraHTTPHeaders(map[string]string{
			"User-Agent": c.randomUserAgent(),
			"Referer":    "https://google.com",
		})
		if err != nil {
			c.logger.Error("could not set extra headers:", err, " | Attempt: ", attempt)
			page.Close()
			time.Sleep(retryDelay)
			continue
		}

		playwrightTimeout, err := strconv.Atoi(utils.GetConfig("PLAYWRIGHT_GOTO_TIMEOUT"))
		if err != nil || playwrightTimeout <= 0 {
			playwrightTimeout = defaultGotoTimeout
		}
		_, err = page.Goto(postURL, playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateNetworkidle,
			Timeout:   playwright.Float(float64(playwrightTimeout)),
		})
		if err != nil {
			c.logger.Error("Error navigating to: ", postURL, " | Attempt: ", attempt, " error: ", err)
			page.Close()
			time.Sleep(retryDelay)
			continue
		}

		content, err := page.Content()
		if err != nil {
			c.logger.Error("could not get page content: ", err, " | Attempt: ", attempt)
			page.Close()
			time.Sleep(retryDelay)
			continue
		}
		page.Close()

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
		if err != nil {
			c.logger.Error("could not parse HTML from: ", postURL, " | Attempt: ", attempt, " error: ", err)
			time.Sleep(retryDelay)
			continue
		}

		// Extract location details
		var locationDetails []string
		doc.Find("nav#UVpPz ul li a").Each(func(i int, s *goquery.Selection) {
			locationDetails = append(locationDetails, strings.TrimSpace(s.Text()))
		})
		if len(locationDetails) > 0 {
			post.City = crawlerModels.City{
				Name:  locationDetails[0], // Assuming the first element is the city name
				Slug:  strings.ToLower(strings.ReplaceAll(locationDetails[0], " ", "-")),
				Level: "city",
			}
			if len(locationDetails) > 1 {
				post.Neighborhood = strings.Join(locationDetails[1:], ", ")
			}
		}

		// Extract property images
		var imageUrls []string
		doc.Find("div.swiper-slide img").Each(func(i int, s *goquery.Selection) {
			if src, exists := s.Attr("src"); exists && src != "" {
				imageUrls = append(imageUrls, src)
			}
		})
		if len(imageUrls) > 0 {
			post.Images = imageUrls
		}

		// Extract title of the post
		if title := doc.Find("h1#listing-title").Text(); title != "" {
			post.Title = strings.TrimSpace(title)
		}

		// Extract price information
		if price := doc.Find("div.tOq3m span strong").Text(); price != "" {
			post.Price = strings.TrimSpace(price)
		}

		// Extract property features
		var features []string
		doc.Find("div.C7Rh9").Each(func(i int, s *goquery.Selection) {
			featureName := s.Find("p._2e124").Text()
			featureValue := s.Find("p._874-x").Text()
			if featureName != "" && featureValue != "" {
				feature := featureName + ": " + featureValue
				features = append(features, strings.TrimSpace(feature))
			}
		})
		if len(features) > 0 {
			post.Features = features
		}

		// Extract description
		if description, err := doc.Find("div.VNOCj div.MQJ5W").Html(); err == nil {
			post.Description = strings.TrimSpace(description)
		}

		// Extract area, rooms, year built, and other relevant information
		doc.Find("div.C7Rh9").Each(func(i int, s *goquery.Selection) {
			featureName := s.Find("p._2e124").Text()
			featureValue := s.Find("p._874-x").Text()
			if featureName != "" && featureValue != "" {
				switch strings.TrimSpace(featureName) {
				case "متراژ":
					post.Area = strings.TrimSpace(featureValue)
				case "سال ساخت":
					post.YearBuilt = strings.TrimSpace(featureValue)
				case "اتاق‌ها":
					post.Rooms = strings.TrimSpace(featureValue)
				case "قیمت هر متر مربع":
					post.PricePerSquareMeter = strings.TrimSpace(featureValue)
				case "طبقه":
					post.Floor = strings.TrimSpace(featureValue)
				}
			}
		})

		// Extract rental specific metadata if applicable
		rentalMetadata := &crawlerModels.RentalMetadata{}
		if capacity := doc.Find("div.rental-capacity").Text(); capacity != "" {
			rentalMetadata.Capacity = strings.TrimSpace(capacity)
		}
		if normalDayPrice := doc.Find("span.normal-day-price").Text(); normalDayPrice != "" {
			rentalMetadata.NormalDayPrice = strings.TrimSpace(normalDayPrice)
		}
		if weekendPrice := doc.Find("span.weekend-price").Text(); weekendPrice != "" {
			rentalMetadata.WeekendPrice = strings.TrimSpace(weekendPrice)
		}
		if holidayPrice := doc.Find("span.holiday-price").Text(); holidayPrice != "" {
			rentalMetadata.HolidayPrice = strings.TrimSpace(holidayPrice)
		}
		if extraPersonCost := doc.Find("span.extra-person-cost").Text(); extraPersonCost != "" {
			rentalMetadata.ExtraPersonCost = strings.TrimSpace(extraPersonCost)
		}

		if rentalMetadata.Capacity != "" || rentalMetadata.NormalDayPrice != "" || rentalMetadata.WeekendPrice != "" || rentalMetadata.HolidayPrice != "" || rentalMetadata.ExtraPersonCost != "" {
			post.RentalMetadata = rentalMetadata
		}

		post.Website = types.Sheypoor

		return post, nil
	}

	c.logger.Error("failed to crawl post details from", postURL, " after ", maxRetries, " attempts")
	return post, fmt.Errorf("failed to crawl post details from %s after %d attempts", postURL, maxRetries)
}

// Helper methods

// Helper function to split a slice into chunks of specified size
func splitIntoChunks(postLinks []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(postLinks); i += chunkSize {
		end := i + chunkSize
		if end > len(postLinks) {
			end = len(postLinks)
		}
		chunks = append(chunks, postLinks[i:end])
	}
	return chunks
}

func (c *SheypoorCrawler) randomUserAgent() string {
	return c.userAgents[rand.Intn(len(c.userAgents))]
}

func (c *SheypoorCrawler) autoScroll(page playwright.Page) ([]string, error) {
	var allLinks []string

	maxScrollAttempts, err := strconv.Atoi(utils.GetConfig("CRAWLER_MAX_SCROLL_ATTEMPTS"))
	if err != nil || maxScrollAttempts <= 0 {
		maxScrollAttempts = defaultMaxScrollAttempts
	}

	scrollAttempts := 0
	noNewContentAttempts := 0

	for scrollAttempts < maxScrollAttempts {
		// Scroll to bottom
		_, err := page.Evaluate(`() => {
            window.scrollTo(0, document.body.scrollHeight);
            return document.querySelectorAll('a.flex').length;
        }`)
		if err != nil {
			c.logger.Error("error scrolling page: ", err)
			return nil, fmt.Errorf("error scrolling page: %w", err)
		}

		// Wait for new content
		time.Sleep(scrollWaitDuration)

		// Get page content
		content, err := page.Content()
		if err != nil {
			c.logger.Error("error getting page content: ", err)
			return nil, fmt.Errorf("error getting page content: %w", err)
		}

		// Parse the page content with goquery
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
		if err != nil {
			c.logger.Error("error parsing page content: ", err)
			return nil, fmt.Errorf("error parsing page content: %w", err)
		}

		// Extract links using the extractPostLinksFromSelection method
		links := c.extractPostLinksFromSelection(doc)

		// Append only new links to the list
		newLinks := make(map[string]bool)
		for _, link := range allLinks {
			newLinks[link] = true
		}
		newLinksFound := false
		for _, link := range links {
			if !newLinks[link] {
				allLinks = append(allLinks, link)
				newLinks[link] = true
				newLinksFound = true
			}
		}

		// Check if new links were found
		if newLinksFound {
			noNewContentAttempts = 0 // Reset counter if new content is found
		} else {
			noNewContentAttempts++
		}

		// Exit if no new content is found after several attempts
		if noNewContentAttempts >= 3 {
			c.logger.Info("No new content found after multiple attempts, stopping scroll.")
			break
		}

		// Attempt to click "load more" button if available
		buttonExists, err := page.Evaluate(`() => {
            const button = document.querySelector('div.post-list__load-more-btn-container-cef96 button');
            if (button) {
                button.click();
                return true;
            }
            return false;
        }`)
		if err != nil {
			c.logger.Error("error clicking load more button: ", err)
			return nil, fmt.Errorf("error clicking load more button: %w", err)
		}

		if buttonExists.(bool) {
			// Wait for new content to load after clicking the button
			time.Sleep(scrollWaitDuration)
		}

		scrollAttempts++
	}

	return allLinks, nil
}

func (c *SheypoorCrawler) extractPostLinksFromSelection(doc *goquery.Document) []string {
	var postLinks []string
	doc.Find("a.flex").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists {
			postLinks = append(postLinks, fmt.Sprintf("%s%s", c.baseURL, link))
		}
	})
	return postLinks
}

func (c *SheypoorCrawler) extractPostDetails(doc *goquery.Document, post *crawlerModels.Post) {
	// Extract the price and area
	post.Price = strings.TrimSpace(doc.Find("span.text-heading-4-bolder").Text())
	post.Area = strings.TrimSpace(doc.Find("span.text-heading-6-bolder").First().Text())

	// Extract neighborhood
	post.Neighborhood = strings.TrimSpace(doc.Find("small.text-heading-6-lighter").First().Text())

	// Extract images
	var images []string
	doc.Find("img.rounded-lg").Each(func(i int, s *goquery.Selection) {
		if src, exists := s.Attr("src"); exists {
			images = append(images, src)
		}
	})
	post.Images = images
}
