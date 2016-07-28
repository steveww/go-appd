package appd

/*
#cgo CFLAGS: -I ${APPD_SDK_HOME}
#cgo LDFLAGS: -L ${APPD_SDK_HOME}/lib -lappdynamics_native_sdk
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <appdynamics.h>

// global configuration
struct appd_config cfg;

void dump_config(struct appd_config *cfg) {
    printf("app name %s\n", cfg->app_name);
    printf("tier name %s\n", cfg->tier_name);
    printf("node name %s\n", cfg->node_name);
    printf("\n");
    printf("controller host %s\n", cfg->controller.host);
    printf("controller port %d\n", cfg->controller.port);
    printf("controller account %s\n", cfg->controller.account);
    printf("controller access key %s\n", cfg->controller.access_key);
    printf("controller use ssl %c\n", cfg->controller.use_ssl);
    printf("\n");
    printf("init timeout ms %d\n", cfg->init_timeout_ms);
    printf("\n");
    printf("http proxy host %s\n", cfg->controller.http_proxy.host);
    printf("http proxy port %d\n", cfg->controller.http_proxy.port);
    printf("http proxy username %s\n", cfg->controller.http_proxy.username);
    printf("http proxy password file %s\n", cfg->controller.http_proxy.password_file);
}

void printpair(char *name, char *value) {
    printf("%s = %s\n", name, value);
}

uintptr_t bt_handle_to_int(appd_bt_handle bt_handle) {
    return (uintptr_t)bt_handle;
}
appd_bt_handle bt_int_to_handle(uintptr_t bt) {
    return (appd_bt_handle)bt;
}

uintptr_t exit_handle_to_int(appd_exitcall_handle exitcall_handle) {
    return (uintptr_t)exitcall_handle;
}
appd_bt_handle exit_int_to_handle(uintptr_t exit) {
    return (appd_exitcall_handle)exit;
}
*/
import "C"

import (
    "log"
    "unsafe"
    "net/http"
    "strconv"
)

/*
* Some constants from appdynamics.h
*/
const CORRELATION_HEADER_NAME = C.APPD_CORRELATION_HEADER_NAME

const BACKEND_HTTP = C.APPD_BACKEND_HTTP
const BACKEND_DB = C.APPD_BACKEND_DB
const BACKEND_CACHE = C.APPD_BACKEND_CACHE
const BACKEND_RABBITMQ = C.APPD_BACKEND_RABBITMQ
const BACKEND_WEBSERVICE = C.APPD_BACKEND_WEBSERVICE

const ERROR_LEVEL_NOTICE = C.APPD_LEVEL_NOTICE
const ERROR_LEVEL_WARNING = C.APPD_LEVEL_WARNING
const ERROR_LEVEL_ERROR = C.APPD_LEVEL_ERROR

const APPD_BT = "APPD_BT"

type ID_properties_map map[string]string

// do not free the C.CString as the config struct is only pointers
func Init(appName string, controllerKey string) {
    log.Println("Init called")

    appName_c := C.CString(appName)
    controllerKey_c := C.CString(controllerKey)

    C.appd_config_init(&C.cfg)
    C.cfg.app_name = appName_c
    C.cfg.controller.access_key = controllerKey_c
}

func SetTierName(name string) {
    C.cfg.tier_name = C.CString(name)
}

func SetNodeName(name string) {
    C.cfg.node_name = C.CString(name)
}

func SetControllerHost(host string) {
    C.cfg.controller.host = C.CString(host)
}

func SetControllerPort(port int16) {
    C.cfg.controller.port = C.ushort(port)
}

func SetControllerAccount(account string) {
    C.cfg.controller.account = C.CString(account)
}

func SetControllerUseSSL(ssl byte) {
    C.cfg.controller.use_ssl = C.char(ssl)
}

// Proxy stuff to be added later
// Windows stuff to be added later

func SetInitTimeout(timeout int) {
    C.cfg.init_timeout_ms = C.int(timeout)
}

func Sdk_init() int {
    //C.dump_config(&C.cfg)
    rc := C.appd_sdk_init(&C.cfg)
    return int(rc)
}

func Sdk_term() {
    C.appd_sdk_term()
}

/*
* BT
*/
func BT_begin(name string, correlation string) uint64 {
    name_c := C.CString(name)
    correlation_c := C.CString(correlation)
    defer C.free(unsafe.Pointer(name_c))
    defer C.free(unsafe.Pointer(correlation_c))

    bt := C.appd_bt_begin(name_c, correlation_c)

    return uint64(C.bt_handle_to_int(bt))
}

func BT_end(bt uint64) {
    C.appd_bt_end(C.bt_int_to_handle(C.uintptr_t(bt)))
}

func BT_set_url(bt uint64, name string) {
    name_c := C.CString(name)
    defer C.free(unsafe.Pointer(name_c))

    C.appd_bt_set_url(C.bt_int_to_handle(C.uintptr_t(bt)), name_c)
}

