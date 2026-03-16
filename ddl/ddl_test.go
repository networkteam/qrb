package ddl_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/networkteam/qrb"
	"github.com/networkteam/qrb/ddl"
	"github.com/networkteam/qrb/internal/testhelper"
)

// --- CREATE TABLE ---

func TestCreateTable(t *testing.T) {
	t.Run("single column", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("users")).
			Column("id", "INTEGER")

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE users (id INTEGER)`,
			nil, q,
		)
	})

	t.Run("multiple columns", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("users")).
			Column("id", "INTEGER").
			Column("name", "TEXT").
			Column("email", "VARCHAR(255)")

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE users (id INTEGER, name TEXT, email VARCHAR(255))`,
			nil, q,
		)
	})

	t.Run("column with NOT NULL", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("users")).
			Column("name", "TEXT").NotNull()

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE users (name TEXT NOT NULL)`,
			nil, q,
		)
	})

	t.Run("column with DEFAULT", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("users")).
			Column("active", "BOOLEAN").Default(qrb.Bool(true))

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE users (active BOOLEAN DEFAULT true)`,
			nil, q,
		)
	})

	t.Run("column with PRIMARY KEY", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("users")).
			Column("id", "INTEGER").PrimaryKey()

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE users (id INTEGER PRIMARY KEY)`,
			nil, q,
		)
	})

	t.Run("column with UNIQUE", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("users")).
			Column("email", "TEXT").Unique()

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE users (email TEXT UNIQUE)`,
			nil, q,
		)
	})

	t.Run("column with CHECK", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("products")).
			Column("price", "NUMERIC").Check(qrb.N("price").Gt(qrb.Int(0)))

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE products (price NUMERIC CHECK (price > 0))`,
			nil, q,
		)
	})

	t.Run("column with REFERENCES", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("orders")).
			Column("user_id", "INTEGER").References(qrb.N("users"), "id")

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE orders (user_id INTEGER REFERENCES users (id))`,
			nil, q,
		)
	})

	t.Run("column with REFERENCES ON DELETE CASCADE", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("orders")).
			Column("user_id", "INTEGER").References(qrb.N("users"), "id").OnDelete("CASCADE")

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE orders (user_id INTEGER REFERENCES users (id) ON DELETE CASCADE)`,
			nil, q,
		)
	})

	t.Run("column with REFERENCES ON DELETE and ON UPDATE", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("orders")).
			Column("user_id", "INTEGER").References(qrb.N("users"), "id").
			OnDelete("CASCADE").OnUpdate("SET NULL")

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE orders (user_id INTEGER REFERENCES users (id) ON DELETE CASCADE ON UPDATE SET NULL)`,
			nil, q,
		)
	})

	t.Run("all column constraints", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("users")).
			Column("id", "SERIAL").PrimaryKey().
			Column("name", "TEXT").NotNull().
			Column("email", "TEXT").NotNull().Unique().
			Column("age", "INTEGER").Check(qrb.N("age").Gte(qrb.Int(0))).
			Column("dept_id", "INTEGER").References(qrb.N("departments"), "id").OnDelete("SET NULL")

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE users (
				id SERIAL PRIMARY KEY,
				name TEXT NOT NULL,
				email TEXT NOT NULL UNIQUE,
				age INTEGER CHECK (age >= 0),
				dept_id INTEGER REFERENCES departments (id) ON DELETE SET NULL
			)`,
			nil, q,
		)
	})

	t.Run("table-level PRIMARY KEY", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("order_items")).
			Column("order_id", "INTEGER").
			Column("item_id", "INTEGER").CreateTableBuilder.
			PrimaryKey("order_id", "item_id")

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE order_items (order_id INTEGER, item_id INTEGER, PRIMARY KEY (order_id, item_id))`,
			nil, q,
		)
	})

	t.Run("table-level UNIQUE", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("users")).
			Column("first_name", "TEXT").
			Column("last_name", "TEXT").CreateTableBuilder.
			Unique("first_name", "last_name")

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE users (first_name TEXT, last_name TEXT, UNIQUE (first_name, last_name))`,
			nil, q,
		)
	})

	t.Run("table-level FOREIGN KEY", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("orders")).
			Column("user_id", "INTEGER").
			ForeignKey("user_id").References(qrb.N("users"), "id")

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE orders (user_id INTEGER, FOREIGN KEY (user_id) REFERENCES users (id))`,
			nil, q,
		)
	})

	t.Run("table-level FOREIGN KEY with ON DELETE and ON UPDATE", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("orders")).
			Column("user_id", "INTEGER").
			ForeignKey("user_id").References(qrb.N("users"), "id").
			OnDelete("CASCADE").OnUpdate("RESTRICT")

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE orders (user_id INTEGER, FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE RESTRICT)`,
			nil, q,
		)
	})

	t.Run("table-level CHECK", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("products")).
			Column("price", "NUMERIC").
			Column("discounted_price", "NUMERIC").CreateTableBuilder.
			Check(qrb.N("discounted_price").Lt(qrb.N("price")))

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE products (price NUMERIC, discounted_price NUMERIC, CHECK (discounted_price < price))`,
			nil, q,
		)
	})

	t.Run("named constraint", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("users")).
			Column("id", "INTEGER").
			Constraint("users_pk").PrimaryKey("id")

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE users (id INTEGER, CONSTRAINT users_pk PRIMARY KEY (id))`,
			nil, q,
		)
	})

	t.Run("named FOREIGN KEY constraint", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("orders")).
			Column("user_id", "INTEGER").
			Constraint("fk_user").ForeignKey("user_id").References(qrb.N("users"), "id")

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE orders (user_id INTEGER, CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id))`,
			nil, q,
		)
	})

	t.Run("named CHECK constraint", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("products")).
			Column("price", "NUMERIC").
			Constraint("positive_price").Check(qrb.N("price").Gt(qrb.Int(0)))

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE products (price NUMERIC, CONSTRAINT positive_price CHECK (price > 0))`,
			nil, q,
		)
	})

	t.Run("IF NOT EXISTS", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("users")).IfNotExists().
			Column("id", "INTEGER")

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE IF NOT EXISTS users (id INTEGER)`,
			nil, q,
		)
	})

	t.Run("UNLOGGED", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("sessions")).Unlogged().
			Column("id", "TEXT")

		testhelper.AssertSQLWriterEquals(t,
			`CREATE UNLOGGED TABLE sessions (id TEXT)`,
			nil, q,
		)
	})

	t.Run("reserved keyword column name", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("t")).
			Column("select", "TEXT").
			Column("from", "TEXT")

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE t ("select" TEXT, "from" TEXT)`,
			nil, q,
		)
	})

	t.Run("schema-qualified table name", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("myschema.users")).
			Column("id", "INTEGER")

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE myschema.users (id INTEGER)`,
			nil, q,
		)
	})

	t.Run("column with DEFAULT expression", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("events")).
			Column("created_at", "TIMESTAMP WITH TIME ZONE").Default(qrb.Func("now"))

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE events (created_at TIMESTAMP WITH TIME ZONE DEFAULT now())`,
			nil, q,
		)
	})

	t.Run("REFERENCES without columns", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("orders")).
			Column("user_id", "INTEGER").References(qrb.N("users"))

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE orders (user_id INTEGER REFERENCES users)`,
			nil, q,
		)
	})

	t.Run("column after references chain continues correctly", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("orders")).
			Column("user_id", "INTEGER").References(qrb.N("users"), "id").OnDelete("CASCADE").
			Column("product_id", "INTEGER")

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE orders (user_id INTEGER REFERENCES users (id) ON DELETE CASCADE, product_id INTEGER)`,
			nil, q,
		)
	})

	t.Run("named UNIQUE constraint", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("users")).
			Column("first_name", "TEXT").
			Column("last_name", "TEXT").
			Constraint("uq_name").Unique("first_name", "last_name")

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE users (first_name TEXT, last_name TEXT, CONSTRAINT uq_name UNIQUE (first_name, last_name))`,
			nil, q,
		)
	})

	t.Run("table-level FOREIGN KEY with ON UPDATE only", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("orders")).
			Column("user_id", "INTEGER").
			ForeignKey("user_id").References(qrb.N("users"), "id").
			OnUpdate("CASCADE")

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE orders (user_id INTEGER, FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE)`,
			nil, q,
		)
	})

	t.Run("columns and table constraints combined", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("order_items")).
			Column("order_id", "INTEGER").NotNull().
			Column("product_id", "INTEGER").NotNull().
			Column("quantity", "INTEGER").Default(qrb.Int(1)).CreateTableBuilder.
			PrimaryKey("order_id", "product_id").
			ForeignKey("order_id").References(qrb.N("orders"), "id").OnDelete("CASCADE").
			ForeignKey("product_id").References(qrb.N("products"), "id").OnDelete("RESTRICT").
			Check(qrb.N("quantity").Gt(qrb.Int(0)))

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE order_items (
				order_id INTEGER NOT NULL,
				product_id INTEGER NOT NULL,
				quantity INTEGER DEFAULT 1,
				PRIMARY KEY (order_id, product_id),
				FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE,
				FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE RESTRICT,
				CHECK (quantity > 0)
			)`,
			nil, q,
		)
	})

	t.Run("LIKE", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("new_table")).
			Like(qrb.N("original_table"))

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE new_table (LIKE original_table)`,
			nil, q,
		)
	})

	t.Run("LIKE INCLUDING ALL", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("target_schema.new_table")).
			Like(qrb.N("source_schema.original_table")).IncludingAll()

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE target_schema.new_table (LIKE source_schema.original_table INCLUDING ALL)`,
			nil, q,
		)
	})

	t.Run("LIKE INCLUDING ALL EXCLUDING INDEXES", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("new_table")).
			Like(qrb.N("original_table")).IncludingAll().ExcludingIndexes()

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE new_table (LIKE original_table INCLUDING ALL EXCLUDING INDEXES)`,
			nil, q,
		)
	})

	t.Run("LIKE INCLUDING ALL with IF NOT EXISTS", func(t *testing.T) {
		q := ddl.CreateTable(qrb.N("target_schema.new_table")).IfNotExists().
			Like(qrb.N("source_schema.original_table")).IncludingAll()

		testhelper.AssertSQLWriterEquals(t,
			`CREATE TABLE IF NOT EXISTS target_schema.new_table (LIKE source_schema.original_table INCLUDING ALL)`,
			nil, q,
		)
	})
}

// --- CREATE SCHEMA ---

func TestCreateSchema(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		q := ddl.CreateSchema(qrb.N("myschema"))

		testhelper.AssertSQLWriterEquals(t,
			`CREATE SCHEMA myschema`,
			nil, q,
		)
	})

	t.Run("IF NOT EXISTS", func(t *testing.T) {
		q := ddl.CreateSchema(qrb.N("myschema")).IfNotExists()

		testhelper.AssertSQLWriterEquals(t,
			`CREATE SCHEMA IF NOT EXISTS myschema`,
			nil, q,
		)
	})

	t.Run("AUTHORIZATION", func(t *testing.T) {
		q := ddl.CreateSchema(qrb.N("myschema")).Authorization("admin")

		testhelper.AssertSQLWriterEquals(t,
			`CREATE SCHEMA myschema AUTHORIZATION admin`,
			nil, q,
		)
	})

	t.Run("IF NOT EXISTS with AUTHORIZATION", func(t *testing.T) {
		q := ddl.CreateSchema(qrb.N("myschema")).IfNotExists().Authorization("admin")

		testhelper.AssertSQLWriterEquals(t,
			`CREATE SCHEMA IF NOT EXISTS myschema AUTHORIZATION admin`,
			nil, q,
		)
	})
}

// --- CREATE INDEX ---

func TestCreateIndex(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		q := ddl.CreateIndex("idx_users_name").
			On(qrb.N("users")).
			Columns(qrb.N("name"))

		testhelper.AssertSQLWriterEquals(t,
			`CREATE INDEX idx_users_name ON users (name)`,
			nil, q,
		)
	})

	t.Run("UNIQUE", func(t *testing.T) {
		q := ddl.CreateIndex("idx_users_email").Unique().
			On(qrb.N("users")).
			Columns(qrb.N("email"))

		testhelper.AssertSQLWriterEquals(t,
			`CREATE UNIQUE INDEX idx_users_email ON users (email)`,
			nil, q,
		)
	})

	t.Run("CONCURRENTLY", func(t *testing.T) {
		q := ddl.CreateIndex("idx_users_name").Concurrently().
			On(qrb.N("users")).
			Columns(qrb.N("name"))

		testhelper.AssertSQLWriterEquals(t,
			`CREATE INDEX CONCURRENTLY idx_users_name ON users (name)`,
			nil, q,
		)
	})

	t.Run("IF NOT EXISTS", func(t *testing.T) {
		q := ddl.CreateIndex("idx_users_name").IfNotExists().
			On(qrb.N("users")).
			Columns(qrb.N("name"))

		testhelper.AssertSQLWriterEquals(t,
			`CREATE INDEX IF NOT EXISTS idx_users_name ON users (name)`,
			nil, q,
		)
	})

	t.Run("USING method", func(t *testing.T) {
		q := ddl.CreateIndex("idx_users_data").
			On(qrb.N("users")).
			Using("gin").
			Columns(qrb.N("data"))

		testhelper.AssertSQLWriterEquals(t,
			`CREATE INDEX idx_users_data ON users USING gin (data)`,
			nil, q,
		)
	})

	t.Run("multiple columns", func(t *testing.T) {
		q := ddl.CreateIndex("idx_users_name_email").
			On(qrb.N("users")).
			Columns(qrb.N("first_name"), qrb.N("last_name"))

		testhelper.AssertSQLWriterEquals(t,
			`CREATE INDEX idx_users_name_email ON users (first_name, last_name)`,
			nil, q,
		)
	})

	t.Run("INCLUDE", func(t *testing.T) {
		q := ddl.CreateIndex("idx_users_name").
			On(qrb.N("users")).
			Columns(qrb.N("name")).
			Include("email")

		testhelper.AssertSQLWriterEquals(t,
			`CREATE INDEX idx_users_name ON users (name) INCLUDE (email)`,
			nil, q,
		)
	})

	t.Run("WHERE partial index", func(t *testing.T) {
		q := ddl.CreateIndex("idx_active_users").
			On(qrb.N("users")).
			Columns(qrb.N("name")).
			Where(qrb.N("active").Eq(qrb.Bool(true)))

		testhelper.AssertSQLWriterEquals(t,
			`CREATE INDEX idx_active_users ON users (name) WHERE active = true`,
			nil, q,
		)
	})

	t.Run("combined options", func(t *testing.T) {
		q := ddl.CreateIndex("idx_active_users_email").
			Unique().Concurrently().IfNotExists().
			On(qrb.N("users")).
			Using("btree").
			Columns(qrb.N("email")).
			Include("name").
			Where(qrb.N("active").Eq(qrb.Bool(true)))

		testhelper.AssertSQLWriterEquals(t,
			`CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx_active_users_email ON users USING btree (email) INCLUDE (name) WHERE active = true`,
			nil, q,
		)
	})

	t.Run("schema-qualified table", func(t *testing.T) {
		q := ddl.CreateIndex("idx_name").
			On(qrb.N("myschema.users")).
			Columns(qrb.N("name"))

		testhelper.AssertSQLWriterEquals(t,
			`CREATE INDEX idx_name ON myschema.users (name)`,
			nil, q,
		)
	})
}

// --- DROP TABLE ---

func TestDropTable(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		q := ddl.DropTable(qrb.N("users"))

		testhelper.AssertSQLWriterEquals(t,
			`DROP TABLE users`,
			nil, q,
		)
	})

	t.Run("IF EXISTS", func(t *testing.T) {
		q := ddl.DropTable(qrb.N("users")).IfExists()

		testhelper.AssertSQLWriterEquals(t,
			`DROP TABLE IF EXISTS users`,
			nil, q,
		)
	})

	t.Run("CASCADE", func(t *testing.T) {
		q := ddl.DropTable(qrb.N("users")).Cascade()

		testhelper.AssertSQLWriterEquals(t,
			`DROP TABLE users CASCADE`,
			nil, q,
		)
	})

	t.Run("RESTRICT", func(t *testing.T) {
		q := ddl.DropTable(qrb.N("users")).Restrict()

		testhelper.AssertSQLWriterEquals(t,
			`DROP TABLE users RESTRICT`,
			nil, q,
		)
	})

	t.Run("multiple targets", func(t *testing.T) {
		q := ddl.DropTable(qrb.N("users"), qrb.N("orders"))

		testhelper.AssertSQLWriterEquals(t,
			`DROP TABLE users, orders`,
			nil, q,
		)
	})

	t.Run("IF EXISTS CASCADE multiple targets", func(t *testing.T) {
		q := ddl.DropTable(qrb.N("users"), qrb.N("orders")).IfExists().Cascade()

		testhelper.AssertSQLWriterEquals(t,
			`DROP TABLE IF EXISTS users, orders CASCADE`,
			nil, q,
		)
	})
}

// --- DROP SCHEMA ---

func TestDropSchema(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		q := ddl.DropSchema(qrb.N("myschema"))

		testhelper.AssertSQLWriterEquals(t,
			`DROP SCHEMA myschema`,
			nil, q,
		)
	})

	t.Run("IF EXISTS", func(t *testing.T) {
		q := ddl.DropSchema(qrb.N("myschema")).IfExists()

		testhelper.AssertSQLWriterEquals(t,
			`DROP SCHEMA IF EXISTS myschema`,
			nil, q,
		)
	})

	t.Run("CASCADE", func(t *testing.T) {
		q := ddl.DropSchema(qrb.N("myschema")).Cascade()

		testhelper.AssertSQLWriterEquals(t,
			`DROP SCHEMA myschema CASCADE`,
			nil, q,
		)
	})

	t.Run("RESTRICT", func(t *testing.T) {
		q := ddl.DropSchema(qrb.N("myschema")).Restrict()

		testhelper.AssertSQLWriterEquals(t,
			`DROP SCHEMA myschema RESTRICT`,
			nil, q,
		)
	})

	t.Run("multiple targets", func(t *testing.T) {
		q := ddl.DropSchema(qrb.N("schema1"), qrb.N("schema2"))

		testhelper.AssertSQLWriterEquals(t,
			`DROP SCHEMA schema1, schema2`,
			nil, q,
		)
	})
}

// --- ALTER TABLE ---

func TestAlterTable(t *testing.T) {
	t.Run("ADD COLUMN", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("users")).
			AddColumn("email", "TEXT")

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE users ADD COLUMN email TEXT`,
			nil, q,
		)
	})

	t.Run("ADD COLUMN with NOT NULL", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("users")).
			AddColumn("email", "TEXT").NotNull()

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE users ADD COLUMN email TEXT NOT NULL`,
			nil, q,
		)
	})

	t.Run("ADD COLUMN with DEFAULT", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("users")).
			AddColumn("active", "BOOLEAN").Default(qrb.Bool(true))

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE users ADD COLUMN active BOOLEAN DEFAULT true`,
			nil, q,
		)
	})

	t.Run("ADD COLUMN with UNIQUE", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("users")).
			AddColumn("email", "TEXT").Unique()

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE users ADD COLUMN email TEXT UNIQUE`,
			nil, q,
		)
	})

	t.Run("ADD COLUMN with REFERENCES", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("orders")).
			AddColumn("user_id", "INTEGER").References(qrb.N("users"), "id")

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE orders ADD COLUMN user_id INTEGER REFERENCES users (id)`,
			nil, q,
		)
	})

	t.Run("ADD COLUMN with REFERENCES ON DELETE CASCADE", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("orders")).
			AddColumn("user_id", "INTEGER").References(qrb.N("users"), "id").OnDelete("CASCADE")

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE orders ADD COLUMN user_id INTEGER REFERENCES users (id) ON DELETE CASCADE`,
			nil, q,
		)
	})

	t.Run("ADD COLUMN IF NOT EXISTS", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("users")).
			AddColumnIfNotExists("email", "TEXT")

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE users ADD COLUMN IF NOT EXISTS email TEXT`,
			nil, q,
		)
	})

	t.Run("DROP COLUMN", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("users")).
			DropColumn("email")

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE users DROP COLUMN email`,
			nil, q,
		)
	})

	t.Run("DROP COLUMN IF EXISTS", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("users")).
			DropColumnIfExists("email")

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE users DROP COLUMN IF EXISTS email`,
			nil, q,
		)
	})

	t.Run("ADD CONSTRAINT PRIMARY KEY", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("users")).
			AddConstraint("users_pk").PrimaryKey("id")

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE users ADD CONSTRAINT users_pk PRIMARY KEY (id)`,
			nil, q,
		)
	})

	t.Run("ADD CONSTRAINT UNIQUE", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("users")).
			AddConstraint("users_email_unique").Unique("email")

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE users ADD CONSTRAINT users_email_unique UNIQUE (email)`,
			nil, q,
		)
	})

	t.Run("ADD CONSTRAINT FOREIGN KEY", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("orders")).
			AddConstraint("fk_user").ForeignKey("user_id").References(qrb.N("users"), "id")

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE orders ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id)`,
			nil, q,
		)
	})

	t.Run("ADD CONSTRAINT FOREIGN KEY with ON DELETE", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("orders")).
			AddConstraint("fk_user").ForeignKey("user_id").References(qrb.N("users"), "id").
			OnDelete("CASCADE")

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE orders ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE`,
			nil, q,
		)
	})

	t.Run("ADD CONSTRAINT CHECK", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("products")).
			AddConstraint("positive_price").Check(qrb.N("price").Gt(qrb.Int(0)))

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE products ADD CONSTRAINT positive_price CHECK (price > 0)`,
			nil, q,
		)
	})

	t.Run("DROP CONSTRAINT", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("users")).
			DropConstraint("users_pk")

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE users DROP CONSTRAINT users_pk`,
			nil, q,
		)
	})

	t.Run("DROP CONSTRAINT IF EXISTS", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("users")).
			DropConstraintIfExists("users_pk")

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE users DROP CONSTRAINT IF EXISTS users_pk`,
			nil, q,
		)
	})

	t.Run("RENAME COLUMN", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("users")).
			RenameColumn("name", "full_name")

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE users RENAME COLUMN name TO full_name`,
			nil, q,
		)
	})

	t.Run("RENAME TO", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("users")).
			RenameTo("people")

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE users RENAME TO people`,
			nil, q,
		)
	})

	t.Run("ALTER COLUMN TYPE", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("users")).
			AlterColumn("name").Type("VARCHAR(255)")

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE users ALTER COLUMN name TYPE VARCHAR(255)`,
			nil, q,
		)
	})

	t.Run("ALTER COLUMN SET DEFAULT", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("users")).
			AlterColumn("active").SetDefault(qrb.Bool(true))

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE users ALTER COLUMN active SET DEFAULT true`,
			nil, q,
		)
	})

	t.Run("ALTER COLUMN DROP DEFAULT", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("users")).
			AlterColumn("active").DropDefault()

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE users ALTER COLUMN active DROP DEFAULT`,
			nil, q,
		)
	})

	t.Run("ALTER COLUMN SET NOT NULL", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("users")).
			AlterColumn("name").SetNotNull()

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE users ALTER COLUMN name SET NOT NULL`,
			nil, q,
		)
	})

	t.Run("ALTER COLUMN DROP NOT NULL", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("users")).
			AlterColumn("name").DropNotNull()

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE users ALTER COLUMN name DROP NOT NULL`,
			nil, q,
		)
	})

	t.Run("IF EXISTS", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("users")).IfExists().
			AddColumn("email", "TEXT")

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE IF EXISTS users ADD COLUMN email TEXT`,
			nil, q,
		)
	})

	t.Run("multiple actions", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("users")).
			AddColumn("email", "TEXT").
			DropColumn("old_email")

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE users ADD COLUMN email TEXT, DROP COLUMN old_email`,
			nil, q,
		)
	})

	t.Run("reserved keyword column names", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("t")).
			RenameColumn("select", "from")

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE t RENAME COLUMN "select" TO "from"`,
			nil, q,
		)
	})

	t.Run("ADD COLUMN with PRIMARY KEY", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("users")).
			AddColumn("id", "SERIAL").PrimaryKey()

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE users ADD COLUMN id SERIAL PRIMARY KEY`,
			nil, q,
		)
	})

	t.Run("ADD COLUMN with CHECK", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("products")).
			AddColumn("price", "NUMERIC").Check(qrb.N("price").Gt(qrb.Int(0)))

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE products ADD COLUMN price NUMERIC CHECK (price > 0)`,
			nil, q,
		)
	})

	t.Run("ADD COLUMN chaining to more actions", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("users")).
			AddColumn("email", "TEXT").NotNull().
			AddColumn("age", "INTEGER").Check(qrb.N("age").Gte(qrb.Int(0)))

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE users ADD COLUMN email TEXT NOT NULL, ADD COLUMN age INTEGER CHECK (age >= 0)`,
			nil, q,
		)
	})

	t.Run("ALTER COLUMN combined with other actions", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("users")).
			AlterColumn("name").Type("VARCHAR(255)").
			AlterColumn("name").SetNotNull().
			AlterColumn("email").SetDefault(qrb.String("unknown"))

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE users ALTER COLUMN name TYPE VARCHAR(255), ALTER COLUMN name SET NOT NULL, ALTER COLUMN email SET DEFAULT 'unknown'`,
			nil, q,
		)
	})

	t.Run("ADD CONSTRAINT FOREIGN KEY with ON DELETE and ON UPDATE", func(t *testing.T) {
		q := ddl.AlterTable(qrb.N("orders")).
			AddConstraint("fk_user").ForeignKey("user_id").References(qrb.N("users"), "id").
			OnDelete("CASCADE").OnUpdate("SET NULL")

		testhelper.AssertSQLWriterEquals(t,
			`ALTER TABLE orders ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE SET NULL`,
			nil, q,
		)
	})
}

// --- CREATE FUNCTION ---

func TestCreateFunction(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		q := ddl.CreateFunction(qrb.N("my_func")).
			Returns("trigger").
			Language("plpgsql").
			Body("BEGIN\n    RETURN NEW;\nEND;")

		sql, args, err := qrb.Build(q).ToSQL()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assert.Empty(t, args)
		assert.Equal(t, "CREATE FUNCTION my_func() RETURNS trigger LANGUAGE plpgsql AS $$\nBEGIN\n    RETURN NEW;\nEND;\n$$", sql)
	})

	t.Run("OR REPLACE with schema-qualified name", func(t *testing.T) {
		q := ddl.CreateFunction(qrb.N("my_schema.trigger_set_timestamp")).
			OrReplace().
			Returns("trigger").
			Language("plpgsql").
			Body("BEGIN\n    IF NEW.updated_at IS NOT DISTINCT FROM OLD.updated_at THEN\n        NEW.updated_at = NOW();\n    END IF;\n    RETURN NEW;\nEND;")

		sql, args, err := qrb.Build(q).ToSQL()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assert.Empty(t, args)
		assert.Equal(t, "CREATE OR REPLACE FUNCTION my_schema.trigger_set_timestamp() RETURNS trigger LANGUAGE plpgsql AS $$\nBEGIN\n    IF NEW.updated_at IS NOT DISTINCT FROM OLD.updated_at THEN\n        NEW.updated_at = NOW();\n    END IF;\n    RETURN NEW;\nEND;\n$$", sql)
	})

	t.Run("custom dollar tag", func(t *testing.T) {
		q := ddl.CreateFunction(qrb.N("my_func")).
			Returns("trigger").
			Language("plpgsql").
			Body("BEGIN\n    RETURN NEW;\nEND;").
			DollarTag("fn")

		sql, args, err := qrb.Build(q).ToSQL()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assert.Empty(t, args)
		assert.Equal(t, "CREATE FUNCTION my_func() RETURNS trigger LANGUAGE plpgsql AS $fn$\nBEGIN\n    RETURN NEW;\nEND;\n$fn$", sql)
	})
}
