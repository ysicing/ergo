// Copyright (c) 2020-2023 ysicing(ysicing@ysicing.cloud) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Y PUBLIC LICENSE 1.0 (YPL 1.0)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package log

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/ysicing/ergo/internal/pkg/util/log/survey"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

func NewUnionLogger(loggers ...Logger) Logger {
	return &unionLogger{
		Loggers: loggers,
	}
}

type unionLogger struct {
	Loggers []Logger

	level logrus.Level
}

func (s *unionLogger) Debug(args ...interface{}) {
	for _, l := range s.Loggers {
		l.Debug(args...)
	}
}

func (s *unionLogger) Debugf(format string, args ...interface{}) {
	for _, l := range s.Loggers {
		l.Debugf(format, args...)
	}
}

func (s *unionLogger) Info(args ...interface{}) {
	for _, l := range s.Loggers {
		l.Info(args...)
	}
}

func (s *unionLogger) Infof(format string, args ...interface{}) {
	for _, l := range s.Loggers {
		l.Infof(format, args...)
	}
}

func (s *unionLogger) Warn(args ...interface{}) {
	for _, l := range s.Loggers {
		l.Warn(args...)
	}
}

func (s *unionLogger) Warnf(format string, args ...interface{}) {
	for _, l := range s.Loggers {
		l.Warnf(format, args...)
	}
}

func (s *unionLogger) Error(args ...interface{}) {
	for _, l := range s.Loggers {
		l.Error(args...)
	}
}

func (s *unionLogger) Errorf(format string, args ...interface{}) {
	for _, l := range s.Loggers {
		l.Errorf(format, args...)
	}
}

func (s *unionLogger) Fatal(args ...interface{}) {
	for _, l := range s.Loggers {
		l.Print(logrus.FatalLevel, args...)
	}
	os.Exit(1)
}

func (s *unionLogger) Fatalf(format string, args ...interface{}) {
	for _, l := range s.Loggers {
		l.Printf(logrus.FatalLevel, format, args...)
	}
	os.Exit(1)
}

func (s *unionLogger) Panic(args ...interface{}) {
	for _, l := range s.Loggers {
		l.Print(logrus.PanicLevel, args...)
	}
	panic(fmt.Sprintln(args...))
}

func (s *unionLogger) Panicf(format string, args ...interface{}) {
	for _, l := range s.Loggers {
		l.Printf(logrus.PanicLevel, format, args...)
	}
	panic(fmt.Sprintln(args...))
}

func (s *unionLogger) Done(args ...interface{}) {
	for _, l := range s.Loggers {
		l.Done(args...)
	}
}

func (s *unionLogger) Donef(format string, args ...interface{}) {
	for _, l := range s.Loggers {
		l.Donef(format, args...)
	}
}

func (s *unionLogger) Print(level logrus.Level, args ...interface{}) {
	switch level {
	case logrus.DebugLevel, logrus.TraceLevel:
		s.Debug(args...)
	case logrus.InfoLevel:
		s.Info(args...)
	case logrus.WarnLevel:
		s.Warn(args...)
	case logrus.ErrorLevel:
		s.Error(args...)
	case logrus.PanicLevel:
		s.Panic(args...)
	case logrus.FatalLevel:
		s.Fatal(args...)
	}
}

func (s *unionLogger) Printf(level logrus.Level, format string, args ...interface{}) {
	switch level {
	case logrus.DebugLevel, logrus.TraceLevel:
		s.Debugf(format, args...)
	case logrus.InfoLevel:
		s.Infof(format, args...)
	case logrus.WarnLevel:
		s.Warnf(format, args...)
	case logrus.ErrorLevel:
		s.Errorf(format, args...)
	case logrus.PanicLevel:
		s.Panicf(format, args...)
	case logrus.FatalLevel:
		s.Fatalf(format, args...)
	}
}

func (s *unionLogger) StartWait(message string) {
	for _, l := range s.Loggers {
		l.StartWait(message)
	}
}

func (s *unionLogger) StopWait() {
	for _, l := range s.Loggers {
		l.StopWait()
	}
}

func (s *unionLogger) Write(message []byte) (int, error) {
	errs := []error{}
	for _, l := range s.Loggers {
		_, err := l.Write(message)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return len(message), utilerrors.NewAggregate(errs)
}

func (s *unionLogger) WriteString(message string) {
	for _, l := range s.Loggers {
		l.WriteString(message)
	}
}

func (s *unionLogger) Question(params *survey.QuestionOptions) (string, error) {
	errs := []error{}
	for _, l := range s.Loggers {
		answer, err := l.Question(params)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		return answer, nil
	}

	return "", utilerrors.NewAggregate(errs)
}

func (s *unionLogger) SetLevel(level logrus.Level) {
	for _, l := range s.Loggers {
		l.SetLevel(level)
	}

	s.level = level
}

func (s *unionLogger) GetLevel() logrus.Level {
	return s.level
}
