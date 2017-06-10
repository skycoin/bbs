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

/***/ "./src/app/app.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("./node_modules/css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, "", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "./src/app/app.component.html":
/***/ (function(module, exports) {

module.exports = "<!--The whole content below can be removed with the new code.-->\n<!--<boards-list (board)=\"openThreads($event)\"></boards-list>\n<threads (thread)=\"openThreadpage($event)\"></threads>\n<threadPage></threadPage>\n<add></add>-->\n<nav class=\"navbar navbar-default navbar-fixed-top\">\n  <div class=\"container-fluid\">\n    <a class=\"navbar-brand\" routerLink=\"/\" >BBS</a>\n    <div class=\" collapse navbar-collapse \" id=\"bs-example-navbar-collapse-1 \">\n      <ul class=\"nav navbar-nav \">\n        <li><a routerLink=\"/\">Board</a></li>\n        <li><a routerLink=\"/add\">Add/Change</a></li>\n      </ul>\n    </div>\n  </div>\n\n</nav>\n<router-outlet></router-outlet>\n"

/***/ }),

/***/ "./src/app/app.component.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("./node_modules/@angular/core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__providers__ = __webpack_require__("./src/providers/index.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__components__ = __webpack_require__("./src/components/index.ts");
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
    function AppComponent(api) {
        this.api = api;
        this.title = 'app';
    }
    AppComponent.prototype.ngOnInit = function () {
        this.api.getBoards();
    };
    AppComponent.prototype.openThreads = function (key) {
        this.threads.start(key);
    };
    AppComponent.prototype.openThreadpage = function (data) {
        this.threadPage.open(data.master, data.ref);
    };
    AppComponent.prototype.test = function () {
        console.log('test');
    };
    return AppComponent;
}());
__decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_15" /* ViewChild */])(__WEBPACK_IMPORTED_MODULE_2__components__["a" /* BoardsListComponent */]),
    __metadata("design:type", typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_2__components__["a" /* BoardsListComponent */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__components__["a" /* BoardsListComponent */]) === "function" && _a || Object)
], AppComponent.prototype, "boards", void 0);
__decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_15" /* ViewChild */])(__WEBPACK_IMPORTED_MODULE_2__components__["b" /* ThreadsComponent */]),
    __metadata("design:type", typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_2__components__["b" /* ThreadsComponent */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__components__["b" /* ThreadsComponent */]) === "function" && _b || Object)
], AppComponent.prototype, "threads", void 0);
__decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_15" /* ViewChild */])(__WEBPACK_IMPORTED_MODULE_2__components__["c" /* ThreadPageComponent */]),
    __metadata("design:type", typeof (_c = typeof __WEBPACK_IMPORTED_MODULE_2__components__["c" /* ThreadPageComponent */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__components__["c" /* ThreadPageComponent */]) === "function" && _c || Object)
], AppComponent.prototype, "threadPage", void 0);
AppComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_14" /* Component */])({
        selector: 'app-root',
        template: __webpack_require__("./src/app/app.component.html"),
        styles: [__webpack_require__("./src/app/app.component.css")]
    }),
    __metadata("design:paramtypes", [typeof (_d = typeof __WEBPACK_IMPORTED_MODULE_1__providers__["a" /* ApiService */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__providers__["a" /* ApiService */]) === "function" && _d || Object])
], AppComponent);

var _a, _b, _c, _d;
//# sourceMappingURL=app.component.js.map

/***/ }),

