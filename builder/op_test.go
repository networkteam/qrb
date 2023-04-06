package builder_test

import (
	"testing"

	"github.com/networkteam/qrb"
	"github.com/networkteam/qrb/internal/testhelper"
)

func TestOp(t *testing.T) {
	t.Run("cast", func(t *testing.T) {
		t.Run("single exp", func(t *testing.T) {
			b := qrb.Arg("foo").Cast("text")

			testhelper.AssertSQLWriterEquals(
				t,
				"$1::text",
				[]any{"foo"},
				b,
			)
		})

		t.Run("combined exp", func(t *testing.T) {
			b := qrb.N("json_column").JsonExtractText(qrb.Arg("my_field")).Cast("int")

			testhelper.AssertSQLWriterEquals(
				t,
				"(json_column ->> $1)::int",
				[]any{"my_field"},
				b,
			)
		})
	})
}
