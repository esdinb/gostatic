
#include <libxml/xmlerror.h>
#include <libxslt/transform.h>
#include <libxslt/xsltutils.h>

void structured_error_func(void *, xmlErrorPtr);
void generic_error_func(void *, const char *message, ...);

void set_xslt_error_func(void *userData);
void set_xslt_transform_error_func(xsltTransformContextPtr ctx, void *userData);
