SELECT 'PROD_UDM_PERSISTENT.' ||
       'PERSISTENT_SFDC.' ||
       'ACCOUNT'                                            AS "record source"
     , 'SalesForce'                                         AS "database"
     , ACCOUNT.ID                                           AS "id"
     , ACCOUNT.ID                                           AS "icm id"
     , ACCOUNT.SITE                                         AS "entity title"
     , ACCOUNT.ID                                           AS "entity number"
     , ACCOUNT.NAME                                         AS "entity name"
     , ACCOUNT.NAME                                         AS "document text"
     , CONCAT_WS('/n',
                 ACCOUNT.PRIMARY_ADDRESS_LOCAL_STREET__C,
                 ACCOUNT.PRIMARY_ADDRESS_LOCAL_CITY__C,
                 ACCOUNT.PRIMARY_ADDRESS_LOCAL_STATE__C,
                 ACCOUNT.PRIMARY_ADDRESS_COUNTRY__C,
                 ACCOUNT.PRIMARY_ADDRESS_LOCAL_ZIP_CODE__C) AS "address"
     , ACCOUNT.PRIMARY_ADDRESS_LOCAL_STREET__C              AS "address street"
     , ACCOUNT.PRIMARY_ADDRESS_LOCAL_CITY__C                AS "address city"
     , ACCOUNT.PRIMARY_ADDRESS_LOCAL_STATE__C               AS "address state or district"
     , ACCOUNT.PRIMARY_ADDRESS_COUNTRY__C                   AS "address country"
     , ACCOUNT.PRIMARY_ADDRESS_LOCAL_ZIP_CODE__C            AS "address zip code"
     , NULL                                                 AS "address po box"
     , ACCOUNT.PHONE                                        AS "phone number"
     , ACCOUNT.FAX                                          AS "fax number"
     , ACCOUNT.WEBSITE                                      AS "email or website"
     , ACCOUNT.CREATEDDATE                                  AS "creation date"
     , USER2.NAME                                           AS "created by"
     , TO_DATE(ACCOUNT.UPDATEDDATETIME)                     AS "document date"
  FROM PROD_UDM_PERSISTENT.PERSISTENT_SFDC.ACCOUNT
           LEFT JOIN PROD_UDM_PERSISTENT.PERSISTENT_SFDC.USER2
                     ON ACCOUNT.CREATEDBYID = USER2.ID