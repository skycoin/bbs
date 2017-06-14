webpackJsonp([1],{

/***/ "./src async recursive":
/***/ (function(module, exports) {

function webpackEmptyContext(req) {
	throw new Error("Cannot find module '" + req + "'.");
}
webpackEmptyContext.keys = function() { return []; };
webpackEmptyContext.resolve = webpackEmptyContext;
module.exports = webpackEmptyContext;
webpackEmptyContext.id = "./src async recursive";

/***/ }),

/***/ "./src/animations/router.animations.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("./node_modules/@angular/core/@angular/core.es5.js");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return slideInLeftAnimation; });

// Component transition animations
var slideInLeftAnimation = __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_12" /* trigger */])('routeAnimation', [
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_13" /* state */])('*', __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_14" /* style */])({
        opacity: 1,
        transform: 'translateX(0)'
    })),
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_15" /* transition */])(':enter', [
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_14" /* style */])({
            opacity: 0,
            transform: 'translateX(-100%)'
        }),
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_16" /* animate */])('0.3s ease-in')
    ]),
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_15" /* transition */])(':leave', [
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_16" /* animate */])('0.3s ease-out', __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_14" /* style */])({
            opacity: 0,
            transform: 'translateX(100%)'
        }))
    ])
]);
//# sourceMappingURL=router.animations.js.map

/***/ }),

/***/ "./src/app/app.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("./node_modules/css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, ".nav-name {\n  color: #d21b1b !important;\n}\n\n.nav-node {\n  color: blue !important;\n  padding: 0;\n}\n\n.nav-name span {\n  margin: 0 5px;\n}\n\n.app-pop {\n  position: fixed;\n  bottom: 0;\n  width: 50%;\n  left: 25%;\n  right: 25%;\n}\n\n.top-btn {\n  position: fixed;\n  z-index: 99;\n  bottom: 2%;\n  right: 2%;\n  cursor: pointer;\n}\n\n.top-btn .top {\n  width: 0;\n  height: 0;\n  border-left: .5em solid transparent;\n  border-right: .5em solid transparent;\n  border-bottom: 1em solid #fff;\n  margin: 0 auto;\n  border-bottom: 1em solid #fff;\n}", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "./src/app/app.component.html":
/***/ (function(module, exports) {

module.exports = "<nav class=\"navbar fixed-top navbar-toggleable-md navbar-light bg-faded\">\n  <button class=\"navbar-toggler navbar-toggler-right\" type=\"button\" data-toggle=\"collapse\" data-target=\"#navbarSupportedContent\"\n    aria-controls=\"navbarSupportedContent\" aria-expanded=\"false\" aria-label=\"Toggle navigation\">\n    <span class=\"navbar-toggler-icon\"></span>\n  </button>\n  <a class=\"navbar-brand\" href=\"#\">BBS</a>\n  <div class=\"collapse navbar-collapse\" id=\"navbarSupportedContent\">\n    <ul class=\"navbar-nav mr-auto\">\n      <li class=\"nav-item\">\n        <a class=\"nav-link\" routerLink=\"/\" routerLinkActive=\"active\">Board <span class=\"sr-only\">(current)</span></a>\n      </li>\n      <li class=\"nav-item\">\n        <a class=\"nav-link\" routerLink=\"/userlist\" routerLinkActive=\"active\">UserList</a>\n      </li>\n      <li class=\"nav-item\">\n        <a class=\"nav-link\" routerLink=\"/conn\" routerLinkActive=\"active\">Connections Manager</a>\n      </li>\n      \n      <li class=\"nav-item\">\n        <a class=\"nav-name nav-link\" routerLink=\"/user\" routerLinkActive=\"active\"><i class=\"fa fa-user\" aria-hidden=\"true\"></i>{{name}}</a>\n      </li>\n    </ul>\n    <span class=\"navbar-text\">\n     <a class=\"nav-node nav-link\" href=\"javascript:void(0);\" *ngIf=\"isMasterNode\">Master Node</a>\n     <a class=\"nav-node nav-link\" href=\"javascript:void(0);\" *ngIf=\"!isMasterNode\">Client Node</a>\n    </span>\n  </div>\n</nav>\n\n<router-outlet></router-outlet>\n<button class=\"top-btn btn btn-primary \" [hidden]=\"!common.topBtn\" (click)=\"common.scrollToTop()\">\n  <div class=\"top\"></div>\n  <span>top</span>\n</button>\n<div class=\"app-pop\">\n  <ngb-alert type=\"{{common.alertType}}\" *ngIf=\"common.alert\" (click)=\"common.alert = false\">\n    {{common.alertMessage}}\n  </ngb-alert>\n</div>\n"

/***/ }),

/***/ "./src/app/app.component.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("./node_modules/@angular/core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__providers__ = __webpack_require__("./src/providers/index.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__components__ = __webpack_require__("./src/components/index.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__angular_router__ = __webpack_require__("./node_modules/@angular/router/@angular/router.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4_rxjs_add_operator_filter__ = __webpack_require__("./node_modules/rxjs/add/operator/filter.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4_rxjs_add_operator_filter___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_4_rxjs_add_operator_filter__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return AppComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};






var AppComponent = (function () {
    function AppComponent(api, user, router, common) {
        this.api = api;
        this.user = user;
        this.router = router;
        this.common = common;
        this.title = 'app';
        this.name = '';
        this.isMasterNode = false;
    }
    AppComponent.prototype.ngOnInit = function () {
        var _this = this;
        this.user.getCurrent().subscribe(function (user) {
            _this.name = user.alias;
        });
        this.api.getStats().subscribe(function (res) {
            _this.isMasterNode = res.node_is_master;
        });
        this.router.events.filter(function (ev) { return ev instanceof __WEBPACK_IMPORTED_MODULE_3__angular_router__["c" /* NavigationStart */]; }).subscribe(function (ev) {
            _this.common.topBtn = false;
        });
    };
    AppComponent.prototype.test = function () {
        console.log('test');
    };
    return AppComponent;
}());
__decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_24" /* ViewChild */])(__WEBPACK_IMPORTED_MODULE_2__components__["a" /* BoardsListComponent */]),
    __metadata("design:type", typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_2__components__["a" /* BoardsListComponent */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__components__["a" /* BoardsListComponent */]) === "function" && _a || Object)
], AppComponent.prototype, "boards", void 0);
__decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_24" /* ViewChild */])(__WEBPACK_IMPORTED_MODULE_2__components__["b" /* ThreadsComponent */]),
    __metadata("design:type", typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_2__components__["b" /* ThreadsComponent */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__components__["b" /* ThreadsComponent */]) === "function" && _b || Object)
], AppComponent.prototype, "threads", void 0);
__decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_24" /* ViewChild */])(__WEBPACK_IMPORTED_MODULE_2__components__["c" /* ThreadPageComponent */]),
    __metadata("design:type", typeof (_c = typeof __WEBPACK_IMPORTED_MODULE_2__components__["c" /* ThreadPageComponent */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__components__["c" /* ThreadPageComponent */]) === "function" && _c || Object)
], AppComponent.prototype, "threadPage", void 0);
AppComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_0" /* Component */])({
        selector: 'app-root',
        template: __webpack_require__("./src/app/app.component.html"),
        styles: [__webpack_require__("./src/app/app.component.css")]
    }),
    __metadata("design:paramtypes", [typeof (_d = typeof __WEBPACK_IMPORTED_MODULE_1__providers__["ApiService"] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__providers__["ApiService"]) === "function" && _d || Object, typeof (_e = typeof __WEBPACK_IMPORTED_MODULE_1__providers__["UserService"] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__providers__["UserService"]) === "function" && _e || Object, typeof (_f = typeof __WEBPACK_IMPORTED_MODULE_3__angular_router__["a" /* Router */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_3__angular_router__["a" /* Router */]) === "function" && _f || Object, typeof (_g = typeof __WEBPACK_IMPORTED_MODULE_1__providers__["CommonService"] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__providers__["CommonService"]) === "function" && _g || Object])
], AppComponent);

var _a, _b, _c, _d, _e, _f, _g;
//# sourceMappingURL=app.component.js.map

/***/ }),

/***/ "./src/app/app.module.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_platform_browser__ = __webpack_require__("./node_modules/@angular/platform-browser/@angular/platform-browser.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_core__ = __webpack_require__("./node_modules/@angular/core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_http__ = __webpack_require__("./node_modules/@angular/http/@angular/http.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__angular_forms__ = __webpack_require__("./node_modules/@angular/forms/@angular/forms.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__ng_bootstrap_ng_bootstrap__ = __webpack_require__("./node_modules/@ng-bootstrap/ng-bootstrap/index.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5__app_component__ = __webpack_require__("./src/app/app.component.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_6__providers__ = __webpack_require__("./src/providers/index.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_7__router_app_router_routing_module__ = __webpack_require__("./src/router/app-router-routing.module.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_8_angular2_froala_wysiwyg__ = __webpack_require__("./node_modules/angular2-froala-wysiwyg/index.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_9__angular_platform_browser_animations__ = __webpack_require__("./node_modules/@angular/platform-browser/@angular/platform-browser/animations.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_10__components__ = __webpack_require__("./src/components/index.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_11__pipes__ = __webpack_require__("./src/pipes/index.ts");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return AppModule; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};












var AppModule = (function () {
    function AppModule() {
    }
    return AppModule;
}());
AppModule = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_1__angular_core__["b" /* NgModule */])({
        imports: [
            __WEBPACK_IMPORTED_MODULE_0__angular_platform_browser__["a" /* BrowserModule */],
            __WEBPACK_IMPORTED_MODULE_9__angular_platform_browser_animations__["a" /* BrowserAnimationsModule */],
            __WEBPACK_IMPORTED_MODULE_2__angular_http__["a" /* HttpModule */],
            __WEBPACK_IMPORTED_MODULE_3__angular_forms__["a" /* FormsModule */],
            __WEBPACK_IMPORTED_MODULE_3__angular_forms__["b" /* ReactiveFormsModule */],
            __WEBPACK_IMPORTED_MODULE_7__router_app_router_routing_module__["a" /* AppRouterRoutingModule */],
            __WEBPACK_IMPORTED_MODULE_4__ng_bootstrap_ng_bootstrap__["a" /* NgbModule */].forRoot(),
            __WEBPACK_IMPORTED_MODULE_8_angular2_froala_wysiwyg__["a" /* FroalaEditorModule */].forRoot(),
            __WEBPACK_IMPORTED_MODULE_8_angular2_froala_wysiwyg__["b" /* FroalaViewModule */].forRoot()
        ],
        declarations: [
            __WEBPACK_IMPORTED_MODULE_5__app_component__["a" /* AppComponent */],
            __WEBPACK_IMPORTED_MODULE_10__components__["a" /* BoardsListComponent */],
            __WEBPACK_IMPORTED_MODULE_10__components__["b" /* ThreadsComponent */],
            __WEBPACK_IMPORTED_MODULE_10__components__["c" /* ThreadPageComponent */],
            __WEBPACK_IMPORTED_MODULE_10__components__["d" /* AddComponent */],
            __WEBPACK_IMPORTED_MODULE_10__components__["e" /* UserlistComponent */],
            __WEBPACK_IMPORTED_MODULE_10__components__["f" /* UserComponent */],
            __WEBPACK_IMPORTED_MODULE_10__components__["g" /* ConnectionComponent */],
            __WEBPACK_IMPORTED_MODULE_10__components__["h" /* AlertComponent */],
            __WEBPACK_IMPORTED_MODULE_11__pipes__["a" /* SafeHTMLPipe */],
        ],
        entryComponents: [__WEBPACK_IMPORTED_MODULE_10__components__["h" /* AlertComponent */]],
        providers: [__WEBPACK_IMPORTED_MODULE_6__providers__["CommonService"], __WEBPACK_IMPORTED_MODULE_6__providers__["ApiService"], __WEBPACK_IMPORTED_MODULE_6__providers__["UserService"], __WEBPACK_IMPORTED_MODULE_6__providers__["ConnectionService"]],
        bootstrap: [__WEBPACK_IMPORTED_MODULE_5__app_component__["a" /* AppComponent */]]
    })
], AppModule);

//# sourceMappingURL=app.module.js.map

/***/ }),

