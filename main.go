package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"unicode"
)

type Result struct {
	Summary   string
	Sentiment string
}

func main() {
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	http.HandleFunc("/", handleHome)

	// Serve static files from the static directory on the filesystem
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Server started on http://localhost:" + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
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
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		return
	}
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
	// Dictionaries with weights for sentiment analysis
	positiveWords := map[string]int{
		"good": 2, "happy": 3, "love": 4, "great": 3,
		"excellent": 5, "positive": 2, "awesome": 4,
	}
	negativeWords := map[string]int{
		"bad": 2, "sad": 3, "hate": 4, "terrible": 5,
		"horrible": 5, "negative": 2,
	}
	neutralWords := map[string]int{
		"okay": 1, "average": 1, "fine": 1,
	}

	// Normalize text and tokenize
	text = strings.ToLower(text)
	words := strings.FieldsFunc(text, func(c rune) bool { return !unicode.IsLetter(c) })

	// Calculate sentiment scores
	posScore, negScore, neutralScore := 0, 0, 0
	for _, word := range words {
		if weight, found := positiveWords[word]; found {
			posScore += weight
		}
		if weight, found := negativeWords[word]; found {
			negScore += weight
		}
		if weight, found := neutralWords[word]; found {
			neutralScore += weight
		}
	}

	// Compute net sentiment score
	totalScore := posScore - negScore
	totalWords := len(words)

	// Normalize sentiment score
	normalizedScore := 0.0
	if totalWords > 0 {
		normalizedScore = float64(totalScore) / float64(totalWords)
	}

	// Categorize sentiment
	if normalizedScore > 0.2 {
		return "Positive"
	} else if normalizedScore < -0.2 {
		return "Negative"
	}
	return "Neutral"
}
