package review

type ReviewUpdateRequest struct {
	Review_picture *string `json:"review_picture,omitempty"`
	Description    string  `validate:"min=5,max=1000" json:"description,omitempty"`
	Rating         int     `json:"rating,omitempty" `
}
