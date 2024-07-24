SELECT 'COMMON_DATA.CENTRAL_CUSTOMER_MASTER' AS "record source"
     , DB                                    AS "database"
     , ICM_ID                                AS "icm id"
     , 'Customer'                            AS "entity title"
     , LFA1_LIFNR                            AS "entity number"
     , RTRIM(CONCAT_WS(' ',
                       LFA1_NAME1,
                       LFA1_NAME2,
                       LFA1_NAME3,
                       LFA1_NAME4))          AS "entity name"
     , RTRIM(CONCAT_WS(' ',
                       LFA1_NAME1,
                       LFA1_NAME2,
                       LFA1_NAME3,
                       LFA1_NAME4))          AS "document text"
     , RTRIM(CONCAT_WS('; ',
                       LFA1_STRAS,
                       LFA1_ORT01,
                       LFA1_REGIO,
                       LFA1_LAND1,
                       LFA1_PSTLZ), '; ')    AS "address"
     , LFA1_STRAS                            AS "address street"
     , LFA1_ORT01                            AS "address city"
     , LFA1_REGIO                            AS "address state or district"
     , "LFA1_LAND1_Concat"                   AS "address country"
     , LFA1_PSTLZ                            AS "address zip code"
     , NULL                                  AS "address po box"
     , LFA1_TELF1                            AS "phone number"
     , LFA1_TELFX                            AS "fax number"
     , NULL                                  AS "email"
     , LFA1_ERDAT                            AS "creation date"
     , LFA1_ERNAM                            AS "created by"
     , CURRENT_DATE()                        AS "document date"
  FROM PROD_LI.COMMON_DATA.CENTRAL_VENDOR_MASTER