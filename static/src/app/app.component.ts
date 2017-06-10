import { Component, OnInit, ViewChild } from '@angular/core';
import { ApiService } from "../providers";
import { BoardsListComponent, ThreadsComponent, ThreadPageComponent } from "../components";
import { UserService, User } from "../providers";
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
  constructor(private api: ApiService, private user: UserService) {
  }
  ngOnInit() {
    this.user.getCurrent().subscribe(user => {
      this.name = user.alias;
    })
  }
  test() {
    console.log('test');
  }
}
