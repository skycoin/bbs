import { Component, OnInit } from '@angular/core';
import { UserService, User } from "../../providers";

@Component({
  selector: 'app-user',
  templateUrl: './user.component.html',
  styleUrls: ['./user.component.css']
})
export class UserComponent implements OnInit {
  private user: User = null;
  closeResult: string;
  constructor(private userService: UserService) { }

  ngOnInit() {
    this.userService.getCurrent().subscribe(user => {
      this.user = user;
    })
  }
}
