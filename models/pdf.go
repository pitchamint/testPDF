package models

type DetailCertificate struct {
	Name     string `from:"name"`
	LastName string `from:"last_name"`
	Course   string `from:"course"`
	Date     string `from:"date"`
	Template string `from:"template"`
}

type PersonUpload struct {
	NameFile string `form:"name_file"`
	File     interface{}
}
