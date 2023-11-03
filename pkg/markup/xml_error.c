
#include <libxml/xmlerror.h>
#include "xml_error.h"

extern void go_error_callback(void *, xmlErrorPtr);

void custom_error_func(void *userData, xmlErrorPtr error) {
    go_error_callback(userData, error);
}
