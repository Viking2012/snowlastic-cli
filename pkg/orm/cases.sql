WITH CASE_DETAILS AS (
    SELECT 
        "Case ID" AS CASE_ID,
        NULLIF("Case Number",'')       AS CASE_NUMBER,
        NULLIF("Alert",'') AS ALERT_LEVEL,
        NULLIF("Branch Number",'') AS BRANCH_NUMBER,
        NULLIF("Details",'') AS CASE_DETAILS,
        NULLIF("Case Type",'') AS CASE_TYPE,
        NULLIF(NULLIF("Action Taken",'- Select One -'),'') AS ACTION_TAKEN,
        "Alleged Incident Date" AS INCIDENT_DATE,
        "Closed" AS CLOSURE_DATE,
        "Case Due Date" AS DUE_DATE,
        "Opened" AS OPEN_DATE,
        "Received/Reported Date" AS REPORTED_DATE,
        NULLIF("City",'') AS CASE_CITY,
        NULLIF("State/Province",'') AS CASE_STATE_PROVINCE,
        NULLIF("Zip/Postal Code",'') AS CASE_POSTAL_CODE,
        NULLIF("Country",'') AS CASE_COUNTRY,
        NULLIF("Region",'') AS CASE_REGION,
        NULLIF("Business",'') AS BUSINESS_AREA,
        NULLIF("Business Line",'') AS DIVISION,
        NULLIF("Primary Issue",'') AS PRIMARY_ISSUE,
        NULLIF("Label for Analytics",'') AS PRIMARY_ISSUE_LAYER1,
        NULLIF("Primary Issue Type Layer 2 Name",'') AS PRIMARY_ISSUE_LAYER2,
        NULLIF("Primary Issue Type Layer 3 Name",'') AS PRIMARY_ISSUE_LAYER3,
        NULLIF(NULLIF("Primary Outcome",'- Select One -'),'') AS PRIMARY_OUTCOME,
        NULLIF("Secondary Issue 1",'') AS SECONDARY_ISSUE,
        NULLIF("Secondary Issue Type Layer 1 Name",'') AS SECONDARY_ISSUE_LAYER1,
        NULLIF("Secondary Issue Type Layer 2 Name",'') AS SECONDARY_ISSUE_LAYER2,
        NULLIF("Secondary Issue Type Layer 3 Name",'') AS SECONDARY_ISSUE_LAYER3,
        NULLIF(NULLIF("Secondary Outcome 1",'- Select One -'),'') AS SECONDARY_OUTCOME,
        NULLIF("Secondary Issue 2",'') AS TERTIARY_ISSUE,
        NULLIF("Tertiary Issue Type Layer 1 Name",'') AS TERTIARY_ISSUE_LAYER1,
        NULLIF("Tertiary Issue Type Layer 2 Name",'') AS TERTIARY_ISSUE_LAYER2,
        NULLIF("Tertiary Issue Type Layer 3 Name",'') AS TERTIARY_ISSUE_LAYER3,
        NULLIF(NULLIF("Secondary Outcome 2",'- Select One -'),'') AS TERTIARY_OUTCOME,
        NULLIF("Email Address",'') AS EMAIL_ADDRESS,
        "Reporter Employee" AS REPORTER_IS_EMPLOYEE,
        NULLIF("Reporter First Name",'') AS REPORTER_NAME_FIRST,
        NULLIF("Reporter Last Name",'')  AS REPORTER_NAME_LAST,
        NULLIF("Status",'') AS CASE_STATUS,
        NULLIF("Disposition",'') AS DISPOSITION,
        CASE WHEN NULLIF("Government Nexus",'0') IS NULL THEN FALSE ELSE TRUE END AS GOVERNMENT_NEXUS,
        NULLIF("Summary",'') AS SUMMARY
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
            'followup_question_id',NULLIF("ID",''),
            'is_comment',NULLIF("Comment",''),
            'followup_question_text',NULLIF("Question",''),
            'followup_question_date',"Question Date",
            'followup_question_asker',NULLIF("Employee Name",''),
            'answer_text',NULLIF("Answer",''),
            'answer_date',"Answer Date"
        )) WITHIN GROUP (ORDER BY "Question Date","ID") AS DETAILS
    FROM DEV_FPA_LI.SQL_NAVEX.FOLLOWUP_QUESTION
    WHERE "Answer Date" IS NOT NULL
    GROUP BY "Case ID"
),
CASE_COMMENTS AS (
    SELECT 
        "Case ID" AS CASE_ID,
        ARRAY_AGG(object_construct(
            'note_id',NULLIF("ID",''),
            'note_text',NULLIF("Note",''),
            'note_date',"Date"
        )) WITHIN GROUP (ORDER BY "Date","ID") AS DETAILS
    FROM DEV_FPA_LI.SQL_NAVEX.FOLLOWUP_COMMENT
    GROUP BY "Case ID"
),
CASE_PARTICIPANTS AS (
    SELECT 
        "Case ID" AS CASE_ID,
        ARRAY_AGG(object_construct(
            'participant_id', NULLIF("Participant ID",''),
            'participant_prefix',NULLIF("Prefix",''),
            'participant_name',NULLIF(TRIM(REPLACE(CONCAT("First Name",' ',"Middle Name",' ',"Last Name"),'  ',' ')),''),
            'participant_hr_id',NULLIF("HR ID",''),
            'participant_relationship',NULLIF("Relationship Name",''),
            'participant_role',NULLIF("Role Name",''),
            'outcome',NULLIF("Outcome Name",'')
        )) WITHIN GROUP (ORDER BY "Participant ID") AS DETAILS
    FROM DEV_FPA_LI.SQL_NAVEX.PARTICIPANTS_EP
    GROUP BY "Case ID"
)
SELECT 
    CD.CASE_ID, 
    to_json(object_construct(
        'case_number',CASE_NUMBER,
        'alert_level',ALERT_LEVEL,
        'branch_number',BRANCH_NUMBER,
        'case_details',CASE_DETAILS,
        'case_type',CASE_TYPE,
        'action_taken',ACTION_TAKEN,
        'incident_date',INCIDENT_DATE,
        'closure_date',CLOSURE_DATE,
        'due_date',DUE_DATE,
        'open_date',OPEN_DATE,
        'reported_date',REPORTED_DATE,
        'case_city',CASE_CITY,
        'case_state_province',CASE_STATE_PROVINCE,
        'case_postal_code',CASE_POSTAL_CODE,
        'case_country',CASE_COUNTRY,
        'case_region',CASE_REGION,
        'business_area',BUSINESS_AREA,
        'division',DIVISION,
        'primary_issue',PRIMARY_ISSUE,
        'primary_issue_layer1',PRIMARY_ISSUE_LAYER1,
        'primary_issue_layer2',PRIMARY_ISSUE_LAYER2,
        'primary_issue_layer3',PRIMARY_ISSUE_LAYER3,
        'primary_outcome',PRIMARY_OUTCOME,
        'secondary_issue',SECONDARY_ISSUE,
        'secondary_issue_layer1',SECONDARY_ISSUE_LAYER1,
        'secondary_issue_layer2',SECONDARY_ISSUE_LAYER2,
        'secondary_issue_layer3',SECONDARY_ISSUE_LAYER3,
        'secondary_outcome',SECONDARY_OUTCOME,
        'tertiary_issue',TERTIARY_ISSUE,
        'tertiary_issue_layer1',TERTIARY_ISSUE_LAYER1,
        'tertiary_issue_layer2',TERTIARY_ISSUE_LAYER2,
        'tertiary_issue_layer3',TERTIARY_ISSUE_LAYER3,
        'tertiary_outcome',TERTIARY_OUTCOME,
        'email_address',EMAIL_ADDRESS,
        'reporter_is_employee',REPORTER_IS_EMPLOYEE,
        'reporter_name_first',REPORTER_NAME_FIRST,
        'reporter_name_last',REPORTER_NAME_LAST,
        'case_status',CASE_STATUS,
        'disposition',DISPOSITION,
        'government_nexus',GOVERNMENT_NEXUS,
        'summary',SUMMARY,
        'case_files',CF.DETAILS,
        'case_questions',CQ.DETAILS,
        'case_comments',CC.DETAILS,
        'case_participants',CP.DETAILS
    )) as DETAILS
FROM CASE_DETAILS CD
LEFT JOIN CASE_FILES CF        ON CD.CASE_ID=CF.CASE_ID
LEFT JOIN CASE_QUESTIONS CQ    ON CD.CASE_ID=CQ.CASE_ID
LEFT JOIN CASE_COMMENTS CC     ON CD.CASE_ID=CC.CASE_ID
LEFT JOIN CASE_PARTICIPANTS CP ON CD.CASE_ID=CP.CASE_ID
ORDER BY TO_NUMBER(CASE_ID)