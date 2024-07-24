SELECT 'COMMON_DATA.CENTRAL_CUSTOMER_MASTER' AS "record source"
     , DB                                    AS "database"
     , ICM_ID                                AS "icm id"
     , 'Customer'                            AS "entity title"
     , KNA1_KUNNR                            AS "entity number"
     , CONCAT_WS(' ',
                 KNA1_NAME1,
                 KNA1_NAME2,
                 KNA1_NAME3,
                 KNA1_NAME4)                 AS "entity name"
     , CONCAT_WS(' ',
                 KNA1_NAME1,
                 KNA1_NAME2,
                 KNA1_NAME3,
                 KNA1_NAME4)                 AS "document text"
     , CONCAT_WS('; ',
                 KNA1_STRAS,
                 KNA1_ORT01,
                 KNA1_REGIO,
                 KNA1_LAND1,
                 KNA1_PSTLZ)                 AS "address"
     , KNA1_STRAS                            AS "address street"
     , KNA1_ORT01                            AS "address city"
     , KNA1_REGIO                            AS "address state or district"
     , "KNA1_LAND1_Concat"                   AS "address country"
     , KNA1_PSTLZ                            AS "address zip code"
     , NULL                                  AS "address po box"
     , KNA1_TELF1                            AS "phone number"
     , KNA1_TELFX                            AS "fax number"
     , NULL                                  AS "email"
     , KNA1_ERDAT                            AS "creation date"
     , KNA1_ERNAM                            AS "created by"
     , CURRENT_DATE()                        AS "document date"
  FROM COMMON_DATA.CENTRAL_CUSTOMER_MASTER

