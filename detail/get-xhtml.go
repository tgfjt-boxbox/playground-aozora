package detail

import (
	"flag"
	"fmt"
	"log"

	"github.com/tgfjt/aozora"
)

func main() {
	// 作家別作品一覧(拡充版)
	listOfWorks, err := aozora.ListOfWorks()
	if err != nil {
		log.Fatal("error:", err)
	}

	// go run detail/get-xhtml.go --id 43737
	cardID := flag.String("id", "", "Card ID")

	flag.Parse()

	result := listOfWorks.Where(func(w aozora.WorkExpanded) bool {
		return w.CardID == fmt.Sprintf("%06s", *cardID)
	})

	if len(result) == 1 {
		fmt.Println(result[0])
	}
}
