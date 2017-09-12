import { Component, OnInit, ViewChild, OnDestroy } from '@angular/core';
import { SocketService, UserService } from '../providers';
import { ImRecentItemComponent } from '../components';

@Component({
  selector: 'app-im',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnDestroy {
  constructor(public ws: SocketService, private user: UserService) {
  }

  ngOnDestroy() {
    this.ws.socket.unsubscribe();
  }
}
