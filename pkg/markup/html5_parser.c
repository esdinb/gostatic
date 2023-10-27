/*
 */
#include "html5_parser.h"

extern void parserError(const char *);

node_stack_t *
node_stack_create()
{
    node_stack_t *stack = lexbor_calloc(1, sizeof(node_stack_t));
    if (stack == NULL) {
        return NULL;
    }

    stack->current = -1;
    stack->capacity = DEFAULT_STACK_CAPACITY;
    stack->nodes = lexbor_calloc(DEFAULT_STACK_CAPACITY, sizeof(xmlNodePtr));

    return stack;
}

inline int
node_stack_push(node_stack_t *stack, xmlNodePtr node)
{
    int next = ++stack->current;
    if (next == stack->capacity) {
        if (next == MAX_STACK_CAPACITY) {
            parserError("max stack capacity exceeded");
        }
        xmlNodePtr *tmp = stack->nodes;
        int capacity = stack->capacity * 2;
        stack->nodes = lexbor_calloc(capacity, sizeof(xmlNodePtr));
        memcpy(stack->nodes, tmp, stack->capacity * sizeof(xmlNodePtr));
        stack->capacity = capacity;
        lexbor_free(tmp);
    }
    stack->nodes[next] = node;
    stack->current = next;

    return 0;
}

inline xmlNodePtr
node_stack_pop(node_stack_t *stack)
{
    if (stack->current == -1) {
        return NULL;
    }

    xmlNodePtr node = stack->nodes[stack->current];
    stack->current--;

    return node;
}

inline xmlNodePtr
node_stack_peek(node_stack_t *stack)
{
    if (stack->current == -1) {
        return NULL;
    }

    return stack->nodes[stack->current];
}

inline int
node_stack_is_empty(node_stack_t *stack)
{
    return stack->current == -1;
}

inline int
node_stack_size(node_stack_t *stack)
{
    return stack->current;
}

inline void
node_stack_destroy(node_stack_t *stack)
{
    lexbor_free(stack->nodes);
    lexbor_free(stack);
}

static lxb_html_token_t *
token_callback(lxb_html_tokenizer_t *tokenizer, lxb_html_token_t *token, void *ctx)
{
    xmlNodePtr node;
    xmlNodePtr text;
    lxb_char_t *name;
    const xmlChar *content;
    lxb_html_token_attr_t *attr;
    const lxb_char_t *attr_name;
    xmlAttrPtr attr_node;
    int is_void;

    lexbor_hash_t *tags = lxb_html_tokenizer_tags(tokenizer);

    html5_parser_context_t *parser_ctx = (html5_parser_context_t*) ctx;

    if (token->tag_id == LXB_TAG__END_OF_FILE) {
        return token;
    }

    if (token->tag_id == LXB_TAG__TEXT) {

        node = xmlNewDocTextLen(parser_ctx->document, (const xmlChar *) token->text_start, (int) (token->text_end - token->text_start));

        html5_parse_append_element(parser_ctx, node);
 
        return token;
    }

    if (token->tag_id == LXB_TAG__EM_COMMENT) {

        node = xmlNewComment(NULL);
        text = xmlNewDocTextLen(parser_ctx->document, (const xmlChar *) token->text_start, (int) (token->text_end - token->text_start));
        if (xmlAddChild(node, text) == NULL) {
            parserError("failed to add comment text");
        }

        html5_parse_append_element(parser_ctx, node);

        return token;
    }

    if (token->tag_id == LXB_TAG__EM_DOCTYPE) {

        if (!node_stack_is_empty(parser_ctx->stack)) {
            parserError("invalid state");
        }

        attr = token->attr_first;
        xmlChar *dtd_id;
        xmlChar *dtd_name;
        xmlDtdPtr dtd;

        dtd_name = xmlStrndup((const xmlChar *) attr->name_begin, (int) (attr->name_end - attr->name_begin));
       
        attr = attr->next;
        if (attr == NULL) {
            dtd = xmlNewDtd(parser_ctx->document, dtd_name, NULL, NULL);
        }
        else if (attr->name->attr_id == LXB_DOM_ATTR_PUBLIC) {
            dtd_id = xmlCharStrndup((const char *) attr->value, attr->value_size);
            dtd = xmlNewDtd(parser_ctx->document, dtd_name, dtd_id, NULL);
            xmlFree(dtd_id);
        }
        else if (attr->name->attr_id == LXB_DOM_ATTR_SYSTEM) {
            dtd_id = xmlCharStrndup((const char *) attr->value, attr->value_size);
            dtd = xmlNewDtd(parser_ctx->document, dtd_name, NULL, dtd_id);
            xmlFree(dtd_id);
        }
        else {
            dtd = xmlNewDtd(parser_ctx->document, dtd_name, NULL, NULL);
        }

        xmlDtdPtr res = (xmlDtdPtr) xmlAddChild((xmlNodePtr) parser_ctx->document, (xmlNodePtr) dtd);
        if (res == NULL) {
            parserError("failed to add DTD node");
        }

        xmlFree(dtd_name);

        return token;
    }

    is_void = lxb_html_tag_is_void(token->tag_id);
    name = (lxb_char_t *) lxb_tag_name_by_id(tags, token->tag_id, NULL);
    if (name == NULL) {
        parserError("failed to get token name");
    }

    if ((token->type & LXB_HTML_TOKEN_TYPE_CLOSE) && !is_void) {
        html5_parse_end_element(parser_ctx, token, name);
    }
    else {
        html5_parse_start_element(parser_ctx, token, name);
    }

    attr = token->attr_first;

    while (attr != NULL) {

        attr_name = lxb_html_token_attr_name(attr, NULL);
        node = node_stack_peek(parser_ctx->stack);
        attr_node = xmlNewNsProp(node, NULL, (xmlChar *) attr_name, NULL);

        if (attr->value_begin) {
            text = xmlNewDocTextLen(parser_ctx->document, (const xmlChar *) attr->value, (int) (attr->value_size));
            if (text == NULL) {
                parserError("failed to create attribute value node");
            }
            xmlAddChild((xmlNodePtr) attr_node, text);
        }

        attr = attr->next;
    }

    if (is_void) {
        html5_parse_end_element(parser_ctx, token, name);
    }
    else if (token->type & LXB_HTML_TOKEN_TYPE_CLOSE_SELF) {
        html5_parse_end_element(parser_ctx, token, name);
    }

    return token;
}

