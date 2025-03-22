package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"github.com/xuri/excelize/v2"
	"google.golang.org/api/option"
)

func singleResponse(resp *genai.GenerateContentResponse) string {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				return fmt.Sprintf("%v", part)
			}
		}
	}
	return ""
}

func main() {
	fmt.Println("Hello, World!")

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	// [START text_gen_text_only_prompt]
	model := client.GenerativeModel("gemini-1.5-flash")

	f, err := excelize.OpenFile("./geminifun/data/dad_library.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	rows, err := f.GetRows("FamousBooks")
	if err != nil {
		fmt.Println(err)
		return
	}

	var title string
	var author string
	var language string
	var country string
	var year string
	var characteristc string

	maxRows := 5

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if i > maxRows {
			break
		}

		for j, colCell := range row {
			switch j {
			case 0:
				country = colCell
			case 1:
				language = colCell
			case 2:
				characteristc = colCell
			case 3:
				author = colCell
			case 4:
				title = colCell
			case 5:
				year = colCell
			case 6:
				continue
			case 7:
				continue
			default:
				log.Fatal("Invalid column number")
			}
		}
		resp, err := model.GenerateContent(ctx, genai.Text("I want you to answer me with just the genre of a book. The book is written by "+author+" and the title (in spanish) is "+title+". The book is written in "+language+" and the country of origin is "+country+". The book was published around "+year+". The book is related to the following characteristic: "+characteristc+"."))
		if err != nil {
			log.Fatal(err)
		}
		f.SetCellStr("FamousBooks", "G"+fmt.Sprint(i+1), strings.TrimSpace(singleResponse(resp)))
		time.Sleep(200 * time.Millisecond)

		resp, err = model.GenerateContent(ctx, genai.Text("I want you to answer me with just a short synopis of a book. The book is written by "+author+" and the title (in spanish) is "+title+". The book is written in "+language+" and the country of origin is "+country+". The book was published around "+year+". The book is related to the following characteristic: "+characteristc+"."))
		if err != nil {
			log.Fatal(err)
		}
		f.SetCellStr("FamousBooks", "H"+fmt.Sprint(i+1), strings.TrimSpace(singleResponse(resp)))
		time.Sleep(200 * time.Millisecond)
	}
	if err := f.Save(); err != nil {
		fmt.Println(err)
		return
	}
}
