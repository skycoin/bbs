import { Component, OnInit, HostBinding } from '@angular/core';
import { UserService, User } from "../../providers";
import { slideInLeftAnimation } from "../../animations/router.animations";
import { NgbModal } from "@ng-bootstrap/ng-bootstrap";


@Component({
  selector: 'app-user',
  templateUrl: './user.component.html',
  styleUrls: ['./user.component.scss'],
  animations: [slideInLeftAnimation]
})
export class UserComponent implements OnInit {
  @HostBinding('@routeAnimation') routeAnimation = true;
  @HostBinding('style.display') display = 'block';
  private user: User = null;
  private switchUserKey: string = '';
  private switchUserList: Array<User> = [];
  constructor(private userService: UserService, private modal: NgbModal) { }

  ngOnInit() {
    this.userService.getCurrent().subscribe(user => {
      this.user = user;
    })
  }
  switchUser(content: any) {
    if (this.switchUserList.length <= 0) {
      this.userService.getAll().subscribe(users => {
        this.switchUserList = users;
        this.switchUserKey = users[0].public_key;
      })
    }
    this.modal.open(content).result.then((result) => {
      if (result) {
        console.log('switchKey:',this.switchUserKey);
        let data = new FormData();
        data.append('user', this.switchUserKey);
        this.userService.setCurrent(data).subscribe(res => {
          location.reload();
        })
      }
    },err => {})
  }
}
