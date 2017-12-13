import {
  Component,
  OnInit,
  ViewChild,
  HostListener,
  QueryList,
  ViewChildren,
  ElementRef,
  Renderer,
  TemplateRef
} from '@angular/core';
import {
  ApiService,
  UserService,
  CommonService,
  Alert,
  LoadingService,
  Popup,
  FollowPageDataInfo,
  FollowPage,
  User,
  LoginInfo
} from '../providers';
import { FixedButtonComponent } from '../components';
import { ToTopComponent } from '../components';
import 'rxjs/add/operator/filter';
import { bounceInAnimation, tabLeftAnimation, tabRightAnimation } from '../animations/common.animations';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/observable/timer'
import 'rxjs/add/operator/do'
import { FormControl, FormGroup, Validators } from '@angular/forms';

// import * as crypto from 'crypto-js';

declare var cipher: any;

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
  animations: [bounceInAnimation, tabLeftAnimation, tabRightAnimation]
})
export class AppComponent implements OnInit {
  @ViewChild(FixedButtonComponent) fb: FixedButtonComponent;
  @ViewChildren('aliasItem') aliasItems: QueryList<ElementRef>;
  @ViewChild('login') loginTemplate: TemplateRef<any>;
  @ViewChild('create') createTemplate: TemplateRef<any>;
  public title = 'app';
  public name = '';
  public isMasterNode = false;
  userName = '';
  alias = '';
  seed = '';
  isLogIn = false;
  navBarBg = 'default-navbar';
  userMenu = false;
  showLoginBox = false;
  userPublicKey = '';
  boardKey = '';
  loginBox = true;
  registerBox = false;
  _orginAutoAilas: Array<User>;
  autoAilas: Array<User>;
  showAilas = false;
  autoAliasIndex = -1;
  userFollow: FollowPageDataInfo = {};
  userList = [];
  selectUser = '';
  selectUserPass = '';
  showPassword = false;
  hasAlias = false;
  loginForm = new FormGroup({
    user: new FormControl('', Validators.required),
    pass: new FormControl('', Validators.required)
  });
  createForm = new FormGroup({
    alias: new FormControl('', Validators.required),
    seed: new FormControl('', Validators.required),
    password: new FormControl('', Validators.required),
    confirmPassword: new FormControl('', Validators.required),
  });
  authPassword = false;
  constructor(
    private api: ApiService,
    public common: CommonService,
    private alert: Alert,
    private loading: LoadingService,
    private pop: Popup,
    private render: Renderer,
    private user: UserService,
  ) {
  }

