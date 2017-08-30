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
  navBarBg = 'default-navbar';
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
    this.pop.open(ToTopComponent);
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
