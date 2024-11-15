package db

import (
	"fmt"
	"github.com/MagicalCrawler/RealEstateApp/types"
	"log"
	"strconv"

	"github.com/MagicalCrawler/RealEstateApp/models"
	"github.com/MagicalCrawler/RealEstateApp/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection() *gorm.DB {
	host := utils.GetConfig("POSTGRES_HOST")
	user := utils.GetConfig("POSTGRES_USER")
	password := utils.GetConfig("POSTGRES_PASSWORD")
	dbname := utils.GetConfig("POSTGRES_DB_NAME")
	port := utils.GetConfig("POSTGRES_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tehran", host, user, password, dbname, port)
	datab, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		panic("Error connecting to database")
	}

	datab.AutoMigrate(&models.User{})
	datab.AutoMigrate(&models.Post{}, &models.PostHistory{}, &models.Bookmark{})
	// Run auto-migrations for FilterItem and WatchList models
	if err := datab.AutoMigrate(&models.FilterItem{}, &models.WatchList{}); err != nil {
		return nil
	}
	seedSuperAdminUser(datab)
	postsSeeds(datab)
	return datab
}

func seedSuperAdminUser(datab *gorm.DB) {
	superAdminTelegramId, _ := strconv.ParseUint(utils.GetConfig("SUPER_ADMIN"), 10, 64)
	superAdminUser := models.User{
		TelegramID: superAdminTelegramId,
		Role:       models.SUPER_ADMIN,
	}
	if err := datab.FirstOrCreate(&superAdminUser, models.User{TelegramID: superAdminTelegramId}).Error; err != nil {
		fmt.Printf("Could not seed super-admin user (%v): %v", superAdminTelegramId, err)
		panic("Could not seed super-admin user")
	}
}