/***/ "./src/components/add/add.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("./node_modules/css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, ".box {\n  width: 80%;\n  min-height: 40%;\n  margin: 10% auto;\n}", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "./src/components/add/add.component.html":
/***/ (function(module, exports) {

module.exports = "<div class=\"container\">\n  <div class=\"page-header\">\n    <h1>Add</h1>\n  </div>\n  <div class=\"form-group\">\n    <label for=\"type\">Select add</label>\n    <select class=\"form-control\" id=\"type\" [ngModel]=\"select\" (ngModelChange)=\"clear($event)\">\n    <option value=\"board\">Add Board</option>\n      <option value=\"thread\">Add Thread</option>\n      <option value=\"post\">Add Post</option>\n      <option value=\"changeBoard\">ChangeBoard</option>\n  </select>\n  </div>\n  <!--<form>-->\n  <div class=\"form-group\" [hidden]=\"select != 'board' && select != 'thread'\">\n    <label for=\"name\">Input board name</label>\n    <input type=\"text\" class=\"form-control\" placeholder=\"name\" id=\"name\" [(ngModel)]=\"form.name\">\n  </div>\n  <div class=\"form-group\" [hidden]=\"select != 'thread' && select != 'post'\">\n    <label for=\"board\">Input board key</label>\n    <input type=\"text\" class=\"form-control\" placeholder=\"board\" id=\"board\" [(ngModel)]=\"form.board\">\n  </div>\n  <div class=\"form-group\" [hidden]=\"select != 'post' && select != 'changeBoard'\">\n    <label for=\"thread\">Input thread key</label>\n    <input type=\"text\" class=\"form-control\" placeholder=\"thread\" id=\"thread\" [(ngModel)]=\"form.thread\">\n  </div>\n  <div class=\"form-group\" [hidden]=\"select != 'post'\">\n    <label for=\"title\">Input title</label>\n    <input type=\"text\" class=\"form-control\" placeholder=\"title\" id=\"title\" [(ngModel)]=\"form.title\">\n  </div>\n  <div class=\"form-group\" [hidden]=\"select != 'post'\">\n    <label for=\"body\">Input Post Body</label>\n    <textarea class=\"form-control\" rows=\"3\" id=\"body\" [(ngModel)]=\"form.body\"></textarea>\n  </div>\n  <div class=\"form-group\" [hidden]=\"select != 'board' && select != 'thread'\">\n    <label for=\"description\">Input description</label>\n    <input type=\"text\" class=\"form-control\" placeholder=\"description\" id=\"description\" [(ngModel)]=\"form.description\">\n  </div>\n  <div class=\"form-group\" [hidden]=\"select != 'board'\">\n    <label for=\"seed\">Input seed</label>\n    <input type=\"text\" class=\"form-control\" placeholder=\"seed\" id=\"seed\" [(ngModel)]=\"form.seed\">\n  </div>\n  <div class=\"form-group\" [hidden]=\"select != 'changeBoard'\">\n    <label for=\"fromBoard\">From Board</label>\n    <input type=\"text\" class=\"form-control\" placeholder=\"fromBoard\" id=\"fromBoard\" [(ngModel)]=\"form.fromBoard\">\n  </div>\n  <div class=\"form-group\" [hidden]=\"select != 'changeBoard'\">\n    <label for=\"toBoard\">To Board</label>\n    <input type=\"text\" class=\"form-control\" placeholder=\"toBoard\" id=\"toBoard\" [(ngModel)]=\"form.toBoard\">\n  </div>\n  <button class=\"btn btn-info\" (click)=\"add($event)\">Submit</button>\n  <!--</form>-->\n\n</div>\n"

/***/ }),

/***/ "./src/components/add/add.component.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("./node_modules/@angular/core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__providers__ = __webpack_require__("./src/providers/index.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_router__ = __webpack_require__("./node_modules/@angular/router/@angular/router.es5.js");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return AddComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};



var AddComponent = (function () {
    function AddComponent(api, router, route) {
        this.api = api;
        this.router = router;
        this.route = route;
        this.select = 'board';
        this.form = {
            name: '',
            description: '',
            board: '',
            thread: '',
            seed: '',
            title: '',
            body: '',
            fromBoard: '',
            toBoard: ''
        };
    }
    AddComponent.prototype.ngOnInit = function () {
        var _this = this;
        this.route.params.subscribe(function (data) {
            if (data['exec']) {
                _this.select = data['exec'];
            }
            _this.form.board = data['board'];
            _this.form.thread = data['thread'];
        });
    };
    AddComponent.prototype.init = function () {
        this.form = {
            name: '',
            description: '',
            board: '',
            thread: '',
            seed: '',
            title: '',
            body: '',
            fromBoard: '',
            toBoard: ''
        };
    };
    AddComponent.prototype.clear = function (ev) {
        this.select = ev;
        this.init();
    };
    AddComponent.prototype.add = function (ev) {
        var _this = this;
        ev.stopImmediatePropagation();
        ev.stopPropagation();
        var data = new FormData();
        // console.log('form:', this.form);
        switch (this.select) {
            case 'board':
                data.append('name', this.form.name);
                data.append('description', this.form.description);
                data.append('seed', this.form.seed);
                this.api.addBoard(data).subscribe(function (res) {
                    alert('add success');
                    _this.init();
                });
                break;
            case 'thread':
                data.append('board', this.form.board);
                data.append('description', this.form.description);
                data.append('name', this.form.name);
                this.api.addThread(data).subscribe(function (res) {
                    alert('add success');
                    _this.init();
                });
                break;
            case 'post':
                data.append('board', this.form.board);
                data.append('thread', this.form.thread);
                data.append('title', this.form.title);
                data.append('body', this.form.body);
                this.api.addPost(data).subscribe(function (res) {
                    alert('add success');
                    _this.init();
                });
                break;
            case 'changeBoard':
                data.append('from_board', this.form.fromBoard);
                data.append('to_board', this.form.toBoard);
                data.append('thread', this.form.thread);
                this.api.importThread(data).subscribe(function (res) {
                    alert('add success');
                    _this.init();
                });
                break;
        }
    };
    return AddComponent;
}());
AddComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_0" /* Component */])({
        selector: 'add',
        template: __webpack_require__("./src/components/add/add.component.html"),
        styles: [__webpack_require__("./src/components/add/add.component.css")]
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__providers__["ApiService"] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__providers__["ApiService"]) === "function" && _a || Object, typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_2__angular_router__["a" /* Router */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__angular_router__["a" /* Router */]) === "function" && _b || Object, typeof (_c = typeof __WEBPACK_IMPORTED_MODULE_2__angular_router__["b" /* ActivatedRoute */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__angular_router__["b" /* ActivatedRoute */]) === "function" && _c || Object])
], AddComponent);

var _a, _b, _c;
//# sourceMappingURL=add.component.js.map

/***/ }),

/***/ "./src/components/alert/alert.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("./node_modules/css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, "", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "./src/components/alert/alert.component.html":
/***/ (function(module, exports) {

module.exports = "<div class=\"modal-header\">\n  <h4 class=\"modal-title\">{{title}}</h4>\n  <button type=\"button\" class=\"close\" aria-label=\"Close\" (click)=\"activeModal.dismiss(false)\">\n        <span aria-hidden=\"true\">&times;</span>\n      </button>\n</div>\n<div class=\"modal-body\">\n  <p>{{body}}</p>\n</div>\n<div class=\"modal-footer\">\n  <button type=\"button\" class=\"btn btn-secondary\" (click)=\"activeModal.close(false)\">Cancel</button>\n  <button type=\"button\" class=\"btn btn-secondary\" (click)=\"activeModal.close(true)\">Submit</button>\n</div>\n"

/***/ }),

/***/ "./src/components/alert/alert.component.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("./node_modules/@angular/core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__ng_bootstrap_ng_bootstrap__ = __webpack_require__("./node_modules/@ng-bootstrap/ng-bootstrap/index.js");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return AlertComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};


var AlertComponent = (function () {
    function AlertComponent(activeModal) {
        this.activeModal = activeModal;
        this.title = '';
        this.body = '';
    }
    AlertComponent.prototype.ngOnInit = function () { };
    return AlertComponent;
}());
AlertComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_0" /* Component */])({
        selector: 'app-alert',
        template: __webpack_require__("./src/components/alert/alert.component.html"),
        styles: [__webpack_require__("./src/components/alert/alert.component.css")]
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__ng_bootstrap_ng_bootstrap__["b" /* NgbActiveModal */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__ng_bootstrap_ng_bootstrap__["b" /* NgbActiveModal */]) === "function" && _a || Object])
], AlertComponent);

var _a;
//# sourceMappingURL=alert.component.js.map

/***/ }),

