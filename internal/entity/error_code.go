package entity

// Tabel untuk menyimpan definisi Error Code (Dictionary)
type ErrorReference struct {
	Code        string `gorm:"primaryKey;type:varchar(20)" json:"code"`
	MessageEN   string `json:"message_en"`
	MessageID   string `json:"message_id"`
	Description string `json:"description"` // Opsional, buat dokumentasi internal
}

// -- Struktur Response Schema (Wrapper) --

type APIResponse struct {
	ErrorSchema  ErrorSchema `json:"error_schema"`
	OutputSchema interface{} `json:"output_schema"` // Data dinamis (bisa null)
}

type ErrorSchema struct {
	ErrorCode    string       `json:"error_code"`
	ErrorMessage ErrorMessage `json:"error_message"`
}

type ErrorMessage struct {
	English    string `json:"english"`
	Indonesian string `json:"indonesian"`
}