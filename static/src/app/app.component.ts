import { Component, OnInit, ViewChild } from '@angular/core';
import { ApiService } from '../providers';
import { LoadingComponent } from '../components';
import { UserService, User, CommonService } from '../providers';
import { Router, NavigationStart } from '@angular/router';
import 'rxjs/add/operator/filter';
@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  @ViewChild(LoadingComponent) loading: LoadingComponent;

  title = 'app';
  name = '';
  isMasterNode = false;
  constructor(
    private api: ApiService,
    private user: UserService,
    private router: Router,
    public common: CommonService) {
  }
  ngOnInit() {
    // const el = document.querySelector('.container');
    // Ps.initialize(el);
    this.common.loading = this.loading;
    this.user.getCurrent().subscribe(user => {
      this.name = user.alias;
    });
    this.api.getStats().subscribe(res => {
      this.isMasterNode = res.node_is_master;
    });
    this.router.events.filter(ev => ev instanceof NavigationStart).subscribe(ev => {
      this.common.topBtn = false;
    })
  }
}
