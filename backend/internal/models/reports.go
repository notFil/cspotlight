package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CSPReport struct {
	ID         string `gorm:"primaryKey;default:uuid_generate_v4()"`
	Age        int
	ReportBody datatypes.JSONType[ReportBody]
	ProjectID  string `gorm:"type:uuid"`
	Type       string `gorm:"type:varchar(20)"`
	URL        string `gorm:"type:varchar(100)"`
	SourceIP   string `gorm:"type:varchar(20)"`
	UserAgent  string `json:"user_agent"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
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
	Age        int        `json:"age"`
	ReportBody ReportBody `json:"reportBody"`
	Type       string     `json:"type"`
	URL        string     `json:"url"`
	UserAgent  string     `json:"user_agent"`
	SourceIP   string     `json:"-"`
}

type CSPReportFetchDTO struct {
	ReportBody ReportBody `json:"reportBody"`
	Directive  string     `json:"directive"`
	URL        string     `json:"url"`
	SourceIP   string     `json:"sourceIP"`
	UserAgent  string     `json:"user_agent"`
	Count      int        `json:"count"`
}

func (r *CSPReportCreateDTO) ToCSPReport() *CSPReport {
	c := &CSPReport{
		Age:       r.Age,
		Type:      r.Type,
		URL:       r.URL,
		UserAgent: r.UserAgent,
		SourceIP:  r.SourceIP,
	}
	c.ReportBody = datatypes.NewJSONType(r.ReportBody)
	return c
}

func (c *CSPReport) ToFetchDTO() *CSPReportFetchDTO {
	r := &CSPReportFetchDTO{
		URL:       c.URL,
		UserAgent: c.UserAgent,
		SourceIP:  c.SourceIP,
	}
	r.ReportBody = c.ReportBody.Data()
	r.Directive = r.ReportBody.EffectiveDirective
	return r
}