/***/ "./src/app/app.module.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_platform_browser__ = __webpack_require__("./node_modules/@angular/platform-browser/@angular/platform-browser.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_core__ = __webpack_require__("./node_modules/@angular/core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_http__ = __webpack_require__("./node_modules/@angular/http/@angular/http.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__angular_forms__ = __webpack_require__("./node_modules/@angular/forms/@angular/forms.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__app_component__ = __webpack_require__("./src/app/app.component.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5__providers__ = __webpack_require__("./src/providers/index.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_6__router_app_router_routing_module__ = __webpack_require__("./src/router/app-router-routing.module.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_7__components__ = __webpack_require__("./src/components/index.ts");
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
        declarations: [
            __WEBPACK_IMPORTED_MODULE_4__app_component__["a" /* AppComponent */], __WEBPACK_IMPORTED_MODULE_7__components__["a" /* BoardsListComponent */], __WEBPACK_IMPORTED_MODULE_7__components__["b" /* ThreadsComponent */], __WEBPACK_IMPORTED_MODULE_7__components__["c" /* ThreadPageComponent */], __WEBPACK_IMPORTED_MODULE_7__components__["d" /* AddComponent */]
        ],
        imports: [
            __WEBPACK_IMPORTED_MODULE_0__angular_platform_browser__["a" /* BrowserModule */], __WEBPACK_IMPORTED_MODULE_2__angular_http__["a" /* HttpModule */], __WEBPACK_IMPORTED_MODULE_3__angular_forms__["a" /* FormsModule */], __WEBPACK_IMPORTED_MODULE_6__router_app_router_routing_module__["a" /* AppRouterRoutingModule */]
        ],
        providers: [__WEBPACK_IMPORTED_MODULE_5__providers__["a" /* ApiService */]],
        bootstrap: [__WEBPACK_IMPORTED_MODULE_4__app_component__["a" /* AppComponent */]]
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

module.exports = "<div class=\"container\">\n  <h1>Add</h1>\n  <div class=\"form-group\">\n    <label for=\"type\">Select add</label>\n    <select class=\"form-control\" id=\"type\" [ngModel]=\"select\" (ngModelChange)=\"clear($event)\">\n    <option value=\"board\">Add Board</option>\n      <option value=\"thread\">Add Thread</option>\n      <option value=\"post\">Add Post</option>\n      <option value=\"changeBoard\">ChangeBoard</option>\n  </select>\n  </div>\n  <!--<form>-->\n  <div class=\"form-group\" [hidden]=\"select != 'board' && select != 'thread'\">\n    <label for=\"name\">Input board name</label>\n    <input type=\"text\" class=\"form-control\" placeholder=\"name\" id=\"name\" [(ngModel)]=\"form.name\">\n  </div>\n  <div class=\"form-group\" [hidden]=\"select != 'thread' && select != 'post'\">\n    <label for=\"board\">Input board key</label>\n    <input type=\"text\" class=\"form-control\" placeholder=\"board\" id=\"board\" [(ngModel)]=\"form.board\">\n  </div>\n  <div class=\"form-group\" [hidden]=\"select != 'post' && select != 'changeBoard'\">\n    <label for=\"thread\">Input thread key</label>\n    <input type=\"text\" class=\"form-control\" placeholder=\"thread\" id=\"thread\" [(ngModel)]=\"form.thread\">\n  </div>\n  <div class=\"form-group\" [hidden]=\"select != 'post'\">\n    <label for=\"title\">Input title</label>\n    <input type=\"text\" class=\"form-control\" placeholder=\"title\" id=\"title\" [(ngModel)]=\"form.title\">\n  </div>\n  <div class=\"form-group\" [hidden]=\"select != 'post'\">\n    <label for=\"body\">Input Post Body</label>\n    <textarea class=\"form-control\" rows=\"3\" id=\"body\" [(ngModel)]=\"form.body\"></textarea>\n  </div>\n  <div class=\"form-group\" [hidden]=\"select != 'board' && select != 'thread'\">\n    <label for=\"description\">Input description</label>\n    <input type=\"text\" class=\"form-control\" placeholder=\"description\" id=\"description\" [(ngModel)]=\"form.description\">\n  </div>\n  <div class=\"form-group\" [hidden]=\"select != 'board'\">\n    <label for=\"seed\">Input seed</label>\n    <input type=\"text\" class=\"form-control\" placeholder=\"seed\" id=\"seed\" [(ngModel)]=\"form.seed\">\n  </div>\n  <div class=\"form-group\" [hidden]=\"select != 'changeBoard'\">\n    <label for=\"fromBoard\">From Board</label>\n    <input type=\"text\" class=\"form-control\" placeholder=\"fromBoard\" id=\"fromBoard\" [(ngModel)]=\"form.fromBoard\">\n  </div>\n  <div class=\"form-group\" [hidden]=\"select != 'changeBoard'\">\n    <label for=\"toBoard\">To Board</label>\n    <input type=\"text\" class=\"form-control\" placeholder=\"toBoard\" id=\"toBoard\" [(ngModel)]=\"form.toBoard\">\n  </div>\n  <button class=\"btn btn-info\" (click)=\"add($event)\">Submit</button>\n  <!--</form>-->\n\n</div>\n"

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
                this.api.addBoard(data).then(function (res) {
                    alert('add success');
                    _this.init();
                }, function (err) {
                    console.error('err:', err);
                });
                break;
            case 'thread':
                data.append('board', this.form.board);
                data.append('description', this.form.description);
                data.append('name', this.form.name);
                this.api.addThread(data).then(function (res) {
                    alert('add success');
                    _this.init();
                    console.log('add success:', res);
                }, function (err) {
                    console.error('err:', err);
                });
                break;
            case 'post':
                data.append('board', this.form.board);
                data.append('thread', this.form.thread);
                data.append('title', this.form.title);
                data.append('body', this.form.body);
                this.api.addPost(data).then(function (res) {
                    console.log('add success:', res);
                    _this.init();
                    alert('add success');
                }, function (err) {
                    console.error('err:', err);
                });
                break;
            case 'changeBoard':
                data.append('from_board', this.form.fromBoard);
                data.append('to_board', this.form.toBoard);
                data.append('thread', this.form.thread);
                this.api.changeThread(data).then(function (res) {
                    console.log('add success:', res);
                    _this.init();
                    alert('add success');
                }, function (err) {
                    console.error('err:', err);
                });
                break;
        }
    };
    return AddComponent;
}());
AddComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_14" /* Component */])({
        selector: 'add',
        template: __webpack_require__("./src/components/add/add.component.html"),
        styles: [__webpack_require__("./src/components/add/add.component.css")]
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__providers__["a" /* ApiService */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__providers__["a" /* ApiService */]) === "function" && _a || Object, typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_2__angular_router__["b" /* Router */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__angular_router__["b" /* Router */]) === "function" && _b || Object, typeof (_c = typeof __WEBPACK_IMPORTED_MODULE_2__angular_router__["c" /* ActivatedRoute */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__angular_router__["c" /* ActivatedRoute */]) === "function" && _c || Object])
], AddComponent);

