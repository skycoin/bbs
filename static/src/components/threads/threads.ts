import { Component, OnInit, ViewEncapsulation, Output, EventEmitter, HostBinding } from '@angular/core';
import { ApiService, Thread, CommonService, Board } from '../../providers';
import { Router, ActivatedRoute } from '@angular/router';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { slideInLeftAnimation } from '../../animations/router.animations';

@Component({
  selector: 'app-threads',
  templateUrl: 'threads.html',
  styleUrls: ['threads.scss'],
  encapsulation: ViewEncapsulation.None,
  animations: [slideInLeftAnimation]
})

export class ThreadsComponent implements OnInit {
  @HostBinding('@routeAnimation') routeAnimation = true;
  @HostBinding('style.display') display = 'block';
  threads: Array<Thread> = [];
  importBoards: Array<Board> = [];
  importBoardKey = '';
  boardKey = '';
  board: Board = null;
  isRoot = false;
  tmpThread: Thread = null;
  public addForm = new FormGroup({
    description: new FormControl('', Validators.required),
    name: new FormControl('', Validators.required)
  });
  @Output() thread: EventEmitter<{ master: string, ref: string }> = new EventEmitter();
  constructor(
    private api: ApiService,
    private router: Router,
    private route: ActivatedRoute,
    private modal: NgbModal,
    private common: CommonService) { }
  ngOnInit() {
    this.route.params.subscribe(res => {
      this.boardKey = res['board'];
      this.init();
      this.api.getStats().subscribe(root => {
        this.isRoot = root;
      });
    })
  }
  initThreads(key) {
    if (key === '') {
      this.common.showErrorAlert('Parameter error!!!');
      return;
    }
    const data = new FormData();
    data.append('board', key);
    this.api.getThreads(data).subscribe(threads => {
      this.threads = threads;
    });
  }
  init() {
    this.common.loading.start();
    const data = new FormData();
    data.append('board', this.boardKey);
    this.api.getBoardPage(data).subscribe(res => {
      this.board = res.board;
      this.threads = res.threads;
      this.common.loading.close();
    })
  }
  openInfo(ev: Event, thread: Thread, content: any) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    this.tmpThread = thread;
    this.modal.open(content, { size: 'lg' });
  }
  openAdd(content) {
    this.addForm.reset();
    this.modal.open(content).result.then((result) => {
      if (result) {
        if (!this.addForm.valid) {
          this.common.showErrorAlert('Parameter error!!!');
          return;
        }
        const data = new FormData();
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
    if (master === '' || ref === '') {
      this.common.showErrorAlert('Parameter error!!!');
      return;
    }
    this.router.navigate(['p', { board: master, thread: ref }], { relativeTo: this.route });
  }
  openImport(ev: Event, threadKey: string, content: any) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    if (!this.isRoot) {
      this.common.showErrorAlert('Only Master Nodes Can Import', 3000);
      return;
    }
    let tmp: Array<Board> = [];
    this.api.getBoards().subscribe(boards => {
      tmp = boards;
      tmp.forEach((el, index) => {
        if (el.public_key === this.boardKey) {
          tmp.splice(index, 1);
        }
      });
      if (tmp.length <= 0) {
        this.common.showErrorAlert('None are suitable');
        return;
      }
      this.importBoards = tmp;
      this.importBoardKey = tmp[0].public_key;
      this.modal.open(content, { size: 'lg' }).result.then(result => {
        if (result) {
          if (this.importBoardKey) {
            const data = new FormData();
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
      }, err => { });
    });
  }
}
