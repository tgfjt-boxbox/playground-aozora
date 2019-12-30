package main

import (
	"context"
	"fmt"
	"log"
	"os/user"
	"sort"
	"strings"

	firebase "firebase.google.com/go"
	"github.com/tgfjt/aozora"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Publisher struct
type Publisher struct {
	// ID string
	Name string `firestore:"name,omitempty"`
}

// initApp init firebase app
func initApp(ctx context.Context) (*firebase.App, error) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	sa := option.WithCredentialsFile(usr.HomeDir + "/masters/forest-night-firebase-admin.json")
	app, err := firebase.NewApp(ctx, nil, sa)

	if err != nil {
		return app, err
	}

	return app, nil
}

func GetUniquePublisher() []string {
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

	fmt.Println("get unique publishers")

	return newList
}

func main() {
	// 作家別作品一覧(拡充版)
	newList := GetUniquePublisher()

	ctx := context.Background()

	app, err := initApp(ctx)
	if err != nil {
		panic(fmt.Sprintf("error initializing app: %v", err))
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
	}
	defer client.Close()

	iter := client.Collection("aozora_publishers").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		fmt.Println(doc.Data())
	}

	for _, name := range newList {
		_, err = client.Collection("aozora_publishers").Doc(name).Set(ctx, Publisher{
			Name: name,
		})
	}

	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
	}

	// fmt.Println(strings.Join(newList, "\n"))
}
