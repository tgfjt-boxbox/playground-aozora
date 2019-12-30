package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/tgfjt/aozora"
)

func main() {
	// 作家別作品一覧(拡充版)
	listOfWorks, err := aozora.ListOfWorks()
	if err != nil {
		log.Fatal("error:", err)
	}

	listOfPublisher := make([]string, 0)

	// とにかく全部いれまくる
	for _, work := range listOfWorks {
		// OriginBook1PublisherName          底本出版社名1
		// ParentOfOriginBook1PublisherName  底本の親本出版社名1
		// OriginBook2PublisherName          底本出版社名2
		// ParentOfOriginBook2PublisherName  底本の親本出版社名2
		if strings.TrimSpace(work.OriginBook1PublisherName) != "" {
			listOfPublisher = append(listOfPublisher, strings.TrimSpace(work.OriginBook1PublisherName))
			if strings.TrimSpace(work.ParentOfOriginBook1PublisherName) != "" {
				listOfPublisher = append(listOfPublisher, strings.TrimSpace(work.ParentOfOriginBook1PublisherName))
			}
		}

		if strings.TrimSpace(work.OriginBook2PublisherName) != "" {
			listOfPublisher = append(listOfPublisher, strings.TrimSpace(work.OriginBook2PublisherName))
			if strings.TrimSpace(work.ParentOfOriginBook2PublisherName) != "" {
				listOfPublisher = append(listOfPublisher, strings.TrimSpace(work.ParentOfOriginBook2PublisherName))
			}
		}
	}

	sort.Strings(listOfPublisher)

	m := make(map[string]struct{})
	newList := make([]string, 0)

	// ユニークな出版社名だけ欲しい
	for _, element := range listOfPublisher {
		if _, ok := m[element]; !ok {
			m[element] = struct{}{}
			newList = append(newList, element)
		}
	}

	fmt.Println(strings.Join(newList, "\n"))
}
