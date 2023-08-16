package types

import (
	"encoding/json"
	"errors"
	"github.com/alexander-orban/icm_goapi/orm"
)

type CaseFile struct {
	FileID   string `json:"file_id"`
	FileName string `json:"file_name"`
}
type CaseFiles []CaseFile

func (c *CaseFiles) IsICMEntity() bool { return true }
func (c *CaseFiles) Scan(src interface{}) error {
	var cf []CaseFile
	if src == nil {
		*c = cf
		return nil
	}
	switch v := src.(type) {
	case string:
		if err := json.Unmarshal([]byte(v), &cf); err != nil {
			return err
		}
	case []byte:
		if err := json.Unmarshal(v, &cf); err != nil {
			return err
		}
	default:
		return errors.New("invalid sql return type for CaseFiles")
	}
	*c = cf
	return nil
}

type CaseQuestion struct {
	QuestionID    string       `json:"followup_question_id,omitempty"`
	IsComment     string       `json:"is_comment,omitempty"`
	QuestionText  string       `json:"followup_question_text,omitempty"`
	QuestionDate  orm.NullTime `json:"followup_question_date,omitempty"`
	QuestionAsker string       `json:"followup_question_asker,omitempty"`
	AnswerText    string       `json:"answer_text,omitempty"`
	AnswerDate    orm.NullTime `json:"answer_date,omitempty"`
}
type CaseQuestions []CaseQuestion

func (c *CaseQuestions) IsICMEntity() bool { return true }
func (c *CaseQuestions) Scan(src interface{}) error {
	var cq []CaseQuestion
	if src == nil {
		*c = cq
		return nil
	}
	switch v := src.(type) {
	case string:
		if err := json.Unmarshal([]byte(v), &cq); err != nil {
			return err
		}
	case []byte:
		if err := json.Unmarshal(v, &cq); err != nil {
			return err
		}
	default:
		return errors.New("invalid sql return type for CaseQuestions")
	}
	*c = cq
	return nil
}

type CaseNote struct {
	NoteID   string       `json:"note_id,omitempty"`
	NoteText string       `json:"note_text,omitempty"`
	NoteDate orm.NullTime `json:"note_date,omitempty"`
}
type CaseNotes []CaseNote

func (c *CaseNotes) IsICMEntity() bool { return true }
func (c *CaseNotes) Scan(src interface{}) error {
	var cn []CaseNote
	if src == nil {
		*c = cn
		return nil
	}
	switch v := src.(type) {
	case string:
		if err := json.Unmarshal([]byte(v), &cn); err != nil {
			return err
		}
	case []byte:
		if err := json.Unmarshal(v, &cn); err != nil {
			return err
		}
	default:
		return errors.New("invalid sql return type for CaseNotes")
	}
	*c = cn
	return nil
}

type CaseParticipant struct {
	ParticipantID           string `json:"participant_id,omitempty"`
	ParticipantPrefix       string `json:"participant_prefix,omitempty"`
	ParticipantName         string `json:"participant_name,omitempty"`
	ParticipantHRID         string `json:"participant_hr_id,omitempty"`
	ParticipantRelationship string `json:"participant_relationship,omitempty"`
	ParticipantRole         string `json:"participant_role,omitempty"`
	Outcome                 string `json:"outcome,omitempty"`
}
type CaseParticipants []CaseParticipant

func (c *CaseParticipants) IsICMEntity() bool { return true }
func (c *CaseParticipants) Scan(src interface{}) error {
	var cp []CaseParticipant
	if src == nil {
		*c = cp
		return nil
	}
	switch v := src.(type) {
	case string:
		if err := json.Unmarshal([]byte(v), &cp); err != nil {
			return err
		}
	case []byte:
		if err := json.Unmarshal(v, &cp); err != nil {
			return err
		}
	default:
		return errors.New("invalid sql return type for CaseParticipants")
	}
	*c = cp
	return nil
}

