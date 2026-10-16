package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/TheAldriguzDev/Bonsai_solutions/packages/discord"
	"github.com/TheAldriguzDev/Bonsai_solutions/packages/log"
)

var location = "Database"

var mutex = &sync.Mutex{}

func StartDataBase(fileName string, products []Product) error {

	if len(products) == 0 {
		return fmt.Errorf("No products has been found!")
	} else {
		err := CheckUpdates(fileName, products)
		if err != nil {
			return err
		}

		err = RemoveProducts(fileName, products)
		if err != nil {
			return err
		}

		return nil
	}
}

// Simple function to read the file
func ReadFile(fileName string) (*ProductList, error) {
	mutex.Lock()
	defer mutex.Unlock()

	file, err := ioutil.ReadFile(fileName)

	if err != nil {
		return nil, err
	}

	productList := &ProductList{}
	err = json.Unmarshal(file, productList)

	if err != nil {
		return nil, err
	}

	return productList, nil
}

// Simple function to write the new file
func WriteFile(fileName string, productList *ProductList) error {
	mutex.Lock()
	defer mutex.Unlock()

	file, err := json.MarshalIndent(productList, "", "	")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fileName, file, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Function to check if a product exist in the product list
func ProductExist(productList *ProductList, product Product) bool {
	for _, p := range productList.Products {
		if p.Sku == product.Sku {
			return true
		}
	}

	return false
}

func CheckUpdates(fileName string, products []Product) error {
	productList, err := ReadFile(fileName)
	if err != nil {
		// If the file does not exist we create a new list
		productList = &ProductList{}
	}

	for _, product := range products {
		if ProductExist(productList, product) {
			// Product exist, we check for changes
			for i, p := range productList.Products {
				if p.Sku == product.Sku {
					if p.Name != product.Name || p.Price != product.Price || p.Image != product.Image || !CheckSize(p.Sizes, product.Sizes) {
						log.Info("Product changes found!", location)
						productList.Products[i] = product
						// Add log webhook, type restock
						SendLog(product)
					}
					break
				}
			}
		} else {
			// The product is new, we add to the list
			log.Info("New product has been found!", location)
			// Add log webhook, type new add
			SendLog(product)
			productList.Products = append(productList.Products, product)
		}
	}

	err = WriteFile(fileName, productList)
	if err != nil {
		return err
	}
	return nil
}

// Function to check if sizes are changed
func CheckSize(sizes1, sizes2 []Size) bool {
	if len(sizes1) != len(sizes2) {
		return false
	}

	for i, size := range sizes1 {
		if size.Size != sizes2[i].Size || size.Value != sizes2[i].Value {
			return false
		}
	}

	return true
}

func RemoveProducts(fileName string, products []Product) error {
	productList, err := ReadFile(fileName)
	if err != nil {
		return err
	}

	var updatedProducts []Product

	for _, product := range productList.Products {
		if ProductExistInList(product, products) {
			updatedProducts = append(updatedProducts, product)
		} else {
			log.General("Product not found, removing...", location)
		}
	}

	productList.Products = updatedProducts

	err = WriteFile(fileName, productList)
	if err != nil {
		return err
	}

	return nil
}

func ProductExistInList(product Product, products []Product) bool {
	for _, p := range products {
		if p.Sku == product.Sku {
			return true
		}
	}
	return false
}

func SendLog(product Product) {
	switch product.Site {
	case "Adidas":
		AdidasLog(product)
	case "Footlocker Restock":
		//FootlockerLog(product)
		break
	case "Footlocker New":
		//FtlNewLog(product)
		break
	case "Zalando":
		ZalandoLog(product)
	case "Bootz":
		BootzLog(product)
	case "Rinascente":
		RinaLog(product)
	case "Mana":
		ManaLog(product)
	case "SkateDeluxe":
		SkateLog(product)
	case "JD Restock":
		JDLog(product)
	case "JD New":
		JDNewAddLog(product)
	case "SIZE New":
		SizeNew(product)
	case "SIZE Restock":
		SizeRestock(product)
	case "FP New":
		FpNew(product)
	case "FP Restock":
		FpRestock(product)
	case "Sotf":
		SotfLog(product)
	case "Asos":
		AsosLog(product)
	case "Namshi":
		NamshiLog(product)
	}
}

func AdidasLog(product Product) {
	if product.Color == "" {
		product.Color = "None"
	}
	var sizeMessage string = ""
	for _, size := range product.Sizes {
		sizeMessage = sizeMessage + fmt.Sprintf("%s [%s]\n", size.Size, size.Value)
	}
	whData := discord.DiscordParams{
		Embeds: []discord.Embed{
			{
				Author: discord.Author{
					Name:    "Hypemedia Monitors",
					IconURL: "https://media.discordapp.net/attachments/906886120502075403/1095465707690524753/hypemedialov2.jpg?width=1170&height=1170",
				},
				Title: product.Name,
				Color: 16776975,
				Url:   fmt.Sprintf("https://www.adidas.it/hm/%s.html", product.Sku),
				Footer: discord.Footer{
					Text: "Hypemedia Monitors",
				},
				Fields: []discord.Field{
					{
						Name:   "Sku",
						Value:  product.Sku,
						Inline: true,
					},
					{
						Name:   "Name",
						Value:  product.Name,
						Inline: true,
					},
					{
						Name:   "Color",
						Value:  product.Color,
						Inline: true,
					},
					{
						Name:   "Price",
						Value:  product.Price,
						Inline: true,
					},
					{
						Name:   "Sizes",
						Value:  sizeMessage,
						Inline: false,
					},
				},
				Thumbnail: discord.Thumbnail{
					URL: product.Image,
				},
			},
		},
	}

	webhookURL := os.Getenv("BONSAI_WEBHOOK_ADIDAS_LOG")
	if webhookURL == "" {
		log.Error("Missing BONSAI_WEBHOOK_ADIDAS_LOG", location)
		return
	}
	err := discord.PostWebhook(webhookURL, whData)
	if err != nil {
		log.Error("An error occured while sending the webhook!", location)
		fmt.Println(whData)
		time.Sleep(5 * time.Second)
		AdidasLog(product)
	} else {
		log.Success("Webhook sent", location)
		return
	}
}

func FootlockerLog(product Product) {
	var sizeMessage string = ""
	for _, size := range product.Sizes {
		sizeMessage = sizeMessage + fmt.Sprintf("%s \n", size.Size)
	}
	whData := discord.DiscordParams{
		Embeds: []discord.Embed{
			{
				Author: discord.Author{
					Name:    "Hypemedia Monitors",
					IconURL: "https://media.discordapp.net/attachments/906886120502075403/1095465707690524753/hypemedialov2.jpg?width=1170&height=1170",
				},
				Title: product.Name,
				Color: 16776975,
				Url:   fmt.Sprintf("https://www.footlocker.it/product/~/%s.html", product.Sku),
				Footer: discord.Footer{
					Text: "Hypemedia Monitors",
				},
				Fields: []discord.Field{
					{
						Name:   "Type",
						Value:  "Restock",
						Inline: true,
					},
					{
						Name:   "Sku",
						Value:  product.Sku,
						Inline: true,
					},
					{
						Name:   "Price",
						Value:  product.Price,
						Inline: true,
					},
					{
						Name:   "Color",
						Value:  product.Color,
						Inline: true,
					},
					{
						Name:   "Sizes",
						Value:  sizeMessage,
						Inline: false,
					},
				},
				Thumbnail: discord.Thumbnail{
					URL: product.Image,
				},
			},
		},
	}

	webhookURL := os.Getenv("BONSAI_WEBHOOK_FOOTLOCKER_LOG")
	if webhookURL == "" {
		log.Error("Missing BONSAI_WEBHOOK_FOOTLOCKER_LOG", location)
		return
	}
	err := discord.PostWebhook(webhookURL, whData)
	if err != nil {
		log.Error("An error occured while sending the webhook!", location)
		fmt.Println(whData)
		time.Sleep(5 * time.Second)
		AdidasLog(product)
	} else {
		log.Success("Webhook sent", location)
		return
	}
}

func FtlNewLog(product Product) {
	whData := discord.DiscordParams{
		Embeds: []discord.Embed{
			{
				Author: discord.Author{
					Name:    "Hypemedia Monitors",
					IconURL: "https://media.discordapp.net/attachments/906886120502075403/1095465707690524753/hypemedialov2.jpg?width=1170&height=1170",
				},
				Title: product.Name,
				Color: 16776975,
				Url:   fmt.Sprintf("https://www.footlocker.it/product/~/%s.html", product.Sku),
				Footer: discord.Footer{
					Text: "Hypemedia Monitors",
				},
				Fields: []discord.Field{
					{
						Name:   "Type",
						Value:  "New Add",
						Inline: true,
					},
					{
						Name:   "Sku",
						Value:  product.Sku,
						Inline: true,
					},
					{
						Name:   "Price",
						Value:  product.Price,
						Inline: true,
					},
				},
				Thumbnail: discord.Thumbnail{
					URL: product.Image,
				},
			},
		},
	}

	webhookURL := os.Getenv("BONSAI_WEBHOOK_FTL_NEW_LOG")
	if webhookURL == "" {
		log.Error("Missing BONSAI_WEBHOOK_FTL_NEW_LOG", location)
		return
	}
	err := discord.PostWebhook(webhookURL, whData)
	if err != nil {
		log.Error("An error occured while sending the webhook!", location)
		fmt.Println(whData)
		time.Sleep(5 * time.Second)
		AdidasLog(product)
	} else {
		log.Success("Webhook sent", location)
		return
	}
}

func ZalandoLog(product Product) {
	var sizeMessage string = ""
	for _, size := range product.Sizes {
		sizeMessage = sizeMessage + fmt.Sprintf("%s [%s]\n", size.Size, size.Value)
	}
	whData := discord.DiscordParams{
		Embeds: []discord.Embed{
			{
				Author: discord.Author{
					Name:    "Hypemedia Monitors",
					IconURL: "https://media.discordapp.net/attachments/906886120502075403/1095465707690524753/hypemedialov2.jpg?width=1170&height=1170",
				},
				Title: product.Name,
				Color: 16776975,
				Url:   fmt.Sprintf("https://www.zalando.it/donna/?q=%s", product.Sku),
				Footer: discord.Footer{
					Text: "Hypemedia Monitors",
				},
				Fields: []discord.Field{
					{
						Name:   "Sku",
						Value:  product.Sku,
						Inline: true,
					},
					{
						Name:   "Name",
						Value:  product.Name,
						Inline: true,
					},
					{
						Name:   "Color",
						Value:  product.Color,
						Inline: true,
					},
					{
						Name:   "Sizes",
						Value:  sizeMessage,
						Inline: false,
					},
				},
				Thumbnail: discord.Thumbnail{
					URL: product.Image,
				},
			},
		},
	}

	webhookURL := os.Getenv("BONSAI_WEBHOOK_ZALANDO_LOG")
	if webhookURL == "" {
		log.Error("Missing BONSAI_WEBHOOK_ZALANDO_LOG", location)
		return
	}
	err := discord.PostWebhook(webhookURL, whData)
	if err != nil {
		log.Error("An error occured while sending the webhook!", location)
		fmt.Println(whData)
		time.Sleep(5 * time.Second)
		AdidasLog(product)
	} else {
		log.Success("Webhook sent", location)
		return
	}
}

func BootzLog(product Product) {
	var sizeMessage string = ""
	for _, size := range product.Sizes {
		sizeMessage = sizeMessage + fmt.Sprintf("%s\n", size.Size)
	}
	whData := discord.DiscordParams{
		Embeds: []discord.Embed{
			{
				Author: discord.Author{
					Name:    "Hypemedia Monitors",
					IconURL: "https://media.discordapp.net/attachments/906886120502075403/1095465707690524753/hypemedialov2.jpg?width=1170&height=1170",
				},
				Title: product.Name,
				Color: 16776975,
				Url:   fmt.Sprintf("https://www.boozt.com/%s", product.Sku),
				Footer: discord.Footer{
					Text: "Hypemedia Monitors",
				},
				Fields: []discord.Field{
					{
						Name:   "Name",
						Value:  product.Name,
						Inline: true,
					},
					{
						Name:   "Color",
						Value:  product.Color,
						Inline: true,
					},
					{
						Name:   "Price",
						Value:  product.Price,
						Inline: true,
					},
					{
						Name:   "Sizes",
						Value:  sizeMessage,
						Inline: false,
					},
				},
				Thumbnail: discord.Thumbnail{
					URL: product.Image,
				},
			},
		},
	}

	webhookURL := os.Getenv("BONSAI_WEBHOOK_BOOTZ_LOG")
	if webhookURL == "" {
		log.Error("Missing BONSAI_WEBHOOK_BOOTZ_LOG", location)
		return
	}
	err := discord.PostWebhook(webhookURL, whData)
	if err != nil {
		log.Error("An error occured while sending the webhook!", location)
		fmt.Println(whData)
		time.Sleep(5 * time.Second)
		AdidasLog(product)
	} else {
		log.Success("Webhook sent", location)
		return
	}
}

func RinaLog(product Product) {
	var sizeMessage string = ""
	for _, size := range product.Sizes {
		sizeMessage = sizeMessage + fmt.Sprintf("%s\n", size.Size)
	}
	whData := discord.DiscordParams{
		Embeds: []discord.Embed{
			{
				Author: discord.Author{
					Name:    "Hypemedia Monitors",
					IconURL: "https://media.discordapp.net/attachments/906886120502075403/1095465707690524753/hypemedialov2.jpg?width=1170&height=1170",
				},
				Title: product.Name,
				Color: 16776975,
				Url:   fmt.Sprintf("https://www.rinascente.it/it/item/%s", product.Sku),
				Footer: discord.Footer{
					Text: "Hypemedia Monitors",
				},
				Fields: []discord.Field{
					{
						Name:   "Name",
						Value:  product.Name,
						Inline: true,
					},
					{
						Name:   "Color",
						Value:  product.Color,
						Inline: true,
					},
					{
						Name:   "Price",
						Value:  product.Price,
						Inline: true,
					},
					{
						Name:   "Sizes",
						Value:  sizeMessage,
						Inline: false,
					},
				},
				Thumbnail: discord.Thumbnail{
					URL: product.Image,
				},
			},
		},
	}

	webhookURL := os.Getenv("BONSAI_WEBHOOK_RINA_LOG")
	if webhookURL == "" {
		log.Error("Missing BONSAI_WEBHOOK_RINA_LOG", location)
		return
	}
	err := discord.PostWebhook(webhookURL, whData)
	if err != nil {
		log.Error("An error occured while sending the webhook!", location)
		fmt.Println(whData)
		time.Sleep(5 * time.Second)
		AdidasLog(product)
	} else {
		log.Success("Webhook sent", location)
		return
	}
}

func ManaLog(product Product) {
	var sizeMessage string = ""
	for _, size := range product.Sizes {
		sizeMessage = sizeMessage + fmt.Sprintf("%s\n", size.Size)
	}
	whData := discord.DiscordParams{
		Embeds: []discord.Embed{
			{
				Author: discord.Author{
					Name:    "Hypemedia Monitors",
					IconURL: "https://media.discordapp.net/attachments/906886120502075403/1095465707690524753/hypemedialov2.jpg?width=1170&height=1170",
				},
				Title: product.Name,
				Color: 16776975,
				Url:   fmt.Sprintf("%s", product.Sku),
				Footer: discord.Footer{
					Text: "Hypemedia Monitors",
				},
				Fields: []discord.Field{
					{
						Name:   "Name",
						Value:  product.Name,
						Inline: true,
					},
					{
						Name:   "Color",
						Value:  product.Color,
						Inline: true,
					},
					{
						Name:   "Price",
						Value:  product.Price,
						Inline: true,
					},
					{
						Name:   "Sizes",
						Value:  sizeMessage,
						Inline: false,
					},
				},
				Thumbnail: discord.Thumbnail{
					URL: product.Image,
				},
			},
		},
	}

	webhookURL := os.Getenv("BONSAI_WEBHOOK_MANA_LOG")
	if webhookURL == "" {
		log.Error("Missing BONSAI_WEBHOOK_MANA_LOG", location)
		return
	}
	err := discord.PostWebhook(webhookURL, whData)
	if err != nil {
		log.Error("An error occured while sending the webhook!", location)
		fmt.Println(whData)
		time.Sleep(5 * time.Second)
		AdidasLog(product)
	} else {
		log.Success("Webhook sent", location)
		return
	}
}

func SkateLog(product Product) {
	var sizeMessage string = ""
	for _, size := range product.Sizes {
		sizeMessage = sizeMessage + fmt.Sprintf("%s\n", size.Size)
	}
	whData := discord.DiscordParams{
		Embeds: []discord.Embed{
			{
				Author: discord.Author{
					Name:    "Hypemedia Monitors",
					IconURL: "https://media.discordapp.net/attachments/906886120502075403/1095465707690524753/hypemedialov2.jpg?width=1170&height=1170",
				},
				Title: product.Name,
				Color: 16776975,
				Url:   fmt.Sprintf("https://www.skatedeluxe.com/it/hm_p%s", product.Sku),
				Footer: discord.Footer{
					Text: "Hypemedia Monitors",
				},
				Fields: []discord.Field{
					{
						Name:   "Name",
						Value:  product.Name,
						Inline: true,
					},
					{
						Name:   "Sku",
						Value:  product.Sku,
						Inline: true,
					},
					{
						Name:   "Sizes",
						Value:  sizeMessage,
						Inline: false,
					},
				},
				Thumbnail: discord.Thumbnail{
					URL: product.Image,
				},
			},
		},
	}

	webhookURL := os.Getenv("BONSAI_WEBHOOK_SKATE_LOG")
	if webhookURL == "" {
		log.Error("Missing BONSAI_WEBHOOK_SKATE_LOG", location)
		return
	}
	err := discord.PostWebhook(webhookURL, whData)
	if err != nil {
		log.Error("An error occured while sending the webhook!", location)
		fmt.Println(whData)
		time.Sleep(5 * time.Second)
		AdidasLog(product)
	} else {
		log.Success("Webhook sent", location)
		return
	}
}

func JDLog(product Product) {
	var sizeMessage string = ""
	var pidsMessage string = ""
	for _, size := range product.Sizes {
		sizeMessage = sizeMessage + fmt.Sprintf("%s\n", size.Size)
		pidsMessage = pidsMessage + fmt.Sprintf("%s\n", size.Value)
	}
	whData := discord.DiscordParams{
		Embeds: []discord.Embed{
			{
				Author: discord.Author{
					Name:    "Hypemedia Monitors",
					IconURL: "https://media.discordapp.net/attachments/906886120502075403/1095465707690524753/hypemedialov2.jpg?width=1170&height=1170",
				},
				Title: product.Name,
				Color: 16776975,
				Url:   fmt.Sprintf("https://www.jdsports.it/product/hm/%s/", product.Sku),
				Footer: discord.Footer{
					Text: "Hypemedia Monitors",
				},
				Fields: []discord.Field{
					{
						Name:   "Name",
						Value:  product.Name,
						Inline: true,
					},
					{
						Name:   "Sku",
						Value:  product.Sku,
						Inline: true,
					},
					{
						Name:   "Price",
						Value:  product.Price,
						Inline: false,
					},
					{
						Name:   "Sizes",
						Value:  sizeMessage,
						Inline: true,
					},
					{
						Name:   "Pids",
						Value:  pidsMessage,
						Inline: true,
					},
				},
				Thumbnail: discord.Thumbnail{
					URL: product.Image,
				},
			},
		},
	}

	webhookURL := os.Getenv("BONSAI_WEBHOOK_J_D_LOG")
	if webhookURL == "" {
		log.Error("Missing BONSAI_WEBHOOK_J_D_LOG", location)
		return
	}
	err := discord.PostWebhook(webhookURL, whData)
	if err != nil {
		log.Error("An error occured while sending the webhook!", location)
		fmt.Println(whData)
		time.Sleep(5 * time.Second)
		AdidasLog(product)
	} else {
		log.Success("Webhook sent", location)
		return
	}
}

func JDNewAddLog(product Product) {
	whData := discord.DiscordParams{
		Embeds: []discord.Embed{
			{
				Author: discord.Author{
					Name:    "Hypemedia Monitors",
					IconURL: "https://media.discordapp.net/attachments/906886120502075403/1095465707690524753/hypemedialov2.jpg?width=1170&height=1170",
				},
				Title: product.Name,
				Color: 16776975,
				Url:   fmt.Sprintf("https://www.jdsports.it/product/hm/%s/", product.Sku),
				Footer: discord.Footer{
					Text: "Hypemedia Monitors",
				},
				Fields: []discord.Field{
					{
						Name:   "Type",
						Value:  "New Add",
						Inline: true,
					},
					{
						Name:   "Name",
						Value:  product.Name,
						Inline: true,
					},
					{
						Name:   "Sku",
						Value:  product.Sku,
						Inline: true,
					},
					{
						Name:   "Price",
						Value:  product.Price,
						Inline: true,
					},
				},
				Thumbnail: discord.Thumbnail{
					URL: product.Image,
				},
			},
		},
	}

	webhookURL := os.Getenv("BONSAI_WEBHOOK_J_D_NEW_ADD_LOG")
	if webhookURL == "" {
		log.Error("Missing BONSAI_WEBHOOK_J_D_NEW_ADD_LOG", location)
		return
	}
	err := discord.PostWebhook(webhookURL, whData)
	if err != nil {
		log.Error("An error occured while sending the webhook!", location)
		fmt.Println(whData)
		time.Sleep(5 * time.Second)
		AdidasLog(product)
	} else {
		log.Success("Webhook sent", location)
		return
	}
}

func SizeRestock(product Product) {
	var sizeMessage string = ""
	var pidsMessage string = ""
	for _, size := range product.Sizes {
		sizeMessage = sizeMessage + fmt.Sprintf("%s\n", size.Size)
		pidsMessage = pidsMessage + fmt.Sprintf("%s\n", size.Value)
	}
	whData := discord.DiscordParams{
		Embeds: []discord.Embed{
			{
				Author: discord.Author{
					Name:    "Hypemedia Monitors",
					IconURL: "https://media.discordapp.net/attachments/906886120502075403/1095465707690524753/hypemedialov2.jpg?width=1170&height=1170",
				},
				Title: product.Name,
				Color: 16776975,
				Url:   fmt.Sprintf("https://www.sizeofficial.it/product/hm/%s/", product.Sku),
				Footer: discord.Footer{
					Text: "Hypemedia Monitors",
				},
				Fields: []discord.Field{
					{
						Name:   "Name",
						Value:  product.Name,
						Inline: true,
					},
					{
						Name:   "Sku",
						Value:  product.Sku,
						Inline: true,
					},
					{
						Name:   "Price",
						Value:  product.Price,
						Inline: false,
					},
					{
						Name:   "Sizes",
						Value:  sizeMessage,
						Inline: true,
					},
					{
						Name:   "Pids",
						Value:  pidsMessage,
						Inline: true,
					},
				},
				Thumbnail: discord.Thumbnail{
					URL: product.Image,
				},
			},
		},
	}

	webhookURL := os.Getenv("BONSAI_WEBHOOK_SIZE_RESTOCK")
	if webhookURL == "" {
		log.Error("Missing BONSAI_WEBHOOK_SIZE_RESTOCK", location)
		return
	}
	err := discord.PostWebhook(webhookURL, whData)
	if err != nil {
		log.Error("An error occured while sending the webhook!", location)
		fmt.Println(whData)
		time.Sleep(5 * time.Second)
		AdidasLog(product)
	} else {
		log.Success("Webhook sent", location)
		return
	}
}

func SizeNew(product Product) {
	whData := discord.DiscordParams{
		Embeds: []discord.Embed{
			{
				Author: discord.Author{
					Name:    "Hypemedia Monitors",
					IconURL: "https://media.discordapp.net/attachments/906886120502075403/1095465707690524753/hypemedialov2.jpg?width=1170&height=1170",
				},
				Title: product.Name,
				Color: 16776975,
				Url:   fmt.Sprintf("https://www.sizeofficial.it/product/hm/%s/", product.Sku),
				Footer: discord.Footer{
					Text: "Hypemedia Monitors",
				},
				Fields: []discord.Field{
					{
						Name:   "Type",
						Value:  "New Add",
						Inline: true,
					},
					{
						Name:   "Name",
						Value:  product.Name,
						Inline: true,
					},
					{
						Name:   "Sku",
						Value:  product.Sku,
						Inline: true,
					},
					{
						Name:   "Price",
						Value:  product.Price,
						Inline: true,
					},
				},
				Thumbnail: discord.Thumbnail{
					URL: product.Image,
				},
			},
		},
	}

	webhookURL := os.Getenv("BONSAI_WEBHOOK_SIZE_NEW")
	if webhookURL == "" {
		log.Error("Missing BONSAI_WEBHOOK_SIZE_NEW", location)
		return
	}
	err := discord.PostWebhook(webhookURL, whData)
	if err != nil {
		log.Error("An error occured while sending the webhook!", location)
		fmt.Println(whData)
		time.Sleep(5 * time.Second)
		AdidasLog(product)
	} else {
		log.Success("Webhook sent", location)
		return
	}
}

func FpRestock(product Product) {
	var sizeMessage string = ""
	var pidsMessage string = ""
	for _, size := range product.Sizes {
		sizeMessage = sizeMessage + fmt.Sprintf("%s\n", size.Size)
		pidsMessage = pidsMessage + fmt.Sprintf("%s\n", size.Value)
	}
	whData := discord.DiscordParams{
		Embeds: []discord.Embed{
			{
				Author: discord.Author{
					Name:    "Hypemedia Monitors",
					IconURL: "https://media.discordapp.net/attachments/906886120502075403/1095465707690524753/hypemedialov2.jpg?width=1170&height=1170",
				},
				Title: product.Name,
				Color: 16776975,
				Url:   fmt.Sprintf("https://www.footpatrol.it/product/hm/%s/", product.Sku),
				Footer: discord.Footer{
					Text: "Hypemedia Monitors",
				},
				Fields: []discord.Field{
					{
						Name:   "Name",
						Value:  product.Name,
						Inline: true,
					},
					{
						Name:   "Sku",
						Value:  product.Sku,
						Inline: true,
					},
					{
						Name:   "Price",
						Value:  product.Price,
						Inline: false,
					},
					{
						Name:   "Sizes",
						Value:  sizeMessage,
						Inline: true,
					},
					{
						Name:   "Pids",
						Value:  pidsMessage,
						Inline: true,
					},
				},
				Thumbnail: discord.Thumbnail{
					URL: product.Image,
				},
			},
		},
	}

	webhookURL := os.Getenv("BONSAI_WEBHOOK_FP_RESTOCK")
	if webhookURL == "" {
		log.Error("Missing BONSAI_WEBHOOK_FP_RESTOCK", location)
		return
	}
	err := discord.PostWebhook(webhookURL, whData)
	if err != nil {
		log.Error("An error occured while sending the webhook!", location)
		fmt.Println(whData)
		time.Sleep(5 * time.Second)
		AdidasLog(product)
	} else {
		log.Success("Webhook sent", location)
		return
	}
}

func FpNew(product Product) {
	whData := discord.DiscordParams{
		Embeds: []discord.Embed{
			{
				Author: discord.Author{
					Name:    "Hypemedia Monitors",
					IconURL: "https://media.discordapp.net/attachments/906886120502075403/1095465707690524753/hypemedialov2.jpg?width=1170&height=1170",
				},
				Title: product.Name,
				Color: 16776975,
				Url:   fmt.Sprintf("https://www.footpatrol.it/product/hm/%s/", product.Sku),
				Footer: discord.Footer{
					Text: "Hypemedia Monitors",
				},
				Fields: []discord.Field{
					{
						Name:   "Type",
						Value:  "New Add",
						Inline: true,
					},
					{
						Name:   "Name",
						Value:  product.Name,
						Inline: true,
					},
					{
						Name:   "Sku",
						Value:  product.Sku,
						Inline: true,
					},
					{
						Name:   "Price",
						Value:  product.Price,
						Inline: true,
					},
				},
				Thumbnail: discord.Thumbnail{
					URL: product.Image,
				},
			},
		},
	}

	webhookURL := os.Getenv("BONSAI_WEBHOOK_FP_NEW")
	if webhookURL == "" {
		log.Error("Missing BONSAI_WEBHOOK_FP_NEW", location)
		return
	}
	err := discord.PostWebhook(webhookURL, whData)
	if err != nil {
		log.Error("An error occured while sending the webhook!", location)
		fmt.Println(whData)
		time.Sleep(5 * time.Second)
		AdidasLog(product)
	} else {
		log.Success("Webhook sent", location)
		return
	}
}

func SotfLog(product Product) {
	var sizeMessage string = ""
	for _, size := range product.Sizes {
		sizeMessage = sizeMessage + fmt.Sprintf("%s\n", size.Size)
	}
	whData := discord.DiscordParams{
		Embeds: []discord.Embed{
			{
				Author: discord.Author{
					Name:    "Hypemedia Monitors",
					IconURL: "https://media.discordapp.net/attachments/906886120502075403/1095465707690524753/hypemedialov2.jpg?width=1170&height=1170",
				},
				Title: product.Name,
				Color: 16776975,
				Url:   fmt.Sprintf("%s", product.Sku),
				Footer: discord.Footer{
					Text: "Hypemedia Monitors",
				},
				Fields: []discord.Field{
					{
						Name:   "Name",
						Value:  product.Name,
						Inline: true,
					},
					{
						Name:   "Price",
						Value:  product.Price,
						Inline: true,
					},
					{
						Name:   "Sizes",
						Value:  sizeMessage,
						Inline: false,
					},
				},
				Thumbnail: discord.Thumbnail{
					URL: product.Image,
				},
			},
		},
	}

	webhookURL := os.Getenv("BONSAI_WEBHOOK_SOTF_LOG")
	if webhookURL == "" {
		log.Error("Missing BONSAI_WEBHOOK_SOTF_LOG", location)
		return
	}
	err := discord.PostWebhook(webhookURL, whData)
	if err != nil {
		log.Error("An error occured while sending the webhook!", location)
		fmt.Println(whData)
		time.Sleep(5 * time.Second)
		AdidasLog(product)
	} else {
		log.Success("Webhook sent", location)
		return
	}
}

func FtlKiller(product Product) {
	var sizeMessage string = ""
	for _, size := range product.Sizes {
		sizeMessage = sizeMessage + fmt.Sprintf("%s\n", size.Size)
	}
	whData := discord.DiscordParams{
		Embeds: []discord.Embed{
			{
				Author: discord.Author{
					Name:    "Hypemedia Monitors",
					IconURL: "https://media.discordapp.net/attachments/906886120502075403/1095465707690524753/hypemedialov2.jpg?width=1170&height=1170",
				},
				Title: product.Name,
				Color: 16776975,
				Url:   fmt.Sprintf("https://www.footlocker.it/product/~/%s.html", product.Sku),
				Footer: discord.Footer{
					Text: "Hypemedia Monitors",
				},
				Fields: []discord.Field{
					{
						Name:   "Sku",
						Value:  product.Sku,
						Inline: true,
					},
					{
						Name:   "Price",
						Value:  product.Price,
						Inline: true,
					},
					{
						Name:   "Color",
						Value:  product.Color,
						Inline: true,
					},
					{
						Name:   "Size",
						Value:  sizeMessage,
						Inline: true,
					},
				},
				Thumbnail: discord.Thumbnail{
					URL: product.Image,
				},
			},
		},
	}

	webhookURL := os.Getenv("BONSAI_WEBHOOK_FTL_KILLER")
	if webhookURL == "" {
		log.Error("Missing BONSAI_WEBHOOK_FTL_KILLER", location)
		return
	}
	err := discord.PostWebhook(webhookURL, whData)
	if err != nil {
		log.Error("An error occured while sending the webhook!", location)
		fmt.Println(whData)
		time.Sleep(5 * time.Second)
		AdidasLog(product)
	} else {
		log.Success("Webhook sent", location)
		return
	}
}

func AsosLog(product Product) {
	var sizeMessage string = ""
	for _, size := range product.Sizes {
		sizeMessage = sizeMessage + fmt.Sprintf("[ATC](https://www.asos.com/it/failed-order/?variants=%s:1) ", size.Value) + fmt.Sprintf("%s\n", size.Size)
	}
	whData := discord.DiscordParams{
		Embeds: []discord.Embed{
			{
				Author: discord.Author{
					Name:    "Hypemedia Monitors",
					IconURL: "https://media.discordapp.net/attachments/906886120502075403/1095465707690524753/hypemedialov2.jpg?width=1170&height=1170",
				},
				Title: product.Name,
				Color: 16776975,
				Url:   fmt.Sprintf("https://www.asos.com/it/hm/prd/%s", product.Sku),
				Footer: discord.Footer{
					Text: "Hypemedia Monitors",
				},
				Fields: []discord.Field{
					{
						Name:   "Sku",
						Value:  product.Sku,
						Inline: true,
					},
					{
						Name:   "Price",
						Value:  product.Price,
						Inline: true,
					},
					{
						Name:   "Color",
						Value:  product.Color,
						Inline: true,
					},
					{
						Name:   "Size",
						Value:  sizeMessage,
						Inline: true,
					},
				},
				Thumbnail: discord.Thumbnail{
					URL: fmt.Sprintf("http://api.lollipopext.com/proxy?url=%s", "https://"+product.Image),
				},
			},
		},
	}

	webhookURL := os.Getenv("BONSAI_WEBHOOK_ASOS_LOG")
	if webhookURL == "" {
		log.Error("Missing BONSAI_WEBHOOK_ASOS_LOG", location)
		return
	}
	err := discord.PostWebhook(webhookURL, whData)
	if err != nil {
		log.Error("An error occured while sending the webhook!", location)
		fmt.Println(whData)
		time.Sleep(5 * time.Second)
		AdidasLog(product)
	} else {
		log.Success("Webhook sent", location)
		return
	}
}

func NamshiLog(product Product) {
	var sizeMessage string = ""
	for _, size := range product.Sizes {
		sizeMessage = sizeMessage + fmt.Sprintf("%s\n", size.Size)
	}
	whData := discord.DiscordParams{
		Embeds: []discord.Embed{
			{
				Author: discord.Author{
					Name:    "Hypemedia Monitors",
					IconURL: "https://media.discordapp.net/attachments/906886120502075403/1095465707690524753/hypemedialov2.jpg?width=1170&height=1170",
				},
				Title: product.Name,
				Color: 16776975,
				Url:   fmt.Sprintf("https://en-global.namshi.com%s", product.Sku),
				Footer: discord.Footer{
					Text: "Hypemedia Monitors",
				},
				Fields: []discord.Field{
					{
						Name:   "Price",
						Value:  product.Price,
						Inline: true,
					},
					{
						Name:   "Size",
						Value:  sizeMessage,
						Inline: true,
					},
				},
				Thumbnail: discord.Thumbnail{
					URL: product.Image,
				},
			},
		},
	}

	webhookURL := os.Getenv("BONSAI_WEBHOOK_NAMSHI_LOG")
	if webhookURL == "" {
		log.Error("Missing BONSAI_WEBHOOK_NAMSHI_LOG", location)
		return
	}
	err := discord.PostWebhook(webhookURL, whData)
	if err != nil {
		log.Error("An error occured while sending the webhook!", location)
		fmt.Println(whData)
		time.Sleep(5 * time.Second)
		AdidasLog(product)
	} else {
		log.Success("Webhook sent", location)
		return
	}
}
