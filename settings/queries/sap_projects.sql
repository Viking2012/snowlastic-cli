SELECT 'PROD_LI.' ||
       'COMMON_DATA.' ||
       'PROJECTS'                     AS "record source"
     , DB                             AS "database"
     , ICM_ID                         AS "icm id"
     , 'Project Definition'           AS "document type"
     , IFNULL(PROJ_PLFAZ, PROJ_ERDAT) AS "document date"
     , PROJ_ERDAT                     AS "creation date"
     , PROJ_ERNAM                     AS "created by"
     , PROJ_POST1                     AS "document text"
  FROM PROD_LI.COMMON_DATA.PROJECTS