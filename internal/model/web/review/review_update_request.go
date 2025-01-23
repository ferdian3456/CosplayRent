package review

type ReviewUpdateRequest struct {
	Review_picture *string `json:"review_picture,omitempty"`
	Description    string  `json:"description,omitempty"`
	Rating         int     `json:"rating,omitempty" `
}
