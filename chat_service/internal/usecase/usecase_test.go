package usecase

import "chat_service/internal/domain"

type mockMessageRepo struct {
	Saved         []domain.Message
	SaveErr       error
	History       []domain.Message
	DeleteMessage error
}
