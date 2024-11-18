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
			// Add more user agents as needed...
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

		// Parse post details
		splitURL := strings.Split(postURL, "/")
		post.ID = splitURL[len(splitURL)-1]
		post.Link = postURL
		post.Title = strings.TrimSpace(doc.Find("h2.text-heading-4-bolder").Text())
		post.Description = strings.TrimSpace(doc.Find("small.text-heading-6-lighter").First().Text())
		post.Website = types.Sheypoor

		// Check if essential details are present
		if post.Title == "" || post.Description == "" {
			c.logger.Error("Missing essential post details for: ", postURL, " | Attempt: ", attempt)
			time.Sleep(retryDelay)
			continue
		}

		// Extract additional details
		c.extractPostDetails(doc, &post)

		// If successful, return
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