/***/ "./src/components/boards/boards-list.component.html":
/***/ (function(module, exports) {

module.exports = "<div class='boards'>\n  <div class='container-fluid'>\n    <div class=\"page-header\">\n      <h1>All Boards</h1>\n    </div>\n    <p class=\"btn-group\"><button [class.disabled]=\"!isRoot\" type=\"button\" class=\"btn btn-outline-primary\" (click)=\"openAdd(add)\">New Board</button></p>\n    <table class=\"table table-hover table-bordered\">\n      <thead>\n        <tr>\n          <!--<th>Board</th>-->\n          <th>Name</th>\n          <th>Description</th>\n          <th>Created</th>\n          <th></th>\n          <th></th>\n        </tr>\n      </thead>\n      <tbody>\n        <tr *ngFor=\"let board of boards;let i = index;\" (click)=\"openThreads($event,board.public_key,board.url)\">\n          <!--<td class=\"url\"><a href=\"javascript:void(0);\" (click)=\"openThreads($event,board.public_key,board.url)\">{{board.url}}</a></td>-->\n          <td class=\"title\" title=\"{{board?.name}}\">{{board?.name}}</td>\n          <td class=\"description\" title=\"{{board?.description}}\">{{board?.description}}</td>\n          <td class=\"created\">{{board?.created / 1000000 | date: 'short'}}</td>\n          <td class=\"subscribe\" title=\"{{board.ui_options?.subscribe ? 'Subscribe':'Unsubscribe'}}\" (click)=\"subscribe($event,board.public_key,i)\"><a href=\"javascript:void(0);\"><i class=\"fa\" [class.fa-star-o]=\"!board.ui_options?.subscribe\" [class.fa-star]=\"board.ui_options?.subscribe\"></i></a></td>\n          <td class=\"subscribe\" title=\"Board Info\" (click)=\"openInfo($event,board,info)\"><a href=\"javascript:void(0);\"><i class=\"fa fa-info-circle\"></i></a></td>\n        </tr>\n      </tbody>\n    </table>\n    <h3 class=\"boardNot\" *ngIf=\"boards?.length == 0\">Not Found Boards</h3>\n\n  </div>\n</div>\n<ng-template #add let-c=\"close\" let-d=\"dismiss\">\n  <div class=\"modal-header\">\n    <h4 class=\"modal-title\">New Board</h4>\n    <button type=\"button\" class=\"close\" aria-label=\"Close\" (click)=\"d('Cross click')\">\n      <span aria-hidden=\"true\">&times;</span>\n    </button>\n  </div>\n  <div class=\"modal-body\">\n    <form [formGroup]=\"addForm\" novalidate>\n      <div class=\"form-group\">\n        <label for=\"name\">Board name</label>\n        <input type=\"text\" class=\"form-control\" placeholder=\"name\" id=\"name\" formControlName=\"name\">\n      </div>\n      <div class=\"form-group\">\n        <label for=\"description\">Board description</label>\n        <textarea class=\"form-control\" rows=\"3\" id=\"description\" formControlName=\"description\"></textarea>\n      </div>\n      <div class=\"form-group\">\n        <label for=\"seed\">Board seed</label>\n        <input type=\"text\" class=\"form-control\" placeholder=\"seed\" id=\"seed\" formControlName=\"seed\">\n      </div>\n    </form>\n  </div>\n  <div class=\"modal-footer\">\n    <button type=\"button\" class=\"btn btn-secondary\" (click)=\"c(false)\">cancel</button>\n    <button type=\"button\" class=\"btn btn-secondary\" (click)=\"c(true)\">submit</button>\n  </div>\n</ng-template>\n\n<ng-template #info let-c=\"close\" let-d=\"dismiss\">\n  <div class=\"modal-header\">\n    <h4 class=\"modal-title\">Board</h4>\n    <button type=\"button\" class=\"close\" aria-label=\"Close\" (click)=\"d('Cross click')\">\n      <span aria-hidden=\"true\">&times;</span>\n    </button>\n  </div>\n  <div class=\"modal-body\">\n    <div class=\"form-group\">\n      <label for=\"tmpName\">Name</label>\n      <input type=\"text\" class=\"form-control\" id=\"tmpName\" [(ngModel)]=\"tmpBoard.name\">\n    </div>\n    <div class=\"form-group\">\n      <label for=\"tmpDescription\">Description</label>\n      <input type=\"text\" class=\"form-control\" id=\"tmpDescription\" [(ngModel)]=\"tmpBoard.description\">\n    </div>\n    <div class=\"form-group\">\n      <label for=\"tmpPublicKey\">Public Key</label>\n      <input type=\"text\" class=\"form-control\" id=\"tmpPublicKey\" [(ngModel)]=\"tmpBoard.public_key\">\n    </div>\n    <div class=\"form-group\">\n      <label for=\"tmpUrl\">Url</label>\n      <input type=\"text\" class=\"form-control\" id=\"tmpUrl\" [(ngModel)]=\"tmpBoard.url\">\n    </div>\n    <div class=\"form-group\">\n      <label for=\"tmpCreated\">Created</label>\n      <input type=\"text\" class=\"form-control\" id=\"tmpCreated\" [(ngModel)]=\"tmpBoard.created\">\n    </div>\n  </div>\n  <div class=\"modal-footer\">\n    <button type=\"button\" class=\"btn btn-secondary\" (click)=\"c(false)\">close</button>\n  </div>\n</ng-template>\n"

/***/ }),

/***/ "./src/components/boards/boards-list.component.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("./node_modules/@angular/core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__providers__ = __webpack_require__("./src/providers/index.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_router__ = __webpack_require__("./node_modules/@angular/router/@angular/router.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__ng_bootstrap_ng_bootstrap__ = __webpack_require__("./node_modules/@ng-bootstrap/ng-bootstrap/index.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__angular_forms__ = __webpack_require__("./node_modules/@angular/forms/@angular/forms.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5__animations_router_animations__ = __webpack_require__("./src/animations/router.animations.ts");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return BoardsListComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};






var BoardsListComponent = (function () {
    function BoardsListComponent(api, user, router, modal, common) {
        this.api = api;
        this.user = user;
        this.router = router;
        this.modal = modal;
        this.common = common;
        this.routeAnimation = true;
        this.display = 'block';
        this.position = 'absolute';
        this.board = new __WEBPACK_IMPORTED_MODULE_0__angular_core__["F" /* EventEmitter */]();
        this.isRoot = false;
        this.boards = [];
        this.addForm = new __WEBPACK_IMPORTED_MODULE_4__angular_forms__["e" /* FormGroup */]({
            name: new __WEBPACK_IMPORTED_MODULE_4__angular_forms__["f" /* FormControl */](),
            description: new __WEBPACK_IMPORTED_MODULE_4__angular_forms__["f" /* FormControl */](),
            seed: new __WEBPACK_IMPORTED_MODULE_4__angular_forms__["f" /* FormControl */]()
        });
        this.tmpBoard = null;
    }
    BoardsListComponent.prototype.ngOnInit = function () {
        var _this = this;
        this.getBoards();
        this.api.getStats().subscribe(function (root) {
            _this.isRoot = root;
        });
    };
    BoardsListComponent.prototype.getBoards = function () {
        var _this = this;
        this.api.getBoards().subscribe(function (boards) {
            _this.boards = boards;
            _this.boards.forEach(function (el) {
                var data = new FormData();
                data.append('board', el.public_key);
                _this.api.getSubscription(data).subscribe(function (res) {
                    if (res.config && res.config.secret_key) {
                        el.ui_options = { subscribe: true };
                    }
                });
            });
        });
    };
    BoardsListComponent.prototype.openInfo = function (ev, board, content) {
        ev.stopImmediatePropagation();
        ev.stopPropagation();
        this.tmpBoard = board;
        this.modal.open(content, { size: 'lg' });
    };
    BoardsListComponent.prototype.openAdd = function (content) {
        var _this = this;
        this.modal.open(content).result.then(function (result) {
            if (result === true) {
                var data = new FormData();
                data.append('name', _this.addForm.get('name').value);
                data.append('description', _this.addForm.get('description').value);
                data.append('seed', _this.addForm.get('seed').value);
                _this.api.addBoard(data).subscribe(function (res) {
                    _this.api.getBoards().subscribe(function (boards) {
                        _this.boards = boards;
                        _this.common.showAlert('Added successfully', 'success', 3000);
                    });
                });
            }
        }, function (err) { });
    };
    BoardsListComponent.prototype.subscribe = function (ev, key, index) {
        var _this = this;
        ev.stopImmediatePropagation();
        ev.stopPropagation();
        var data = new FormData();
        data.append('board', key);
        if (!this.boards[index].ui_options.subscribe) {
            this.api.subscribe(data).subscribe(function (isOk) {
                var options = { subscribe: isOk };
                _this.boards[index].ui_options = options;
                _this.common.showAlert('Subscribe successfully', 'success', 3000);
            });
        }
        else {
            this.api.unSubscribe(data).subscribe(function (isOk) {
                if (isOk) {
                    _this.boards[index].ui_options.subscribe = false;
                    _this.common.showAlert('Unsubscribe successfully', 'success', 3000);
                    _this.getBoards();
                }
            });
        }
    };
    BoardsListComponent.prototype.openThreads = function (ev, key, url) {
        ev.stopImmediatePropagation();
        ev.stopPropagation();
        this.router.navigate(['/threads', { board: key }]);
        // this.board.emit(this.boards[0].public_key);
    };
    BoardsListComponent.prototype.getDismissReason = function (reason) {
        console.log('get dismiss reason:', reason);
    };
    return BoardsListComponent;
}());
__decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_11" /* HostBinding */])('@routeAnimation'),
    __metadata("design:type", Object)
], BoardsListComponent.prototype, "routeAnimation", void 0);
__decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_11" /* HostBinding */])('style.display'),
    __metadata("design:type", Object)
], BoardsListComponent.prototype, "display", void 0);
__decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_11" /* HostBinding */])('style.position'),
    __metadata("design:type", Object)
], BoardsListComponent.prototype, "position", void 0);
__decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_1" /* Output */])(),
    __metadata("design:type", typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_0__angular_core__["F" /* EventEmitter */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_0__angular_core__["F" /* EventEmitter */]) === "function" && _a || Object)
], BoardsListComponent.prototype, "board", void 0);
BoardsListComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_0" /* Component */])({
        selector: 'boards-list',
        template: __webpack_require__("./src/components/boards/boards-list.component.html"),
        styles: [__webpack_require__("./src/components/boards/boards.css")],
        encapsulation: __WEBPACK_IMPORTED_MODULE_0__angular_core__["q" /* ViewEncapsulation */].None,
        animations: [__WEBPACK_IMPORTED_MODULE_5__animations_router_animations__["a" /* slideInLeftAnimation */]],
    }),
    __metadata("design:paramtypes", [typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_1__providers__["ApiService"] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__providers__["ApiService"]) === "function" && _b || Object, typeof (_c = typeof __WEBPACK_IMPORTED_MODULE_1__providers__["UserService"] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__providers__["UserService"]) === "function" && _c || Object, typeof (_d = typeof __WEBPACK_IMPORTED_MODULE_2__angular_router__["a" /* Router */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__angular_router__["a" /* Router */]) === "function" && _d || Object, typeof (_e = typeof __WEBPACK_IMPORTED_MODULE_3__ng_bootstrap_ng_bootstrap__["c" /* NgbModal */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_3__ng_bootstrap_ng_bootstrap__["c" /* NgbModal */]) === "function" && _e || Object, typeof (_f = typeof __WEBPACK_IMPORTED_MODULE_1__providers__["CommonService"] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__providers__["CommonService"]) === "function" && _f || Object])
], BoardsListComponent);

var _a, _b, _c, _d, _e, _f;
//# sourceMappingURL=boards-list.component.js.map

/***/ }),

/***/ "./src/components/boards/boards.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("./node_modules/css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, "boards-list {\n  display: block;\n  width: 100%;\n}\n\n.boards .btn-group {\n  float: right;\n}\n\n.boards .page-header {\n  border-bottom: none;\n}\n.boards .page-header >h1 {\n  text-align: center;\n  color: red;\n}\n\n.boards .boardNot {\n  width: 100%;\n  text-align: center;\n}\n\ntable>thead>tr {\n  background-image: linear-gradient(to bottom, #d9edf7 0, #b9def0 100%);\n}\n\ntable>tbody>tr {\n  cursor: pointer;\n}\n\ntable>tbody>tr>td {\n  /*cursor: pointer;*/\n  vertical-align: middle !important;\n}\n\ntable .url {\n  max-width: 90px;\n}\ntable .url > a {\n  text-decoration: underline;\n}\ntable .url > a:hover {\n  color: red;\n}\n\ntable .description,\ntable .title {\n  max-width: 180px;\n  overflow: hidden;\n  text-overflow: ellipsis;\n  white-space: nowrap;\n}\ntable .subscribe {\n  width: 25px;\n  text-align: center;\n}\ntable .subscribe i {\n  margin-left: 0;\n  color: red;\n}\ntable .created {\n  width: 200px;\n  text-align: center;\n}\n", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "./src/components/connection/connection.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("./node_modules/css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, ".connection .btn-group {\n  float: right;\n}\n.connection .index {\n  width: 20px;\n}\n.connection .del {\n  width: 20px;\n  text-align: center;\n  cursor: pointer;\n}\ntable>thead>tr {\n  background-image: linear-gradient(to bottom, #d9edf7 0, #b9def0 100%);\n}\n.connection .del:hover {\n  color: red;\n}", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "./src/components/connection/connection.component.html":
/***/ (function(module, exports) {

module.exports = "<div class=\"connection\">\n  <div class=\"container-fluid\">\n    <div class=\"page-header\">\n      <h1>All Connection</h1>\n      <p class=\"btn-group\">\n        <button type=\"button\" class=\"btn btn-outline-primary\" (click)=\"openAdd(form)\">New Connection</button>\n      </p>\n    </div>\n    <table class=\"table table-bordered\">\n      <thead>\n        <tr>\n          <th></th>\n          <th>Url</th>\n          <th>Remove</th>\n        </tr>\n      </thead>\n      <tbody>\n        <tr *ngFor=\"let item of list;let i =index;\">\n          <td class=\"index\">{{i + 1}}</td>\n          <td class=\"url\"><a href=\"javascript:void(0);\" (click)=\"openThreads($event,board.public_key,board.url)\">{{item}}</a></td>\n          <td class=\"del\" (click)=\"remove(item)\"><i class=\"fa fa-trash\" aria-hidden=\"true\"></i></td>\n        </tr>\n      </tbody>\n    </table>\n    <h3 class=\"not-found\" *ngIf=\"list?.length == 0\">Not Found Connections</h3>\n  </div>\n</div>\n\n<ng-template #form let-c=\"close\" let-d=\"dismiss\">\n  <div class=\"modal-header\">\n    <h4 class=\"modal-title\">New Board</h4>\n    <button type=\"button\" class=\"close\" aria-label=\"Close\" (click)=\"d('Cross click')\">\n      <span aria-hidden=\"true\">&times;</span>\n    </button>\n  </div>\n  <div class=\"modal-body\">\n    <input [(ngModel)]=\"addUrl\" type=\"text\" class=\"form-control\" placeholder=\"E.g: [::1]:7452\" aria-describedby=\"basic-addon1\">\n  </div>\n  <div class=\"modal-footer\">\n    <button type=\"button\" class=\"btn btn-secondary\" (click)=\"c(false)\">cancel</button>\n    <button type=\"button\" class=\"btn btn-secondary\" (click)=\"c(true)\">submit</button>\n  </div>\n</ng-template>\n"

/***/ }),

