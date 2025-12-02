package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CSPReport struct {
	ID          string `gorm:"primaryKey;default:uuid_generate_v4()"`
	Age         int
	Body        datatypes.JSONType[ReportBody]
	Disposition string `gorm:"type:varchar(20)"`
	Directive   string `gorm:"type:varchar(50)"`
	BlockedURL  string `gorm:"type:text"`
	DocumentURL string `gorm:"type:text"`
	ProjectID   string `gorm:"type:uuid"`
	Type        string `gorm:"type:varchar(20)"`
	URL         string `gorm:"type:varchar(100)"`
	SourceIP    string `gorm:"type:varchar(20)"`
	UserAgent   string `json:"user_agent"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type ReportBody struct {
	BlockedURL         string `json:"blockedURL"`
	ColumnNumber       int    `json:"columnNumber"`
	Disposition        string `json:"disposition"`
	DocumentURL        string `json:"documentURL"`
	EffectiveDirective string `json:"effectiveDirective"`
	LineNumber         int    `json:"lineNumber"`
	OriginalPolicy     string `json:"originalPolicy"`
	Referrer           string `json:"referrer"`
	Sample             string `json:"sample"`
	SourceFile         string `json:"sourceFile"`
	ScriptSample       string `json:"scriptSample"`
}

type CSPReportCreateDTO struct {
	Age       int        `json:"age"`
	Body      ReportBody `json:"body"`
	Type      string     `json:"type"`
	URL       string     `json:"url"`
	UserAgent string     `json:"user_agent"`
	SourceIP  string     `json:"-"`
}

type CSPReportFetchDTO struct {
	Body        datatypes.JSONType[ReportBody] `json:"body"`
	Directive   string                         `json:"directive"`
	URL         string                         `json:"url"`
	BlockedURL  string                         `json:"blockedURL"`
	DocumentURL string                         `json:"documentURL"`
	Disposition string                         `json:"disposition"`
	SourceIP    string                         `json:"sourceIP"`
	UserAgent   string                         `json:"userAgent"`
	LastSeen    time.Time                      `json:"lastSeen"`
	Count       int                            `json:"count"`
}

func (r *CSPReportCreateDTO) ToCSPReport() *CSPReport {
	c := &CSPReport{
		Age:       r.Age,
		Type:      r.Type,
		URL:       r.URL,
		UserAgent: r.UserAgent,
		SourceIP:  r.SourceIP,
	}
	c.Body = datatypes.NewJSONType(r.Body)
	c.Directive = r.Body.EffectiveDirective
	c.Disposition = r.Body.Disposition
	c.BlockedURL = r.Body.BlockedURL
	c.DocumentURL = r.Body.DocumentURL
	return c
}
