package main

import "strconv"

type ircMode int

const (
	ircBold      ircMode = 2
	ircColor     ircMode = 3
	ircItalic    ircMode = 19
	ircCReset    ircMode = 15
	ircCReverse  ircMode = 22
	ircUnderline ircMode = 31
)

type iColor int

const (
	cWhite iColor = iota
	cBlack
	cBlue
	cGreen
	cRed
	cBrown
	cPurple
	cOrange
	cYellow
	cLime
	cTeal
	cCyan
	cRoyal
	cPink
	cGrey
	cSilver
	cNone = -1
)

func setIrcMode(mode ircMode) string {
	return string(byte(mode))
}

func setIrcColor(fgColor iColor, bgColor iColor) string {
	if bgColor != cNone {
		return setIrcMode(ircColor) + strconv.Itoa(int(fgColor)) + "," + strconv.Itoa(int(bgColor))
	} else {
		return setIrcMode(ircColor) + strconv.Itoa(int(fgColor))
	}
}

func colorMsg(msg string, fgColor iColor, bgColor iColor) string {
	return setIrcColor(fgColor, bgColor) + msg + setIrcMode(ircCReset)
}
