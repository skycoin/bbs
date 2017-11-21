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
  FollowPageData,
  UserService,
  Votes
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
  threadPage: ThreadPage;
  postForm = new FormGroup({
    name: new FormControl('', Validators.required),
    body: new FormControl('', Validators.required),
  });
  showUserInfoMenu = false;
  userTag = '';
  threadVoteStyle = '';
  postVoteStyle = '';
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
    private loading: LoadingService,
    private user: UserService) {
  }

  ngOnInit() {
    this.route.queryParams.subscribe(res => {
      this.boardKey = res['boardKey'];
      this.threadKey = res['thread_ref'];
      this.getThreadPage(this.boardKey, this.threadKey);
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

  countVote(votes: Votes) {
    if (votes) {
      const num = votes.up_votes.count - votes.down_votes.count;
      if (num < 1000) {
        return num;
      }
      return (num / 1000) + 'k'
    }
    return null;
  }

  addThreadVote(mode: number, ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    // Init Thread Votes
    if (!this.threadPage.data.thread.votes) {
      this.threadPage.data.thread.votes = {
        up_votes: { count: 0, voted: false },
        down_votes: { count: 0, voted: false }
      }
    }
    if (mode > 0) {
      this.threadVoteStyle = 'up-vote';
      this.threadPage.data.thread.votes.up_votes.count += 1;
    } else {
      this.threadVoteStyle = 'down-vote';
      this.threadPage.data.thread.votes.down_votes.count += 1;
    }
    const jsonStr = {
      type: `${this.api.version},thread_vote`,
      ts: new Date().getTime() * 1000000,
      of_board: this.boardKey,
      of_thread: this.threadKey,
      value: mode,
      creator: this.user.loginInfo.PublicKey,
    };

    this.api.submit(JSON.stringify(jsonStr), this.user.loginInfo.SecKey).subscribe(voteRes => {
      if (voteRes.okay) {
        this.threadPage.data.thread.votes = voteRes.data.votes;
      }
    }, err => {
      this.threadVoteStyle = '';
      if (mode > 0) {
        this.threadPage.data.thread.votes.up_votes.count += 1;
      } else {
        this.threadPage.data.thread.votes.down_votes.count += 1;
      }
    })
  }
  addUserVote(ev: Event, ofUser: string, mode: string) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.loading.start();
    const jsonStr = {
      type: `${this.api.version},user_vote`,
      ts: new Date().getTime() * 1000000,
      of_board: this.boardKey,
      of_user: ofUser,
      value: mode,
      creator: this.user.loginInfo.PublicKey,
    }
    this.api.submit(JSON.stringify(jsonStr), this.user.loginInfo.SecKey).subscribe(result => {
      if (result.okay) {
        this.userTag = '';
      }
      this.loading.close();
    }, err => {
      this.loading.close();
    })
  }

  addPostVote(mode: number, post: Post, ev: Event) {
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
    console.log('test post:', post);
    if (mode > 0) {
      this.postVoteStyle = 'up-vote';
      post.votes.up_votes.count += 1;
    } else {
      this.postVoteStyle = 'down-vote';
      post.votes.down_votes.count += 1;
    }
    const jsonStr = {
      type: `${this.api.version},post_vote`,
      ts: new Date().getTime() * 1000000,
      of_board: this.boardKey,
      of_thread: post.body.of_thread,
      of_post: post.header.hash,
      value: mode,
      creator: this.user.loginInfo.PublicKey,
    }
    this.api.submit(JSON.stringify(jsonStr), this.user.loginInfo.SecKey).subscribe(res => {
      console.log('add post vote: ', res);
      if (res.okay) {
        post.votes = res.data.votes;
      }
    }, err => {
      this.postVoteStyle = '';
      if (mode > 0) {
        post.votes.up_votes.count += 1;
      } else {
        post.votes.down_votes.count += 1;
      }
    })
  }

  openReply(content) {
    if (this.user.loginInfo) {
      this.postForm.reset();
      this.pop.open(content).result.then((result) => {
        if (result) {
          if (!this.postForm.valid) {
            this.alert.error({ content: 'title and content can not be empty' });
            return;
          }
          const jsonStr = {
            type: `${this.api.version},post`,
            name: this.postForm.get('name').value,
            body: this.common.replaceHtmlEnter(this.postForm.get('body').value),
            creator: this.user.loginInfo.PublicKey,
            ts: new Date().getTime() * 1000000,
            of_board: this.boardKey,
            of_thread: this.threadKey
          };
          this.loading.start();
          this.api.submit(JSON.stringify(jsonStr), this.user.loginInfo.SecKey).subscribe((res: ThreadPage) => {
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
  getThreadPage(boardKey, ref: string) {
    if (boardKey === '' || ref === '') {
      this.alert.error({ content: 'Parameter error!!!' });
      return;
    }
    const data = new FormData();
    data.append('board_public_key', boardKey);
    data.append('thread_ref', ref);
    this.api.getThreadpage(data).subscribe(res => {
      this.threadPage = res;
    }, err => {
      // this.router.navigate(['']);
    });
  }

}
