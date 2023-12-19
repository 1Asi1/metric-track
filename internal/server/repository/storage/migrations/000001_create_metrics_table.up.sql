BEGIN TRANSACTION;

   CREATE TABLE tbl_metrics(
       id text primary key not null,
       gauge double precision null,
       counter int null
   );

COMMIT;