package main

import (
    "context"
    "fmt"
    "log"
    "google.golang.org/genai"
)

func main100() {
    ctx := context.Background()
    // The client gets the API key from the environment variable `GEMINI_API_KEY`.
    client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: "AIzaSyB9aVrch2a0sHPjzIdrBTpbmC4JHQ9BHKM"})
    if err != nil {
        log.Fatal(err)
    }

    result, err := client.Models.GenerateContent(
        ctx,
        "gemini-2.5-flash",
        genai.Text("سلام می توانی صدا را بررسی کنی و فایل صدا بفرستی "),
        nil,
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result.())
}