package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Test method signatures GetV2MethodFromID
*/
func TestGetV2MethodFromID(t *testing.T) {
	t.Run("Test SwapExactTokensForTokens", func(t *testing.T) {

		method, ok := GetV2MethodFromID("0x38ed1739")
		assert.True(t, ok)
		assert.Equal(t, SwapExactTokensForTokens, method)
	})
}