var _a, _b, _c;
//# sourceMappingURL=add.component.js.map

/***/ }),

/***/ "./src/components/boards/boards-list.component.html":
/***/ (function(module, exports) {

module.exports = "<div class='boards'>\n  <div class='container-fluid'>\n    <h1 class=\"title\">All Board</h1>\n    <table class=\"table table-hover\">\n      <thead>\n        <tr>\n          <th>Url</th>\n          <th>Title</th>\n          <th>Description</th>\n          <th>Created</th>\n        </tr>\n      </thead>\n      <tbody>\n        <tr *ngFor=\"let board of boards\" (click)=\"openThreads($event,board.public_key)\">\n          <td class=\"url\">{{board.url}}</td>\n          <td class=\"title\" title=\"{{board.name}}\">{{board.name}}</td>\n          <td class=\"description\" title=\"{{board.description}}\">{{board.description}}</td>\n          <td class=\"created\">{{board.created / 1000000 | date: 'short'}}</td>\n        </tr>\n\n      </tbody>\n\n    </table>\n    <h3 class=\"boardNot\" *ngIf=\"boards?.length == 0\">Not Found Board</h3>\n\n    <!--<div *ngFor=\"let board of boards\" class=\"item\" (click)=\"openThreads($event,board.public_key)\">\n      <div class=\"content\">\n        <h3 class=\"name\">{{board.name}}</h3>\n        <div class=\"description\">{{board.description}}</div>\n        <div class=\"created\">{{board.created / 1000000 | date: 'medium'}}</div>\n      </div>\n    </div>-->\n  </div>\n</div>\n"

/***/ }),

