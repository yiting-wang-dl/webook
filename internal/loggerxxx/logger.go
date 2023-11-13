package loggerxxx

import "go.uber.org/zap"

// in main(), Logger = xxx
var Logger *zap.Logger

var CommonLogger *zap.Logger
var SensitiveLogger *zap.Logger
