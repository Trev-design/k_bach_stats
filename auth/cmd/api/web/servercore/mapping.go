package servercore

import (
	"auth_server/cmd/api/domain/types"
	"auth_server/cmd/api/web/webtypes"
)

func toNewAccountDTO(representation *webtypes.NewAccountRepresentation) *types.NewAccountRequestDTO {
	return &types.NewAccountRequestDTO{
		Name:         representation.Name,
		Email:        representation.Email,
		Password:     representation.Password,
		Confirmation: representation.Confirmation,
	}
}
