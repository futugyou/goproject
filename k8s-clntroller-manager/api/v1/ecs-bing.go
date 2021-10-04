package v1

type EcsBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EcsBindSpec   `json:"spec,omitempty"`
	Status EcsBindStatus `json:"status,omitempty"`
}

type EcsBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []EcsBinding `json:"items"`
}

type EcsBindSpec struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	NodeName string `json:"node_name"`
	InnerIp  string `json:"inner_ip"`
}

type EcsBindStatus struct {
}

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "hopegi.com", Version: "v1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}
	//SchemeBuilder = runtime.NewSchemeBuilder(AddKnownTypes)

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme

	EscBindVersionKind = schema.GroupVersionKind{Group: GroupVersion.Group, Version: GroupVersion.Version, Kind: "EcsBinding"}
)

func init() {

	SchemeBuilder.Register(&EcsBinding{}, &EcsBindingList{})

}
