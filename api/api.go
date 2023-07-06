package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func ApiServer() {
	app := fiber.New()

	app.Get("/redirect", func(c *fiber.Ctx) error {
		site := c.Query("site")
		sku := c.Query("sku")

		// Logic to get the site redirect

		return c.Redirect(GetSiteUrl(site, sku))
	})

	app.Listen(":8080")
}

func GetSiteUrl(site string, sku string) string {
	switch site {
	case "footlocker":
		return fmt.Sprintf("https://www.footlocker.it/product/~/%s.html", sku)
	case "adidas":
		return fmt.Sprintf("https://www.adidas.it/hm/%s.html", sku)
	case "bootz":
		return fmt.Sprintf("https://www.boozt.com/%s", sku)
	case "mana":
		return fmt.Sprintf("%s", sku)
	case "jd":
		return fmt.Sprintf("https://www.jdsports.it/product/hm/%s/", sku)
	case "size":
		return fmt.Sprintf("https://www.sizeofficial.it/product/hm/%s/", sku)
	case "footpatrol":
		return fmt.Sprintf("https://www.footpatrol.it/product/hm/%s/", sku)
	case "rina":
		return fmt.Sprintf("https://www.rinascente.it/it/item/%s", sku)
	case "skate":
		return fmt.Sprintf("https://www.skatedeluxe.com/it/hm_p%s", sku)
	case "sotf":
		return fmt.Sprintf("%s", sku)
	case "lando":
		return fmt.Sprintf("https://www.zalando.it/donna/?q=%s", sku)
	default:
		return fmt.Sprintf("https://www.google.com/")
	}
}
