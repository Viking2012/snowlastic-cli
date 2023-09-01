package types

import (
	"database/sql"
	icmEncoding "github.com/alexander-orban/icm_goapi/encoding"
	icmOrm "github.com/alexander-orban/icm_goapi/orm"
	"github.com/spf13/viper"
	"strings"
)

type PurchaseOrder struct {
	icmOrm.PurchaseOrder
}

func (po *PurchaseOrder) IsDocument()   {}
func (po *PurchaseOrder) GetID() string { return po.PurchaseOrder.GetID() }
func (po *PurchaseOrder) GetQuery(schema, table string) string {
	var database = strings.Join([]string{
		viper.GetString("snowflakeDatabase"),
		schema,
		table,
	}, ".")
	return string(icmEncoding.MarshalToSelect(&po.PurchaseOrder, database, false))
}
func (po *PurchaseOrder) ScanFrom(rows *sql.Rows) error {
	return rows.Scan(
		po.PurchaseOrder.Database,
		po.PurchaseOrder.PoNumber,
		po.PurchaseOrder.PurchasingDocumentLineItem,
		po.PurchaseOrder.AcctAssignmentLine,
		po.PurchaseOrder.CompanyCodeId,
		po.PurchaseOrder.LocalCurrency,
		po.PurchaseOrder.LocalLanguage,
		po.PurchaseOrder.ChartOfAccounts,
		po.PurchaseOrder.CompanyCode,
		po.PurchaseOrder.PurchasingDocumentCategoryCode,
		po.PurchaseOrder.PurchasingDocumentCategory,
		po.PurchaseOrder.PurchasingDocumentTypeCode,
		po.PurchaseOrder.PurchasingDocumentType,
		po.PurchaseOrder.ControlIndicatorForPurchasingDocumentType,
		po.PurchaseOrder.PoDeletionIndicator,
		po.PurchaseOrder.PoStatusCode,
		po.PurchaseOrder.PoStatus,
		po.PurchaseOrder.PoCreationDate,
		po.PurchaseOrder.PoCreatedBy,
		po.PurchaseOrder.VendorNumber,
		po.PurchaseOrder.Vendor,
		po.PurchaseOrder.CountryCode,
		po.PurchaseOrder.AccountGroupCode,
		po.PurchaseOrder.AccountGroup,
		po.PurchaseOrder.TermsOfPaymentKey,
		po.PurchaseOrder.CashPromptPaymentDiscount1Days,
		po.PurchaseOrder.CashPromptPaymentDiscount2Days,
		po.PurchaseOrder.CashPromptPaymentDiscount3Days,
		po.PurchaseOrder.CashDiscountPercentage1,
		po.PurchaseOrder.CashDiscountPercentage2,
		po.PurchaseOrder.PurchasingOrganizationCode,
		po.PurchaseOrder.PurchasingOrganization,
		po.PurchaseOrder.PurchasingGroupCode,
		po.PurchaseOrder.PurchasingGroup,
		po.PurchaseOrder.PoCurrency,
		po.PurchaseOrder.ExchangeRate,
		po.PurchaseOrder.PurchasingDocumentDate,
		po.PurchaseOrder.SupplyingVendorCode,
		po.PurchaseOrder.NumberOfPrincipalPurchaseAgreementPOHeader,
		po.PurchaseOrder.IncotermsPart1PO,
		po.PurchaseOrder.IncotermsPart2PO,
		po.PurchaseOrder.ReleaseGroup,
		po.PurchaseOrder.ReleaseStrategy,
		po.PurchaseOrder.PurchasingDocumentReleaseIndicator,
		po.PurchaseOrder.PoReleaseStatus,
		po.PurchaseOrder.MostRecentApproverCode,
		po.PurchaseOrder.MostRecentApprover,
		po.PurchaseOrder.NextApproverCode,
		po.PurchaseOrder.NextApprover,
		po.PurchaseOrder.ReleaseNotYetCompletelyEffected,
		po.PurchaseOrder.AddressId,
		po.PurchaseOrder.VatNumber,
		po.PurchaseOrder.DeletionIndicatorInPurchasingDocumentItem,
		po.PurchaseOrder.RfqStatus,
		po.PurchaseOrder.PurchasingDocumentItemChangeDate,
		po.PurchaseOrder.ShortTextForPurchasingDocumentItem,
		po.PurchaseOrder.MaterialNumber,
		po.PurchaseOrder.Ekpo_Matnr_Concat,
		po.PurchaseOrder.Mpn,
		po.PurchaseOrder.PlantCode,
		po.PurchaseOrder.Plant,
		po.PurchaseOrder.StorageLocation,
		po.PurchaseOrder.MaterialGroupCode,
		po.PurchaseOrder.MaterialGroup,
		po.PurchaseOrder.PoInfoRecord,
		po.PurchaseOrder.VendorMaterialNumber,
		po.PurchaseOrder.PurchaseOrderUnitOfMeasure,
		po.PurchaseOrder.OverdeliveryToleranceLimit,
		po.PurchaseOrder.IndicatorUnlimitedOverdeliveryAllowed,
		po.PurchaseOrder.UnderdeliveryToleranceLimit,
		po.PurchaseOrder.DeliveryCompletedIndicator,
		po.PurchaseOrder.FinalInvoiceIndicator,
		po.PurchaseOrder.ItemCategoryInPurchasingDocument,
		po.PurchaseOrder.PoLineItemAccountAssignmentCategory,
		po.PurchaseOrder.ConsumptionPosting,
		po.PurchaseOrder.DistributionIndicatorForMultipleAccountAssignment,
		po.PurchaseOrder.GoodsReceiptIndicator,
		po.PurchaseOrder.GoodsReceiptNonValuated,
		po.PurchaseOrder.InvoiceReceiptIndicator,
		po.PurchaseOrder.GrBasedInvoiceVerification,
		po.PurchaseOrder.AcknowledgementRequired,
		po.PurchaseOrder.PoLineAddressId,
		po.PurchaseOrder.Ekpo_Adrn2,
		po.PurchaseOrder.ConfirmationControl,
		po.PurchaseOrder.Incoterms1,
		po.PurchaseOrder.Incoterms2,
		po.PurchaseOrder.Statistical,
		po.PurchaseOrder.PrNumber,
		po.PurchaseOrder.PrItem,
		po.PurchaseOrder.ReturnsItem,
		po.PurchaseOrder.FirstDeliveryDate,
		po.PurchaseOrder.LastDeliveryDate,
		po.PurchaseOrder.Eket_Slfdt_Min,
		po.PurchaseOrder.Eket_Slfdt_Max,
		po.PurchaseOrder.ScheduledQuantity,
		po.PurchaseOrder.AcctAssignmentDeleted,
		po.PurchaseOrder.AcctAssignmentCreatedOn,
		po.PurchaseOrder.Ekkn_Vproz,
		po.PurchaseOrder.GlAccountId,
		po.PurchaseOrder.CostCenter,
		po.PurchaseOrder.Ekkn_Aufnr,
		po.PurchaseOrder.Ekkn_Kstrg,
		po.PurchaseOrder.WbsElementId,
		po.PurchaseOrder.ProjectNumber,
		po.PurchaseOrder.ProjectName,
		po.PurchaseOrder.Network,
		po.PurchaseOrder.RoutingLine,
		po.PurchaseOrder.RoutingNumber,
		po.PurchaseOrder.Ekkn_Vbeln,
		po.PurchaseOrder.Ekkn_Vbelp,
		po.PurchaseOrder.Ekkn_Anln1,
		po.PurchaseOrder.Ekkn_Anln2,
		po.PurchaseOrder.Quantity,
		po.PurchaseOrder.Ekkn_Qty_Base,
		po.PurchaseOrder.Value,
		po.PurchaseOrder.ValueInLc,
		po.PurchaseOrder.ValueInUsd,
		po.PurchaseOrder.ExchangeRateToLc,
		po.PurchaseOrder.ExchangeRateToUsd,
		po.PurchaseOrder.AcctAssignmentRatio,
		po.PurchaseOrder.ProfitCenter,
		po.PurchaseOrder.Profitcenter_Org1,
		po.PurchaseOrder.Profitcenter_Org1_Concat,
		po.PurchaseOrder.Profitcenter_Org2,
		po.PurchaseOrder.Profitcenter_Org2_Concat,
		po.PurchaseOrder.Profitcenter_Org3,
		po.PurchaseOrder.Profitcenter_Org3_Concat,
		po.PurchaseOrder.AbacusCountry,
		po.PurchaseOrder.ReportingUnit,
		po.PurchaseOrder.InvoicedQuantity,
		po.PurchaseOrder.InvoicedValueInLc,
		po.PurchaseOrder.FirstInvoicePostingDate,
		po.PurchaseOrder.FirstInvoiceCreatedOn,
		po.PurchaseOrder.LastInvoicePostingDate,
		po.PurchaseOrder.LastInvoiceCreatedOn,
		po.PurchaseOrder.FirstInvoiceDocDate,
		po.PurchaseOrder.LastInvoiceDocDate,
		po.PurchaseOrder.ReceiptedQuantity,
		po.PurchaseOrder.ReceiptedValueInLc,
		po.PurchaseOrder.FirstReceiptPostingDate,
		po.PurchaseOrder.FirstReceiptCreatedOn,
		po.PurchaseOrder.LastReceiptPostingDate,
		po.PurchaseOrder.LastReceiptCreatedOn,
		po.PurchaseOrder.FirstReceiptDocDate,
		po.PurchaseOrder.LastReceiptDocDate,
		po.PurchaseOrder.ApprovalsRequired,
		po.PurchaseOrder.ID,
		po.PurchaseOrder.ApprovalDescriptive,
		po.PurchaseOrder.TotalApprovals,
		po.PurchaseOrder.Flags,
	)
}
func (po *PurchaseOrder) New() SnowlasticDocument {
	return &PurchaseOrder{icmOrm.PurchaseOrder{Flags: &icmOrm.PurchaseOrderFlags{}}}
}

//type _ struct {
//	icmOrm._
//}
//
//func (_ *_) IsDocument()                   {}
//func (_ *_) GetID() string                 { return _.GetID() }
//func (_ *_) GetQuery(schema string, table string) string        { return "" }
//func (_ *_) ScanFrom(rows *sql.Rows) error { return rows.Scan() }
//func (_ *_) New() SnowlasticDocument       { return nil }
