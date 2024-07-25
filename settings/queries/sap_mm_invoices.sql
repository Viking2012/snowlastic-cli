SELECT 'COMMON_DATA.MM_INVOICES'       AS "record source"
     , MM_INVOICES.DB                  AS "database"
     , ICM_ID || '|' || RSEG_BUZEI     AS "id"
     , ICM_ID                          AS "icm id"
     , RSEG_BUZEI                      AS "document line number"
     , MM_INVOICES.ICM_VENDOR_ID       AS "entity icm id"
     , 'Vendor'                        AS "entity title"
     , RBKP_LIFNR                      AS "entity number"
     , RTRIM(CONCAT_WS(' ',
                       VM.LFA1_NAME1,
                       VM.LFA1_NAME2,
                       VM.LFA1_NAME3,
                       VM.LFA1_NAME4)) AS "entity name"
     , RBKP_BLART                      AS "document category code"
     , RTRIM(T003T.LTEXT)              AS "document category"
     , RBKP_VGART                      AS "document type code"
     , (CASE
            WHEN RBKP_VGART = 'RP' THEN 'Invoice receipt'
            WHEN RBKP_VGART = 'KP' THEN 'Account maintenance'
            WHEN RBKP_VGART = 'KS' THEN 'Account maintenance reversal'
            WHEN RBKP_VGART = 'PR' THEN 'Price change'
            WHEN RBKP_VGART = 'BL' THEN 'Material debit'
            WHEN RBKP_VGART = 'PF' THEN 'Std cost estim. release'
            WHEN RBKP_VGART = 'RD' THEN 'Logistics invoice'
            WHEN RBKP_VGART = 'ML' THEN 'Material ledger settlement'
            WHEN RBKP_VGART = 'MI' THEN 'Material ledger initialization'
            WHEN RBKP_VGART = 'RS' THEN 'Logistics invoice, cancel'
    END)                               AS "document type"
     , RBKP_CPUDT                      AS "creation date"
     , RBKP_USNAM                      AS "created by"
     , RBKP_BLDAT                      AS "document date"
     , "AmountUSD"                     AS "document value"
     , RBKP_WAERS                      AS "document currency"
     , RSEG_RBCO_BPMNG                 AS "document quantity"
     , NULL                            AS "document text"
     , EKPO_MATKL                      AS "material group code"
     , "EKPO_MATKL_Concat"             AS "material group"
     , RSEG_MATNR                      AS "material code"
     , "RSEG_MATNR_Concat"             AS "material"
     , RBKP_BUKRS                      AS "company code id"
     , "RSEG_BUKRS_Concat"             AS "company code"
     , "ProfitCenter_ORG1"             AS "organization, level 1 code"
     , "ProfitCenter_ORG1_Concat"      AS "organization, level 1"
     , "ProfitCenter_ORG2"             AS "organization, level 2 code"
     , "ProfitCenter_ORG2_Concat"      AS "organization, level 2"
     , "ProfitCenter_ORG3"             AS "organization, level 3 code"
     , "ProfitCenter_ORG3_Concat"      AS "organization, level 3"
     , NULL                            AS "deletion indicator"
  FROM COMMON_DATA.MM_INVOICES
           LEFT JOIN PROD_UDM_PERSISTENT.PERSISTENT_SAP.T003T
                     ON MM_INVOICES.SOURCESYSTEMID = T003T.SOURCESYSTEMID
                         AND MM_INVOICES.RBKP_BLART = T003T.BLART
           LEFT JOIN COMMON_DATA.CENTRAL_VENDOR_MASTER VM
                     ON MM_INVOICES.ICM_VENDOR_ID = VM.ICM_VENDOR_ID