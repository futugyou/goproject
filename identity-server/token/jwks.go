package token

type JwkModel struct {
	ID      string `bson:"_id" json:"id"`
	Payload string `bson:"payload" json:"payload"`
}

func (JwkModel) GetType() string {
	return "jwks"
}