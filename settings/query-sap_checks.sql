SELECT CONCAT_WS('|', '{{.schema}}', PAYR.ZBUKR, PAYR.HBKID, PAYR.HKTID, PAYR.RZAWE, PAYR.CHECT) AS ID,
       '{{.schema}}'                                                                             AS "Database",
       PAYR.ZBUKR                                                                                AS "Company Code ID",
       T001.BUTXT                                                                                AS "Company Code",
       'PAYR'                                                                                    AS "Document Category Code",
       'Check'                                                                                   AS "Document Category",
       PAYR.RZAWE                                                                                AS "Document Type Code",
       T042Z.TEXT1                                                                               AS "Document Type",
       PAYR.VOIDR                                                                                AS "Deletion Indicator",
       PAYR.LAUFD                                                                                AS "Creation Date",
       PAYR.PRIUS                                                                                AS "Created By",
       (CASE
            WHEN PAYR.LIFNR != ' ' THEN 'Vendor'
            WHEN PAYR.KUNNR != ' ' THEN 'Customer'
            ELSE 'Unknown' END)                                                                  AS "Entity Type",
       (CASE
            WHEN PAYR.LIFNR != ' ' THEN CONCAT_WS('|', '{{.schema}}', LFA1.LIFNR)
            WHEN PAYR.KUNNR != ' ' THEN CONCAT_WS('|', '{{.schema}}', KNA1.KUNNR)
            ELSE NULL END)                                                                       AS "Entity ICM ID",
       (CASE
            WHEN PAYR.LIFNR != ' ' THEN PAYR.LIFNR
            WHEN PAYR.KUNNR != ' ' THEN PAYR.KUNNR
            ELSE NULL END)                                                                       AS "Entity Number",
       (CASE
            WHEN PAYR.LIFNR != ' ' THEN CONCAT_WS(', ', LFA1.LIFNR, LFA1.NAME1)
            WHEN PAYR.KUNNR != ' ' THEN CONCAT_WS(', ', KNA1.KUNNR, KNA1.NAME1)
            ELSE NULL END)                                                                       AS "Entity",
       TIMESTAMP_NTZ_FROM_PARTS(PAYR.PRIDT, PAYR.PRITI)                                          AS "Document Date",
       TRIM(CONCAT_WS(', ',
                      IFNULL(NULLIF(ZANRE, ''), ''),
                      IFNULL(NULLIF(ZNME1, ''), ''),
                      IFNULL(NULLIF(ZNME2, ''), ''),
                      IFNULL(NULLIF(ZNME3, ''), ''),
                      IFNULL(NULLIF(ZNME4, ''), '')), ', ')                                      AS "Document Text",
       WAERS                                                                                     AS "Document Currency",
       RWBTR                                                                                     AS "Document Value",
       NULL                                                                                      AS "Material Code",
       NULL                                                                                      AS "Material",
       NULL                                                                                      AS "Material Group Code",
       NULL                                                                                      AS "Material Group",
       NULL                                                                                      AS "Organization, Level 1 Code",
       NULL                                                                                      AS "Organization, Level 1",
       NULL                                                                                      AS "Organization, Level 2 Code",
       NULL                                                                                      AS "Organization, Level 2",
       NULL                                                                                      AS "Organization, Level 3 Code",
       NULL                                                                                      AS "Organization, Level 3"
FROM {{.database}}.{{.schema}}.PAYR
         LEFT JOIN {{.database}}.{{.schema}}.T001
                   ON PAYR.MANDT = T001.MANDT
                       AND PAYR.ZBUKR = T001.BUKRS
         LEFT JOIN {{.database}}.{{.schema}}.T042Z
                   ON PAYR.MANDT = T042Z.MANDT AND PAYR.RZAWE = T042Z.ZLSCH AND T001.LAND1 = T042Z.LAND1
         LEFT JOIN {{.database}}.{{.schema}}.LFA1
                   ON PAYR.MANDT = LFA1.MANDT
                       AND PAYR.LIFNR = LFA1.LIFNR
         LEFT JOIN {{.database}}.{{.schema}}.KNA1
                   ON PAYR.MANDT = KNA1.MANDT
                       AND PAYR.KUNNR = KNA1.KUNNR