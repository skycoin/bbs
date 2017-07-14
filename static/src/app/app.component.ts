import { Component, OnInit, ViewChild } from '@angular/core';
import { ApiService, CommonService, UserService } from '../providers';
import { LoadingComponent, FixedButtonComponent } from '../components';
import { NavigationStart, Router } from '@angular/router';
import 'rxjs/add/operator/filter';


@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit {
  @ViewChild(LoadingComponent) loading: LoadingComponent;
  @ViewChild(FixedButtonComponent) fb: FixedButtonComponent;
  public title = 'app';
  public name = '';
  public isMasterNode = false;

  constructor(
    private api: ApiService,
    private user: UserService,
    private router: Router,
    public common: CommonService) {
  }

  ngOnInit() {
    this.common.fb = this.fb;
    this.common.loading = this.loading;
    this.user.getCurrent().subscribe(user => {
      this.name = user.alias;
    });
    this.api.getStats().subscribe(res => {
      this.isMasterNode = res.node_is_master;
    });
    this.router.events.filter(ev => ev instanceof NavigationStart).subscribe(ev => {
      this.common.topBtn = false;
    });
  }
}
