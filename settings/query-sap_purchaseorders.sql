SELECT CONCAT_WS('|', '{{.schema}}', EKPO.EBELN, EKPO.EBELP) AS ID,
       '{{.schema}}'                                         AS "Database",
       EKPO.BUKRS                                            AS "Company Code ID",
       T001.BUTXT                                            AS "Company Code",
       EKKO.BSTYP                                            AS "Document Category Code",
       (CASE
            WHEN EKKO.BSTYP = 'A' THEN 'Request for Quotation'
            WHEN EKKO.BSTYP = 'F' THEN 'Purchase Order'
            WHEN EKKO.BSTYP = 'K' THEN 'Contract'
            WHEN EKKO.BSTYP = 'L' THEN 'Scheduling Agreement'
            ELSE 'Unknown' END)                              AS "Document Category",
       EKKO.BSART                                            AS "Document Type Code",
       T161T.BATXT                                           AS "Document Type",
       EKPO.LOEKZ                                            AS "Deletion Indicator",
       EKKO.AEDAT                                            AS "Creation Date",
       EKKO.ERNAM                                            AS "Created By",
       'Vendor'                                              AS "Entity Type",
       CONCAT_WS('|', '{{.schema}}', EKKO.LIFNR)             AS "Entity ICM ID",
       EKKO.LIFNR                                            AS "Entity Number",
       CONCAT_WS(', ', EKKO.LIFNR, LFA1.NAME1)               AS "Entity",
       EKKO.BEDAT                                            AS "Document Date",
       EKPO.TXZ01                                            AS "Document Text",
       EKPO.MATNR                                            AS "Material Code",
       MAKT.MAKTX                                            AS "Material",
       EKPO.MATKL                                            AS "Material Group Code",
       T023T.WGBEZ60                                         AS "Material Group",
       PCM.ORG1                                              AS "Organization, Level 1 Code",
       PCM.ORG1_CONCAT                                       AS "Organization, Level 1",
       PCM.ORG2                                              AS "Organization, Level 2 Code",
       PCM.ORG2_CONCAT                                       AS "Organization, Level 2",
       PCM.ORG3                                              AS "Organization, Level 3 Code",
       PCM.ORG3_CONCAT                                       AS "Organization, Level 3"
FROM {{.database}}.{{.schema}}.EKKO
         INNER JOIN {{.database}}.{{.schema}}.EKPO
                    ON EKKO.MANDT = EKPO.MANDT
                        AND EKKO.EBELN = EKPO.EBELN
         LEFT JOIN {{.database}}.{{.schema}}.T001
                   ON EKPO.MANDT = T001.MANDT
                       AND EKPO.BUKRS = T001.BUKRS
         LEFT JOIN {{.database}}.{{.schema}}.T161T
                   ON EKKO.MANDT = T161T.MANDT
                       AND EKKO.BSTYP = T161T.BSTYP
                       AND EKKO.BSART = T161T.BSART
         LEFT JOIN {{.database}}.{{.schema}}.LFA1
                   ON EKKO.MANDT = LFA1.MANDT
                       AND EKKO.LIFNR = LFA1.LIFNR
         LEFT JOIN {{.database}}.{{.schema}}.MAKT
                   ON EKPO.MANDT = MAKT.MANDT AND EKPO.MATNR = MAKT.MATNR
         LEFT JOIN {{.database}}.{{.schema}}.T023T
                   ON EKPO.MANDT = T023T.MANDT AND EKPO.MATKL = T023T.MATKL
         LEFT JOIN {{.database}}.BA_MASTERDATA.PROFITCENTERMAPPING PCM
                   ON EKPO.BUKRS = PCM.BUKRS
                       AND EKPO.KO_PRCTR = PCM.PRCTR
WHERE T161T.SPRAS = 'E'
  AND MAKT.SPRAS = 'E'
  AND T023T.SPRAS = 'E'