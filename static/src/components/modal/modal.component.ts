import { Component, Input } from '@angular/core';
import { NgbActiveModal } from '@ng-bootstrap/ng-bootstrap';
import { UserService } from "../../providers";

@Component({
  selector: 'app-modal',
  templateUrl: './modal.component.html',
  styleUrls: ['./modal.component.css']
})
export class ModalComponent {
  name: string = '';
  constructor(public activeModal: NgbActiveModal, private user: UserService) { }
  exec(ok: boolean) {
    let data = { name: this.name, ok: ok };
    this.activeModal.close(data);
  }
}
