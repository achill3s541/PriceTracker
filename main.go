package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type shop struct {
	Shops []*product `json:"Shops"`
}

type product struct {
	ShopsName       string  `json:"ShopsName"`
	Variant         string  `json:"Variant"`
	Price           float64 `json:"Price"`
	AddressURL      string  `json:"URL"`
	LastUpadateDate string  `json:"LastUpdateDate"`
	PriceAlert      float64 `json:"PriceAlert"`
}

// var compareVariantFromJSON []string
// var compareURLFromJSON []string
func parseContent(website string, filename string, time string, priceAlert []float64) ([]float64, []string, string, error) {
	var gotPriceFromContent []float64
	var comparePriceFromContent []float64
	var gotVariantFromContent []string
	var compareVariantFromContent []string
	url, err := http.Get(website)
	if err != nil {
		return nil, nil, "", fmt.Errorf("the website doesn't work, check the URL is properly")
	}
	defer url.Body.Close()
	// The website's contentet is saved to the variable
	page, err := goquery.NewDocumentFromReader(url.Body)
	if err != nil {
		return nil, nil, "", fmt.Errorf("cannot to read the website's content, try again later")
	}
	page.Find(".z-price__amount").Each(func(i int, pageOutput *goquery.Selection) {
		// This regular expresion is used to find a price and convert it to float64
		value := pageOutput.Text()
		regex := regexp.MustCompile(`[zł].`)
		value = regex.ReplaceAllString(value, "")
		value = strings.Replace(value, ",", ".", 1)
		value = strings.TrimSpace(value)
		stringPriceFromWebsite, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Fatal(err)
		}
		gotPriceFromContent = append(gotPriceFromContent, stringPriceFromWebsite)
	})
	page.Find(".Variant_variantWrapper__eUlYB").Each(func(i int, pageOutput *goquery.Selection) {
		//This regular expresion is used to find and to remove the product's number from the product's name
		regex := regexp.MustCompile(`(\s[g]\s*\d*.\d*.\d*)|([)]\d*.\d*)`)
		gotVariantFromContent = append(gotVariantFromContent, pageOutput.Text())
		gotVariantFromContent[i] = regex.ReplaceAllString(gotVariantFromContent[i], "")
	})
	//The website's content is converting to JSON's format.
	var shops []*product
	for i := range gotVariantFromContent {
		shops, comparePriceFromContent, compareVariantFromContent = append(shops, &product{ShopsName: "zooplus", Variant: gotVariantFromContent[i], Price: gotPriceFromContent[i], AddressURL: string(website), LastUpadateDate: time, PriceAlert: priceAlert[i]}), append(gotPriceFromContent), append(gotVariantFromContent)

	}
	pageContent, err := json.Marshal(shop{Shops: shops})
	if err != nil {
		return nil, nil, "", fmt.Errorf("Test1")
	}
	//The website's content is wroten to JSON's file.
	err = os.WriteFile(filename, pageContent, 0644)
	if err != nil {
		return nil, nil, "", fmt.Errorf("Test2")
	}
	return comparePriceFromContent, compareVariantFromContent, website, nil
}

func readingJSONFile(filename string) ([]float64, []float64, []string, error) {
	var comparePriceFromJSON []float64
	var priceAlertFromJSON []float64
	var variantFromJSON []string

	// This fucntion reads data from JSON's file,
	//	it verifies the JSON's file exists
	checkFile, err := os.Open(filename)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("file %s does not exist", filename)
	}
	defer checkFile.Close()
	fileJson, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("cannot open the file %s", filename)
	}
	conentFileJson := shop{}
	err = json.Unmarshal([]byte(fileJson), &conentFileJson)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(conentFileJson.Shops); i++ {
		//The product's price is got from every product in the file
		comparePriceFromJSON = append(comparePriceFromJSON, conentFileJson.Shops[i].Price)
		priceAlertFromJSON = append(priceAlertFromJSON, conentFileJson.Shops[i].PriceAlert)
		variantFromJSON = append(variantFromJSON, conentFileJson.Shops[i].Variant)
	}
	return priceAlertFromJSON, comparePriceFromJSON, variantFromJSON, nil
}

func compareContToJSON(variantJSON []string, variant []string, priceJSON []float64, priceContent []float64, priceAlert []float64, website string) error {
	// This function compares JSON's data to the fresh Website's data, if the old price has changed, the message will be displayed and the email will be sent.
	for i := range variant {
		if priceJSON[i] > priceContent[i] {
			email := fmt.Sprintf("Dla produktu %v w wariancie %v pojawiła się tańsza wersja.\nStara cena to: %0.2f\nAktualna cena to: %0.2f \n", website, variant[i], priceJSON[i], priceContent[i])
			err := emailSender(email, "newPrice")
			if err != nil {
				return err
			}
		} else if priceAlert[i] > priceContent[i] && variantJSON[i] == variant[i] {
			email := fmt.Sprintf("Produkt %v w wariancie %v osiągnął cenę ze wskazanego progu.\nStara cena to: %0.2f\nAktualna cena to: %0.2f \n", website, variant[i], priceJSON[i], priceContent[i])
			err := emailSender(email, "alertPrice")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func emailSender(messageInput string, subjectPrefix string) error {
	//This function is responsible for sending email if the Webiste's price is lower then JSON's price.
	var subject string
	authentication := smtp.PlainAuth("", "<fill_email_addres>", "<fill_password>", "smtp.gmail.com")
	sendingTo := []string{"<fill_email_addres>"}
	sender := fmt.Sprintf("From: <%s>\r\n", "<fill_email_addres>")
	receiver := fmt.Sprintf("To: <%s>\r\n", "<fill_email_addres>")
	if subjectPrefix == "alertPrice" {
		subject = "Subject: [Alert!] The price has reached the alarm value \r\n"
	} else {
		subject = "Subject: The price on the website is lower \r\n"
	}
	body := messageInput + "\r\n"
	//This variable builds a request for email
	messaging := sender + receiver + subject + "\r\n" + body
	err := smtp.SendMail("smtp.gmail.com:587", authentication, "<fill_email_addres>", sendingTo, []byte(messaging))
	if err != nil {
		return fmt.Errorf("the system cannot send email message: %s", err)
	}
	fmt.Print("The message has been sent to receiver\n")
	return err
}

func main() {
	var website string
	currentTime := time.Now().Format("02-01-2006 15:04:05")
	priceAlertFromJSON, comparePriceFromJSON, variantFromJSON, err := readingJSONFile("tracker_output.json")
	if err != nil {
		fmt.Println(err)
	}
	comparePriceFromContent, compareVariantFromContent, website, err := parseContent("https://www.zooplus.pl/shop/koty/zwirek_dla_kota/benek/compact/335113", "tracker_output.json", currentTime, priceAlertFromJSON)
	if err != nil {
		fmt.Println(err)
	}
	err = compareContToJSON(variantFromJSON, compareVariantFromContent, comparePriceFromJSON, comparePriceFromContent, priceAlertFromJSON, website)
	if err != nil {
		fmt.Println(err)
	}
}
