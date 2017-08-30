import { Component, HostBinding, OnInit, ViewEncapsulation, ViewChild, TemplateRef, AfterViewInit } from '@angular/core';
import { ApiService, CommonService, UserService, AllBoards, Alert, Popup, Dialog, LoadingService } from '../../providers';
import { Board } from '../../providers/api/msg';
import { Router } from '@angular/router';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { AlertComponent } from '../alert/alert.component';
import { FabComponent } from '../fab/fab.component';
import { slideInLeftAnimation } from '../../animations/router.animations';
import { flyInOutAnimation } from '../../animations/common.animations';

@Component({
  selector: 'app-boardslist',
  templateUrl: 'boards-list.component.html',
  styleUrls: ['boards-list.scss'],
  encapsulation: ViewEncapsulation.None,
  animations: [slideInLeftAnimation, flyInOutAnimation],
})
export class BoardsListComponent implements OnInit, AfterViewInit {
  @HostBinding('@routeAnimation') routeAnimation = true;
  @HostBinding('style.display') display = 'block';
  @ViewChild('fab') fabBtnTemplate: TemplateRef<any>;
  sort = 'asc';
  isRoot = false;
  boards: Array<Board> = [];
  remoteBoards: Array<Board> = [];
  subscribeForm = new FormGroup({
    address: new FormControl('', Validators.required),
    board: new FormControl('', Validators.required),
  });
  addressForm = new FormGroup({
    ip: new FormControl(''),
    port: new FormControl('', Validators.required),
  });
  addForm = new FormGroup({
    name: new FormControl('', Validators.required),
    body: new FormControl('', Validators.required),
    seed: new FormControl({ value: '', disabled: true }, Validators.required),
    submission_addresses: new FormControl('')
  });
  tmpBoard: Board = null;
  regexpStr = new RegExp('<br\s*/?>', 'g');
  replaceStr = '';
  constructor(private api: ApiService,
    private user: UserService,
    private router: Router,
    private modal: NgbModal,
    public common: CommonService,
    private pop: Popup,
    private alert: Alert,
    private dialog: Dialog,
    private loading: LoadingService) {
  }

  ngOnInit(): void {
    this.getBoards();
    this.api.getStats().subscribe(status => {
      this.isRoot = status.node_is_master;
    });

  }
  ngAfterViewInit() {
    // setTimeout(() => {
    //   this.dialog.open();
    // }, 10);
    this.pop.open(this.fabBtnTemplate);
  }

  setSort() {
    this.sort = this.sort === 'desc' ? 'asc' : 'desc';
  }

  getBoards() {
    this.api.getBoards().subscribe((allBoards: AllBoards) => {
      if (!allBoards.okay || allBoards.data.master_boards.length <= 0) {
        return;
      }
      this.boards = allBoards.data.master_boards;
      this.remoteBoards = allBoards.data.remote_boards;
    });
  }

  addAddress(content: any, key: string) {
    this.addressForm.reset();
    if (key === '') {
      this.alert.error({ content: 'The Key can not be empty!!!' });
      return;
    }
    this.modal.open(content, { windowClass: 'multi-modal' }).result.then((reslut) => {
      if (reslut) {
        if (!this.addressForm.valid) {
          this.alert.error({ content: 'The Port Or Url can not be empty!!!' });
          return;
        }
        const data = new FormData();
        data.append('board_public_key', key);
        let ip = this.addressForm.get('ip').value;
        if (ip === '' || !ip) {
          ip = '[::]:'
        } else {
          ip = ip + ':';
        }
        data.append('address', ip + this.addressForm.get('port').value);
        this.loading.start();
        this.api.newSubmissionAddress(data).subscribe(res => {
          this.tmpBoard.submission_addresses = res.data.board.submission_addresses;
          this.loading.close();
        });
      }
    });
  }
  openURL(ev: Event) {
    ev.preventDefault();
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    window.open(ev.target['href'], '_blank');
  }
  openInfo(ev: Event, board: Board, content: any) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    if (!board) {
      // this.common.showErrorAlert('Failed to get info!!');
      return;
    }
    this.tmpBoard = board;
    this.modal.open(content, { size: 'lg' });
  }
  trackBoards(index, board) {
    return board ? board.public_key : undefined;
  }
  openAdd(content) {
    // this.pop.open(content);
    this.addForm.reset();
    this.api.newSeed().subscribe(seed => {
      this.addForm.patchValue({ seed: seed.data });
      this.modal.open(content).result.then((result) => {
        if (result === true) {
          if (!this.addForm.valid) {
            this.alert.error({ content: 'Parameter error' });
            return;
          }
          const data = new FormData();
          data.append('seed', this.addForm.get('seed').value);
          data.append('name', this.addForm.get('name').value.trim());
          data.append('body', this.common.replaceURL(this.common.replaceHtmlEnter(this.addForm.get('body').value)));
          data.append('submission_addresses', this.addForm.get('submission_addresses').value);
          this.api.addBoard(data).subscribe(res => {
            this.getBoards();
            this.alert.success({ content: 'Added Successfully' });
          });
        }
      }, err => {
      });
    }, err => {
      // this.common.showErrorAlert('Unable to create,Please try again later');
    })
  }

  delAddress(ev: Event, key: string, address: string) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    if (key === '' || address === '') {
      // this.common.showErrorAlert('The key and address can not be empty');
      return;
    }
    const data = new FormData();
    data.append('board_public_key', key);
    data.append('address', address);
    const modalRef = this.modal.open(AlertComponent, { windowClass: 'multi-modal' });
    modalRef.componentInstance.title = 'Delete Address';
    modalRef.componentInstance.body = 'Do you delete the address?';
    modalRef.result.then(result => {
      if (result) {
        this.loading.start();
        this.api.delSubmissionAddress(data).subscribe(res => {
          this.tmpBoard.submission_addresses = res.data.board.submission_addresses;
          this.loading.close();
        });
      }
    }, err => { });
  }

  subscribe(content: any) {
    this.subscribeForm.reset();
    this.modal.open(content).result.then(result => {
      if (result) {
        if (!this.subscribeForm.valid) {
          // this.common.showErrorAlert('The Board Key Or Address can not be empty!!!');
          return;
        }
        const data = new FormData();
        data.append('address', this.subscribeForm.get('address').value);
        data.append('board', this.subscribeForm.get('board').value);
        this.api.subscribe(data).subscribe(isOk => {
          if (isOk) {
            // this.common.showSucceedAlert('Subscribed successfully');
            this.getBoards();
          }
        });
      }
    }, err => { });
  }

  delBoard(ev: Event, boardKey: string) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    const modalRef = this.modal.open(AlertComponent);
    modalRef.componentInstance.title = 'Delete Board';
    modalRef.componentInstance.body = 'Do you delete the board?';
    modalRef.result.then(result => {
      if (result) {
        if (boardKey === '') {
          this.alert.error({ content: 'Delete failed' });
          return;
        }
        const data = new FormData();
        data.append('board_public_key', boardKey);
        this.api.delBoard(data).subscribe(res => {
          if (res.okay) {
            this.alert.success({ content: 'Delete successfully' });
            this.boards = res.data.master_boards;
          }
        });
      }
    }, err => { })
  }

  openThreads(ev: Event, key: string) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    if (!key) {
      this.alert.error({ content: 'Abnormal parameters!!!' });
      return;
    }
    this.router.navigate(['/threads', { boardKey: key }]);
  }
}
