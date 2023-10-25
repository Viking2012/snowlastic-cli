SELECT CONCAT_WS('|', '{{.schema}}', RBCO.MANDT, RBCO.BELNR, RBCO.GJAHR, RBCO.BUZEI, RBCO.COBL_NR) AS ID,
       '{{.schema}}'                                                                               AS "Database",
       RBCO.BUKRS                                                                                  AS "Company Code ID",
       T001.BUTXT                                                                                  AS "Company Code",
       RBKP.BLART                                                                                  AS "Document Category Code",
       T003t.LTEXT                                                                                 AS "Document Category",
       RBKP.VGART                                                                                  AS "Document Type Code",
       (CASE
            WHEN RBKP.VGART = 'RP' THEN 'Invoice receipt'
            WHEN RBKP.VGART = 'KP' THEN 'Account maintenance'
            WHEN RBKP.VGART = 'KS' THEN 'Account maintenance reversal'
            WHEN RBKP.VGART = 'PR' THEN 'Price change'
            WHEN RBKP.VGART = 'BL' THEN 'Material debit'
            WHEN RBKP.VGART = 'PF' THEN 'Std cost estim. release'
            WHEN RBKP.VGART = 'RD' THEN 'Logistics invoice'
            WHEN RBKP.VGART = 'ML' THEN 'Material ledger settlement'
            WHEN RBKP.VGART = 'MI' THEN 'Material ledger initialization'
            WHEN RBKP.VGART = 'RS' THEN 'Logistics invoice, cancel'
           END)                                                                                    AS "Document Type",
       NULL                                                                                        AS "Deletion Indicator",
       TIMESTAMP_NTZ_FROM_PARTS(RBKP.CPUDT, RBKP.CPUTM)                                            AS "Creation Date",
       RBKP.USNAM                                                                                  AS "Created By",
       'Vendor'                                                                                    AS "Entity Type",
       CONCAT_WS('|', '{{.schema}}', NULLIF(RBKP.LIFNR, ' '))                                      AS "Entity ICM ID",
       NULLIF(RBKP.LIFNR, ' ')                                                                     AS "Entity Number",
       CONCAT_WS(', ', RBKP.LIFNR, LFA1.NAME1)                                                     AS "Entity",
       RBKP.BLDAT                                                                                  AS "Document Date",
       TRIM(CONCAT_WS('\n',
                      IFNULL(CONCAT('HEADER: ', NULLIF(RBKP.BKTXT, ' '), '\n'), ''),
                      IFNULL(CONCAT('ITEM: ', NULLIF(RBKP.SGTXT, ' '), '\n'), '')
                ), '\n')                                                                           AS "Document Text",
       NULL                                                                                        AS "Material Code",
       NULL                                                                                        AS "Material",
       NULL                                                                                        AS "Material Group Code",
       NULL                                                                                        AS "Material Group",
       PCM.ORG1                                                                                    AS "Organization, Level 1 Code",
       PCM.ORG1_CONCAT                                                                             AS "Organization, Level 1",
       PCM.ORG2                                                                                    AS "Organization, Level 2 Code",
       PCM.ORG2_CONCAT                                                                             AS "Organization, Level 2",
       PCM.ORG3                                                                                    AS "Organization, Level 3 Code",
       PCM.ORG3_CONCAT                                                                             AS "Organization, Level 3"
FROM {{.database}}.{{.schema}}.RBKP
         INNER JOIN {{.database}}.{{.schema}}.RBCO
                    ON RBKP.MANDT = RBCO.MANDT
                        AND RBKP.BELNR = RBCO.BELNR
                        AND RBKP.GJAHR = RBCO.GJAHR
         LEFT JOIN {{.database}}.{{.schema}}.T001
                   ON RBCO.MANDT = T001.MANDT
                       AND RBCO.BUKRS = T001.BUKRS
         LEFT JOIN {{.database}}.{{.schema}}.T003T
                   ON RBKP.MANDT = T003T.MANDT
                       AND RBKP.BLART = T003T.BLART
         LEFT JOIN {{.database}}.{{.schema}}.LFA1
                   ON RBKP.MANDT = LFA1.MANDT
                       AND RBKP.LIFNR = LFA1.LIFNR
         LEFT JOIN {{.database}}.BA_MASTERDATA.PROFITCENTERMAPPING PCM
                   ON RBCO.BUKRS = PCM.BUKRS
                       AND RBCO.PRCTR = PCM.PRCTR
WHERE T003T.SPRAS = 'E'