/***/ "./src/components/connection/connection.component.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("./node_modules/@angular/core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__providers__ = __webpack_require__("./src/providers/index.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__ng_bootstrap_ng_bootstrap__ = __webpack_require__("./node_modules/@ng-bootstrap/ng-bootstrap/index.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__animations_router_animations__ = __webpack_require__("./src/animations/router.animations.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__components__ = __webpack_require__("./src/components/index.ts");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return ConnectionComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};





var ConnectionComponent = (function () {
    function ConnectionComponent(conn, modal, common) {
        this.conn = conn;
        this.modal = modal;
        this.common = common;
        this.routeAnimation = true;
        this.display = 'block';
        this.list = [];
        this.addUrl = '';
    }
    ConnectionComponent.prototype.ngOnInit = function () {
        this.getAllConnections();
    };
    ConnectionComponent.prototype.getAllConnections = function () {
        var _this = this;
        this.conn.getAllConnections().subscribe(function (list) {
            _this.list = list;
        });
    };
    ConnectionComponent.prototype.openAdd = function (content) {
        var _this = this;
        this.addUrl = '';
        this.modal.open(content).result.then(function (result) {
            if (result) {
                if (!_this.addUrl) {
                    _this.common.showAlert('The link can not be empty', 'danger', 3000);
                    return;
                }
                var data = new FormData();
                data.append('address', _this.addUrl);
                _this.conn.addConnection(data).subscribe(function (isOk) {
                    if (isOk) {
                        _this.getAllConnections();
                        _this.common.showAlert('The connection was added successfully', 'success', 3000);
                    }
                });
            }
        }, function (err) { });
    };
    ConnectionComponent.prototype.remove = function (address) {
        var _this = this;
        var modalRef = this.modal.open(__WEBPACK_IMPORTED_MODULE_4__components__["h" /* AlertComponent */]);
        modalRef.componentInstance.title = 'Delete Connection';
        modalRef.componentInstance.body = "Do you delete the connection?";
        modalRef.result.then(function (result) {
            if (result) {
                var data = new FormData();
                data.append('address', address);
                _this.conn.removeConnection(data).subscribe(function (isOk) {
                    if (isOk) {
                        _this.getAllConnections();
                        _this.common.showAlert('The connection has been deleted', 'success', 3000);
                    }
                });
            }
        });
    };
    return ConnectionComponent;
}());
__decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_11" /* HostBinding */])('@routeAnimation'),
    __metadata("design:type", Object)
], ConnectionComponent.prototype, "routeAnimation", void 0);
__decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_11" /* HostBinding */])('style.display'),
    __metadata("design:type", Object)
], ConnectionComponent.prototype, "display", void 0);
ConnectionComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_0" /* Component */])({
        selector: 'app-connection',
        template: __webpack_require__("./src/components/connection/connection.component.html"),
        styles: [__webpack_require__("./src/components/connection/connection.component.css")],
        encapsulation: __WEBPACK_IMPORTED_MODULE_0__angular_core__["q" /* ViewEncapsulation */].None,
        animations: [__WEBPACK_IMPORTED_MODULE_3__animations_router_animations__["a" /* slideInLeftAnimation */]]
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__providers__["ConnectionService"] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__providers__["ConnectionService"]) === "function" && _a || Object, typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_2__ng_bootstrap_ng_bootstrap__["c" /* NgbModal */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__ng_bootstrap_ng_bootstrap__["c" /* NgbModal */]) === "function" && _b || Object, typeof (_c = typeof __WEBPACK_IMPORTED_MODULE_1__providers__["CommonService"] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__providers__["CommonService"]) === "function" && _c || Object])
], ConnectionComponent);

var _a, _b, _c;
//# sourceMappingURL=connection.component.js.map

/***/ }),

/***/ "./src/components/index.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__boards_boards_list_component__ = __webpack_require__("./src/components/boards/boards-list.component.ts");
/* harmony namespace reexport (by used) */ __webpack_require__.d(__webpack_exports__, "a", function() { return __WEBPACK_IMPORTED_MODULE_0__boards_boards_list_component__["a"]; });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__threads_threads__ = __webpack_require__("./src/components/threads/threads.ts");
/* harmony namespace reexport (by used) */ __webpack_require__.d(__webpack_exports__, "b", function() { return __WEBPACK_IMPORTED_MODULE_1__threads_threads__["a"]; });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__threadPage_threadPage__ = __webpack_require__("./src/components/threadPage/threadPage.ts");
/* harmony namespace reexport (by used) */ __webpack_require__.d(__webpack_exports__, "c", function() { return __WEBPACK_IMPORTED_MODULE_2__threadPage_threadPage__["a"]; });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__add_add_component__ = __webpack_require__("./src/components/add/add.component.ts");
/* harmony namespace reexport (by used) */ __webpack_require__.d(__webpack_exports__, "d", function() { return __WEBPACK_IMPORTED_MODULE_3__add_add_component__["a"]; });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__userlist_userlist_component__ = __webpack_require__("./src/components/userlist/userlist.component.ts");
/* harmony namespace reexport (by used) */ __webpack_require__.d(__webpack_exports__, "e", function() { return __WEBPACK_IMPORTED_MODULE_4__userlist_userlist_component__["a"]; });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5__user_user_component__ = __webpack_require__("./src/components/user/user.component.ts");
/* harmony namespace reexport (by used) */ __webpack_require__.d(__webpack_exports__, "f", function() { return __WEBPACK_IMPORTED_MODULE_5__user_user_component__["a"]; });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_6__connection_connection_component__ = __webpack_require__("./src/components/connection/connection.component.ts");
/* harmony namespace reexport (by used) */ __webpack_require__.d(__webpack_exports__, "g", function() { return __WEBPACK_IMPORTED_MODULE_6__connection_connection_component__["a"]; });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_7__alert_alert_component__ = __webpack_require__("./src/components/alert/alert.component.ts");
/* harmony namespace reexport (by used) */ __webpack_require__.d(__webpack_exports__, "h", function() { return __WEBPACK_IMPORTED_MODULE_7__alert_alert_component__["a"]; });







// Modal

//# sourceMappingURL=index.js.map

/***/ }),

/***/ "./src/components/threadPage/threadPage.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("./node_modules/css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, "threadPage {\n  width: 100%;\n}\n\n.thread .thread-description {\n  min-height: 200px;\n}\n.thread .btn-box {\n  margin: 1rem 0;\n  text-align: right;\n}\n.post {\n  /*max-width: 50%;*/\n  margin: 10px 0;\n  background-color: #e5ecfd;\n}\n", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "./src/components/threadPage/threadPage.html":
/***/ (function(module, exports) {

module.exports = "<div class=\"container-fluid\">\n  <div class=\"card thread\" *ngIf=\"data\">\n    <div class=\"card-block\">\n      <h3 class=\"card-title\">{{data.thread.name}}</h3>\n      <p class=\"thread-description\">{{data.thread.description}}</p>\n      <p class=\"btn-box\"><button class=\"btn btn-info\" (click)=\"openReply(addPost)\">reply</button></p>\n    </div>\n  </div>\n  <div class=\"card post\" *ngFor=\"let item of data.posts\">\n    <div class=\"card-block\">\n      <h5 class=\"card-title\">{{item.title}}</h5>\n      <a class=\"card-text\">{{item.author}}</a>\n      <p class=\"card-text\" [innerHTML]=\"item.body | safeHtml\"></p>\n      <p class=\"card-text\">{{item.created / 1000000 | date:'yMMMdjms'}}</p>\n    </div>\n  </div>\n</div>\n\n\n<ng-template #addPost let-c=\"close\" let-d=\"dismiss\">\n  <div class=\"modal-header\">\n    <h4 class=\"modal-title\">New Post</h4>\n    <button type=\"button\" class=\"close\" aria-label=\"Close\" (click)=\"d('Cross click')\">\n      <span aria-hidden=\"true\">&times;</span>\n    </button>\n  </div>\n  <div class=\"modal-body\">\n\n    <form [formGroup]=\"postForm\" novalidate>\n      <div class=\"form-group\">\n        <label for=\"title\">Title</label>\n        <input type=\"text\" class=\"form-control \" placeholder=\"title\" id=\"title\" formControlName=\"title\">\n      </div>\n      <!--<div [froalaEditor]=\"editorOptions\"></div>-->\n      <div class=\"form-group\">\n        <label for=\"body\">Content</label>\n        <textarea [froalaEditor]=\"editorOptions\" formControlName=\"body\"></textarea>\n        <!--<textarea class=\"form-control\" rows=\"3\" id=\"body\" formControlName=\"body\"></textarea>-->\n      </div>\n    </form>\n  </div>\n  <div class=\"modal-footer\">\n    <button type=\"button\" class=\"btn btn-secondary\" (click)=\"c(false)\">cancel</button>\n    <button type=\"button\" class=\"btn btn-primary\" (click)=\"c(true)\">submit</button>\n  </div>\n</ng-template>\n"

/***/ }),

