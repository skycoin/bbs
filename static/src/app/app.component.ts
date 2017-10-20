import { Component, OnInit, ViewChild, HostListener, QueryList, ViewChildren, ElementRef, Renderer } from '@angular/core';
import {
  ApiService,
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
  public title = 'app';
  public name = '';
  public isMasterNode = false;
  userName = 'LogIn';
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
  userList = []
  selectUser = '';
  constructor(
    private api: ApiService,
    public common: CommonService,
    private alert: Alert,
    private loading: LoadingService,
    private pop: Popup,
    private render: Renderer) {
  }

  ngOnInit() {
    // const seed = cipher.generateSeed();
    // const keypair = cipher.generateKeyPair(seed)
    // const encryptText = crypto.AES.encrypt('nyf', '123456')
    // console.log('crypto:', encryptText.toString());
    // const bytes = crypto.AES.decrypt(encryptText.toString(), '123456');
    // const plaintext = bytes.toString(crypto.enc.Utf8);
    // console.log('plaintext:', plaintext);
    // this.api.getAllUser().subscribe(res => {
    //   this.autoAilas = res.data.users;
    //   this._orginAutoAilas = this.autoAilas;
    // })
    this.common.fb = this.fb;
    // this.api.getStats().subscribe(stats => {
    //   this.isMasterNode = stats.node_is_master;
    // });
    Observable.timer(10).subscribe(() => {
      this.pop.open(ToTopComponent, { isDialog: false });
    });
    // this.api.getSessionInfo().subscribe(info => {
    //   if (info.okay) {
    //     if (info.data.session && info.data.logged_in) {
    //       this.isLogIn = info.data.logged_in;
    //       this.userName = info.data.session.user.alias;
    //       this.userPublicKey = info.data.session.user.public_key;
    //       ApiService.userInfo = info.data.session.user;
    //     }
    //   }
    // })
  }
  openRegister(ev: Event, content: any) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.pop.open(content, { isDialog: true, canClickBackdrop: false }).result.then(result => {
    }, err => { });
  }
  openLogin(ev: Event, content: any) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.pop.open(content, { isDialog: true, canClickBackdrop: false }).result.then(result => {
    }, err => { });
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
    this.api.logout().subscribe(res => {
      if (res.okay) {
        this.userName = 'LogIn';
        this.isLogIn = res.data.logged_in;
        this.userMenu = false;
      }
    })
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
