SELECT 'PROD_UDM_PERSISTENT.PERSISTENT_SFDC.EMAILMESSAGE'         AS "record source"
     , 'SalesForce'                                               AS "database"
     , EMAILMESSAGE.ID                                            AS "icm id"
     , 'Email Message'                                            AS "document type"
     , IFNULL(EMAILMESSAGE.MESSAGEDATE, EMAILMESSAGE.CREATEDDATE) AS "document date"
     , EMAILMESSAGE.CREATEDDATE                                   AS "creation date"
     , USER2.NAME                                                 AS "created by"
     , CONCAT_WS('\n',
                 'FROM: ' || IFNULL(EMAILMESSAGE.FROMNAME, '')
                     || ' <'
                     || IFNULL(EMAILMESSAGE.FROMADDRESS, '')
                     || '>',
                 'TO: ' || IFNULL(EMAILMESSAGE.TOADDRESS, ''),
                 'CC: ' || IFNULL(EMAILMESSAGE.CCADDRESS, ''),
                 'BCC: ' || IFNULL(EMAILMESSAGE.BCCADDRESS, ''),
                 'SUBJECT: ' || IFNULL(EMAILMESSAGE.SUBJECT, ''),
                 'BODY: ' || IFNULL(EMAILMESSAGE.TEXTBODY, '')
       )                                                          AS "document text"
  FROM PROD_UDM_PERSISTENT.PERSISTENT_SFDC.EMAILMESSAGE
           LEFT JOIN PROD_UDM_PERSISTENT.PERSISTENT_SFDC.USER2
                     ON EMAILMESSAGE.CREATEDBYID = USER2.ID