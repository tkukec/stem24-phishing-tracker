package domain

import (
	"github.com/asseco-voice/agent-management/domain/models"
)

type SipPayload struct {
	activity          *models.Activity
	action            string
	number            string
	agentId           string
	hangupReason      string
	hangupExplanation string
}

func NewSipPayload(
	activity *models.Activity,
	action string,
	number string,
	agentId string,
	hangupReason string,
	hangupExplanation string,
) *SipPayload {
	return &SipPayload{
		activity:          activity,
		action:            action,
		number:            number,
		agentId:           agentId,
		hangupReason:      hangupReason,
		hangupExplanation: hangupExplanation,
	}
}

func (s SipPayload) Activity() *models.Activity {
	return s.activity
}

func (s SipPayload) Action() string {
	return s.action
}

func (s SipPayload) Number() string {
	return s.number
}

func (s SipPayload) AgentId() string {
	return s.agentId
}

func (s SipPayload) HangupExplanation() string {
	return s.hangupExplanation
}

func (s SipPayload) HangupReason() string {
	return s.hangupReason
}