/***/ "./src/components/threadPage/threadPage.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("./node_modules/@angular/core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__providers__ = __webpack_require__("./src/providers/index.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_router__ = __webpack_require__("./node_modules/@angular/router/@angular/router.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__angular_forms__ = __webpack_require__("./node_modules/@angular/forms/@angular/forms.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__ng_bootstrap_ng_bootstrap__ = __webpack_require__("./node_modules/@ng-bootstrap/ng-bootstrap/index.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5__animations_router_animations__ = __webpack_require__("./src/animations/router.animations.ts");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return ThreadPageComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};






var ThreadPageComponent = (function () {
    function ThreadPageComponent(api, router, route, modal, common) {
        this.api = api;
        this.router = router;
        this.route = route;
        this.modal = modal;
        this.common = common;
        this.routeAnimation = true;
        this.display = 'block';
        this.boardKey = '';
        this.threadKey = '';
        this.data = { posts: [], thread: { name: '', description: '' } };
        this.postForm = new __WEBPACK_IMPORTED_MODULE_3__angular_forms__["e" /* FormGroup */]({
            title: new __WEBPACK_IMPORTED_MODULE_3__angular_forms__["f" /* FormControl */](),
            body: new __WEBPACK_IMPORTED_MODULE_3__angular_forms__["f" /* FormControl */]()
        });
        this.editorOptions = {
            placeholderText: 'Edit Your Content Here!',
            // toolbarButtons: ['bold', 'italic', 'underline', 'strikeThrough', 'subscript', 'superscript', '|', 'fontFamily', 'fontSize', 'color', 'inlineStyle', 'paragraphStyle', '|', 'paragraphFormat', 'align', 'formatOL', 'formatUL', 'outdent', 'indent', 'quote', '-', 'insertLink', 'insertImage', 'insertVideo', 'insertFile', 'insertTable', '|', 'emoticons', 'specialCharacters', 'insertHR', 'selectAll', 'clearFormatting', '|', 'print', 'spellChecker', 'help', 'html', '|', 'undo', 'redo'],
            toolbarButtons: ['bold', 'italic', 'underline', 'strikeThrough', 'subscript', 'superscript', '|', 'fontFamily', 'fontSize', 'color', 'inlineStyle', 'paragraphStyle', '|', 'paragraphFormat', 'align', 'formatOL', 'formatUL', 'outdent', 'indent', 'quote', '-', 'insertLink', '|', 'emoticons', 'specialCharacters', 'insertHR', 'selectAll', 'clearFormatting', '|', 'print', 'spellChecker', 'help', 'html', '|', 'undo', 'redo'],
            heightMin: 200,
            events: {},
        };
    }
    ThreadPageComponent.prototype.ngOnInit = function () {
        var _this = this;
        this.route.params.subscribe(function (res) {
            _this.boardKey = res['board'];
            _this.threadKey = res['thread'];
            _this.open(_this.boardKey, _this.threadKey);
        });
    };
    ThreadPageComponent.prototype.openReply = function (content) {
        var _this = this;
        this.postForm.reset();
        this.modal.open(content, { backdrop: 'static', size: 'lg' }).result.then(function (result) {
            var data = new FormData();
            if (result) {
                data.append('board', _this.boardKey);
                data.append('thread', _this.threadKey);
                data.append('title', _this.postForm.get('title').value);
                data.append('body', _this.postForm.get('body').value);
                _this.api.addPost(data).subscribe(function (post) {
                    if (post) {
                        console.log('add post successfully:', post);
                        _this.data.posts.unshift(post);
                        _this.common.showAlert('Added successfully', 'success', 3000);
                    }
                });
            }
        }, function (err) { });
    };
    ThreadPageComponent.prototype.reply = function () {
        if (!this.boardKey || !this.threadKey) {
            alert('Will not be able to post');
            return;
        }
        this.router.navigate(['/add', { exec: 'post', board: this.boardKey, thread: this.threadKey }]);
    };
    ThreadPageComponent.prototype.open = function (master, ref) {
        var _this = this;
        var data = new FormData();
        data.append('board', master);
        data.append('thread', ref);
        this.api.getThreadpage(data).subscribe(function (data) {
            _this.data = data;
        });
    };
    ThreadPageComponent.prototype.windowScroll = function (event) {
        this.common.showOrHideToTopBtn();
    };
    return ThreadPageComponent;
}());
__decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_11" /* HostBinding */])('@routeAnimation'),
    __metadata("design:type", Object)
], ThreadPageComponent.prototype, "routeAnimation", void 0);
__decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_11" /* HostBinding */])('style.display'),
    __metadata("design:type", Object)
], ThreadPageComponent.prototype, "display", void 0);
__decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_17" /* HostListener */])('window:scroll', ['$event']),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Object]),
    __metadata("design:returntype", void 0)
], ThreadPageComponent.prototype, "windowScroll", null);
ThreadPageComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_0" /* Component */])({
        selector: 'threadPage',
        template: __webpack_require__("./src/components/threadPage/threadPage.html"),
        styles: [__webpack_require__("./src/components/threadPage/threadPage.css")],
        encapsulation: __WEBPACK_IMPORTED_MODULE_0__angular_core__["q" /* ViewEncapsulation */].None,
        animations: [__WEBPACK_IMPORTED_MODULE_5__animations_router_animations__["a" /* slideInLeftAnimation */]]
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__providers__["ApiService"] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__providers__["ApiService"]) === "function" && _a || Object, typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_2__angular_router__["a" /* Router */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__angular_router__["a" /* Router */]) === "function" && _b || Object, typeof (_c = typeof __WEBPACK_IMPORTED_MODULE_2__angular_router__["b" /* ActivatedRoute */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__angular_router__["b" /* ActivatedRoute */]) === "function" && _c || Object, typeof (_d = typeof __WEBPACK_IMPORTED_MODULE_4__ng_bootstrap_ng_bootstrap__["c" /* NgbModal */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_4__ng_bootstrap_ng_bootstrap__["c" /* NgbModal */]) === "function" && _d || Object, typeof (_e = typeof __WEBPACK_IMPORTED_MODULE_1__providers__["CommonService"] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__providers__["CommonService"]) === "function" && _e || Object])
], ThreadPageComponent);

var _a, _b, _c, _d, _e;
//# sourceMappingURL=threadPage.js.map

/***/ }),

/***/ "./src/components/threads/threads.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("./node_modules/css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, "threads {\n  width: 100%;\n}\n.threads .btn-group {\n  float: right;\n}\n.threads table>thead>tr {\n  background-image: linear-gradient(to bottom, #d9edf7 0, #b9def0 100%);\n}\n.threads .page-header {\n  border-bottom: none;\n}\n.threads .page-header>h1 {\n  color: red;\n  text-align: center;\n}\n\n.threads table .name {\n  min-width: 10px;\n}\n.threads table .transfer {\n  width: 25px;\n  cursor: pointer;\n}\n.threads table .transfer i {\n  margin-left: 0;\n}\n.threads table .transfer i:hover {\n  color: red;\n}\n.threads table .board {\n  min-width: 100px;\n}\n.threads table .board > a {\n  text-decoration: underline;\n}\n.threads table .board >a:hover {\n  color: red;\n}\n", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "./src/components/threads/threads.html":
/***/ (function(module, exports) {

module.exports = "<div class=\"threads\">\n  <div class=\"container-fluid\">\n    <!--<ol class=\"breadcrumb\">\n      <li><a outerLink=\"/\">All Board</a></li>\n      <li><a outerLink=\"/\">Home</a></li>\n    </ol>-->\n    <div class=\"page-header\">\n      <h1>All Threads</h1>\n    </div>\n    <p class=\"btn-group\"><button type=\"button\" class=\"btn btn-outline-primary\" (click)=\"openAdd(content)\">New Thread</button></p>\n\n    <table class=\"table table-hover table-bordered\">\n      <thead>\n        <tr>\n          <th>Thread</th>\n          <th>Board</th>\n          <th>Name</th>\n          <th></th>\n        </tr>\n      </thead>\n      <tbody>\n        <tr *ngFor=\"let thread of threads\">\n          <td class=\"name\" title=\"{{thread.name}}\">{{thread.name}}</td>\n          <td class=\"board\" title=\"{{url}}\"><a href=\"javascript:void(0);\" (click)=\"open(thread?.master_board,thread?.ref)\">{{board.url}}</a></td>\n          <td class=\"description\" title=\"{{thread.description}}\">{{thread.description}}</td>\n          <td class=\"transfer\" (click)=\"openImport(importBox,thread.ref)\"><i title=\"Import Thread\" class=\"fa fa-exchange\"></i></td>\n        </tr>\n      </tbody>\n    </table>\n    <h3 class=\"not-found\" *ngIf=\"threads?.length == 0\">Not Found Threads</h3>\n  </div>\n</div>\n\n<ng-template #content let-c=\"close\" let-d=\"dismiss\">\n  <div class=\"modal-header\">\n    <h4 class=\"modal-title\">New Thread</h4>\n    <button type=\"button\" class=\"close\" aria-label=\"Close\" (click)=\"d('Cross click')\">\n      <span aria-hidden=\"true\">&times;</span>\n    </button>\n  </div>\n  <div class=\"modal-body\">\n    <form [formGroup]=\"addForm\" novalidate>\n      <div class=\"form-group\">\n        <label for=\"name\">Thread name</label>\n        <input type=\"text\" class=\"form-control\" placeholder=\"name\" id=\"name\" formControlName=\"name\">\n      </div>\n      <div class=\"form-group\">\n        <label for=\"description\">Thread description</label>\n        <input type=\"text\" class=\"form-control\" placeholder=\"description\" id=\"description\" formControlName=\"description\">\n      </div>\n    </form>\n  </div>\n  <div class=\"modal-footer\">\n    <button type=\"button\" class=\"btn btn-secondary\" (click)=\"c(false)\">cancel</button>\n    <button type=\"button\" class=\"btn btn-secondary\" (click)=\"c(true)\">submit</button>\n  </div>\n</ng-template>\n\n<ng-template #importBox let-c=\"close\" let-d=\"dismiss\">\n  <div class=\"modal-header\">\n    <h4 class=\"modal-title\">Import Thread</h4>\n    <button type=\"button\" class=\"close\" aria-label=\"Close\" (click)=\"d('Cross click')\">\n      <span aria-hidden=\"true\">&times;</span>\n    </button>\n  </div>\n  <div class=\"modal-body\">\n    <div class=\"form-group\">\n      <label for=\"toBoard\">Choose to import</label>\n      <select class=\"form-control\" id=\"toBoard\" [(ngModel)]=\"importBoardKey\" placeholder=\"Choose to move\">\n        <option *ngFor=\"let board of importBoards\" value=\"{{board.public_key}}\">{{board.name}}</option>\n    </select>\n      <!--<input type=\"text\" class=\"form-control\" placeholder=\"toBoard\" id=\"toBoard\">-->\n    </div>\n  </div>\n  <div class=\"modal-footer\">\n    <button type=\"button\" class=\"btn btn-secondary\" (click)=\"c(false)\">cancel</button>\n    <button type=\"button\" class=\"btn btn-secondary\" (click)=\"c(true)\">submit</button>\n  </div>\n</ng-template>\n"

/***/ }),

