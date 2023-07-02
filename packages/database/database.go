package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		break
	case "Footlocker":
		FootlockerLog(product)
		break
	case "Zalando":
		ZalandoLog(product)
		break
	case "Bootz":
		BootzLog(product)
		break
	case "Rinascente":
		RinaLog(product)
		break
	case "Mana":
		ManaLog(product)
		break
	case "SkateDeluxe":
		SkateLog(product)
		break
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
				Url:   fmt.Sprintf("https://www.adidas.it/hm/B22716.html%s", product.Sku),
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

	err := discord.PostWebhook("https://discord.com/api/webhooks/1123628879416930434/fWivbwuzHBWBJhZ9SU3BVSFLvIGnQCyloDD34BS5eMPEHu810ueucqupKxyaHxdhjqV8", whData)
	//err := discord.PostWebhook("https://discord.com/api/webhooks/1119683086507057242/jVRDjqmMVBmlelNnUrw5Hkr9vc8fg-v8yrz7_1Qkvny7w2dSn6uJPL9WqG7vkcc-qWQa", whData)
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
				},
				Thumbnail: discord.Thumbnail{
					URL: product.Image,
				},
			},
		},
	}

	err := discord.PostWebhook("https://discord.com/api/webhooks/1119684786529779742/kwvf9VL0T-pmSh8M9juBBJSBfRS_zMdvSEquM8xSLlqoACcRJgPDSHSpiZbelzPu3qGN", whData)
	//err := discord.PostWebhook("https://discord.com/api/webhooks/1119683086507057242/jVRDjqmMVBmlelNnUrw5Hkr9vc8fg-v8yrz7_1Qkvny7w2dSn6uJPL9WqG7vkcc-qWQa", whData)
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

	err := discord.PostWebhook("https://discord.com/api/webhooks/1120085184163885247/_boNPxubvjk50fAhIjVmwpIwLFDlWYMsr_FM7ZsrToUWEkv16pImp9FAB_MaQAaI9Zo-", whData)
	//err := discord.PostWebhook("https://discord.com/api/webhooks/1119683086507057242/jVRDjqmMVBmlelNnUrw5Hkr9vc8fg-v8yrz7_1Qkvny7w2dSn6uJPL9WqG7vkcc-qWQa", whData)
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

	err := discord.PostWebhook("https://discord.com/api/webhooks/1124468356641206405/I0A0H1cFBTr0_XNEcC6RXGVjw5eIYZj8_njszBQMtyOL4t3JBR0L9jkoqYuzy3TJvwtY", whData)
	//err := discord.PostWebhook("https://discord.com/api/webhooks/1119683086507057242/jVRDjqmMVBmlelNnUrw5Hkr9vc8fg-v8yrz7_1Qkvny7w2dSn6uJPL9WqG7vkcc-qWQa", whData)
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

	err := discord.PostWebhook("https://discord.com/api/webhooks/1124468356641206405/I0A0H1cFBTr0_XNEcC6RXGVjw5eIYZj8_njszBQMtyOL4t3JBR0L9jkoqYuzy3TJvwtY", whData)
	//err := discord.PostWebhook("https://discord.com/api/webhooks/1124657248522752091/LZC0u58A4FvxTLIiV8PL4dow5wRfslch_AVM72c8SqYdmeOnDRA_FYwWUmSN5OccdXIX", whData)
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

	err := discord.PostWebhook("https://discord.com/api/webhooks/1124468356641206405/I0A0H1cFBTr0_XNEcC6RXGVjw5eIYZj8_njszBQMtyOL4t3JBR0L9jkoqYuzy3TJvwtY", whData)
	//err := discord.PostWebhook("https://discord.com/api/webhooks/1124657248522752091/LZC0u58A4FvxTLIiV8PL4dow5wRfslch_AVM72c8SqYdmeOnDRA_FYwWUmSN5OccdXIX", whData)
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

	err := discord.PostWebhook("https://discord.com/api/webhooks/1104086275276951644/LXydWEdohhMoGBcG9Rx1Kzh7vGztwvCZkmpcFhDWMzN_CoTtUu1HxLekH039q79kXagG", whData)
	//err := discord.PostWebhook("https://discord.com/api/webhooks/1124657248522752091/LZC0u58A4FvxTLIiV8PL4dow5wRfslch_AVM72c8SqYdmeOnDRA_FYwWUmSN5OccdXIX", whData)
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
