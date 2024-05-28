package platform

import "github.com/futugyou/infr-project/domain"

type IPlatformRepository interface {
	domain.IRepository[*Platform]
}