/***/ "./src/components/threads/threads.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("./node_modules/@angular/core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__providers__ = __webpack_require__("./src/providers/index.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_router__ = __webpack_require__("./node_modules/@angular/router/@angular/router.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__angular_forms__ = __webpack_require__("./node_modules/@angular/forms/@angular/forms.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__ng_bootstrap_ng_bootstrap__ = __webpack_require__("./node_modules/@ng-bootstrap/ng-bootstrap/index.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5__animations_router_animations__ = __webpack_require__("./src/animations/router.animations.ts");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return ThreadsComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};






var ThreadsComponent = (function () {
    function ThreadsComponent(api, router, route, modal, common) {
        this.api = api;
        this.router = router;
        this.route = route;
        this.modal = modal;
        this.common = common;
        this.routeAnimation = true;
        this.display = 'block';
        this.thread = new __WEBPACK_IMPORTED_MODULE_0__angular_core__["F" /* EventEmitter */]();
        this.threads = [];
        this.importBoards = [];
        this.importBoardKey = '';
        this.boardKey = '';
        this.board = null;
        this.addForm = new __WEBPACK_IMPORTED_MODULE_3__angular_forms__["e" /* FormGroup */]({
            description: new __WEBPACK_IMPORTED_MODULE_3__angular_forms__["f" /* FormControl */](),
            name: new __WEBPACK_IMPORTED_MODULE_3__angular_forms__["f" /* FormControl */]()
        });
    }
    ThreadsComponent.prototype.ngOnInit = function () {
        var _this = this;
        this.route.params.subscribe(function (res) {
            // this.url = res['url'];
            _this.boardKey = res['board'];
            _this.init();
        });
    };
    ThreadsComponent.prototype.initThreads = function (key) {
        var _this = this;
        var data = new FormData();
        data.append('board', key);
        this.api.getThreads(data).subscribe(function (threads) {
            _this.threads = threads;
        });
    };
    ThreadsComponent.prototype.init = function () {
        var _this = this;
        var data = new FormData();
        data.append('board', this.boardKey);
        this.api.getBoardPage(data).subscribe(function (data) {
            _this.board = data.board;
            _this.threads = data.threads;
        });
    };
    ThreadsComponent.prototype.openAdd = function (content) {
        var _this = this;
        this.modal.open(content).result.then(function (result) {
            if (result) {
                var data = new FormData();
                data.append('board', _this.boardKey);
                data.append('description', _this.addForm.get('description').value);
                data.append('name', _this.addForm.get('name').value);
                _this.api.addThread(data).subscribe(function (thread) {
                    _this.threads.unshift(thread);
                    _this.common.showAlert('Added successfully', 'success', 3000);
                });
            }
        }, function (err) { });
    };
    ThreadsComponent.prototype.open = function (master, ref) {
        this.router.navigate(['p', { board: master, thread: ref }], { relativeTo: this.route });
    };
    ThreadsComponent.prototype.openImport = function (content, threadKey) {
        var _this = this;
        if (this.importBoards.length <= 0) {
            this.api.getBoards().subscribe(function (boards) {
                _this.importBoards = boards;
                _this.importBoardKey = boards[0].public_key;
            });
        }
        this.modal.open(content, { size: 'lg' }).result.then(function (result) {
            if (result) {
                if (_this.importBoardKey) {
                    var data = new FormData();
                    data.append('from_board', _this.boardKey);
                    data.append('thread', threadKey);
                    data.append('to_board', _this.importBoardKey);
                    _this.api.importThread(data).subscribe(function (res) {
                        console.log('transfer thread:', res);
                        _this.common.showAlert('successfully', 'success', 3000);
                        _this.initThreads(_this.boardKey);
                    });
                }
            }
        }, function (err) { });
    };
    return ThreadsComponent;
}());
__decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_11" /* HostBinding */])('@routeAnimation'),
    __metadata("design:type", Object)
], ThreadsComponent.prototype, "routeAnimation", void 0);
__decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_11" /* HostBinding */])('style.display'),
    __metadata("design:type", Object)
], ThreadsComponent.prototype, "display", void 0);
__decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_1" /* Output */])(),
    __metadata("design:type", typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_0__angular_core__["F" /* EventEmitter */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_0__angular_core__["F" /* EventEmitter */]) === "function" && _a || Object)
], ThreadsComponent.prototype, "thread", void 0);
ThreadsComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_0" /* Component */])({
        selector: 'threads',
        template: __webpack_require__("./src/components/threads/threads.html"),
        styles: [__webpack_require__("./src/components/threads/threads.css")],
        encapsulation: __WEBPACK_IMPORTED_MODULE_0__angular_core__["q" /* ViewEncapsulation */].None,
        animations: [__WEBPACK_IMPORTED_MODULE_5__animations_router_animations__["a" /* slideInLeftAnimation */]]
    }),
    __metadata("design:paramtypes", [typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_1__providers__["ApiService"] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__providers__["ApiService"]) === "function" && _b || Object, typeof (_c = typeof __WEBPACK_IMPORTED_MODULE_2__angular_router__["a" /* Router */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__angular_router__["a" /* Router */]) === "function" && _c || Object, typeof (_d = typeof __WEBPACK_IMPORTED_MODULE_2__angular_router__["b" /* ActivatedRoute */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__angular_router__["b" /* ActivatedRoute */]) === "function" && _d || Object, typeof (_e = typeof __WEBPACK_IMPORTED_MODULE_4__ng_bootstrap_ng_bootstrap__["c" /* NgbModal */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_4__ng_bootstrap_ng_bootstrap__["c" /* NgbModal */]) === "function" && _e || Object, typeof (_f = typeof __WEBPACK_IMPORTED_MODULE_1__providers__["CommonService"] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__providers__["CommonService"]) === "function" && _f || Object])
], ThreadsComponent);

var _a, _b, _c, _d, _e, _f;
//# sourceMappingURL=threads.js.map

/***/ }),

/***/ "./src/components/user/user.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("./node_modules/css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, ".user-box {\n  display: block;\n  width: 100%;\n}\n\n.user-box .item span {\n  margin: 0 10px;\n}\n.user-box .item .exec {\n  font-size: 12px;\n  text-decoration: underline;\n}\n.user-box .item i {\n  margin-right: 5px;\n}", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "./src/components/user/user.component.html":
/***/ (function(module, exports) {

module.exports = "<ng-template #content let-c=\"close\" let-d=\"dismiss\">\n  <div class=\"modal-header\">\n    <h4 class=\"modal-title\">Edit</h4>\n    <button type=\"button\" class=\"close\" aria-label=\"Close\" (click)=\"d('Cross click')\">\n      <span aria-hidden=\"true\">&times;</span>\n    </button>\n  </div>\n  <div class=\"modal-body\">\n    <input type=\"text\" class=\"form-control\" placeholder=\"Name\" aria-describedby=\"basic-addon1\">\n  </div>\n  <div class=\"modal-footer\">\n    <button type=\"button\" class=\"btn btn-secondary\" (click)=\"edit()\">Close</button>\n    <button type=\"button\" class=\"btn btn-secondary\" (click)=\"edit()\">Edit</button>\n  </div>\n</ng-template>\n<div class=\"user-box\">\n  <div class=\"container\">\n    <div class=\"card\">\n      <div class=\"card-header\">\n        User Info\n      </div>\n      <div class=\"card-block\">\n        <p class=\"item\"><i class=\"fa fa-user\" aria-hidden=\"true\"></i>Name:<span class=\"name\" >{{user?.alias}}</span><a href=\"javascript:void(0);\" (click)=\"switchUser(switch)\">switch</a></p>\n        <p class=\"item\"><i class=\"fa fa-suitcase\" aria-hidden=\"true\"></i>Master:<span class=\"master\">{{user?.master}}</span></p>\n        <p class=\"item\"><i class=\"fa fa-key\" aria-hidden=\"true\"></i> Key:<span class=\"key\">{{user?.public_key}}</span></p>\n      </div>\n    </div>\n  </div>\n</div>\n<ng-template #switch let-c=\"close\" let-d=\"dismiss\">\n  <div class=\"modal-header\">\n    <h4 class=\"modal-title\">Switch Users</h4>\n    <button type=\"button\" class=\"close\" aria-label=\"Close\" (click)=\"d('Cross click')\">\n      <span aria-hidden=\"true\">&times;</span>\n    </button>\n  </div>\n  <div class=\"modal-body\">\n    <div class=\"form-group\">\n      <label for=\"user\">Choose to move</label>\n      <select class=\"form-control\" id=\"user\" [(ngModel)]=\"switchUserKey\" placeholder=\"Choose to move\">\n        <option *ngFor=\"let user of switchUserList\" value=\"{{user.public_key}}\">{{user.alias}}</option>\n    </select>\n    </div>\n  </div>\n  <div class=\"modal-footer\">\n    <button type=\"button\" class=\"btn btn-secondary\" (click)=\"c(false)\">cancel</button>\n    <button type=\"button\" class=\"btn btn-secondary\" (click)=\"c(true)\">switch</button>\n  </div>\n</ng-template>\n"

/***/ }),

/***/ "./src/components/user/user.component.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("./node_modules/@angular/core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__providers__ = __webpack_require__("./src/providers/index.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__animations_router_animations__ = __webpack_require__("./src/animations/router.animations.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__ng_bootstrap_ng_bootstrap__ = __webpack_require__("./node_modules/@ng-bootstrap/ng-bootstrap/index.js");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return UserComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};




var UserComponent = (function () {
    function UserComponent(userService, modal) {
        this.userService = userService;
        this.modal = modal;
        this.routeAnimation = true;
        this.display = 'block';
        this.user = null;
        this.switchUserKey = '';
        this.switchUserList = [];
    }
    UserComponent.prototype.ngOnInit = function () {
        var _this = this;
        this.userService.getCurrent().subscribe(function (user) {
            _this.user = user;
        });
    };
    UserComponent.prototype.switchUser = function (content) {
        var _this = this;
        if (this.switchUserList.length <= 0) {
            this.userService.getAll().subscribe(function (users) {
                _this.switchUserList = users;
                _this.switchUserKey = users[0].public_key;
            });
        }
        this.modal.open(content).result.then(function (result) {
            if (result) {
                console.log('switchKey:', _this.switchUserKey);
                var data = new FormData();
                data.append('user', _this.switchUserKey);
                _this.userService.setCurrent(data).subscribe(function (res) {
                    location.reload();
                });
            }
        }, function (err) { });
    };
    return UserComponent;
}());
__decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_11" /* HostBinding */])('@routeAnimation'),
    __metadata("design:type", Object)
], UserComponent.prototype, "routeAnimation", void 0);
__decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_11" /* HostBinding */])('style.display'),
    __metadata("design:type", Object)
], UserComponent.prototype, "display", void 0);
UserComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_0" /* Component */])({
        selector: 'app-user',
        template: __webpack_require__("./src/components/user/user.component.html"),
        styles: [__webpack_require__("./src/components/user/user.component.css")],
        animations: [__WEBPACK_IMPORTED_MODULE_2__animations_router_animations__["a" /* slideInLeftAnimation */]]
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__providers__["UserService"] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__providers__["UserService"]) === "function" && _a || Object, typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_3__ng_bootstrap_ng_bootstrap__["c" /* NgbModal */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_3__ng_bootstrap_ng_bootstrap__["c" /* NgbModal */]) === "function" && _b || Object])
], UserComponent);

var _a, _b;
//# sourceMappingURL=user.component.js.map

/***/ }),

/***/ "./src/components/userlist/userlist.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("./node_modules/css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, ".user-list {}\n\n.user-list table .alias {\n  min-width: 10px;\n}\n\n.user-list table .master {\n  min-width: 10px;\n}\n\n.user-list table .key {\n  min-width: 10px;\n}\n.user-list table .edit,\n.user-list table .del {\n  width: 20px;\n  cursor: pointer;\n  text-align: center;\n}\n.user-list table .edit:hover,\n.user-list table .del:hover {\n  color: red;\n}\n", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "./src/components/userlist/userlist.component.html":
/***/ (function(module, exports) {

module.exports = "<div class=\"user-list\">\n  <table class=\"table table-hover table-bordered\">\n    <thead>\n      <tr>\n        <th>Name</th>\n        <th>Master</th>\n        <th>Public Key</th>\n        <th>Edit</th>\n        <th>Remove</th>\n      </tr>\n    </thead>\n    <tbody>\n      <tr *ngFor=\"let user of userlist\">\n        <td class=\"alias\">{{user.alias}}</td>\n        <td class=\"master\">{{user.master}}</td>\n        <td class=\"key\">{{user.public_key}}</td>\n        <td class=\"edit\" (click)=\"openEdit(editBox,user.public_key)\"><i class=\"fa fa-pencil\" aria-hidden=\"true\"></i></td>\n        <td class=\"del\" (click)=\"remove($event,user.public_key)\"><i class=\"fa fa-trash\" aria-hidden=\"true\"></i></td>\n      </tr>\n    </tbody>\n  </table>\n</div>\n<ng-template #editBox let-c=\"close\" let-d=\"dismiss\">\n  <div class=\"modal-header\">\n    <h4 class=\"modal-title\">Edit</h4>\n    <button type=\"button\" class=\"close\" aria-label=\"Close\" (click)=\"c(false)\">\n        <span aria-hidden=\"true\">&times;</span>\n      </button>\n  </div>\n  <div class=\"modal-body\">\n    <input [(ngModel)]=\"editName\" type=\"text\" class=\"form-control\" placeholder=\"Username\" aria-describedby=\"basic-addon1\">\n  </div>\n  <div class=\"modal-footer\">\n    <button type=\"button\" class=\"btn btn-secondary\" (click)=\"c(false)\">cancel</button>\n    <button type=\"button\" class=\"btn btn-secondary\" (click)=\"c(true)\">submit</button>\n  </div>\n</ng-template>\n"

/***/ }),

