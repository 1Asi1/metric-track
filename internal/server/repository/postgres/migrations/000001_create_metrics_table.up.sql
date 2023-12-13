BEGIN TRANSACTION;

   CREATE TABLE tbl_metrics(
       id text not null,
       gauge double precision,
       counter int
   );

COMMIT;