func postsSeeds(datab *gorm.DB) {
	datab.Unscoped().Where("1 = 1").Delete(&models.PostHistory{})
	datab.Unscoped().Where("1 = 1").Delete(&models.Post{})
	postRepository := NewPostRepository(datab)
	crawlInfo := models.CrawlHistory{}
	if err := datab.Create(&crawlInfo).Error; err != nil {
		fmt.Printf("Could not seed crawl history (%v): %v", crawlInfo, err)
	}
	post1 := models.Post{
		UniqueCode: "wZ0kfXs_",
	}

	if err := postRepository.PostIsExist(post1); err == false {
		datab.Create(&post1)
	} else {
		log.Println("Post already exists")
	}

	postHistory1 := models.PostHistory{
		Post:           post1,
		Title:          "شاهین، ۶۲متر، ۶ساله (تمامی اطلاعات و عکسها واقعی)",
		PostURL:        "https://divar.ir/v/%D8%B4%D8%A7%D9%87%DB%8C%D9%86-%DB%B6%DB%B2%D9%85%D8%AA%D8%B1-%DB%B6%D8%B3%D8%A7%D9%84%D9%87-%D8%AA%D9%85%D8%A7%D9%85%DB%8C-%D8%A7%D8%B7%D9%84%D8%A7%D8%B9%D8%A7%D8%AA-%D9%88-%D8%B9%DA%A9%D8%B3%D9%87%D8%A7-%D9%88%D8%A7%D9%82%D8%B9%DB%8C/wZ0kfXs_",
		Price:          6900000000,
		City:           "تهران",
		Neighborhood:   "پونک",
		Area:           63,
		BedroomNum:     1,
		Age:            5,
		FloorsNum:      6,
		BuyMode:        types.Shopping,
		Building:       types.Apartment,
		HasStorage:     true,
		HsaParking:     true,
		HasElevator:    true,
		ImageURL:       "https://s100.divarcdn.com/static/photo/neda/post/YjOoz1jbtpZUd2h7EqFl_Q/3078dfe4-3251-43b2-9404-04baecfa031b.jpg",
		Description:    "توضیحات\n❌❌❌ فیلم ، عکس و مشخصات ۱۰۰٪ واقعی ❌❌❌\n\n‼️‼️‼️یکی از جذابترین یکخوابهای منطقه ‼️‼️‼️\n\n☑️ نور و نقشه سوپر استثنایی\n\n☑️ داخل واحد در حد کلیدنخورده (کلا یکسال ساکن داشته)\n\n☑️ دسترسی عالی به اتوبانها، مراکز خرید، بیمارستان و هر آنچه که برای یه زندگی آروم نیازمندش هستین\n\n☑️ فایل کاملا شخصی ، بازدید آزاد ( ۸صبح تا ۱۰شب )\n\n☑️ قابلیت دریافت ۱ میلیارد وام\n\n☑️ قابلیت ۹۰۰ میلیون رهن کامل ( کمتر از یک هفته )\n\n☑️ نقدینگی لازم برای شما »»»»»» ۵ میلیارد !!!!!!\n\n❌❌❌❌ مالک فروشنده قطعی و واقعی ❌❌❌❌\n\n\nفایلهای مشابه؛\n\n۶۴متر ، ۱۲ ساله ( بلوار اباذر )\n۵۸ متر ، ۵ ساله ( بلوار فردوس )\n۶۰ متر ، ۱۰ ساله ( جنت آباد جنوبی )\n۶۱ متر ، ۹ ساله ( ستاری ، مهستان )\n۶۰ متر ، نوساز ( باکس پونک )\n۵۵ متر ، ۱۴ ساله ( کاشانی ، آلاله )\n\n\n✍️ برای دریافت اطلاعات بیشتر لطفاً تماس بگیرید\n\n❇️ املاک بزرگ باراد\nرمضانی",
		CrawlHistory:   crawlInfo,
		CrawlHistoryID: crawlInfo.ID,
	}

	if err := postRepository.PostHistoryIsExist(postHistory1); err == false {
		datab.Create(&postHistory1)
	} else {
		log.Println("PostHistory already exists")
	}

	post2 := models.Post{
		UniqueCode: "wZ0o_N3w",
	}
	if err := postRepository.PostIsExist(post2); err == false {
		datab.Create(&post2)
	} else {
		log.Println("Post already exists")
	}

	postHistory2 := models.PostHistory{
		Post:           post2,
		Title:          "آپارتمان 62 متری طالقانی نزدیک پارک هنرمندان",
		PostID:         post2.ID,
		PostURL:        "https://divar.ir/v/%D8%A2%D9%BE%D8%A7%D8%B1%D8%AA%D9%85%D8%A7%D9%86-62-%D9%85%D8%AA%D8%B1%DB%8C-%D8%B7%D8%A7%D9%84%D9%82%D8%A7%D9%86%DB%8C-%D9%86%D8%B2%D8%AF%DB%8C%DA%A9-%D9%BE%D8%A7%D8%B1%DA%A9-%D9%87%D9%86%D8%B1%D9%85%D9%86%D8%AF%D8%A7%D9%86/wZ0o_N3w",
		Price:          7000000000,
		City:           "تهران",
		Neighborhood:   "ایرانشهر",
		Area:           62,
		BedroomNum:     1,
		BuyMode:        types.Shopping,
		Building:       types.Apartment,
		Age:            9,
		FloorsNum:      5,
		HasStorage:     true,
		HsaParking:     true,
		HasElevator:    true,
		ImageURL:       "https://s100.divarcdn.com/static/photo/neda/post/xMScyIZiNPYUP9L_fAJrpw/5b5d3bc8-d2bd-44bc-84e4-8a512ef2387e.jpg",
		Description:    "توضیحات\nآپارتمان 62 متری\nیک خواب دارای کمد دیواری\nآسانسور\nپارکینگ\nانباری\nبالکن\nسرویس ایرانی و فرنگی\nآشپزخانه MDF\nایفون تصویری\nدرب ریموت کنترل\nنما سنگ\nساختمان جنوبی\nنورگیر و رو به آفتاب\nدسترسی سریع به مترو طالقانی پارک هنرمندان\nتخفیف هنگام قرارداد\nمشاور املاک\nخانم قمری\n\n\n",
		CrawlHistory:   crawlInfo,
		CrawlHistoryID: crawlInfo.ID,
	}
	if err := postRepository.PostHistoryIsExist(postHistory2); err == false {
		datab.Create(&postHistory2)
	}

	post3 := models.Post{UniqueCode: "wZ0YoPIi"}
	if err := postRepository.PostIsExist(post3); err == false {
		datab.Create(&post3)
	} else {
		log.Printf("Post already exists")
	}
	postHistory3 := models.PostHistory{
		Post:           post3,
		PostID:         post3.ID,
		Title:          "اپارتمان ۱۱۲متری فول امکانت تک واحدی وحدت اسلامی",
		PostURL:        "https://divar.ir/v/%D8%A7%D9%BE%D8%A7%D8%B1%D8%AA%D9%85%D8%A7%D9%86-%DB%B1%DB%B1%DB%B2%D9%85%D8%AA%D8%B1%DB%8C-%D9%81%D9%88%D9%84-%D8%A7%D9%85%DA%A9%D8%A7%D9%86%D8%AA-%D8%AA%DA%A9-%D9%88%D8%A7%D8%AD%D8%AF%DB%8C-%D9%88%D8%AD%D8%AF%D8%AA-%D8%A7%D8%B3%D9%84%D8%A7%D9%85%DB%8C/wZ0YoPIi",
		Price:          7550000000,
		City:           "تهران",
		Neighborhood:   "امیریه",
		Area:           112,
		BedroomNum:     2,
		BuyMode:        types.Shopping,
		Building:       types.Apartment,
		Age:            5,
		FloorsNum:      6,
		HasStorage:     false,
		HsaParking:     true,
		HasElevator:    true,
		ImageURL:       "https://s100.divarcdn.com/static/photo/neda/post/ac9wb68uPe6R7B374NLECw/72a83800-d2b0-45c6-8f4b-aba2be95853d.jpg",
		Description:    "توضیحات\nکل طبقه ۶واحد تک واحد\nخوش نقشه بدون پرتی\n۳بهر نور همسایهای عالی محله دنج ارام\nاین ملک شخصی میباشد\nداری مستجر\nکابینت های گلاس\nیکی از بهترین کوچه شاپور می باشد",
		CrawlHistory:   crawlInfo,
		CrawlHistoryID: crawlInfo.ID,
	}
	if err := postRepository.PostHistoryIsExist(postHistory3); err == false {
		datab.Create(&postHistory3)
	} else {
		log.Printf("PostHistory already exists")
	}

	post4 := models.Post{UniqueCode: "wZQUyrBv"}
	if err := postRepository.PostIsExist(post4); err == false {
		datab.Create(&post4)
	} else {
		log.Printf("Post already exists")
	}
	postHistory4 := models.PostHistory{
		Post:           post4,
		PostID:         post4.ID,
		Title:          "۲۰۰متر/۳جهت ویو/۲پارکینگ/تکواحدی/کم سن",
		PostURL:        "https://divar.ir/v/%DB%B2%DB%B0%DB%B0%D9%85%D8%AA%D8%B1-%DB%B3%D8%AC%D9%87%D8%AA-%D9%88%DB%8C%D9%88-%DB%B2%D9%BE%D8%A7%D8%B1%DA%A9%DB%8C%D9%86%DA%AF-%D8%AA%DA%A9%D9%88%D8%A7%D8%AD%D8%AF%DB%8C-%DA%A9%D9%85-%D8%B3%D9%86/wZQUyrBv",
		Price:          28000000000,
		City:           "تهران",
		Neighborhood:   "پونک",
		Area:           200,
		BedroomNum:     3,
		BuyMode:        types.Shopping,
		Building:       types.Apartment,
		Age:            7,
		FloorsNum:      6,
		HasStorage:     true,
		HsaParking:     true,
		HasElevator:    true,
		ImageURL:       "https://s100.divarcdn.com/static/photo/afra/post/2W3iDULmlVViIh1KfSMHnQ/932cab09-f723-460c-9f01-5f9032d632b7.jpg",
		Description:    "با درود فراوان\n\nملک ۲کله و ۳نبش و شخصی ساز است\nواقع در یکی از لوکیشن های خوب منطقه۵\n\nقابل توجه مشاورین محترم ملک به هیچ وجه کارشناسی ندارد ، از آوردن همکاران خود به عنوان مشتری جدا خودداری کنید چون قطع همکاری میشه ، برای دریافت فیلم و عکسهای واحد جهت کارشناسی و معرفی به مشتری ، در واتساپ پیام بدهید با ارسال نام خودتون و املاکتون.\nکمیسیون فروش این ملک ۱٪ تقدیم میگردد و تخفیف باتوجه به شرایط پرداختی خریدار محترم منظور می‌گردد.\n\nبا احترام گلنام",
		CrawlHistory:   crawlInfo,
		CrawlHistoryID: crawlInfo.ID,
	}
	if err := postRepository.PostHistoryIsExist(postHistory4); err == false {
		datab.Create(&postHistory4)
	} else {
		log.Printf("PostHistory already exists")
	}

	post5 := models.Post{UniqueCode: "wZ0ATY3W"}
	if err := postRepository.PostIsExist(post5); err == false {
		datab.Create(&post5)
	} else {
		log.Printf("Post already exists")
	}
	postHistory5 := models.PostHistory{
		Post:           post5,
		PostID:         post5.ID,
		Title:          "آپارتمان ۱۰۰ متری ، طبقه دوم",
		PostURL:        "https://divar.ir/v/%D8%A2%D9%BE%D8%A7%D8%B1%D8%AA%D9%85%D8%A7%D9%86-%DB%B1%DB%B0%DB%B0-%D9%85%D8%AA%D8%B1%DB%8C-%D8%B7%D8%A8%D9%82%D9%87-%D8%AF%D9%88%D9%85/wZ0ATY3W",
		Price:          8000000000,
		City:           "تهران",
		Neighborhood:   "بهارستان",
		Area:           100,
		BedroomNum:     2,
		BuyMode:        types.Shopping,
		Building:       types.Apartment,
		Age:            0,
		FloorsNum:      2,
		HasStorage:     false,
		HsaParking:     false,
		HasElevator:    false,
		ImageURL:       "https://s100.divarcdn.com/static/photo/neda/post/i8c7Vnn3BsKLu6dIYq9m5g/56a67e14-59be-4804-bd1b-49a175c31c47.jpg",
		Description:    "سند مسکونی موقعیت اداری\nدر دو طبقه ، دو واحدی\nنوساز ، کلید نخورده\nواقع در بالای فروشگاه\nآگهی توسط مالک درج شده ، فروشنده بدون واسطه",
		CrawlHistory:   crawlInfo,
		CrawlHistoryID: crawlInfo.ID,
	}
	if err := postRepository.PostHistoryIsExist(postHistory5); err == false {
		datab.Create(&postHistory5)
	} else {
		log.Printf("PostHistory already exists")
	}

	post6 := models.Post{UniqueCode: "wZyQVpfx"}
	if err := postRepository.PostIsExist(post6); err == false {
		datab.Create(&post6)
	} else {
		log.Printf("Post already exists")
	}
	postHistory6 := models.PostHistory{
		Post:           post6,
		PostID:         post6.ID,
		Title:          "آپارتمان ۷۰متری ،دوخوابه",
		PostURL:        "https://divar.ir/v/%D8%A2%D9%BE%D8%A7%D8%B1%D8%AA%D9%85%D8%A7%D9%86-%DB%B7%DB%B0%D9%85%D8%AA%D8%B1%DB%8C-%D8%AF%D9%88%D8%AE%D9%88%D8%A7%D8%A8%D9%87/wZyQVpfx",
		Price:          1850000000,
		City:           "تهران",
		Neighborhood:   "بهارستان",
		Area:           70,
		BedroomNum:     2,
		BuyMode:        types.Shopping,
		Building:       types.Apartment,
		Age:            13,
		FloorsNum:      5,
		HasStorage:     true,
		HsaParking:     false,
		HasElevator:    true,
		ImageURL:       "https://s100.divarcdn.com/static/photo/neda/post/7Ih77OjS2R4oNOuhL_vMHA/4199785b-cd97-4003-a162-8fdbe39f24a9.jpg",
		Description:    "آپارتمان کاملا نو سازی شده ،کاغذ دیواری ،کابینت بندی ،هود وسینک، کم دیواری ،کاشی وسنگ و روشویی توالت ،کولرآبی نو، در تراس داخل آشپزخانه دوجداره بزرگ ویکسره می‌باشد،سندمنگوله دار ،\nآدرس شهرستان بهارستان ،شهرگلستان",
		CrawlHistory:   crawlInfo,
		CrawlHistoryID: crawlInfo.ID,
	}
	if err := postRepository.PostHistoryIsExist(postHistory6); err == false {
		datab.Create(&postHistory6)
	} else {
		log.Printf("PostHistory already exists")
	}
}
