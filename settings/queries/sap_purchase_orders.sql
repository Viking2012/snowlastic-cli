SELECT 'PROD_LI.' ||
       'COMMON_DATA.' ||
       'PURCHASE_ORDERS'           AS "record source"
     , DB                          AS "database"
     , ICM_ID || '|' || EKPO_EBELP AS "id"
     , ICM_ID                      AS "icm id"
     , EKPO_EBELP                  AS "document line number"
     , ICM_VENDOR_ID               AS "entity icm id"
     , 'Vendor'                    AS "entity title"
     , EKKO_LIFNR                  AS "entity number"
     , "EKKO_LIFNR_Concat"         AS "entity name"
     , EKKO_BSTYP                  AS "document category code"
     , "EKKO_BSTYP_Concat"         AS "document category"
     , EKKO_BSART                  AS "document type code"
     , "EKKO_BSART_Concat"         AS "document type"
     , EKKN_AEDAT                  AS "creation date"
     , EKKO_ERNAM                  AS "created by"
     , EKKO_BEDAT                  AS "document date"
     , "EKKN_Value_USD"            AS "document value"
     , EKKO_WAERS                  AS "document currency"
     , EKET_MENGE                  AS "document quantity"
     , EKPO_TXZ01                  AS "document text"
     , EKPO_MATKL                  AS "material group code"
     , "EKPO_MATKL_Concat"         AS "material group"
     , EKPO_MATNR                  AS "material code"
     , "EKPO_MATNR_Concat"         AS "material"
     , EKKO_BUKRS                  AS "company code id"
     , "EKKO_BUKRS_Concat"         AS "company code"
     , "ProfitCenter_ORG1"         AS "organization, level 1 code"
     , "ProfitCenter_ORG1_Concat"  AS "organization, level 1"
     , "ProfitCenter_ORG2"         AS "organization, level 2 code"
     , "ProfitCenter_ORG2_Concat"  AS "organization, level 2"
     , "ProfitCenter_ORG3"         AS "organization, level 3 code"
     , "ProfitCenter_ORG3_Concat"  AS "organization, level 3"
     , EKKO_LOEKZ                  AS "deletion indicator"
  FROM PROD_LI.COMMON_DATA.PURCHASE_ORDERS