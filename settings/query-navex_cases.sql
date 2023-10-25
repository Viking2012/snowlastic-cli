WITH CASE_DETAILS AS (
    SELECT
        "Case ID"                                       AS CASE_ID,
        "Case Number"                                   AS CASE_NUMBER,
        "Alert"                                         AS ALERT_LEVEL,
        "Branch Number"                                 AS BRANCH_NUMBER,
        "Details"                                       AS CASE_DETAILS,
        "Case Type"                                     AS CASE_TYPE,
        NULLIF("Action Taken", '- Select One -')        AS ACTION_TAKEN,
        "Alleged Incident Date"                         AS INCIDENT_DATE,
        "Closed"                                        AS CLOSURE_DATE,
        "Case Due Date"                                 AS DUE_DATE,
        "Opened"                                        AS OPEN_DATE,
        "Received/Reported Date"                        AS REPORTED_DATE,
        "City"                                          AS CASE_CITY,
        "State/Province"                                AS CASE_STATE_PROVINCE,
        "Zip/Postal Code"                               AS CASE_POSTAL_CODE,
        "Country"                                       AS CASE_COUNTRY,
        "Region"                                        AS CASE_REGION,
        "Business"                                      AS BUSINESS_AREA,
        "Business Line"                                 AS DIVISION,
        "Primary Issue"                                 AS PRIMARY_ISSUE,
        "Label for Analytics"                           AS PRIMARY_ISSUE_LAYER1,
        "Primary Issue Type Layer 2 Name"               AS PRIMARY_ISSUE_LAYER2,
        "Primary Issue Type Layer 3 Name"               AS PRIMARY_ISSUE_LAYER3,
        NULLIF("Primary Outcome", '- Select One -')     AS PRIMARY_OUTCOME,
        "Secondary Issue 1"                             AS SECONDARY_ISSUE,
        "Secondary Issue Type Layer 1 Name"             AS SECONDARY_ISSUE_LAYER1,
        "Secondary Issue Type Layer 2 Name"             AS SECONDARY_ISSUE_LAYER2,
        "Secondary Issue Type Layer 3 Name"             AS SECONDARY_ISSUE_LAYER3,
        NULLIF("Secondary Outcome 1", '- Select One -') AS SECONDARY_OUTCOME,
        "Secondary Issue 2"                             AS TERTIARY_ISSUE,
        "Tertiary Issue Type Layer 1 Name"              AS TERTIARY_ISSUE_LAYER1,
        "Tertiary Issue Type Layer 2 Name"              AS TERTIARY_ISSUE_LAYER2,
        "Tertiary Issue Type Layer 3 Name"              AS TERTIARY_ISSUE_LAYER3,
        NULLIF("Secondary Outcome 2", '- Select One -') AS TERTIARY_OUTCOME,
        "Email Address"                                 AS EMAIL_ADDRESS,
        "Reporter Employee"                             AS REPORTER_IS_EMPLOYEE,
        "Reporter First Name"                           AS REPORTER_NAME_FIRST,
        "Reporter Last Name"                            AS REPORTER_NAME_LAST,
        "Status"                                        AS CASE_STATUS,
        "Disposition"                                   AS DISPOSITION,
        CASE
            WHEN NULLIF("Government Nexus", '0') IS NULL 
            THEN FALSE
            ELSE TRUE
        END                                             AS GOVERNMENT_NEXUS,
        "Summary"                                       AS SUMMARY
    FROM
        DEV_FPA_LI.SQL_NAVEX.CASES_EP
),
CASE_FILES AS (
    SELECT
        "Case ID"                                       AS CASE_ID,
        ARRAY_AGG(object_construct(
            'file_id',NULLIF("Case File ID", ''),
            'file_name',NULLIF("Filename", '')
        )) WITHIN GROUP (ORDER BY"Case File ID")        AS DETAILS
    FROM
        DEV_FPA_LI.SQL_NAVEX.CASE_FILE_METADATA
    GROUP BY
        "Case ID"
),
CASE_QUESTIONS AS (
    SELECT
        "Case ID"                                       AS CASE_ID,
        ARRAY_AGG(object_construct(
            'followup_question_id',"ID",
            'is_comment',"Comment",
            'followup_question_text',"Question",
            'followup_question_date',TO_VARCHAR("Question Date", 'YYYY-MM-DDThh:mi:ssZ'),
            'followup_question_asker',"Employee Name",
            'answer_text',
            "Answer",
            'answer_date',
            TO_VARCHAR("Answer Date", 'YYYY-MM-DDThh:mi:ssZ')
        )) WITHIN GROUP (ORDER BY "Question Date","ID") AS DETAILS
    FROM
        DEV_FPA_LI.SQL_NAVEX.FOLLOWUP_QUESTION
    WHERE
        "Answer Date" IS NOT NULL
    GROUP BY
        "Case ID"
),
CASE_COMMENTS AS (
    SELECT
        "Case ID"                                       AS CASE_ID,
        ARRAY_AGG(object_construct(
                'note_id',"ID",
                'note_text',"Note",
                'note_date',TO_VARCHAR("Date", 'YYYY-MM-DDThh:mi:ssZ')
        )) WITHIN GROUP (ORDER BY "Date","ID")          AS DETAILS
    FROM
        DEV_FPA_LI.SQL_NAVEX.FOLLOWUP_COMMENT
    GROUP BY
        "Case ID"
),
CASE_PARTICIPANTS AS (
    SELECT
        "Case ID"                                       AS CASE_ID,
        ARRAY_AGG(
            object_construct(
                'participant_id',"Participant ID",
                'participant_prefix',"Prefix",
                'participant_name',TRIM(REPLACE(CONCAT("First Name",' ',"Middle Name",' ',"Last Name"),'  ',' ')),
                'participant_hr_id',"HR ID",
                'participant_relationship',"Relationship Name",
                'participant_role',"Role Name",
                'outcome',"Outcome Name"
        )) WITHIN GROUP (ORDER BY "Participant ID")     AS DETAILS
    FROM
        DEV_FPA_LI.SQL_NAVEX.PARTICIPANTS_EP
    GROUP BY
        "Case ID"
)
SELECT
    CD.CASE_ID                  AS case_id,
    CD.CASE_NUMBER              AS case_number,
    CD.ALERT_LEVEL              AS alert_level,
    CD.BRANCH_NUMBER            AS branch_number,
    CD.CASE_DETAILS             AS case_details,
    CD.CASE_TYPE                AS case_type,
    CD.ACTION_TAKEN             AS action_taken,
    CD.INCIDENT_DATE            AS incident_date,
    CD.CLOSURE_DATE             AS closure_date,
    CD.DUE_DATE                 AS due_date,
    CD.OPEN_DATE                AS open_date,
    CD.REPORTED_DATE            AS reported_date,
    CD.CASE_CITY                AS case_city,
    CD.CASE_STATE_PROVINCE      AS case_state_province,
    CD.CASE_POSTAL_CODE         AS case_postal_code,
    CD.CASE_COUNTRY             AS case_country,
    CD.CASE_REGION              AS case_region,
    CD.BUSINESS_AREA            AS business_area,
    CD.DIVISION                 AS division,
    CD.PRIMARY_ISSUE            AS primary_issue,
    CD.PRIMARY_ISSUE_LAYER1     AS primary_issue_layer1,
    CD.PRIMARY_ISSUE_LAYER2     AS primary_issue_layer2,
    CD.PRIMARY_ISSUE_LAYER3     AS primary_issue_layer3,
    CD.PRIMARY_OUTCOME          AS primary_outcome,
    CD.SECONDARY_ISSUE          AS secondary_issue,
    CD.SECONDARY_ISSUE_LAYER1   AS secondary_issue_layer1,
    CD.SECONDARY_ISSUE_LAYER2   AS secondary_issue_layer2,
    CD.SECONDARY_ISSUE_LAYER3   AS secondary_issue_layer3,
    CD.SECONDARY_OUTCOME        AS secondary_outcome,
    CD.TERTIARY_ISSUE           AS tertiary_issue,
    CD.TERTIARY_ISSUE_LAYER1    AS tertiary_issue_layer1,
    CD.TERTIARY_ISSUE_LAYER2    AS tertiary_issue_layer2,
    CD.TERTIARY_ISSUE_LAYER3    AS tertiary_issue_layer3,
    CD.TERTIARY_OUTCOME         AS tertiary_outcome,
    CD.EMAIL_ADDRESS            AS email_address,
    CD.REPORTER_IS_EMPLOYEE     AS reporter_is_employee,
    CD.REPORTER_NAME_FIRST      AS reporter_name_first,
    CD.REPORTER_NAME_LAST       AS reporter_name_last,
    CD.CASE_STATUS              AS case_status,
    CD.DISPOSITION              AS disposition,
    CD.GOVERNMENT_NEXUS         AS government_nexus,
    CD.SUMMARY                  AS summary,
    CF.DETAILS                  AS case_files,
    CQ.DETAILS                  AS case_questions,
    CC.DETAILS                  AS case_comments,
    CP.DETAILS                  AS case_participants
FROM
    CASE_DETAILS CD
    LEFT JOIN CASE_FILES CF         ON CD.CASE_ID = CF.CASE_ID
    LEFT JOIN CASE_QUESTIONS CQ     ON CD.CASE_ID = CQ.CASE_ID
    LEFT JOIN CASE_COMMENTS CC      ON CD.CASE_ID = CC.CASE_ID
    LEFT JOIN CASE_PARTICIPANTS CP  ON CD.CASE_ID = CP.CASE_ID