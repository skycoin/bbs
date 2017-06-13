import { Component, OnInit, ViewEncapsulation, Output, EventEmitter } from '@angular/core';
import { ApiService, Thread, CommonService, Board } from "../../providers";
import { Router, ActivatedRoute } from "@angular/router";
import { FormControl, FormGroup } from '@angular/forms';
import { NgbModal, ModalDismissReasons } from "@ng-bootstrap/ng-bootstrap";

@Component({
  selector: 'threads',
  templateUrl: 'threads.html',
  styleUrls: ['threads.css'],
  encapsulation: ViewEncapsulation.None
})

export class ThreadsComponent implements OnInit {
  @Output() thread: EventEmitter<{ master: string, ref: string }> = new EventEmitter();
  constructor(
    private api: ApiService,
    private router: Router,
    private route: ActivatedRoute,
    private modal: NgbModal,
    private common: CommonService) { }
  threads: Array<Thread> = [];
  transferBoards: Array<Board> = [];
  transferBoardKey: string = '';
  boardKey: string = '';
  url: string = '';
  addForm = new FormGroup({
    description: new FormControl(),
    name: new FormControl()
  });
  ngOnInit() {
    this.route.params.subscribe(res => {
      this.url = res['url'];
      this.boardKey = res['board'];
      this.start(this.boardKey);
    })
  }
  start(key) {
    this.api.getThreads(key).subscribe(threads => {
      this.threads = threads;
    });
  }
  openAdd(content) {
    this.modal.open(content).result.then((result) => {
      if (result) {
        let data = new FormData();
        data.append('board', this.boardKey);
        data.append('description', this.addForm.get('description').value);
        data.append('name', this.addForm.get('name').value);
        this.api.addThread(data).subscribe(thread => {
          this.threads.unshift(thread);
          this.common.showAlert('Added successfully', 'success', 3000);
        });
      }
    }, err => { });
  }
  open(master, ref: string) {
    this.router.navigate(['p', { board: master, thread: ref }], { relativeTo: this.route });
  }
  openTransfer(content: any, threadKey: string) {
    if (this.transferBoards.length <= 0) {
      this.api.getBoards().subscribe(boards => {
        this.transferBoards = boards;
      });
    }
    this.modal.open(content, { size: 'lg' }).result.then(result => {
      if (result) {
        if (this.transferBoardKey) {
          let data = new FormData();
          data.append('from_board', this.boardKey);
          data.append('thread', threadKey);
          data.append('to_board', this.transferBoardKey);
          this.api.importThread(data).subscribe(res => {
            console.log('transfer thread:',res);
            this.common.showAlert('successfully', 'success', 3000);
            this.start(this.boardKey);
          })
        }
      }
    });
  }
}