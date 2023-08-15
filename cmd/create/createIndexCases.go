package create

import (
	"bytes"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
)

var ()

func indexCase(c *elasticsearch.Client) error {
	res, err := c.Indices.Delete([]string{"case"})
	if err != nil {
		return fmt.Errorf("cannot delete index: %s", err)
	}
	if res.IsError() {
		log.Println("error when deleting index", res.String())
	} else {
		log.Println(res.String())
	}

	b := []byte(caseIndex)
	res, err = c.Indices.Create("case", c.Indices.Create.WithBody(bytes.NewReader(b)))
	if err != nil {
		return fmt.Errorf("cannot create index: %s", err)
	}

	if res.IsError() {
		return fmt.Errorf("cannot create index, got an error response code: %s\n", res.String())
	}
	log.Println("successfully created an index")
	return nil
}

const caseIndex string = `{
  "settings": {
    "analysis": {
      "analyzer": {
        "caseAnalyzer": {
          "char_filter": ["html_strip"],
          "tokenizer": "standard",
          "filter": ["lowercase","stop","asciifolding","stemmer","edge_ngram"]
        }
      },
      "normalizer": {
        "caseNormalizer": {
          "type": "custom",
          "filter": ["lowercase","asciifolding"]
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "action_taken": 			{"type": "keyword",	"normalizer": "caseNormalizer"},
      "alert_level": 			{"type": "keyword",	"normalizer": "caseNormalizer"},
      "branch_number": 			{"type": "keyword",	"normalizer": "caseNormalizer"},
      "business_case": 			{"type": "keyword",	"normalizer": "caseNormalizer"},
      "case_city": 				{"type": "keyword",	"normalizer": "caseNormalizer"},
      "case_comments": 			{"type": "text"                                   },
      "case_country": 			{"type": "keyword",	"normalizer": "caseNormalizer"},
      "case_details": 			{"type": "text"                                   },
      "case_files": 			{"type": "text"                                   },
      "case_number": 			{"type": "keyword",	"normalizer": "caseNormalizer"},
      "case_participants": 		{"type": "text"                                   },
      "case_postal_code": 		{"type": "keyword",	"normalizer": "caseNormalizer"},
      "case_questions": 		{"type": "text"                                   },
      "case_region": 			{"type": "keyword",	"normalizer": "caseNormalizer"},
      "case_state_province": 	{"type": "keyword",	"normalizer": "caseNormalizer"},
      "case_status": 			{"type": "keyword",	"normalizer": "caseNormalizer"},
      "case_type": 				{"type": "text"                                   },
      "closure_date": 			{"type": "keyword",	"normalizer": "caseNormalizer"},
      "disposition": 			{"type": "keyword",	"normalizer": "caseNormalizer"},
      "division": 				{"type": "keyword",	"normalizer": "caseNormalizer"},
      "due_date": 				{"type": "keyword",	"normalizer": "caseNormalizer"},
      "email_address": 			{"type": "keyword",	"normalizer": "caseNormalizer"},
      "government_nexus": 		{"type": "keyword",	"normalizer": "caseNormalizer"},
      "incident_date": 			{"type": "keyword",	"normalizer": "caseNormalizer"},
      "open_date": 				{"type": "keyword",	"normalizer": "caseNormalizer"},
      "primary_issue": 			{"type": "keyword",	"normalizer": "caseNormalizer"},
      "primary_issue_layer1": 	{"type": "keyword",	"normalizer": "caseNormalizer"},
      "primary_issue_layer2": 	{"type": "keyword",	"normalizer": "caseNormalizer"},
      "primary_issue_layer3": 	{"type": "keyword",	"normalizer": "caseNormalizer"},
      "primary_outcome": 		{"type": "keyword",	"normalizer": "caseNormalizer"},
      "reported_date": 			{"type": "keyword",	"normalizer": "caseNormalizer"},
      "reporter_is_employee": 	{"type": "keyword",	"normalizer": "caseNormalizer"},
      "reporter_name_first": 	{"type": "keyword",	"normalizer": "caseNormalizer"},
      "reporter_name_last": 	{"type": "keyword",	"normalizer": "caseNormalizer"},
      "secondary_issue": 		{"type": "keyword",	"normalizer": "caseNormalizer"},
      "secondary_issue_layer1": {"type": "keyword",	"normalizer": "caseNormalizer"},
      "secondary_issue_layer2": {"type": "keyword",	"normalizer": "caseNormalizer"},
      "secondary_issue_layer3": {"type": "keyword",	"normalizer": "caseNormalizer"},
      "secondary_outcome": 		{"type": "keyword",	"normalizer": "caseNormalizer"},
      "summary": 				{"type": "text"                                   },
      "tertiary_issue": 		{"type": "keyword",	"normalizer": "caseNormalizer"},
      "tertiary_issue_layer1": 	{"type": "keyword",	"normalizer": "caseNormalizer"},
      "tertiary_issue_layer2": 	{"type": "keyword",	"normalizer": "caseNormalizer"},
      "tertiary_issue_layer3": 	{"type": "keyword",	"normalizer": "caseNormalizer"},
      "tertiary_outcome": 		{"type": "keyword",	"normalizer": "caseNormalizer"}
    }
  }
}`
