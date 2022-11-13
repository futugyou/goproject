package v1

import (
	v1 "github.com/futugyousuzu/k8s-controller-manager/api/v1"
	"github.com/futugyousuzu/k8s-controller-manager/client/versiond/scheme"
)

type EcsBindingClient struct {
}
type EcsBindingInterface interface {
	Create(*v1.EcsBinding) (*v1.EcsBinding, error)
	Update(*v1.EcsBinding) (*v1.EcsBinding, error)
	Delete(string, *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	List(metav1.ListOptions) (*v1.EcsBindingList, error)
	Get(name string, options metav1.GetOptions) (*v1.EcsBinding, error)
	Watch(options metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.EcsBinding, err error)
}

func (c *ecsBinding) List(opts metav1.ListOptions) (*v1.EcsBindingList, error) {
	result := &v1.EcsBindingList{}
	err := c.client.Get().
		Resource("ecsbinding").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	for _, o := range result.Items {
		o.SetGroupVersionKind(v1.EscBindVersionKind)
	}
	return result, err
}
