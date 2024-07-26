SELECT 'PROD_LI.' ||
       'COMMON_DATA.' ||
       'VENDOR_FI'                   AS "record source"
     , DB                            AS "database"
     , ICM_ID || '|' || BSAK_BUZEI   AS "id"
     , ICM_ID                        AS "icm id"
     , BSAK_BUZEI                    AS "document line number"
     , ICM_VENDOR_ID                 AS "entity icm id"
     , 'Vendor'                      AS "entity title"
     , BSAK_LIFNR                    AS "entity number"
     , "BSAK_LIFNR_Concat"           AS "entity name"
     , NULL                          AS "document category code"
     , NULL                          AS "document category"
     , BSAK_BLART                    AS "document type code"
     , "BSAK_BLART_Concat"           AS "document type"
     , BSAK_CPUDT                    AS "creation date"
     , BKPF_USNAM                    AS "created by"
     , BSAK_BLDAT                    AS "document date"
     , FAGLFLEXA_HSL_USD             AS "document value"
     , BSAK_WAERS                    AS "document currency"
     , NULL                          AS "document quantity"
     , TRIM(
        CONCAT_WS(
                '\n',
                IFNULL(CONCAT('HEADER: ', NULLIF(BKPF_BKTXT, ''), '\n'), ''),
                IFNULL(CONCAT('ITEM: ', NULLIF(BSAK_SGTXT, ''), '\n'), '')),
        '\n')                        AS "document text"
     , NULL                          AS "material group code"
     , NULL                          AS "material group"
     , NULL                          AS "material code"
     , NULL                          AS "material"
     , BSAK_BUKRS                    AS "company code id"
     , "BSAK_BUKRS_Concat"           AS "company code"
     , FAGLFLEXA_PRCTR_ORG1          AS "organization, level 1 code"
     , "FAGLFLEXA_PRCTR_ORG1_Concat" AS "organization, level 1"
     , FAGLFLEXA_PRCTR_ORG2          AS "organization, level 2 code"
     , "FAGLFLEXA_PRCTR_ORG2_Concat" AS "organization, level 2"
     , FAGLFLEXA_PRCTR_ORG3          AS "organization, level 3 code"
     , "FAGLFLEXA_PRCTR_ORG3_Concat" AS "organization, level 3"
     , NULL                          AS "deletion indicator"
  FROM PROD_LI.COMMON_DATA.VENDOR_FI