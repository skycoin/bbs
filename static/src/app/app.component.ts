import { Component, OnInit, ViewChild, HostListener } from '@angular/core';
import { ApiService, CommonService, UserService, Alert, LoadingService, Popup } from '../providers';
import { FixedButtonComponent } from '../components';
import { ToTopComponent } from '../components';
import 'rxjs/add/operator/filter';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit {
  @ViewChild(FixedButtonComponent) fb: FixedButtonComponent;
  public title = 'app';
  public name = '';
  public isMasterNode = false;
  userName = 'LogIn';
  isLogIn = false;
  navBarBg = 'default-navbar';
  userMenu = false;
  constructor(
    private api: ApiService,
    private user: UserService,
    public common: CommonService,
    private alert: Alert,
    private loading: LoadingService,
    private pop: Popup) {
  }

  ngOnInit() {
    // this.loading.start();
    this.common.fb = this.fb;
    this.api.getStats().subscribe(stats => {
      this.isMasterNode = stats.node_is_master;
    });
    this.pop.open(ToTopComponent, false);
    this.api.getSessionInfo().subscribe(info => {
      if (info.okay) {
        if (info.data.session) {
          this.isLogIn = info.data.logged_in;
          this.userName = info.data.session.user.alias;
        }
      }
    })
  }
  userAction(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    if (!this.isLogIn) {
      this.login();
    } else {
      this.showUserMenu();
    }
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
  login() {
    const data = new FormData();
    data.append('alias', 'angular4');
    this.api.login(data).subscribe(res => {
      if (res.okay) {
        console.log('res:', res);
        this.isLogIn = res.data.logged_in;
        this.userName = res.data.session.user.alias;
      }
    })
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
