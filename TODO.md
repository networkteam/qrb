# TODOs

* Insert, update, delete statements:
    * [x] Add `Insert` statement
      * [x] Support `WITH` queries
      * [ ] Support `OVERRIDING { SYSTEM | USER } VALUE` clause
      * [x] Support `DEFAULT VALUES`
      * [x] Support `.Query` to add a `SELECT` statement
      * [ ] Support `ON CONFLICT` clause
      * [x] Support `RETURNING` clause
    * [x] Add `Update` statement
      * [x] Support `WITH` queries
      * [ ] Suppport `SetColumnList` to set column names from expressions or a sub-select
      * [x] Support `FROM` clause for joins
      * [ ] Support `WHERE CURRENT OF cursor_name` clause
      * [x] Support `RETURNING` clause
    * [x] Add `Delete` statement
      * [ ] Support `WITH` queries
      * [ ] Support `ONLY` and `table_name *`
      * [ ] Support `USING` clause
      * [ ] Support `WHERE CURRENT OF cursor_name` clause
      * [ ] Support `RETURNING` clause
* Select:
    * [ ] Support locking clauses
    * [ ] Support window functions
* [x] Reduce exported types on `qrb` package
* [x] Check if we want to add `.As()` to `N` to improve select lists and from clauses
    * Not sure for now, since output names and aliases, column aliases and column definitions differ between select
      lists, from items and functions