/***/ "./src/components/boards/boards-list.component.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("./node_modules/@angular/core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__providers__ = __webpack_require__("./src/providers/index.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_router__ = __webpack_require__("./node_modules/@angular/router/@angular/router.es5.js");
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
    function BoardsListComponent(api, router) {
        this.api = api;
        this.router = router;
        this.board = new __WEBPACK_IMPORTED_MODULE_0__angular_core__["F" /* EventEmitter */]();
        this.boardsTitle = 'Boards List';
        this.boards = [];
    }
    BoardsListComponent.prototype.ngOnInit = function () {
        var _this = this;
        this.api.getBoards().then(function (data) {
            _this.boards = data;
        });
    };
    BoardsListComponent.prototype.openThreads = function (ev, key) {
        ev.stopImmediatePropagation();
        ev.stopPropagation();
        this.router.navigate(['/threads', { board: key }]);
        // this.board.emit(this.boards[0].public_key);
    };
    return BoardsListComponent;
}());
__decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_5" /* Output */])(),
    __metadata("design:type", typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_0__angular_core__["F" /* EventEmitter */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_0__angular_core__["F" /* EventEmitter */]) === "function" && _a || Object)
], BoardsListComponent.prototype, "board", void 0);
BoardsListComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_14" /* Component */])({
        selector: 'boards-list',
        template: __webpack_require__("./src/components/boards/boards-list.component.html"),
        styles: [__webpack_require__("./src/components/boards/boards.css")],
        encapsulation: __WEBPACK_IMPORTED_MODULE_0__angular_core__["q" /* ViewEncapsulation */].None,
    }),
    __metadata("design:paramtypes", [typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_1__providers__["a" /* ApiService */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__providers__["a" /* ApiService */]) === "function" && _b || Object, typeof (_c = typeof __WEBPACK_IMPORTED_MODULE_2__angular_router__["b" /* Router */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__angular_router__["b" /* Router */]) === "function" && _c || Object])
], BoardsListComponent);

var _a, _b, _c;
//# sourceMappingURL=boards-list.component.js.map

/***/ }),

/***/ "./src/components/boards/boards.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("./node_modules/css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, "boards-list {\n  display: block;\n  width: 100%;\n}\n\n\n/*.container-fluid {\n  padding: 0;\n}*/\n\n.boards h1 {\n  text-align: center;\n  color: red;\n}\n\n.boards .boardNot {\n  width: 100%;\n  text-align: center;\n}\n\ntable>thead>tr {\n  background-image: linear-gradient(to bottom, #d9edf7 0, #b9def0 100%);\n}\n\ntable>tbody>tr {\n  border-bottom: 1px solid #ccc;\n}\n\ntable>tbody>tr>td {\n  cursor: pointer;\n  vertical-align: middle !important;\n}\n\ntable .url {\n  max-width: 90px;\n}\n\ntable .description,\ntable .title {\n  max-width: 180px;\n  overflow: hidden;\n  text-overflow: ellipsis;\n  white-space: nowrap;\n}\n\ntable .created {\n  width: 200px;\n}\n", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

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




//# sourceMappingURL=index.js.map

/***/ }),

/***/ "./src/components/threadPage/threadPage.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("./node_modules/css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, "threadPage {\n  width: 100%;\n}\n\n.thread {\n  border-bottom: .1rem solid #000;\n}\n.thread .thread-description {\n  min-height: 200px;\n}\n.thread .btn-box {\n  margin: 1rem 0;\n  text-align: right;\n}\n.posts {\n  /*border: .1rem solid #ccc;*/\n  margin: 1rem 0;\n  padding: 1rem;\n  background-color: hsla(0, 0%, 80%, 0.3);\n}\n.posts .post-author {\n  color: red;\n  text-decoration: underline;\n  cursor: pointer;\n}\n.posts .post-body {\n  margin: 1rem 0;\n  font-weight: 500;\n}", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "./src/components/threadPage/threadPage.html":
/***/ (function(module, exports) {

module.exports = "<div class=\"container-fluid\">\n  <div class=\"thread\">\n    <h3>{{data.thread.name}}</h3>\n    <p class=\"thread-description\">{{data.thread.description}}</p>\n    <div class=\"btn-box\"><button class=\"btn btn-info\" (click)=\"reply()\">reply</button></div>\n  </div>\n  <div class=\"posts\" *ngFor=\"let item of data.posts\">\n    <h5 class=\"post-title\">{{item.title}}</h5>\n    <a class=\"post-author\">{{item.author}}</a>\n    <p class=\"post-body\">{{item.body}}</p>\n    <p class=\"post-created\">{{item.created / 1000000 | date:'yMMMdjms'}}</p>\n  </div>\n</div>\n"

/***/ }),

