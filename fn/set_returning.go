package fn

import "github.com/networkteam/qrb/builder"

// GenerateSeries builds a generate_series functional call.
//
//	generate_series ( start, stop ) → setof integer
//	generate_series ( start, stop ) → setof bigint
//	generate_series ( start, stop, step ) → setof integer
//	generate_series ( start, stop, step ) → setof bigint
//	generate_series ( start, stop, step interval ) → setof timestamp
//	generate_series ( start, stop, step interval ) → setof timestamp with time zone
//
// Generates a series of values from start to stop, with a step size of one if not specified, or step if specified.
//
// Examples from PostgreSQL documentation:
//
//	SELECT * FROM generate_series(2,4);
//	 generate_series
//	-----------------
//	               2
//	               3
//	               4
//	(3 rows)
//
//	SELECT * FROM generate_series(5,1,-2);
//	 generate_series
//	-----------------
//	               5
//	               3
//	               1
//	(3 rows)
//
//	SELECT * FROM generate_series(4,3);
//	 generate_series
//	-----------------
//	(0 rows)
//
//	SELECT current_date + s.a AS dates FROM generate_series(0,14,7) AS s(a);
//	   dates
//	------------
//	 2004-02-05
//	 2004-02-12
//	 2004-02-19
//	(3 rows)
//
//	SELECT * FROM generate_series('2008-03-01 00:00'::timestamp,
//	                              '2008-03-04 12:00', '10 hours');
//	   generate_series
//	---------------------
//	 2008-03-01 00:00:00
//	 2008-03-01 10:00:00
//	 2008-03-01 20:00:00
//	 2008-03-02 06:00:00
//	 2008-03-02 16:00:00
//	 2008-03-03 02:00:00
//	 2008-03-03 12:00:00
//	 2008-03-03 22:00:00
//	 2008-03-04 08:00:00
//	(9 rows)
func GenerateSeries(start builder.Exp, stop builder.Exp, step ...builder.Exp) builder.FuncBuilder {
	args := []builder.Exp{start, stop}
	if len(step) > 0 {
		args = append(args, step[0])
	}
	return builder.Func("generate_series", args...)
}

// GenerateSubscripts builds a generate_subscripts functional call.
//
//	generate_subscripts ( array anyarray, dim integer ) → setof integer
//	generate_subscripts ( array anyarray, dim integer, reverse boolean ) → setof integer
//
// Generates a series comprising the subscripts of the dim'th dimension of the given array.
// When reverse is true, returns the series in reverse order.
//
// Examples from PostgreSQL documentation:
//
//	SELECT generate_subscripts('{NULL,1,NULL,2}'::int[], 1) AS s;
//	 s
//	---
//	 1
//	 2
//	 3
//	 4
//	(4 rows)
//
//	-- unnest a 2D array
//	CREATE OR REPLACE FUNCTION unnest2(anyarray)
//	RETURNS SETOF anyelement AS $$
//	select $1[i][j]
//	   from generate_subscripts($1,1) g1(i),
//	        generate_subscripts($1,2) g2(j);
//	$$ LANGUAGE sql IMMUTABLE;
//	CREATE FUNCTION
//	SELECT * FROM unnest2(ARRAY[[1,2],[3,4]]);
//	 unnest2
//	---------
//	       1
//	       2
//	       3
//	       4
//	(4 rows)
//
// When using generate_subscripts in the FROM clause, it's useful to also have the array value itself:
//
//	-- set returning function WITH ORDINALITY
//	SELECT a AS array, s AS subscript, a[s] AS value
//	FROM (SELECT generate_subscripts(a, 1) AS s, a FROM arrays) foo;
//	     array     | subscript | value
//	---------------+-----------+-------
//	 {-1,-2}       |         1 |    -1
//	 {-1,-2}       |         2 |    -2
//	 {100,200,300} |         1 |   100
//	 {100,200,300} |         2 |   200
//	 {100,200,300} |         3 |   300
//	(5 rows)
func GenerateSubscripts(array builder.Exp, dim builder.Exp, reverse ...builder.Exp) builder.FuncBuilder {
	args := []builder.Exp{array, dim}
	if len(reverse) > 0 {
		args = append(args, reverse[0])
	}
	return builder.Func("generate_subscripts", args...)
}
