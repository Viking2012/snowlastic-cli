SELECT DB                    AS "database"
     , CONCAT_WS('|',
                 DB,
                 PAYR_ZBUKR,
                 PAYR_HBKID,
                 PAYR_HKTID,
                 PAYR_RZAWE,
                 PAYR_CHECT) AS "document id"
     , NULL                  AS "document line number"
     , PAYR_ZBUKR            AS "company code id"
     , NULL                  AS "company code"
     , NULL                  AS "document category code"
     , 'Check'               AS "document category"
     , PAYR_RZAWE            AS "document type code"
     , NULL                  AS "document type"
     , NULL                  AS "deletion indicator"
     , PAYR_LAUFD            AS "creation date"
     , NULL                  AS "created by"
     , 'Vendor'              AS "entity type"
     , CONCAT_WS('|',
                 DB,
                 PAYR_LIFNR) AS "entity icm id"
     , PAYR_LIFNR            AS "entity number"
     , CONCAT_WS(' ',
                 PAYR_ZNME1,
                 PAYR_ZNME2,
                 PAYR_ZNME3,
                 PAYR_ZNME4) AS "entity"
     , PAYR_ZALDT            AS "document date"
     , CONCAT_WS(' ',
                 'Pay To:',
                 PAYR_ZNME1,
                 PAYR_ZNME2,
                 PAYR_ZNME3,
                 PAYR_ZNME4) AS "document text"
     , PAYR_WAERS            AS "document currency"
     , PAYR_RWBTR_USD        AS "document value"
     , NULL                  AS "document quantity"
     , NULL                  AS "material code"
     , NULL                  AS "material"
     , NULL                  AS "material group code"
     , NULL                  AS "material group"
     , NULL                  AS "organization, level 1 code"
     , NULL                  AS "organization, level 1"
     , NULL                  AS "organization, level 2 code"
     , NULL                  AS "organization, level 2"
     , NULL                  AS "organization, level 3 code"
     , NULL                  AS "organization, level 3"
  FROM PROD_LI.COMMON_DATA.CHECKS
 WHERE PAYR_LIFNR != ''
 UNION
SELECT DB                    AS "database"
     , CONCAT_WS('|',
                 DB,
                 PAYR_ZBUKR,
                 PAYR_HBKID,
                 PAYR_HKTID,
                 PAYR_RZAWE,
                 PAYR_CHECT) AS "document id"
     , NULL                  AS "document line number"
     , PAYR_ZBUKR            AS "company code id"
     , NULL                  AS "company code"
     , NULL                  AS "document category code"
     , 'Check'               AS "document category"
     , PAYR_RZAWE            AS "document type code"
     , NULL                  AS "document type"
     , NULL                  AS "deletion indicator"
     , PAYR_LAUFD            AS "creation date"
     , NULL                  AS "created by"
     , 'Customer'            AS "entity type"
     , CONCAT_WS('|',
                 DB,
                 PAYR_KUNNR) AS "entity icm id"
     , PAYR_KUNNR            AS "entity number"
     , CONCAT_WS(' ',
                 PAYR_ZNME1,
                 PAYR_ZNME2,
                 PAYR_ZNME3,
                 PAYR_ZNME4) AS "entity"
     , PAYR_ZALDT            AS "document date"
     , CONCAT_WS(' ',
                 'Pay To:',
                 PAYR_ZNME1,
                 PAYR_ZNME2,
                 PAYR_ZNME3,
                 PAYR_ZNME4) AS "document text"
     , PAYR_WAERS            AS "document currency"
     , PAYR_RWBTR_USD        AS "document value"
     , NULL                  AS "document quantity"
     , NULL                  AS "material code"
     , NULL                  AS "material"
     , NULL                  AS "material group code"
     , NULL                  AS "material group"
     , NULL                  AS "organization, level 1 code"
     , NULL                  AS "organization, level 1"
     , NULL                  AS "organization, level 2 code"
     , NULL                  AS "organization, level 2"
     , NULL                  AS "organization, level 3 code"
     , NULL                  AS "organization, level 3"
  FROM PROD_LI.COMMON_DATA.CHECKS
 WHERE PAYR_LIFNR = ''