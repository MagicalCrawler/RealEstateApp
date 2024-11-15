package divar

import (
	"context"
	"fmt"
	crawlerModels "github.com/MagicalCrawler/RealEstateApp/models/crawler"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
)

// DivarCrawler implements the Crawler interface for the Divar website
type DivarCrawler struct {
	baseURL    string
	userAgents []string
}

// NewDivarCrawler creates a new instance of DivarCrawler
func NewDivarCrawler() *DivarCrawler {
	return &DivarCrawler{
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

// Crawl fetches posts for a given city
func (c *DivarCrawler) Crawl(ctx context.Context, city crawlerModels.City) ([]crawlerModels.Post, error) {
	pageLimit := 2 // Adjust as needed
	cityURL := fmt.Sprintf("%s/s/%s/real-estate", c.baseURL, city.Slug)
	paginationURLs := c.getPaginationURLs(cityURL, pageLimit)
	var allPosts []crawlerModels.Post

	// Initialize Playwright
	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("could not start Playwright: %w", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		return nil, fmt.Errorf("could not launch browser: %w", err)
	}
	defer browser.Close()

	for _, pageURL := range paginationURLs {
		select {
		case <-ctx.Done():
			return allPosts, ctx.Err()
		default:
		}

		fmt.Printf("Crawling page: %s\n", pageURL)
		page, err := browser.NewPage()
		if err != nil {
			return nil, fmt.Errorf("could not create new page: %w", err)
		}

		err = page.SetExtraHTTPHeaders(map[string]string{
			"User-Agent": c.randomUserAgent(),
			"Referer":    "https://google.com",
		})
		if err != nil {
			return nil, fmt.Errorf("could not set extra headers: %w", err)
		}

		_, err = page.Goto(pageURL, playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateNetworkidle,
			Timeout:   playwright.Float(90000),
		})
		if err != nil {
			log.Printf("Error navigating to %s: %v", pageURL, err)
			page.Close()
			continue
		}

		// Scroll to load all content
		err = c.autoScroll(page)
		if err != nil {
			log.Printf("Error during auto-scrolling: %v", err)
		}

		content, err := page.Content()
		if err != nil {
			page.Close()
			return nil, fmt.Errorf("could not get page content: %w", err)
		}
		page.Close()

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
		if err != nil {
			return nil, fmt.Errorf("failed to parse HTML from %s: %w", pageURL, err)
		}

		postLinks := c.extractPostLinksFromSelection(doc)
		if len(postLinks) == 0 {
			log.Printf("No posts found on %s", pageURL)
			continue
		}

		for _, link := range postLinks {
			select {
			case <-ctx.Done():
				return allPosts, ctx.Err()
			default:
			}

			post, err := c.CrawlPostDetails(ctx, link)
			if err != nil {
				log.Printf("Error crawling post details from %s: %v", link, err)
				continue
			}
			post.City = city
			allPosts = append(allPosts, post)
		}

		// Random sleep to mimic human behavior
		time.Sleep(time.Duration(rand.Intn(7)+3) * time.Second)
	}

	return allPosts, nil
}

// CrawlPostDetails fetches details for a single post
func (c *DivarCrawler) CrawlPostDetails(ctx context.Context, postURL string) (crawlerModels.Post, error) {
	var post crawlerModels.Post

	// Initialize Playwright
	pw, err := playwright.Run()
	if err != nil {
		return post, fmt.Errorf("could not start Playwright: %w", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		return post, fmt.Errorf("could not launch browser: %w", err)
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		return post, fmt.Errorf("could not create new page: %w", err)
	}
	defer page.Close()

	err = page.SetExtraHTTPHeaders(map[string]string{
		"User-Agent": c.randomUserAgent(),
		"Referer":    "https://google.com",
	})
	if err != nil {
		return post, fmt.Errorf("could not set extra headers: %w", err)
	}

	_, err = page.Goto(postURL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
		Timeout:   playwright.Float(90000),
	})
	if err != nil {
		log.Printf("Error navigating to %s: %v", postURL, err)
		return post, err
	}

	content, err := page.Content()
	if err != nil {
		return post, fmt.Errorf("could not get page content: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return post, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Parse post details
	splitURL := strings.Split(postURL, "/")
	post.ID = splitURL[len(splitURL)-1]
	post.Link = postURL
	post.Title = strings.TrimSpace(doc.Find("h1.kt-page-title__title").Text())
	post.Description = strings.TrimSpace(doc.Find("div.post-page__section--padded").Text())

	// Extract additional details
	c.extractPostDetails(doc, &post)

	return post, nil
}

// Helper methods

func (c *DivarCrawler) randomUserAgent() string {
	return c.userAgents[rand.Intn(len(c.userAgents))]
}

func (c *DivarCrawler) autoScroll(page playwright.Page) error {
	maxScrollAttempts := 10
	scrollAttempts := 0
	prevPostCount := 0

	for scrollAttempts < maxScrollAttempts {
		// Scroll to bottom
		_, err := page.Evaluate(`() => {
            window.scrollTo(0, document.body.scrollHeight);
            return document.querySelectorAll('div.kt-post-card__body').length;
        }`)
		if err != nil {
			return fmt.Errorf("error scrolling page: %w", err)
		}

		// Wait for new content
		time.Sleep(2 * time.Second)

		// Check number of posts
		currentPostCount, err := page.Evaluate(`() => {
            return document.querySelectorAll('div.kt-post-card__body').length;
        }`)
		if err != nil {
			return fmt.Errorf("error counting posts: %w", err)
		}

		var postCount int
		switch v := currentPostCount.(type) {
		case float64:
			postCount = int(v)
		case int:
			postCount = v
		case string:
			postCount, _ = strconv.Atoi(v)
		default:
			log.Printf("Unexpected type for post count: %T", currentPostCount)
			break
		}

		if postCount == prevPostCount {
			break
		}

		prevPostCount = postCount
		scrollAttempts++
	}

	return nil
}

func (c *DivarCrawler) getPaginationURLs(base string, limit int) []string {
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

func (c *DivarCrawler) extractPostLinksFromSelection(doc *goquery.Document) []string {
	var postLinks []string
	doc.Find("div.kt-post-card__body").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Parent().Attr("href")
		if exists {
			postLinks = append(postLinks, fmt.Sprintf("%s%s", c.baseURL, link))
		}
	})
	return postLinks
}

func (c *DivarCrawler) extractPostDetails(doc *goquery.Document, post *crawlerModels.Post) {
	isRental := doc.Find("div.kt-base-row:contains('ودیعه')").Length() > 0 ||
		doc.Find("div.kt-base-row:contains('اجاره')").Length() > 0

	if isRental {
		isDailyRental := doc.Find("div:contains('روزانه')").Length() > 0 ||
			doc.Find("div:contains('شب')").Length() > 0
		if isDailyRental {
			rentalMetadata := &crawlerModels.RentalMetadata{}
			doc.Find("div.kt-base-row").Each(func(i int, s *goquery.Selection) {
				title := s.Find("p.kt-unexpandable-row__title").Text()
				value := s.Find("p.kt-unexpandable-row__value").Text()

				switch strings.TrimSpace(title) {
				case "ظرفیت":
					rentalMetadata.Capacity = strings.TrimSpace(value)
				case "روزهای عادی":
					rentalMetadata.NormalDayPrice = strings.TrimSpace(value)
				case "آخر هفته":
					rentalMetadata.WeekendPrice = strings.TrimSpace(value)
				case "تعطیلات و مناسبت‌ها":
					rentalMetadata.HolidayPrice = strings.TrimSpace(value)
				case "هزینهٔ هر نفرِ اضافه":
					rentalMetadata.ExtraPersonCost = strings.TrimSpace(value)
				}
			})
			post.RentalMetadata = rentalMetadata
		} else {
			doc.Find("div.kt-base-row").Each(func(i int, s *goquery.Selection) {
				title := s.Find("p.kt-unexpandable-row__title").Text()
				value := s.Find("p.kt-unexpandable-row__value").Text()

				switch strings.TrimSpace(title) {
				case "ودیعه":
					post.Deposit = strings.TrimSpace(value)
				case "اجارهٔ ماهانه":
					post.MonthlyRent = strings.TrimSpace(value)
				case "قیمت کل":
					post.TotalPrice = strings.TrimSpace(value)
				case "قیمت هر متر":
					post.PricePerSquareMeter = strings.TrimSpace(value)
				case "طبقه":
					post.Floor = strings.TrimSpace(value)
				case "ودیعه و اجاره":
					post.DepositOnRentDesc = strings.TrimSpace(value)
				}
			})
		}
	} else {
		doc.Find("div.kt-base-row").Each(func(i int, s *goquery.Selection) {
			title := s.Find("p.kt-unexpandable-row__title").Text()
			value := s.Find("p.kt-unexpandable-row__value").Text()

			switch strings.TrimSpace(title) {
			case "قیمت کل":
				post.TotalPrice = strings.TrimSpace(value)
			case "قیمت هر متر":
				post.PricePerSquareMeter = strings.TrimSpace(value)
			case "طبقه":
				post.Floor = strings.TrimSpace(value)
			}
		})
	}

	// Extract area, year built, rooms
	doc.Find("thead + tbody tr.kt-group-row__data-row").Each(func(i int, s *goquery.Selection) {
		columns := s.Find("td.kt-group-row-item__value.kt-group-row-item--info-row")

		if columns.Length() >= 1 {
			post.Area = strings.TrimSpace(columns.Eq(0).Text())
		}
		if columns.Length() >= 2 {
			post.YearBuilt = strings.TrimSpace(columns.Eq(1).Text())
		}
		if columns.Length() >= 3 {
			post.Rooms = strings.TrimSpace(columns.Eq(2).Text())
		}
	})

	// Extract features
	var features []string
	doc.Find("table.kt-group-row").Last().Find("tbody tr.kt-group-row__data-row td.kt-group-row-item__value").Each(func(i int, s *goquery.Selection) {
		if !s.HasClass("kt-group-row-item--disabled") && s.HasClass("kt-body--stable") {
			feature := strings.TrimSpace(s.Text())
			if feature != "" {
				features = append(features, feature)
			}
		}
	})
	post.Features = features

	// Extract images
	var images []string
	doc.Find("div.kt-base-carousel__slide img.kt-image-block__image").Each(func(i int, s *goquery.Selection) {
		if src, exists := s.Attr("src"); exists {
			images = append(images, src)
		}
	})
	post.Images = images
}
