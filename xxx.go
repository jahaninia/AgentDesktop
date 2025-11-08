package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
)

func main() {
	// 🔐 کلید API
	apiKey := "AIzaSyB9aVrch2a0sHPjzIdrBTpbmC4JHQ9BHKM"


	// 🎤 فایل صوتی ورودی را بخوان
	audioBytes, err := ioutil.ReadFile("Recording.wav") // یا mp3
	if err != nil {
		panic(err)
	}
	audioBase64 := base64.StdEncoding.EncodeToString(audioBytes)

	// 🧠 پرسش فارسی برای مدل
	promptText := "متن این صدا را تحلیل کن و به زبان فارسی با صدای مهربان پاسخ بده."

	// 🧱 بدنه درخواست برای Gemini
	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"role": "user",
				"parts": []map[string]interface{}{
					{"mime_type": "audio/wav", "data": audioBase64},
					{"text": promptText},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			// در این بخش به مدل می‌گوئیم که خروجی صوتی برگرداند
			"responseModalities": []string{"AUDIO"},
			"audioConfig": map[string]interface{}{
				"voice":            "basic", // صدای ساده برای پاسخ
				"sampleRateHertz":  16000,
				"audioEncoding":    "wav", // می‌تواند mp3 هم باشد
			},
		},
	}

	body, _ := json.Marshal(payload)

	// 🚀 ارسال به Gemini API
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent"
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respBody)) // برای بررسی پاسخ خام

	// 🔍 استخراج صوت خروجی از پاسخ JSON
	var result map[string]interface{}
	json.Unmarshal(respBody, &result)

	// این مسیر در بدنهٔ پاسخ ممکن است تغییر کند بسته به مدل انتخابی
	parts := result["candidates"].([]interface{})[0].(map[string]interface{})["content"].(map[string]interface{})["parts"].([]interface{})
	audioData := parts[0].(map[string]interface{})["data"].(string)

	audioDecoded, _ := base64.StdEncoding.DecodeString(audioData)

	// 💾 ذخیره فایل خروجی
	ioutil.WriteFile("reply.wav", audioDecoded, 0644)
	fmt.Println("✅ پاسخ صوتی در فایل reply.wav ذخیره شد.")

	// 🔊 پخش فایل خروجی در لینوکس/macOS
	exec.Command("play", "reply.wav").Run()   // برای Linux (اگر sox نصب داری)
	// یا exec.Command("afplay", "reply.wav").Run() // برای macOS
}
