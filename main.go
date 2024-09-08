package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"html/template"
	"unicode"
)

type Result struct {
	Summary   string
	Sentiment string
}

func main() {
  http.HandleFunc("/", handleHome)
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

  log.Println("Server started on http://localhost:8080")
  log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		renderTemplate(w, Result{})
		return
	}

	inputText := r.FormValue("inputText")
	var fullText string

	if isURL(inputText) {
		fetchedText, err := fetchTextFromURL(inputText)
		if err != nil {
			http.Error(w, "Unable to fetch the URL", http.StatusBadRequest)
			return
		}
		fullText = fetchedText
	} else {
		fullText = inputText
	}

	summary := summarizeText(fullText)
	sentiment := analyzeSentiment(fullText)

	renderTemplate(w, Result{
		Summary:   summary,
		Sentiment: sentiment,
	})
}

func renderTemplate(w http.ResponseWriter, data Result) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, data)
}

func isURL(input string) bool {
	_, err := url.ParseRequestURI(input)
	return err == nil
}

func fetchTextFromURL(inputURL string) (string, error) {
	response, err := http.Get(inputURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch the URL: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read the response body: %v", err)
	}

	return extractTextFromHTML(string(body)), nil
}

func extractTextFromHTML(htmlContent string) string {
	// A very basic and rudimentary way to extract visible text from HTML
	// by removing all the tags. This is not a perfect solution, but
	// it works for this example without using a third-party library.
	inTag := false
	var output strings.Builder
	for _, ch := range htmlContent {
		if ch == '<' {
			inTag = true
		} else if ch == '>' {
			inTag = false
		} else if !inTag {
			output.WriteRune(ch)
		}
	}
	return strings.TrimSpace(output.String())
}

func summarizeText(text string) string {
	words := strings.Fields(text)
	if len(words) > 20 {
		return strings.Join(words[:20], " ") + "..."
	}
	return text
}

func analyzeSentiment(text string) string {
	positiveWords := []string{"good", "happy", "love", "great", "excellent", "positive", "awesome"}
	negativeWords := []string{"bad", "sad", "hate", "terrible", "horrible", "negative"}

	text = strings.ToLower(text)
	posCount, negCount := 0, 0

	for _, word := range strings.FieldsFunc(text, func(c rune) bool { return !unicode.IsLetter(c) }) {
		for _, pos := range positiveWords {
			if word == pos {
				posCount++
			}
		}
		for _, neg := range negativeWords {
			if word == neg {
				negCount++
			}
		}
	}

	if posCount > negCount {
		return "Positive"
	} else if negCount > posCount {
		return "Negative"
	}
	return "Neutral"
}
