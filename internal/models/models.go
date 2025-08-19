package models

// ExtractedInfo holds the structured data we capture from PDFs and save to Excel.
type ExtractedInfo struct {
  SourceFiles              []string
  ContractDate             string
  ActDate                  string
  BuyerRepresentative      string
  SellerRepresentative     string
  BuyerPoAFrom            string
  BuyerPoATo              string
  SellerPoAFrom           string
  SellerPoATo             string

  // New fields
  ContractAmountRubles     string
  VehiclePassportIssueDate string
  VIN                      string
  BuyerCompany             string
  VehicleModel             string
  SellerCompanyDKP         string
} 