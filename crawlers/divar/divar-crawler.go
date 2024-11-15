package divar

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
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
func NewDivarRealEstateCrawler(url string) *DivarRealEstateCrawler {
	return &DivarRealEstateCrawler{
		baseURL: url,
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

// autoScroll scrolls the page to load all content
func (c *DivarRealEstateCrawler) autoScroll(page playwright.Page) error {
	// Maximum scroll attempts
	maxScrollAttempts := 10
	scrollAttempts := 0

	// Previous number of posts to check if new posts are loaded
	prevPostCount := 0

	for scrollAttempts < maxScrollAttempts {
		// Scroll to bottom of the page
		_, err := page.Evaluate(`() => {
            window.scrollTo(0, document.body.scrollHeight);
            return document.querySelectorAll('div.kt-post-card__body').length;
        }`)
		if err != nil {
			return fmt.Errorf("error scrolling page: %w", err)
		}

		// Wait for potential new content to load
		time.Sleep(2 * time.Second)

		// Check number of posts
		currentPostCount, err := page.Evaluate(`() => {
            return document.querySelectorAll('div.kt-post-card__body').length;
        }`)
		if err != nil {
			return fmt.Errorf("error counting posts: %w", err)
		}

		// Safely convert post count
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

		// If no new posts loaded, stop scrolling
		if postCount == prevPostCount {
			break
		}

		prevPostCount = postCount
		scrollAttempts++
	}

	return nil
}

// Modify the Crawl method to handle no posts scenario
func (c *DivarRealEstateCrawler) Crawl(pageLimit int, city crawlers.City) ([]crawlers.Post, error) {
	paginationURLs := getPaginationURLs(c.baseURL, pageLimit)
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
		maxRetries := 3
		for attempt := 1; attempt <= maxRetries; attempt++ {
			fmt.Printf("Crawling page: %s (Attempt %d)\n", pageURL, attempt)

			page, err := browser.NewPage()
			if err != nil {
				return nil, fmt.Errorf("could not create new page: %w", err)
			}

			err = page.SetExtraHTTPHeaders(map[string]string{
				"User-Agent": c.RandomUserAgent(),
				"Referer":    "https://google.com",
			})

			if err != nil {
				return nil, fmt.Errorf("could not set extra headers: %w", err)
			}

			if _, err = page.Goto(pageURL, playwright.PageGotoOptions{
				WaitUntil: playwright.WaitUntilStateNetworkidle,
				Timeout:   playwright.Float(90000),
			}); err != nil {
				log.Printf("retrying due to timeout or navigation error for %s: %v", pageURL, err)
				page.Close()
				time.Sleep(5 * time.Second)
				continue
			}

			// Scroll to load all content
			err = c.autoScroll(page)
			if err != nil {
				log.Printf("error during auto-scrolling: %v", err)
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
				log.Printf("No posts found for city %s on attempt %d, retrying...",
					strings.Split(pageURL, "/")[4], attempt)

				if attempt == maxRetries {
					log.Printf("Failed to find posts after %d attempts for %s", maxRetries, pageURL)
					break
				}

				time.Sleep(time.Duration(attempt*5) * time.Second)
				continue
			}

			for _, link := range postLinks {
				post, err := c.CrawlPostDetails(link)
				if err != nil {
					log.Printf("error crawling post details from %s: %v", link, err)
					continue
				}
				post.City = city
				allPosts = append(allPosts, post)
			}

			time.Sleep(time.Duration(rand.Intn(7)+3) * time.Second)
			break
		}
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

// extractPostLinksFromSelection extracts only post links from the main listing page
func (c *DivarRealEstateCrawler) extractPostLinksFromSelection(doc *goquery.Document) []string {
	var postLinks []string
	doc.Find("div.kt-post-card__body").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Parent().Attr("href")
		if exists {
			postLinks = append(postLinks, fmt.Sprintf("https://divar.ir%s", link))
		}
	})
	return postLinks
}

// CrawlPostDetails extracts each post details and fill-out model data
func (c *DivarRealEstateCrawler) CrawlPostDetails(postURL string) (crawlers.Post, error) {
	var post crawlers.Post
	pw, err := playwright.Run()
	if err != nil {
		return post, fmt.Errorf("could not start Playwright: %v", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		return post, fmt.Errorf("could not launch browser: %v", err)
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		return post, fmt.Errorf("could not create new page: %w", err)
	}
	defer page.Close()

	page.SetExtraHTTPHeaders(map[string]string{
		"User-Agent": c.RandomUserAgent(),
		"Referer":    "https://google.com",
	})

	if _, err = page.Goto(postURL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
		Timeout:   playwright.Float(90000),
	}); err != nil {
		log.Printf("retrying due to timeout or navigation error for %s: %v", postURL, err)
		time.Sleep(5 * time.Second)
	}

	content, err := page.Content()
	if err != nil {
		return post, fmt.Errorf("could not get page content: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return post, fmt.Errorf("failed to parse HTML: %w", err)
	}

	splitURL := strings.Split(postURL, "/")
	post.ID = splitURL[len(splitURL)-1]

	post.Title = strings.TrimSpace(doc.Find("h1.kt-page-title__title").Text())
	post.Description = strings.TrimSpace(doc.Find("div.post-page__section--padded").Text())

	isRental := doc.Find("div.kt-base-row:contains('ودیعه')").Length() > 0 || doc.Find("div.kt-base-row:contains('اجاره')").Length() > 0

	if isRental {
		isDailyRental := doc.Find("div:contains('روزانه')").Length() > 0 || doc.Find("div:contains('شب')").Length() > 0
		if isDailyRental {
			rentalMetadata := &crawlers.RentalMetadata{}
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

	var images []string
	doc.Find("div.kt-base-carousel__slide img.kt-image-block__image").Each(func(i int, s *goquery.Selection) {
		if src, exists := s.Attr("src"); exists {
			images = append(images, src)
		}
	})
	post.Images = images

	post.Link = postURL

	return post, nil
}
