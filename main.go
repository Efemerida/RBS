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
	"sync"
	"time"
)

func main() {
	beginTime := time.Now()

	params := readFlugs()
	if params == nil {
		return
	}

	fmt.Println("Запуск")
	err := readLinks(params[0], params[1])
	if err != nil {
		return
	}
	endTime := time.Now()

	fmt.Printf("Время работы программы: %s\n", endTime.Sub(beginTime))

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

func readLinks(path string, savePath string) error {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Неудалось считать файл")
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var wg sync.WaitGroup
	for scanner.Scan() {
		lineStr := scanner.Text()
		lineStr = strings.TrimSpace(lineStr)
		if lineStr != "" {
			wg.Add(1)
			go func() {
				defer wg.Done()
				url := fmt.Sprintf("https://%s", lineStr)
				response, err := http.Get(url)
				if err != nil {
					fmt.Printf("%s: ошибка запроса\n", lineStr)
					return
				}

				defer response.Body.Close()

				doc, _ := io.ReadAll(response.Body)
				err = saveHtml(savePath, lineStr, string(doc))
				if err != nil {
					panic("can't save html file")
				}
				fmt.Printf("%s: успешно\n", lineStr)
			}()

		}
	}

	wg.Wait()
	return nil
}
