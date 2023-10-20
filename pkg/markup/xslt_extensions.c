
#include "xslt_extensions.h"

extern void FormatDateCallback(xmlXPathParserContextPtr ctx, int nargs);

void registerExtensionFunctions(xsltTransformContextPtr ctx) {
	xsltRegisterExtFunction(ctx, BAD_CAST "format-date", (xmlChar *)GOSTATIC_NAMESPACE, &FormatDateCallback);
}



