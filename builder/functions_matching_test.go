package builder_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/networkteam/qrb"
)

func TestMatching(t *testing.T) {
	t.Run("like", func(t *testing.T) {
		b := qrb.N("name").Like(qrb.String("foo%"))

		sql, args, _ := qrb.Build(b).ToSQL()

		assert.Equal(t, "name LIKE 'foo%'", sql)
		assert.Empty(t, args)
	})

	t.Run("similar to", func(t *testing.T) {
		b := qrb.N("name").SimilarTo(qrb.String("%(b|d)%"))

		sql, args, _ := qrb.Build(b).ToSQL()

		assert.Equal(t, "name SIMILAR TO '%(b|d)%'", sql)
		assert.Empty(t, args)
	})
}
