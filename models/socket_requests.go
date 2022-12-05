package models

type SUser struct {
	ID   string   `json:"id"`
	User UserJSON `json:"user"`
}

type SRRegister struct {
	Token string `json:"token"`
}

func (srRegister *SRRegister) Validate() []string {
	errors := Validator([]ValidatorClosure{
		func() string {
			return NullValidator("token", srRegister.Token)
		},
	})

	return errors
}