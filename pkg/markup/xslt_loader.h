
#include <libxslt/documents.h>
#include <libxml/tree.h>

void save_default_loader();

xmlDocPtr custom_loader(const xmlChar * URI, xmlDictPtr dict, int options, void * ctxt, xsltLoadType type);

xmlDocPtr default_loader(const xmlChar * URI, xmlDictPtr dict, int options, void * ctxt, xsltLoadType type);

