import { Component, EventEmitter, HostBinding, OnInit, Output, ViewEncapsulation, ViewChild, TemplateRef } from '@angular/core';
import { ApiService, Board, CommonService, Thread, Alert, BoardPage, Popup } from '../../providers';
import { ActivatedRoute, Router, ParamMap } from '@angular/router';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { slideInLeftAnimation } from '../../animations/router.animations';
import { flyInOutAnimation, bounceInAnimation } from '../../animations/common.animations';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/observable/timer'
import 'rxjs/add/operator/switchMap';

@Component({
  selector: 'app-threads',
  templateUrl: 'threads.component.html',
  styleUrls: ['threads.component.scss'],
  encapsulation: ViewEncapsulation.None,
  animations: [slideInLeftAnimation, flyInOutAnimation, bounceInAnimation],
})

export class ThreadsComponent implements OnInit {
  @HostBinding('@routeAnimation') routeAnimation = true;
  @HostBinding('style.display') display = 'block';
  @ViewChild('fab') fabBtnTemplate: TemplateRef<any>;
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
    name: new FormControl('', Validators.required),
  });
  @Output() thread: EventEmitter<{ master: string, ref: string }> = new EventEmitter();

  constructor(private api: ApiService,
    private router: Router,
    private route: ActivatedRoute,
    private common: CommonService,
    private alert: Alert,
    private pop: Popup) {
  }

  ngOnInit() {
    this.route.queryParams.subscribe(params => {
      this.boardKey = params['boardKey'];
      this.init();
    })
    Observable.timer(10).subscribe(() => {
      this.pop.open(this.fabBtnTemplate, { isDialog: false });
    });
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
    this.pop.open(content);
  }

  openAdd(content) {
    this.api.getSessionInfo().subscribe(info => {
      if (info.data.logged_in) {
        this.addForm.reset();
        this.pop.open(content).result.then((result) => {
          if (result) {
            if (!this.addForm.valid) {
              this.alert.error({ content: 'Parameter error!!!' });
              return;
            }
            const data = new FormData();
            data.append('board_public_key', this.boardKey);
            data.append('body', this.common.replaceHtmlEnter(this.addForm.get('body').value));
            data.append('name', this.addForm.get('name').value);
            this.api.newThread(data).subscribe(threadRes => {
              this.threads = threadRes.data.threads;
              this.alert.success({ content: 'Added successfully' });
            });
          }
        }, err => {
        });
      } else {
        this.alert.warning({ content: 'Please Login' });
      }
    })

  }

  open(ref: string) {
    if (this.boardKey === '' || ref === '') {
      // this.common.showErrorAlert('Parameter error!!!');
      return;
    }
    this.router.navigate(['/threads/p'], { queryParams: { boardKey: this.boardKey, thread_ref: ref } });
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
      this.pop.open(content).result.then(result => {
        if (result) {
          if (this.importBoardKey) {
            const data = new FormData();
            data.append('from_board', this.boardKey);
            data.append('thread', threadKey);
            data.append('to_board', this.importBoardKey);
            this.api.importThread(data).subscribe(res => {
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
