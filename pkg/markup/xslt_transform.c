
// https://groups.google.com/g/golang-nuts/c/pQueMFdY0mk/m/OAX5-Fqus0UJ

#include "xslt_transform.h"

char **makeParamsArray(int size) {
    return calloc(sizeof(char *), size);
}

void freeParamsArray(char **list, int size) {
    int i;
    for (i = 0; i < size; i++) {
        free(list[i]);
    }
    free(list);
}

void setParamsElement(char **list, char *value, int idx) {
    list[idx] = value;
}

