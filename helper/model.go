package helper

import (
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
