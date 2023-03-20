package genErr

var (
	ErrRecievingData = New("Can not receive image data")
	ErrWritingData   = New("Can not write image data")
	ErrFailedToSave  = New("Can not save image")
	ErrOpeningFile   = New("Can not open image file")
	ErrClosingStream = New("Can not close stream")
	ErrSendingImage  = New("Can not send image")
	ErrReadingFile   = New("Can not read file to bytes")
	ErrSendingData   = New("Can not send data to client")
	ErrGetingInfo    = New("Can not get information about files")
	ErrSendingInfo   = New("Can not send info about files")
)
