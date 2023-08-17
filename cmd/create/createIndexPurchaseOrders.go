package create

import (
	"bytes"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
)

func indexPurchaseOrder(c *elasticsearch.Client) error {
	res, err := c.Indices.Delete([]string{"purchaseOrders"})
	if err != nil {
		return fmt.Errorf("cannot delete index: %s", err)
	}
	if res.IsError() {
		log.Println("error when deleting index", res.String())
	} else {
		log.Println(res.String())
	}

	b := []byte(purchaseOrderIndex)
	res, err = c.Indices.Create("purchaseOrders", c.Indices.Create.WithBody(bytes.NewReader(b)))
	if err != nil {
		return fmt.Errorf("cannot create index: %s", err)
	}

	if res.IsError() {
		return fmt.Errorf("cannot create index, got an error response code: %s\n", res.String())
	}
	log.Println("successfully created an index")
	return nil
}

var purchaseOrderIndex = `{
  "settings": {
    "analysis": {
      "analyzer": {
        "poAnalyzer": {
          "char_filter": ["html_strip"],
          "tokenizer": "standard",
          "filter": ["lowercase","stop","asciifolding","stemmer","edge_ngram"]
        }
      },
      "normalizer": {
        "poNormalizer": {
          "type": "custom",
          "filter": ["lowercase","asciifolding"]
        }
      } 
    }
  },
  "mappings": {
    "properties": {
      "Database":                                               {"type": "keyword",	"normalizer": "poNormalizer"},
      "PO Number":                                              {"type": "keyword",	"normalizer": "poNormalizer"},
      "Purchasing Document Line Item":                          {"type": "keyword",	"normalizer": "poNormalizer"},
      "Acct Assignment Line":                                   {"type": "keyword",	"normalizer": "poNormalizer"},
      "Company Code ID":                                        {"type": "keyword",	"normalizer": "poNormalizer"},
      "Local Currency":                                         {"type": "keyword",	"normalizer": "poNormalizer"},
      "Local Language":                                         {"type": "keyword",	"normalizer": "poNormalizer"},
      "Chart of Accounts":                                      {"type": "keyword",	"normalizer": "poNormalizer"},
      "Company Code":                                           {"type": "keyword",	"normalizer": "poNormalizer"},
      "Purchasing Document Category Code":                      {"type": "keyword",	"normalizer": "poNormalizer"},
      "Purchasing Document Category":                           {"type": "keyword",	"normalizer": "poNormalizer"},
      "Purchasing Document Type Code":                          {"type": "keyword",	"normalizer": "poNormalizer"},
      "Purchasing Document Type":                               {"type": "keyword",	"normalizer": "poNormalizer"},
      "Control indicator for purchasing document type":         {"type": "keyword",	"normalizer": "poNormalizer"},
      "PO Deletion Indicator":                                  {"type": "keyword",	"normalizer": "poNormalizer"},
      "PO Status Code":                                         {"type": "keyword",	"normalizer": "poNormalizer"},
      "PO Status":                                              {"type": "keyword",	"normalizer": "poNormalizer"},
      "PO Creation Date":                                       {"type": "date",	"normalizer": "poNormalizer"},
      "PO Created By":                                          {"type": "keyword",	"normalizer": "poNormalizer"},
      "Vendor Number":                                          {"type": "keyword",	"normalizer": "poNormalizer"},
      "Vendor":                                                 {"type": "keyword",	"normalizer": "poNormalizer"},
      "Vendor Country Code":                                    {"type": "keyword",	"normalizer": "poNormalizer"},
      "Vendor Account Group Code":                              {"type": "keyword",	"normalizer": "poNormalizer"},
      "Vendor Account Group":                                   {"type": "keyword",	"normalizer": "poNormalizer"},
      "Terms of Payment Key":                                   {"type": "keyword",	"normalizer": "poNormalizer"},
      "Cash Prompt Payment Discount (1) Days":                  {"type": "number",	"normalizer": "poNormalizer"},
      "Cash Prompt Payment Discount (2) Days":                  {"type": "number",	"normalizer": "poNormalizer"},
      "Cash Prompt Payment Discount (3) Days":                  {"type": "number",	"normalizer": "poNormalizer"},
      "Cash Discount Percentage 1":                             {"type": "number",	"normalizer": "poNormalizer"},
      "Cash Discount Percentage 2":                             {"type": "number",	"normalizer": "poNormalizer"},
      "Purchasing Organization Code":                           {"type": "keyword",	"normalizer": "poNormalizer"},
      "Purchasing Organization":                                {"type": "keyword",	"normalizer": "poNormalizer"},
      "Purchasing Group Code":                                  {"type": "keyword",	"normalizer": "poNormalizer"},
      "Purchasing Group":                                       {"type": "keyword",	"normalizer": "poNormalizer"},
      "PO Currency":                                            {"type": "keyword",	"normalizer": "poNormalizer"},
      "Exchange Rate":                                          {"type": "number",	"normalizer": "poNormalizer"},
      "Purchasing Document Date":                               {"type": "date",	"normalizer": "poNormalizer"},
      "Supplying Vendor Code":                                  {"type": "keyword",	"normalizer": "poNormalizer"},
      "Number of Principal Purchase Agreement (PO Header)":     {"type": "keyword",	"normalizer": "poNormalizer"},
      "Incoterms part 1 (PO)":                                  {"type": "keyword",	"normalizer": "poNormalizer"},
      "Incoterms part 2 (PO)":                                  {"type": "keyword",	"normalizer": "poNormalizer"},
      "Release Group":                                          {"type": "keyword",	"normalizer": "poNormalizer"},
      "Release Strategy":                                       {"type": "keyword",	"normalizer": "poNormalizer"},
      "Purchasing Document Release Indicator":                  {"type": "keyword",	"normalizer": "poNormalizer"},
      "PO Release Status":                                      {"type": "keyword",	"normalizer": "poNormalizer"},
      "Most Recent Approver Code":                              {"type": "keyword",	"normalizer": "poNormalizer"},
      "Most Recent Approver":                                   {"type": "keyword",	"normalizer": "poNormalizer"},
      "Next Approver Code":                                     {"type": "keyword",	"normalizer": "poNormalizer"},
      "Next Approver":                                          {"type": "keyword",	"normalizer": "poNormalizer"},
      "Release Not Yet Completely Effected":                    {"type": "keyword",	"normalizer": "poNormalizer"},
      "Address ID":                                             {"type": "keyword",	"normalizer": "poNormalizer"},
      "VAT Number":                                             {"type": "keyword",	"normalizer": "poNormalizer"},
      "Deletion Indicator in Purchasing Document Item":         {"type": "keyword",	"normalizer": "poNormalizer"},
      "RFQ status":                                             {"type": "keyword",	"normalizer": "poNormalizer"},
      "Purchasing Document Item Change Date":                   {"type": "date",	"normalizer": "poNormalizer"},
      "Short Text for Purchasing Document Item":                {"type": "text",	"normalizer": "poNormalizer"},
      "Material Number":                                        {"type": "keyword",	"normalizer": "poNormalizer"},
      "EKPO_MATNR_Concat":                                      {"type": "keyword",	"normalizer": "poNormalizer"},
      "MPN":                                                    {"type": "keyword",	"normalizer": "poNormalizer"},
      "Plant Code":                                             {"type": "keyword",	"normalizer": "poNormalizer"},
      "Plant":                                                  {"type": "keyword",	"normalizer": "poNormalizer"},
      "Storage Location":                                       {"type": "keyword",	"normalizer": "poNormalizer"},
      "Material Group Code":                                    {"type": "keyword",	"normalizer": "poNormalizer"},
      "Material Group":                                         {"type": "keyword",	"normalizer": "poNormalizer"},
      "PO Info Record":                                         {"type": "keyword",	"normalizer": "poNormalizer"},
      "Vendor Material Number":                                 {"type": "keyword",	"normalizer": "poNormalizer"},
      "Purchase Order Unit of Measure":                         {"type": "keyword",	"normalizer": "poNormalizer"},
      "Overdelivery Tolerance Limit":                           {"type": "number",	"normalizer": "poNormalizer"},
      "Indicator Unlimited Overdelivery Allowed":               {"type": "keyword",	"normalizer": "poNormalizer"},
      "Underdelivery Tolerance Limit":                          {"type": "number",	"normalizer": "poNormalizer"},
      "Delivery Completed Indicator":                           {"type": "keyword",	"normalizer": "poNormalizer"},
      "Final Invoice Indicator":                                {"type": "keyword",	"normalizer": "poNormalizer"},
      "Item Category in Purchasing Document":                   {"type": "keyword",	"normalizer": "poNormalizer"},
      "PO Line Item Account Assignment Category":               {"type": "keyword",	"normalizer": "poNormalizer"},
      "Consumption Posting":                                    {"type": "keyword",	"normalizer": "poNormalizer"},
      "Distribution indicator for multiple account assignment": {"type": "keyword",	"normalizer": "poNormalizer"},
      "Goods Receipt Indicator":                                {"type": "keyword",	"normalizer": "poNormalizer"},
      "Goods Receipt Non-Valuated":                             {"type": "keyword",	"normalizer": "poNormalizer"},
      "Invoice Receipt Indicator":                              {"type": "keyword",	"normalizer": "poNormalizer"},
      "GR-Based Invoice Verification":                          {"type": "keyword",	"normalizer": "poNormalizer"},
      "Acknowledgement Required":                               {"type": "keyword",	"normalizer": "poNormalizer"},
      "PO Line Address ID":                                     {"type": "keyword",	"normalizer": "poNormalizer"},
      "EKPO_ADRN2":                                             {"type": "keyword",	"normalizer": "poNormalizer"},
      "Confirmation Control":                                   {"type": "keyword",	"normalizer": "poNormalizer"},
      "Incoterms 1":                                            {"type": "keyword",	"normalizer": "poNormalizer"},
      "Incoterms 2":                                            {"type": "keyword",	"normalizer": "poNormalizer"},
      "Statistical":                                            {"type": "keyword",	"normalizer": "poNormalizer"},
      "PR Number":                                              {"type": "keyword",	"normalizer": "poNormalizer"},
      "PR Item":                                                {"type": "keyword",	"normalizer": "poNormalizer"},
      "Returns Item":                                           {"type": "keyword",	"normalizer": "poNormalizer"},
      "First Delivery Date":                                    {"type": "date",	"normalizer": "poNormalizer"},
      "Last Delivery Date":                                     {"type": "date",	"normalizer": "poNormalizer"},
      "EKET_SLFDT_Min":                                         {"type": "date",	"normalizer": "poNormalizer"},
      "EKET_SLFDT_Max":                                         {"type": "date",	"normalizer": "poNormalizer"},
      "Scheduled Quantity":                                     {"type": "number",	"normalizer": "poNormalizer"},
      "Acct Assignment Deleted":                                {"type": "keyword",	"normalizer": "poNormalizer"},
      "Acct Assignment Created On":                             {"type": "date",	"normalizer": "poNormalizer"},
      "EKKN_VPROZ":                                             {"type": "number",	"normalizer": "poNormalizer"},
      "GL Account ID":                                          {"type": "keyword",	"normalizer": "poNormalizer"},
      "Cost Center":                                            {"type": "keyword",	"normalizer": "poNormalizer"},
      "EKKN_AUFNR":                                             {"type": "keyword",	"normalizer": "poNormalizer"},
      "EKKN_KSTRG":                                             {"type": "keyword",	"normalizer": "poNormalizer"},
      "WBS Element ID":                                         {"type": "keyword",	"normalizer": "poNormalizer"},
      "Project Number":                                         {"type": "keyword",	"normalizer": "poNormalizer"},
      "PROJ_Concat":                                            {"type": "keyword",	"normalizer": "poNormalizer"},
      "Network":                                                {"type": "keyword",	"normalizer": "poNormalizer"},
      "Routing Line":                                           {"type": "keyword",	"normalizer": "poNormalizer"},
      "Routing Number":                                         {"type": "keyword",	"normalizer": "poNormalizer"},
      "EKKN_VBELN":                                             {"type": "keyword",	"normalizer": "poNormalizer"},
      "EKKN_VBELP":                                             {"type": "keyword",	"normalizer": "poNormalizer"},
      "EKKN_ANLN1":                                             {"type": "keyword",	"normalizer": "poNormalizer"},
      "EKKN_ANLN2":                                             {"type": "keyword",	"normalizer": "poNormalizer"},
      "Quantity":                                               {"type": "number",	"normalizer": "poNormalizer"},
      "EKKN_Qty_Base":                                          {"type": "number",	"normalizer": "poNormalizer"},
      "Value":                                                  {"type": "number",	"normalizer": "poNormalizer"},
      "Value in LC":                                            {"type": "number",	"normalizer": "poNormalizer"},
      "Value in USD":                                           {"type": "number",	"normalizer": "poNormalizer"},
      "Exchange Rate to LC":                                    {"type": "number",	"normalizer": "poNormalizer"},
      "Exchange Rate to USD":                                   {"type": "number",	"normalizer": "poNormalizer"},
      "Acct Assignment Ratio":                                  {"type": "number",	"normalizer": "poNormalizer"},
      "Profit Center":                                          {"type": "keyword",	"normalizer": "poNormalizer"},
      "ProfitCenter_ORG1":                                      {"type": "keyword",	"normalizer": "poNormalizer"},
      "ProfitCenter_ORG1_Concat":                               {"type": "keyword",	"normalizer": "poNormalizer"},
      "ProfitCenter_ORG2":                                      {"type": "keyword",	"normalizer": "poNormalizer"},
      "ProfitCenter_ORG2_Concat":                               {"type": "keyword",	"normalizer": "poNormalizer"},
      "ProfitCenter_ORG3":                                      {"type": "keyword",	"normalizer": "poNormalizer"},
      "ProfitCenter_ORG3_Concat":                               {"type": "keyword",	"normalizer": "poNormalizer"},
      "Abacus Country":                                         {"type": "keyword",	"normalizer": "poNormalizer"},
      "Reporting Unit":                                         {"type": "keyword",	"normalizer": "poNormalizer"},
      "Invoiced Quantity":                                      {"type": "number",	"normalizer": "poNormalizer"},
      "Invoiced Value in LC":                                   {"type": "number",	"normalizer": "poNormalizer"},
      "First Invoice Posting Date":                             {"type": "date",	"normalizer": "poNormalizer"},
      "First Invoice Created On":                               {"type": "date",	"normalizer": "poNormalizer"},
      "Last Invoice Posting Date":                              {"type": "date",	"normalizer": "poNormalizer"},
      "Last Invoice Created On":                                {"type": "date",	"normalizer": "poNormalizer"},
      "First Invoice Doc Date":                                 {"type": "date",	"normalizer": "poNormalizer"},
      "Last Invoice Doc Date":                                  {"type": "date",	"normalizer": "poNormalizer"},
      "Receipted Quantity":                                     {"type": "number",	"normalizer": "poNormalizer"},
      "Receipted Value in LC":                                  {"type": "number",	"normalizer": "poNormalizer"},
      "First Receipt Posting Date":                             {"type": "date",	"normalizer": "poNormalizer"},
      "First Receipt Created On":                               {"type": "date",	"normalizer": "poNormalizer"},
      "Last Receipt Posting Date":                              {"type": "date",	"normalizer": "poNormalizer"},
      "Last Receipt Created On":                                {"type": "date",	"normalizer": "poNormalizer"},
      "First Receipt Doc Date":                                 {"type": "date",	"normalizer": "poNormalizer"},
      "Last Receipt Doc Date":                                  {"type": "date",	"normalizer": "poNormalizer"},
      "Approvals Required":                                     {"type": "number",	"normalizer": "poNormalizer"},
      "ICM_VENDOR_ID":                                          {"type": "keyword",	"normalizer": "poNormalizer"},
      "Aproval Descriptive":                                    {"type": "keyword",	"normalizer": "poNormalizer"},
      "Total Approvals":                                        {"type": "number",	"normalizer": "poNormalizer"},
      "Flags":                                                  {"type":  "nested","properties": {
                                                                  "Flag: NoInvoiceNeeded":          {"type": "bool"},
                                                                  "Flag: NoReceiptNeeded":          {"type": "bool"},
                                                                  "Flag: NoInvoiceOrReceiptNeeded": {"type": "bool"}
                                                                }}
    }
  }
}
`
