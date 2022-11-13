package ecsbind

type Interface interface {
	// V1 provides access to shared informers for resources in V1.
	V1() v1.Interface
}
