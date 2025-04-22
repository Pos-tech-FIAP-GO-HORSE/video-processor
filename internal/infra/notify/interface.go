package notify

import "video-processor/internal/service"

type NotifierInterface interface {
	NotifyResult(event service.NotificationEvent) error
}
