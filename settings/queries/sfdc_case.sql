SELECT 'PROD_UDM_PERSISTENT.' ||
       'PERSISTENT_SFDC.' ||
       'CASE2'                    AS "record source"
     , 'SalesForce'               AS "database"
     , CASE2.ID                   AS "id"
     , CASE2.ID                   AS "icm id"
     , 'Case'                     AS "document type"
     , TO_DATE(CASE2.CREATEDDATE) AS "document date"
     , TO_DATE(CASE2.CREATEDDATE) AS "creation date"
     , CREATOR.NAME               AS "created by"
     , CONCAT_WS('\n',
                 'CASE NUMBER: ' || IFNULL(CASE2.CASENUMBER, ''),
                 'FROM: ' || IFNULL(CASE2.CONTACT_NAME__C, '') || ' <' || IFNULL(CASE2.CONTACT_EMAIL__C, '') || '>',
                 'ACCOUNT: ' || IFNULL(CASE2.ACCOUNT_NAME_FORMULA__C, ''),
                 'TO: ' || IFNULL(OWNERS.NAME, '') || '<' || IFNULL(OWNERS.EMAIL, '') || '>',
                 'DIVISION: ' || IFNULL(CASE2.DIVISION__C, ''),
                 'BUSINESS UNIT: ' || IFNULL(CASE2.BUSINESS_UNIT__C, ''),
                 'PG: ' || IFNULL(CASE2.PRODUCT_GROUP__C,''),
                 'CATEGORY: ' || RTRIM(IFNULL(CASE2.CATEGORY__C, '') || ';' ||
                                       IFNULL(CASE2.CATEGORY_2__C, '') || ';' ||
                                       IFNULL(CASE2.CATEGORY_3__C, ''), ';'),
                 'ORIGIN: ' || IFNULL(CASE2.ORIGIN, ''),
                 'STATUS: ' || IFNULL(CASE2.STATUS, ''),
                 'SUBJECT: ' || IFNULL(CASE2.SUBJECT, ''),
                 'DESCRIPTION: ' || IFNULL(CASE2.DESCRIPTION, '')
       )                          AS "document text"
  FROM PROD_UDM_PERSISTENT.PERSISTENT_SFDC.CASE2
           LEFT JOIN PROD_UDM_PERSISTENT.PERSISTENT_SFDC.USER2 OWNERS
                     ON CASE2.OWNERID = OWNERS.ID
           LEFT JOIN PROD_UDM_PERSISTENT.PERSISTENT_SFDC.USER2 CREATOR
                     ON CASE2.CREATEDBYID = CREATOR.ID