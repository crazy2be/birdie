include $(GOROOT)/src/Make.inc

TARG=birdie
GOFILES=\
	conn.go\
	packet.go\
	types.go

include $(GOROOT)/src/Make.pkg