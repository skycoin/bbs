import { Component, OnInit, ViewChild } from '@angular/core';
import { ApiService } from "../providers";
import { BoardsListComponent, ThreadsComponent, ThreadPageComponent } from "../components";
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
  constructor(private api: ApiService) {
  }
  ngOnInit() {
    this.api.getBoards();
  }
  openThreads(key: string) {
    this.threads.start(key);
  }
  openThreadpage(data: { master: string, ref: string }) {
    this.threadPage.open(data.master, data.ref);
  }
  test() {
    console.log('test');
  }
}
