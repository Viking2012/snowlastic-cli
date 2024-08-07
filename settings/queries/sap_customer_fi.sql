SELECT 'PROD_LI.' ||
       'COMMON_DATA.' ||
       'CUSTOMER_FI'     AS "record source"
     , DB                            AS "database"
     , ICM_ID || '|' || BSAD_BUZEI   AS "id"
     , ICM_ID                        AS "icm id"
     , BSAD_BUZEI                    AS "document line number"
     , ICM_CUSTOMER_ID               AS "entity icm id"
     , 'Customer'                    AS "entity title"
     , BSAD_KUNNR                    AS "entity number"
     , "BSAD_KUNNR_Concat"           AS "entity name"
     , NULL                          AS "document category code"
     , NULL                          AS "document category"
     , BSAD_BLART                    AS "document type code"
     , "BSAD_BLART_Concat"           AS "document type"
     , BSAD_CPUDT                    AS "creation date"
     , BKPF_USNAM                    AS "created by"
     , BSAD_BLDAT                    AS "document date"
     , FAGLFLEXA_HSL_USD             AS "document value"
     , BSAD_WAERS                    AS "document currency"
     , NULL                          AS "document quantity"
     , TRIM(
        CONCAT_WS(
                '\n',
                IFNULL(CONCAT('HEADER: ', NULLIF(BKPF_BKTXT, ''), '\n'), ''),
                IFNULL(CONCAT('ITEM: ', NULLIF(BSAD_SGTXT, ''), '\n'), '')),
        '\n')                        AS "document text"
     , NULL                          AS "material group code"
     , NULL                          AS "material group"
     , NULL                          AS "material code"
     , NULL                          AS "material"
     , BSAD_BUKRS                    AS "company code id"
     , "BSAD_BUKRS_Concat"           AS "company code"
     , FAGLFLEXA_PRCTR_ORG1          AS "organization, level 1 code"
     , "FAGLFLEXA_PRCTR_ORG1_Concat" AS "organization, level 1"
     , FAGLFLEXA_PRCTR_ORG2          AS "organization, level 2 code"
     , "FAGLFLEXA_PRCTR_ORG2_Concat" AS "organization, level 2"
     , FAGLFLEXA_PRCTR_ORG3          AS "organization, level 3 code"
     , "FAGLFLEXA_PRCTR_ORG3_Concat" AS "organization, level 3"
     , NULL                          AS "deletion indicator"
  FROM PROD_LI.COMMON_DATA.CUSTOMER_FI