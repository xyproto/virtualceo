package main

import (
	"fmt"
	"log"
	"net/http"
)

type Message struct {
	Subject string `json:"subject"`
	From    string `json:"from"`
	Body    string `json:"body"`
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/newmessage", newMessageHandler)

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Message Form</title>
    <script src="https://unpkg.com/htmx.org@2.0.2"></script>
    <style>
        @import url('https://fonts.googleapis.com/css2?family=Montserrat:wght@400;600&family=Roboto:wght@300;500&display=swap');
        body {
            font-family: 'Montserrat', sans-serif;
            background-color: #f0f2f5;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            margin: 0;
        }
        form {
            background: white;
            padding: 20px;
            border-radius: 10px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
            animation: fadeIn 1s ease-in-out;
            width: 100%;
            max-width: 500px;
        }
        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(-20px); }
            to { opacity: 1; transform: translateY(0); }
        }
        label {
            display: block;
            margin-bottom: 8px;
            font-weight: 600;
            color: #333;
        }
        input, textarea, button {
            width: 100%;
            padding: 10px;
            margin-bottom: 20px;
            border: 1px solid #ccc;
            border-radius: 5px;
            font-family: 'Roboto', sans-serif;
        }
        textarea {
            resize: vertical;
            height: 150px; /* Increased height for larger message area */
        }
        button {
            background-color: #4CAF50;
            color: white;
            font-size: 16px;
            border: none;
            cursor: pointer;
            transition: background-color 0.3s;
        }
        button:hover {
            background-color: #45a049;
        }
    </style>
</head>
<body>
    <form hx-post="/api/newmessage" hx-trigger="submit" hx-target="#response" hx-swap="innerHTML">
        <label for="subject">Subject:</label>
        <input type="text" id="subject" name="subject" required>
        <label for="from">From:</label>
        <input type="email" id="from" name="from" required>
        <label for="body">Message:</label>
        <textarea id="body" name="body" required></textarea>
        <button type="submit">Send</button>
    </form>
    <div id="response"></div>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func newMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	msg := Message{
		Subject: r.FormValue("subject"),
		From:    r.FormValue("from"),
		Body:    r.FormValue("body"),
	}

	err = processMessage(msg)
	if err != nil {
		http.Error(w, "Failed to process message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Message sent successfully!"))
}
