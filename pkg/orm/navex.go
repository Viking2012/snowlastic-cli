package types

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/alexander-orban/icm_goapi/orm"
)

type CaseFile struct {
	FileID   string `json:"file_id"`
	FileName string `json:"file_name"`
}

func (c *CaseFile) IsICMEntity() bool { return true }
func (c *CaseFile) Scan(src interface{}) error {
	switch v := src.(type) {
	case string:
		return json.Unmarshal([]byte(v), c)
	case []byte:
		return json.Unmarshal(v, c)
	default:
		return errors.New("invalid sql return type for CaseFile")
	}
}

type CaseQuestion struct {
	QuestionID    string       `json:"followup_question_id"`
	IsComment     bool         `json:"is_comment"`
	QuestionText  string       `json:"followup_question_text"`
	QuestionDate  orm.NullTime `json:"followup_question_date"`
	QuestionAsker string       `json:"followup_question_asker"`
	AnswerText    string       `json:"answer_text"`
	AnswerDate    orm.NullTime `json:"answer_date"`
}

func (c *CaseQuestion) IsICMEntity() bool { return true }
func (c *CaseQuestion) Scan(src interface{}) error {
	switch v := src.(type) {
	case string:
		return json.Unmarshal([]byte(v), c)
	case []byte:
		return json.Unmarshal(v, c)
	default:
		return errors.New("invalid sql return type for CaseQuestion")
	}
}

type CaseNote struct {
	NoteID   string       `json:"note_id"`
	NoteText string       `json:"note_text"`
	NoteDate orm.NullTime `json:"note_date"`
}

func (c *CaseNote) IsICMEntity() bool { return true }
func (c *CaseNote) Scan(src interface{}) error {
	switch v := src.(type) {
	case string:
		return json.Unmarshal([]byte(v), c)
	case []byte:
		return json.Unmarshal(v, c)
	default:
		return errors.New("invalid sql return type for CaseNote")
	}
}

type CaseParticipant struct {
	ParticipantID           string `json:"participant_id"`
	ParticipantPrefix       string `json:"participant_prefix"`
	ParticipantName         string `json:"participant_name"`
	ParticipantHRID         string `json:"participant_hr_id"`
	ParticipantRelationship string `json:"participant_relationship"`
	ParticipantRole         string `json:"participant_role"`
	Outcome                 string `json:"outcome"`
}

func (c *CaseParticipant) IsICMEntity() bool { return true }
func (c *CaseParticipant) Scan(src interface{}) error {
	switch v := src.(type) {
	case string:
		return json.Unmarshal([]byte(v), c)
	case []byte:
		return json.Unmarshal(v, c)
	default:
		return errors.New("invalid sql return type for CaseParticipant")
	}
}

type Case struct {
	CaseNumber           string            `json:"case_number"`
	Alert                string            `json:"alert_level"`
	CaseBranchNumber     string            `json:"branch_number"`
	CaseDetails          string            `json:"case_details"`
	CaseType             string            `json:"case_type"`
	ActionTaken          string            `json:"action_taken"`
	IncidentDate         orm.NullTime      `json:"incident_date"`
	ClosureDate          orm.NullTime      `json:"closure_date"`
	DueDate              orm.NullTime      `json:"due_date"`
	OpenDate             orm.NullTime      `json:"open_date"`
	ReportDate           orm.NullTime      `json:"reported_date"`
	CaseCity             string            `json:"case_city"`
	CaseStateProvince    string            `json:"case_state_province"`
	CasePostalCode       string            `json:"case_postal_code"`
	CaseCountry          string            `json:"case_country"`
	CaseRegion           string            `json:"case_region"`
	BusinessArea         string            `json:"business_case"`
	Division             string            `json:"division"`
	PrimaryIssue         string            `json:"primary_issue"`
	PrimaryIssueLayer1   string            `json:"primary_issue_layer1"`
	PrimaryIssueLayer2   string            `json:"primary_issue_layer2"`
	PrimaryIssueLayer3   string            `json:"primary_issue_layer3"`
	PrimaryOutcome       string            `json:"primary_outcome"`
	SecondaryIssue       string            `json:"secondary_issue"`
	SecondaryIssueLayer1 string            `json:"secondary_issue_layer1"`
	SecondaryIssueLayer2 string            `json:"secondary_issue_layer2"`
	SecondaryIssueLayer3 string            `json:"secondary_issue_layer3"`
	SecondaryOutcome     string            `json:"secondary_outcome"`
	TertiaryIssue        string            `json:"tertiary_issue"`
	TertiaryIssueLayer1  string            `json:"tertiary_issue_layer1"`
	TertiaryIssueLayer2  string            `json:"tertiary_issue_layer2"`
	TertiaryIssueLayer3  string            `json:"tertiary_issue_layer3"`
	TertiaryOutcome      string            `json:"tertiary_outcome"`
	EmailAddress         string            `json:"email_address"`
	ReporterIsEmployee   bool              `json:"reporter_is_employee"`
	ReporterNameFirst    string            `json:"reporter_name_first"`
	ReporterNameLast     string            `json:"reporter_name_last"`
	CaseStatus           string            `json:"case_status"`
	Disposition          bool              `json:"disposition"`
	GovernmentNexus      string            `json:"government_nexus"`
	Summary              string            `json:"summary"`
	CaseFiles            []CaseFile        `json:"case_files"`
	CaseQuestions        []CaseQuestion    `json:"case_questions"`
	CaseNotes            []CaseNote        `json:"case_comments"`
	CaseParticipants     []CaseParticipant `json:"case_participants"`
}

func (c *Case) IsICMEntity() bool { return true }
func (c *Case) Scan(src interface{}) error {
	switch v := src.(type) {
	case string:
		return json.Unmarshal([]byte(v), c)
	case []byte:
		return json.Unmarshal(v, c)
	default:
		return errors.New("invalid sql return type for CaseParticipant")
	}
}

type CaseSqlResponse struct {
	ID      string
	Details Case
}

func (c *CaseSqlResponse) GetID() string {
	return c.ID
}
func (c *CaseSqlResponse) GetDetails() interface{} {
	return c.Details
}
func (c *CaseSqlResponse) LoadFromSql(row *sql.Rows) error {
	err := row.Scan(&c.ID, &c.Details)
	return err
}
