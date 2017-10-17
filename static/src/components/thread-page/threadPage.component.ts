import {
  Component,
  HostBinding,
  HostListener,
  OnInit,
  ViewEncapsulation,
  ViewChild,
  AfterViewInit,
  TemplateRef,
  ElementRef,
  ViewChildren,
  QueryList
} from '@angular/core';
import {
  ApiService,
  CommonService,
  ThreadPage,
  Post,
  VotesSummary,
  Thread,
  Alert,
  Popup,
  LoadingService,
  FollowPage,
  FollowPageData
} from '../../providers';
import { ActivatedRoute, Router } from '@angular/router';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { slideInLeftAnimation } from '../../animations/router.animations';
import { flyInOutAnimation, bounceInAnimation } from '../../animations/common.animations';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/observable/timer'
import 'rxjs/add/operator/filter';

@Component({
  selector: 'app-threadpage',
  templateUrl: 'threadPage.component.html',
  styleUrls: ['threadPage.component.scss'],
  encapsulation: ViewEncapsulation.None,
  animations: [slideInLeftAnimation, flyInOutAnimation, bounceInAnimation],
})

export class ThreadPageComponent implements OnInit {
  @HostBinding('@routeAnimation') routeAnimation = true;
  @HostBinding('style.display') display = 'block';
  @ViewChild('editor') editor: ElementRef;
  @ViewChild('fab') fab: TemplateRef<any>;
  @ViewChildren('post') posts: QueryList<ElementRef>;
  sort = 'esc';
  boardKey = '';
  threadKey = '';
  threadPk = '';
  threadPage: ThreadPage;
  postForm = new FormGroup({
    name: new FormControl('', Validators.required),
    body: new FormControl('', Validators.required),
  });
  showUserInfoMenu = false;
  userTag = '';
  editorOptions = {
    placeholderText: 'Edit Your Content Here!',
    quickInsertButtons: ['table', 'ul', 'ol', 'hr'],
    toolbarButtons: [
      'bold',
      'italic',
      'underline',
      'strikeThrough',
      'subscript',
      'superscript',
      '|',
      'fontFamily',
      'fontSize',
      'color',
      'inlineStyle',
      'paragraphStyle',
      '|',
      'paragraphFormat',
      'align',
      'formatOL',
      'formatUL',
      'outdent',
      'indent',
      'quote',
      '-',
      'emoticons',
      'insertLink',
      '|',
      'insertHR',
      'selectAll',
      'clearFormatting',
      '|',
      'print',
      'spellChecker',
      'help',
      'html',
      '|',
      'undo',
      'redo'],
    heightMin: 200,
    events: {},
  };

  constructor(private api: ApiService,
    private router: Router,
    private route: ActivatedRoute,
    private common: CommonService,
    private alert: Alert,
    private pop: Popup,
    private loading: LoadingService) {
  }

