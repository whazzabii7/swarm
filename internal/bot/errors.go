package bot

import "errors"

var ErrCreateStderrFailed = errors.New("failed to create stderr pipe")
var ErrCreateStdinFailed = errors.New("failed to create stdin pipe")
var ErrCreateStdoutFailed = errors.New("failed to create stdout pipe")
var ErrStartBotFailed = errors.New("failed to start bot")
