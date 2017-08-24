import { Component, EventEmitter, HostBinding, OnInit, Output, ViewEncapsulation } from '@angular/core';
import { ApiService, Board, CommonService, Thread, Alert, BoardPage } from '../../providers';
import { ActivatedRoute, Router, ParamMap } from '@angular/router';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { slideInLeftAnimation } from '../../animations/router.animations';
import { flyInOutAnimation } from '../../animations/common.animations';
import 'rxjs/add/operator/switchMap';

@Component({
  selector: 'app-threads',
  templateUrl: 'threads.component.html',
  styleUrls: ['threads.component.scss'],
  encapsulation: ViewEncapsulation.None,
  animations: [slideInLeftAnimation, flyInOutAnimation],
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
  sort = 'asc';

  public addForm = new FormGroup({
    body: new FormControl('', Validators.required),
    title: new FormControl('', Validators.required),
  });
  @Output() thread: EventEmitter<{ master: string, ref: string }> = new EventEmitter();

  constructor(private api: ApiService,
    private router: Router,
    private route: ActivatedRoute,
    private modal: NgbModal,
    private common: CommonService,
    private alert: Alert) {
  }

  ngOnInit() {
    this.route.params.subscribe(params => {
      this.boardKey = params['boardKey'];
      this.init();
    })
    // this.route.params.subscribe(res => {
    //   this.boardKey = res.params['boardKey'];
    //   this.init();
    //   // this.api.getStats().subscribe(root => {
    //   //   this.isRoot = root;
    //   // });
    // });
  }
  trackThreads(index, thread) {
    return thread ? thread.reference : undefined;
  }
  initThreads(key) {
    if (key === '') {
      this.alert.error({ content: 'Parameter error!!!' });
      return;
    }
    const data = new FormData();
    data.append('board', key);
    this.api.getThreads(data).subscribe(threads => {
      this.threads = threads;
    });
  }

  init() {
    const data = new FormData();
    data.append('board_public_key', this.boardKey);
    this.api.getBoardPage(data).subscribe((res: BoardPage) => {
      this.board = res.data.board;
      this.threads = res.data.threads;
    }, err => {
      this.router.navigate(['']);
    });
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
          this.alert.error({ content: 'Parameter error!!!' });
          return;
        }
        const data = new FormData();
        data.append('board_public_key', this.boardKey);
        data.append('body', this.common.replaceHtmlEnter(this.addForm.get('body').value));
        data.append('title', this.addForm.get('title').value);
        this.api.addThread(data).subscribe(threadRes => {
          console.log('thread:', threadRes);
          this.threads = threadRes.data.threads;
          this.alert.success({ content: 'Added successfully' });
        });
      }
    }, err => {
    });
  }

  open(ref: string) {
    if (this.boardKey === '' || ref === '') {
      // this.common.showErrorAlert('Parameter error!!!');
      return;
    }
    this.router.navigate(['/threads/p', { board: this.boardKey, thread: ref }]);
  }

  openImport(ev: Event, threadKey: string, content: any) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    if (!this.isRoot) {
      // this.common.showErrorAlert('Only Master Nodes Can Import', 3000);
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
        // this.common.showErrorAlert('None are suitable');
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
              // this.common.showAlert('successfully', 'success', 3000);
              this.initThreads(this.boardKey);
            });
          }
        }
      }, err => {
      });
    });
  }
}
