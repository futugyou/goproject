package scheme

import "fmt"

var scheme = runtime.NewScheme()
var Codecs = serializer.NewCodecFactory(scheme)
var ParameterCodec = runtime.NewParameterCodec(scheme)

func init() {
	metav1.AddToGroupVersion(scheme, schema.GroupVersion{Version: "v1"})
	if err := AddToScheme(scheme); err != nil {
		fmt.Println("error to AddToScheme ", err)
	}
}

func AddToScheme(scheme *runtime.Scheme) error {
	return ecsv1.AddToScheme(scheme)
}
