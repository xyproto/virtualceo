package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/xyproto/ollamaclient/v2"
)

func GetTools() (ollamaclient.Tool, ollamaclient.Tool, ollamaclient.Tool, ollamaclient.Tool, ollamaclient.Tool, ollamaclient.Tool) {
	var toolSendEmail, toolDoNothing, toolAddToWorkLog, toolScheduleMeeting, toolGenerateReport, toolReplyEmail ollamaclient.Tool

	// Define the tools
	json.Unmarshal(json.RawMessage(`{
		"type": "function",
		"function": {
			"name": "send_email",
			"description": "Send an email to someone",
			"parameters": {
				"type": "object",
				"properties": {
					"to": {
						"type": "string",
						"description": "The recipient of the email"
					},
					"subject": {
						"type": "string",
						"description": "The subject of the email"
					},
					"body": {
						"type": "string",
						"description": "The body of the email"
					}
				},
				"required": ["to", "subject", "body"]
			}
		}
	}`), &toolSendEmail)

	json.Unmarshal(json.RawMessage(`{
		"type": "function",
		"function": {
			"name": "do_nothing",
			"description": "A function that does nothing and just prints a message"
		}
	}`), &toolDoNothing)

	// Updated tool: add_to_work_log with HELLO print statement
	json.Unmarshal(json.RawMessage(`{
		"type": "function",
		"function": {
			"name": "add_to_work_log",
			"description": "Logs an event with a summary and timestamp",
			"parameters": {
				"type": "object",
				"properties": {
					"timestamp": {
						"type": "string",
						"description": "The timestamp of the event"
					},
					"summary": {
						"type": "string",
						"description": "A summary of the event without timestamp"
					}
				},
				"required": ["timestamp", "summary"]
			}
		}
	}`), &toolAddToWorkLog)

	// New tool: Schedule a meeting
	json.Unmarshal(json.RawMessage(`{
		"type": "function",
		"function": {
			"name": "schedule_meeting",
			"description": "Schedule a meeting with specified participants",
			"parameters": {
				"type": "object",
				"properties": {
					"participants": {
						"type": "array",
						"items": {
							"type": "string"
						},
						"description": "List of participants for the meeting"
					},
					"time": {
						"type": "string",
						"description": "The time for the meeting"
					}
				},
				"required": ["participants", "time"]
			}
		}
	}`), &toolScheduleMeeting)

	// New tool: Generate a report
	json.Unmarshal(json.RawMessage(`{
		"type": "function",
		"function": {
			"name": "generate_report",
			"description": "Generate a report based on given criteria",
			"parameters": {
				"type": "object",
				"properties": {
					"criteria": {
						"type": "string",
						"description": "The criteria for the report"
					}
				},
				"required": ["criteria"]
			}
		}
	}`), &toolGenerateReport)

	// New tool: Reply email
	json.Unmarshal(json.RawMessage(`{
		"type": "function",
		"function": {
			"name": "reply_email",
			"description": "Craft and send a clever reply email to the customer",
			"parameters": {
				"type": "object",
				"properties": {
					"to": {
						"type": "string",
						"description": "The recipient of the reply email"
					},
					"subject": {
						"type": "string",
						"description": "The subject of the reply email"
					},
					"body": {
						"type": "string",
						"description": "The body of the reply email"
					}
				},
				"required": ["to", "subject", "body"]
			}
		}
	}`), &toolReplyEmail)

	return toolSendEmail, toolDoNothing, toolAddToWorkLog, toolScheduleMeeting, toolGenerateReport, toolReplyEmail
}

func sendEmail(to, subject, body string) {
	log.Printf("sendEmail called with to=%s, subject=%s, body=%s", to, subject, body)
}

func doNothing() {
	log.Println("doNothing called")
}

func addToWorkLog(timestamp, summary string) {
	log.Printf("addToWorkLog called with timestamp=%s, summary=%s", timestamp, summary)
	fmt.Println("HELLO")
}

func scheduleMeeting(participants []string, time string) {
	log.Printf("scheduleMeeting called with participants=%v, time=%s", participants, time)
}

func generateReport(criteria string) {
	log.Printf("generateReport called with criteria=%s", criteria)
}

func replyEmail(to, subject, body string) {
	log.Printf("replyEmail called with to=%s, subject=%s, body=%s", to, subject, body)
}
