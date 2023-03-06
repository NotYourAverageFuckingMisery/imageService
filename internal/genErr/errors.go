package genErr

var (
	ErrRecievingChunkData = New("Can not receive chunk data")
	ErrWritingChunkData   = New("Can not write chunk data")
	ErrFailedToSave       = New("Can not save image")
	ErrOpeningFile        = New("Can not open image file")
	ErrClosingStream      = New("Can not close stream")
	ErrSendingImage       = New("Can not send image")
	ErrReadingChunk       = New("Can not read chunk to buffer")
	ErrSendingChunkData   = New("Can not send chunk data to client")
	ErrGetingInfo         = New("Can not get information about files")
	ErrSendingInfo        = New("Can not send info about files")
)
