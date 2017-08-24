import { Component, HostBinding, OnInit } from '@angular/core';
import { CommonService, User, UserService } from '../../providers';
import { slideInLeftAnimation } from '../../animations/router.animations';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';

@Component({
  selector: 'app-user',
  templateUrl: './user.component.html',
  styleUrls: ['./user.component.scss'],
  animations: [slideInLeftAnimation],
})
export class UserComponent implements OnInit {
  @HostBinding('@routeAnimation') routeAnimation = true;
  @HostBinding('style.display') display = 'block';
  user: User = null;
  switchUserKey = '';
  switchUserList: Array<User> = [];

  constructor(private userService: UserService, private modal: NgbModal, public common: CommonService) {
  }

  ngOnInit() {
    this.userService.getCurrent().subscribe(user => {
      this.user = user;
    });
  }

  copy(ev) {
    if (ev) {
      // this.common.showSucceedAlert('Copy Successful');
    } else {
      // this.common.showErrorAlert('Copy Failed');
    }
  }

  switchUser(content: any) {
    if (this.switchUserList.length <= 0) {
      this.userService.getAll().subscribe(users => {
        this.switchUserList = users;
        this.switchUserKey = users[0].public_key;
      });
    }
    this.modal.open(content).result.then((result) => {
      if (result) {
        console.log('switchKey:', this.switchUserKey);
        const data = new FormData();
        data.append('user', this.switchUserKey);
        this.userService.setCurrent(data).subscribe(res => {
          location.reload();
        });
      }
    }, err => {
    });
  }
}
