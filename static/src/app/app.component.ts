import { Component, OnInit, ViewChild, HostListener } from '@angular/core';
import { ApiService, CommonService, Alert, LoadingService, Popup, FollowPageDataInfo, FollowPage } from '../providers';
import { FixedButtonComponent } from '../components';
import { ToTopComponent } from '../components';
import 'rxjs/add/operator/filter';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { bounceInAnimation } from '../animations/common.animations';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
  animations: [bounceInAnimation]
})
export class AppComponent implements OnInit {
  @ViewChild(FixedButtonComponent) fb: FixedButtonComponent;
  public title = 'app';
  public name = '';
  public isMasterNode = false;
  userName = 'LogIn';
  alias = '';
  isLogIn = false;
  navBarBg = 'default-navbar';
  userMenu = false;
  showLoginBox = false;
  userPublicKey = '';
  boardKey = '';
  userFollow: FollowPageDataInfo = {};
  constructor(
    private api: ApiService,
    public common: CommonService,
    private alert: Alert,
    private loading: LoadingService,
    private pop: Popup,
    private modal: NgbModal) {
  }

  ngOnInit() {
    this.common.fb = this.fb;
    this.api.getStats().subscribe(stats => {
      this.isMasterNode = stats.node_is_master;
    });
    this.pop.open(ToTopComponent, false);
    this.api.getSessionInfo().subscribe(info => {
      console.log('user info:', info);
      if (info.okay) {
        if (info.data.session && info.data.logged_in) {
          this.isLogIn = info.data.logged_in;
          this.userName = info.data.session.user.alias;
          this.userPublicKey = info.data.session.user.public_key;
        }
      }
    })
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
      this.showLoginBox = true;
    } else {
      this.showUserMenu();
    }
  }
  login(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    if (!this.alias) {
      return;
    }
    const data = new FormData();
    data.append('alias', this.alias);
    this.api.login(data).subscribe(res => {
      if (res.okay) {
        this.isLogIn = res.data.logged_in;
        this.userName = res.data.session.user.alias;
        this.userPublicKey = res.data.session.user.public_key;
      }
      this.showLoginBox = false;
    })
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
    if (!this.boardKey || !this.userPublicKey) {
      this.alert.error({ content: 'Please go to a board' });
      return;
    }
    const data = new FormData();
    data.append('board_public_key', this.boardKey);
    data.append('user_public_key', this.userPublicKey);
    this.api.getFollowPage(data).subscribe((page: FollowPage) => {
      if (page.okay) {
        this.userFollow = page.data.follow_page;
      }
    })
    this.modal.open(content, { size: 'lg' });
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
