#include <libxml/xmlwriter.h>
#include <libxml/tree.h>
#include <lexbor/html/tokenizer.h>
#include <lexbor/html/tag.h>

#define DEFAULT_STACK_CAPACITY 32
#define MAX_STACK_CAPACITY 256

#define PARSE_OK 0
#define PARSE_MALFORMED_INPUT 1
#define PARSE_MISMATCHED_ELEMENTS 2

struct node_stack {
    int current;
    int capacity;
    xmlNodePtr *nodes;
};

typedef struct node_stack node_stack_t;

#define STACK_APPEND(stack, node) (xmlAddChild(STACK_PEEK(stack), node) == NULL ? -1 : 0)
#define STACK_CURRENT(stack) (stack->current)
#define STACK_POP(stack) (STACK_EMPTY(stack) ? NULL : stack->nodes[stack->current--])
#define STACK_PEEK(stack) (stack->current == -1 ? NULL : stack->nodes[stack->current])
#define STACK_EMPTY(stack) (stack->current == -1)

struct html5_parser_context {
    xmlDocPtr document;
    xmlNodePtr node;
    node_stack_t *stack;
    lxb_html_tokenizer_t *tokenizer;
};

typedef struct html5_parser_context html5_parser_context_t;

node_stack_t *node_stack_create();
int node_stack_push(node_stack_t *, xmlNodePtr);
xmlNodePtr node_stack_pop(node_stack_t *);
xmlNodePtr node_stack_peek(node_stack_t *);
int node_stack_is_empty(node_stack_t *);
int node_stack_size(node_stack_t *);
void node_stack_destroy(node_stack_t *);

html5_parser_context_t *html5_create_parser_context(xmlDocPtr doc, xmlNodePtr node);
void html5_destroy_parser_context(html5_parser_context_t *);
int html5_parse_chunk(html5_parser_context_t *, char *, size_t);
int html5_parse_end_document(html5_parser_context_t *);

int html5_parse_append_element(html5_parser_context_t *, xmlNodePtr);
int html5_parse_start_element(html5_parser_context_t *, lxb_html_token_t *, lxb_char_t *);
int html5_parse_end_element(html5_parser_context_t *, lxb_html_token_t *, lxb_char_t *);

