package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	beginTime := time.Now()

	params := readFlugs()
	if params == nil {
		return
	}

	fmt.Println("Запуск")
	var strings = readLinks(params[0])
	err := connect(strings, params[1])
	if err != nil {
		return
	}
	endTime := time.Now()

	fmt.Printf("Время работы программы: %s\n", endTime.Sub(beginTime))

}

func connect(links []string, savePath string) error {
	for _, s := range links {
		fmt.Printf("%s: ", s)
		url := fmt.Sprintf("https://%s", s)

		response, err := http.Get(url)
		if err != nil {
			fmt.Printf("ошибка запроса\n")
			continue
		}

		defer response.Body.Close()

		doc, _ := io.ReadAll(response.Body)
		err = saveHtml(savePath, s, string(doc))
		if err != nil {
			return errors.New("can't save html file")
		}
		fmt.Printf("успешно\n")

	}
	return nil

}

func saveHtml(path string, name string, html string) error {

	pathFile := fmt.Sprintf("%s/%s.html", path, name)

	file, err := os.Create(pathFile)
	if err != nil {

		err = os.MkdirAll(path, 0777)

		if err != nil {
			fmt.Printf("Неудалось создать файл")
			return errors.New("can't create file")
		}
		file.Close()
		file, err = os.Create(pathFile)
		if err != nil {
			fmt.Printf("Неудалось создать файл")
			return errors.New("can't create file")
		}
	}
	defer file.Close()
	file.WriteString(html)
	return nil
}

func readFlugs() []string {
	pathLinks := flag.String("src", "input.txt", "path txt file")
	pathDic := flag.String("dst", "targetDir", "path txt file")

	flag.Parse()

	result := []string{*pathLinks, *pathDic}
	return result
}

func readLinks(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Неудалось считать файл")
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var links []string

	for scanner.Scan() {
		lineStr := scanner.Text()
		lineStr = strings.TrimSpace(lineStr)
		if lineStr != "" {
			links = append(links, lineStr)
		}
	}

	return links

}
