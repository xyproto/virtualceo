package main

import (
	"fmt"
	"log"
	"time"

	"github.com/xyproto/ollamaclient/v2"
)

func processMessage(msg Message) error {
	oc := ollamaclient.New("llama3.1")
	oc.Verbose = true

	err := oc.PullIfNeeded(true)
	if err != nil {
		return fmt.Errorf("failed to pull model: %w", err)
	}

	if found, err := oc.Has("llama3.1"); err != nil || !found {
		return fmt.Errorf("expected to have 'llama3.1' model downloaded, but it's not present")
	}

	// Improved instructions for the virtual CEO
	oc.SetSystemPrompt(`You are a helpful and intelligent virtual CEO. Your goal is to manage tasks efficiently and respond to messages appropriately.
When you receive a message:
- If the message asks to send an email, use the "send_email" tool.
- If the message asks to log an event, use the "add_to_work_log" tool.
- If the message asks to do nothing, use the "do_nothing" tool.
- If the message requires a reply, craft a clever and helpful response using the "reply_email" tool.
- Always log significant actions with proper context.`)
	oc.SetRandom()

	toolSendEmail, toolDoNothing, toolAddToWorkLog, toolScheduleMeeting, toolGenerateReport, toolReplyEmail := GetTools()

	oc.SetTool(toolSendEmail)
	oc.SetTool(toolDoNothing)
	oc.SetTool(toolAddToWorkLog)
	oc.SetTool(toolScheduleMeeting)
	oc.SetTool(toolGenerateReport)
	oc.SetTool(toolReplyEmail)

	// Call the send_email tool
	emailCall := fmt.Sprintf(`send_email {"to": "%s", "subject": "%s", "body": "%s"}`, msg.From, msg.Subject, msg.Body)

	log.Printf("Sending to llama3.1: %s", emailCall)
	generatedOutput := oc.MustGetChatResponse(emailCall)
	if generatedOutput.Error != "" {
		return fmt.Errorf(generatedOutput.Error)
	}

	// Process each tool call
	for _, toolCall := range generatedOutput.ToolCalls {
		log.Printf("Tool called: %s", toolCall.Function.Name)
		// Invoke the corresponding function based on the tool call
		switch toolCall.Function.Name {
		case "send_email":
			to := toolCall.Function.Arguments["to"].(string)
			subject := toolCall.Function.Arguments["subject"].(string)
			body := toolCall.Function.Arguments["body"].(string)
			log.Printf("Calling sendEmail with arguments: to=%s, subject=%s, body=%s", to, subject, body)
			sendEmail(to, subject, body)
		case "do_nothing":
			log.Println("Calling doNothing")
			doNothing()
		case "add_to_work_log":
			timestamp := toolCall.Function.Arguments["timestamp"].(string)
			summary := toolCall.Function.Arguments["summary"].(string)
			log.Printf("Calling addToWorkLog with arguments: timestamp=%s, summary=%s", timestamp, summary)
			addToWorkLog(timestamp, summary)
		case "schedule_meeting":
			participants := convertToStringSlice(toolCall.Function.Arguments["participants"])
			time := toolCall.Function.Arguments["time"].(string)
			log.Printf("Calling scheduleMeeting with arguments: participants=%v, time=%s", participants, time)
			scheduleMeeting(participants, time)
		case "generate_report":
			criteria := toolCall.Function.Arguments["criteria"].(string)
			log.Printf("Calling generateReport with argument: criteria=%s", criteria)
			generateReport(criteria)
		case "reply_email":
			to := toolCall.Function.Arguments["to"].(string)
			subject := toolCall.Function.Arguments["subject"].(string)
			body := toolCall.Function.Arguments["body"].(string)
			log.Printf("Calling replyEmail with arguments: to=%s, subject=%s, body=%s", to, subject, body)
			replyEmail(to, subject, body)
		default:
			log.Printf("Unknown tool call: %s", toolCall.Function.Name)
		}
	}

	// Log the message sending event without including the timestamp in the summary
	logEvent := fmt.Sprintf("Message from %s with subject %s", msg.From, msg.Subject)
	logCall := fmt.Sprintf(`add_to_work_log {"timestamp": "%s", "summary": "%s"}`, time.Now().Format(time.RFC3339), logEvent)

	generatedOutput = oc.MustGetChatResponse(logCall)
	if generatedOutput.Error != "" {
		return fmt.Errorf(generatedOutput.Error)
	}

	log.Println("Tool function add_to_work_log called")

	return nil
}

// Helper function to convert interface{} to []string
func convertToStringSlice(input interface{}) []string {
	if input == nil {
		return nil
	}
	interfaceSlice := input.([]interface{})
	stringSlice := make([]string, len(interfaceSlice))
	for i, v := range interfaceSlice {
		stringSlice[i] = v.(string)
	}
	return stringSlice
}