/***/ "./src/components/userlist/userlist.component.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("./node_modules/@angular/core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__providers__ = __webpack_require__("./src/providers/index.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__ng_bootstrap_ng_bootstrap__ = __webpack_require__("./node_modules/@ng-bootstrap/ng-bootstrap/index.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__animations_router_animations__ = __webpack_require__("./src/animations/router.animations.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__alert_alert_component__ = __webpack_require__("./src/components/alert/alert.component.ts");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return UserlistComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};





var UserlistComponent = (function () {
    function UserlistComponent(user, modal, common) {
        this.user = user;
        this.modal = modal;
        this.common = common;
        this.routeAnimation = true;
        this.display = 'block';
        // @HostBinding('style.position') position = 'absolute';
        this.userlist = [];
        this.editName = '';
    }
    UserlistComponent.prototype.ngOnInit = function () {
        var _this = this;
        this.user.getAll().subscribe(function (userlist) {
            _this.userlist = userlist;
        });
    };
    UserlistComponent.prototype.openEdit = function (content, key) {
        var _this = this;
        var modalRef = this.modal.open(content).result.then(function (result) {
            if (result) {
                _this.edit(_this.editName, key);
            }
        });
    };
    UserlistComponent.prototype.edit = function (name, key) {
        var _this = this;
        var data = new FormData();
        data.append('alias', name);
        data.append('user', key);
        this.user.newOrModifyUser(data).subscribe(function (res) {
            _this.userlist = [];
            _this.user.getAll().subscribe(function (userlist) {
                _this.userlist = userlist;
                _this.common.showAlert('successfully modified', 'success', 3000);
            });
        });
    };
    UserlistComponent.prototype.remove = function (ev, key) {
        var _this = this;
        ev.stopImmediatePropagation();
        ev.stopPropagation();
        var data = new FormData();
        data.append('user', key);
        var modalRef = this.modal.open(__WEBPACK_IMPORTED_MODULE_4__alert_alert_component__["a" /* AlertComponent */]);
        modalRef.componentInstance.title = 'Delete User';
        modalRef.componentInstance.body = 'Do you delete the user?';
        modalRef.result.then(function (result) {
            if (result) {
                _this.user.remove(data).subscribe(function (isOk) {
                    if (isOk) {
                        _this.userlist = [];
                        _this.user.getAll().subscribe(function (userlist) {
                            _this.userlist = userlist;
                            _this.common.showAlert('successfully deleted', 'success', 1000);
                        });
                    }
                    else {
                        _this.common.showAlert('failed to delete', 'success', 1000);
                    }
                });
            }
        });
    };
    UserlistComponent.prototype.windowScroll = function (event) {
        this.common.showOrHideToTopBtn();
    };
    return UserlistComponent;
}());
__decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_11" /* HostBinding */])('@routeAnimation'),
    __metadata("design:type", Object)
], UserlistComponent.prototype, "routeAnimation", void 0);
__decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_11" /* HostBinding */])('style.display'),
    __metadata("design:type", Object)
], UserlistComponent.prototype, "display", void 0);
__decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_17" /* HostListener */])('window:scroll', ['$event']),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Object]),
    __metadata("design:returntype", void 0)
], UserlistComponent.prototype, "windowScroll", null);
UserlistComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_0" /* Component */])({
        selector: 'app-userlist',
        template: __webpack_require__("./src/components/userlist/userlist.component.html"),
        styles: [__webpack_require__("./src/components/userlist/userlist.component.css")],
        animations: [__WEBPACK_IMPORTED_MODULE_3__animations_router_animations__["a" /* slideInLeftAnimation */]],
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__providers__["UserService"] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__providers__["UserService"]) === "function" && _a || Object, typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_2__ng_bootstrap_ng_bootstrap__["c" /* NgbModal */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__ng_bootstrap_ng_bootstrap__["c" /* NgbModal */]) === "function" && _b || Object, typeof (_c = typeof __WEBPACK_IMPORTED_MODULE_1__providers__["CommonService"] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__providers__["CommonService"]) === "function" && _c || Object])
], UserlistComponent);

var _a, _b, _c;
//# sourceMappingURL=userlist.component.js.map

/***/ }),

/***/ "./src/environments/environment.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return environment; });
// The file contents for the current environment will overwrite these during build.
// The build system defaults to the dev environment which uses `environment.ts`, but if you do
// `ng build --env=prod` then `environment.prod.ts` will be used instead.
// The list of which env maps to which file can be found in `.angular-cli.json`.
// The file contents for the current environment will overwrite these during build.
var environment = {
    production: false
};
//# sourceMappingURL=environment.js.map

/***/ }),

/***/ "./src/main.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
Object.defineProperty(__webpack_exports__, "__esModule", { value: true });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("./node_modules/@angular/core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_platform_browser_dynamic__ = __webpack_require__("./node_modules/@angular/platform-browser-dynamic/@angular/platform-browser-dynamic.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__app_app_module__ = __webpack_require__("./src/app/app.module.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__environments_environment__ = __webpack_require__("./src/environments/environment.ts");




if (__WEBPACK_IMPORTED_MODULE_3__environments_environment__["a" /* environment */].production) {
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["a" /* enableProdMode */])();
}
__webpack_require__.i(__WEBPACK_IMPORTED_MODULE_1__angular_platform_browser_dynamic__["a" /* platformBrowserDynamic */])().bootstrapModule(__WEBPACK_IMPORTED_MODULE_2__app_app_module__["a" /* AppModule */]);
//# sourceMappingURL=main.js.map

/***/ }),

/***/ "./src/pipes/index.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__safe_safe_html_pipe__ = __webpack_require__("./src/pipes/safe/safe-html.pipe.ts");
/* harmony namespace reexport (by used) */ __webpack_require__.d(__webpack_exports__, "a", function() { return __WEBPACK_IMPORTED_MODULE_0__safe_safe_html_pipe__["a"]; });

//# sourceMappingURL=index.js.map

/***/ }),

/***/ "./src/pipes/safe/safe-html.pipe.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("./node_modules/@angular/core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_platform_browser__ = __webpack_require__("./node_modules/@angular/platform-browser/@angular/platform-browser.es5.js");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return SafeHTMLPipe; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};


var SafeHTMLPipe = (function () {
    function SafeHTMLPipe(sanitizer) {
        this.sanitizer = sanitizer;
    }
    SafeHTMLPipe.prototype.transform = function (html) {
        return this.sanitizer.bypassSecurityTrustHtml(html);
    };
    return SafeHTMLPipe;
}());
SafeHTMLPipe = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Y" /* Pipe */])({
        name: 'safeHtml'
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__angular_platform_browser__["c" /* DomSanitizer */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__angular_platform_browser__["c" /* DomSanitizer */]) === "function" && _a || Object])
], SafeHTMLPipe);

var _a;
//# sourceMappingURL=safe-html.pipe.js.map

/***/ }),

/***/ "./src/providers/api/api.service.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("./node_modules/@angular/core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_http__ = __webpack_require__("./node_modules/@angular/http/@angular/http.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__common_common_service__ = __webpack_require__("./src/providers/common/common.service.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_rxjs_add_operator_map__ = __webpack_require__("./node_modules/rxjs/add/operator/map.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_rxjs_add_operator_map___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_3_rxjs_add_operator_map__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4_rxjs_add_operator_catch__ = __webpack_require__("./node_modules/rxjs/add/operator/catch.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4_rxjs_add_operator_catch___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_4_rxjs_add_operator_catch__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return ApiService; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};





var ApiService = (function () {
    function ApiService(http, common) {
        this.http = http;
        this.common = common;
        this.base_url = 'http://127.0.0.1:7410/api/';
    }
    ApiService.prototype.getSubscriptions = function () {
        return this.common.handleGet(this.base_url + "get_subscriptions");
    };
    ApiService.prototype.getSubscription = function (data) {
        return this.common.handlePost(this.base_url + 'get_subscription', data);
    };
    ApiService.prototype.subscribe = function (data) {
        return this.common.handlePost(this.base_url + 'subscribe', data);
    };
    ApiService.prototype.unSubscribe = function (data) {
        return this.common.handlePost(this.base_url + 'unsubscribe', data);
    };
    ApiService.prototype.getStats = function () {
        return this.common.handleGet(this.base_url + 'get_stats');
    };
    ApiService.prototype.getThreads = function (data) {
        return this.common.handlePost(this.base_url + 'get_threads', data);
    };
    ApiService.prototype.getBoards = function () {
        return this.common.handleGet(this.base_url + 'get_boards');
    };
    ApiService.prototype.getPosts = function (data) {
        return this.common.handlePost(this.base_url + 'get_posts', data);
    };
    ApiService.prototype.getBoardPage = function (data) {
        return this.common.handlePost(this.base_url + 'get_boardpage', data);
    };
    ApiService.prototype.getThreadpage = function (data) {
        return this.common.handlePost(this.base_url + 'get_threadpage', data);
    };
    ApiService.prototype.addBoard = function (data) {
        return this.common.handlePost(this.base_url + 'new_board', data);
    };
    ApiService.prototype.addThread = function (data) {
        return this.common.handlePost(this.base_url + 'new_thread', data);
    };
    ApiService.prototype.addPost = function (data) {
        return this.common.handlePost(this.base_url + 'new_post', data);
    };
    ApiService.prototype.importThread = function (data) {
        return this.common.handlePost(this.base_url + 'import_thread', data);
    };
    return ApiService;
}());
ApiService = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["c" /* Injectable */])(),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__angular_http__["b" /* Http */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__angular_http__["b" /* Http */]) === "function" && _a || Object, typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_2__common_common_service__["a" /* CommonService */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__common_common_service__["a" /* CommonService */]) === "function" && _b || Object])
], ApiService);

var _a, _b;
//# sourceMappingURL=api.service.js.map

/***/ }),

/***/ "./src/providers/api/msg.ts":
/***/ (function(module, exports) {

//# sourceMappingURL=msg.js.map

/***/ }),

