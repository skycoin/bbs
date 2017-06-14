import { Component, OnInit, ViewEncapsulation, Output, EventEmitter, HostBinding } from '@angular/core';
import { ApiService, Thread, CommonService, Board } from "../../providers";
import { Router, ActivatedRoute } from "@angular/router";
import { FormControl, FormGroup } from '@angular/forms';
import { NgbModal } from "@ng-bootstrap/ng-bootstrap";
import { slideInLeftAnimation } from "../../animations/router.animations";

@Component({
  selector: 'threads',
  templateUrl: 'threads.html',
  styleUrls: ['threads.css'],
  encapsulation: ViewEncapsulation.None,
  animations: [slideInLeftAnimation]
})

export class ThreadsComponent implements OnInit {
  @HostBinding('@routeAnimation') routeAnimation = true;
  @HostBinding('style.display') display = 'block';
  @Output() thread: EventEmitter<{ master: string, ref: string }> = new EventEmitter();
  constructor(
    private api: ApiService,
    private router: Router,
    private route: ActivatedRoute,
    private modal: NgbModal,
    private common: CommonService) { }
  private threads: Array<Thread> = [];
  private importBoards: Array<Board> = [];
  private importBoardKey: string = '';
  private boardKey: string = '';
  private board: Board = null;
  private addForm = new FormGroup({
    description: new FormControl(),
    name: new FormControl()
  });
  ngOnInit() {
    this.route.params.subscribe(res => {
      // this.url = res['url'];
      this.boardKey = res['board'];
      this.init();
    })
  }
  initThreads(key) {
    let data = new FormData();
    data.append('board',key);
    this.api.getThreads(data).subscribe(threads => {
      this.threads = threads;
    });
  }
  init() {
    let data = new FormData();
    data.append('board', this.boardKey);
    this.api.getBoardPage(data).subscribe(data => {
      this.board = data.board;
      this.threads = data.threads;
    })
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
  openImport(content: any, threadKey: string) {
    if (this.importBoards.length <= 0) {
      this.api.getBoards().subscribe(boards => {
        this.importBoards = boards;
        this.importBoardKey = boards[0].public_key;
      });
    }
    this.modal.open(content, { size: 'lg' }).result.then(result => {
      if (result) {
        if (this.importBoardKey) {
          let data = new FormData();
          data.append('from_board', this.boardKey);
          data.append('thread', threadKey);
          data.append('to_board', this.importBoardKey);
          this.api.importThread(data).subscribe(res => {
            console.log('transfer thread:', res);
            this.common.showAlert('successfully', 'success', 3000);
            this.initThreads(this.boardKey);
          })
        }
      }
    },err => {});
  }
}