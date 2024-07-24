SELECT 'PROD_LI.' ||
       'COMMON_DATA.' ||
       'REQUISITIONS'    AS "record source"
     , DB                            AS "database"
     , CONCAT_WS('|', DB, BANFN)     AS "icm id"
     , BNFPO                         AS "document line number"
     , ICM_VENDOR_ID                 AS "entity icm id"
     , 'Vendor'                      AS "entity title"
     , LIFNR                         AS "entity number"
     , RTRIM(CONCAT_WS(' ',
                       LIFNR_NAME1,
                       LIFNR_NAME2)) AS "entity name"
     , BSART                         AS "document category code"
     , BSART_BATXT                   AS "document category"
     , BSTYP                         AS "document type code"
     , "BSTYP_Concat"                AS "document type"
     , ERDAT                         AS "creation date"
     , ERNAM                         AS "created by"
     , BADAT                         AS "document date"
     , "Value_USD"                   AS "document value"
     , WAERS                         AS "document currency"
     , "Qty_Base"                    AS "document quantity"
     , TXZ01                         AS "document text"
     , MATKL                         AS "material group code"
     , MATKL_WGBEZ                   AS "material group"
     , MATNR                         AS "material code"
     , MATNR                         AS "material"
     , WERKS_EKORG                   AS "company code id"
     , WERKS_EKOTX                   AS "company code"
     , ORG1                          AS "organization, level 1 code"
     , "ORG1_Concat"                 AS "organization, level 1"
     , ORG2                          AS "organization, level 2 code"
     , "ORG2_Concat"                 AS "organization, level 2"
     , ORG3                          AS "organization, level 3 code"
     , "ORG3_Concat"                 AS "organization, level 3"
     , LOEKZ                         AS "deletion indicator"
  FROM PROD_LI.COMMON_DATA.REQUISITIONS