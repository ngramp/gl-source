package database

import (
	"gorm.io/gorm"
	"time"
)

// Address represents the address fields in the CSV data.
type Address struct {
	CompanyNumber string `gorm:"primarykey"`
	CareOf        string `csv:"RegAddress.CareOf"`
	POBox         string `csv:"RegAddress.POBox"`
	AddressLine1  string `csv:"RegAddress.AddressLine1"`
	AddressLine2  string `csv:"RegAddress.AddressLine2"`
	PostTown      string `csv:"RegAddress.PostTown"`
	County        string `csv:"RegAddress.County"`
	Country       string `csv:"RegAddress.Country"`
	PostCode      string `csv:"RegAddress.PostCode"`
}

// Accounts represents the accounts fields in the CSV data.
type Accounts struct {
	CompanyNumber   string `gorm:"primarykey"`
	AccountRefDay   int    `csv:"Accounts.AccountRefDay"`
	AccountRefMonth int    `csv:"Accounts.AccountRefMonth"`
	NextDueDate     string `csv:"Accounts.NextDueDate"`
	LastMadeUpDate  string `csv:"Accounts.LastMadeUpDate"`
	AccountCategory string `csv:"Accounts.AccountCategory"`
}

// SICCode represents the SIC code fields in the CSV data.
type SICCode struct {
	CompanyNumber string `gorm:"primarykey"`
	SicText1      string `csv:"SICCode.SicText_1"`
	SicText2      string `csv:"SICCode.SicText_2"`
	SicText3      string `csv:"SICCode.SicText_3"`
	SicText4      string `csv:"SICCode.SicText_4"`
}

// LimitedPartnerships represents the limited partnerships fields in the CSV data.
type LimitedPartnerships struct {
	CompanyNumber  string `gorm:"primarykey"`
	NumGenPartners int    `csv:"LimitedPartnerships.NumGenPartners"`
	NumLimPartners int    `csv:"LimitedPartnerships.NumLimPartners"`
}

// PreviousName represents the previous name fields in the CSV data.
type PreviousName struct {
	CompanyNumber string    `gorm:"primarykey"`
	CONDATE       time.Time `gorm:"primarykey" csv:"PreviousName.CONDATE"`
	CompanyName   string    `csv:"PreviousName.CompanyName"`
}

// Mortgages represents the mortgages section in the CSV data.
type Mortgages struct {
	CompanyNumber        string `gorm:"primarykey"`
	NumMortCharges       int    `csv:"Mortgages.NumMortCharges"`
	NumMortOutstanding   int    `csv:"Mortgages.NumMortOutstanding"`
	NumMortPartSatisfied int    `csv:"Mortgages.NumMortPartSatisfied"`
	NumMortSatisfied     int    `csv:"Mortgages.NumMortSatisfied"`
}

// Returns represents the returns section in the CSV data.
type Returns struct {
	CompanyNumber  string `gorm:"primarykey"`
	NextDueDate    string `csv:"Returns.NextDueDate"`
	LastMadeUpDate string `csv:"Returns.LastMadeUpDate"`
}

// Company represents the main structure for the CSV data.
type Company struct {
	CompanyNumber     string `gorm:"primaryKey" csv:"CompanyNumber"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
	CompanyName       string         `gorm:"index" csv:"CompanyName"`
	RegAddress        Address        `gorm:"foreignKey:CompanyNumber"`
	CompanyCategory   string         `csv:"CompanyCategory"`
	CompanyStatus     string         `csv:"CompanyStatus"`
	CountryOfOrigin   string         `csv:"CountryOfOrigin"`
	DissolutionDate   time.Time      `csv:"DissolutionDate"`
	IncorporationDate time.Time      `csv:"IncorporationDate"`
	//Accounts               Accounts `gorm:"foreignKey:CompanyNumber"`
	//Mortgages              Mortgages `gorm:"foreignKey:CompanyNumber"`
	SICCode SICCode `gorm:"foreignKey:CompanyNumber"`
	//LimitedPartnerships    LimitedPartnerships
	URI           string         `csv:"URI"`
	PreviousNames []PreviousName `gorm:"foreignKey:CompanyNumber"`
	//ConfStmtNextDueDate    string         `csv:"ConfStmtNextDueDate"`
	//ConfStmtLastMadeUpDate string         `csv:"ConfStmtLastMadeUpDate"`
	//Returns                Returns `gorm:"foreignKey:CompanyNumber"`
}

func (company *Company) MapCSVData(record []string) error {
	// Map the CSV data into the Company struct using struct tags.
	company.CompanyName = record[0]
	company.CompanyNumber = record[1]
	address := Address{
		CompanyNumber: record[1],
		CareOf:        record[2],
		POBox:         record[3],
		AddressLine1:  record[4],
		AddressLine2:  record[5],
		PostTown:      record[6],
		County:        record[7],
		Country:       record[8],
		PostCode:      record[9],
	}
	company.RegAddress = address
	company.CompanyCategory = record[10]
	company.CompanyStatus = record[11]
	company.CountryOfOrigin = record[12]
	dDate := record[13]
	if dDate != "" {
		company.DissolutionDate, _ = time.Parse("02/01/2006", dDate)
	}
	company.IncorporationDate, _ = time.Parse("02/01/2006", record[14])

	//company.Accounts.AccountRefDay, _ = strconv.Atoi(record[15])
	//company.Accounts.AccountRefMonth, _ = strconv.Atoi(record[16])
	//company.Accounts.NextDueDate = record[17]
	//company.Accounts.LastMadeUpDate = record[18]
	//company.Accounts.AccountCategory = record[19]

	//company.Returns.NextDueDate = record[20]
	//company.Returns.LastMadeUpDate = record[21]

	//company.Mortgages.NumMortCharges, _ = strconv.Atoi(record[22])
	//company.Mortgages.NumMortOutstanding, _ = strconv.Atoi(record[23])
	//company.Mortgages.NumMortPartSatisfied, _ = strconv.Atoi(record[24])
	//company.Mortgages.NumMortSatisfied, _ = strconv.Atoi(record[25])

	sicCode := SICCode{
		SicText1: record[26],
		SicText2: record[27],
		SicText3: record[28],
		SicText4: record[29],
	}
	company.SICCode = sicCode
	//company.LimitedPartnerships.NumGenPartners, _ = strconv.Atoi(record[30])
	//company.LimitedPartnerships.NumLimPartners, _ = strconv.Atoi(record[31])
	company.URI = record[32]

	//Loop to populate the PreviousName slice.
	var previousNames []PreviousName
	var conText, name string
	for i := 0; i < 10; i++ {
		conText = record[33+i*2]
		name = record[34+i*2]
		if name != "" && conText != "" {
			conDate, _ := time.Parse("02/01/2006", conText)
			previousName := PreviousName{
				CompanyNumber: record[1],
				CONDATE:       conDate,
				CompanyName:   name,
			}
			previousNames = append(previousNames, previousName)
		}
	}
	company.PreviousNames = previousNames
	//company.ConfStmtNextDueDate = record[53]
	//company.ConfStmtLastMadeUpDate = record[54]

	return nil
}
