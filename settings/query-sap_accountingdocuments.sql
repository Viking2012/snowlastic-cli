SELECT CONCAT_WS('|', '{{.schema}}', BSEG.BUKRS, BSEG.BELNR, BSEG.GJAHR, BSEG.BUZEI) AS ID,
       '{{.schema}}'                                                                 AS "Database",
       BSEG.BUKRS                                                                    AS "Company Code ID",
       T001.BUTXT                                                                    AS "Company Code",
       BKPF.DOCCAT                                                                   AS "Document Category Code",
       (CASE
            WHEN BKPF.DOCCAT = 'INVREC' THEN 'Invoice Receipt'
            WHEN BKPF.DOCCAT = 'INVRED' THEN 'Invoice Reduction'
            WHEN BKPF.DOCCAT = 'REVAL' THEN 'Revaluation'
            WHEN BKPF.DOCCAT = 'PURACC' THEN 'Purchase Account Processing'
            WHEN BKPF.DOCCAT = 'RET' THEN 'Retention'
            WHEN BKPF.DOCCAT = 'DPC_MM' THEN 'Down Payment Clearing'
            WHEN BKPF.DOCCAT = 'PEVACR' THEN 'Period-End Valuation - Posting Accrual Document'
            WHEN BKPF.DOCCAT = 'PEVRST' THEN 'Period-End Valuation - Resetting Accrual Document'
            ELSE NULL END)                                                           AS "Document Category",
       BKPF.BLART                                                                    AS "Document Type Code",
       T003T.LTEXT                                                                   AS "Document Type",
       NULL                                                                          AS "Deletion Indicator",
       TIMESTAMP_NTZ_FROM_PARTS(BKPF.CPUDT, BKPF.CPUTM)                              AS "Creation Date",
       BKPF.USNAM                                                                    AS "Created By",
       (CASE
            WHEN BSEG.LIFNR != '' THEN 'Vendor'
            WHEN BSEG.KUNNR != '' THEN 'Customer'
            ELSE 'Unknown' END)                                                      AS "Entity Type",
       (CASE
            WHEN BSEG.LIFNR != '' THEN CONCAT_WS('|', '{{.schema}}', BSEG.LIFNR)
            WHEN BSEG.KUNNR != '' THEN CONCAT_WS('|', '{{.schema}}', BSEG.KUNNR)
            ELSE 'Unknown' END)                                                      AS "Entity ICM ID",
       (CASE
            WHEN BSEG.LIFNR != '' THEN BSEG.LIFNR
            WHEN BSEG.KUNNR != '' THEN BSEG.KUNNR
            ELSE NULL END)                                                           AS "Entity Number",
       (CASE
            WHEN BSEG.LIFNR != '' THEN CONCAT_WS(', ', BSEG.LIFNR, LFA1.NAME1)
            WHEN BSEG.KUNNR != '' THEN CONCAT_WS(', ', BSEG.LIFNR, KNA1.NAME1)
            ELSE NULL END)                                                           AS "Entity",
       BKPF.BUDAT                                                                    AS "Document Date",
       TRIM(CONCAT_WS('\n',
                      IFNULL(CONCAT('HEADER: ', NULLIF(BKPF.BKTXT, ''), '\n'), ''),
                      IFNULL(CONCAT('ITEM: ', NULLIF(BSEG.SGTXT, ''), '\n'), '')
                ), '\n')                                                             AS "Document Text",
       BKPF.WAERS                                                                    AS "Document Currency",
       BSEG.WRBTR                                                                    AS "Document Value",
       NULL                                                                          AS "Material Code",
       NULL                                                                          AS "Material",
       NULL                                                                          AS "Material Group Code",
       NULL                                                                          AS "Material Group",
       PCM.ORG1                                                                      AS "Organization, Level 1 Code",
       PCM.ORG1_CONCAT                                                               AS "Organization, Level 1",
       PCM.ORG2                                                                      AS "Organization, Level 2 Code",
       PCM.ORG2_CONCAT                                                               AS "Organization, Level 2",
       PCM.ORG3                                                                      AS "Organization, Level 3 Code",
       PCM.ORG3_CONCAT                                                               AS "Organization, Level 3"
FROM {{.database}}.{{.schema}}.BKPF
         INNER JOIN {{.database}}.{{.schema}}.BSEG
                    ON BKPF.MANDT = BSEG.MANDT
                        AND BKPF.BUKRS = BSEG.BUKRS
                        AND BKPF.BELNR = BSEG.BELNR
                        AND BKPF.GJAHR = BSEG.GJAHR
         LEFT JOIN {{.database}}.{{.schema}}.T001
                   ON BSEG.MANDT = T001.MANDT
                       AND BSEG.BUKRS = T001.BUKRS
         LEFT JOIN {{.database}}.{{.schema}}.T003T
                   ON BSEG.MANDT = T003T.MANDT
                       AND BKPF.BLART = T003T.BLART
         LEFT JOIN {{.database}}.{{.schema}}.LFA1
                   ON BSEG.MANDT = LFA1.MANDT
                       AND BSEG.LIFNR = LFA1.LIFNR
         LEFT JOIN {{.database}}.{{.schema}}.KNA1
                   ON BSEG.MANDT = KNA1.MANDT
                       AND BSEG.KUNNR = KNA1.KUNNR
         LEFT JOIN (SELECT DISTINCT
                        BUKRS,PRCTR,
                        ORG1,ORG1_CONCAT,
                        ORG2,ORG2_CONCAT,
                        ORG3,ORG3_CONCAT
                    FROM {{.database}}.BA_MASTERDATA.PROFITCENTERMAPPING) PCM
                   ON BKPF.BUKRS = PCM.BUKRS
                       AND BSEG.PRCTR = PCM.PRCTR
WHERE T003T.SPRAS = 'E'
  AND BSEG.LIFNR != ''
  AND BSEG.KUNNR != ''