/***/ "./src/providers/common/common.service.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_http__ = __webpack_require__("./node_modules/@angular/http/@angular/http.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_core__ = __webpack_require__("./node_modules/@angular/core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs_add_observable_throw__ = __webpack_require__("./node_modules/rxjs/add/observable/throw.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs_add_observable_throw___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_rxjs_add_observable_throw__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_rxjs_add_operator_filter__ = __webpack_require__("./node_modules/rxjs/add/operator/filter.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_rxjs_add_operator_filter___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_3_rxjs_add_operator_filter__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4_rxjs_Observable__ = __webpack_require__("./node_modules/rxjs/Observable.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4_rxjs_Observable___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_4_rxjs_Observable__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return CommonService; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};





var CommonService = (function () {
    function CommonService(http) {
        this.http = http;
        this.alertType = 'info';
        this.alertMessage = '';
        this.alert = false;
        this.topBtn = false;
    }
    CommonService.prototype.handleError = function (error) {
        console.error('Error:', error.json() || 'Server error', 'danger');
        this.showAlert(error.json() || 'Server error', 'danger', 3000);
        return __WEBPACK_IMPORTED_MODULE_4_rxjs_Observable__["Observable"].throw(error.json() || 'Server error');
    };
    CommonService.prototype.handleGet = function (url) {
        var _this = this;
        if (!url) {
            return __WEBPACK_IMPORTED_MODULE_4_rxjs_Observable__["Observable"].throw('The connection is empty');
        }
        return this.http.get(url).filter(function (res) { return res.status === 200; }).map(function (res) { return res.json(); }).catch(function (err) { return _this.handleError(err); });
    };
    CommonService.prototype.handlePost = function (url, data) {
        var _this = this;
        if (!url || !data) {
            return __WEBPACK_IMPORTED_MODULE_4_rxjs_Observable__["Observable"].throw('Parameters and connections can not be empty');
        }
        return this.http.post(url, data).filter(function (res) { return res.status === 200; }).map(function (res) { return res.json(); }).catch(function (err) { return _this.handleError(err); });
    };
    CommonService.prototype.showAlert = function (message, type, timeout) {
        var _this = this;
        this.alert = false;
        this.alertMessage = message;
        if (type) {
            this.alertType = type;
        }
        if (timeout > 0) {
            setTimeout(function () {
                _this.alert = false;
            }, timeout);
        }
        this.alert = true;
    };
    CommonService.prototype.showOrHideToTopBtn = function () {
        var pos = (document.documentElement.scrollTop || document.body.scrollTop) + document.documentElement.offsetHeight;
        var max = document.documentElement.scrollHeight;
        if (pos > (max / 3)) {
            this.topBtn = true;
        }
        else {
            this.topBtn = false;
        }
    };
    CommonService.prototype.scrollToTop = function () {
        window.scrollTo(0, 0);
    };
    return CommonService;
}());
CommonService = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_1__angular_core__["c" /* Injectable */])(),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_0__angular_http__["b" /* Http */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_0__angular_http__["b" /* Http */]) === "function" && _a || Object])
], CommonService);

var _a;
//# sourceMappingURL=common.service.js.map

/***/ }),

/***/ "./src/providers/connection/connection.service.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("./node_modules/@angular/core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_http__ = __webpack_require__("./node_modules/@angular/http/@angular/http.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs_add_operator_map__ = __webpack_require__("./node_modules/rxjs/add/operator/map.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs_add_operator_map___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_rxjs_add_operator_map__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_rxjs_add_operator_catch__ = __webpack_require__("./node_modules/rxjs/add/operator/catch.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_rxjs_add_operator_catch___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_3_rxjs_add_operator_catch__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__common_common_service__ = __webpack_require__("./src/providers/common/common.service.ts");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return ConnectionService; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};





var ConnectionService = (function () {
    function ConnectionService(http, common) {
        this.http = http;
        this.common = common;
        this.base_url = 'http://127.0.0.1:7410/api/connections/';
    }
    ConnectionService.prototype.getAllConnections = function () {
        var _this = this;
        return this.http.get(this.base_url + 'get_all').map(function (res) { return res.json(); }).catch(function (err) { return _this.common.handleError(err); });
    };
    ConnectionService.prototype.addConnection = function (data) {
        var _this = this;
        return this.http.post(this.base_url + 'new', data).map(function (res) { return res.json(); }).catch(function (err) { return _this.common.handleError(err); });
    };
    ConnectionService.prototype.removeConnection = function (data) {
        var _this = this;
        return this.http.post(this.base_url + 'remove', data).map(function (res) { return res.json(); }).catch(function (err) { return _this.common.handleError(err); });
    };
    return ConnectionService;
}());
ConnectionService = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["c" /* Injectable */])(),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__angular_http__["b" /* Http */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__angular_http__["b" /* Http */]) === "function" && _a || Object, typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_4__common_common_service__["a" /* CommonService */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_4__common_common_service__["a" /* CommonService */]) === "function" && _b || Object])
], ConnectionService);

var _a, _b;
//# sourceMappingURL=connection.service.js.map

/***/ }),

/***/ "./src/providers/index.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__common_common_service__ = __webpack_require__("./src/providers/common/common.service.ts");
/* harmony namespace reexport (by used) */ __webpack_require__.d(__webpack_exports__, "CommonService", function() { return __WEBPACK_IMPORTED_MODULE_0__common_common_service__["a"]; });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__api_api_service__ = __webpack_require__("./src/providers/api/api.service.ts");
/* harmony namespace reexport (by used) */ __webpack_require__.d(__webpack_exports__, "ApiService", function() { return __WEBPACK_IMPORTED_MODULE_1__api_api_service__["a"]; });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__api_msg__ = __webpack_require__("./src/providers/api/msg.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__api_msg___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2__api_msg__);
/* harmony namespace reexport (by used) */ if(__webpack_require__.o(__WEBPACK_IMPORTED_MODULE_2__api_msg__, "UserService")) __webpack_require__.d(__webpack_exports__, "UserService", function() { return __WEBPACK_IMPORTED_MODULE_2__api_msg__["UserService"]; });
/* harmony namespace reexport (by used) */ if(__webpack_require__.o(__WEBPACK_IMPORTED_MODULE_2__api_msg__, "ConnectionService")) __webpack_require__.d(__webpack_exports__, "ConnectionService", function() { return __WEBPACK_IMPORTED_MODULE_2__api_msg__["ConnectionService"]; });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__user_user_service__ = __webpack_require__("./src/providers/user/user.service.ts");
/* harmony namespace reexport (by used) */ __webpack_require__.d(__webpack_exports__, "UserService", function() { return __WEBPACK_IMPORTED_MODULE_3__user_user_service__["a"]; });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__user_user_msg__ = __webpack_require__("./src/providers/user/user.msg.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__user_user_msg___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_4__user_user_msg__);
/* harmony namespace reexport (by used) */ if(__webpack_require__.o(__WEBPACK_IMPORTED_MODULE_4__user_user_msg__, "ConnectionService")) __webpack_require__.d(__webpack_exports__, "ConnectionService", function() { return __WEBPACK_IMPORTED_MODULE_4__user_user_msg__["ConnectionService"]; });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5__connection_connection_service__ = __webpack_require__("./src/providers/connection/connection.service.ts");
/* harmony namespace reexport (by used) */ __webpack_require__.d(__webpack_exports__, "ConnectionService", function() { return __WEBPACK_IMPORTED_MODULE_5__connection_connection_service__["a"]; });






//# sourceMappingURL=index.js.map

/***/ }),

/***/ "./src/providers/user/user.msg.ts":
/***/ (function(module, exports) {

//# sourceMappingURL=user.msg.js.map

/***/ }),

/***/ "./src/providers/user/user.service.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("./node_modules/@angular/core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_http__ = __webpack_require__("./node_modules/@angular/http/@angular/http.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs_add_operator_map__ = __webpack_require__("./node_modules/rxjs/add/operator/map.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs_add_operator_map___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_rxjs_add_operator_map__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_rxjs_add_operator_catch__ = __webpack_require__("./node_modules/rxjs/add/operator/catch.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_rxjs_add_operator_catch___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_3_rxjs_add_operator_catch__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__common_common_service__ = __webpack_require__("./src/providers/common/common.service.ts");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return UserService; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};





var UserService = (function () {
    function UserService(http, common) {
        this.http = http;
        this.common = common;
        this.base_url = 'http://127.0.0.1:7410/api/users/';
    }
    UserService.prototype.getCurrent = function () {
        var _this = this;
        return this.http.get(this.base_url + 'get_current')
            .map(function (response) { return response.json(); })
            .catch(function (err) { return _this.common.handleError(err); });
    };
    UserService.prototype.getAllMasters = function () {
        var _this = this;
        return this.http.get(this.base_url + 'get_masters').map(function (res) { return res.json(); }).catch(function (err) { return _this.common.handleError(err); });
    };
    UserService.prototype.getAll = function () {
        var _this = this;
        return this.http.get(this.base_url + 'get_all').map(function (res) { return res.json(); }).catch(function (err) { return _this.common.handleError(err); });
    };
    UserService.prototype.setCurrent = function (data) {
        var _this = this;
        return this.http.post(this.base_url + 'set_current', data).map(function (res) { return res.json(); }).catch(function (err) { return _this.common.handleError(err); });
    };
    UserService.prototype.newMaster = function (data) {
        var _this = this;
        return this.http.post(this.base_url + 'new_master', data).map(function (res) { return res.json(); }).catch(function (err) { return _this.common.handleError(err); });
    };
    UserService.prototype.newOrModifyUser = function (data) {
        var _this = this;
        return this.http.post(this.base_url + 'new', data).map(function (res) { return res.json(); }).catch(function (err) { return _this.common.handleError(err); });
    };
    UserService.prototype.remove = function (data) {
        var _this = this;
        return this.http.post(this.base_url + 'remove', data).map(function (res) { return res.json(); }).catch(function (err) { return _this.common.handleError(err); });
    };
    return UserService;
}());
UserService = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["c" /* Injectable */])(),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__angular_http__["b" /* Http */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__angular_http__["b" /* Http */]) === "function" && _a || Object, typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_4__common_common_service__["a" /* CommonService */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_4__common_common_service__["a" /* CommonService */]) === "function" && _b || Object])
], UserService);

var _a, _b;
//# sourceMappingURL=user.service.js.map

/***/ }),

/***/ "./src/router/app-router-routing.module.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("./node_modules/@angular/core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_router__ = __webpack_require__("./node_modules/@angular/router/@angular/router.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__components__ = __webpack_require__("./src/components/index.ts");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return AppRouterRoutingModule; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};



var routes = [
    { path: '', component: __WEBPACK_IMPORTED_MODULE_2__components__["a" /* BoardsListComponent */], pathMatch: 'full' },
    {
        path: 'threads', children: [
            { path: '', component: __WEBPACK_IMPORTED_MODULE_2__components__["b" /* ThreadsComponent */] },
            { path: 'p', component: __WEBPACK_IMPORTED_MODULE_2__components__["c" /* ThreadPageComponent */] },
        ]
    },
    // { path: 'threads', component: ThreadsComponent },
    { path: 'add', component: __WEBPACK_IMPORTED_MODULE_2__components__["d" /* AddComponent */] },
    { path: 'userlist', component: __WEBPACK_IMPORTED_MODULE_2__components__["e" /* UserlistComponent */] },
    { path: 'user', component: __WEBPACK_IMPORTED_MODULE_2__components__["f" /* UserComponent */] },
    { path: 'conn', component: __WEBPACK_IMPORTED_MODULE_2__components__["g" /* ConnectionComponent */] },
    { path: '**', redirectTo: '' }
];
var AppRouterRoutingModule = (function () {
    function AppRouterRoutingModule() {
    }
    return AppRouterRoutingModule;
}());
AppRouterRoutingModule = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["b" /* NgModule */])({
        imports: [__WEBPACK_IMPORTED_MODULE_1__angular_router__["d" /* RouterModule */].forRoot(routes)],
        exports: [__WEBPACK_IMPORTED_MODULE_1__angular_router__["d" /* RouterModule */]],
    })
], AppRouterRoutingModule);

//# sourceMappingURL=app-router-routing.module.js.map

/***/ }),

/***/ 1:
/***/ (function(module, exports, __webpack_require__) {

module.exports = __webpack_require__("./src/main.ts");


/***/ })

},[1]);
//# sourceMappingURL=main.bundle.js.map