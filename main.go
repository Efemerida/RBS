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
	fmt.Printf("\nЗапуск\n\n")

	//метка старта программы
	beginTime := time.Now()

	//чтение флагов
	pathLinks, pathDirectory, readFlugsErr := readFlugs()
	if readFlugsErr != nil {
		panic(readFlugsErr)
	}

	//парсинг ссылок, получение html и его сохранение
	errReadLinksGetAndSaveHtml := readLinksGetAndSaveHtml(*pathLinks, *pathDirectory)
	if errReadLinksGetAndSaveHtml != nil {
		panic(errReadLinksGetAndSaveHtml)
	}

	//метка завершения программы
	endTime := time.Now()

	fmt.Printf("Время работы программы: %s\n", endTime.Sub(beginTime))

}

// saveHtml - сохранение html документа
func saveHtml(savePath string, nameFile string, htmlDocument string) error {

	//формирования пути сохранения
	pathFile := fmt.Sprintf("%s/%s.html", savePath, nameFile)

	//открытие файла
	file, errOpenFile := os.Create(pathFile)
	if errOpenFile != nil {

		//создание пути сохранения, если он не найден
		errMakePath := os.MkdirAll(savePath, 0777)

		if errMakePath != nil {
			errorStr := fmt.Sprintf("Неудалось создать файл %s\nОшибка:%s\n\n", pathFile, errMakePath)
			return errors.New(errorStr)
		}

		//повторная попытка окрытия
		file, errOpenFile = os.Create(pathFile)
		if errOpenFile != nil {
			errorStr := fmt.Sprintf("Неудалось создать файл %s\nОшибка:%s\n\n", pathFile, errOpenFile)
			return errors.New(errorStr)
		}
	}
	defer file.Close()

	file.WriteString(htmlDocument)
	return nil
}

// readFlugs - чтение флагов
func readFlugs() (*string, *string, error) {
	pathLinks := flag.String("src", "", "Путь на файл с ссылками")
	pathDirectory := flag.String("dst", "", "Путь для сохранения html документов")

	flag.Parse()

	//если флаг на файл с ссылками не установлен
	if *pathLinks == "" {
		flag.PrintDefaults()
		return nil, nil, errors.New("не указан файл, содержащий ссылки")
	}

	//если флаг на директорию для сохранения не установлен
	if *pathDirectory == "" {
		fmt.Printf("Вы не указали директорию сохранения файлов.\nБудет выбрана директория по умолчанию (./targetDir)\n\n")
		*pathDirectory = "targetDir"
	}

	return pathLinks, pathDirectory, nil
}

// readLinksGetAndSaveHtml - чтение ссылок из файла, получение по ним html документа и запись его в файл
func readLinksGetAndSaveHtml(pathInputFile string, savePath string) error {

	//открытие файла с ссылками
	file, err := os.Open(pathInputFile)
	if err != nil {
		errorStr := fmt.Sprintf("Неудалось считать файл %s\n\n", pathInputFile)
		return errors.New(errorStr)
	}
	defer file.Close()

	var wg sync.WaitGroup

	//чтение из файла
	var countlinks map[string]bool = map[string]bool{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		link := scanner.Text()
		link = strings.TrimSpace(link)

		//если не пустая строка в качестве ссылки,
		//то выполнение запроса и сохранение html документа
		if link != "" {

			if _, ok := countlinks[link]; ok {
				continue
			}
			countlinks[link] = true
			wg.Add(1)
			go getAndSaveHtmlFromLink(link, savePath, &wg)
		}
	}

	wg.Wait()
	return nil
}

// getAndSaveHtmlFromLink - получение и сохранение html по ссылке
func getAndSaveHtmlFromLink(link string, savePath string, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	//формирование url
	link = strings.Trim(link, "htps:/")
	url := fmt.Sprintf("https://%s", link)

	//выполнение запроса и получение html
	response, errGet := http.Get(url)
	if errGet != nil {
		fmt.Printf("%s: ошибка запроса:\n%s\n\n", link, errGet)
		return
	}

	//чтение html документа из тела запроса и его сохранение
	defer response.Body.Close()
	responceBody, _ := io.ReadAll(response.Body)
	errSaveHtml := saveHtml(savePath, link, string(responceBody))
	if errSaveHtml != nil {
		fmt.Printf("%s: неудалось создать файл:\n%s\n\n", link, errSaveHtml)
		return
	}

	fmt.Printf("%s: успешно\n\n", link)
}
