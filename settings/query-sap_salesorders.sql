SELECT CONCAT_WS('|', '{{.schema}}', VBAP.VBELN, VBAP.POSNR) AS ID,
       '{{.schema}}'                                         AS "Database",
       VBAK.BUKRS_VF                                         AS "Company Code ID",
       T001.BUTXT                                            AS "Company Code",
       VBAK.VBTYP                                            AS "Document Category Code",
       (CASE
            WHEN VBAK.VBTYP = 'A' THEN 'Inquiry'
            WHEN VBAK.VBTYP = 'B' THEN 'Quotation'
            WHEN VBAK.VBTYP = 'C' THEN 'Order'
            WHEN VBAK.VBTYP = 'D' THEN 'Item proposal'
            WHEN VBAK.VBTYP = 'E' THEN 'Scheduling agreement'
            WHEN VBAK.VBTYP = 'F' THEN 'Scheduling agreement with external service agent'
            WHEN VBAK.VBTYP = 'G' THEN 'Contract'
            WHEN VBAK.VBTYP = 'H' THEN 'Returns'
            WHEN VBAK.VBTYP = 'I' THEN 'Order w/o charge'
            WHEN VBAK.VBTYP = 'J' THEN 'Delivery'
            WHEN VBAK.VBTYP = 'K' THEN 'Credit memo request'
            WHEN VBAK.VBTYP = 'L' THEN 'Debit memo request'
            WHEN VBAK.VBTYP = 'M' THEN 'Invoice'
            WHEN VBAK.VBTYP = 'N' THEN 'Invoice cancellation'
            WHEN VBAK.VBTYP = 'O' THEN 'Credit memo'
            WHEN VBAK.VBTYP = 'P' THEN 'Debit memo'
            WHEN VBAK.VBTYP = 'Q' THEN 'WMS transfer order'
            WHEN VBAK.VBTYP = 'R' THEN 'Goods movement'
            WHEN VBAK.VBTYP = 'S' THEN 'Credit memo cancellation'
            WHEN VBAK.VBTYP = 'T' THEN 'Returns delivery for order'
            WHEN VBAK.VBTYP = 'U' THEN 'Pro forma invoice'
            WHEN VBAK.VBTYP = 'V' THEN 'Purchase Order'
            WHEN VBAK.VBTYP = 'W' THEN 'Independent reqts plan'
            WHEN VBAK.VBTYP = 'X' THEN 'Handling unit'
            WHEN VBAK.VBTYP = '0' THEN 'Master contract'
            WHEN VBAK.VBTYP = '1' THEN 'Sales activities (CAS)'
            WHEN VBAK.VBTYP = '2' THEN 'External transaction'
            WHEN VBAK.VBTYP = '3' THEN 'Invoice list'
            WHEN VBAK.VBTYP = '4' THEN 'Credit memo list'
            WHEN VBAK.VBTYP = '5' THEN 'Intercompany invoice'
            WHEN VBAK.VBTYP = '6' THEN 'Intercompany credit memo'
            WHEN VBAK.VBTYP = '7' THEN 'Delivery/shipping notification'
            WHEN VBAK.VBTYP = '8' THEN 'Shipment'
            WHEN VBAK.VBTYP = 'a' THEN 'Shipment costs'
            WHEN VBAK.VBTYP = 'b' THEN 'CRM Opportunity'
            WHEN VBAK.VBTYP = 'c' THEN 'Unverified delivery'
            WHEN VBAK.VBTYP = 'd' THEN 'Trading Contract'
            WHEN VBAK.VBTYP = 'e' THEN 'Allocation table'
            WHEN VBAK.VBTYP = 'f' THEN 'Additional Billing Documents'
            WHEN VBAK.VBTYP = 'g' THEN 'Rough Goods Receipt (only IS-Retail)'
            WHEN VBAK.VBTYP = 'h' THEN 'Cancel Goods Issue'
            WHEN VBAK.VBTYP = 'i' THEN 'Goods receipt'
            WHEN VBAK.VBTYP = 'j' THEN 'JIT call'
            WHEN VBAK.VBTYP = 'n' THEN 'Reserved'
            WHEN VBAK.VBTYP = 'o' THEN 'Reserved'
            WHEN VBAK.VBTYP = 'p' THEN 'Goods Movement (Documentation)'
            WHEN VBAK.VBTYP = 'q' THEN 'Reserved'
            WHEN VBAK.VBTYP = 'r' THEN 'TD Transport (only IS-Oil)'
            WHEN VBAK.VBTYP = 's' THEN 'Load Confirmation, Reposting (Only IS-Oil)'
            WHEN VBAK.VBTYP = 't' THEN 'Gain / Loss (Only IS-Oil)'
            WHEN VBAK.VBTYP = 'u' THEN 'Reentry into Storage (Only IS-Oil)'
            WHEN VBAK.VBTYP = 'v' THEN 'Data Collation (only IS-Oil)'
            WHEN VBAK.VBTYP = 'w' THEN 'Reservation (Only IS-Oil)'
            WHEN VBAK.VBTYP = 'x' THEN 'Load Confirmation, Goods Receipt (Only IS-Oil)'
            WHEN VBAK.VBTYP = '$' THEN '(AFS)'
            WHEN VBAK.VBTYP = ':' THEN 'Service Order'
            WHEN VBAK.VBTYP = '.' THEN 'Service Notification'
            WHEN VBAK.VBTYP = '&' THEN 'Warehouse Document'
            WHEN VBAK.VBTYP = '*' THEN 'Pick Order'
            WHEN VBAK.VBTYP = ',' THEN 'Shipment Document'
            WHEN VBAK.VBTYP = '^' THEN 'Reserved'
            WHEN VBAK.VBTYP = '|' THEN 'Reserved'
            WHEN VBAK.VBTYP = 'k' THEN 'Agency Document'
            ELSE NULL END)                                   AS "Document Category",
       VBAK.AUART                                            AS "Document Type Code",
       TVAKT.BEZEI                                           AS "Document Type",
       VBAP.ABGRU                                            AS "Deletion Indicator",
       VBAP.ERDAT                                            AS "Creation Date",
       VBAP.ERNAM                                            AS "Created By",
       'Customer'                                            AS "Entity Type",
       CONCAT_WS('|', '{{.schema}}', VBAK.KUNNR)             AS "Entity ICM ID",
       VBAK.KUNNR                                            AS "Entity Number",
       CONCAT_WS(', ', KNA1.KUNNR, KNA1.NAME1)               AS "Entity",
       VBAK.AUDAT                                            AS "Document Date",
       VBAP.ARKTX                                            AS "Document Text",
       VBAP.WAERK                                            AS "Document Currency",
       VBAP.NETWR                                            AS "Document Value",
       VBAP.MATNR                                            AS "Material Code",
       MAKT.MAKTX                                            AS "Material",
       VBAP.MATKL                                            AS "Material Group Code",
       T023T.WGBEZ60                                         AS "Material Group",
       PCM.ORG1                                              AS "Organization, Level 1 Code",
       PCM.ORG1_CONCAT                                       AS "Organization, Level 1",
       PCM.ORG2                                              AS "Organization, Level 2 Code",
       PCM.ORG2_CONCAT                                       AS "Organization, Level 2",
       PCM.ORG3                                              AS "Organization, Level 3 Code",
       PCM.ORG3_CONCAT                                       AS "Organization, Level 3"
