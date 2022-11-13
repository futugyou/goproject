package versiond

import "github.com/futugyousuzu/k8s-controller-manager/client/versiond/scheme"

type clientSet struct {
}

func NewForConfig(c *rest.Config) (*EcsV1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &EcsV1Client{client}, nil

}

func setConfigDefaults(config *rest.Config) error {
	gv := ecsv1.GroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}
