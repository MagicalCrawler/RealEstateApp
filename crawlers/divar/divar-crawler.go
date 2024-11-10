package divar

import (
	"Bootcamp/Projects/RealEstateApp/crawlers"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type DivarRealEstateCrawler struct {
	baseURL    string
	httpClient *http.Client
	userAgent  string
}

func NewDivarRealEstateCrawler() *DivarRealEstateCrawler {

	return &DivarRealEstateCrawler{
		baseURL: "https://divar.ir/s/tehran/real-estate",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36",
	}
}

// Crawl fetches and extracts posts from the main page
func (c *DivarRealEstateCrawler) Crawl(ctx context.Context) ([]crawlers.Post, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var posts []crawlers.Post
	doc.Find("div.kt-post-card__body").Each(func(i int, s *goquery.Selection) {
		post := c.extractPostFromSelection(s)
		posts = append(posts, post)
	})

	return posts, nil

}

// CrawlPostDetails fetches additional details for a specific post
func (c *DivarRealEstateCrawler) CrawlPostDetails(ctx context.Context, post crawlers.Post) (crawlers.Post, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", post.Link, nil)
	if err != nil {
		return post, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	resp, err := c.httpClient.Do(req)
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

	return post, nil
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

	image, exists := s.Parent().Find("img").Attr("src")
	if !exists {
		image = "No Image"
	}

	return crawlers.Post{
		Title: title,
		Price: price,
		Link:  fmt.Sprintf("https://divar.ir%s", link),
		Image: image,
	}
}
