SELECT 'PROD_UDM_PERSISTENT.' ||
       'PERSISTENT_SFDC.' ||
       'USER2'                                 AS "record source"
     , 'SalesForce'                            AS "database"
     , USER2.ID                                AS "id"
     , USER2.ID                                AS "icm id"
     , USER2.TITLE                             AS "entity title"
     , USER2.ID                                AS "entity number"
     , USER2.NAME                              AS "entity name"
     , USER2.USERNAME                          AS "document text"
     , CONCAT_WS(' ',
                 IFNULL(USER2.STREET, ''),
                 IFNULL(USER2.CITY, ''),
                 IFNULL(USER2.STATE, ''),
                 IFNULL(USER2.COUNTRY, ''),
                 IFNULL(USER2.POSTALCODE, '')) AS "address"
     , USER2.STREET                            AS "address street"
     , USER2.CITY                              AS "address city"
     , USER2.STATE                             AS "address state or district"
     , USER2.COUNTRY                           AS "address country"
     , USER2.POSTALCODE                        AS "address zip code"
     , NULL                                    AS "address po box"
     , USER2.PHONE                             AS "phone number"
     , USER2.FAX                               AS "fax number"
     , USER2.EMAIL                             AS "email or website"
     , TO_DATE(USER2.CREATEDDATE)                       AS "creation date"
     , CREATE_USER.NAME                        AS "created by"
     , TO_DATE(USER2.UPDATEDDATETIME)          AS "document date"
  FROM PROD_UDM_PERSISTENT.PERSISTENT_SFDC.USER2
           LEFT JOIN PROD_UDM_PERSISTENT.PERSISTENT_SFDC.USER2 CREATE_USER
                     ON USER2.CREATEDBYID = CREATE_USER.ID