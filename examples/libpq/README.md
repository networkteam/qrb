# Example for qrbsql using github.com/lib/pq

## Running

Create a PostreSQL database and import the schema:

```bash
createdb qrb-examples
psql -d qrb-examples -f ../schema.sql
```

Run the example:

```bash
go run . books list --author "Harp"
```
