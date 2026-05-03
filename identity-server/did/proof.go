package did

type ISignatureProof interface {
	GetType() string
	GetVerificationMethod() string
	GetTransformationMethod() string
	GetHashingMethod() string
	ComputeProof(proof DataIntegrityProof, payload []byte, asymmetricKey IAsymmetricKey, alg string) error
	GetSignature(proof DataIntegrityProof) ([]byte, error)
}
