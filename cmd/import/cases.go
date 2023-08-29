/*
Copyright © 2023 Alexander Orban <alexander.orban@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package _import

import (
	"database/sql"
	"errors"
	"fmt"
	icm_orm "github.com/alexander-orban/icm_goapi/orm"
	"github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"math"
	es "snowlastic-cli/pkg/es"
	orm "snowlastic-cli/pkg/orm"
	"snowlastic-cli/pkg/snowflake"
	"time"
)

// import/casesCmd represents the import/cases command
var casesCmd = &cobra.Command{
	Use:   "cases",
	Short: "Index all Navex Cases contained in the snowflake database",
	Long:  ``,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !viper.IsSet("esClient") {
			return errors.New("elasticsearch was somehow not created by the `import` command prior to running `demos`")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			indexName = "cases"

			db   *sql.DB
			rows *sql.Rows
			c    *elasticsearch.Client

			docs       = make(chan icm_orm.ICMEntity, es.BulkInsertSize)
			numErrors  int64
			numIndexed int64
			rowCount   int64
			numBatches float64

			err error
		)

		log.Println("connecting to database")
		db, err = snowflake.NewDB(snowflake.Config{
			Account:   viper.GetString("snowflakeAccount"),
			Warehouse: viper.GetString("snowflakeWarehouse"),
			Database:  viper.GetString("snowflakeDatabase"),
			Schema:    "SQL_NAVEX",
			User:      viper.GetString("snowflakeUser"),
			Password:  viper.GetString("snowflakePassword"),
			Role:      viper.GetString("snowflakeRole"),
		})
		if err != nil {
			return err
		}
		defer db.Close()

		var countQuery = "SELECT COUNT(1) FROM (" + string(caseQuery) + ")"
		rows, err = db.Query(countQuery)
		if err != nil {
			return err
		}
		rows.Next()
		err = rows.Scan(&rowCount)
		if err != nil {
			return err
		}
		numBatches = math.Ceil(float64(rowCount) / es.BulkInsertSize)

		start := time.Now().UTC()
		go func() {
			log.Println("reading cases from database")
			rows, err = db.Query(caseQuery)
			if err != nil {
				log.Fatal(err)
			}
			for rows.Next() {
				var c orm.Case
				if err := c.ScanFrom(rows); err != nil {
					log.Fatal(err)
				}
				//fmt.Printf("id: %-15s notes: %4d questions: %4d files: %4d participants: %4d\n", c.GetID(), len(c.CaseNotes), len(c.CaseQuestions), len(c.CaseFiles), len(c.CaseParticipants))
				docs <- icm_orm.ICMEntity(&c)
			}
			close(docs)
		}()

		// Get the generated elasticsearch client
		c, err = getElasticClient()
		if err != nil {
			return err
		}
		batches := es.BatchEntities(docs, es.BulkInsertSize)
		log.Println("indexing cases")
		numIndexed, numErrors, err = es.BulkImport(c, batches, indexName, int64(numBatches))
		if err != nil {
			return err
		}

		dur := time.Since(start)
		if numErrors > 0 {
			return errors.New(fmt.Sprintf(
				"Indexed [%s] documents with [%s] errors in %s (%s docs/sec)",
				humanize.Comma(int64(numIndexed)),
				humanize.Comma(int64(numErrors)),
				dur.Truncate(time.Millisecond),
				humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(numIndexed))),
			))
		} else {
			log.Printf(
				"Sucessfuly indexed [%s] documents in %s (%s docs/sec)",
				humanize.Comma(int64(numIndexed)),
				dur.Truncate(time.Millisecond),
				humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(numIndexed))),
			)
		}
		return nil
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// import/casesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// import/casesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

const caseQuery string = `WITH CASE_DETAILS AS (
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
LEFT JOIN CASE_PARTICIPANTS CP ON CD.CASE_ID=CP.CASE_ID
`