/***/ "./src/components/threadPage/threadPage.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("./node_modules/@angular/core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__providers__ = __webpack_require__("./src/providers/index.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_router__ = __webpack_require__("./node_modules/@angular/router/@angular/router.es5.js");
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
    function ThreadPageComponent(api, router, route) {
        this.api = api;
        this.router = router;
        this.route = route;
        this.boardKey = '';
        this.threadKey = '';
        this.data = { posts: [], thread: { name: '', description: '' } };
    }
    ThreadPageComponent.prototype.ngOnInit = function () {
        var _this = this;
        this.route.params.subscribe(function (res) {
            _this.boardKey = res['board'];
            _this.threadKey = res['thread'];
            _this.open(_this.boardKey, _this.threadKey);
        });
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
        console.warn('open:', master);
        this.api.getThreadpage(master, ref).then(function (data) {
            console.warn('get threads2:', data);
            _this.data = data;
            console.log('this data:', _this.data);
        });
    };
    return ThreadPageComponent;
}());
ThreadPageComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_14" /* Component */])({
        selector: 'threadPage',
        template: __webpack_require__("./src/components/threadPage/threadPage.html"),
        styles: [__webpack_require__("./src/components/threadPage/threadPage.css")]
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__providers__["a" /* ApiService */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__providers__["a" /* ApiService */]) === "function" && _a || Object, typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_2__angular_router__["b" /* Router */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__angular_router__["b" /* Router */]) === "function" && _b || Object, typeof (_c = typeof __WEBPACK_IMPORTED_MODULE_2__angular_router__["c" /* ActivatedRoute */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__angular_router__["c" /* ActivatedRoute */]) === "function" && _c || Object])
], ThreadPageComponent);

var _a, _b, _c;
//# sourceMappingURL=threadPage.js.map

/***/ }),

/***/ "./src/components/threads/threads.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("./node_modules/css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, "threads {\n  width: 100%;\n}\n\n.threads {\n  width: 100%;\n}\n.threads .item {\n  display: block;\n  width: 100%;\n  min-height: 200px;\n  text-decoration: none;\n  /*border-top: .1rem solid #ccc;*/\n  border-bottom: .1rem solid #ccc;\n  cursor: pointer;\n}\n.threads .item .description {\n  color: #999;\n}", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "./src/components/threads/threads.html":
/***/ (function(module, exports) {

module.exports = "<div class=\"threads\">\n  <div class=\"container-fluid\">\n    <!--<ol class=\"breadcrumb\">\n      <li><a outerLink=\"/\">All Board</a></li>\n      <li><a outerLink=\"/\">Home</a></li>\n    </ol>-->\n    <div *ngFor=\"let item of threads\" class=\"item\" (click)=\"open(item.master_board,item.ref)\">\n      <h3 class=\"name\">{{item.name}}</h3>\n      <p class=\"description\">{{item.description}}</p>\n    </div>\n    <h3 class=\"boardNot\" *ngIf=\"threads?.length == 0\">Not Found Threads</h3>\n\n  </div>\n</div>\n"

/***/ }),

/***/ "./src/components/threads/threads.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("./node_modules/@angular/core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__providers__ = __webpack_require__("./src/providers/index.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_router__ = __webpack_require__("./node_modules/@angular/router/@angular/router.es5.js");
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
    function ThreadsComponent(api, router, route) {
        this.api = api;
        this.router = router;
        this.route = route;
        this.thread = new __WEBPACK_IMPORTED_MODULE_0__angular_core__["F" /* EventEmitter */]();
        this.threads = [];
    }
    ThreadsComponent.prototype.ngOnInit = function () {
        var _this = this;
        this.route.params.subscribe(function (res) {
            _this.start(res['board']);
        });
    };
    ThreadsComponent.prototype.start = function (key) {
        var _this = this;
        this.api.getThreads(key).then(function (data) {
            console.warn('get threads:', data);
            _this.threads = data;
        });
    };
    ThreadsComponent.prototype.open = function (master, ref) {
        this.router.navigate(['p', { board: master, thread: ref }], { relativeTo: this.route });
        // this.thread.emit({ master: master, ref: ref });
    };
    return ThreadsComponent;
}());
__decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_5" /* Output */])(),
    __metadata("design:type", typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_0__angular_core__["F" /* EventEmitter */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_0__angular_core__["F" /* EventEmitter */]) === "function" && _a || Object)
], ThreadsComponent.prototype, "thread", void 0);
ThreadsComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_14" /* Component */])({
        selector: 'threads',
        template: __webpack_require__("./src/components/threads/threads.html"),
        styles: [__webpack_require__("./src/components/threads/threads.css")],
        encapsulation: __WEBPACK_IMPORTED_MODULE_0__angular_core__["q" /* ViewEncapsulation */].None
    }),
    __metadata("design:paramtypes", [typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_1__providers__["a" /* ApiService */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__providers__["a" /* ApiService */]) === "function" && _b || Object, typeof (_c = typeof __WEBPACK_IMPORTED_MODULE_2__angular_router__["b" /* Router */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__angular_router__["b" /* Router */]) === "function" && _c || Object, typeof (_d = typeof __WEBPACK_IMPORTED_MODULE_2__angular_router__["c" /* ActivatedRoute */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__angular_router__["c" /* ActivatedRoute */]) === "function" && _d || Object])
], ThreadsComponent);

