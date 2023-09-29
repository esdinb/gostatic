
#include "xslt_loader.h"

extern xmlDocPtr go_loader_callback(const xmlChar * URI, xmlDictPtr dict, int options, void * ctxt, xsltLoadType type);

static xsltDocLoaderFunc xslt_loader;

void save_default_loader() {
    xslt_loader = xsltDocDefaultLoader;
}

xmlDocPtr custom_loader(const xmlChar * URI, xmlDictPtr dict, int options, void * ctxt, xsltLoadType type) {
    return go_loader_callback(URI, dict, options, ctxt, type);
}

xmlDocPtr default_loader(const xmlChar * URI, xmlDictPtr dict, int options, void * ctxt, xsltLoadType type) {
    return xslt_loader(URI, dict, options, ctxt, type);
}