FROM {{.database}}.{{.schema}}.VBAK
         INNER JOIN {{.database}}.{{.schema}}.VBAP
            ON VBAK.MANDT = VBAP.MANDT AND VBAK.VBELN = VBAP.VBELN
         LEFT JOIN {{.database}}.{{.schema}}.T001
             ON VBAK.MANDT = T001.MANDT AND VBAK.BUKRS_VF = T001.BUKRS
         LEFT JOIN {{.database}}.{{.schema}}.TVAKT
             ON VBAK.MANDT = TVAKT.MANDT AND VBAK.AUART = TVAKT.AUART
         LEFT JOIN {{.database}}.{{.schema}}.KNA1
             ON VBAK.MANDT = KNA1.MANDT AND VBAK.KUNNR = KNA1.KUNNR
         LEFT JOIN DEV_FPA_LI.BA_MASTERDATA.PROFITCENTERMAPPING PCM
             ON VBAP.PRCTR = PCM.PRCTR AND VBAK.BUKRS_VF = PCM.BUKRS
         LEFT JOIN {{.database}}.{{.schema}}.MAKT
             ON VBAP.MANDT = MAKT.MANDT AND VBAP.MATNR = MAKT.MATNR
         LEFT JOIN {{.database}}.{{.schema}}.T023T
             ON VBAP.MANDT = T023T.MANDT AND VBAP.MATKL = T023T.MATKL
WHERE TVAKT.SPRAS = 'E'
  AND  MAKT.SPRAS = 'E'
  AND T023T.SPRAS = 'E'