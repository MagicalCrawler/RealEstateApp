package divar

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/MagicalCrawler/RealEstateApp/crawlers"
	"github.com/PuerkitoBio/goquery"
)

type DivarRealEstateCrawler struct {
	baseURL    string
	httpClient *http.Client
	userAgent  string
}

// NewDivarRealEstateCrawler creates a new instance of DivarRealEstateCrawler
func NewDivarRealEstateCrawler() *DivarRealEstateCrawler {
	return &DivarRealEstateCrawler{
		baseURL: "https://divar.ir/s/tehran/real-estate",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36",
	}
}

// RetryableRequest is a helper function to make HTTP requests with retries
func (c *DivarRealEstateCrawler) RetryableRequest(ctx context.Context, req *http.Request, retries int, wait time.Duration) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i <= retries; i++ {
		resp, err = c.httpClient.Do(req)
		if err == nil && resp.StatusCode != http.StatusTooManyRequests {
			return resp, nil
		}

		if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
			log.Printf("Received 429 status code, on try %d. Waiting for %s before retrying...", i, wait)
			time.Sleep(wait)
		} else {
			break
		}
	}

	return resp, err
}

// Modify the Crawl method to support pagination
func (c *DivarRealEstateCrawler) Crawl(ctx context.Context, pageLimit int) ([]crawlers.Post, error) {
	paginationURLs := getPaginationURLs(c.baseURL, pageLimit)

	var allPosts []crawlers.Post

	for _, pageURL := range paginationURLs {
		fmt.Println("Crawling page: ", pageURL)

		req, err := http.NewRequestWithContext(ctx, "GET", pageURL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request for %s: %w", pageURL, err)
		}
		req.Header.Set("User-Agent", c.userAgent)

		resp, err := c.RetryableRequest(ctx, req, 5, 1*time.Second)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch data from %s: %w", pageURL, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code for %s: %d", pageURL, resp.StatusCode)
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to parse HTML from %s: %w", pageURL, err)
		}

		doc.Find("div.kt-post-card__body").Each(func(i int, s *goquery.Selection) {
			post := c.extractPostFromSelection(s)
			allPosts = append(allPosts, post)
		})

		time.Sleep(1 * time.Second)
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

// CrawlPostDetails fetches additional details for a specific post
func (c *DivarRealEstateCrawler) CrawlPostDetails(ctx context.Context, post crawlers.Post) (crawlers.Post, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", post.Link, nil)
	if err != nil {
		return post, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.RetryableRequest(ctx, req, 5, 1*time.Second)
	if err != nil {
		return post, fmt.Errorf("failed to fetch post details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return post, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return post, fmt.Errorf("failed to parse HTML: %w", err)
	}

	title := strings.TrimSpace(doc.Find("h1.kt-page-title__title").Text())
	if title != "" {
		post.Title = title
	}

	price := strings.TrimSpace(doc.Find("p.kt-unexpandable-row__value").First().Text())
	if price != "" {
		post.Price = price
	}

	images := extractImageURLs(doc)
	if len(images) > 0 {
		post.Images = images
	}

	return post, nil
}

func extractImageURLs(doc *goquery.Document) []string {
	var images []string
	doc.Find("div.kt-base-carousel__slide img.kt-image-block__image").Each(func(i int, s *goquery.Selection) {
		if src, exists := s.Attr("src"); exists {
			images = append(images, src)
		}
	})
	return images
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