type Case struct {
	CaseID               string           `json:"case_id"`
	CaseNumber           string           `json:"case_number"`
	Alert                string           `json:"alert_level,omitempty"`
	CaseBranchNumber     string           `json:"branch_number,omitempty"`
	CaseDetails          string           `json:"case_details,omitempty"`
	CaseType             string           `json:"case_type,omitempty"`
	ActionTaken          string           `json:"action_taken,omitempty"`
	IncidentDate         orm.NullTime     `json:"incident_date,omitempty"`
	ClosureDate          orm.NullTime     `json:"closure_date,omitempty"`
	DueDate              orm.NullTime     `json:"due_date,omitempty"`
	OpenDate             orm.NullTime     `json:"open_date,omitempty"`
	ReportDate           orm.NullTime     `json:"reported_date,omitempty"`
	CaseCity             string           `json:"case_city,omitempty"`
	CaseStateProvince    string           `json:"case_state_province,omitempty"`
	CasePostalCode       string           `json:"case_postal_code,omitempty"`
	CaseCountry          string           `json:"case_country,omitempty"`
	CaseRegion           string           `json:"case_region,omitempty"`
	BusinessArea         string           `json:"business_case,omitempty"`
	Division             string           `json:"division,omitempty"`
	PrimaryIssue         string           `json:"primary_issue,omitempty"`
	PrimaryIssueLayer1   string           `json:"primary_issue_layer1,omitempty"`
	PrimaryIssueLayer2   string           `json:"primary_issue_layer2,omitempty"`
	PrimaryIssueLayer3   string           `json:"primary_issue_layer3,omitempty"`
	PrimaryOutcome       string           `json:"primary_outcome,omitempty"`
	SecondaryIssue       string           `json:"secondary_issue,omitempty"`
	SecondaryIssueLayer1 string           `json:"secondary_issue_layer1,omitempty"`
	SecondaryIssueLayer2 string           `json:"secondary_issue_layer2,omitempty"`
	SecondaryIssueLayer3 string           `json:"secondary_issue_layer3,omitempty"`
	SecondaryOutcome     string           `json:"secondary_outcome,omitempty"`
	TertiaryIssue        string           `json:"tertiary_issue,omitempty"`
	TertiaryIssueLayer1  string           `json:"tertiary_issue_layer1,omitempty"`
	TertiaryIssueLayer2  string           `json:"tertiary_issue_layer2,omitempty"`
	TertiaryIssueLayer3  string           `json:"tertiary_issue_layer3,omitempty"`
	TertiaryOutcome      string           `json:"tertiary_outcome,omitempty"`
	EmailAddress         string           `json:"email_address,omitempty"`
	ReporterIsEmployee   bool             `json:"reporter_is_employee,omitempty"`
	ReporterNameFirst    string           `json:"reporter_name_first,omitempty"`
	ReporterNameLast     string           `json:"reporter_name_last,omitempty"`
	CaseStatus           string           `json:"case_status,omitempty"`
	Disposition          string           `json:"disposition,omitempty"`
	GovernmentNexus      string           `json:"government_nexus,omitempty"`
	Summary              string           `json:"summary,omitempty"`
	CaseFiles            CaseFiles        `json:"case_files,omitempty"`
	CaseQuestions        CaseQuestions    `json:"case_questions,omitempty"`
	CaseNotes            CaseNotes        `json:"case_comments,omitempty"`
	CaseParticipants     CaseParticipants `json:"case_participants,omitempty"`
}

func (c *Case) IsICMEntity() bool { return true }
func (c *Case) GetID() string {
	return c.CaseNumber
}
func (c *Case) Scan(src interface{}) error {
	switch v := src.(type) {
	case string:
		return json.Unmarshal([]byte(v), c)
	case []byte:
		return json.Unmarshal(v, c)
	default:
		return errors.New("invalid sql return type for Case")
	}
}
