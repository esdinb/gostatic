#include <libxml/xmlwriter.h>
#include <libxml/tree.h>
#include <lexbor/html/tokenizer.h>

#define FAILED(...)                                                            \
    do {                                                                       \
        fprintf(stderr, __VA_ARGS__);                                          \
        fprintf(stderr, "\n");                                                 \
        exit(EXIT_FAILURE);                                                    \
    }                                                                          \
    while (0)

struct html5_parser_context {
    xmlTextWriterPtr writer;
    lxb_html_tokenizer_t *tokenizer;
};

typedef struct html5_parser_context html5_parser_context_t;

html5_parser_context_t *html5_create_parser_context(xmlDocPtr doc, xmlNodePtr node);
void html5_destroy_parser_context(html5_parser_context_t *);
int html5_parse_chunk(html5_parser_context_t *, char *, size_t);
int html5_parse_end(html5_parser_context_t *);


