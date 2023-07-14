CREATE OR REPLACE FUNCTION f_num_nulls(_tbl regclass)
  RETURNS SETOF int AS
$func$
BEGIN
   RETURN QUERY EXECUTE format(
      'SELECT num_nulls(%s) FROM %s'
    , (SELECT string_agg(quote_ident(attname), ', ')  -- column list
       FROM   pg_attribute
       WHERE  attrelid = _tbl
       AND    NOT attisdropped    -- no dropped (dead) columns
       AND    attnum > 0)         -- no system columns
    , _tbl
   );
END
$func$  LANGUAGE plpgsql;