html5_parser_context_t *
html5_create_parser_context(xmlDocPtr doc, xmlNodePtr node)
{
    html5_parser_context_t *ctx;
    lxb_status_t status;

    ctx = lexbor_calloc(1, sizeof(html5_parser_context_t));
    if (ctx == NULL) {
        parserError("failed to create parser context");
    }

    ctx->tokenizer = lxb_html_tokenizer_create();

    status = lxb_html_tokenizer_init(ctx->tokenizer);
    if (status != LXB_STATUS_OK) {
        parserError("failed to create tokenizer object");
    }

    ctx->document = doc;
    ctx->node = node;

    ctx->stack = node_stack_create();
    if (ctx->stack == NULL) {
        parserError("failed to create context node stack");
    }

    lxb_html_tokenizer_callback_token_done_set(ctx->tokenizer, token_callback, ctx);

    status = lxb_html_tokenizer_begin(ctx->tokenizer);
    if (status != LXB_STATUS_OK) {
        parserError("failed to prepare tokenizer object for parsing");
    }

    return ctx;
}

void
html5_destroy_parser_context(html5_parser_context_t *ctx)
{
    node_stack_destroy(ctx->stack);
    lxb_html_tokenizer_destroy(ctx->tokenizer);
}

int
html5_parse_chunk(html5_parser_context_t *ctx, char *data, size_t len)
{
    lxb_status_t status;

    status = lxb_html_tokenizer_chunk(ctx->tokenizer, (lxb_char_t *) data, len);
    if (status != LXB_STATUS_OK) {
        parserError("failed to parse the html data");
    }
    return (int) len;
}

int
html5_parse_end_document(html5_parser_context_t *ctx)
{
    lxb_status_t status;

    status = lxb_html_tokenizer_end(ctx->tokenizer);
    if (status != LXB_STATUS_OK) {
        parserError("failed to end parsing html data");
    }

    return 0;
}

int
html5_parse_append_element(html5_parser_context_t *ctx, xmlNodePtr node)
{
    xmlNodePtr parent = (xmlNodePtr) node_stack_peek(ctx->stack);
    xmlNodePtr res = xmlAddChild(parent, node);
    if (res == NULL) {
        return -1;
    }
    return 0;
}

int
html5_parse_start_element(html5_parser_context_t *ctx, lxb_html_token_t *token, lxb_char_t *name)
{
    xmlNodePtr node = xmlNewDocNode(ctx->document, NULL, (const xmlChar *) name, NULL);
    if (node == NULL) {
        parserError("failed to create new node");
    }

    if (node_stack_is_empty(ctx->stack)) {
        if (ctx->node == NULL) {
            xmlDocSetRootElement(ctx->document, node);
        }
        else {
            node_stack_push(ctx->stack, ctx->node);
        }
    }

    node_stack_push(ctx->stack, node);

    return 0;
}

int
html5_parse_end_element(html5_parser_context_t *ctx, lxb_html_token_t *token, lxb_char_t *name)
{
    xmlNodePtr node;
   
    node  = node_stack_pop(ctx->stack);
    if (node == NULL) {
        parserError("unmatched element");
    }
    
    if (xmlStrcmp(node->name, (const xmlChar *) name)) {
        parserError("start and end element name mismatch");
    }
 
    if (node_stack_is_empty(ctx->stack)) {
        return 0;
    }

    int res = html5_parse_append_element(ctx, node);
    if (res < 0) {
        parserError("unexpected end of document");
    }

    return 0;
}

