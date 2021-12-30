package repository

type DocumentsTypes string

const (
	CNH       DocumentsTypes = "CNH"
	RG        DocumentsTypes = "RG"
	Passaport DocumentsTypes = "Passaport"
)

func (a DocumentsTypes) String() string {
	return string(a)
}

type AddressTypes string

const (
	Personal AddressTypes = "Personal"
	Company  AddressTypes = "Company"
	Other    AddressTypes = "Other"
)

func (a AddressTypes) String() string {
	return string(a)
}

type PeopleType string

const (
	People      PeopleType = "People"
	LegalEntity PeopleType = "Company"
)

func (a PeopleType) String() string {
	return string(a)
}

type AccountBankType string

const (
	CC AccountBankType = "Conta Corrente"
	CP AccountBankType = "Conta Poupaca"
	CI AccountBankType = "Conta Investimento"
)

func (a AccountBankType) String() string {
	return string(a)
}

type MaritalStatus string

const (
	Simgle  MaritalStatus = "Solteiro(a)"
	Married MaritalStatus = "Casado(a)"
	Widow   MaritalStatus = "Viúva"
)

func (a MaritalStatus) String() string {
	return string(a)
}
