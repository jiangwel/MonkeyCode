package domain
type ParseResult struct {
	FilePath   string `json:"file_path"`
	Definition string `json:"definition"` 
	Error      string `json:"error,omitempty"` 
	Success    bool   `json:"success"`
}
type SaveAstReq struct { 
	UserID    string   `json:"user_id" validate:"required"` 
	ProjectID string   `json:"project_id" validate:"required"`
	Files     []string `json:"files" validate:"required"`
}