package dto

// Store(ctx context.Context, data dto.NotificationRequest) error
// Finds(ctx context.Context, param dto.ListNotificationRequest) (response []domain.Notification, err error)
// Update(ctx context.Context, id string, data dto.NotificationRequest) error

type NotificationRequest struct {
	NotificationType   string `json:"notification_type"`
	TargetID           uint64 `json:"target_id"`
	Title              string `json:"title"`
	Message            string `json:"message"`
	AttachmentImageUrl string `json:"attachment_image_url"`
}

type ListNotificationRequest struct {
	UserID uint64 `json:"user_id" form:"user_id"`
	Offset uint64 `json:"offset" form:"offset"`
	Limit  uint64 `json:"limit" form:"limit"`
	IsRead bool   `json:"is_read" form:"is_read"`
}
