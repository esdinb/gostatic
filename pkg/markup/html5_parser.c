/*
 */
#include "html5_parser.h"

static lxb_html_token_t *
token_callback(lxb_html_tokenizer_t *tokenizer, lxb_html_token_t *token, void *ctx)
{
    const lxb_char_t *name;
    const lxb_char_t *attr_name;
    lexbor_hash_t *tags = lxb_html_tokenizer_tags(tokenizer);
    xmlTextWriterPtr writer = (xmlTextWriterPtr) ctx;

    lxb_html_token_attr_t *attr;

    if (token->tag_id == LXB_TAG__END_OF_FILE) {
        return token;
    }

    if (token->tag_id == LXB_TAG__TEXT) {
        xmlTextWriterWriteFormatString(writer, "%.*s", (int) (token->text_end - token->text_start), token->text_start);

        return token;
    }

    if (token->tag_id == LXB_TAG__EM_COMMENT) {
        xmlTextWriterWriteFormatRaw(writer, "<!%.*s>", (int) (token->text_end - token->text_start), token->text_start);
        return token;
    }

    name = lxb_tag_name_by_id(tags, token->tag_id, NULL);
    if (name == NULL) {
        FAILED("Failed to get token name");
    }

    if (token->type & LXB_HTML_TOKEN_TYPE_CLOSE) {
        xmlTextWriterEndElement(writer);
    }
    else {
        xmlTextWriterStartElement(writer, (xmlChar *) name);
    }

    attr = token->attr_first;

    while (attr != NULL) {

        attr_name = lxb_html_token_attr_name(attr, NULL);
        xmlTextWriterStartAttribute(writer, attr_name);

        if (attr->value_begin) {
            xmlTextWriterWriteFormatRaw(writer, "%.*s", (int) (attr->value_size), attr->value);
        }

        xmlTextWriterEndAttribute(writer);

        attr = attr->next;
    }

    if (lxb_html_tag_is_void(token->tag_id)) {
        xmlTextWriterEndElement(writer);
    }

    if (token->type & LXB_HTML_TOKEN_TYPE_CLOSE_SELF) {
        xmlTextWriterEndElement(writer);
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
        FAILED("Failed to create parser context");
    }

    ctx->tokenizer = lxb_html_tokenizer_create();

    status = lxb_html_tokenizer_init(ctx->tokenizer);
    if (status != LXB_STATUS_OK) {
        FAILED("Failed to create tokenizer object");
    }

    ctx->writer = xmlNewTextWriterTree(doc, node, 0);
    if (ctx->writer == NULL) {
        FAILED("Failed to create writer object");
    }

    lxb_html_tokenizer_callback_token_done_set(ctx->tokenizer, token_callback, ctx->writer);

    status = lxb_html_tokenizer_begin(ctx->tokenizer);
    if (status != LXB_STATUS_OK) {
        FAILED("Failed to prepare tokenizer object for parsing");
    }

    return ctx;
}

void
html5_destroy_parser_context(html5_parser_context_t *ctx)
{
    lxb_html_tokenizer_t *tokenizer;

    lxb_html_tokenizer_destroy(ctx->tokenizer);
}

int
html5_parse_chunk(html5_parser_context_t *ctx, char *data, size_t len)
{
    lxb_status_t status;

    status = lxb_html_tokenizer_chunk(ctx->tokenizer, (lxb_char_t *) data, len);
    if (status != LXB_STATUS_OK) {
        FAILED("Failed to parse the html data");
    }
    return (int) len;
}

int
html5_parse_end(html5_parser_context_t *ctx)
{
    lxb_status_t status;
    int result;

    status = lxb_html_tokenizer_end(ctx->tokenizer);
    if (status != LXB_STATUS_OK) {
        FAILED("Failed to ending of parsing the html data");
    }

    result = xmlTextWriterFlush(ctx->writer);
    if (result < 0) {
        FAILED("Failed to end document");
    }

    return 0;
}

