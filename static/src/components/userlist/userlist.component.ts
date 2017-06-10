import { Component, OnInit } from '@angular/core';
import { UserService, User } from "../../providers";

@Component({
  selector: 'app-userlist',
  templateUrl: './userlist.component.html',
  styleUrls: ['./userlist.component.css']
})
export class UserlistComponent implements OnInit {
  userlist: Array<User> = [];
  constructor(private user: UserService) { }
  ngOnInit() {
    this.user.getAll().subscribe(userlist => {
      this.userlist = userlist;
    })
  }
  remove(ev: Event, key: string) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    let data = new FormData();
    data.append('user', key);
    this.user.remove(data).subscribe(isOk => {
      if (isOk) {
        this.userlist = [];
        this.user.getAll().subscribe(userlist => {
          this.userlist = userlist;
        })
      }
    })
  }
}
