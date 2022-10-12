package hospital

import "time"

type Name struct {
	FullName  string
	Initial   string
	Firstname string
	Lastname  string
}

func NewName(init, first, last string) *Name {
	return &Name{
		FullName:  parseFullName(init, first, last),
		Initial:   init,
		Firstname: first,
		Lastname:  last,
	}
}

type Patient struct {
	BirthDate   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PassportId  *string
	NameEN      *Name
	NameTH      *Name
	NationalId  *string
	Id          string
	Nationality string
	PhoneNumber string
	BloodType   BloodType
	Height      float64
	Weight      float64
}

type Doctor struct {
	Id            string
	NameEN        *Name
	NameTH        *Name
	Position      string
	ProfilePicURL string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Username      string
}

type InvoiceOverview struct {
	CreatedAt     time.Time
	AppointmentID string
	PatientID     string
	Id            int
	Total         float64
	Paid          bool
}

type AppointmentOverview struct {
	Id            string            `json:"id"`
	Detail        string            `json:"detail"`
	StartDateTime time.Time         `json:"start_date_time"`
	EndDateTime   time.Time         `json:"end_date_time"`
	Status        AppointmentStatus `json:"status"`
	Doctor        DoctorOverview    `json:"doctor"`
	Patient       PatientOverview   `json:"patient"`
}
type DoctorOverview struct {
	ID            string `json:"id"`
	FullName      string `json:"full_name"`
	Position      string `json:"position"`
	ProfilePicURL string `json:"profile_pic_url"`
}
type PatientOverview struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
}

type Appointment struct {
	Id              string            `json:"id"`
	PatientID       string            `json:"patient_id"`
	StartDateTime   time.Time         `json:"start_date_time"`
	EndDateTime     time.Time         `json:"end_date_time"`
	NextAppointment *time.Time        `json:"next_appointment"`
	Detail          string            `json:"detail"`
	Status          AppointmentStatus `json:"status"`
	Doctor          DoctorOverview    `json:"doctor"`
	Invoice         *Invoice          `json:"invoice"`
	Prescriptions   []*Prescription   `json:"prescriptions"`
}
type Invoice struct {
	InvoiceItems     []*InvoiceItem     `json:"invoice_items"`
	InvoiceDiscounts []*InvoiceDiscount `json:"invoice_discounts"`
	Id               int                `json:"id"`
	Total            float64            `json:"total"`
	Paid             bool               `json:"paid"`
}
type InvoiceItem struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}
type InvoiceDiscount struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}
type Prescription struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PictureURL  string `json:"picture_url"`
	Amount      int    `json:"amount"`
}
