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
type CaseFiles []CaseFile

func (c *CaseFiles) IsDocument()                    {}
func (c *CaseFiles) GetID() string                  { return "" }
func (c *CaseFiles) GetQuery(string, string) string { return "" }
func (c *CaseFiles) ScanFrom(*sql.Rows) error       { return nil }
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

func (c *CaseQuestions) IsDocument()                    {}
func (c *CaseQuestions) GetID() string                  { return "" }
func (c *CaseQuestions) GetQuery(string, string) string { return "" }
func (c *CaseQuestions) ScanFrom(*sql.Rows) error       { return nil }
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

func (c *CaseNotes) IsDocument()                    {}
func (c *CaseNotes) GetID() string                  { return "" }
func (c *CaseNotes) GetQuery(string, string) string { return "" }
func (c *CaseNotes) ScanFrom(*sql.Rows) error       { return nil }
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

func (c *CaseParticipants) IsDocument()                    {}
func (c *CaseParticipants) GetID() string                  { return "" }
func (c *CaseParticipants) GetQuery(string, string) string { return "" }
func (c *CaseParticipants) ScanFrom(*sql.Rows) error       { return nil }
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

func (c *Case) IsDocument() {}
func (c *Case) GetID() string {
	return c.CaseNumber
}
func (c *Case) GetQuery(string, string) string {
	return `WITH CASE_DETAILS AS (
    SELECT 
        "Case ID" AS CASE_ID,
        "Case Number"       AS CASE_NUMBER,
        "Alert" AS ALERT_LEVEL,
        "Branch Number" AS BRANCH_NUMBER,
        "Details" AS CASE_DETAILS,
        "Case Type" AS CASE_TYPE,
        IFNULL(NULLIF("Action Taken",'- Select One -'),'') AS ACTION_TAKEN,
        "Alleged Incident Date" AS INCIDENT_DATE,
        "Closed" AS CLOSURE_DATE,
        "Case Due Date" AS DUE_DATE,
        "Opened" AS OPEN_DATE,
        "Received/Reported Date" AS REPORTED_DATE,
        "City" AS CASE_CITY,
        "State/Province" AS CASE_STATE_PROVINCE,
        "Zip/Postal Code" AS CASE_POSTAL_CODE,
        "Country" AS CASE_COUNTRY,
        "Region" AS CASE_REGION,
        "Business" AS BUSINESS_AREA,
        "Business Line" AS DIVISION,
        "Primary Issue" AS PRIMARY_ISSUE,
        "Label for Analytics" AS PRIMARY_ISSUE_LAYER1,
        "Primary Issue Type Layer 2 Name" AS PRIMARY_ISSUE_LAYER2,
        "Primary Issue Type Layer 3 Name" AS PRIMARY_ISSUE_LAYER3,
        IFNULL(NULLIF("Primary Outcome",'- Select One -'),'') AS PRIMARY_OUTCOME,
        "Secondary Issue 1" AS SECONDARY_ISSUE,
        "Secondary Issue Type Layer 1 Name" AS SECONDARY_ISSUE_LAYER1,
        "Secondary Issue Type Layer 2 Name" AS SECONDARY_ISSUE_LAYER2,
        "Secondary Issue Type Layer 3 Name" AS SECONDARY_ISSUE_LAYER3,
        IFNULL(NULLIF("Secondary Outcome 1",'- Select One -'),'') AS SECONDARY_OUTCOME,
        "Secondary Issue 2" AS TERTIARY_ISSUE,
        "Tertiary Issue Type Layer 1 Name" AS TERTIARY_ISSUE_LAYER1,
        "Tertiary Issue Type Layer 2 Name" AS TERTIARY_ISSUE_LAYER2,
        "Tertiary Issue Type Layer 3 Name" AS TERTIARY_ISSUE_LAYER3,
        IFNULL(NULLIF("Secondary Outcome 2",'- Select One -'),'') AS TERTIARY_OUTCOME,
        "Email Address" AS EMAIL_ADDRESS,
        "Reporter Employee" AS REPORTER_IS_EMPLOYEE,
        "Reporter First Name" AS REPORTER_NAME_FIRST,
        "Reporter Last Name"  AS REPORTER_NAME_LAST,
        "Status" AS CASE_STATUS,
        IFNULL("Disposition",'') AS DISPOSITION,
        CASE WHEN NULLIF("Government Nexus",'0') IS NULL THEN FALSE ELSE TRUE END AS GOVERNMENT_NEXUS,
        "Summary" AS SUMMARY
    FROM DEV_FPA_LI.SQL_NAVEX.CASES_EP
),
CASE_FILES AS (
    SELECT 
        "Case ID" AS CASE_ID, 
        ARRAY_AGG(object_construct(
            'file_id',NULLIF("Case File ID",''),
            'file_name',NULLIF("Filename",'')
        )) WITHIN GROUP (ORDER BY "Case File ID") AS DETAILS
    FROM DEV_FPA_LI.SQL_NAVEX.CASE_FILE_METADATA
    GROUP BY "Case ID"
),
CASE_QUESTIONS AS (
    SELECT 
        "Case ID" AS CASE_ID, 
        ARRAY_AGG(object_construct(
            'followup_question_id',"ID",
            'is_comment',"Comment",
            'followup_question_text',"Question",
            'followup_question_date',TO_VARCHAR("Question Date", 'YYYY-MM-DDThh:mi:ssZ'),
            'followup_question_asker',"Employee Name",
            'answer_text',"Answer",
            'answer_date',TO_VARCHAR("Answer Date", 'YYYY-MM-DDThh:mi:ssZ')
        )) WITHIN GROUP (ORDER BY "Question Date","ID") AS DETAILS
    FROM DEV_FPA_LI.SQL_NAVEX.FOLLOWUP_QUESTION
    WHERE "Answer Date" IS NOT NULL
    GROUP BY "Case ID"
),
CASE_COMMENTS AS (
    SELECT 
        "Case ID" AS CASE_ID,
        ARRAY_AGG(object_construct(
            'note_id',"ID",
            'note_text',"Note",
            'note_date',TO_VARCHAR("Date", 'YYYY-MM-DDThh:mi:ssZ')
        )) WITHIN GROUP (ORDER BY "Date","ID") AS DETAILS
    FROM DEV_FPA_LI.SQL_NAVEX.FOLLOWUP_COMMENT
    GROUP BY "Case ID"
),
CASE_PARTICIPANTS AS (
    SELECT 
        "Case ID" AS CASE_ID,
        ARRAY_AGG(object_construct(
            'participant_id', "Participant ID",
            'participant_prefix',"Prefix",
            'participant_name',TRIM(REPLACE(CONCAT("First Name",' ',"Middle Name",' ',"Last Name"),'  ',' ')),
            'participant_hr_id',"HR ID",
            'participant_relationship',"Relationship Name",
            'participant_role',"Role Name",
            'outcome',"Outcome Name"
        )) WITHIN GROUP (ORDER BY "Participant ID") AS DETAILS
    FROM DEV_FPA_LI.SQL_NAVEX.PARTICIPANTS_EP
    GROUP BY "Case ID"
)
SELECT 
    CD.CASE_ID as case_id,
    CD.CASE_NUMBER AS case_number,
    CD.ALERT_LEVEL AS alert_level,
    CD.BRANCH_NUMBER AS branch_number,
    CD.CASE_DETAILS AS case_details,
    CD.CASE_TYPE AS case_type,
    CD.ACTION_TAKEN AS action_taken,
    CD.INCIDENT_DATE AS incident_date,
    CD.CLOSURE_DATE AS closure_date,
    CD.DUE_DATE AS due_date,
    CD.OPEN_DATE AS open_date,
    CD.REPORTED_DATE AS reported_date,
    CD.CASE_CITY AS case_city,
    CD.CASE_STATE_PROVINCE AS case_state_province,
    CD.CASE_POSTAL_CODE AS case_postal_code,
    CD.CASE_COUNTRY AS case_country,
    CD.CASE_REGION AS case_region,
    CD.BUSINESS_AREA AS business_area,
    CD.DIVISION AS division,
    CD.PRIMARY_ISSUE AS primary_issue,
    CD.PRIMARY_ISSUE_LAYER1 AS primary_issue_layer1,
    CD.PRIMARY_ISSUE_LAYER2 AS primary_issue_layer2,
    CD.PRIMARY_ISSUE_LAYER3 AS primary_issue_layer3,
    CD.PRIMARY_OUTCOME AS primary_outcome,
    CD.SECONDARY_ISSUE AS secondary_issue,
    CD.SECONDARY_ISSUE_LAYER1 AS secondary_issue_layer1,
    CD.SECONDARY_ISSUE_LAYER2 AS secondary_issue_layer2,
    CD.SECONDARY_ISSUE_LAYER3 AS secondary_issue_layer3,
    CD.SECONDARY_OUTCOME AS secondary_outcome,
    CD.TERTIARY_ISSUE AS tertiary_issue,
    CD.TERTIARY_ISSUE_LAYER1 AS tertiary_issue_layer1,
    CD.TERTIARY_ISSUE_LAYER2 AS tertiary_issue_layer2,
    CD.TERTIARY_ISSUE_LAYER3 AS tertiary_issue_layer3,
    CD.TERTIARY_OUTCOME AS tertiary_outcome,
    CD.EMAIL_ADDRESS AS email_address,
    CD.REPORTER_IS_EMPLOYEE AS reporter_is_employee,
    CD.REPORTER_NAME_FIRST AS reporter_name_first,
    CD.REPORTER_NAME_LAST AS reporter_name_last,
    CD.CASE_STATUS AS case_status,
    CD.DISPOSITION AS disposition,
    CD.GOVERNMENT_NEXUS AS government_nexus,
    CD.SUMMARY AS summary,
    CF.DETAILS AS case_files,
    CQ.DETAILS AS case_questions,
    CC.DETAILS AS case_comments,
    CP.DETAILS AS case_participants
FROM CASE_DETAILS CD
LEFT JOIN CASE_FILES CF        ON CD.CASE_ID=CF.CASE_ID
LEFT JOIN CASE_QUESTIONS CQ    ON CD.CASE_ID=CQ.CASE_ID
LEFT JOIN CASE_COMMENTS CC     ON CD.CASE_ID=CC.CASE_ID
LEFT JOIN CASE_PARTICIPANTS CP ON CD.CASE_ID=CP.CASE_ID`
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
func (c *Case) ScanFrom(rows *sql.Rows) error {
	return rows.Scan(
		&c.CaseID,
		&c.CaseNumber,
		&c.Alert,
		&c.CaseBranchNumber,
		&c.CaseDetails,
		&c.CaseType,
		&c.ActionTaken,
		&c.IncidentDate,
		&c.ClosureDate,
		&c.DueDate,
		&c.OpenDate,
		&c.ReportDate,
		&c.CaseCity,
		&c.CaseStateProvince,
		&c.CasePostalCode,
		&c.CaseCountry,
		&c.CaseRegion,
		&c.BusinessArea,
		&c.Division,
		&c.PrimaryIssue,
		&c.PrimaryIssueLayer1,
		&c.PrimaryIssueLayer2,
		&c.PrimaryIssueLayer3,
		&c.PrimaryOutcome,
		&c.SecondaryIssue,
		&c.SecondaryIssueLayer1,
		&c.SecondaryIssueLayer2,
		&c.SecondaryIssueLayer3,
		&c.SecondaryOutcome,
		&c.TertiaryIssue,
		&c.TertiaryIssueLayer1,
		&c.TertiaryIssueLayer2,
		&c.TertiaryIssueLayer3,
		&c.TertiaryOutcome,
		&c.EmailAddress,
		&c.ReporterIsEmployee,
		&c.ReporterNameFirst,
		&c.ReporterNameLast,
		&c.CaseStatus,
		&c.Disposition,
		&c.GovernmentNexus,
		&c.Summary,
		&c.CaseFiles,
		&c.CaseQuestions,
		&c.CaseNotes,
		&c.CaseParticipants,
	)
}
func (c *Case) New() SnowlasticDocument { return new(Case) }
