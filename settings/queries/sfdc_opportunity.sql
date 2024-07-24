SELECT 'PROD_UDM_PERSISTENT.PERSISTENT_SFDC.OPPORTUNITY' AS "record source"
     , 'SalesForce'                                      AS "database"
     , OPP.ID                                            AS "icm id"
     , 'Opportunity'                                     AS "document type"
     , TO_DATE(OPP.CREATEDDATE)                          AS "document date"
     , TO_DATE(OPP.CREATEDDATE)                          AS "creation date"
     , CREATE_USER.NAME                                  AS "created by"
     , CONCAT_WS('\n',
                 'FROM: ' || IFNULL(OPP.BUSINESS_UNIT_LONG__C, ''),
                 'ACCOUNT: ' || IFNULL(OPP_ACCOUNT.NAME, ''),
                 'END USER: ' || IFNULL(END_ACCOUNT.NAME, ''),
                 'NUMBER: ' || IFNULL(OPP.OPPORTUNITY_NUMBER__C, ''),
                 'STAGE: ' || IFNULL(OPP.STAGENAME, ''),
                 'NAME: ' || IFNULL(OPP.NAME, ''),
                 'DESCRIPTION: ' || IFNULL(OPP.DESCRIPTION, '')
       )                                                 AS "document text"
  FROM PROD_UDM_PERSISTENT.PERSISTENT_SFDC.OPPORTUNITY OPP
           LEFT JOIN PROD_UDM_PERSISTENT.PERSISTENT_SFDC.USER2 CREATE_USER
                     ON OPP.CREATEDBYID = CREATE_USER.ID
           LEFT JOIN PROD_UDM_PERSISTENT.PERSISTENT_SFDC.ACCOUNT OPP_ACCOUNT
                     ON OPP.ACCOUNTID = OPP_ACCOUNT.ID
           LEFT JOIN PROD_UDM_PERSISTENT.PERSISTENT_SFDC.ACCOUNT END_ACCOUNT
                     ON OPP.ACCOUNT_END_USER__C = END_ACCOUNT.ID