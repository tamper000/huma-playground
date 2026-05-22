package handler

// Create link
type CreateShortLink struct {
	Body struct {
		URL string `json:"url" doc:"original URL to shorten" example:"https://youtu.be/12345bdkj" pattern:"^https://.*" maxLength:"150"`
	}
}

type ShortLinkOutput struct {
	Body struct {
		ID string `json:"id" example:"a1b2c3d4" doc:"short link ID"`
	}
}

type GetLink struct {
	ID string `path:"id" doc:"Short URL ID" example:"12345bd" minLength:"7" maxLength:"7"`
}

type InfoLink struct {
	Body struct {
		Link string `json:"link" doc:"Original link" example:"https://youtu.be/12345bdkj" minLength:"7" maxLength:"7"`
	}
}

type RedirectOutput struct {
	Status   int    `header:"-"`
	Location string `header:"Location"`
}
