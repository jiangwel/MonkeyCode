package consts

type ReportAction string

const (
	ReportActionAccept      ReportAction = "accept"
	ReportActionSuggest     ReportAction = "suggest"
	ReportActionFileWritten ReportAction = "file_written"
)
