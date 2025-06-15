package domain

// type ValidateUserDataResp struct {
// 	NationalID   string `json:"national_id"`
// 	FullName     string `json:"full_name"`
// 	LegalName    string `json:"legal_name"`
// 	BirthOfPlace string `json:"birth_of_place"`
// 	BirthOfDate  string `json:"birth_of_date"`
// 	Salary       string `json:"salary"`
// }

type ValidateUserReq struct {
	NationalID      string `json:"national_id"`
	FullName        string `json:"full_name"`
	LegalName       string `json:"legal_name"`
	BirthOfPlace    string `json:"birth_of_place"`
	BirthOfDate     string `json:"birth_of_date"`
	Salary          string `json:"salary"`
	NationalIDPhoto []byte `json:"national_id_photo"`
	UserPhoto       []byte `json:"user_photo"`
}

type KYCValidateNationalIDReq struct {
	NationalID  string `json:"nik"`
	LegalName   string `json:"name"`
	DateOfBirth string `json:"birth_date"`
	ReferenceID string `json:"reference_id"`
}

type KYCData struct {
	NationalID  bool   `json:"nik"`
	LegalName   bool   `json:"name"`
	DateOfBirth bool   `json:"birth_date"`
	Gender      bool   `json:"gender,omitempty"`
	SalaryUper  string `json:"salary_upper,omitempty"`
	SalaryLower string `json:"salary_lower,omitempty"`
	ReferenceID string `json:"reference_id"`
}
type KYCValidateNationalIDResp struct {
	Message string  `json:"message"`
	Data    KYCData `json:"data"`
}

type KYCValidateSalaryResp = KYCValidateNationalIDResp

type KYCValidateSalaryReq struct {
	NationalID string `json:"national_id"`
	LegalName  string `json:"name"`
	// Gender      string `json:"gender"`
	Salary      string `json:"salary"`
	ReferenceID string `json:"reference_id"`
}

type KYCValidatePhotoReq struct {
	NationalID      string `json:"national_id"`
	LegalName       string `json:"name"`
	NationalIDPhoto []byte `json:"national_id_photo"`
	UserPhoto       []byte `json:"user_photo"`
	ReferenceID     string `json:"reference_id"`
}

type KYCValidatePhotoResp struct {
	Message string `json:"message"`
	Data    struct {
		Status string `json:"status"`
	}
}

type CreateLoanReq struct {
	NationalID      string  `json:"national_id"`
	LegalName       string  `json:"legal_name"`
	ContractNumber  string  `json:"contract_number"`
	Amount          float64 `json:"amount"`
	DownPayment     float64 `json:"down_payment"`
	OTRAmount       float64
	PrincipalAmount float64
	AssetName       string `json:"asset_name"`
	LoanTypeID      int    `json:"loan_type_id"`
	LimitTypeID     int    `json:"limit_type_id"`
	Status          string `json:"status"`
	StartDate       string `json:"start_date"`
	InterestRate    int    `json:"interest_rate"`
	NationalIDPhoto []byte `json:"national_id_photo"`
	UserPhoto       []byte `json:"user_photo"`
	Salary          string `json:"salary"`
	BirthOfDate     string `json:"birth_of_date"`
}

type GeneralResponse struct {
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}
