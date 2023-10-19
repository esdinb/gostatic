
#include "xml_save.h"

extern int ioWrite(void *, const char *, int len);
extern int ioClose(void *);

static int ioNop(void *ctx) { return 0; }

xmlSaveCtxtPtr saveToIO(void *ioctx, const char *encoding, int options) {
	return xmlSaveToIO((xmlOutputWriteCallback) ioWrite, (xmlOutputCloseCallback) ioNop, ioctx, encoding, options);
}

xmlSaveCtxtPtr saveToIOCloser(void *ioctx, const char *encoding, int options) {
	return xmlSaveToIO((xmlOutputWriteCallback) ioWrite, (xmlOutputCloseCallback) ioClose, ioctx, encoding, options);
}
