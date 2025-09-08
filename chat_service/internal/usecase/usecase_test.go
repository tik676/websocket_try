package usecase

import (
	"chat_service/internal/domain"
	"fmt"
	"testing"
)

type mockMessageRepo struct {
	Saved      []domain.Message
	SaveErr    error
	History    []domain.Message
	HistoryErr error
	DeleteErr  error
}

func (m *mockMessageRepo) SaveMessage(msg domain.Message) (domain.Message, error) {
	m.Saved = append(m.Saved, msg)
	if m.SaveErr != nil {
		return domain.Message{}, m.SaveErr
	}
	return msg, nil

}

func (m *mockMessageRepo) MessageHistory(limit, offset int64) ([]domain.Message, error) {
	if m.HistoryErr != nil {
		return []domain.Message{}, m.HistoryErr
	}
	return m.History, nil
}

func (m *mockMessageRepo) DeleteMessage(id int64) error {
	return m.DeleteErr
}

func TestSendMessage(t *testing.T) {
	mockRepo := &mockMessageRepo{}
	uc := NewUseCase(mockRepo)

	msg := domain.Message{
		UserID:   1,
		Username: "user",
		Content:  "test message",
		Role:     "user",
	}

	result, err := uc.SendMessage(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Content != msg.Content {
		t.Errorf("expected content %s, got %s", msg.Content, result.Content)
	}

	if len(mockRepo.Saved) != 1 {
		t.Error("message should be saved")
	}

	_, err = uc.SendMessage(domain.Message{UserID: 1, Content: ""})
	if err == nil || err.Error() != "error message can't be empty" {
		t.Error("expected validation error for empty content")
	}

	_, err = uc.SendMessage(domain.Message{UserID: 0, Content: "111"})
	if err == nil || err.Error() != "user id cannot be zero" {
		t.Errorf("expected validation error for zero id")
	}
}

func TestGetMessages(t *testing.T) {
	testMessages := []domain.Message{
		{ID: 1, Content: "msg1"},
		{ID: 2, Content: "msg2"},
	}

	mockRepo := &mockMessageRepo{History: testMessages}
	uc := NewUseCase(mockRepo)

	msgs, err := uc.GetMessages(10, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(msgs) != 2 {
		t.Errorf("expected 2 messages got: %d", len(msgs))
	}

	if msgs[0].Content != "msg1" {
		t.Errorf("expected firs message content `msg1`,got %s", msgs[0].Content)
	}

	msgs, err = uc.GetMessages(0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	msgs, err = uc.GetMessages(10, -5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mockRepo.HistoryErr = fmt.Errorf("some db error")
	_, err = uc.GetMessages(10, 0)
	if err == nil {
		t.Fatalf("expected error from repo, got nil")
	}

}

func TestDeleteMessage(t *testing.T) {
	mockRepo := &mockMessageRepo{}
	uc := NewUseCase(mockRepo)

	err := uc.DeleteMessage(1)
	if err != nil {
		t.Fatalf("expected nil, got: %v", err)
	}

	mockRepo.DeleteErr = fmt.Errorf("some db error")
	err = uc.DeleteMessage(2)
	if err == nil {
		t.Errorf("expected error from repo, got nil")
	}
}
