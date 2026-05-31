package bot

import "errors"

var ErrFailCreateStderr = errors.New("failed to create stderr pipe")
var ErrFailCreateStdin = errors.New("failed to create stdin pipe")
var ErrFailCreateStdout = errors.New("failed to create stdout pipe")
var ErrFailStartBot = errors.New("failed to start bot")
var ErrFailUnpackArgs = errors.New("failed to unpack arguments")
