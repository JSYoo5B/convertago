package googlechat_test

import (
	"testing"

	"github.com/JSYoo5B/convertago/googlechat"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestChipList_Validate(t *testing.T) {
	v := validator.New()
	require.NoError(t, v.Struct(googlechat.ChipList{Chips: []googlechat.Chip{{Label: "Chip"}}}))
	require.NoError(t, v.Struct(googlechat.Chip{Icon: &googlechat.Icon{KnownIcon: "EMAIL"}}))
	require.Error(t, v.Struct(googlechat.ChipList{}))
	require.Error(t, v.Struct(googlechat.Chip{}))
	require.Error(t, v.Struct(googlechat.ChipList{Layout: "GRID", Chips: []googlechat.Chip{{Label: "Chip"}}}))
}
