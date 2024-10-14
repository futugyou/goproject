package platform

type PlatformProvider interface {
	privatePlatformProvider()
	String() string
}

type platformProvider string

func (c platformProvider) privatePlatformProvider() {}

func (c platformProvider) String() string {
	return string(c)
}

const (
	PlatformProviderVercel   platformProvider = "vercel"
	PlatformProviderGithub   platformProvider = "github"
	PlatformProviderCircleci platformProvider = "circleci"
	PlatformProviderOther    platformProvider = "other"
)

func GetPlatformProvider(rType string) PlatformProvider {
	switch rType {
	case "vercel":
		return PlatformProviderVercel
	case "github":
		return PlatformProviderGithub
	case "circleci":
		return PlatformProviderCircleci
	default:
		return PlatformProviderOther
	}
}
