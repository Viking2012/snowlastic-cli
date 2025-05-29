SELECT 'PROD_LI.' ||
       'COMMON_DATA.' ||
       'CHECKS'  AS "record source"
     , DB                    AS "database"
     , CONCAT_WS('|',
                 DB,
                 PAYR_ZBUKR,
                 PAYR_HBKID,
                 PAYR_HKTID,
                 PAYR_RZAWE,
                 PAYR_CHECT) AS "icm id"
     , "icm id"              as "id"
     , NULL                  AS "document line number"
     , CONCAT_WS('|',
                 DB,
                 PAYR_LIFNR) AS "entity icm id"
     , 'Vendor'              AS "entity title"
     , PAYR_LIFNR            AS "entity number"
     , CONCAT_WS(' ',
                 PAYR_ZNME1,
                 PAYR_ZNME2,
                 PAYR_ZNME3,
                 PAYR_ZNME4) AS "entity name"
     , NULL                  AS "document category code"
     , 'Check'               AS "document category"
     , PAYR_RZAWE            AS "document type code"
     , NULL                  AS "document type"
     , PAYR_LAUFD            AS "creation date"
     , NULL                  AS "created by"
     , PAYR_ZALDT            AS "document date"
     , PAYR_RWBTR_USD        AS "document value"
     , PAYR_WAERS            AS "document currency"
     , NULL                  AS "document quantity"
     , CONCAT_WS(' ',
                 'Pay To:',
                 PAYR_ZNME1,
                 PAYR_ZNME2,
                 PAYR_ZNME3,
                 PAYR_ZNME4) AS "document text"
     , NULL                  AS "material group code"
     , NULL                  AS "material group"
     , NULL                  AS "material code"
     , NULL                  AS "material"
     , PAYR_ZBUKR            AS "company code id"
     , NULL                  AS "company code"
     , NULL                  AS "organization, level 1 code"
     , NULL                  AS "organization, level 1"
     , NULL                  AS "organization, level 2 code"
     , NULL                  AS "organization, level 2"
     , NULL                  AS "organization, level 3 code"
     , NULL                  AS "organization, level 3"
     , NULL                  AS "deletion indicator"
  FROM PROD_LI.COMMON_DATA.CHECKS
 WHERE PAYR_LIFNR != ''
 UNION ALL
SELECT 'COMMON_DATA.CHECKS'  AS "record source"
     , DB                    AS "database"
     , CONCAT_WS('|',
                 DB,
                 PAYR_ZBUKR,
                 PAYR_HBKID,
                 PAYR_HKTID,
                 PAYR_RZAWE,
                 PAYR_CHECT) AS "icm id"
     , "icm id"              as "id"
     , NULL                  AS "document line number"
     , CONCAT_WS('|',
                 DB,
                 PAYR_KUNNR) AS "entity icm id"
     , 'Customer'            AS "entity title"
     , PAYR_KUNNR            AS "entity number"
     , CONCAT_WS(' ',
                 PAYR_ZNME1,
                 PAYR_ZNME2,
                 PAYR_ZNME3,
                 PAYR_ZNME4) AS "entity name"
     , NULL                  AS "document category code"
     , 'Check'               AS "document category"
     , PAYR_RZAWE            AS "document type code"
     , NULL                  AS "document type"
     , PAYR_LAUFD            AS "creation date"
     , NULL                  AS "created by"
     , PAYR_ZALDT            AS "document date"
     , PAYR_RWBTR_USD        AS "document value"
     , PAYR_WAERS            AS "document currency"
     , NULL                  AS "document quantity"
     , CONCAT_WS(' ',
                 'Pay To:',
                 PAYR_ZNME1,
                 PAYR_ZNME2,
                 PAYR_ZNME3,
                 PAYR_ZNME4) AS "document text"
     , NULL                  AS "material group code"
     , NULL                  AS "material group"
     , NULL                  AS "material code"
     , NULL                  AS "material"
     , PAYR_ZBUKR            AS "company code id"
     , NULL                  AS "company code"
     , NULL                  AS "organization, level 1 code"
     , NULL                  AS "organization, level 1"
     , NULL                  AS "organization, level 2 code"
     , NULL                  AS "organization, level 2"
     , NULL                  AS "organization, level 3 code"
     , NULL                  AS "organization, level 3"
     , NULL                  AS "deletion indicator"
  FROM PROD_LI.COMMON_DATA.CHECKS
 WHERE PAYR_LIFNR = ''