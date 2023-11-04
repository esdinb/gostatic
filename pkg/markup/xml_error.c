
#include "xml_error.h"

extern void go_error_callback(void *, xmlErrorPtr);
extern void go_error_print_callback(void *, const char *);

void
structured_error_func(void *userData, xmlErrorPtr error) {
    go_error_callback(userData, error);
}

#define ERROR_BUFFER_SIZE 1024

void
generic_error_func(void *userData, const char *error, ...) {
    va_list args;
    static char message[ERROR_BUFFER_SIZE];
    int res;

    va_start(args, error);
    res = vsnprintf(message, ERROR_BUFFER_SIZE, (char *) error, args);
    go_error_print_callback(userData, message);
    va_end(args);
}

void
set_xslt_error_func(void *userData) {
    xsltSetGenericErrorFunc(userData, (xmlGenericErrorFunc) generic_error_func);
}

void
set_xslt_transform_error_func(xsltTransformContextPtr ctx, void *userData) {
    xsltSetTransformErrorFunc(ctx, userData, (xmlGenericErrorFunc) generic_error_func);
}
