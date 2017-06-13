import { Component, OnInit, ViewChild } from '@angular/core';
import { ApiService } from "../providers";
import { BoardsListComponent, ThreadsComponent, ThreadPageComponent } from "../components";
import { UserService, User, CommonService } from "../providers";
@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  @ViewChild(BoardsListComponent) boards: BoardsListComponent;
  @ViewChild(ThreadsComponent) threads: ThreadsComponent;
  @ViewChild(ThreadPageComponent) threadPage: ThreadPageComponent
  title = 'app';
  name: string = '';
  isMasterNode: boolean = false;
  constructor(private api: ApiService, private user: UserService, public common: CommonService) {
  }
  ngOnInit() {
    this.user.getCurrent().subscribe(user => {
      this.name = user.alias;
    });
    this.api.getStats().subscribe(res => {
      this.isMasterNode = res.node_is_master;
    });
  }
  test() {
    console.log('test');
  }
}