  ngOnInit() {
    this.common.fb = this.fb;
    this.userList = this.user.getUserList();
    const loginInfo = this.user.getTmpItem();
    if (loginInfo) {
      this.userName = loginInfo.name;
    }
    // if (!loginInfo.data) {
    Observable.timer(10).subscribe(() => {
      this.loginForm.patchValue({ user: this.userName })
      this.pop.open(this.loginTemplate, { isDialog: true, canClickBackdrop: false }).result.then(result => {
        if (result === true) {
          const user = this.loginForm.get('user').value;
          const hash = this.user.getItem(user);
          this.user.decrypt(hash, this.loginForm.get('pass').value).subscribe((info: any) => {
            if (info) {
              this.userName = user;
              this.user.setTmpItem(this.userName);
              this.user.loginInfo = info;
              this.alert.success({ content: 'Authentication is successful' });
            } else {
              this.alert.error({ content: 'Password error or system is busy, please try again later' });
            }
          })
        } else if (result === 'create') {
          this.createForm.reset();
          this.hasAlias = false;
          this.authPassword = false;
          this.createForm.patchValue({ seed: this.user.newSeed() });
          this.pop.open(this.createTemplate, { isDialog: true, canClickBackdrop: false }).result.then(createRsult => {
            if (createRsult) {
              this.user.newKeyPair(this.createForm.get('seed').value).then(json => {
                const password = this.createForm.get('password').value;
                const alias = this.createForm.get('alias').value;
                this.user.encrypt(JSON.stringify(json), password).do(() => {
                }).subscribe(data => {
                  this.user.setItem(alias, data);
                  this.userList = this.user.getUserList();
                  this.alert.success({ content: 'Successfully created' });
                })
              })
            }
          }, err => {
            console.log('error:', err);
          });
        }
      });
    });
    // } else {
    //   this.user.loginInfo = loginInfo.data;
    // }
    Observable.timer(10).subscribe(() => {
      this.pop.open(ToTopComponent, { isDialog: false });
    });
  }
  test(content) {
    this.loading.start();
  }
  isShowPassword(ev: Event, input: any) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.showPassword = !this.showPassword;
    if (this.showPassword) {
      input.type = 'text';
    } else {
      input.type = 'password';
    }
  }
  checkPassword() {
    const pass = this.createForm.get('password').value;
    const confirmPassword = this.createForm.get('confirmPassword').value;
    this.authPassword = !(pass === confirmPassword);
  }

  checkAlias() {
    const item = this.user.getItem(this.createForm.get('alias').value)
    if (item) {
      this.hasAlias = true;
      return
    }
    this.hasAlias = false;
  }
  openRegister(ev: Event, content: any) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.showPassword = false;
    this.pop.open(content, { isDialog: true, canClickBackdrop: false }).result.then(result => {
    }, err => { });
  }
  openLogin(ev: Event, login, create: any, isSwitch: boolean = false) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    if (this.user.loginInfo && !isSwitch) {
      this.userMenu = !this.userMenu;
      return;
    }
    this.userMenu = false;
    this.loginForm.reset();
    this.showPassword = false;
    // if (this.userName) {
    //   this.loginForm.patchValue({ user: this.userName })
    // }
    this.pop.open(login, { isDialog: true, canClickBackdrop: false }).result.then(result => {
      if (result === 'create') {
        this.createForm.reset();
        this.hasAlias = false;
        this.createForm.patchValue({ seed: this.user.newSeed() });
        this.pop.open(create, { isDialog: true, canClickBackdrop: false }).result.then(createRsult => {
          if (createRsult) {
            this.user.newKeyPair(this.createForm.get('seed').value).then(json => {
              const password = this.createForm.get('password').value;
              const alias = this.createForm.get('alias').value;
              this.user.encrypt(JSON.stringify(json), password).do(() => {
              }).subscribe(data => {
                this.user.setItem(alias, data);
                this.userList = this.user.getUserList();
                this.loading.close();
              })
            })
          }
        }, err => {
          console.log('error:', err);
        });
      } else if (result === true) {
        const user = this.loginForm.get('user').value;
        const hash = this.user.getItem(user);
        this.user.decrypt(hash, this.loginForm.get('pass').value).subscribe((loginInfo: any) => {
          if (loginInfo) {
            this.userName = user;
            this.user.setTmpItem(this.userName);
            this.user.loginInfo = loginInfo;
            this.alert.success({ content: 'Authentication is successful' });
          } else {
            this.alert.error({ content: 'Password error or system is busy, please try again later' });
          }
        }, err => {
          console.log('error:', err);
        })
      }
    }, err => { });
  }
  newSeed(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.createForm.patchValue({ seed: this.user.newSeed() });
  }
  selectAlias(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.alias = ev.target['textContent'];
    this.showAilas = false;
  }
  switchTab(tab: string) {
    switch (tab) {
      case 'login':
        this.loginBox = true;
        this.registerBox = false;
        break;
      case 'register':
        this.loginBox = false;
        this.registerBox = true;
        break;
      default:
        break;
    }
  }
  onKeyUp(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    const items = this.aliasItems['_results'];
    const length = this.aliasItems.length - 1;
    switch (ev['keyCode']) {
      case 13:
        if (!this.showAilas) {
          this.login();
        } else {
          this.alias = items[this.autoAliasIndex].nativeElement.textContent;
          this.showAilas = false;
          this.autoAliasIndex = -1;
        }
        break;
      case 38:
        if (this.autoAliasIndex <= -1) {
          this.render.setElementClass(this.aliasItems.last.nativeElement, 'active-alias', true);
          this.autoAliasIndex = length;
        } else {
          this.render.setElementClass(items[this.autoAliasIndex].nativeElement, 'active-alias', false);
          if (this.autoAliasIndex <= 0) {
            this.autoAliasIndex = length;
          } else {
            this.autoAliasIndex -= 1;
          }
          this.render.setElementClass(items[this.autoAliasIndex].nativeElement, 'active-alias', true);
        }
        items[0].nativeElement.parentNode.parentNode.parentNode.scrollTop = this.autoAliasIndex * items[0].nativeElement.clientHeight
        break;
      case 40:
        if (this.autoAliasIndex <= -1) {
          this.render.setElementClass(this.aliasItems.first.nativeElement, 'active-alias', true);
          this.autoAliasIndex = 0;
        } else {
          this.render.setElementClass(items[this.autoAliasIndex].nativeElement, 'active-alias', false);
          if (this.autoAliasIndex >= length) {
            this.autoAliasIndex = 0;
          } else {
            this.autoAliasIndex += 1;
          }
          this.render.setElementClass(items[this.autoAliasIndex].nativeElement, 'active-alias', true);
        }
        items[0].nativeElement.parentNode.parentNode.parentNode.scrollTop = this.autoAliasIndex * items[0].nativeElement.clientHeight
        break;
      default:
        this.showAilas = true;
        this.autoAilas = this._orginAutoAilas;
        this.filterAlias();
        break;
    }
  }
  filterAlias() {
    const tmp: Array<User> = [];
    this.autoAilas.forEach(el => {
      if (el.alias.indexOf(this.alias) > -1) {
        tmp.push(el);
      }
    })
    this.autoAilas = tmp;
    // return this.autoAilas.filter()
  }

  userAction(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    const url = new URL(location.href);
    this.boardKey = url.searchParams.get('boardKey');
    if (this.showLoginBox) {
      this.showLoginBox = false;
      return;
    }
    if (!this.isLogIn) {
      this.alias = '';
      this.seed = '';
      this.showLoginBox = true;
      setTimeout(() => {
        this.autoAilas = this._orginAutoAilas;
        this.showAilas = true;
      }, 300)
      this.loginBox = true;
      this.registerBox = false;
    } else {
      this.showUserMenu();
    }
  }
  login() {
    if (!this.alias) {
      this.alert.error({ content: 'The alias can not empty!' });
      return;
    }
    this.startLogin();
  }
  startLogin() {
    const data = new FormData();
    data.append('alias', this.alias);
    this.api.login(data).subscribe((loginInfo: LoginInfo) => {
      if (loginInfo.okay) {
        this.isLogIn = loginInfo.data.logged_in;
        this.userName = loginInfo.data.session.user.alias;
        this.userPublicKey = loginInfo.data.session.user.public_key;
        this.alias = '';
        ApiService.userInfo = loginInfo.data.session.user;
      }
      this.showLoginBox = false;
    })
  }
  register(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    if (!this.alias) {
      this.alert.error({ content: 'The alias can not empty!' });
      return;
    }
    this.api.newSeed().subscribe(seed => {
      if (seed.okay) {
        const data = new FormData();
        data.append('alias', this.alias);
        data.append('seed', seed.data);
        this.api.newUser(data).subscribe(userData => {
          if (userData.okay) {
            this.startLogin();
          }
        })
      }
    });

  }
  showUserMenu() {
    this.userMenu = !this.userMenu;
  }
  logout(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.userName = 'Login';
    this.user.loginInfo = null;
    this.userMenu = false;
  }

  openFollow(ev: Event, content: any) {
    if (!this.boardKey || !ApiService.userInfo) {
      this.alert.error({ content: 'Please go to a board' });
      return;
    }
    const data = new FormData();
    data.append('board_public_key', this.boardKey);
    data.append('user_public_key', ApiService.userInfo.public_key);
    this.api.getFollowPage(data).subscribe((page: FollowPage) => {
      if (page.okay) {
        this.userFollow = page.data.follow_page;
      }
    })
    this.pop.open(content);
  }
  @HostListener('window:scroll', ['$event'])
  windowScroll(event) {
    const pos = (document.documentElement.scrollTop || document.body.scrollTop) + document.documentElement.offsetHeight;
    const max = document.documentElement.scrollHeight;
    const clientHeight = document.documentElement.clientHeight;
    const distance = max - pos;
    const enableScroll = max - clientHeight - 10;
    if (distance < enableScroll) {
      this.navBarBg = 'after-navbar';
    } else if (distance >= enableScroll) {
      this.navBarBg = 'default-navbar';
    }
  }
}
