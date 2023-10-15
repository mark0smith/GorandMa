package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
)

// main is the main function of the program
func downloadFile(webUrl, filename string) {
	resp, err := http.Get(webUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// create a new file with the file name in the current directory
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// copy the response body to the file using io.Copy function
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		panic(err)
	}

	// print a message indicating that the file has been downloaded successfully
	fmt.Println("Downloaded", filename)
}

func parseRelativeUrl(baseUrl, relativeUrl string) string {
	base, err := url.Parse(baseUrl)
	if err != nil {
		log.Fatal(err)
	}
	ref, err := url.Parse(relativeUrl)
	if err != nil {
		log.Fatal(err)
	}
	u := base.ResolveReference(ref)
	return u.String()
}

func FileExist(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		// path/to/whatever exists
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does *not* exist
		return false
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		return false
	}
}

func formatText(text string) string {
	re := regexp.MustCompile(`[^\p{Han}\w\-《》：]`)
	return re.ReplaceAllString(text, "")
}

func downloadWtc(webUrl string) {
	resp, err := soup.Get(webUrl)
	if err != nil {
		os.Exit(1)
	}
	doc := soup.HTMLParse(resp)

	wtcTheme := doc.FindAll("h1")
	var folderSlice []string
	for _, theme := range wtcTheme {
		folderSlice = append(folderSlice, theme.Text())
	}
	folderName := ""
	if len(folderSlice) == 0 {
		if strings.Contains(webUrl, "Past") {
			folderName = strings.Split(webUrl, "/")[4]
			folderName = strings.Split(folderName, ".")[0]
			folderName = folderName + "年安天网络安全冬训营"
		} else {
			datetime := time.Now()
			currentYear := datetime.Format("2006")
			folderName = fmt.Sprintf("%s年安天网络安全冬训营", currentYear)
			fmt.Println(folderName)
		}

	} else {
		folderName = strings.Join(folderSlice, "——")
	}

	os.MkdirAll(folderName, os.ModePerm)

	// links := doc.Find("tbody").FindAll("tr")
	// for _, link := range links {
	// 	tr1 := link.Children()
	// 	href := tr1[1].Find("a")
	// 	if href.NodeValue == "a" {
	// 		filenameP1 := strings.TrimSpace(tr1[1].Text())
	// 		filenameP2 := tr1[1].FindNextElementSibling().Text()
	// 		filename := fmt.Sprintf("%s（%s）.pdf", filenameP1, filenameP2)
	// 		filename = strings.ReplaceAll(filename, "/", "_")

	// 		filename = fmt.Sprintf("./%s/%s", folderName, filename)

	// 		relativeUrl := href.Attrs()["href"]
	// 		fileUrl := parseRelativeUrl(webUrl, relativeUrl)
	// 		fmt.Printf("[+]Downloading %s from %s ...\n", filename, fileUrl)
	// 		downloadFile(fileUrl, filename)
	// 	}

	// }

	topics := doc.Find("tbody").FindAll("td")
	for _, topic := range topics {

		if topic.Find("a").NodeValue != "a" {
			continue
		}

		topicText := strings.TrimSpace(topic.Text())
		topicText = formatText(topicText)

		author := topic.FindNextElementSibling()
		authorText := ""
		filename := ""
		if author.NodeValue == "td" {
			authorText = author.Text()
			authorText = formatText(authorText)
			filename = fmt.Sprintf("%s（%s）.pdf", topicText, authorText)
		} else {
			filename = fmt.Sprintf("%s.pdf", topicText)
		}
		filename = strings.ReplaceAll(filename, "/", "_")
		filename = fmt.Sprintf("./%s/%s", folderName, filename)

		link := topic.Find("a")
		if link.NodeValue == "a" {
			relativeUrl := link.Attrs()["href"]
			fileUrl := parseRelativeUrl(webUrl, relativeUrl)

			if !FileExist(filename) {
				fmt.Printf("[+]Downloading %s from %s ...\n", filename, fileUrl)
				downloadFile(fileUrl, filename)
			}

		}

	}

}

func main() {

	var webUrl []string

	webUrl = append(webUrl, "https://wtc.antiy.cn/Past/2014.html")
	webUrl = append(webUrl, "https://wtc.antiy.cn/Past/2015.html")
	webUrl = append(webUrl, "https://wtc.antiy.cn/Past/2016.html")
	webUrl = append(webUrl, "https://wtc.antiy.cn/Past/2017.html")
	webUrl = append(webUrl, "https://wtc.antiy.cn/Past/2018.html")
	webUrl = append(webUrl, "https://wtc.antiy.cn/Past/2019.html")
	webUrl = append(webUrl, "https://wtc.antiy.cn/Past/2021.html")
	webUrl = append(webUrl, "https://wtc.antiy.cn/Past/2022.html")
	webUrl = append(webUrl, "https://wtc.antiy.cn")

	for _, url := range webUrl {
		downloadWtc(url)
	}
}
