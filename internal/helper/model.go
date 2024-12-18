package helper

import (
	"cosplayrent/internal/model/web/costume"
	"cosplayrent/internal/model/web/user"
)

func ToUserResponse(userResponse user.UserResponse) user.UserResponse {
	return user.UserResponse{
		Id:              userResponse.Id,
		Name:            userResponse.Name,
		Email:           userResponse.Email,
		Address:         userResponse.Address,
		Profile_picture: userResponse.Profile_picture,
		Created_at:      userResponse.Created_at,
	}
}

func ToUserResponses(userResponse []user.UserResponse) []user.UserResponse {
	var UserResponse []user.UserResponse
	for _, userIteration := range userResponse {
		UserResponse = append(UserResponse, ToUserResponse(userIteration))
	}

	return UserResponse
}

func ToCostumeResponses(costumeResponse costume.CostumeResponse) costume.CostumeResponse {
	return costume.CostumeResponse{
		Id:              costumeResponse.Id,
		User_id:         costumeResponse.User_id,
		Username:        costumeResponse.Username,
		Profile_picture: costumeResponse.Profile_picture,
		Name:            costumeResponse.Name,
		Description:     costumeResponse.Description,
		Bahan:           costumeResponse.Bahan,
		Ukuran:          costumeResponse.Ukuran,
		Berat:           costumeResponse.Berat,
		Kategori:        costumeResponse.Kategori,
		Price:           costumeResponse.Price,
		Picture:         costumeResponse.Picture,
		Available:       costumeResponse.Available,
		Created_at:      costumeResponse.Created_at,
	}
}

func ToCostumeResponse(costumeResponse []costume.CostumeResponse) []costume.CostumeResponse {
	var CostumeResponse []costume.CostumeResponse
	for _, costumeIteration := range costumeResponse {
		CostumeResponse = append(CostumeResponse, ToCostumeResponses(costumeIteration))
	}

	return CostumeResponse
}