func BT_is_snapshotting(bt uint64) int {
    result := C.appd_bt_is_snapshotting(C.bt_int_to_handle(C.uintptr_t(bt)))
    return int(result)
}

func BT_add_user_data(bt uint64, key string, value string) {
    key_c := C.CString(key)
    value_c := C.CString(value)
    defer C.free(unsafe.Pointer(key_c))
    defer C.free(unsafe.Pointer(value_c))

    C.appd_bt_add_user_data(C.bt_int_to_handle(C.uintptr_t(bt)), key_c, value_c)
}

func BT_add_error(bt uint64, level uint32, message string, mark_bt_as_error int) {
    message_c := C.CString(message)
    defer C.free(unsafe.Pointer(message_c))

    C.appd_bt_add_error(C.bt_int_to_handle(C.uintptr_t(bt)), level, message_c, C.int(mark_bt_as_error))
}

/*
* Backend
*/
func Backend_declare(betype string, name string) {
    betype_c := C.CString(betype)
    name_c := C.CString(name)
    defer C.free(unsafe.Pointer(betype_c))
    defer C.free(unsafe.Pointer(name_c))

    log.Println("Backend decalre", betype, name)
    C.appd_backend_declare(C.CString(betype), C.CString(name))
}

func Backend_set_identifying_property(name string, key string, value string) int {
    name_c := C.CString(name)
    key_c := C.CString(key)
    value_c := C.CString(value)
    defer C.free(unsafe.Pointer(name_c))
    defer C.free(unsafe.Pointer(key_c))
    defer C.free(unsafe.Pointer(value_c))

    rc := C.appd_backend_set_identifying_property(C.CString(name), C.CString(key), C.CString(value))
    return int(rc)
}

func Backend_set_identifying_properties(name string, props ID_properties_map) int {
    var rc C.int
    name_c := C.CString(name)
    defer C.free(unsafe.Pointer(name_c))

    for key, value := range props {
        key_c := C.CString(key)
        value_c := C.CString(value)
        rc = C.appd_backend_set_identifying_property(name_c, key_c, value_c)
        C.free(unsafe.Pointer(key_c))
        C.free(unsafe.Pointer(value_c))
        if(rc != 0) {
            break
        }
    }
    return int(rc)
}

func Backend_add(name string) int {
    name_c := C.CString(name)
    defer C.free(unsafe.Pointer(name_c))

    rc := C.appd_backend_add(name_c)
    return int(rc)
}

func Backend_prevent_agent_resolution(name string) int {
    name_c := C.CString(name)
    defer C.free(unsafe.Pointer(name_c))

    rc := C.appd_backend_prevent_agent_resolution(name_c)
    return int(rc)
}

/*
* Exit
*/
func Exitcall_begin(bt uint64, name string) uint64 {
    name_c := C.CString(name)
    defer C.free(unsafe.Pointer(name_c))

    bt_h := C.bt_int_to_handle(C.uintptr_t(bt))
    exit := C.appd_exitcall_begin(bt_h, name_c)

    return uint64(C.exit_handle_to_int(exit))
}

func Exitcall_end(exit uint64) {
    exit_h := C.exit_int_to_handle(C.uintptr_t(exit))
    C.appd_exitcall_end(exit_h)
}

func Exitcall_set_details(exit uint64, details string) int {
    details_c := C.CString(details)
    defer C.free(unsafe.Pointer(details_c))

    rc := C.appd_exitcall_set_details(C.exit_int_to_handle(C.uintptr_t(exit)), details_c)
    return int(rc)
}

func Exitcall_add_error(exit uint64, error_level uint32, message string, mark_bt_as_error int) {
    message_c := C.CString(message)
    defer C.free(unsafe.Pointer(message_c))

    C.appd_exitcall_add_error(C.exit_int_to_handle(C.uintptr_t(exit)), error_level, message_c, C.int(mark_bt_as_error))
}

func Exitcall_get_correlation_header(exit uint64) string {
    header := C.appd_exitcall_get_correlation_header(C.exit_int_to_handle(C.uintptr_t(exit)))
    return C.GoString(header)
}

/*
* HTTP wrappers
*/
/*
* Wrappers for http handlers
* name is the name used for the BT
*/
func WrapHandle(name string, pattern string, handler http.Handler) (string, http.Handler) {
    return pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // get correlation id from header
        appd_correlation := r.Header.Get(CORRELATION_HEADER_NAME)
        bt := BT_begin(name, appd_correlation)
        defer BT_end(bt)
        BT_set_url(bt, r.URL.String())
        r.Header.Add(APPD_BT, strconv.FormatUint(bt, 10))
        handler.ServeHTTP(w, r)
    })
}

func WrapHandleFunc(name string, pattern string, handler func(http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request)) {
    p, h := WrapHandle(name, pattern, http.HandlerFunc(handler))
    return p, func(w http.ResponseWriter, r *http.Request) {
        h.ServeHTTP(w, r)
    }
}

