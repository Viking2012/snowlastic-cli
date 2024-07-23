SELECT DB                     AS "database"
     , ICM_ID                 AS "document id"
     , NULL                   AS "document line number"
     , REGUP_BUKRS            AS "company code id"
     , T001.CONCAT            AS "company code"
     , NULL                   AS "document category code"
     , NULL                   AS "document category"
     , 'REGUH'                AS "document type code"
     , 'Payment Run'          AS "document type"
     , REGUH_XVORL            AS "deletion indicator"
     , REGUH_LAUFD            AS "creation date"
     , NULL                   AS "created by"
     , 'Vendor'               AS "entity type"
     , ICM_VENDOR_ID          AS "entity icm id"
     , REGUH_LIFNR            AS "entity number"
     , "LFA1_LIFNR_Concat"    AS "entity"
     , REGUH_ZALDT            AS "document date"
     , CONCAT_WS(' ',
                 REGUH_ZNME1,
                 REGUH_ZNME2,
                 REGUH_ZNME3,
                 REGUH_ZNME4) AS "document text"
     , REGUH_WAERS            AS "document currency"
     , REGUH_RWBTR_USD        AS "document value"
     , NULL                   AS "document quantity"
     , NULL                   AS "material code"
     , NULL                   AS "material"
     , NULL                   AS "material group code"
     , NULL                   AS "material group"
     , NULL                   AS "organization, level 1 code"
     , NULL                   AS "organization, level 1"
     , NULL                   AS "organization, level 2 code"
     , NULL                   AS "organization, level 2"
     , NULL                   AS "organization, level 3 code"
     , NULL                   AS "organization, level 3"
  FROM PROD_LI.COMMON_DATA.PAYMENT_RUNS
           INNER JOIN PROD_UDM.CONFIGURATION.SOURCE_SYSTEM SS
                      ON PAYMENT_RUNS.DB = RTRIM(SS.SYSTEM_NAME)
           LEFT JOIN COMMON_DATA_STAGING.T001
                     ON SS.SYSTEM_ID = T001.SOURCESYSTEMID AND
                        PAYMENT_RUNS.REGUP_BUKRS = T001.BUKRS
 WHERE RTRIM(REGUH_LIFNR) != ''
 UNION ALL
SELECT DB                     AS "database"
     , ICM_ID                 AS "document id"
     , NULL                   AS "document line number"
     , REGUP_BUKRS            AS "company code id"
     , T001.CONCAT            AS "company code"
     , NULL                   AS "document category code"
     , NULL                   AS "document category"
     , 'REGUH'                AS "document type code"
     , 'Payment Run'          AS "document type"
     , REGUH_XVORL            AS "deletion indicator"
     , REGUH_LAUFD            AS "creation date"
     , NULL                   AS "created by"
     , 'Vendor'               AS "entity type"
     , ICM_VENDOR_ID          AS "entity icm id"
     , REGUH_LIFNR            AS "entity number"
     , "LFA1_LIFNR_Concat"    AS "entity"
     , REGUH_ZALDT            AS "document date"
     , CONCAT_WS(' ',
                 REGUH_ZNME1,
                 REGUH_ZNME2,
                 REGUH_ZNME3,
                 REGUH_ZNME4) AS "document text"
     , REGUH_WAERS            AS "document currency"
     , REGUH_RWBTR_USD        AS "document value"
     , NULL                   AS "document quantity"
     , NULL                   AS "material code"
     , NULL                   AS "material"
     , NULL                   AS "material group code"
     , NULL                   AS "material group"
     , NULL                   AS "organization, level 1 code"
     , NULL                   AS "organization, level 1"
     , NULL                   AS "organization, level 2 code"
     , NULL                   AS "organization, level 2"
     , NULL                   AS "organization, level 3 code"
     , NULL                   AS "organization, level 3"
  FROM PROD_LI.COMMON_DATA.PAYMENT_RUNS
           INNER JOIN PROD_UDM.CONFIGURATION.SOURCE_SYSTEM SS
                      ON PAYMENT_RUNS.DB = RTRIM(SS.SYSTEM_NAME)
           LEFT JOIN COMMON_DATA_STAGING.T001
                     ON SS.SYSTEM_ID = T001.SOURCESYSTEMID AND
                        PAYMENT_RUNS.REGUP_BUKRS = T001.BUKRS
 WHERE RTRIM(REGUH_LIFNR) = ''