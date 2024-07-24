SELECT 'COMMON_DATA.PAYMENT_RUNS'    AS "record source"
     , DB                            AS "database"
     , ICM_ID                        AS "icm id"
     , NULL                          AS "document line number"
     , ICM_VENDOR_ID                 AS "entity icm id"
     , 'Vendor'                      AS "entity title"
     , REGUH_LIFNR                   AS "entity number"
     , "LFA1_LIFNR_Concat"           AS "entity name"
     , NULL                          AS "document category code"
     , NULL                          AS "document category"
     , 'REGUH'                       AS "document type code"
     , 'Payment Run'                 AS "document type"
     , REGUH_LAUFD                   AS "creation date"
     , NULL                          AS "created by"
     , REGUH_ZALDT                   AS "document date"
     , REGUH_RWBTR_USD               AS "document value"
     , REGUH_WAERS                   AS "document currency"
     , NULL                          AS "document quantity"
     , RTRIM(CONCAT_WS(' ',
                       REGUH_ZNME1,
                       REGUH_ZNME2,
                       REGUH_ZNME3,
                       REGUH_ZNME4)) AS "document text"
     , NULL                          AS "material group code"
     , NULL                          AS "material group"
     , NULL                          AS "material code"
     , NULL                          AS "material"
     , REGUP_BUKRS                   AS "company code id"
     , T001.CONCAT                   AS "company code"
     , PCM.ORG1                      AS "organization, level 1 code"
     , PCM.ORG1_CONCAT               AS "organization, level 1"
     , PCM.ORG2                      AS "organization, level 2 code"
     , PCM.ORG2_CONCAT               AS "organization, level 2"
     , PCM.ORG3                      AS "organization, level 3 code"
     , PCM.ORG3_CONCAT               AS "organization, level 3"
     , REGUH_XVORL                   AS "deletion indicator"
  FROM COMMON_DATA.PAYMENT_RUNS
           INNER JOIN PROD_UDM.CONFIGURATION.SOURCE_SYSTEM SS
                      ON PAYMENT_RUNS.DB = RTRIM(SS.SYSTEM_NAME)
           LEFT JOIN COMMON_DATA_STAGING.T001
                     ON SS.SYSTEM_ID = T001.SOURCESYSTEMID AND
                        PAYMENT_RUNS.REGUP_BUKRS = T001.BUKRS
           LEFT JOIN BA_MASTERDATA.PROFITCENTERMAPPING PCM
                     ON PAYMENT_RUNS.REGUP_BUKRS = PCM.BUKRS AND
                        PAYMENT_RUNS.REGUP_PRCTR = PCM.PRCTR
 WHERE RTRIM(REGUH_LIFNR) != ''
 UNION ALL
SELECT 'COMMON_DATA.PAYMENT_RUNS'    AS "record source"
     , DB                            AS "database"
     , ICM_ID                        AS "icm id"
     , NULL                          AS "document line number"
     , ICM_VENDOR_ID                 AS "entity icm id"
     , 'Vendor'                      AS "entity title"
     , REGUH_LIFNR                   AS "entity number"
     , "LFA1_LIFNR_Concat"           AS "entity name"
     , NULL                          AS "document category code"
     , NULL                          AS "document category"
     , 'REGUH'                       AS "document type code"
     , 'Payment Run'                 AS "document type"
     , REGUH_LAUFD                   AS "creation date"
     , NULL                          AS "created by"
     , REGUH_ZALDT                   AS "document date"
     , REGUH_RWBTR_USD               AS "document value"
     , REGUH_WAERS                   AS "document currency"
     , NULL                          AS "document quantity"
     , RTRIM(CONCAT_WS(' ',
                       REGUH_ZNME1,
                       REGUH_ZNME2,
                       REGUH_ZNME3,
                       REGUH_ZNME4)) AS "document text"
     , NULL                          AS "material group code"
     , NULL                          AS "material group"
     , NULL                          AS "material code"
     , NULL                          AS "material"
     , REGUP_BUKRS                   AS "company code id"
     , T001.CONCAT                   AS "company code"
     , PCM.ORG1                      AS "organization, level 1 code"
     , PCM.ORG1_CONCAT               AS "organization, level 1"
     , PCM.ORG2                      AS "organization, level 2 code"
     , PCM.ORG2_CONCAT               AS "organization, level 2"
     , PCM.ORG3                      AS "organization, level 3 code"
     , PCM.ORG3_CONCAT               AS "organization, level 3"
     , REGUH_XVORL                   AS "deletion indicator"
  FROM COMMON_DATA.PAYMENT_RUNS
           INNER JOIN PROD_UDM.CONFIGURATION.SOURCE_SYSTEM SS
                      ON PAYMENT_RUNS.DB = RTRIM(SS.SYSTEM_NAME)
           LEFT JOIN COMMON_DATA_STAGING.T001
                     ON SS.SYSTEM_ID = T001.SOURCESYSTEMID AND
                        PAYMENT_RUNS.REGUP_BUKRS = T001.BUKRS
           LEFT JOIN BA_MASTERDATA.PROFITCENTERMAPPING PCM
                     ON PAYMENT_RUNS.REGUP_BUKRS = PCM.BUKRS AND
                        PAYMENT_RUNS.REGUP_PRCTR = PCM.PRCTR
 WHERE RTRIM(REGUH_LIFNR) = ''