var _a, _b, _c, _d;
//# sourceMappingURL=threads.js.map

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

/***/ "./src/providers/api/api.service.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("./node_modules/@angular/core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_http__ = __webpack_require__("./node_modules/@angular/http/@angular/http.es5.js");
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
    function ApiService(http) {
        this.http = http;
        this.base_url = 'http://127.0.0.1:7410/api/';
    }
    ApiService.prototype.getThreads = function (key) {
        var data = new FormData();
        data.append('board', key);
        return this.handlePost(this.base_url + 'get_threads', data);
    };
    ApiService.prototype.getBoards = function () {
        return this.handleGet(this.base_url + 'get_boards');
    };
    ApiService.prototype.getPosts = function (masterKey, sub) {
        var form = new FormData();
        form.append('board', masterKey);
        form.append('thread', sub);
        return this.handlePost(this.base_url + 'get_posts', form);
    };
    // get_threadpage
    ApiService.prototype.getThreadpage = function (masterKey, sub) {
        var form = new FormData();
        form.append('board', masterKey);
        form.append('thread', sub);
        return this.handlePost(this.base_url + 'get_threadpage', form);
    };
    ApiService.prototype.addBoard = function (data) {
        return this.handlePost(this.base_url + 'new_board', data);
    };
    ApiService.prototype.addThread = function (data) {
        return this.handlePost(this.base_url + 'new_thread', data);
    };
    ApiService.prototype.addPost = function (data) {
        return this.handlePost(this.base_url + 'new_post', data);
    };
    ApiService.prototype.changeThread = function (data) {
        return this.handlePost(this.base_url + 'import_thread', data);
    };
    ApiService.prototype.handlePost = function (url, data) {
        var _this = this;
        return new Promise(function (resolve, reject) {
            _this.http.post(url, data).subscribe(function (res) {
                if (res.status != 200) {
                    reject(new Error('Request Fail'));
                }
                resolve(res.json());
            }, function (err) {
                reject(err);
            });
        });
    };
    ApiService.prototype.handleGet = function (url) {
        var _this = this;
        return new Promise(function (resolve, reject) {
            _this.http.get(url).subscribe(function (res) {
                resolve(res.json());
            }, function (err) {
                reject(err);
            });
        });
    };
    return ApiService;
}());
ApiService = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["c" /* Injectable */])(),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__angular_http__["b" /* Http */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__angular_http__["b" /* Http */]) === "function" && _a || Object])
], ApiService);

var _a;
//# sourceMappingURL=api.service.js.map

/***/ }),

/***/ "./src/providers/index.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__api_api_service__ = __webpack_require__("./src/providers/api/api.service.ts");
/* harmony namespace reexport (by used) */ __webpack_require__.d(__webpack_exports__, "a", function() { return __WEBPACK_IMPORTED_MODULE_0__api_api_service__["a"]; });

//# sourceMappingURL=index.js.map

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
    { path: '', component: __WEBPACK_IMPORTED_MODULE_2__components__["a" /* BoardsListComponent */] },
    {
        path: 'threads', children: [
            { path: '', component: __WEBPACK_IMPORTED_MODULE_2__components__["b" /* ThreadsComponent */] },
            { path: 'p', component: __WEBPACK_IMPORTED_MODULE_2__components__["c" /* ThreadPageComponent */] },
        ]
    },
    // { path: 'threads', component: ThreadsComponent },
    { path: 'add', component: __WEBPACK_IMPORTED_MODULE_2__components__["d" /* AddComponent */] },
    { path: '**', redirectTo: '' }
];
var AppRouterRoutingModule = (function () {
    function AppRouterRoutingModule() {
    }
    return AppRouterRoutingModule;
}());
AppRouterRoutingModule = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["b" /* NgModule */])({
        imports: [__WEBPACK_IMPORTED_MODULE_1__angular_router__["a" /* RouterModule */].forRoot(routes)],
        exports: [__WEBPACK_IMPORTED_MODULE_1__angular_router__["a" /* RouterModule */]],
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