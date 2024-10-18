package helper

import (
	"cosplayrent/model/web/costume"
	"cosplayrent/model/web/user"
)

func ToUserResponse(userResponse user.UserResponse) user.UserResponse {
	return user.UserResponse{
		Id:              userResponse.Id,
		Name:            userResponse.Name,
		Email:           userResponse.Email,
		Role:            userResponse.Role,
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
		Id:          costumeResponse.Id,
		User_id:     costumeResponse.User_id,
		Name:        costumeResponse.Name,
		Description: costumeResponse.Description,
		Price:       costumeResponse.Price,
		Picture:     costumeResponse.Picture,
		Available:   costumeResponse.Available,
		Created_at:  costumeResponse.Created_at,
	}
}

func ToCostumeResponse(costumeResponse []costume.CostumeResponse) []costume.CostumeResponse {
	var CostumeResponse []costume.CostumeResponse
	for _, costumeIteration := range costumeResponse {
		CostumeResponse = append(CostumeResponse, ToCostumeResponses(costumeIteration))
	}

	return CostumeResponse
}