  ngOnInit() {
    this.route.queryParams.subscribe(res => {
      this.boardKey = res['boardKey'];
      this.threadKey = res['thread_ref'];
      this.threadPk = res['thread_pk'];
      this.open(this.boardKey, this.threadKey);
    });
    Observable.timer(10).subscribe(() => {
      this.pop.open(this.fab, { isDialog: false });
    });
  }
  showUserMenu(post: Post, ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    if (post.creatorMenu) {
      post.creatorMenu = false;
      return;
    }
    post.creatorMenu = true;
    if (!post.body.creator) {
      return;
    }
    this.showUserInfoMenu = true;
  }
  Menu(ev: Event, post: Post) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    if (!post.voteMenu) {
      post.voteMenu = true;
    } else {
      post.voteMenu = false;
    }
  }
  public setSort() {
    this.sort = this.sort === 'desc' ? 'asc' : 'desc';
  }
  trackPosts(index, post: Post) {
    return post ? post.header.hash : undefined;
  }
  addThreadVote(mode: string, ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    if (!this.threadPage.data.thread.votes) {
      this.threadPage.data.thread.votes = {
        up_votes: { count: 0, voted: false },
        down_votes: { count: 0, voted: false }
      }
    }
    if (mode === '-1') {
      this.threadPage.data.thread.votes.down_votes.count += 1;
    } else {
      this.threadPage.data.thread.votes.up_votes.count += 1;
    }
    // const jsonStr = {
    //   creator: ApiService.userInfo.public_key,
    //   of_board: this.boardKey,
    //   of_thread: this.threadKey,
    //   mode: mode
    // }
    const data = new FormData();
    data.append('board_public_key', this.boardKey);
    data.append('thread_ref', this.threadKey);
    data.append('mode', mode);
    this.api.addOldThreadVote(data).subscribe(voteRes => {
      if (voteRes.okay) {
        this.threadPage.data.thread.votes = voteRes.data.votes;
      }
    }, err => {
      if (mode === '-1') {
        this.threadPage.data.thread.votes.down_votes.count -= 1;
      } else {
        this.threadPage.data.thread.votes.up_votes.count -= 1;
      }
    })
  }
  addUserVote(ev: Event, post: Post, mode: string) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    const data = new FormData();
    data.append('board_public_key', this.boardKey);
    data.append('user_public_key', post.body.creator);
    data.append('mode', mode);
    this.loading.start();
    this.api.addUserVote(data).subscribe(result => {
      if (result.okay) {
        // this.userFollow = result.data;
        this.userTag = '';
      }
      this.loading.close();
    }, err => {
      this.loading.close();
    })
    post.creatorMenu = false;
  }
  addPostVote(mode: string, post: Post, ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    post.voteMenu = false;
    if (post.uiOptions !== undefined && post.uiOptions.voted !== undefined && post.uiOptions.voted) {
      return;
    }
    if (!post.votes) {
      post.votes = {
        up_votes: { count: 0, voted: false },
        down_votes: { count: 0, voted: false }
      }
    }
    if (mode === '-1') {
      post.votes.down_votes.count += 1;
    } else {
      post.votes.up_votes.count += 1;
    }
    const jsonStr = {
      creator: ApiService.userInfo.public_key,
      of_board: this.boardKey,
      of_thread: post.body.of_thread,
      mode: mode
    }
    const data = new FormData();
    data.append('board_public_key', this.boardKey);
    data.append('post_ref', post.body.of_thread);
    data.append('mode', mode);
    this.api.addOldPostVote(data).subscribe(res => {
      console.log('add post vote: ', res);
      if (res.okay) {
        post.votes = res.data.votes;
      }
    }, err => {
      if (mode === '-1') {
        post.votes.down_votes.count -= 1;
      } else {
        post.votes.up_votes.count -= 1;
      }
    })
  }
  openReply(content) {
    this.api.getSessionInfo().subscribe(info => {
      if (info.data.logged_in) {
        this.postForm.reset();
        this.pop.open(content).result.then((result) => {
          if (result) {
            if (!this.postForm.valid) {
              this.alert.error({ content: 'title and content can not be empty' });
              return;
            }
            const jsonStr = {
              name: this.postForm.get('name').value,
              body: this.common.replaceHtmlEnter(this.postForm.get('body').value),
              creator: ApiService.userInfo.public_key,
              of_board: this.boardKey,
              of_thread: this.threadKey
            };
            this.loading.start();
            this.api.newPost(JSON.stringify(jsonStr)).subscribe((res: ThreadPage) => {
              console.log('new post:', res);
              if (res.okay) {
                this.threadPage.data.posts = res.data.posts;
                this.alert.success({ content: 'Added successfully' });
                this.loading.close();
              }
            });
          }
        }, err => {
        });
      } else {
        this.alert.warning({ content: 'Please Login at first' });
      }
    });


  }
  PostAuthorMenu(post: Post, ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    if (!post.uiOptions) {
      post.uiOptions = { menu: true };
    } else {
      post.uiOptions.menu = !post.uiOptions.menu;
    }
  }
  open(boardKey, ref: string) {
    if (boardKey === '' || ref === '') {
      this.alert.error({ content: 'Parameter error!!!' });
      return;
    }
    const data = new FormData();
    data.append('board_public_key', boardKey);
    data.append('thread_ref', ref);
    this.api.getThreadpage(data).subscribe(res => {
      console.log('res: ', res);
      this.threadPage = res;
    }, err => {
      // this.router.navigate(['']);
    });
  }

}
